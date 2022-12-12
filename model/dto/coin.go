package dto

type Coin struct {
	CoinID   string  `json:"coin_id"`
	CoinName string  `json:"coin_name"`
	Balance  float64 `json:"balance"`
	Note     string  `json:"note"`
}

type CoinUpdateBalance struct {
	CoinID  string  `validate:"required" json:"coin_id"`
	Balance float64 `validate:"required" json:"balance"`
	Types   string  `validate:"required" json:"type"`
}
