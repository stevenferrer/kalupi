package account

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/sf9v/kalupi/currency"
)

// createAccountRequest is a create account request
type createAccountRequest struct {
	AccountID AccountID         `json:"account_id"`
	Currency  currency.Currency `json:"currency"`
}

// createAccountResponse is a create account response
type createAccountResponse struct {
	Err error `json:"error,omitempty"`
}

func (r createAccountResponse) error() error { return r.Err }

// newCreateAccountEndpoint returns a create account endpoint
func newCreateAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createAccountRequest)
		accnt := Account{
			AccountID: req.AccountID,
			Currency:  req.Currency,
		}

		err := s.CreateAccount(ctx, accnt)
		return createAccountResponse{Err: err}, nil
	}
}

// getAccountRequest is a get account request
type getAccountRequest struct {
	AccountID AccountID
}

// getAccountResponse is a get account response
type getAccountResponse struct {
	Account *Account `json:"account,omitempty"`
	Err     error    `json:"error,omitempty"`
}

func (r getAccountResponse) error() error { return r.Err }

// newGetAccountEndpoint returns a new account endpoint
func newGetAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getAccountRequest)
		accnt, err := s.GetAccount(ctx, req.AccountID)
		return getAccountResponse{Account: accnt, Err: err}, nil
	}
}

// listAccountsRequest is a list accounts request.
// Empty for now but can have other params such as limit.
type listAccountsRequest struct{}

// ListAccountsResponse is a list accounts response.
type listAccountsResponse struct {
	Accounts []*Account `json:"accounts,omitempty"`
	Err      error      `json:"error,omitempty"`
}

func (r listAccountsResponse) error() error { return r.Err }

// newListAccountsEndpoint returns a list accounts endpoint
func newListAccountsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		accnts, err := s.ListAccounts(ctx)
		return listAccountsResponse{Accounts: accnts, Err: err}, nil
	}
}
