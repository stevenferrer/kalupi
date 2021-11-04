package balance

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/stevenferrer/kalupi/account"
	"github.com/stevenferrer/kalupi/etc/tx"
)

// Service is an account balance service
type Service interface {
	// GetAccntBal retrieives the account balance
	GetAccntBal(context.Context, account.AccountID) (*account.Balance, error)
}

// service an a balance service implementation
type service struct {
	balRepo Repository
}

var _ Service = (*service)(nil)

// NewService takes a balance repository and returns a balance service
func NewService(balRepo Repository) Service {
	return &service{balRepo: balRepo}
}

// GetAccntBal retreives the account balance
func (s *service) GetAccntBal(ctx context.Context, accntID account.AccountID) (accntBal *account.Balance, err error) {
	var tx tx.Tx
	tx, err = s.balRepo.BeginTx(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "begin tx")
	}
	defer func() {
		// rollback if there are errors
		if err != nil {
			_ = tx.Rollback()
			return
		}

		// commit if no errors
		if commitErr := tx.Commit(); commitErr != nil {
			err = multierr.Combine(err, commitErr)
		}
	}()

	accntBal, err = s.balRepo.GetAccntBal(ctx, tx, accntID)
	if err != nil {
		err = errors.Wrap(err, "get accnt bal")
		return
	}

	return accntBal, nil
}
