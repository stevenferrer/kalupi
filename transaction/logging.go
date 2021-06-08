package transaction

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

type loggingService struct {
	logger log.Logger
	s      Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) MakeDeposit(ctx context.Context, dp DepositXact) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "make_deposit",
			"account_id", dp.AccountID,
			"amount", dp.Amount,
			"err", err,
		)
	}(time.Now())

	return s.s.MakeDeposit(ctx, dp)
}

func (s *loggingService) MakeWithdrawal(ctx context.Context, wd WithdrawalXact) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "make_withdrawal",
			"account_id", wd.AccountID,
			"amount", wd.Amount,
			"err", err,
		)
	}(time.Now())

	return s.s.MakeWithdrawal(ctx, wd)
}

func (s *loggingService) MakeTransfer(ctx context.Context, tr TransferXact) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "make_transfer",
			"from_account", tr.FromAccount,
			"to_account", tr.ToAccount,
			"amount", tr.Amount,
			"err", err,
		)
	}(time.Now())

	return s.s.MakeTransfer(ctx, tr)
}

func (s *loggingService) ListTransfers(ctx context.Context) (xacts []*Transaction, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "list_transfers",
			"count", len(xacts),
			"err", err,
		)
	}(time.Now())

	return s.s.ListTransfers(ctx)

}
