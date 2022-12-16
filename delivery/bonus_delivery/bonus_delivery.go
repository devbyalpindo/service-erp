package bonus_delivery

import "github.com/gin-gonic/gin"

type BonusDelivery interface {
	AddBonus(*gin.Context)
	GetAllBonus(*gin.Context)
	UpdateBonus(*gin.Context)
}
