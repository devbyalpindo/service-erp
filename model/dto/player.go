package dto

type Player struct {
	PlayerID   string `json:"player_id"`
	PlayerName string `json:"player_name"`
	Recid      string `json:"recid"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type PlayerBankPlayer struct {
	PlayerID       string       `json:"player_id"`
	PlayerName     string       `json:"player_name"`
	CreatedAt      string       `json:"created_at"`
	UpdatedAt      string       `json:"updated_at"`
	ListBankPlayer []BankPlayer `json:"list_bank_player"`
}

type AddPlayer struct {
	PlayerID      string `validate:"required" json:"player_id"`
	PlayerName    string `validate:"required" json:"player_name"`
	BankName      string `validate:"required" json:"bank_name"`
	AccountName   string `validate:"required" json:"account_name"`
	AccountNumber string `validate:"required" json:"account_number"`
	Category      string `validate:"required" json:"category"`
}

type UpdatePlayer struct {
	PlayerID   string `validate:"required" json:"player_id"`
	PlayerName string `validate:"required" json:"player_name"`
}

type BulkInsertPlayer struct {
	Username         string `validate:"required" json:"username"`
	FullName         string `json:"full_Name"`
	Recid            string `json:"recid"`
	RegistrationDate string `json:"registration_Date"`
}

type BulkInsertBankPlayer struct {
	PlayerID      string `validate:"required" json:"player_id"`
	BankName      string `validate:"required" json:"bank_Name"`
	AccountName   string `validate:"required" json:"bank_Account_Name"`
	AccountNumber string `validate:"required" json:"bank_Account_Number"`
}
