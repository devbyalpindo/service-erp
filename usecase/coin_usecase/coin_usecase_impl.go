package coin_usecase

import (
	"erp-service/helper"
	"erp-service/model/dto"
	"erp-service/model/entity"
	"erp-service/repository/coin_repository"
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CoinUsecaseImpl struct {
	CoinRepository coin_repository.CoinRepository
	Validate       *validator.Validate
}

func NewCoinUsecase(coinRepository coin_repository.CoinRepository, validate *validator.Validate) CoinUsecase {
	return &CoinUsecaseImpl{
		CoinRepository: coinRepository,
		Validate:       validate,
	}
}

func (usecase *CoinUsecaseImpl) GetCoin() dto.Response {
	coinList, err := usecase.CoinRepository.GetCoin()
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return helper.ResponseError("failed", "Data not found", 404)
	} else if err != nil {
		return helper.ResponseError("failed", err, 500)
	}
	helper.PanicIfError(err)
	response := []dto.Coin{}
	for _, coin := range coinList {
		responseData := dto.Coin{
			CoinID:   coin.CoinID,
			CoinName: coin.CoinName,
			Balance:  coin.Balance,
			Note:     coin.Note,
		}
		response = append(response, responseData)
	}

	var result map[string]any = make(map[string]any)
	result["list_coin"] = response

	return helper.ResponseSuccess("ok", nil, result, 200)
}

func (usecase *CoinUsecaseImpl) GetDetailCoin() dto.Response {
	coin, err := usecase.CoinRepository.GetDetailCoin()
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return helper.ResponseError("failed", "Data not found", 404)
	} else if err != nil {
		return helper.ResponseError("failed", err, 500)
	}
	helper.PanicIfError(err)
	response := dto.Coin{
		CoinID:   coin.CoinID,
		CoinName: coin.CoinName,
		Balance:  coin.Balance,
	}

	var result map[string]any = make(map[string]any)
	result["coin"] = response

	return helper.ResponseSuccess("ok", nil, result, 200)
}

func (usecase *CoinUsecaseImpl) UpdateCoinBalance(body dto.CoinUpdateBalance) dto.Response {
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

	payloadCoin := &entity.Coin{
		CoinID:    body.CoinID,
		Balance:   body.Balance,
		UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	coinID, lastBalance, err := usecase.CoinRepository.UpdateBalanceCoin(payloadCoin, body.Types)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return helper.ResponseError("failed", err.Error(), 404)
	}

	if err != nil {
		return helper.ResponseError("failed", err.Error(), 400)
	}

	helper.PanicIfError(err)

	var typeMutation string
	var desc string

	if body.Types == "PLUS" {
		typeMutation = "CREDIT"
		desc = "Saldo coin ditambahkan sebesar " + fmt.Sprintf("%.2f", body.Balance)
	}

	if body.Types == "MINUS" {
		typeMutation = "DEBET"
		desc = "Saldo coin dikurangi sebesar " + fmt.Sprintf("%.2f", body.Balance)
	}

	payloadMutation := entity.MutationCoin{
		MutationCoinID:    uuid.New().String(),
		Type:              typeMutation,
		Ammount:           body.Balance,
		LastBalance:       lastBalance,
		Description:       desc,
		IsTransactionBank: true,
		CreatedAt:         time.Now().Format("2006-01-02 15:04:05"),
	}

	mutationID, err := usecase.CoinRepository.TransactionCoin(payloadMutation)
	if err != nil {
		return helper.ResponseError("failed", err.Error(), 400)
	}

	return helper.ResponseSuccess("ok", nil, map[string]any{"bank_id": coinID, "mutation_coin_id": mutationID}, 200)
}

func (usecase *CoinUsecaseImpl) GetMutation(body dto.GetMutationCoin) dto.Response {
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

	listMutation, err := usecase.CoinRepository.GetMutation(body.Type, body.IsTransactionBank, body.Limit, body.Offset, body.DateFrom, body.DateTo)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return helper.ResponseError("failed", "Data not found", 404)
	} else if err != nil {
		return helper.ResponseError("failed", err, 500)
	}
	helper.PanicIfError(err)

	response := []entity.MutationCoin{}
	for _, item := range listMutation.Mutation {
		timeCreated, _ := time.Parse(time.RFC3339, item.CreatedAt)
		responseData := entity.MutationCoin{
			MutationCoinID: item.MutationCoinID,
			Type:           item.Type,
			Ammount:        item.Ammount,
			LastBalance:    item.LastBalance,
			Description:    item.Description,
			CreatedAt:      timeCreated.Format("2006-01-02 15:04:05"),
		}
		response = append(response, responseData)
	}

	var result map[string]any = make(map[string]any)
	result["total"] = listMutation.Total
	result["mutasi"] = response

	return helper.ResponseSuccess("ok", nil, result, 200)
}
