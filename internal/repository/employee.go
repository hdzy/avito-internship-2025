package repository

import (
	"avito-internship-2025/internal/entity"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	_ "github.com/jmoiron/sqlx"
)

type EmployeeRepository struct {
	DB *sqlx.DB
}

func NewEmployeeRepository(db *sqlx.DB) *EmployeeRepository {
	return &EmployeeRepository{DB: db}
}

// GetEmployeeByUsername ищет сотрудника по username
func (r *EmployeeRepository) GetEmployeeByUsername(username string) (*entity.Employee, error) {
	employee := &entity.Employee{}
	query := "SELECT id, username, coins, created_at FROM employees WHERE username = $1"
	err := r.DB.Get(employee, query, username)
	if err != nil {
		// Пользователя с таким username нет в базе данных
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return employee, nil
}

// CreateEmployee создаёт нового сотрудника
func (r *EmployeeRepository) CreateEmployee(username string) (*entity.Employee, error) {
	var id int
	query := "INSERT INTO employees (username, coins) VALUES ($1, 1000) RETURNING id"
	err := r.DB.Get(&id, query, username)
	if err != nil {
		return nil, err
	}
	return r.GetEmployeeByUsername(username)
}
