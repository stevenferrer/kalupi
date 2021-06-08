package transaction

import (
	"context"
	"fmt"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"go.uber.org/multierr"

	"github.com/sf9v/kalupi/account"
	"github.com/sf9v/kalupi/balance"
	"github.com/sf9v/kalupi/etc/tx"
	"github.com/sf9v/kalupi/ledger"
)

// TODO: return the XactNo??

// Service is the transaction service
type Service interface {
	MakeDeposit(context.Context, DepositXact) error
	MakeWithdrawal(context.Context, WithdrawalXact) error
	MakeTransfer(context.Context, TransferXact) error
}

// DepositXact is a deposit transaction
type DepositXact struct {
	AccountID account.AccountID
	Amount    decimal.Decimal
}

// WithdrawalXact is a withdrawal transaction
type WithdrawalXact struct {
	AccountID account.AccountID
	Amount    decimal.Decimal
}

// TransferXact is a transfer transaction
type TransferXact struct {
	FromAccount account.AccountID
	ToAccount   account.AccountID
	Amount      decimal.Decimal
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

	xactNo, err := newXactNo()
	if err != nil {
		return errors.Wrap(err, "new xact no")
	}

	tx, err := s.xactRepo.BeginTx(ctx)
	if err != nil {
		return errors.Wrap(err, "begin tx")
	}

	err = s.xactRepo.CreateXact(ctx, tx, Transaction{
		XactNo:      xactNo,
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

func (s *service) MakeWithdrawal(ctx context.Context, wd WithdrawalXact) (err error) {
	var accnt *account.Account
	accnt, err = s.accountRepo.GetAccount(ctx, wd.AccountID)
	if err != nil {
		return errors.Wrap(err, "get account")
	}

	var cashLedgerNo ledger.LedgerNo
	cashLedgerNo, err = ledger.GetCashLedgerNo(accnt.Currency)
	if err != nil {
		return errors.Wrap(err, "get cash ledger no")
	}

	var xactNo XactNo
	xactNo, err = newXactNo()
	if err != nil {
		return errors.Wrap(err, "new xact no")
	}

	var tx tx.Tx
	tx, err = s.xactRepo.BeginTx(ctx)
	if err != nil {
		return errors.Wrap(err, "begin tx")
	}
	defer func() {
		// rollback if there are errors
		if err != nil {
			if rollBackErr := tx.Rollback(); rollBackErr != nil {
				err = multierr.Combine(err, rollBackErr)
				return
			}
		}

		// commit if no errors
		if commitErr := tx.Commit(); commitErr != nil {
			err = multierr.Combine(err, commitErr)
		}
	}()

	var accntBal *account.Balance
	accntBal, err = s.balRepo.GetAccntBal(ctx, tx, accnt.AccountID)
	if err != nil {
		err = errors.Wrap(err, "get account balance")
		return
	}

	if wd.Amount.GreaterThan(accntBal.CurrentBal) {
		err = ErrInsufficientBalance
		return
	}

	err = s.xactRepo.CreateXact(ctx, tx, Transaction{
		XactNo:      xactNo,
		LedgerNo:    cashLedgerNo,
		XactType:    XactTypeCredit, // credit ledger's cash
		AccountID:   accnt.AccountID,
		XactTypeExt: XactTypeExtWithdrawal, // debit account's cash
		Amount:      wd.Amount,
		Desc:        fmt.Sprintf("Cash withdrawal from %s", accnt.AccountID),
	})
	if err != nil {
		err = errors.Wrap(err, "create wd xact")
		return
	}

	return nil
}

func (s *service) MakeTransfer(ctx context.Context, tr TransferXact) (err error) {
	var from *account.Account
	from, err = s.accountRepo.GetAccount(ctx, tr.FromAccount)
	if err != nil {
		return errors.Wrap(err, "get from account")
	}

	var to *account.Account
	to, err = s.accountRepo.GetAccount(ctx, tr.ToAccount)
	if err != nil {
		err = errors.Wrap(err, "get to account")
		return
	}

	// validate that two accounts have the same currency
	if from.Currency != to.Currency {
		err = ErrMustHaveSameCurrency
		return
	}

	var cashLedgerNo ledger.LedgerNo
	cashLedgerNo, err = ledger.GetCashLedgerNo(from.Currency)
	if err != nil {
		err = errors.Wrap(err, "get cash ledger no")
		return
	}

	var xactNo XactNo
	xactNo, err = newXactNo()
	if err != nil {
		err = errors.Wrap(err, "new xact no")
		return
	}

	var tx tx.Tx
	tx, err = s.xactRepo.BeginTx(ctx)
	if err != nil {
		err = errors.Wrap(err, "begin tx")
		return
	}
	defer func() {
		// rollback if there are errors
		if err != nil {
			if rollBackErr := tx.Rollback(); rollBackErr != nil {
				err = multierr.Combine(err, rollBackErr)
				return
			}
		}

		// commit if no errors
		if commitErr := tx.Commit(); commitErr != nil {
			err = multierr.Combine(err, commitErr)
		}
	}()

	var fromBal *account.Balance
	fromBal, err = s.balRepo.GetAccntBal(ctx, tx, from.AccountID)
	if err != nil {
		err = errors.Wrap(err, "get from account balance")
		return
	}

	// sending account must have sufficient balance
	if tr.Amount.GreaterThan(fromBal.CurrentBal) {
		err = ErrInsufficientBalance
		return
	}

	// debit the sending account
	err = s.xactRepo.CreateXact(ctx, tx, Transaction{
		XactNo:      xactNo,
		LedgerNo:    cashLedgerNo,
		XactType:    XactTypeCredit, // credit ledger's cash
		AccountID:   from.AccountID,
		XactTypeExt: XactTypeExtSndTransfer, // debit sending account's cash
		Amount:      tr.Amount,
		Desc:        fmt.Sprintf("Outgoing cash transfer to %s", to.AccountID),
	})
	if err != nil {
		err = errors.Wrap(err, "create snd xact")
		return
	}

	// credit the receiving account
	err = s.xactRepo.CreateXact(ctx, tx, Transaction{
		XactNo:      xactNo,
		LedgerNo:    cashLedgerNo,
		XactType:    XactTypeDebit, // debit ledger's cash
		AccountID:   to.AccountID,
		XactTypeExt: XactTypeExtRcvTransfer, // credit receiving account's cash
		Amount:      tr.Amount,
		Desc:        fmt.Sprintf("Incoming cash transfer from %s", from.AccountID),
	})
	if err != nil {
		err = errors.Wrap(err, "create rcv xact")
		return
	}

	return nil
}

const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func newXactNo() (XactNo, error) {
	xactNoStr, err := gonanoid.Generate(alphabet, 12)
	if err != nil {
		return "", errors.Wrap(err, "generate")
	}

	return XactNo(xactNoStr), nil
}
