package player_delivery

import (
	"erp-service/helper"
	"erp-service/model/dto"
	"erp-service/usecase/activity_log_usecase"
	"erp-service/usecase/player_usecase"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type PlayerDeliveryImpl struct {
	usecase player_usecase.PlayerUsecase
	log     activity_log_usecase.ActivityLogUsecase
}

func NewPlayerDelivery(logUsecase player_usecase.PlayerUsecase, log activity_log_usecase.ActivityLogUsecase) PlayerDelivery {
	return &PlayerDeliveryImpl{logUsecase, log}
}

func (res *PlayerDeliveryImpl) GetAllPlayer(c *gin.Context) {

	response := res.usecase.GetAllPlayer()
	if response.StatusCode != 200 {
		c.JSON(response.StatusCode, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (res *PlayerDeliveryImpl) AddPlayer(c *gin.Context) {
	playerReq := dto.AddPlayer{}
	if err := c.ShouldBindJSON(&playerReq); err != nil {
		errorRes := helper.ResponseError("Bad Request", err.Error(), 400)
		c.JSON(errorRes.StatusCode, errorRes)
		return
	}

	response := res.usecase.AddPlayer(playerReq)
	if response.StatusCode != 201 {
		c.JSON(response.StatusCode, response)
		return
	}

	userID, _ := c.Get("user_id")
	userName, _ := c.Get("username")

	logBody := dto.ActivityLog{
		UserID:        userID.(string),
		IsTransaction: false,
		Description:   userName.(string) + " telah menambahkan player " + playerReq.PlayerID,
		CreatedAt:     time.Now().Format("2006-01-02 15:04:05"),
	}

	_, errors := res.log.AddActivity(logBody)

	if errors != nil {
		helper.PanicIfError(errors)
	}

	c.JSON(http.StatusCreated, response)
}

func (res *PlayerDeliveryImpl) AddBankPlayer(c *gin.Context) {
	playerReq := dto.AddBankPlayer{}
	if err := c.ShouldBindJSON(&playerReq); err != nil {
		errorRes := helper.ResponseError("Bad Request", err.Error(), 400)
		c.JSON(errorRes.StatusCode, errorRes)
		return
	}

	response := res.usecase.AddPlayerBank(playerReq)
	if response.StatusCode != 201 {
		c.JSON(response.StatusCode, response)
		return
	}

	userID, _ := c.Get("user_id")
	userName, _ := c.Get("username")

	logBody := dto.ActivityLog{
		UserID:        userID.(string),
		IsTransaction: false,
		Description:   userName.(string) + " telah menambahkan bank " + playerReq.BankName + " pada player " + playerReq.PlayerID,
		CreatedAt:     time.Now().Format("2006-01-02 15:04:05"),
	}

	_, errors := res.log.AddActivity(logBody)

	if errors != nil {
		helper.PanicIfError(errors)
	}

	c.JSON(http.StatusCreated, response)
}
