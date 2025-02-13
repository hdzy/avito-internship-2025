package repository

import (
	"avito-internship-2025/internal/entity"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	_ "github.com/jmoiron/sqlx"
	"log/slog"
)

type EmployeeRepository struct {
	DB     *sqlx.DB
	Logger *slog.Logger
}

func NewEmployeeRepository(db *sqlx.DB, logger *slog.Logger) *EmployeeRepository {
	return &EmployeeRepository{
		DB:     db,
		Logger: logger,
	}
}

// GetEmployeeByUsername ищет сотрудника по username
func (r *EmployeeRepository) GetEmployeeByUsername(username string) (*entity.Employee, error) {
	employee := &entity.Employee{}
	query := "SELECT id, username, password_hash, coins, created_at FROM employees WHERE username = $1"
	err := r.DB.Get(employee, query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return employee, nil
}

// CreateEmployee создаёт нового сотрудника
func (r *EmployeeRepository) CreateEmployee(username string, passwordHash string) (*entity.Employee, error) {
	var id int
	query := "INSERT INTO employees (username, password_hash) VALUES ($1, $2) RETURNING id"
	err := r.DB.Get(&id, query, username, passwordHash)
	if err != nil {
		return nil, err
	}
	r.Logger.Info("Сотрудник создан", slog.String("username", username), slog.Int("id", id))
	return r.GetEmployeeByUsername(username)
}

// UpdateEmployeeBalance обновляет баланс сотрудника
func (r *EmployeeRepository) UpdateEmployeeBalance(tx *sqlx.Tx, employeeID, newBalance int) error {
	query := "UPDATE employees SET coins = $1 WHERE id = $2"
	_, err := tx.Exec(query, newBalance, employeeID)
	if err != nil {
		r.Logger.Error("Ошибка обновления баланса", slog.Any("error", err))
	}
	return err
}
