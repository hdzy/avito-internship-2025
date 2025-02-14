package entity

// InventoryItem элемент купленного мерча
type InventoryItem struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

// TransactionHistoryItem элемент истории переводов
type TransactionHistoryItem struct {
	FromUser string `json:"fromUser,omitempty"`
	ToUser   string `json:"toUser,omitempty"`
	Amount   int    `json:"amount"`
}

// CoinHistory история переводов
type CoinHistory struct {
	Received []TransactionHistoryItem `json:"received"`
	Sent     []TransactionHistoryItem `json:"sent"`
}

// InfoResponse /api/info
type InfoResponse struct {
	Coins       int             `json:"coins"`
	Inventory   []InventoryItem `json:"inventory"`
	CoinHistory CoinHistory     `json:"coinHistory"`
}
