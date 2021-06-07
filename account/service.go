package account

import (
	"context"

	"github.com/pkg/errors"
)

// Service is an account service
type Service interface {
	CreateAccount(context.Context, Account) (AccountID, error)
	GetAccount(context.Context, AccountID) (*Account, error)
	ListAccounts(context.Context) ([]*Account, error)
}

// service is an implementation of account service
type service struct {
	accountRepo Repository
}

var _ Service = (*service)(nil)

// NewService takes an account repository and returns an account service
func NewService(accountRepo Repository) Service {
	return &service{accountRepo: accountRepo}
}

func (s *service) CreateAccount(ctx context.Context,
	accnt Account) (AccountID, error) {
	id, err := s.accountRepo.CreateAccount(ctx, accnt)
	if err != nil {
		return "", errors.Wrap(err, "repo create account")
	}

	return id, nil
}

func (s *service) GetAccount(ctx context.Context,
	accntID AccountID) (*Account, error) {
	accnt, err := s.accountRepo.GetAccount(ctx, accntID)
	if err != nil {
		return nil, errors.Wrap(err, "repo get account")
	}

	return accnt, nil
}

func (s *service) ListAccounts(ctx context.Context) ([]*Account, error) {
	accnts, err := s.accountRepo.ListAccounts(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "repo list accounts")
	}

	return accnts, nil
}
