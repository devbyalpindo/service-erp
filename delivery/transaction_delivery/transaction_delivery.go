package transaction_delivery

import "github.com/gin-gonic/gin"

type TransactionDelivery interface {
	AddTransaction(*gin.Context)
	GetAllTransaction(*gin.Context)
	UpdateTransaction(*gin.Context)
	CanceledTransaction(*gin.Context)
}
