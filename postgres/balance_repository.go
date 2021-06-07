package postgres

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sf9v/kalupi/account"
	"github.com/sf9v/kalupi/balance"
	"github.com/sf9v/kalupi/etc/tx"
)

type BalanceRepository struct{ db *sql.DB }

var _ balance.Repository = (*BalanceRepository)(nil)

func NewBalanceRepository(db *sql.DB) *BalanceRepository {
	return &BalanceRepository{db: db}
}

func (br *BalanceRepository) BeginTx(ctx context.Context) (tx.Tx, error) {
	return br.db.BeginTx(ctx, nil)
}

func (br *BalanceRepository) GetAccntBal(ctx context.Context, tx tx.Tx, accntID account.AccountID) (*account.Balance, error) {
	txx, ok := tx.(*sql.Tx)
	if !ok {
		return nil, errors.New("expecting tx to be *sql.Tx")
	}

	// This will block when somebody is trying to insert in account_transactions
	// i.e. share lock conflicts with row exclusive mode, hence, this will block
	_, err := txx.ExecContext(ctx, "lock table account_transactions in share mode")
	if err != nil {
		return nil, errors.Wrap(err, "acquire share lock")
	}

	stmnt := `select account_id, total_debit, total_credit, current_balance, ts 
		from account_balances where account_id=$1`

	var accntBal account.Balance
	err = txx.QueryRowContext(ctx, stmnt, accntID).Scan(
		&accntBal.AccntID,
		&accntBal.TotalDebit,
		&accntBal.TotalCredit,
		&accntBal.CurrentBal,
		&accntBal.Ts,
	)
	if err != nil {
		// no transaction record yet
		if err == sql.ErrNoRows {
			return &account.Balance{AccntID: accntID}, nil
		}
		return nil, errors.Wrap(err, "query row context")
	}

	return &accntBal, nil
}
