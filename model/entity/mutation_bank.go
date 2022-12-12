package entity

type MutationBank struct {
	MutationBankID string  `json:"mutation_bank_id"`
	BankID         string  `json:"bank_id"`
	Type           string  `json:"type"`
	Ammount        float32 `json:"ammount"`
	LastBalance    float32 `json:"last_balance"`
	Description    string  `json:"description"`
	CreatedAt      string  `json:"created_at"`
}

type BankJoinMutation struct {
	MutationBankID string  `json:"mutation_bank_id"`
	BankID         string  `json:"bank_id"`
	BankName       string  `json:"bank_name"`
	AccountNumber  string  `json:"account_number"`
	Type           string  `json:"type"`
	Ammount        float32 `json:"ammount"`
	LastBalance    float32 `json:"last_balance"`
	Description    string  `json:"description"`
}

type BankithTotal struct {
	Total    int64              `json:"total"`
	Mutation []BankJoinMutation `json:"mutasi"`
}
