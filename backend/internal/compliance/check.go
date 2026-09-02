package compliance

import (
	"fmt"
	"math/big"
	"strings"

	"portfolio-compliance-demo/backend/internal/domain"
)

func Check(portfolio domain.Portfolio, trade domain.Trade, limitBasisPoints int32) (domain.Result, error) {
	if trade.AmountMinor <= 0 || (trade.Side != domain.Buy && trade.Side != domain.Sell) || portfolio.NAVMinor <= 0 {
		return domain.Result{}, domain.ErrInvalidTrade
	}
	var holding *domain.Holding
	for i := range portfolio.Holdings {
		if portfolio.Holdings[i].Instrument.ID == trade.InstrumentID {
			holding = &portfolio.Holdings[i]
			break
		}
	}
	if holding == nil {
		return domain.Result{}, domain.ErrInstrumentNotFound
	}

	projected := holding.ValueMinor
	if trade.Side == domain.Buy {
		projected += trade.AmountMinor
	} else {
		projected -= trade.AmountMinor
	}
	if projected < 0 {
		return domain.Result{}, domain.ErrInvalidTrade
	}

	left := new(big.Int).Mul(big.NewInt(projected), big.NewInt(10_000))
	right := new(big.Int).Mul(big.NewInt(portfolio.NAVMinor), big.NewInt(int64(limitBasisPoints)))
	passed := left.Cmp(right) <= 0
	reason, prefix, final := "WITHIN_LIMIT", "PASS", "The maximum permitted exposure is"
	if !passed {
		reason, prefix, final = "LIMIT_EXCEEDED", "FAIL", "This exceeds the"
	}
	currentPercent, projectedPercent := FormatPercent(holding.ValueMinor, portfolio.NAVMinor), FormatPercent(projected, portfolio.NAVMinor)
	limitPercent := formatBasisPoints(limitBasisPoints)
	explanation := fmt.Sprintf("%s: %s would move from %s to %s.\n%s %s limit.", prefix, holding.Instrument.Name, currentPercent, projectedPercent, final, limitPercent)
	if passed {
		explanation = fmt.Sprintf("%s: %s would move from %s to %s.\nThe maximum permitted exposure is %s.", prefix, holding.Instrument.Name, currentPercent, projectedPercent, limitPercent)
	}
	return domain.Result{Passed: passed, ReasonCode: reason, CurrentValueMinor: holding.ValueMinor, ProjectedValueMinor: projected, NAVMinor: portfolio.NAVMinor, LimitBasisPoints: limitBasisPoints, InstrumentName: holding.Instrument.Name, Explanation: explanation}, nil
}

func FormatPercent(value, nav int64) string {
	if nav <= 0 {
		return "0.00%"
	}
	scaled := new(big.Int).Quo(new(big.Int).Mul(big.NewInt(value), big.NewInt(100_000_000)), big.NewInt(nav)).Int64()
	whole, fraction := scaled/1_000_000, scaled%1_000_000
	decimal := fmt.Sprintf("%06d", fraction)
	decimal = strings.TrimRight(decimal, "0")
	if len(decimal) < 2 {
		decimal += strings.Repeat("0", 2-len(decimal))
	}
	return fmt.Sprintf("%d.%s%%", whole, decimal)
}

func formatBasisPoints(bp int32) string { return fmt.Sprintf("%d.%02d%%", bp/100, bp%100) }
