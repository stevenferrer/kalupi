// Package txdb is used for database testing
package txdb

import (
	"database/sql"
	"os"

	"github.com/DATA-DOG/go-txdb"

	// postgres driver
	_ "github.com/lib/pq"
)

const (
	driver  = "txdb"
	dialect = "postgres"
)

// dsn is the txdb connection string
var dsn string

func init() {
	dsn = os.Getenv("TEST_PG_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	}

	txdb.Register(driver, dialect, dsn)
}

// Open opens a txdb connection
func Open() (*sql.DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// MustOpen opens a txdb connection and panics if an error occurerd.
func MustOpen() *sql.DB {
	db, err := Open()
	if err != nil {
		panic(err)
	}

	return db
}
