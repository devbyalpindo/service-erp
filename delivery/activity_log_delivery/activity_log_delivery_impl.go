package activity_log_delivery

import (
	"erp-service/usecase/activity_log_usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ActivityLogDeliveryImpl struct {
	usecase activity_log_usecase.ActivityLogUsecase
}

func NewActivityLogDelivery(logUsecase activity_log_usecase.ActivityLogUsecase) ActivityLogDelivery {
	return &ActivityLogDeliveryImpl{logUsecase}
}

func (res *ActivityLogDeliveryImpl) GetActivity(c *gin.Context) {
	limit := c.Query("limit")
	offset := c.Query("offset")
	limits, _ := strconv.Atoi(limit)
	offsets, _ := strconv.Atoi(offset)

	response := res.usecase.GetActivity(limits, offsets)
	if response.StatusCode != 200 {
		c.JSON(response.StatusCode, response)
		return
	}

	c.JSON(http.StatusOK, response)
}
