package user_repository

import (
	"erp-service/helper"
	"erp-service/model/entity"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) UserRepository {
	return &UserRepositoryImpl{DB: DB}
}

func (repository *UserRepositoryImpl) AddUser(user *entity.User) (*string, error) {
	if err := repository.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user.UserID, nil
}

func (repository *UserRepositoryImpl) LoginUsers(username string) (*entity.User, error) {
	user := entity.User{}
	result := repository.DB.Where("username = ? AND status = ?", username, "ACTIVE").Find(&user)

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &user, nil
}

func (repository *UserRepositoryImpl) GetAllUser() ([]entity.UserRole, error) {
	user := []entity.UserRole{}
	err := repository.DB.Table("users").Select("users.user_id, users.username, users.phone_number, users.role_id, roles.role_name, users.created_at, users.updated_at").Joins("inner join roles on roles.role_id = users.role_id").Where("users.status = ?", "ACTIVE").Find(&user).Error
	helper.PanicIfError(err)
	if len(user) <= 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return user, nil
}

func (repository *UserRepositoryImpl) GetAllRole() ([]entity.Role, error) {
	role := []entity.Role{}
	err := repository.DB.Model(&entity.Role{}).Scan(&role).Error
	helper.PanicIfError(err)
	if len(role) <= 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return role, nil
}

func (repository *UserRepositoryImpl) GetRoleById(id string) (*entity.Role, error) {
	role := entity.Role{}
	result := repository.DB.Where("role_id = ?", id).Find(&role)

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &role, nil
}

func (repository *UserRepositoryImpl) GetUserDetail(id string) (*entity.User, *entity.Role, error) {
	user := entity.User{}
	result := repository.DB.Where("user_id = ? AND status = ?", id, "ACTIVE").Find(&user)

	if result.RowsAffected == 0 {
		return nil, nil, gorm.ErrRecordNotFound
	}

	role := entity.Role{}
	results := repository.DB.Where("role_id = ?", user.RoleID).Find(&role)

	if results.RowsAffected == 0 {
		return nil, nil, gorm.ErrRecordNotFound
	}

	return &user, &role, nil
}

func (repository *UserRepositoryImpl) CheckExistUser(username string) bool {
	user := entity.User{}
	result := repository.DB.Where("username = ?", username).Find(&user)

	return result.RowsAffected == 0
}

func (repository *UserRepositoryImpl) DeleteUsers(id string) (*string, error) {
	result := repository.DB.Model(&entity.User{}).Where("user_id = ?", id).Update("status", "INACTIVE")

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &id, nil
}

func (repository *UserRepositoryImpl) ChangePassword(id string, password string) (*string, error) {
	result := repository.DB.Model(&entity.User{}).Where("user_id = ?", id).Update("password", password)

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &id, nil
}
