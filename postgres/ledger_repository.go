package postgres

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/stevenferrer/kalupi/ledger"
)

// LedgerRepository implements the ledger repository
// interface and uses postgres as back-end
type LedgerRepository struct{ db *sql.DB }

var _ ledger.Repository = (*LedgerRepository)(nil)

// NewLedgerRepository returns a new ledger repository
func NewLedgerRepository(db *sql.DB) *LedgerRepository {
	return &LedgerRepository{db: db}
}

// CreateLedgersIfNotExists will create the ledgers if it doesn't exists in database yet
func (lr *LedgerRepository) CreateLedgersIfNotExists(ctx context.Context, lgs ...ledger.Ledger) error {
	for _, lg := range lgs {
		exist, err := lr.isLedgerExists(ctx, lg.LedgerNo)
		if err != nil {
			return errors.Wrap(err, "is ledger exist")
		}

		if !exist {
			err = lr.createLedger(ctx, lg)
			if err != nil {
				return errors.Wrap(err, "create ledger")
			}
		}
	}

	return nil
}

// GetLedger retrieves the ledger
func (lr *LedgerRepository) GetLedger(ctx context.Context, ledgerNo ledger.LedgerNo) (*ledger.Ledger, error) {
	stmnt := `select ledger_no, account_type, currency, name
		from ledgers where ledger_no = $1`
	var lg ledger.Ledger
	err := lr.db.QueryRowContext(ctx, stmnt, ledgerNo).
		Scan(&lg.LedgerNo, &lg.AccountType, &lg.Currency, &lg.Name)
	if err != nil {
		return nil, errors.Wrap(err, "query row context")
	}

	return &lg, nil
}

// ListLedgers retrieves the list of ledgers
func (lr *LedgerRepository) ListLedgers(ctx context.Context) ([]*ledger.Ledger, error) {
	stmnt := `select ledger_no, account_type, currency, name from ledgers`

	rows, err := lr.db.QueryContext(ctx, stmnt)
	if err != nil {
		return nil, errors.Wrap(err, "query context")
	}
	defer rows.Close()

	lgs := []*ledger.Ledger{}
	for rows.Next() {
		var lg ledger.Ledger
		err = rows.Scan(&lg.LedgerNo, &lg.AccountType, &lg.Currency, &lg.Name)
		if err != nil {
			return nil, errors.Wrap(err, "row scan")
		}
		lgs = append(lgs, &lg)
	}

	return lgs, nil
}

// createLedger is a helper method for creating a ledger
func (lr *LedgerRepository) createLedger(ctx context.Context, lg ledger.Ledger) error {
	stmnt := `insert into ledgers (ledger_no, account_type, currency, name)
		values ($1, $2, $3, $4)`
	_, err := lr.db.ExecContext(ctx, stmnt, lg.LedgerNo,
		lg.AccountType, lg.Currency, lg.Name)
	if err != nil {
		return errors.Wrap(err, "exec context")
	}

	return nil
}

// isLedgerExists is a helper method for checking ledger existence
func (lr *LedgerRepository) isLedgerExists(ctx context.Context, ledgerNo ledger.LedgerNo) (bool, error) {
	stmnt := "select exists(select 1 from ledgers where ledger_no=$1)"
	var exists bool
	err := lr.db.QueryRowContext(ctx, stmnt, ledgerNo).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, errors.Wrap(err, "query row context")
	}

	return exists, nil
}
