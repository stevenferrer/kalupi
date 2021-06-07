package postgres

import (
	"database/sql"

	// postgres driver
	_ "github.com/lib/pq"
	"github.com/lopezator/migrator"
)

// defaultOpts is the migration options
var defaultOpts = []migrator.Option{migrator.WithLogger(newNopLogger())}

// nopLogger is a nop migrator.Logger
type nopLogger struct{}

func newNopLogger() migrator.Logger {
	return &nopLogger{}
}

func (l *nopLogger) Printf(string, ...interface{}) {}

// Migrate migrates the database to the latest version
func Migrate(db *sql.DB, opts ...migrator.Option) error {
	if len(opts) == 0 {
		opts = defaultOpts
	}

	opts = append(opts, migrations)

	m, err := migrator.New(opts...)
	if err != nil {
		return err
	}

	return m.Migrate(db)
}

// migrations are list of database migrations
var migrations = migrator.Migrations(
	&migrator.Migration{
		Name: "create accounts table",
		Func: func(tx *sql.Tx) error {
			stmnt := `create table accounts (
				account_id varchar(64) primary key,
				currency varchar(3) not null
			)`
			if _, err := tx.Exec(stmnt); err != nil {
				return err
			}

			return nil
		},
	},
	&migrator.Migration{
		Name: "create ledgers table",
		Func: func(tx *sql.Tx) error {
			stmnt := `create table ledgers (
				ledger_no varchar(64) primary key,
				account_type varchar(3) not null,
				currency varchar(3) not null,
				name varchar(64) not null
			)`
			if _, err := tx.Exec(stmnt); err != nil {
				return err
			}

			return nil
		},
	},

	&migrator.Migration{
		Name: "create account_transactions table",
		Func: func(tx *sql.Tx) error {
			stmnt := `create table account_transactions (
				xact_no varchar not null, -- reference number
				ledger_no varchar(64) not null, -- fk
				xact_type varchar(3) not null,  
				account_id varchar(64) not null, -- fk
				xact_type_ext varchar(3) not null,
				amount numeric(15, 4),
				"desc" text not null default '',
				ts timestamptz not null default now(),
				constraint fk_ledger
					foreign key (ledger_no)
						references ledgers(ledger_no),
				constraint fk_account
					foreign key (account_id)
						references accounts(account_id)
				-- add composite primary keys??
			)`
			if _, err := tx.Exec(stmnt); err != nil {
				return err
			}

			return nil
		},
	},

	&migrator.Migration{
		Name: "create account_balances view",
		Func: func(tx *sql.Tx) error {
			stmnt := `create view account_balances as 
				select 
					account_id,
					coalesce((
						select sum(amount) from account_transactions
						where account_id=at.account_id and
							xact_type_ext in ('RTr', 'Dp')
					), 0) as total_credit,
					coalesce((
						select sum(amount) from account_transactions
						where account_id=at.account_id and
							xact_type_ext in ('STr','Wd')
					), 0) as total_debit,
					coalesce((
						select sum(amount) from account_transactions
						where account_id=at.account_id and
							xact_type_ext in ('RTr', 'Dp')
					), 0) - coalesce((
						select sum(amount) from account_transactions
						where account_id=at.account_id and
							xact_type_ext in ('STr','Wd')
					), 0) as current_balance,
					now() as ts
				from  account_transactions at
			`
			if _, err := tx.Exec(stmnt); err != nil {
				return err
			}

			return nil
		},
	},
)
