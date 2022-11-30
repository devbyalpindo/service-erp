package user_delivery

import "github.com/gin-gonic/gin"

type UserDelivery interface {
	AddUser(*gin.Context)
	GetAllUser(*gin.Context)
	GetAllRole(*gin.Context)
	UserLogin(*gin.Context)
	DeleteUsers(*gin.Context)
}
