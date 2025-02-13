package service

import (
	"errors"
	"fmt"

	"avito-internship-2025/internal/entity"
	"avito-internship-2025/internal/repository"
	"github.com/jmoiron/sqlx"
	"log/slog"
)

type CoinService struct {
	DB              *sqlx.DB
	EmployeeRepo    *repository.EmployeeRepository
	TransactionRepo *repository.TransactionRepository
	Logger          *slog.Logger
}

func NewCoinService(db *sqlx.DB, empRepo *repository.EmployeeRepository, txRepo *repository.TransactionRepository, logger *slog.Logger) *CoinService {
	return &CoinService{
		DB:              db,
		EmployeeRepo:    empRepo,
		TransactionRepo: txRepo,
		Logger:          logger,
	}
}

// TransferCoins переводит coins от fromUser к toUser.
func (s *CoinService) TransferCoins(fromUsername, toUsername string, amount int) error {
	if amount <= 0 {
		return errors.New("сумма перевода должна быть положительной")
	}

	tx, err := s.DB.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	fromEmp, err := s.EmployeeRepo.GetEmployeeByUsername(fromUsername)
	if err != nil {
		return err
	}
	if fromEmp == nil {
		return fmt.Errorf("отправитель %s не найден", fromUsername)
	}

	toEmp, err := s.EmployeeRepo.GetEmployeeByUsername(toUsername)
	if err != nil {
		return err
	}
	if toEmp == nil {
		return fmt.Errorf("получатель %s не найден", toUsername)
	}

	if fromEmp.Coins < amount {
		return fmt.Errorf("недостаточно монет: доступно %d, требуется %d", fromEmp.Coins, amount)
	}

	newFromBalance := fromEmp.Coins - amount
	newToBalance := toEmp.Coins + amount

	if err = s.EmployeeRepo.UpdateEmployeeBalance(tx, fromEmp.ID, newFromBalance); err != nil {
		return err
	}
	if err = s.EmployeeRepo.UpdateEmployeeBalance(tx, toEmp.ID, newToBalance); err != nil {
		return err
	}

	counterparty := toUsername
	tran := &entity.Transaction{
		EmployeeID:   fromEmp.ID,
		Type:         "transfer",
		Amount:       amount,
		Counterparty: &counterparty,
	}
	if err = s.TransactionRepo.CreateTransaction(tx, tran); err != nil {
		return err
	}

	s.Logger.Info("Перевод монет выполнен успешно", slog.String("from", fromUsername), slog.String("to", toUsername), slog.Int("amount", amount))
	return nil
}

// BuyMerch функция покупки мерча по его названию
func (s *CoinService) BuyMerch(username, itemName string) error {
	var merch entity.MerchItem
	query := "SELECT id, name, price, created_at FROM merch_items WHERE name = $1"
	if err := s.DB.Get(&merch, query, itemName); err != nil {
		return err
	}

	emp, err := s.EmployeeRepo.GetEmployeeByUsername(username)
	if err != nil {
		return err
	}
	if emp == nil {
		return fmt.Errorf("сотрудник %s не найден", username)
	}
	if emp.Coins < merch.Price {
		return fmt.Errorf("недостаточно монет для покупки %s: требуется %d, доступно %d", merch.Name, merch.Price, emp.Coins)
	}

	tx, err := s.DB.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	newBalance := emp.Coins - merch.Price
	if err = s.EmployeeRepo.UpdateEmployeeBalance(tx, emp.ID, newBalance); err != nil {
		return err
	}

	tran := &entity.Transaction{
		EmployeeID: emp.ID,
		Type:       "purchase",
		Amount:     merch.Price,
		MerchID:    &merch.ID,
	}
	if err = s.TransactionRepo.CreateTransaction(tx, tran); err != nil {
		return err
	}

	s.Logger.Info("Покупка мерча выполнена успешно", slog.String("username", username), slog.String("item", merch.Name))
	return nil
}
