package transaction

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

// loggingService is a logging service middleware
type loggingService struct {
	logger log.Logger
	s      Service
}

// NewLoggingService returns a new logging service middleware
func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

// MakeDeposit logs the deposit params
func (s *loggingService) MakeDeposit(ctx context.Context, dp DepositXact) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "make_deposit",
			"account_id", dp.AccountID,
			"amount", dp.Amount,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return s.s.MakeDeposit(ctx, dp)
}

// MakeWithdrawal logs the withdrawal params
func (s *loggingService) MakeWithdrawal(ctx context.Context, wd WithdrawalXact) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "make_withdrawal",
			"account_id", wd.AccountID,
			"amount", wd.Amount,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return s.s.MakeWithdrawal(ctx, wd)
}

// MakeTransfer logs the transfer params
func (s *loggingService) MakeTransfer(ctx context.Context, tr TransferXact) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "make_transfer",
			"from_account", tr.FromAccount,
			"to_account", tr.ToAccount,
			"amount", tr.Amount,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return s.s.MakeTransfer(ctx, tr)
}

// ListTransfers logs the list transfers params
func (s *loggingService) ListTransfers(ctx context.Context) (xacts []*Transaction, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "list_transfers",
			"count", len(xacts),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return s.s.ListTransfers(ctx)

}
