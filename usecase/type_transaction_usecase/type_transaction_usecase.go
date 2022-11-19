package type_transaction_usecase

import "erp-service/model/dto"

type TypeTransactionUsecase interface {
	GetAllType() dto.Response
}
