package entity

import "time"

// MerchItem представляет элемент мерча
type MerchItem struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Price     int       `db:"price" json:"price"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
