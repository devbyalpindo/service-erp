package type_transaction_delivery

import "github.com/gin-gonic/gin"

type TypeTransactionDelivery interface {
	GetAllType(*gin.Context)
}
