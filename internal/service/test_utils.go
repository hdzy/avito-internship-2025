package service

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

// setupTestDB поднимает sqlmock и возвращает *sqlx.DB и сам sqlmock
func setupTestDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "Не удалось создать sqlmock")

	// Оборачиваем sqlmock в sqlx
	xdb := sqlx.NewDb(db, "sqlmock")
	return xdb, mock
}
