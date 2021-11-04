package postgres_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stevenferrer/kalupi/etc/txdb"
	"github.com/stevenferrer/kalupi/postgres"
)

func TestMigrate(t *testing.T) {
	db := txdb.MustOpen()
	defer db.Close()
	err := postgres.Migrate(db)
	require.NoError(t, err, "migrate should not error")
}
