package transport

import (
	"context"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	compliancev1 "portfolio-compliance-demo/backend/generated/compliance/v1"
	"portfolio-compliance-demo/backend/internal/repository"
)

func clientForTest(t *testing.T) compliancev1.ComplianceServiceClient {
	t.Helper()
	lis := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	compliancev1.RegisterComplianceServiceServer(server, NewService(repository.NewSeeded()))
	go func() { _ = server.Serve(lis) }()
	t.Cleanup(server.Stop)
	conn, err := grpc.NewClient("passthrough:///bufnet", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return lis.Dial() }), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	return compliancev1.NewComplianceServiceClient(conn)
}

func TestServiceIntegration(t *testing.T) {
	client := clientForTest(t)
	portfolio, err := client.GetPortfolio(context.Background(), &compliancev1.GetPortfolioRequest{PortfolioId: "portfolio-demo-001"})
	if err != nil {
		t.Fatal(err)
	}
	if got := len(portfolio.GetPortfolio().GetHoldings()); got != 3 {
		t.Fatalf("got %d holdings", got)
	}

	response, err := client.CheckTrade(context.Background(), &compliancev1.CheckTradeRequest{Trade: &compliancev1.ProposedTrade{PortfolioId: "portfolio-demo-001", InstrumentId: "alpha", Side: compliancev1.TradeSide_TRADE_SIDE_BUY, Amount: &compliancev1.Money{MinorUnits: 1_000_000}}})
	if err != nil {
		t.Fatal(err)
	}
	if response.GetResult().GetStatus() != compliancev1.ComplianceStatus_COMPLIANCE_STATUS_PASS {
		t.Fatalf("got %s", response.GetResult().GetStatus())
	}
}

func TestServiceStatusCodes(t *testing.T) {
	client := clientForTest(t)
	for _, tt := range []struct {
		name    string
		request *compliancev1.CheckTradeRequest
		want    codes.Code
	}{
		{"malformed", &compliancev1.CheckTradeRequest{}, codes.InvalidArgument},
		{"unknown portfolio", &compliancev1.CheckTradeRequest{Trade: &compliancev1.ProposedTrade{PortfolioId: "missing", InstrumentId: "alpha", Side: compliancev1.TradeSide_TRADE_SIDE_BUY, Amount: &compliancev1.Money{MinorUnits: 100}}}, codes.NotFound},
		{"unknown instrument", &compliancev1.CheckTradeRequest{Trade: &compliancev1.ProposedTrade{PortfolioId: "portfolio-demo-001", InstrumentId: "missing", Side: compliancev1.TradeSide_TRADE_SIDE_BUY, Amount: &compliancev1.Money{MinorUnits: 100}}}, codes.NotFound},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.CheckTrade(context.Background(), tt.request)
			if status.Code(err) != tt.want {
				t.Fatalf("got %v", err)
			}
		})
	}
}
