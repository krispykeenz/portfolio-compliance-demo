package domain

import "errors"

var (
	ErrInvalidTrade       = errors.New("invalid trade")
	ErrPortfolioNotFound  = errors.New("portfolio not found")
	ErrInstrumentNotFound = errors.New("instrument not found")
)

type TradeSide int

const (
	Buy TradeSide = iota + 1
	Sell
)

type Instrument struct{ ID, Name string }
type Holding struct {
	Instrument Instrument
	ValueMinor int64
}
type Portfolio struct {
	ID, Name             string
	NAVMinor             int64
	Holdings             []Holding
	RemainderDescription string
}
type Trade struct {
	PortfolioID, InstrumentID string
	Side                      TradeSide
	AmountMinor               int64
}
type Result struct {
	Passed                                           bool
	ReasonCode                                       string
	CurrentValueMinor, ProjectedValueMinor, NAVMinor int64
	LimitBasisPoints                                 int32
	InstrumentName, Explanation                      string
}
