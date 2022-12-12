package dto

import "erp-service/model/entity"

type Dashboard struct {
	TransactionValue  []TransactionValue  `json:"transaction_value"`
	CoinGame          *entity.Coin        `json:"coin"`
	TopPlayerDeposit  []TopPlayerDeposit  `json:"top_player_deposit"`
	TopPlayerWithdraw []TopPlayerWithdraw `json:"top_player_withdraw"`
}

type TransactionValue struct {
	Total        float64 `json:"total"`
	Value        string  `json:"value"`
	TotalAmmount float64 `json:"total_ammount"`
}

type TopPlayerDeposit struct {
	PlayerID     string  `json:"player_id"`
	PlayerName   string  `json:"player_name"`
	TotalDeposit float64 `json:"total_deposit"`
}

type TopPlayerWithdraw struct {
	PlayerID      string  `json:"player_id"`
	PlayerName    string  `json:"player_name"`
	TotalWithdraw float64 `json:"total_withdraw"`
}
