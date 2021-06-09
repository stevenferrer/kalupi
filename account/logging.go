package account

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

// loggingService is a service logging middleware
type loggingService struct {
	logger log.Logger
	s      Service
}

// NewLoggingService returns a logging service middleware
func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger: logger, s: s}
}

// CreateAccount logs the create account params
func (s *loggingService) CreateAccount(ctx context.Context, accnt Account) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "create_account",
			"account_id", accnt.AccountID,
			"currency", accnt.Currency,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return s.s.CreateAccount(ctx, accnt)
}

// GetAccount logs the get account params
func (s *loggingService) GetAccount(ctx context.Context, accntID AccountID) (accnt *Account, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "get_account",
			"account_id", accntID,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return s.s.GetAccount(ctx, accntID)
}

// ListAccounts logs the list account params
func (s *loggingService) ListAccounts(ctx context.Context) (accnts []*Account, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "list_accounts",
			"count", len(accnts),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return s.s.ListAccounts(ctx)
}
