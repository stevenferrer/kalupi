package transaction

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/shopspring/decimal"
	"github.com/stevenferrer/kalupi/account"
)

// depositRequest is a deposit request
type depositRequest struct {
	AccountID account.AccountID `json:"account_id"`
	Amount    decimal.Decimal   `json:"amount"`
}

// depositResponse is a deposit response
type depositResponse struct {
	Err error `json:"error,omitempty"`
}

func (r depositResponse) error() error { return r.Err }

// newDeposit endpoint returns a deposit endpoint
func newDepositEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(depositRequest)
		err := s.MakeDeposit(ctx, DepositXact(req))
		return depositResponse{Err: err}, nil
	}
}

// withdrawalRequest is a withdrawal request
type withdrawalRequest struct {
	AccountID account.AccountID `json:"account_id"`
	Amount    decimal.Decimal   `json:"amount"`
}

// withdrawalResponse is a withdrawal response
type withdrawalResponse struct {
	Err error `json:"error,omitempty"`
}

func (r withdrawalResponse) error() error { return r.Err }

// newWithdrawalEndpoint returns a withdrawal endpoint
func newWithdrawalEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(withdrawalRequest)
		err := s.MakeWithdrawal(ctx, WithdrawalXact(req))
		return withdrawalResponse{Err: err}, nil
	}
}

// paymentRequest is a payment request
type paymentRequest struct {
	FromAccount account.AccountID `json:"from_account"`
	ToAccount   account.AccountID `json:"to_account"`
	Amount      decimal.Decimal   `json:"amount"`
}

// paymentResponse is a payment response
type paymentResponse struct {
	Err error `json:"error,omitempty"`
}

func (r paymentResponse) error() error { return r.Err }

// newPaymentEndpoint returns a payment endpoint
func newPaymentEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(paymentRequest)
		err := s.MakeTransfer(ctx, TransferXact(req))
		return paymentResponse{Err: err}, nil
	}
}

// listPaymentsRequest is list payments request.
// Empty for now but could contain other params such as limit.
type listPaymentsRequest struct{}

// listPaymentsResponse is a list payments response
type listPaymentsResponse struct {
	Payments []*Payment `json:"payments"`
	Err      error      `json:"error,omitempty"`
}

// newListPaymentsEndpoint returns a list payments endpoint
func newListPaymentsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		xacts, err := s.ListTransfers(ctx)
		if err != nil {
			return listPaymentsResponse{Err: err}, nil
		}
		payments := xactsToPayments(xacts)
		return listPaymentsResponse{Payments: payments}, nil
	}
}
