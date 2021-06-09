package service

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/sf9v/kalupi/account"
	"github.com/sf9v/kalupi/balance"
)

// service is an account service implementation
type service struct {
	accountRepo account.Repository
	balService  balance.Service
}

var _ account.Service = (*service)(nil)

// New takes an account and balance repository and returns an account service
func New(accountRepo account.Repository, balService balance.Service) account.Service {
	return &service{accountRepo: accountRepo, balService: balService}
}

// CreateAccount creates a new account
func (s *service) CreateAccount(ctx context.Context, accnt account.Account) error {
	err := accnt.Validate()
	if err != nil {
		return multierr.Combine(account.ErrValidation, err)
	}

	exists, err := s.accountRepo.IsAccountExists(ctx, accnt.AccountID)
	if err != nil {
		return errors.Wrap(err, "is account exists")
	}

	if exists {
		return account.ErrAccountAlreadyExists
	}

	_, err = s.accountRepo.CreateAccount(ctx, accnt)
	if err != nil {
		return errors.Wrap(err, "repo create account")
	}

	return nil
}

// GetAccount retrievies an account via account id
func (s *service) GetAccount(ctx context.Context,
	accntID account.AccountID) (*account.Account, error) {
	err := accntID.Validate()
	if err != nil {
		return nil, account.ErrValidation
	}

	exists, err := s.accountRepo.IsAccountExists(ctx, accntID)
	if err != nil {
		return nil, errors.Wrap(err, "is account exists")
	}

	if !exists {
		return nil, account.ErrAccountNotFound
	}

	accnt, err := s.accountRepo.GetAccount(ctx, accntID)
	if err != nil {
		return nil, errors.Wrap(err, "repo get account")
	}

	bal, err := s.balService.GetAccntBal(ctx, accnt.AccountID)
	if err != nil {
		return nil, errors.Wrap(err, "get account balance")
	}

	accnt.Balance = bal.CurrentBal

	return accnt, nil
}

func (s *service) ListAccounts(ctx context.Context) ([]*account.Account, error) {
	accnts, err := s.accountRepo.ListAccounts(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "repo list accounts")
	}

	for _, accnt := range accnts {
		bal, err := s.balService.GetAccntBal(ctx, accnt.AccountID)
		if err != nil {
			return nil, errors.Wrap(err, "get account balance")
		}
		accnt.Balance = bal.CurrentBal
	}

	return accnts, nil
}
