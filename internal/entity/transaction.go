package entity

import "time"

// Transaction представляет транзакцию
type Transaction struct {
	ID           int       `db:"id" json:"id"`
	EmployeeID   int       `db:"employee_id" json:"employee_id"`
	Type         string    `db:"type" json:"type"` // transfer / purchase
	Amount       int       `db:"amount" json:"amount"`
	MerchID      *int      `db:"merch_id" json:"merch_id,omitempty"`         // при покупке
	Counterparty *string   `db:"counterparty" json:"counterparty,omitempty"` // при переводе
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}
