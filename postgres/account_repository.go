package postgres

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/sf9v/kalupi/account"
)

// AccountRepository implements the
type AccountRepository struct{ db *sql.DB }

var _ account.Repository = (*AccountRepository)(nil)

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (ar *AccountRepository) CreateAccount(ctx context.Context, ac account.Account) (account.AccountID, error) {
	stmnt := "insert into accounts (account_id, currency) values ($1, $2)"
	_, err := ar.db.ExecContext(ctx, stmnt, ac.AccountID, ac.Currency)
	if err != nil {
		return "", errors.Wrap(err, "exec context")
	}

	return ac.AccountID, nil
}

func (ar *AccountRepository) GetAccount(ctx context.Context, acID account.AccountID) (*account.Account, error) {
	stmnt := `select account_id, currency from accounts
		where account_id = $1`

	var ac account.Account
	err := ar.db.QueryRowContext(ctx, stmnt, acID).
		Scan(&ac.AccountID, &ac.Currency)
	if err != nil {
		return nil, errors.Wrap(err, "query row context")
	}

	return &ac, nil
}

func (ar *AccountRepository) ListAccounts(ctx context.Context) ([]*account.Account, error) {
	stmnt := `select account_id, currency from accounts`

	rows, err := ar.db.QueryContext(ctx, stmnt)
	if err != nil {
		return nil, errors.Wrap(err, "query context")
	}
	defer rows.Close()

	acs := []*account.Account{}
	for rows.Next() {
		var ac account.Account
		err = rows.Scan(&ac.AccountID, &ac.Currency)
		if err != nil {
			return nil, errors.Wrap(err, "row scan")
		}
		acs = append(acs, &ac)
	}

	return acs, nil
}
