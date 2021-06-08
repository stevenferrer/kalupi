package account

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/sf9v/kalupi/currency"
)

type createAccountRequest struct {
	AccountID AccountID         `json:"account_id"`
	Currency  currency.Currency `json:"currency"`
}

type createAccountResponse struct {
	Err error `json:"error,omitempty"`
}

func (r createAccountResponse) error() error { return r.Err }

func newCreateAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createAccountRequest)
		ac := Account{
			AccountID: req.AccountID,
			Currency:  req.Currency,
		}

		err := s.CreateAccount(ctx, ac)
		return createAccountResponse{Err: err}, nil
	}
}

type getAccountRequest struct {
	AccountID AccountID
}

type getAccountResponse struct {
	Account *Account `json:"account,omitempty"`
	Err     error    `json:"error,omitempty"`
}

func (r getAccountResponse) error() error { return r.Err }

func newGetAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getAccountRequest)
		accnt, err := s.GetAccount(ctx, req.AccountID)
		return getAccountResponse{Account: accnt, Err: err}, nil
	}
}

type listAccountsRequest struct{}

type listAccountsResponse struct {
	Accounts []*Account `json:"accounts,omitempty"`
	Err      error      `json:"error,omitempty"`
}

func (r listAccountsResponse) error() error { return r.Err }

func newListAccountsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		accnts, err := s.ListAccounts(ctx)
		return listAccountsResponse{Accounts: accnts, Err: err}, nil
	}
}
