package coin_delivery

import (
	"erp-service/helper"
	"erp-service/model/dto"
	"erp-service/usecase/activity_log_usecase"
	"erp-service/usecase/coin_usecase"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/paimanbandi/rupiah"
)

type CoinDeliveryImpl struct {
	usecase coin_usecase.CoinUsecase
	log     activity_log_usecase.ActivityLogUsecase
}

func NewCoinDelivery(coinUsecase coin_usecase.CoinUsecase, log activity_log_usecase.ActivityLogUsecase) CoinDelivery {
	return &CoinDeliveryImpl{coinUsecase, log}
}

func (res *CoinDeliveryImpl) GetCoin(c *gin.Context) {

	response := res.usecase.GetCoin()
	if response.StatusCode != 200 {
		c.JSON(response.StatusCode, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (res *CoinDeliveryImpl) GetDetailCoin(c *gin.Context) {

	response := res.usecase.GetDetailCoin()
	if response.StatusCode != 200 {
		c.JSON(response.StatusCode, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (res *CoinDeliveryImpl) UpdateCoinBalance(c *gin.Context) {
	role, _ := c.Get("role")

	coinRequest := dto.CoinUpdateBalance{}
	if err := c.ShouldBindJSON(&coinRequest); err != nil {
		errorRes := helper.ResponseError("Bad Request", err.Error(), 400)
		c.JSON(errorRes.StatusCode, errorRes)
		return
	}

	if coinRequest.Types == "MINUS" && role != "ADMIN" {
		errorRes := helper.ResponseError("Forbidden", "You have no access to do this action", 403)
		c.JSON(errorRes.StatusCode, errorRes)
		return
	}

	response := res.usecase.UpdateCoinBalance(coinRequest)
	if response.StatusCode != 200 {
		c.JSON(response.StatusCode, response)
		return
	}

	userID, _ := c.Get("user_id")
	userName, _ := c.Get("username")

	var logBody dto.ActivityLog

	if coinRequest.Types == "MINUS" {
		logBody = dto.ActivityLog{
			UserID:        userID.(string),
			IsTransaction: false,
			Description:   userName.(string) + " telah mengurangi saldo coin sebesar " + rupiah.FormatFloat64ToRp(coinRequest.Balance),
			CreatedAt:     time.Now().Format("2006-01-02 15:04:05"),
		}
	}

	if coinRequest.Types == "PLUS" {
		logBody = dto.ActivityLog{
			UserID:        userID.(string),
			IsTransaction: false,
			Description:   userName.(string) + " telah menambah saldo coin sebesar " + rupiah.FormatFloat64ToRp(coinRequest.Balance),
			CreatedAt:     time.Now().Format("2006-01-02 15:04:05"),
		}
	}

	_, errors := res.log.AddActivity(logBody)

	if errors != nil {
		helper.PanicIfError(errors)
	}

	c.JSON(http.StatusOK, response)
}
