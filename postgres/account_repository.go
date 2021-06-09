package postgres

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/sf9v/kalupi/account"
)

// AccountRepository implements the account repository
// interface and uses postgres as back-end
type AccountRepository struct{ db *sql.DB }

var _ account.Repository = (*AccountRepository)(nil)

// NewAccountRepository returns an account repository
func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

// CreateAccount creates an account
func (ar *AccountRepository) CreateAccount(ctx context.Context, accnt account.Account) (account.AccountID, error) {
	stmnt := "insert into accounts (account_id, currency) values ($1, $2)"
	_, err := ar.db.ExecContext(ctx, stmnt, accnt.AccountID, accnt.Currency)
	if err != nil {
		return "", errors.Wrap(err, "exec context")
	}

	return accnt.AccountID, nil
}

// GetAccount retrieves an account
func (ar *AccountRepository) GetAccount(ctx context.Context, accntID account.AccountID) (*account.Account, error) {
	stmnt := `select account_id, currency from accounts
		where account_id = $1`

	var ac account.Account
	err := ar.db.QueryRowContext(ctx, stmnt, accntID).
		Scan(&ac.AccountID, &ac.Currency)
	if err != nil {
		return nil, errors.Wrap(err, "query row context")
	}

	return &ac, nil
}

// ListAccounts retrieves the list of accounts
func (ar *AccountRepository) ListAccounts(ctx context.Context) ([]*account.Account, error) {
	stmnt := `select account_id, currency from accounts`

	rows, err := ar.db.QueryContext(ctx, stmnt)
	if err != nil {
		return nil, errors.Wrap(err, "query context")
	}
	defer rows.Close()

	accnts := []*account.Account{}
	for rows.Next() {
		var accnt account.Account
		err = rows.Scan(&accnt.AccountID, &accnt.Currency)
		if err != nil {
			return nil, errors.Wrap(err, "row scan")
		}
		accnts = append(accnts, &accnt)
	}

	return accnts, nil
}

// IsAccountExists returns true if an account exists
func (ar *AccountRepository) IsAccountExists(ctx context.Context, accntID account.AccountID) (bool, error) {
	stmnt := "select exists(select 1 from accounts where account_id=$1)"
	var exists bool
	err := ar.db.QueryRowContext(ctx, stmnt, accntID).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, errors.Wrap(err, "query row context")
	}

	return exists, nil
}
