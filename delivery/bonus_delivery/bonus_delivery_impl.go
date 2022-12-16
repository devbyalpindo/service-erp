package bonus_delivery

import (
	"erp-service/helper"
	"erp-service/model/dto"
	"erp-service/usecase/activity_log_usecase"
	"erp-service/usecase/bonus_usecase"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/paimanbandi/rupiah"
)

type BonusDeliveryImpl struct {
	usecase bonus_usecase.BonusUsecase
	log     activity_log_usecase.ActivityLogUsecase
}

func NewBonusDelivery(bonusUsecase bonus_usecase.BonusUsecase, log activity_log_usecase.ActivityLogUsecase) BonusDelivery {
	return &BonusDeliveryImpl{bonusUsecase, log}
}

func (res *BonusDeliveryImpl) AddBonus(c *gin.Context) {
	bonusReq := dto.BonusAdd{}
	if err := c.ShouldBindJSON(&bonusReq); err != nil {
		errorRes := helper.ResponseError("Bad Request", err.Error(), 400)
		c.JSON(errorRes.StatusCode, errorRes)
		return
	}

	response := res.usecase.AddBonus(bonusReq)
	if response.StatusCode != 201 {
		c.JSON(response.StatusCode, response)
		return
	}

	userID, _ := c.Get("user_id")
	userName, _ := c.Get("username")

	logBody := dto.ActivityLog{
		UserID:        userID.(string),
		IsTransaction: false,
		Description:   userName.(string) + " telah menambahkan bonus " + bonusReq.Type + " dengan jumlah " + rupiah.FormatFloat64ToRp(bonusReq.Ammount),
		CreatedAt:     time.Now().Format("2006-01-02 15:04:05"),
	}

	_, errors := res.log.AddActivity(logBody)

	if errors != nil {
		helper.PanicIfError(errors)
	}

	c.JSON(http.StatusCreated, response)
}

func (res *BonusDeliveryImpl) GetAllBonus(c *gin.Context) {
	limit := c.Query("limit")
	offset := c.Query("offset")
	types := c.Query("type")
	limits, _ := strconv.Atoi(limit)
	offsets, _ := strconv.Atoi(offset)

	dateFrom := c.Query("dateFrom")
	dateTo := c.Query("dateTo")

	response := res.usecase.GetAllBonus(types, dateFrom, dateTo, limits, offsets)
	if response.StatusCode != 200 {
		c.JSON(response.StatusCode, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (res *BonusDeliveryImpl) UpdateBonus(c *gin.Context) {
	id := c.Param("id")
	bonusReq := dto.BonusAdd{}
	if err := c.ShouldBindJSON(&bonusReq); err != nil {
		errorRes := helper.ResponseError("Bad Request", err.Error(), 400)
		c.JSON(errorRes.StatusCode, errorRes)
		return
	}

	response := res.usecase.UpdateBonus(id, bonusReq)
	if response.StatusCode != 200 {
		c.JSON(response.StatusCode, response)
		return
	}

	userID, _ := c.Get("user_id")
	userName, _ := c.Get("username")

	logBody := dto.ActivityLog{
		UserID:        userID.(string),
		IsTransaction: false,
		Description:   userName.(string) + " telah merubah bonus dengan id " + id + " dengan jumlah " + rupiah.FormatFloat64ToRp(bonusReq.Ammount),
		CreatedAt:     time.Now().Format("2006-01-02 15:04:05"),
	}

	_, errors := res.log.AddActivity(logBody)

	if errors != nil {
		helper.PanicIfError(errors)
	}

	c.JSON(http.StatusCreated, response)
}
