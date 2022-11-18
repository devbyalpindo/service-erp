package dto

type UserLogin struct {
	Username string `validate:"required" json:"username"`
	Password string `validate:"required" json:"password"`
}

type UserAdd struct {
	Username    string `validate:"required" json:"username"`
	Password    string `validate:"required" json:"password"`
	PhoneNumber string `validate:"required" json:"phone_number"`
	RoleID      string `validate:"required" json:"role_id"`
}

type UserRole struct {
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	PhoneNumber string `json:"phone_number"`
	RoleID      string `json:"role_id"`
	RoleName    string `json:"role_name"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
