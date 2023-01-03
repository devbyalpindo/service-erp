package dto

type Bank struct {
	BankID        string  `json:"bank_id"`
	BankName      string  `json:"bank_name"`
	AccountName   string  `json:"account_name"`
	Category      string  `json:"category"`
	AccountNumber string  `json:"account_number"`
	Balance       float64 `json:"balance"`
	Active        *bool   `json:"active"`
	Ibanking      string  `json:"ibanking"`
	CodeAccess    string  `json:"code_access"`
	Pin           string  `json:"pin"`
}

type BankAdd struct {
	BankName      string  `validate:"required" json:"bank_name"`
	AccountName   string  `validate:"required" json:"account_name"`
	Category      string  `validate:"required" json:"category"`
	AccountNumber string  `validate:"required" json:"account_number"`
	Balance       float64 `validate:"required" json:"balance"`
	Active        *bool   `validate:"required" json:"active"`
	Ibanking      string  `json:"ibanking"`
	CodeAccess    string  `json:"code_access"`
	Pin           string  `json:"pin"`
}

type BankUpdate struct {
	BankName      string `validate:"required" json:"bank_name"`
	AccountName   string `validate:"required" json:"account_name"`
	Category      string `json:"category"`
	AccountNumber string `validate:"required" json:"account_number"`
	Active        *bool  `validate:"required" json:"active"`
	Ibanking      string `json:"ibanking"`
	CodeAccess    string `json:"code_access"`
	Pin           string `json:"pin"`
}

type BankUpdateBalance struct {
	BankID  string  `validate:"required" json:"bank_id"`
	Types   string  `validate:"required" json:"type"`
	Balance float64 `validate:"required" json:"balance"`
}

type BankTransfer struct {
	FromBankID string  `validate:"required" json:"from_bank_id"`
	ToBankID   string  `validate:"required" json:"to_bank_id"`
	Balance    float64 `validate:"min=0" json:"balance"`
	AdminFee   float64 `validate:"min=0" json:"admin_fee"`
}

type GetMutationBank struct {
	BankID            string `json:"bank_id"`
	Type              string `json:"type"`
	IsTransactionBank bool   `json:"is_transaction_bank"`
	DateFrom          string `validate:"required" json:"date_from"`
	DateTo            string `validate:"required" json:"date_to"`
	Limit             int    `validate:"min=0" json:"limit"`
	Offset            int    `validate:"min=0" json:"offset"`
}
