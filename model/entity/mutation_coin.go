package entity

type MutationCoin struct {
	MutationCoinID    string  `json:"mutation_coin_id"`
	Type              string  `json:"type"`
	Ammount           float64 `json:"ammount"`
	LastBalance       float64 `json:"last_balance"`
	Description       string  `json:"description"`
	IsTransactionBank bool    `json:"is_transaction_bank"`
	CreatedAt         string  `json:"created_at"`
}

type CoinWithTotal struct {
	Total    int64          `json:"total"`
	Mutation []MutationCoin `json:"mutasi"`
}
