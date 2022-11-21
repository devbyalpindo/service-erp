package transaction_delivery

import (
	"erp-service/helper"
	"erp-service/model/dto"
	"erp-service/usecase/activity_log_usecase"
	"erp-service/usecase/transaction_usecase"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TransactionDeliveryImpl struct {
	usecase transaction_usecase.TransactionUsecase
	log     activity_log_usecase.ActivityLogUsecase
}

func NewTransactionDelivery(trx transaction_usecase.TransactionUsecase, log activity_log_usecase.ActivityLogUsecase) TransactionDelivery {
	return &TransactionDeliveryImpl{trx, log}
}

func (res *TransactionDeliveryImpl) AddTransaction(c *gin.Context) {
	trxReq := dto.AddTransaction{}
	userID, _ := c.Get("user_id")
	userName, _ := c.Get("username")

	if err := c.ShouldBindJSON(&trxReq); err != nil {
		errorRes := helper.ResponseError("Bad Request", err.Error(), 400)
		c.JSON(errorRes.StatusCode, errorRes)
		return
	}

	response := res.usecase.AddTransaction(userID.(string), trxReq)
	if response.StatusCode != 201 {
		c.JSON(response.StatusCode, response)
		return
	}

	logBody := dto.ActivityLog{
		UserID:        userID.(string),
		IsTransaction: true,
		TransactionID: response.Data["id"].(string),
		Description:   userName.(string) + " telah melakukan transaksi",
		CreatedAt:     time.Now().Format("2006-01-02 15:04:05"),
	}

	_, errors := res.log.AddActivity(logBody)

	if errors != nil {
		helper.PanicIfError(errors)
	}

	c.JSON(http.StatusCreated, response)
}

func (res *TransactionDeliveryImpl) GetAllTransaction(c *gin.Context) {
	limit := c.Query("limit")
	offset := c.Query("offset")
	limits, _ := strconv.Atoi(limit)
	offsets, _ := strconv.Atoi(offset)

	dateFrom := c.Query("dateFrom")
	dateTo := c.Query("dateTo")

	roleName, _ := c.Get("role")

	response := res.usecase.GetAllTransaction(roleName.(string), limits, offsets, dateFrom, dateTo)
	if response.StatusCode != 200 {
		c.JSON(response.StatusCode, response)
		return
	}

	c.JSON(http.StatusOK, response)
}
