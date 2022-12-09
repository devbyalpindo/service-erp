package dto

import "erp-service/model/entity"

type Bonus struct {
	BonusID   string  `json:"bonus_id"`
	Type      string  `json:"type"`
	Ammount   float32 `json:"ammount"`
	CreatedAt string  `json:"created_at"`
}

type BonusWithTotal struct {
	Total      int64          `json:"total"`
	TotalBonus float32        `json:"total_bonus"`
	Bonus      []entity.Bonus `json:"list_bonus"`
}

type BonusAdd struct {
	Type    string  `validate:"required" json:"type"`
	Ammount float32 `validate:"min=0" json:"ammount"`
}
