package bank_usecase

import (
	"erp-service/helper"
	"erp-service/model/dto"
	"erp-service/model/entity"
	"erp-service/repository/bank_repository"
	"errors"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BankUsecaseImpl struct {
	BankRepository bank_repository.BankRepository
	Validate       *validator.Validate
}

func NewBankUsecase(bankRepository bank_repository.BankRepository, validate *validator.Validate) BankUsecase {
	return &BankUsecaseImpl{
		BankRepository: bankRepository,
		Validate:       validate,
	}
}

func (usecase *BankUsecaseImpl) AddBank(body dto.BankAdd) dto.Response {
	err := usecase.Validate.Struct(body)

	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]dto.CustomMessageError, len(ve))
			for i, fe := range ve {
				out[i] = dto.CustomMessageError{
					Field:    fe.Field(),
					Messsage: helper.MessageError(fe.Tag()),
				}
			}
			return helper.ResponseError("failed", out, 400)
		}

	}

	createID := uuid.New().String()
	helper.PanicIfError(err)

	payloadBank := &entity.Bank{
		BankID:        createID,
		BankName:      body.BankName,
		AccountName:   body.AccountName,
		Category:      body.Category,
		AccountNumber: body.AccountNumber,
		Balance:       body.Balance,
		Active:        true,
		Ibanking:      body.Ibanking,
		CodeAccess:    body.CodeAccess,
		Pin:           body.Pin,
		CreatedAt:     time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt:     time.Now().Format("2006-01-02 15:04:05"),
	}

	bank, err := usecase.BankRepository.AddBank(payloadBank)
	helper.PanicIfError(err)

	return helper.ResponseSuccess("ok", nil, map[string]interface{}{"id": bank}, 201)

}

func (usecase *BankUsecaseImpl) GetAllBank() dto.Response {
	bankList, err := usecase.BankRepository.GetAllBank()
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return helper.ResponseError("failed", "Data not found", 404)
	} else if err != nil {
		return helper.ResponseError("failed", err, 500)
	}
	helper.PanicIfError(err)
	response := []dto.Bank{}
	for _, bank := range bankList {
		responseData := dto.Bank{
			BankID:        bank.BankID,
			BankName:      bank.BankName,
			AccountName:   bank.AccountName,
			Category:      bank.Category,
			AccountNumber: bank.AccountNumber,
			Balance:       bank.Balance,
			Active:        bank.Active,
			Ibanking:      bank.Ibanking,
			CodeAccess:    bank.CodeAccess,
			Pin:           bank.Pin,
		}
		response = append(response, responseData)
	}

	var result map[string]any = make(map[string]any)
	result["list_bank"] = response

	return helper.ResponseSuccess("ok", nil, result, 200)
}

func (usecase *BankUsecaseImpl) UpdateBank(id string, body dto.BankUpdate) dto.Response {
	err := usecase.Validate.Struct(body)

	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]dto.CustomMessageError, len(ve))
			for i, fe := range ve {
				out[i] = dto.CustomMessageError{
					Field:    fe.Field(),
					Messsage: helper.MessageError(fe.Tag()),
				}
			}
			return helper.ResponseError("failed", out, 400)
		}

	}

	payloadBank := &entity.Bank{
		BankName:      body.BankName,
		AccountName:   body.AccountName,
		Category:      body.Category,
		AccountNumber: body.AccountNumber,
		Active:        body.Active,
		Ibanking:      body.Ibanking,
		CodeAccess:    body.CodeAccess,
		Pin:           body.Pin,
		UpdatedAt:     time.Now().Format("2006-01-02 15:04:05"),
	}

	bankID, err := usecase.BankRepository.UpdateBank(id, payloadBank)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return helper.ResponseError("failed", err.Error(), 404)
	}
	helper.PanicIfError(err)

	return helper.ResponseSuccess("ok", nil, map[string]any{"bank_id": bankID}, 200)
}

func (usecase *BankUsecaseImpl) UpdateBankBalance(body dto.BankUpdateBalance) dto.Response {
	err := usecase.Validate.Struct(body)

	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]dto.CustomMessageError, len(ve))
			for i, fe := range ve {
				out[i] = dto.CustomMessageError{
					Field:    fe.Field(),
					Messsage: helper.MessageError(fe.Tag()),
				}
			}
			return helper.ResponseError("failed", out, 400)
		}

	}

	payloadBank := &entity.Bank{
		BankID:    body.BankID,
		Balance:   body.Balance,
		UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	bankID, err := usecase.BankRepository.UpdateBalanceBank(payloadBank, body.Types)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println(err.Error())
		return helper.ResponseError("failed", err.Error(), 404)
	}

	if err != nil {
		return helper.ResponseError("failed", err.Error(), 400)
	}

	helper.PanicIfError(err)

	return helper.ResponseSuccess("ok", nil, map[string]any{"bank_id": bankID}, 200)
}
