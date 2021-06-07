package balance

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/sf9v/kalupi/account"
)

type Service interface {
	GetAccntBal(context.Context, account.AccountID) (*account.Balance, error)
}

type service struct {
	balRepo Repository
}

var _ Service = (*service)(nil)

func NewService(balRepo Repository) Service {
	return &service{balRepo: balRepo}
}

func (s *service) GetAccntBal(ctx context.Context, accntID account.AccountID) (*account.Balance, error) {
	tx, err := s.balRepo.BeginTx(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "begin tx")
	}

	accntBal, err := s.balRepo.GetAccntBal(ctx, tx, accntID)
	if err != nil {
		err = errors.Wrap(err, "get accnt bal")
		txErr := tx.Rollback()
		return nil, multierr.Combine(err, txErr)
	}

	err = tx.Commit()
	if err != nil {
		return nil, errors.Wrap(err, "tx commit")
	}

	return accntBal, nil
}
