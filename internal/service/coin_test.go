// internal/service/coin_test.go
package service

import (
	"database/sql"
	"fmt"
	"io"
	"regexp"
	"testing"
	"time"

	"avito-internship-2025/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
)

func TestCoinService_TransferCoins(t *testing.T) {
	// Создаем тестовую базу через sqlmock и оборачиваем её через sqlx
	db, mock := setupTestDB(t)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))
	empRepo := repository.NewEmployeeRepository(db, logger)
	txRepo := repository.NewTransactionRepository(db, logger)
	merchRepo := repository.NewMerchRepository(db, logger)
	coinService := NewCoinService(empRepo, txRepo, merchRepo, logger)

	// Определяем ожидаемые запросы
	queryEmployee := regexp.QuoteMeta("SELECT id, username, password_hash, coins, created_at FROM employees WHERE username = $1")
	updateEmployee := regexp.QuoteMeta("UPDATE employees SET coins = $1 WHERE id = $2")
	insertTransaction := regexp.QuoteMeta("INSERT INTO transactions (employee_id, type, amount, merch_id, counterparty, created_at) VALUES")

	t.Run("Успешный перевод монет", func(t *testing.T) {
		fromUser := "alice"
		toUser := "bob"
		amount := 100

		mock.ExpectBegin()

		mock.ExpectQuery(queryEmployee).
			WithArgs(fromUser).
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password_hash", "coins", "created_at"}).
				AddRow(1, fromUser, "hash", 300, time.Now()))

		mock.ExpectQuery(queryEmployee).
			WithArgs(toUser).
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password_hash", "coins", "created_at"}).
				AddRow(2, toUser, "hash2", 50, time.Now()))

		mock.ExpectExec(updateEmployee).
			WithArgs(200, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec(updateEmployee).
			WithArgs(150, 2).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec(insertTransaction).
			WithArgs(1, "transfer", amount, nil, toUser).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit()

		err := coinService.TransferCoins(fromUser, toUser, amount)
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Недостаточно монет", func(t *testing.T) {
		// Здесь функция вернет ошибку из-за недостаточного баланса, но, согласно текущей реализации,
		// транзакция завершается commit, а не rollback.
		mock.ExpectBegin()

		mock.ExpectQuery(queryEmployee).
			WithArgs("alice").
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password_hash", "coins", "created_at"}).
				AddRow(1, "alice", "hash", 10, time.Now()))

		mock.ExpectQuery(queryEmployee).
			WithArgs("bob").
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password_hash", "coins", "created_at"}).
				AddRow(2, "bob", "hash2", 50, time.Now()))

		// Ожидаем commit (из-за отсутствия именованных возвращаемых значений функция вызывает commit даже при ошибке)
		mock.ExpectCommit()

		err := coinService.TransferCoins("alice", "bob", 100)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "недостаточно монет")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Отправитель не найден", func(t *testing.T) {
		mock.ExpectBegin()

		// Возвращаем пустой результат для запроса отправителя
		mock.ExpectQuery(queryEmployee).
			WithArgs("alice").
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password_hash", "coins", "created_at"}))

		mock.ExpectCommit()

		err := coinService.TransferCoins("alice", "bob", 50)
		require.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("отправитель %s не найден", "alice"))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCoinService_BuyMerch(t *testing.T) {
	db, mock := setupTestDB(t)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))
	empRepo := repository.NewEmployeeRepository(db, logger)
	txRepo := repository.NewTransactionRepository(db, logger)
	merchRepo := repository.NewMerchRepository(db, logger)
	coinService := NewCoinService(empRepo, txRepo, merchRepo, logger)

	// Определяем ожидаемые запросы
	queryMerch := regexp.QuoteMeta("SELECT id, name, price, created_at FROM merch_items WHERE name = $1")
	queryEmployee := regexp.QuoteMeta("SELECT id, username, password_hash, coins, created_at FROM employees WHERE username = $1")
	updateEmployee := regexp.QuoteMeta("UPDATE employees SET coins = $1 WHERE id = $2")
	insertTransaction := regexp.QuoteMeta("INSERT INTO transactions (employee_id, type, amount, merch_id, counterparty, created_at) VALUES")

	t.Run("Успешная покупка мерча", func(t *testing.T) {
		username := "alice"
		itemName := "T-Shirt"
		price := 200

		mock.ExpectQuery(queryMerch).
			WithArgs(itemName).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price", "created_at"}).
				AddRow(10, itemName, price, time.Now()))

		mock.ExpectQuery(queryEmployee).
			WithArgs(username).
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password_hash", "coins", "created_at"}).
				AddRow(1, username, "hash", 1000, time.Now()))

		mock.ExpectBegin()

		mock.ExpectExec(updateEmployee).
			WithArgs(800, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec(insertTransaction).
			WithArgs(1, "purchase", price, 10, nil).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit()

		err := coinService.BuyMerch(username, itemName)
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Товар не найден", func(t *testing.T) {
		mock.ExpectQuery(queryMerch).
			WithArgs("unknown_item").
			WillReturnError(sql.ErrNoRows)

		err := coinService.BuyMerch("alice", "unknown_item")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no rows")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Сотрудник не найден", func(t *testing.T) {
		mock.ExpectQuery(queryMerch).
			WithArgs("T-Shirt").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price", "created_at"}).
				AddRow(10, "T-Shirt", 200, time.Now()))
		mock.ExpectQuery(queryEmployee).
			WithArgs("unknown_user").
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password_hash", "coins", "created_at"}))

		err := coinService.BuyMerch("unknown_user", "T-Shirt")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "не найден")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
