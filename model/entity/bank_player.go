package entity

type BankPlayer struct {
	BankPlayerID  string `json:"bank_player_id"`
	PlayerID      string `json:"player_id"`
	BankName      string `json:"bank_name"`
	AccountNumber string `json:"account_number"`
	Category      string `json:"category"`
	CreatedAt     string `json:"created_at"`
}
