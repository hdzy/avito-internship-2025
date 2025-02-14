package service

import (
	"fmt"

	"avito-internship-2025/internal/entity"
	"avito-internship-2025/internal/repository"
	"log/slog"
)

type InfoService struct {
	EmployeeRepo    *repository.EmployeeRepository
	TransactionRepo *repository.TransactionRepository
	MerchRepo       *repository.MerchRepository
	Logger          *slog.Logger
}

func NewInfoService(empRepo *repository.EmployeeRepository, txRepo *repository.TransactionRepository, merchRepo *repository.MerchRepository, logger *slog.Logger) *InfoService {
	return &InfoService{
		EmployeeRepo:    empRepo,
		TransactionRepo: txRepo,
		MerchRepo:       merchRepo,
		Logger:          logger,
	}
}

func (s *InfoService) GetInfo(username string) (*entity.InfoResponse, error) {
	emp, err := s.EmployeeRepo.GetEmployeeByUsername(username)
	if err != nil {
		s.Logger.Error("Ошибка получения сотрудника", slog.Any("error", err))
		return nil, err
	}
	if emp == nil {
		return nil, fmt.Errorf("сотрудник %s не найден", username)
	}

	transactions, err := s.TransactionRepo.GetTransactionsByEmployee(emp.ID)
	if err != nil {
		s.Logger.Error("Ошибка получения транзакций", slog.Any("error", err))
		return nil, err
	}

	info := &entity.InfoResponse{
		Coins:     emp.Coins,
		Inventory: []entity.InventoryItem{},
		CoinHistory: entity.CoinHistory{
			Received: []entity.TransactionHistoryItem{},
			Sent:     []entity.TransactionHistoryItem{},
		},
	}

	for _, t := range transactions {
		if t.Type == "transfer" {
			info.CoinHistory.Sent = append(info.CoinHistory.Sent, entity.TransactionHistoryItem{
				ToUser: func() string {
					if t.Counterparty != nil {
						return *t.Counterparty
					}
					return ""
				}(),
				Amount: t.Amount,
			})
		} else if t.Type == "purchase" {
			var merchName string
			if t.MerchID != nil {
				merch, err := s.MerchRepo.GetMerchByID(*t.MerchID)
				if err != nil {
					s.Logger.Warn("Не удалось получить информацию о мерче", slog.Any("error", err))
					merchName = "unknown"
				} else {
					merchName = merch.Name
				}
			} else {
				merchName = "unknown"
			}

			found := false
			for i, inv := range info.Inventory {
				if inv.Type == merchName {
					info.Inventory[i].Quantity++
					found = true
					break
				}
			}
			if !found {
				info.Inventory = append(info.Inventory, entity.InventoryItem{
					Type:     merchName,
					Quantity: 1,
				})
			}
		}
	}

	return info, nil
}
