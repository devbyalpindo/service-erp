package user_delivery

import (
	"erp-service/helper"
	"erp-service/model/dto"
	"erp-service/usecase/activity_log_usecase"
	"erp-service/usecase/user_usecase"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserDeliveryImpl struct {
	usecase user_usecase.UserUsecase
	log     activity_log_usecase.ActivityLogUsecase
}

func NewUserDelivery(userUsecase user_usecase.UserUsecase, log activity_log_usecase.ActivityLogUsecase) UserDelivery {
	return &UserDeliveryImpl{userUsecase, log}
}

func (res *UserDeliveryImpl) AddUser(c *gin.Context) {
	userReq := dto.UserAdd{}
	if err := c.ShouldBindJSON(&userReq); err != nil {
		errorRes := helper.ResponseError("Bad Request", err.Error(), 400)
		c.JSON(errorRes.StatusCode, errorRes)
		return
	}

	response := res.usecase.AddUser(userReq)
	if response.StatusCode != 200 {
		c.JSON(response.StatusCode, response)
		return
	}

	userID, _ := c.Get("user_id")
	userName, _ := c.Get("username")

	logBody := dto.ActivityLog{
		UserID:        userID.(string),
		IsTransaction: true,
		Description:   userName.(string) + " telah menambahkan user " + userReq.Username,
		CreatedAt:     time.Now().Format("2006-01-02 15:04:05"),
	}

	_, errors := res.log.AddActivity(logBody)

	if errors != nil {
		helper.PanicIfError(errors)
	}

	c.JSON(http.StatusCreated, response)
}

func (res *UserDeliveryImpl) GetAllUser(c *gin.Context) {

	response := res.usecase.GetAllUser()
	if response.Status != "ok" {
		c.JSON(response.StatusCode, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (res *UserDeliveryImpl) GetAllRole(c *gin.Context) {

	response := res.usecase.GetAllRole()
	if response.Status != "ok" {
		c.JSON(response.StatusCode, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (res *UserDeliveryImpl) UserLogin(c *gin.Context) {
	userLogin := dto.UserLogin{}
	if err := c.ShouldBindJSON(&userLogin); err != nil {
		errorRes := helper.ResponseError("Bad Request", "Please fill username and password", 400)
		c.JSON(errorRes.StatusCode, errorRes)
		return
	}

	response := res.usecase.LoginUser(userLogin)

	if response.StatusCode != 200 {
		c.JSON(response.StatusCode, response)
		return
	}

	logBody := dto.ActivityLog{
		IsTransaction: false,
		Description:   userLogin.Username + " telah login",
		CreatedAt:     time.Now().Format("2006-01-02 15:04:05"),
	}

	_, errors := res.log.AddActivity(logBody)

	if errors != nil {
		helper.PanicIfError(errors)
	}

	c.JSON(response.StatusCode, response)
}

func (res *UserDeliveryImpl) DeleteUsers(c *gin.Context) {
	id := c.Param("id")
	response := res.usecase.DeleteUsers(id)
	if response.Status != "ok" {
		c.JSON(response.StatusCode, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (res *UserDeliveryImpl) ChangePassword(c *gin.Context) {
	user := dto.UserChangePassword{}
	if err := c.ShouldBindJSON(&user); err != nil {
		errorRes := helper.ResponseError("Bad Request", "Please fill form", 400)
		c.JSON(errorRes.StatusCode, errorRes)
		return
	}
	response := res.usecase.ChangePassword(user)
	if response.Status != "ok" {
		c.JSON(response.StatusCode, response)
		return
	}
	c.JSON(http.StatusOK, response)
}
