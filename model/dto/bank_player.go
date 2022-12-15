package dto

type BankPlayer struct {
	BankPlayerID  string `json:"bank_player_id"`
	PlayerID      string `json:"player_id"`
	BankName      string `json:"bank_name"`
	AccountName   string `json:"account_name"`
	AccountNumber string `json:"account_number"`
	Category      string `json:"category"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type AddBankPlayer struct {
	PlayerID      string `validate:"required" json:"player_id"`
	BankName      string `validate:"required" json:"bank_name"`
	AccountName   string `validate:"required" json:"account_name"`
	AccountNumber string `validate:"required" json:"account_number"`
}

type UpdateBankPlayer struct {
	BankPlayerID  string `validate:"required" json:"bank_player_id"`
	BankName      string `json:"bank_name"`
	AccountName   string `json:"account_name"`
	AccountNumber string `json:"account_number"`
	Category      string `json:"category"`
}
