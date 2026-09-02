package transport

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	compliancev1 "portfolio-compliance-demo/backend/generated/compliance/v1"
	"portfolio-compliance-demo/backend/internal/compliance"
	"portfolio-compliance-demo/backend/internal/domain"
	"portfolio-compliance-demo/backend/internal/repository"
)

type Service struct {
	compliancev1.UnimplementedComplianceServiceServer
	repo *repository.Seeded
}

func NewService(repo *repository.Seeded) *Service { return &Service{repo: repo} }

func (s *Service) GetPortfolio(ctx context.Context, req *compliancev1.GetPortfolioRequest) (*compliancev1.GetPortfolioResponse, error) {
	if req == nil || req.GetPortfolioId() == "" {
		return nil, status.Error(codes.InvalidArgument, "portfolio_id is required")
	}
	p, err := s.repo.Get(ctx, req.GetPortfolioId())
	if err != nil {
		return nil, grpcError(err)
	}
	return &compliancev1.GetPortfolioResponse{Portfolio: portfolioToProto(p)}, nil
}

func (s *Service) CheckTrade(ctx context.Context, req *compliancev1.CheckTradeRequest) (*compliancev1.CheckTradeResponse, error) {
	if req == nil || req.GetTrade() == nil || req.GetTrade().GetAmount() == nil {
		return nil, status.Error(codes.InvalidArgument, "trade and amount are required")
	}
	t := req.GetTrade()
	if t.GetPortfolioId() == "" || t.GetInstrumentId() == "" {
		return nil, status.Error(codes.InvalidArgument, "portfolio_id and instrument_id are required")
	}
	p, err := s.repo.Get(ctx, t.GetPortfolioId())
	if err != nil {
		return nil, grpcError(err)
	}
	side := domain.TradeSide(0)
	if t.GetSide() == compliancev1.TradeSide_TRADE_SIDE_BUY {
		side = domain.Buy
	}
	if t.GetSide() == compliancev1.TradeSide_TRADE_SIDE_SELL {
		side = domain.Sell
	}
	result, err := compliance.Check(p, domain.Trade{PortfolioID: t.GetPortfolioId(), InstrumentID: t.GetInstrumentId(), Side: side, AmountMinor: t.GetAmount().GetMinorUnits()}, 1000)
	if err != nil {
		return nil, grpcError(err)
	}
	statusValue, reason := compliancev1.ComplianceStatus_COMPLIANCE_STATUS_FAIL, compliancev1.ReasonCode_REASON_CODE_LIMIT_EXCEEDED
	if result.Passed {
		statusValue, reason = compliancev1.ComplianceStatus_COMPLIANCE_STATUS_PASS, compliancev1.ReasonCode_REASON_CODE_WITHIN_LIMIT
	}
	return &compliancev1.CheckTradeResponse{Result: &compliancev1.RuleResult{
		RuleId: "MAX_INSTRUMENT_EXPOSURE", Status: statusValue, ReasonCode: reason,
		CurrentExposure: exposure(result.CurrentValueMinor, result.NAVMinor), ProjectedExposure: exposure(result.ProjectedValueMinor, result.NAVMinor),
		LimitBasisPoints: result.LimitBasisPoints, LimitFormattedPercent: "10.00%", Explanation: result.Explanation,
	}}, nil
}

func portfolioToProto(p domain.Portfolio) *compliancev1.Portfolio {
	holdings := make([]*compliancev1.Holding, 0, len(p.Holdings))
	for _, h := range p.Holdings {
		holdings = append(holdings, &compliancev1.Holding{Instrument: &compliancev1.Instrument{Id: h.Instrument.ID, Name: h.Instrument.Name}, CurrentValue: money(h.ValueMinor), CurrentExposure: exposure(h.ValueMinor, p.NAVMinor)})
	}
	return &compliancev1.Portfolio{Id: p.ID, Name: p.Name, TotalNav: money(p.NAVMinor), Holdings: holdings, RemainderDescription: p.RemainderDescription}
}

func money(minor int64) *compliancev1.Money {
	return &compliancev1.Money{MinorUnits: minor, Currency: "ZAR", Formatted: fmt.Sprintf("%d.%02d", minor/100, minor%100)}
}
func exposure(value, nav int64) *compliancev1.Exposure {
	return &compliancev1.Exposure{InstrumentValueMinorUnits: value, PortfolioNavMinorUnits: nav, FormattedPercent: compliance.FormatPercent(value, nav)}
}
func grpcError(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidTrade):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrPortfolioNotFound), errors.Is(err, domain.ErrInstrumentNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
