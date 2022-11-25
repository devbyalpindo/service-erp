package entity

type Player struct {
	PlayerID   string       `json:"player_id"`
	PlayerName string       `json:"player_name"`
	BankPlayer []BankPlayer `gorm:"foreignKey:PlayerID;references:PlayerID" json:"bank_player"`
	CreatedAt  string       `json:"created_at"`
}

type PlayerBankPlayer struct {
	PlayerID       string       `json:"player_id"`
	PlayerName     string       `json:"player_name"`
	CreatedAt      string       `json:"created_at"`
	ListBankPlayer []BankPlayer `gorm:"foreignKey:PlayerID"`
}
