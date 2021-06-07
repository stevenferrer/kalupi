package postgres

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/sf9v/kalupi/etc/tx"
	"github.com/sf9v/kalupi/transaction"
)

type XactRepository struct{ db *sql.DB }

var _ transaction.Repository = (*XactRepository)(nil)

func NewXactRepository(db *sql.DB) *XactRepository {
	return &XactRepository{db: db}
}

func (tr *XactRepository) BeginTx(ctx context.Context) (tx.Tx, error) {
	return tr.db.BeginTx(ctx, nil)
}

func (tr *XactRepository) CreateXact(ctx context.Context, tx tx.Tx, xact transaction.Transaction) error {
	txx, ok := tx.(*sql.Tx)
	if !ok {
		return errors.New("expecting tx to be *sql.Tx")
	}

	// _, err := txx.ExecContext(ctx, "lock table account_transactions in exclusive mode")
	// if err != nil {
	// 	return errors.Wrap(err, "acquire exclusive lock")
	// }

	stmnt := `insert into account_transactions (
			xact_no, ledger_no, xact_type,
			account_id, xact_type_ext, amount, "desc"
		) values ($1, $2, $3, $4, $5, $6, $7)`
	_, err := txx.ExecContext(ctx, stmnt,
		xact.XactNo, xact.LedgerNo, xact.XactType,
		xact.AccountID, xact.XactTypeExt,
		xact.Amount, xact.Desc,
	)
	if err != nil {
		return errors.Wrap(err, "exec context")
	}

	return nil
}

func (tr *XactRepository) ListXacts(ctx context.Context) ([]*transaction.Transaction, error) {
	stmnt := `select xact_no, ledger_no, xact_type, account_id, 
		xact_type_ext, amount, "desc", ts from account_transactions`

	rows, err := tr.db.QueryContext(ctx, stmnt)
	if err != nil {
		return nil, errors.Wrap(err, "query row context")
	}
	defer rows.Close()

	xacts := []*transaction.Transaction{}
	for rows.Next() {
		var xact transaction.Transaction
		err = rows.Scan(
			&xact.XactNo,
			&xact.LedgerNo,
			&xact.XactType,
			&xact.AccountID,
			&xact.XactTypeExt,
			&xact.Amount,
			&xact.Desc,
			&xact.Ts,
		)
		if err != nil {
			return nil, errors.Wrap(err, "row scan")
		}

		xacts = append(xacts, &xact)
	}

	return xacts, nil
}
