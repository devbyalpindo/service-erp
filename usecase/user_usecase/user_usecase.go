package user_usecase

import "erp-service/model/dto"

type UserUsecase interface {
	AddUser(dto.UserAdd) dto.Response
	GetAllUser() dto.Response
	GetAllRole() dto.Response
	LoginUser(dto.UserLogin) dto.Response
	DeleteUsers(id string) dto.Response
}
