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

type GetMutationCoin struct {
	Type              string `json:"type"`
	IsTransactionBank bool   `json:"is_transaction_bank"`
	DateFrom          string `validate:"required" json:"date_from"`
	DateTo            string `validate:"required" json:"date_to"`
	Limit             int    `validate:"min=0" json:"limit"`
	Offset            int    `validate:"min=0" json:"offset"`
}
