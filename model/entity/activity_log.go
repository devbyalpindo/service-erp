package entity

type ActivityLog struct {
	ActivityID    string `json:"activity_id"`
	UserID        string `json:"user_id"`
	IsTransaction bool   `json:"is_transaction"`
	TransactionID string `json:"transaction_id"`
	Description   string `json:"description"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}
