package entity

type Coin struct {
	CoinID    string  `json:"coin_id"`
	CoinName  string  `json:"coin_name"`
	Balance   float64 `json:"balance"`
	Note      string  `json:"note"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}
