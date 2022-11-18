package dto

type Bank struct {
	BankID        string  `json:"bank_id"`
	BankName      string  `json:"bank_name"`
	AccountName   string  `json:"account_name"`
	Category      string  `json:"category"`
	AccountNumber string  `json:"account_number"`
	Balance       float32 `json:"balance"`
	Active        bool    `json:"active"`
	Ibanking      string  `json:"ibanking"`
	CodeAccess    string  `json:"code_access"`
	Pin           string  `json:"pin"`
}

type BankAdd struct {
	BankName      string  `validate:"required" json:"bank_name"`
	AccountName   string  `validate:"required" json:"account_name"`
	Category      string  `validate:"required" json:"category"`
	AccountNumber string  `validate:"required" json:"account_number"`
	Balance       float32 `validate:"required" json:"balance"`
	Active        bool    `validate:"required" json:"active"`
	Ibanking      string  `json:"ibanking"`
	CodeAccess    string  `json:"code_access"`
	Pin           string  `json:"pin"`
}

type BankUpdate struct {
	BankName      string `validate:"required" json:"bank_name"`
	AccountName   string `validate:"required" json:"account_name"`
	Category      string `json:"category"`
	AccountNumber string `validate:"required" json:"account_number"`
	Active        bool   `validate:"required" json:"active"`
	Ibanking      string `json:"ibanking"`
	CodeAccess    string `json:"code_access"`
	Pin           string `json:"pin"`
}

type BankUpdateBalance struct {
	BankID  string  `validate:"required" json:"bank_id"`
	Types   string  `validate:"required" json:"type"`
	Balance float32 `validate:"required" json:"balance"`
}
