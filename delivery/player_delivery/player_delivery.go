package player_delivery

import "github.com/gin-gonic/gin"

type PlayerDelivery interface {
	GetAllPlayer(*gin.Context)
	AddPlayer(*gin.Context)
	AddBankPlayer(*gin.Context)
	UpdatePlayer(*gin.Context)
	UpdateBankPlayer(*gin.Context)
}
