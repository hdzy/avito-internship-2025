package service

import (
	"database/sql"
	"io"
	"testing"
	"time"

	"avito-internship-2025/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
)

func TestInfoService_GetInfo(t *testing.T) {
	db, mock := setupTestDB(t)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))
	empRepo := repository.NewEmployeeRepository(db, logger)
	txRepo := repository.NewTransactionRepository(db, logger)
	merchRepo := repository.NewMerchRepository(db, logger)

	infoService := NewInfoService(empRepo, txRepo, merchRepo, logger)

	t.Run("Успешное получение информации", func(t *testing.T) {
		username := "alice"
		employeeID := 1

		// 1) SELECT employee
		mock.ExpectQuery(`SELECT id, username, password_hash, coins, created_at FROM employees WHERE username = \$1`).
			WithArgs(username).
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password_hash", "coins", "created_at"}).
				AddRow(employeeID, username, "hash", 500, time.Now()))

		// 2) SELECT transactions
		mock.ExpectQuery(`SELECT id, employee_id, type, amount, merch_id, counterparty, created_at FROM transactions WHERE employee_id = \$1`).
			WithArgs(employeeID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "employee_id", "type", "amount", "merch_id", "counterparty", "created_at"}).
				AddRow(1, employeeID, "transfer", 100, nil, "bob", time.Now()).
				AddRow(2, employeeID, "purchase", 200, 10, nil, time.Now()),
			)

		mock.ExpectQuery(`SELECT id, name, price, created_at FROM merch_items WHERE id = \$1`).
			WithArgs(10).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price", "created_at"}).
				AddRow(10, "T-Shirt", 200, time.Now()))

		resp, err := infoService.GetInfo(username)
		require.NoError(t, err)
		assert.Equal(t, 500, resp.Coins)
		assert.Len(t, resp.CoinHistory.Sent, 1)
		assert.Equal(t, "bob", resp.CoinHistory.Sent[0].ToUser)
		assert.Len(t, resp.Inventory, 1)
		assert.Equal(t, "T-Shirt", resp.Inventory[0].Type)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Сотрудник не найден", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, username, password_hash, coins, created_at FROM employees WHERE username = \$1`).
			WithArgs("unknown").
			WillReturnError(sql.ErrNoRows)

		resp, err := infoService.GetInfo("unknown")
		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "не найден")

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
