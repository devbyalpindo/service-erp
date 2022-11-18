package dto

import (
	"erp-service/model/entity"
)

type ActivityLog struct {
	UserID        string `json:"user_id"`
	IsTransaction bool   `json:"is_transaction"`
	TransactionID string `json:"transaction_id"`
	Description   string `json:"description"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type ActivityLogWithTotal struct {
	Total       int64 `json:"total"`
	ActivityLog []entity.ActivityLog
}
