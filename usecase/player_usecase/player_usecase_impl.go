package player_usecase

import (
	"erp-service/helper"
	"erp-service/model/dto"
	"erp-service/model/entity"
	"erp-service/repository/player_repository"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PlayerUsecaseImpl struct {
	PlayerRepository player_repository.PlayerRepository
	Validate         *validator.Validate
}

func NewPlayerUsecase(playerRepository player_repository.PlayerRepository, validate *validator.Validate) PlayerUsecase {
	return &PlayerUsecaseImpl{
		PlayerRepository: playerRepository,
		Validate:         validate,
	}
}

func (usecase *PlayerUsecaseImpl) GetAllPlayer() dto.Response {
	playerList, err := usecase.PlayerRepository.GetAllPlayer()
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return helper.ResponseError("failed", "Data not found", 404)
	} else if err != nil {
		return helper.ResponseError("failed", err, 500)
	}
	helper.PanicIfError(err)

	var result map[string]any = make(map[string]any)
	result["list_player"] = playerList

	return helper.ResponseSuccess("ok", nil, result, 200)
}

func (usecase *PlayerUsecaseImpl) AddPlayer(body dto.AddPlayer) dto.Response {
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
	helper.PanicIfError(err)

	exist, _ := usecase.PlayerRepository.GetPlayerByID(body.PlayerID)

	if !exist {
		return helper.ResponseError("failed", "PlayerID sudah terdaftar", 400)
	}

	payloadPlayer := &entity.Player{
		PlayerID:   body.PlayerID,
		PlayerName: body.PlayerName,
		CreatedAt:  time.Now().Format("2006-01-02 15:04:05"),
	}

	createID := uuid.New().String()

	payloadBankPlayer := &entity.BankPlayer{
		BankPlayerID:  createID,
		PlayerID:      body.PlayerID,
		BankName:      body.BankName,
		AccountNumber: body.AccountNumber,
		CreatedAt:     time.Now().Format("2006-01-02 15:04:05"),
	}

	player, err := usecase.PlayerRepository.AddPlayer(payloadPlayer)
	helper.PanicIfError(err)

	bankPlayer, err := usecase.PlayerRepository.AddBankPlayer(payloadBankPlayer)
	helper.PanicIfError(err)

	return helper.ResponseSuccess("ok", nil, map[string]interface{}{"id": player, "id_bank": bankPlayer}, 201)
}

func (usecase *PlayerUsecaseImpl) AddPlayerBank(body dto.AddBankPlayer) dto.Response {
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
	helper.PanicIfError(err)
	createID := uuid.New().String()

	exist, _ := usecase.PlayerRepository.GetPlayerByID(body.PlayerID)

	if exist {
		return helper.ResponseError("failed", "PlayerID tidak ditemukan", 400)
	}

	bankPlayer := &entity.BankPlayer{
		BankPlayerID:  createID,
		PlayerID:      body.PlayerID,
		BankName:      body.BankName,
		AccountNumber: body.AccountNumber,
		CreatedAt:     time.Now().Format("2006-01-02 15:04:05"),
	}

	idBank, err := usecase.PlayerRepository.AddBankPlayer(bankPlayer)
	helper.PanicIfError(err)

	return helper.ResponseSuccess("ok", nil, map[string]interface{}{"id": idBank}, 201)
}
