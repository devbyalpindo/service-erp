package entity

type Transaction struct {
	TransactionID   string  `json:"transaction_id"`
	UserID          string  `json:"user_id"`
	PlayerID        string  `json:"player_id"`
	BankPlayerID    string  `json:"bank_player_id"`
	BankID          string  `json:"bank_id"`
	TypeID          string  `json:"type_id"`
	Ammount         float64 `json:"ammount"`
	AdminFee        float64 `json:"admin_fee"`
	LastBalanceCoin float64 `json:"last_balance_coin"`
	LastBalanceBank float64 `json:"last_balance_bank"`
	Status          string  `json:"status"`
	Note            string  `json:"note"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

type TransactionJoin struct {
	TransactionID       string  `json:"transaction_id"`
	UserID              string  `json:"user_id"`
	PlayerID            string  `json:"player_id"`
	PlayerName          string  `json:"player_name"`
	BankPlayerName      string  `json:"bank_player_name"`
	AccountNumberPlayer string  `json:"account_number_player"`
	BankID              string  `json:"bank_id"`
	BankName            string  `json:"bank_name"`
	AccountNumberBank   string  `json:"account_number_bank"`
	TypeID              string  `json:"type_id"`
	TypeTransaction     string  `json:"type_transaction"`
	Ammount             float64 `json:"ammount"`
	AdminFee            float64 `json:"admin_fee"`
	LastBalanceCoin     float64 `json:"last_balance_coin"`
	LastBalanceBank     float64 `json:"last_balance_bank"`
	Status              string  `json:"status"`
	Note                string  `json:"string"`
	CreatedBy           string  `json:"created_by"`
	CreatedAt           string  `json:"created_at"`
	UpdatedAt           string  `json:"updated_at"`
}
