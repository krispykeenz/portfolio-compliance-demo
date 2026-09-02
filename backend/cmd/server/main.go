package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	compliancev1 "portfolio-compliance-demo/backend/generated/compliance/v1"
	"portfolio-compliance-demo/backend/internal/repository"
	"portfolio-compliance-demo/backend/internal/transport"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(loggingInterceptor(logger)))
	compliancev1.RegisterComplianceServiceServer(grpcServer, transport.NewService(repository.NewSeeded()))

	grpcAddress, webAddress := env("GRPC_ADDRESS", ":9090"), env("WEB_ADDRESS", ":8080")
	grpcListener, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		logger.Error("listen native grpc", "error", err)
		os.Exit(1)
	}
	wrapped := grpcweb.WrapServer(grpcServer, grpcweb.WithOriginFunc(allowedOrigin))
	webServer := &http.Server{Addr: webAddress, Handler: originHeader(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/healthz" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"status":"ok"}`))
			return
		}
		if wrapped.IsGrpcWebRequest(r) || wrapped.IsAcceptableGrpcCorsRequest(r) {
			wrapped.ServeHTTP(w, r)
			return
		}
		http.NotFound(w, r)
	})), ReadHeaderTimeout: 5 * time.Second}

	go func() {
		logger.Info("native grpc listening", "address", grpcAddress)
		if err := grpcServer.Serve(grpcListener); err != nil {
			logger.Error("native grpc stopped", "error", err)
		}
	}()
	go func() {
		logger.Info("grpc-web listening", "address", webAddress)
		if err := webServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Error("grpc-web stopped", "error", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	grpcServer.GracefulStop()
	_ = webServer.Shutdown(shutdownCtx)
}

func allowedOrigin(origin string) bool {
	return origin == env("ALLOWED_ORIGIN", "http://localhost:5173")
}
func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
func originHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if secret := os.Getenv("ORIGIN_VERIFY_SECRET"); secret != "" && r.URL.Path != "/healthz" && r.Header.Get("X-Origin-Verify") != secret {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func requestID(ctx context.Context) string {
	if values := metadata.ValueFromIncomingContext(ctx, "x-request-id"); len(values) > 0 && values[0] != "" {
		return values[0]
	}
	b := make([]byte, 8)
	if _, err := rand.Read(b); err == nil {
		return hex.EncodeToString(b)
	}
	return "unavailable"
}
func loggingInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start, id := time.Now(), requestID(ctx)
		response, err := handler(ctx, req)
		logger.Info("rpc", "operation", strings.TrimPrefix(info.FullMethod, "/"), "request_id", id, "status", status.Code(err).String(), "duration_ms", time.Since(start).Milliseconds())
		return response, err
	}
}
