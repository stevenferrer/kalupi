package transaction

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/sf9v/kalupi/account"
	"github.com/shopspring/decimal"
)

type depositRequest struct {
	AccountID account.AccountID `json:"account_id"`
	Amount    decimal.Decimal   `json:"amount"`
}

type depositResponse struct {
	Err error `json:"error,omitempty"`
}

func (r depositResponse) error() error { return r.Err }

func newDepositEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(depositRequest)
		err := s.MakeDeposit(ctx, DepositXact(req))
		return depositResponse{Err: err}, nil
	}
}

type withdrawalRequest struct {
	AccountID account.AccountID `json:"account_id"`
	Amount    decimal.Decimal   `json:"amount"`
}

type withdrawalResponse struct {
	Err error `json:"error,omitempty"`
}

func (r withdrawalResponse) error() error { return r.Err }

func newWithdrawalEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(withdrawalRequest)
		err := s.MakeWithdrawal(ctx, WithdrawalXact(req))
		return withdrawalResponse{Err: err}, nil
	}
}

type paymentRequest struct {
	FromAccount account.AccountID `json:"from_account"`
	ToAccount   account.AccountID `json:"to_account"`
	Amount      decimal.Decimal   `json:"amount"`
}

type paymentResponse struct {
	Err error `json:"error,omitempty"`
}

func (r paymentResponse) error() error { return r.Err }

func newPaymentEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(paymentRequest)
		err := s.MakeTransfer(ctx, TransferXact(req))
		return paymentResponse{Err: err}, nil
	}
}

type listPaymentsRequest struct{}

type listPaymentsResponse struct {
	Payments []*Payment `json:"payments"`
	Err      error      `json:"error,omitempty"`
}

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
