package entity

type Bonus struct {
	BonusID   string  `json:"bonus_id"`
	Type      string  `json:"type"`
	Ammount   float64 `json:"ammount"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}
