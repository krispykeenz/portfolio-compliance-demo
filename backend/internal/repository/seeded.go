package repository

import (
	"context"
	"portfolio-compliance-demo/backend/internal/domain"
)

type Seeded struct{ portfolio domain.Portfolio }

func NewSeeded() *Seeded {
	return &Seeded{portfolio: domain.Portfolio{
		ID: "portfolio-demo-001", Name: "Interview Demo Portfolio", NAVMinor: 100_000_000,
		Holdings: []domain.Holding{
			{Instrument: domain.Instrument{ID: "alpha", Name: "Alpha Global Equity Fund"}, ValueMinor: 9_000_000},
			{Instrument: domain.Instrument{ID: "beta", Name: "Beta Income Fund"}, ValueMinor: 7_500_000},
			{Instrument: domain.Instrument{ID: "gamma", Name: "Gamma Property Fund"}, ValueMinor: 5_000_000},
		},
		RemainderDescription: "Diversified assets outside the three holdings shown.",
	}}
}

func (s *Seeded) Get(ctx context.Context, id string) (domain.Portfolio, error) {
	if err := ctx.Err(); err != nil {
		return domain.Portfolio{}, err
	}
	if id != s.portfolio.ID {
		return domain.Portfolio{}, domain.ErrPortfolioNotFound
	}
	return s.portfolio, nil
}
