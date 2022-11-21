package dashboard_delivery

import "github.com/gin-gonic/gin"

type DashboardDelivery interface {
	GetDashboard(*gin.Context)
}
