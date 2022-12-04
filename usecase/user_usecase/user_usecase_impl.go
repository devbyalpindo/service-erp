package user_usecase

import (
	"erp-service/helper"
	"erp-service/model/dto"
	"erp-service/model/entity"
	"erp-service/repository/user_repository"
	"erp-service/usecase/jwt_usecase"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserUsecaseImpl struct {
	UserRepository user_repository.UserRepository
	JwtUsecase     jwt_usecase.JwtUsecase
	Validate       *validator.Validate
}

func NewUserUsecase(userRepository user_repository.UserRepository, jwtUsecase jwt_usecase.JwtUsecase, validate *validator.Validate) UserUsecase {
	return &UserUsecaseImpl{
		UserRepository: userRepository,
		JwtUsecase:     jwtUsecase,
		Validate:       validate,
	}
}

func (usecase *UserUsecaseImpl) AddUser(body dto.UserAdd) dto.Response {
	err := usecase.Validate.Struct(body)

	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]dto.CustomMessageError, len(ve))
			for i, fe := range ve {
				out[i] = dto.CustomMessageError{
					Field:    fe.Field(),
					Messsage: helper.MessageError(fe.Tag()),
				}
			}
			return helper.ResponseError("failed", out, 400)
		}

	}

	checkUser := usecase.UserRepository.CheckExistUser(body.Username)

	if !checkUser {
		return helper.ResponseError("failed", "Username telah terdaftar", 400)
	}

	checkRole, err := usecase.UserRepository.GetRoleById(body.RoleID)

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return helper.ResponseError("failed", "Role tidak ditemukan", 404)
	}

	createID := uuid.New().String()
	encryptPwd, err := helper.HashPassword(body.Password)
	helper.PanicIfError(err)

	payloadUser := &entity.User{
		UserID:      createID,
		Username:    body.Username,
		Password:    encryptPwd,
		PhoneNumber: body.PhoneNumber,
		RoleID:      checkRole.RoleID,
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt:   time.Now().Format("2006-01-02 15:04:05"),
	}

	user, err := usecase.UserRepository.AddUser(payloadUser)
	helper.PanicIfError(err)

	return helper.ResponseSuccess("ok", nil, map[string]interface{}{"id": user}, 201)

}

func (usecase *UserUsecaseImpl) GetAllUser() dto.Response {
	userList, err := usecase.UserRepository.GetAllUser()
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return helper.ResponseError("Data not found", err, 404)
	} else if err != nil {
		return helper.ResponseError("Internal server error", err, 500)
	}
	helper.PanicIfError(err)

	response := []dto.UserRole{}

	for _, user := range userList {
		timeCreated, _ := time.Parse(time.RFC3339, user.CreatedAt)
		responseData := dto.UserRole{
			UserID:      user.UserID,
			Username:    user.Username,
			PhoneNumber: user.PhoneNumber,
			RoleID:      user.RoleID,
			RoleName:    user.RoleName,
			CreatedAt:   timeCreated.Format("2006-01-02 15:04:05"),
			UpdatedAt:   user.UpdatedAt,
		}
		response = append(response, responseData)
	}

	var result map[string]any = make(map[string]any)
	result["list_user"] = response

	return helper.ResponseSuccess("ok", nil, result, 200)
}

func (usecase *UserUsecaseImpl) GetAllRole() dto.Response {
	roleList, err := usecase.UserRepository.GetAllRole()
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return helper.ResponseError("Data not found", err, 404)
	} else if err != nil {
		return helper.ResponseError("Internal server error", err, 500)
	}
	helper.PanicIfError(err)

	response := []dto.Role{}

	for _, role := range roleList {
		responseData := dto.Role{
			RoleID:   role.RoleID,
			RoleName: role.RoleName,
		}
		response = append(response, responseData)
	}

	var result map[string]any = make(map[string]any)
	result["list_role"] = response

	return helper.ResponseSuccess("ok", nil, result, 200)

}

func (usecase *UserUsecaseImpl) LoginUser(userPayload dto.UserLogin) dto.Response {
	user, err := usecase.UserRepository.LoginUsers(userPayload.Username)

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return helper.ResponseError("failed", "Username or password incorrect", 404)
	}

	errPassword := helper.CheckPasswordHash(userPayload.Password, user.Password)

	if errPassword != nil {
		return helper.ResponseError("failed", "Username or password incorrect", 400)
	}

	jwt, err := usecase.JwtUsecase.GenerateToken(user.UserID, user.Username, user.RoleID)

	if err != nil {
		return helper.ResponseError("failed", "Wrong personal number / password", 404)
	}

	return helper.ResponseSuccess("ok", nil, map[string]interface{}{"user_id": user.UserID, "token": jwt}, 200)

}

func (usecase *UserUsecaseImpl) DeleteUsers(id string) dto.Response {
	user, err := usecase.UserRepository.DeleteUsers(id)

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return helper.ResponseError("failed", "User not found", 404)
	}

	return helper.ResponseSuccess("ok", nil, map[string]interface{}{"user_id": user}, 200)

}

func (usecase *UserUsecaseImpl) ChangePassword(body dto.UserChangePassword) dto.Response {
	user, _, err := usecase.UserRepository.GetUserDetail(body.UserID)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return helper.ResponseError("failed", "User not found", 404)
	}

	errPassword := helper.CheckPasswordHash(body.OldPassword, user.Password)

	if errPassword != nil {
		return helper.ResponseError("failed", "Password incorrect", 400)
	}

	encryptPwd, err := helper.HashPassword(body.NewPassword)
	helper.PanicIfError(err)

	id, err := usecase.UserRepository.ChangePassword(body.UserID, encryptPwd)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return helper.ResponseError("failed", "Failed change password", 404)
	}

	return helper.ResponseSuccess("ok", nil, map[string]interface{}{"user_id": id}, 200)
}
