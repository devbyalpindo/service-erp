package bank_delivery

import "github.com/gin-gonic/gin"

type BankDelivery interface {
	AddBank(*gin.Context)
	GetAllBank(*gin.Context)
	UpdateBank(*gin.Context)
	UpdateBankBalance(*gin.Context)
}
