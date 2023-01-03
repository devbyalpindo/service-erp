package entity

type MutationBank struct {
	MutationBankID    string  `json:"mutation_bank_id"`
	BankID            string  `json:"bank_id"`
	Type              string  `json:"type"`
	Ammount           float64 `json:"ammount"`
	LastBalance       float64 `json:"last_balance"`
	Description       string  `json:"description"`
	IsTransactionBank bool    `json:"is_transaction_bank"`
	CreatedAt         string  `json:"created_at"`
}

type BankJoinMutation struct {
	MutationBankID string  `json:"mutation_bank_id"`
	BankID         string  `json:"bank_id"`
	BankName       string  `json:"bank_name"`
	AccountNumber  string  `json:"account_number"`
	Type           string  `json:"type"`
	Ammount        float64 `json:"ammount"`
	LastBalance    float64 `json:"last_balance"`
	Description    string  `json:"description"`
	CreatedAt      string  `json:"created_at"`
}

type BankithTotal struct {
	Total    int64              `json:"total"`
	Mutation []BankJoinMutation `json:"mutasi"`
}
