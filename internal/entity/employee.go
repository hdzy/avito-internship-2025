package entity

import "time"

// Employee представляет сотрудника Avito
type Employee struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Coins     int       `json:"coins"`
	CreatedAt time.Time `json:"created_at"`
}
