package transaction

import (
	"context"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
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
	ListTransfers(context.Context) ([]*Transaction, error)
}

// DepositXact is a deposit transaction
type DepositXact struct {
	AccountID account.AccountID
	Amount    decimal.Decimal
}

func (dp DepositXact) Validate() error {
	return validation.Errors{
		"account_id": dp.AccountID.Validate(),
		"amount": validation.Validate(dp.Amount,
			validation.By(nonZeroDecimal),
			validation.By(nonNegativeDecimal),
		),
	}.Filter()
}

// WithdrawalXact is a withdrawal transaction
type WithdrawalXact struct {
	AccountID account.AccountID
	Amount    decimal.Decimal
}

func (wd WithdrawalXact) Validate() error {
	return validation.Errors{
		"account_id": wd.AccountID.Validate(),
		"amount": validation.Validate(wd.Amount,
			validation.By(nonZeroDecimal),
			validation.By(nonNegativeDecimal),
		),
	}.Filter()
}

// TransferXact is a transfer transaction
type TransferXact struct {
	FromAccount account.AccountID
	ToAccount   account.AccountID
	Amount      decimal.Decimal
}

func (tr TransferXact) Validate() error {
	return validation.Errors{
		"from_account": tr.FromAccount.Validate(),
		"to_account":   tr.ToAccount.Validate(),
		"amount": validation.Validate(tr.Amount,
			validation.By(nonZeroDecimal),
			validation.By(nonNegativeDecimal),
		),
	}.Filter()
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

func (s *service) MakeDeposit(ctx context.Context, dp DepositXact) (err error) {
	err = dp.Validate()
	if err != nil {
		return multierr.Combine(ErrValidation, err)
	}

	var exists bool
	exists, err = s.accountRepo.IsAccountExists(ctx, dp.AccountID)
	if err != nil {
		return errors.Wrap(err, "is account exists")
	}

	if !exists {
		return account.ErrAccountNotFound
	}

	var accnt *account.Account
	accnt, err = s.accountRepo.GetAccount(ctx, dp.AccountID)
	if err != nil {
		return errors.Wrap(err, "get account")
	}

	var cashLedgerNo ledger.LedgerNo
	cashLedgerNo, err = ledger.GetCashLedgerNo(accnt.Currency)
	if err != nil {
		return errors.Wrap(err, "get cash ledger no")
	}

	var xactNo XactNo
	xactNo, err = NewXactNo()
	if err != nil {
		return errors.Wrap(err, "new xact number")
	}

	tx, err := s.xactRepo.BeginTx(ctx)
	if err != nil {
		return errors.Wrap(err, "begin tx")

	}
	defer func() {
		// rollback if there are errors
		if err != nil {
			_ = tx.Rollback()
			return
		}

		// commit if no errors
		if commitErr := tx.Commit(); commitErr != nil {
			err = multierr.Combine(err, commitErr)
		}
	}()

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
		return
	}

	return nil
}

func (s *service) MakeWithdrawal(ctx context.Context, wd WithdrawalXact) (err error) {
	err = wd.Validate()
	if err != nil {
		return multierr.Combine(ErrValidation, err)
	}

	var exists bool
	exists, err = s.accountRepo.IsAccountExists(ctx, wd.AccountID)
	if err != nil {
		return errors.Wrap(err, "is account exists")
	}

	if !exists {
		return account.ErrAccountNotFound
	}

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
	xactNo, err = NewXactNo()
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
			_ = tx.Rollback()
			return
		}

		// commit if no errors
		if commitErr := tx.Commit(); commitErr != nil {
			err = multierr.Combine(err, commitErr)
		}
	}()

	var bal *account.Balance
	bal, err = s.balRepo.GetAccntBal(ctx, tx, accnt.AccountID)
	if err != nil {
		err = errors.Wrap(err, "get account balance")
		return
	}

	if wd.Amount.GreaterThan(bal.CurrentBal) {
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
	err = tr.Validate()
	if err != nil {
		return multierr.Combine(ErrValidation, err)
	}

	// verify sending and receiving account exists
	fromExists, err := s.accountRepo.IsAccountExists(ctx, tr.FromAccount)
	if err != nil {
		return errors.Wrap(err, "is sending account exists")
	}

	if !fromExists {
		return ErrSendingAccountNotFound
	}

	toExists, err := s.accountRepo.IsAccountExists(ctx, tr.ToAccount)
	if err != nil {
		return errors.Wrap(err, "is receiving account exists")
	}

	if !toExists {
		return ErrReceivingAccountNotFound
	}

	var from *account.Account
	from, err = s.accountRepo.GetAccount(ctx, tr.FromAccount)
	if err != nil {
		return errors.Wrap(err, "get from account")
	}

	var to *account.Account
	to, err = s.accountRepo.GetAccount(ctx, tr.ToAccount)
	if err != nil {
		return errors.Wrap(err, "get to account")
	}

	// validate that two accounts have the same currency
	if from.Currency != to.Currency {
		return errors.Wrap(ErrDifferentCurrencies, "sending and receiving account have different currencies")
	}

	var cashLedgerNo ledger.LedgerNo
	cashLedgerNo, err = ledger.GetCashLedgerNo(from.Currency)
	if err != nil {
		return errors.Wrap(err, "get cash ledger no")
	}

	var xactNo XactNo
	xactNo, err = NewXactNo()
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
			_ = tx.Rollback()
			return
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

func (s *service) ListTransfers(ctx context.Context) ([]*Transaction, error) {
	return s.xactRepo.ListTransfers(ctx)
}

const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func NewXactNo() (XactNo, error) {
	xactNoStr, err := gonanoid.Generate(alphabet, 12)
	if err != nil {
		return "", errors.Wrap(err, "generate")
	}

	return XactNo(xactNoStr), nil
}
