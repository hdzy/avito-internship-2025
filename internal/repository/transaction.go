package repository

import (
	"avito-internship-2025/internal/entity"
	"github.com/jmoiron/sqlx"
	"log/slog"
)

type TransactionRepository struct {
	DB     *sqlx.DB
	Logger *slog.Logger
}

func NewTransactionRepository(db *sqlx.DB, logger *slog.Logger) *TransactionRepository {
	return &TransactionRepository{
		DB:     db,
		Logger: logger,
	}
}

// CreateTransaction создает запись транзакции в базе данных
func (r *TransactionRepository) CreateTransaction(tx *sqlx.Tx, t *entity.Transaction) error {
	query := `INSERT INTO transactions (employee_id, type, amount, merch_id, counterparty, created_at)
              VALUES ($1, $2, $3, $4, $5, DEFAULT)`
	_, err := tx.Exec(query, t.EmployeeID, t.Type, t.Amount, t.MerchID, t.Counterparty)
	if err != nil {
		r.Logger.Error("Ошибка создания транзакции", slog.Any("error", err))
	}
	return err
}

func (r *TransactionRepository) GetTransactionsByEmployee(employeeID int) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	query := "SELECT id, employee_id, type, amount, merch_id, counterparty, created_at FROM transactions WHERE employee_id = $1"
	err := r.DB.Select(&transactions, query, employeeID)
	if err != nil {
		r.Logger.Error("Ошибка получения транзакций", slog.Any("error", err))
		return nil, err
	}
	return transactions, nil
}
