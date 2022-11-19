package type_transaction_usecase

import (
	"erp-service/helper"
	"erp-service/model/dto"
	"erp-service/repository/type_repository"
	"errors"

	"gorm.io/gorm"
)

type TypeTransactionUsecaseImpl struct {
	TypeRepository type_repository.TypeRepository
}

func NewTypeTransactionUsecase(typeRepository type_repository.TypeRepository) TypeTransactionUsecase {
	return &TypeTransactionUsecaseImpl{
		TypeRepository: typeRepository,
	}
}

func (usecase *TypeTransactionUsecaseImpl) GetAllType() dto.Response {
	typeList, err := usecase.TypeRepository.GetAllType()
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return helper.ResponseError("failed", "Data not found", 404)
	} else if err != nil {
		return helper.ResponseError("failed", err, 500)
	}
	helper.PanicIfError(err)
	response := []dto.TypeTransaction{}
	for _, coin := range typeList {
		responseData := dto.TypeTransaction{
			TypeID:          coin.TypeID,
			TypeTransaction: coin.TypeTransaction,
		}
		response = append(response, responseData)
	}

	var result map[string]any = make(map[string]any)
	result["list_type"] = response

	return helper.ResponseSuccess("ok", nil, result, 200)
}
