package ledger

import (
	"context"
)

// Service is a ledger service
type Service interface {
	// CreateCashLedgers creates the internal cash ledger accounts (USD, EUR, etc.)
	CreateCashLedgers(context.Context) error
}

// service is a ledger service implementation
type service struct {
	ledgerRepo Repository
}

var _ Service = (*service)(nil)

// NewService takes a ledger repository and returns a ledger service
func NewService(ledgerRepo Repository) Service {
	return &service{ledgerRepo: ledgerRepo}
}

// CreateCashLedgers creates the internal cash ledger accounts
func (s *service) CreateCashLedgers(ctx context.Context) error {
	return s.ledgerRepo.CreateLedgersIfNotExists(ctx, cashLedgers[:]...)
}
