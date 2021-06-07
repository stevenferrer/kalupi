package ledger

import (
	"context"
)

// Service is a ledger service
type Service interface {
	// CreateCashLedgers will create the internal cash ledger accounts (usd, eur, etc.) if not exist
	CreateCashLedgers(context.Context) error
}

type service struct {
	ledgerRepo Repository
}

var _ Service = (*service)(nil)

func NewService(ledgerRepo Repository) Service {
	return &service{ledgerRepo: ledgerRepo}
}

func (s *service) CreateCashLedgers(ctx context.Context) error {
	return s.ledgerRepo.CreateLedgersIfNotExist(ctx, cashLedgers[:]...)
}
