package type_transaction_delivery

import (
	"erp-service/usecase/type_transaction_usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TypeTransactionDeliveryImpl struct {
	usecase type_transaction_usecase.TypeTransactionUsecase
}

func NewTypeTransactionDelivery(typeUsecase type_transaction_usecase.TypeTransactionUsecase) TypeTransactionDelivery {
	return &TypeTransactionDeliveryImpl{usecase: typeUsecase}
}

func (res *TypeTransactionDeliveryImpl) GetAllType(c *gin.Context) {

	response := res.usecase.GetAllType()
	if response.StatusCode != 200 {
		c.JSON(response.StatusCode, response)
		return
	}

	c.JSON(http.StatusOK, response)
}
