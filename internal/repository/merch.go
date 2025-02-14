package repository

import (
	"avito-internship-2025/internal/entity"
	"github.com/jmoiron/sqlx"
	"log/slog"
)

type MerchRepository struct {
	DB     *sqlx.DB
	Logger *slog.Logger
}

func NewMerchRepository(db *sqlx.DB, logger *slog.Logger) *MerchRepository {
	return &MerchRepository{
		DB:     db,
		Logger: logger,
	}
}

// GetMerchByID возвращает мерч по его ID.
func (r *MerchRepository) GetMerchByID(id int) (*entity.MerchItem, error) {
	var merch entity.MerchItem
	query := "SELECT id, name, price, created_at FROM merch_items WHERE id = $1"
	err := r.DB.Get(&merch, query, id)
	if err != nil {
		r.Logger.Error("Ошибка получения мерча", slog.Any("error", err))
		return nil, err
	}
	return &merch, nil
}

// GetMerchByID возвращает мерч по его Name
func (r *MerchRepository) GetMerchByName(name string) (*entity.MerchItem, error) {
	var merch entity.MerchItem
	query := "SELECT id, name, price, created_at FROM merch_items WHERE name = $1"
	err := r.DB.Get(&merch, query, name)
	if err != nil {
		r.Logger.Error("Ошибка получения мерча", slog.Any("error", err))
		return nil, err
	}
	return &merch, nil
}
