package bank_delivery

import (
	"erp-service/helper"
	"erp-service/model/dto"
	"erp-service/usecase/activity_log_usecase"
	"erp-service/usecase/bank_usecase"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/paimanbandi/rupiah"
)

type BankDeliveryImpl struct {
	usecase bank_usecase.BankUsecase
	log     activity_log_usecase.ActivityLogUsecase
}

func NewBankDelivery(bankUsecase bank_usecase.BankUsecase, log activity_log_usecase.ActivityLogUsecase) BankDelivery {
	return &BankDeliveryImpl{bankUsecase, log}
}

func (res *BankDeliveryImpl) AddBank(c *gin.Context) {
	bankReq := dto.BankAdd{}
	if err := c.ShouldBindJSON(&bankReq); err != nil {
		errorRes := helper.ResponseError("Bad Request", err.Error(), 400)
		c.JSON(errorRes.StatusCode, errorRes)
		return
	}

	response := res.usecase.AddBank(bankReq)
	if response.StatusCode != 201 {
		c.JSON(response.StatusCode, response)
		return
	}

	userID, _ := c.Get("user_id")
	userName, _ := c.Get("username")

	logBody := dto.ActivityLog{
		UserID:        userID.(string),
		IsTransaction: false,
		Description:   userName.(string) + " telah menambahkan akun bank " + bankReq.BankName + " dengan account number " + bankReq.AccountNumber,
		CreatedAt:     time.Now().Format("2006-01-02 15:04:05"),
	}

	_, errors := res.log.AddActivity(logBody)

	if errors != nil {
		helper.PanicIfError(errors)
	}

	c.JSON(http.StatusCreated, response)
}

func (res *BankDeliveryImpl) GetAllBank(c *gin.Context) {

	response := res.usecase.GetAllBank()
	if response.StatusCode != 200 {
		c.JSON(response.StatusCode, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (res *BankDeliveryImpl) UpdateBank(c *gin.Context) {
	id := c.Param("id")

	bankRequest := dto.BankUpdate{}
	if err := c.ShouldBindJSON(&bankRequest); err != nil {
		errorRes := helper.ResponseError("Bad Request", err.Error(), 400)
		c.JSON(errorRes.StatusCode, errorRes)
		return
	}

	response := res.usecase.UpdateBank(id, bankRequest)
	if response.StatusCode != 200 {
		c.JSON(response.StatusCode, response)
		return
	}

	userID, _ := c.Get("user_id")
	userName, _ := c.Get("username")

	logBody := dto.ActivityLog{
		UserID:        userID.(string),
		IsTransaction: false,
		Description:   userName.(string) + " telah merubah akun bank " + bankRequest.BankName + " dengan account number " + bankRequest.AccountNumber,
		CreatedAt:     time.Now().Format("2006-01-02 15:04:05"),
	}

	_, errors := res.log.AddActivity(logBody)

	if errors != nil {
		helper.PanicIfError(errors)
	}

	c.JSON(http.StatusOK, response)
}

func (res *BankDeliveryImpl) UpdateBankBalance(c *gin.Context) {
	role, _ := c.Get("role")
	bankRequest := dto.BankUpdateBalance{}
	if err := c.ShouldBindJSON(&bankRequest); err != nil {
		errorRes := helper.ResponseError("Bad Request", err.Error(), 400)
		c.JSON(errorRes.StatusCode, errorRes)
		return
	}

	if bankRequest.Types == "MINUS" && role != "ADMIN" {
		errorRes := helper.ResponseError("Forbidden", "You have no access to do this action", 403)
		c.JSON(errorRes.StatusCode, errorRes)
		return
	}

	response := res.usecase.UpdateBankBalance(bankRequest)
	if response.StatusCode != 200 {
		c.JSON(response.StatusCode, response)
		return
	}
	userID, _ := c.Get("user_id")
	userName, _ := c.Get("username")

	var logBody dto.ActivityLog

	if bankRequest.Types == "MINUS" {
		logBody = dto.ActivityLog{
			UserID:        userID.(string),
			IsTransaction: false,
			Description:   userName.(string) + " telah mengurangi saldo bank " + bankRequest.BankID + " sebesar " + rupiah.FormatFloat64ToRp(bankRequest.Balance),
			CreatedAt:     time.Now().Format("2006-01-02 15:04:05"),
		}
	}

	if bankRequest.Types == "PLUS" {
		logBody = dto.ActivityLog{
			UserID:        userID.(string),
			IsTransaction: false,
			Description:   userName.(string) + " telah menambah saldo bank " + bankRequest.BankID + " sebesar " + rupiah.FormatFloat64ToRp(bankRequest.Balance),
			CreatedAt:     time.Now().Format("2006-01-02 15:04:05"),
		}
	}

	_, errors := res.log.AddActivity(logBody)

	if errors != nil {
		helper.PanicIfError(errors)
	}

	c.JSON(http.StatusOK, response)
}

func (res *BankDeliveryImpl) TransferToBank(c *gin.Context) {

	bankRequest := dto.BankTransfer{}
	if err := c.ShouldBindJSON(&bankRequest); err != nil {
		errorRes := helper.ResponseError("Bad Request", err.Error(), 400)
		c.JSON(errorRes.StatusCode, errorRes)
		return
	}

	response := res.usecase.TransferToBank(bankRequest)
	if response.StatusCode != 200 {
		c.JSON(response.StatusCode, response)
		return
	}
	userID, _ := c.Get("user_id")
	userName, _ := c.Get("username")

	logBody := dto.ActivityLog{
		UserID:        userID.(string),
		IsTransaction: false,
		Description:   userName.(string) + " telah transfer saldo bank dari bank " + bankRequest.FromBankID + " ke bank " + bankRequest.ToBankID + " sebesar " + rupiah.FormatFloat64ToRp(bankRequest.Balance),
		CreatedAt:     time.Now().Format("2006-01-02 15:04:05"),
	}

	_, errors := res.log.AddActivity(logBody)

	if errors != nil {
		helper.PanicIfError(errors)
	}

	c.JSON(http.StatusOK, response)
}

func (res *BankDeliveryImpl) GetMutation(c *gin.Context) {
	mutationRequest := dto.GetMutationBank{}
	if err := c.ShouldBindJSON(&mutationRequest); err != nil {
		errorRes := helper.ResponseError("Bad Request", err.Error(), 400)
		c.JSON(errorRes.StatusCode, errorRes)
		return
	}

	response := res.usecase.GetMutation(mutationRequest)
	if response.StatusCode != 200 {
		c.JSON(response.StatusCode, response)
		return
	}

	c.JSON(http.StatusOK, response)
}
