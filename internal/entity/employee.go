package entity

import "time"

// Employee представляет сотрудника Avito
type Employee struct {
	ID           int       `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Coins        int       `db:"coins" json:"coins"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}
