package transaction

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"go.uber.org/multierr"

	"github.com/sf9v/kalupi/account"
	"github.com/sf9v/kalupi/balance"
	"github.com/sf9v/kalupi/ledger"
)

// Service is the transaction service
type Service interface {
	MakeDeposit(context.Context, DepositXact) error
	MakeWithdrawal(context.Context, WithdrawalXact) error
	MakeTransfer(context.Context, TransferXact) error
}

// DepositXact is a deposit transaction
type DepositXact struct {
	XactNo    XactNo
	AccountID account.AccountID
	Amount    decimal.Decimal
}

// WithdrawalXact is a withdrawal transaction
type WithdrawalXact struct {
	XactNo    XactNo
	AccountID account.AccountID
	Amount    decimal.Decimal
}

// TransferXact is a transfer transaction
type TransferXact struct {
	XactNo    XactNo
	FromAccnt account.AccountID
	ToAccnt   account.AccountID
	Amount    decimal.Decimal
}

type service struct {
	accountRepo account.Repository
	ledgerRepo  ledger.Repository
	xactRepo    Repository
	balRepo     balance.Repository
}

var _ Service = (*service)(nil)

func NewService(
	accountRepo account.Repository,
	ledgerRepo ledger.Repository,
	xactRepo Repository,
	balRepo balance.Repository,
) Service {
	return &service{
		accountRepo: accountRepo,
		ledgerRepo:  ledgerRepo,
		xactRepo:    xactRepo,
		balRepo:     balRepo,
	}
}

func (s *service) MakeDeposit(ctx context.Context, dp DepositXact) error {
	accnt, err := s.accountRepo.GetAccount(ctx, dp.AccountID)
	if err != nil {
		return errors.Wrap(err, "get account")
	}

	cashLedgerNo, err := ledger.GetCashLedgerNo(accnt.Currency)
	if err != nil {
		return errors.Wrap(err, "get cash ledger no")
	}

	tx, err := s.xactRepo.BeginTx(ctx)
	if err != nil {
		return errors.Wrap(err, "begin tx")
	}

	err = s.xactRepo.CreateXact(ctx, tx, Transaction{
		XactNo:      dp.XactNo,
		LedgerNo:    cashLedgerNo,
		XactType:    XactTypeDebit, // debit ledger's cash
		AccountID:   accnt.AccountID,
		XactTypeExt: XactTypeExtDeposit, // credit account's cash
		Amount:      dp.Amount,
		Desc:        fmt.Sprintf("Cash deposit from %s", accnt.AccountID),
	})
	if err != nil {
		err = errors.Wrap(err, "create dp xact")
		txErr := tx.Rollback()
		return multierr.Combine(err, txErr)
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "commit tx")
	}

	return nil
}

func (s *service) MakeWithdrawal(ctx context.Context, wd WithdrawalXact) error {
	// TODO: Validate balance

	accnt, err := s.accountRepo.GetAccount(ctx, wd.AccountID)
	if err != nil {
		return errors.Wrap(err, "get account")
	}

	cashLedgerNo, err := ledger.GetCashLedgerNo(accnt.Currency)
	if err != nil {
		return errors.Wrap(err, "get cash ledger no")
	}

	tx, err := s.xactRepo.BeginTx(ctx)
	if err != nil {
		return errors.Wrap(err, "begin tx")
	}

	accntBal, err := s.balRepo.GetAccntBal(ctx, tx, accnt.AccountID)
	if err != nil {
		err = errors.Wrap(err, "get account balance")
		txErr := tx.Rollback()
		return multierr.Combine(err, txErr)
	}

	if wd.Amount.GreaterThan(accntBal.CurrentBal) {
		txErr := tx.Rollback()
		return multierr.Combine(ErrInsufficientBalance, txErr)
	}

	err = s.xactRepo.CreateXact(ctx, tx, Transaction{
		XactNo:      wd.XactNo,
		LedgerNo:    cashLedgerNo,
		XactType:    XactTypeCredit, // credit ledger's cash
		AccountID:   accnt.AccountID,
		XactTypeExt: XactTypeExtWithdrawal, // debit account's cash
		Amount:      wd.Amount,
		Desc:        fmt.Sprintf("Cash withdrawal from %s", accnt.AccountID),
	})
	if err != nil {
		err = errors.Wrap(err, "create wd xact")
		txErr := tx.Rollback()
		return multierr.Combine(err, txErr)
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "commit tx")
	}

	return nil
}

func (s *service) MakeTransfer(ctx context.Context, tr TransferXact) error {
	// TODO: validate balance

	fromAccnt, err := s.accountRepo.GetAccount(ctx, tr.FromAccnt)
	if err != nil {
		return errors.Wrap(err, "get from account")
	}

	toAccnt, err := s.accountRepo.GetAccount(ctx, tr.ToAccnt)
	if err != nil {
		return errors.Wrap(err, "get to account")
	}

	// validate that two accounts have the same currency
	if fromAccnt.Currency != toAccnt.Currency {
		return ErrMustHaveSameCurrency
	}

	cashLedgerNo, err := ledger.GetCashLedgerNo(fromAccnt.Currency)
	if err != nil {
		return errors.Wrap(err, "get cash ledger no")
	}

	tx, err := s.xactRepo.BeginTx(ctx)
	if err != nil {
		return errors.Wrap(err, "begin tx")
	}
	// TODO: Use deferred rollback and commit??

	frmAccntBal, err := s.balRepo.GetAccntBal(ctx, tx, fromAccnt.AccountID)
	if err != nil {
		err = errors.Wrap(err, "get from account balance")
		txErr := tx.Rollback()
		return multierr.Combine(err, txErr)
	}

	if tr.Amount.GreaterThan(frmAccntBal.CurrentBal) {
		txErr := tx.Rollback()
		return multierr.Combine(ErrInsufficientBalance, txErr)
	}

	// debit the sending account
	err = s.xactRepo.CreateXact(ctx, tx, Transaction{
		XactNo:      tr.XactNo,
		LedgerNo:    cashLedgerNo,
		XactType:    XactTypeCredit, // credit ledger's cash
		AccountID:   fromAccnt.AccountID,
		XactTypeExt: XactTypeExtSndTransfer, // debit sending account's cash
		Amount:      tr.Amount,
		Desc:        fmt.Sprintf("Outgoing cash transfer to %s", toAccnt.AccountID),
	})
	if err != nil {
		err = errors.Wrap(err, "create snd xact")
		txErr := tx.Rollback()
		return multierr.Combine(err, txErr)
	}

	// credit the receiving account
	err = s.xactRepo.CreateXact(ctx, tx, Transaction{
		XactNo:      tr.XactNo,
		LedgerNo:    cashLedgerNo,
		XactType:    XactTypeDebit, // debit ledger's cash
		AccountID:   toAccnt.AccountID,
		XactTypeExt: XactTypeExtRcvTransfer, // credit receiving account's cash
		Amount:      tr.Amount,
		Desc:        fmt.Sprintf("Incoming cash transfer from %s", fromAccnt.AccountID),
	})
	if err != nil {
		err = errors.Wrap(err, "create rcv xact")
		txErr := tx.Rollback()
		return multierr.Combine(err, txErr)
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "commit tx")
	}

	return nil
}
