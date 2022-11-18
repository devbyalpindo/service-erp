package activity_log_delivery

import "github.com/gin-gonic/gin"

type ActivityLogDelivery interface {
	GetActivity(*gin.Context)
}
