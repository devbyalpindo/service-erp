package dto

type Player struct {
	PlayerID   string `json:"player_id"`
	PlayerName string `json:"player_name"`
	CreatedAt  string `json:"created_at"`
}

type PlayerBankPlayer struct {
	PlayerID       string       `json:"player_id"`
	PlayerName     string       `json:"player_name"`
	CreatedAt      string       `json:"created_at"`
	ListBankPlayer []BankPlayer `json:"list_bank_player"`
}

type AddPlayer struct {
	PlayerID      string `validate:"required" json:"player_id"`
	PlayerName    string `validate:"required" json:"player_name"`
	BankName      string `validate:"required" json:"bank_name"`
	AccountNumber string `validate:"required" json:"account_number"`
	Category      string `validate:"required" json:"category"`
}
