package user_repository

import "erp-service/model/entity"

type UserRepository interface {
	AddUser(*entity.User) (*string, error)
	GetAllUser() ([]entity.UserRole, error)
	GetAllRole() ([]entity.Role, error)
	LoginUsers(username string) (*entity.User, error)
	GetRoleById(id string) (*entity.Role, error)
	GetUserDetail(id string) (*entity.User, *entity.Role, error)
}
