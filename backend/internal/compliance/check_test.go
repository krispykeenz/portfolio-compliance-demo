package compliance

import (
	"errors"
	"strings"
	"testing"

	"portfolio-compliance-demo/backend/internal/domain"
)

func demoPortfolio() domain.Portfolio {
	return domain.Portfolio{ID: "portfolio-demo-001", NAVMinor: 100_000_000, Holdings: []domain.Holding{{Instrument: domain.Instrument{ID: "alpha", Name: "Alpha Global Equity Fund"}, ValueMinor: 9_000_000}}}
}

func TestCheckTradeBoundaries(t *testing.T) {
	tests := []struct {
		name       string
		amount     int64
		wantPass   bool
		wantReason string
	}{
		{"below", 999_900, true, "WITHIN_LIMIT"},
		{"exact", 1_000_000, true, "WITHIN_LIMIT"},
		{"above", 1_000_100, false, "LIMIT_EXCEEDED"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Check(demoPortfolio(), domain.Trade{PortfolioID: "portfolio-demo-001", InstrumentID: "alpha", Side: domain.Buy, AmountMinor: tt.amount}, 1000)
			if err != nil {
				t.Fatal(err)
			}
			if got.Passed != tt.wantPass || got.ReasonCode != tt.wantReason {
				t.Fatalf("got pass=%v reason=%s", got.Passed, got.ReasonCode)
			}
			if got.Passed != strings.HasPrefix(got.Explanation, "PASS:") {
				t.Fatalf("explanation disagrees: %q", got.Explanation)
			}
		})
	}
}

func TestCheckTradeValidation(t *testing.T) {
	tests := []struct {
		name, instrument string
		side             domain.TradeSide
		amount           int64
		want             error
	}{
		{"zero", "alpha", domain.Buy, 0, domain.ErrInvalidTrade},
		{"negative", "alpha", domain.Buy, -1, domain.ErrInvalidTrade},
		{"unknown", "missing", domain.Buy, 100, domain.ErrInstrumentNotFound},
		{"sell underflow", "alpha", domain.Sell, 9_000_001, domain.ErrInvalidTrade},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Check(demoPortfolio(), domain.Trade{InstrumentID: tt.instrument, Side: tt.side, AmountMinor: tt.amount}, 1000)
			if !errors.Is(err, tt.want) {
				t.Fatalf("got %v want %v", err, tt.want)
			}
		})
	}
}

func TestFormatPercent(t *testing.T) {
	for _, tt := range []struct {
		value, nav int64
		want       string
	}{{9_000_000, 100_000_000, "9.00%"}, {10_000_000, 100_000_000, "10.00%"}, {10_000_100, 100_000_000, "10.0001%"}} {
		if got := FormatPercent(tt.value, tt.nav); got != tt.want {
			t.Fatalf("got %s want %s", got, tt.want)
		}
	}
}
