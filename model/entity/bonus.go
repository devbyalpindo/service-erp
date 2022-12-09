package entity

type Bonus struct {
	BonusID   string  `json:"bonus_id"`
	Type      string  `json:"type"`
	Ammount   float32 `json:"ammount"`
	CreatedAt string  `json:"created_at"`
}
