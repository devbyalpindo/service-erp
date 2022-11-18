package dto

import "erp-service/model/entity"

type Transaction struct {
	TransactionID   string  `json:"transaction_id"`
	UserID          string  `json:"user_id"`
	PlayerName      string  `json:"player_name"`
	PlayerID        string  `json:"player_id"`
	BankPlayer      string  `json:"bank_player"`
	AccountNumber   string  `json:"account_number"`
	BankID          string  `json:"bank_id"`
	TypeID          string  `json:"type_id"`
	Ammount         float32 `json:"ammount"`
	AdminFee        float32 `json:"admin_fee"`
	LastBalanceCoin float32 `json:"last_balance_coin"`
	LastBalanceBank float32 `json:"last_balance_bank"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}
type AddTransaction struct {
	PlayerName    string  `validate:"required" json:"player_name"`
	PlayerID      string  `validate:"required" json:"player_id"`
	BankPlayer    string  `validate:"required" json:"bank_player"`
	AccountNumber string  `validate:"required" json:"account_number"`
	BankID        string  `validate:"required" json:"bank_id"`
	TypeID        string  `validate:"required" json:"type_id"`
	Ammount       float32 `validate:"min=0" json:"ammount"`
	AdminFee      float32 `validate:"min=0" json:"admin_fee"`
}

type TransactionWithTotal struct {
	Total       int64 `json:"total"`
	Transaction []entity.TransactionJoin
}

type TransactionJoin struct {
	TransactionID     string  `json:"transaction_id"`
	UserID            string  `json:"user_id"`
	PlayerName        string  `json:"player_name"`
	PlayerID          string  `json:"player_id"`
	BankPlayer        string  `json:"bank_player"`
	AccountNumber     string  `json:"account_number"`
	BankID            string  `json:"bank_id"`
	BankName          string  `json:"bank_name"`
	AccountNumberBank string  `json:"account_number_bank"`
	TypeID            string  `json:"type_id"`
	TypeTransaction   string  `json:"type_transaction"`
	Ammount           float32 `json:"ammount"`
	AdminFee          float32 `json:"admin_fee"`
	LastBalanceCoin   float32 `json:"last_balance_coin"`
	LastBalanceBank   float32 `json:"last_balance_bank"`
	CreatedBy         string  `json:"created_by"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}