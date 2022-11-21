package entity

type Bank struct {
	BankID        string  `json:"bank_id"`
	BankName      string  `json:"bank_name"`
	AccountName   string  `json:"account_name"`
	Category      string  `json:"category"`
	AccountNumber string  `json:"account_number"`
	Balance       float32 `json:"balance"`
	Active        *bool   `json:"active"`
	Ibanking      string  `json:"ibanking"`
	CodeAccess    string  `json:"code_access"`
	Pin           string  `json:"pin"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}
