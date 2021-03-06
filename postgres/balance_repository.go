package postgres

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/stevenferrer/kalupi/account"
	"github.com/stevenferrer/kalupi/balance"
	"github.com/stevenferrer/kalupi/etc/tx"
)

// BalanceRepository implements the balance repository
// interface and uses postgres as back-end
type BalanceRepository struct{ db *sql.DB }

var _ balance.Repository = (*BalanceRepository)(nil)

// NewBalanceRepository returns a balance repository
func NewBalanceRepository(db *sql.DB) *BalanceRepository {
	return &BalanceRepository{db: db}
}

// BeginTx begins a new tx
func (br *BalanceRepository) BeginTx(ctx context.Context) (tx.Tx, error) {
	return br.db.BeginTx(ctx, nil)
}

// GetAccntBal retrieves the account balance within tx
func (br *BalanceRepository) GetAccntBal(ctx context.Context, tx tx.Tx, accntID account.AccountID) (*account.Balance, error) {
	txx, ok := tx.(*sql.Tx)
	if !ok {
		return nil, errors.New("expecting tx to be *sql.Tx")
	}

	// This will block when somebody is trying to insert in account_transactions
	// i.e. share lock mode conflicts with row exclusive mode, hence, the select
	// statement to account_balances view should block
	_, err := txx.ExecContext(ctx, "lock table account_transactions in share mode")
	if err != nil {
		return nil, errors.Wrap(err, "acquire share lock")
	}

	stmnt := `select account_id, total_debit, total_credit, current_balance, ts 
		from account_balances where account_id=$1`

	var accntBal account.Balance
	err = txx.QueryRowContext(ctx, stmnt, accntID).Scan(
		&accntBal.AccountID,
		&accntBal.TotalDebit,
		&accntBal.TotalCredit,
		&accntBal.CurrentBal,
		&accntBal.Ts,
	)
	if err != nil {
		// no transaction record yet
		if err == sql.ErrNoRows {
			return &account.Balance{AccountID: accntID}, nil
		}
		return nil, errors.Wrap(err, "query row context")
	}

	return &accntBal, nil
}
