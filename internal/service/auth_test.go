package service

import (
	"database/sql"
	"errors"
	"io"
	"testing"
	"time"

	"avito-internship-2025/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

func TestAuthService_Authenticate(t *testing.T) {
	db, mock := setupTestDB(t)
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))
	employeeRepo := repository.NewEmployeeRepository(db, logger)

	authService := NewAuthService(employeeRepo, "test-secret", logger)

	t.Run("Существующий пользователь + правильный пароль", func(t *testing.T) {
		username := "alice"
		password := "correct_password"
		hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		mock.ExpectQuery(`SELECT id, username, password_hash, coins, created_at FROM employees WHERE username = \$1`).
			WithArgs(username).
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password_hash", "coins", "created_at"}).
				AddRow(1, username, string(hashed), 1000, time.Now()))

		req := AuthRequest{Username: username, Password: password}
		resp, err := authService.Authenticate(req)
		require.NoError(t, err)
		assert.NotEmpty(t, resp.Token, "JWT-токен не должен быть пустым")

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Существующий пользователь + неверный пароль", func(t *testing.T) {
		username := "bob"
		password := "wrong_password"
		hashed, _ := bcrypt.GenerateFromPassword([]byte("correct_password"), bcrypt.DefaultCost)

		mock.ExpectQuery(`SELECT id, username, password_hash, coins, created_at FROM employees WHERE username = \$1`).
			WithArgs(username).
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password_hash", "coins", "created_at"}).
				AddRow(2, username, string(hashed), 500, time.Now()))

		req := AuthRequest{Username: username, Password: password}
		resp, err := authService.Authenticate(req)
		require.Error(t, err)
		assert.Nil(t, resp)
		var vErr *jwt.ValidationError
		assert.True(t, errors.As(err, &vErr), "Ожидаем ErrInvalidCredentials")

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Новый пользователь — создаём запись", func(t *testing.T) {
		username := "charlie"
		password := "some_password"

		mock.ExpectQuery(`SELECT id, username, password_hash, coins, created_at FROM employees WHERE username = \$1`).
			WithArgs(username).
			WillReturnError(sql.ErrNoRows)

		mock.ExpectQuery(`INSERT INTO employees \(username, password_hash\) VALUES \(\$1, \$2\) RETURNING id`).
			WithArgs(username, sqlmock.AnyArg()). // проверяем тип
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(3))

		mock.ExpectQuery(`SELECT id, username, password_hash, coins, created_at FROM employees WHERE username = \$1`).
			WithArgs(username).
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password_hash", "coins", "created_at"}).
				AddRow(3, username, "somehash", 0, time.Now()))

		req := AuthRequest{Username: username, Password: password}
		resp, err := authService.Authenticate(req)
		require.NoError(t, err)
		assert.NotEmpty(t, resp.Token)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
