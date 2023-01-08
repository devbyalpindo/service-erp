package bonus_usecase

import (
	"erp-service/helper"
	"erp-service/model/dto"
	"erp-service/model/entity"
	"erp-service/repository/bonus_repository"
	"erp-service/repository/coin_repository"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

type BonusUsecaseImpl struct {
	BonusRepository bonus_repository.BonusRepository
	CoinRepository  coin_repository.CoinRepository
	Validate        *validator.Validate
}

func NewBonusUsecase(bonusRepository bonus_repository.BonusRepository, coinRepository coin_repository.CoinRepository, validate *validator.Validate) BonusUsecase {
	return &BonusUsecaseImpl{
		BonusRepository: bonusRepository,
		CoinRepository:  coinRepository,
		Validate:        validate,
	}
}

func (usecase *BonusUsecaseImpl) AddBonus(body dto.BonusAdd) dto.Response {
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

	typeValid := []string{"CASHBACK", "ROLLING"}

	if !slices.Contains(typeValid, body.Type) {
		return helper.ResponseError("bad request", "Type Only CASHBACK AND ROLLING", 400)
	}

	coin, err := usecase.CoinRepository.GetDetailCoin()
	if err != nil {
		return helper.ResponseError("failed", "coin not found", 404)
	}

	createID := uuid.New().String()
	helper.PanicIfError(err)

	payloadBonus := &entity.Bonus{
		BonusID:   createID,
		Type:      body.Type,
		Ammount:   body.Ammount,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	bonus, err := usecase.BonusRepository.AddBonus(payloadBonus, coin.Balance-payloadBonus.Ammount)
	helper.PanicIfError(err)

	typeMutation := "DEBET"
	desc := "Bonus " + body.Type

	payloadMutation := entity.MutationCoin{
		MutationCoinID:    uuid.New().String(),
		Type:              typeMutation,
		Ammount:           body.Ammount,
		LastBalance:       coin.Balance - payloadBonus.Ammount,
		Description:       desc,
		IsTransactionBank: true,
		CreatedAt:         time.Now().Format("2006-01-02 15:04:05"),
	}

	mutationID, err := usecase.CoinRepository.TransactionCoin(payloadMutation)
	if err != nil {
		return helper.ResponseError("failed", err.Error(), 400)
	}

	return helper.ResponseSuccess("ok", nil, map[string]interface{}{"id": bonus, "mutation_id": mutationID}, 201)

}

func (usecase *BonusUsecaseImpl) GetAllBonus(types string, dateFrom string, dateTo string, limit int, offset int) dto.Response {
	bonusList, err := usecase.BonusRepository.GetAllBonus(types, dateFrom, dateTo, limit, offset)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return helper.ResponseError("failed", "Data not found", 404)
	} else if err != nil {
		return helper.ResponseError("failed", err, 500)
	}
	helper.PanicIfError(err)

	response := []dto.Bonus{}
	for _, bns := range bonusList.Bonus {
		timeCreated, _ := time.Parse(time.RFC3339, bns.CreatedAt)
		timeUpdated, _ := time.Parse(time.RFC3339, bns.UpdatedAt)
		responseData := dto.Bonus{
			BonusID:   bns.BonusID,
			Type:      bns.Type,
			Ammount:   bns.Ammount,
			CreatedAt: timeCreated.Format("2006-01-02 15:04:05"),
			UpdatedAt: timeUpdated.Format("2006-01-02 15:04:05"),
		}
		response = append(response, responseData)
	}

	var result map[string]any = make(map[string]any)
	result["total"] = bonusList.Total
	result["total_bonus"] = bonusList.TotalBonus
	result["bonus_list"] = response

	return helper.ResponseSuccess("ok", nil, result, 200)
}

func (usecase *BonusUsecaseImpl) UpdateBonus(id string, body dto.BonusAdd) dto.Response {
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

	typeValid := []string{"CASHBACK", "ROLLING"}

	if !slices.Contains(typeValid, body.Type) {
		return helper.ResponseError("bad request", "Type Only CASHBACK AND ROLLING", 400)
	}

	coin, err := usecase.CoinRepository.GetDetailCoin()
	if err != nil {
		return helper.ResponseError("failed", "coin not found", 404)
	}

	bonus, err := usecase.BonusRepository.GetBonusByID(id)
	if err != nil {
		return helper.ResponseError("failed", "coin not found", 404)
	}
	helper.PanicIfError(err)

	balanceCoin := coin.Balance + (bonus.Ammount - body.Ammount)

	bonusResult, err := usecase.BonusRepository.UpdateBonus(bonus.BonusID, body.Type, body.Ammount, balanceCoin)
	helper.PanicIfError(err)

	return helper.ResponseSuccess("ok", nil, map[string]interface{}{"id": bonusResult}, 200)

}
