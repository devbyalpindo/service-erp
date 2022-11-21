package dashboard_delivery

import (
	"erp-service/usecase/dashboard_usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardDeliveryImpl struct {
	usecase dashboard_usecase.DashboardUsecase
}

func NewDashboardDelivery(dashboardUsecase dashboard_usecase.DashboardUsecase) DashboardDelivery {
	return &DashboardDeliveryImpl{usecase: dashboardUsecase}
}

func (res *DashboardDeliveryImpl) GetDashboard(c *gin.Context) {

	dateFrom := c.Query("dateFrom")
	dateTo := c.Query("dateTo")

	response := res.usecase.GetDashboard(dateFrom, dateTo)
	if response.StatusCode != 200 {
		c.JSON(response.StatusCode, response)
		return
	}

	c.JSON(http.StatusOK, response)
}
