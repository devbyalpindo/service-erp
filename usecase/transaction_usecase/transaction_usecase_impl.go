package transaction_usecase

import (
	"erp-service/helper"
	"erp-service/model/dto"
	"erp-service/model/entity"
	"erp-service/repository/bank_repository"
	"erp-service/repository/coin_repository"
	"erp-service/repository/player_repository"
	"erp-service/repository/transaction_repository"
	"erp-service/repository/type_repository"
	"errors"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionUsecaseImpl struct {
	TrxRepository    transaction_repository.TransactionRepository
	PlayerRepository player_repository.PlayerRepository
	CoinRepository   coin_repository.CoinRepository
	BankRepository   bank_repository.BankRepository
	TypeRepository   type_repository.TypeRepository
	Validate         *validator.Validate
}

func NewTransactionUsecase(trxRepository transaction_repository.TransactionRepository, coinRepository coin_repository.CoinRepository, bankRepository bank_repository.BankRepository, typeRepository type_repository.TypeRepository, playerRepository player_repository.PlayerRepository, validate *validator.Validate) TransactionUsecase {
	return &TransactionUsecaseImpl{
		TrxRepository:    trxRepository,
		PlayerRepository: playerRepository,
		CoinRepository:   coinRepository,
		BankRepository:   bankRepository,
		TypeRepository:   typeRepository,
		Validate:         validate,
	}
}

func (usecase *TransactionUsecaseImpl) AddTransaction(userID string, body dto.AddTransaction) dto.Response {
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

	types, err := usecase.TypeRepository.GetDetailType(body.TypeID)
	if err != nil {
		return helper.ResponseError("failed", "type not found", 404)
	}

	if types.TypeTransaction == "WITHDRAW" && body.Status == "PENDING" {
		return helper.ResponseError("failed", "status PENDING only DEPOSIT", 400)
	}

	if body.Status == "COMPLETED" && len(body.PlayerID) == 0 {
		return helper.ResponseError("failed", "Player ID is required", 400)
	}

	if body.Status == "COMPLETED" {
		player, _ := usecase.PlayerRepository.GetPlayerByID(body.PlayerID)
		if player {
			return helper.ResponseError("failed", "player not found", 404)
		}

		_, err := usecase.PlayerRepository.GetPlayerBankByID(body.BankPlayerID)
		if err != nil {
			return helper.ResponseError("failed", "bank player not found", 404)
		}
	}

	coin, err := usecase.CoinRepository.GetDetailCoin()
	if err != nil {
		return helper.ResponseError("failed", "coin not found", 404)
	}

	bank, err := usecase.BankRepository.GetDetailBank(body.BankID)
	if err != nil {
		return helper.ResponseError("failed", "bank not found", 404)
	}

	payloadTrx := &entity.Transaction{
		TransactionID:   createID,
		UserID:          userID,
		PlayerID:        body.PlayerID,
		BankPlayerID:    body.BankPlayerID,
		BankID:          body.BankID,
		TypeID:          body.TypeID,
		Ammount:         body.Ammount,
		AdminFee:        body.AdminFee,
		LastBalanceCoin: coin.Balance,
		LastBalanceBank: bank.Balance,
		Status:          body.Status,
		CreatedAt:       time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt:       time.Now().Format("2006-01-02 15:04:05"),
	}

	if types.TypeTransaction == "WITHDRAW" {
		if bank.Balance < (body.Ammount + body.AdminFee) {
			return helper.ResponseError("failed", "saldo bank tidak mencukupi", 400)
		}
	}

	var bankBalance float32
	var coinBalance float32
	var typeMutation string

	switch types.TypeTransaction {
	case "WITHDRAW":
		bankBalance = bank.Balance - (body.Ammount + body.AdminFee)
		coinBalance = coin.Balance + body.Ammount
		typeMutation = "DEBET"
	case "DEPOSIT":
		bankBalance = bank.Balance + (body.Ammount + body.AdminFee)
		coinBalance = coin.Balance - body.Ammount
		typeMutation = "CREDIT"
	case "BONUS":
		bankBalance = bank.Balance
		coinBalance = coin.Balance + body.Ammount
		typeMutation = "NOT"
	default:
		return helper.ResponseError("failed", "transaksi gagal", 400)
	}

	trx, err := usecase.TrxRepository.AddTransaction(payloadTrx, coinBalance, bankBalance, typeMutation)
	helper.PanicIfError(err)

	return helper.ResponseSuccess("ok", nil, map[string]interface{}{"id": trx}, 201)
}

func (usecase *TransactionUsecaseImpl) GetAllTransaction(roleName string, limit int, offset int, dateFrom string, dateTo string) dto.Response {
	trxList, err := usecase.TrxRepository.GetAllTransaction(roleName, limit, offset, dateFrom, dateTo)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return helper.ResponseError("failed", "Data not found", 404)
	} else if err != nil {
		return helper.ResponseError("failed", err, 500)
	}
	helper.PanicIfError(err)
	response := []dto.TransactionJoin{}
	for _, trx := range trxList.Transaction {
		timeCreated, _ := time.Parse(time.RFC3339, trx.CreatedAt)
		responseData := dto.TransactionJoin{
			TransactionID:       trx.TransactionID,
			UserID:              trx.UserID,
			PlayerName:          trx.PlayerName,
			PlayerID:            trx.PlayerID,
			BankPlayerName:      trx.BankPlayerName,
			AccountNumberPlayer: trx.AccountNumberPlayer,
			BankID:              trx.BankID,
			BankName:            trx.BankName,
			AccountNumberBank:   trx.AccountNumberBank,
			TypeID:              trx.TypeID,
			TypeTransaction:     trx.TypeTransaction,
			Ammount:             trx.Ammount,
			AdminFee:            trx.AdminFee,
			LastBalanceCoin:     trx.LastBalanceCoin,
			LastBalanceBank:     trx.LastBalanceBank,
			Status:              trx.Status,
			CreatedBy:           trx.CreatedBy,
			CreatedAt:           timeCreated.Format("2006-01-02 15:04:05"),
		}
		response = append(response, responseData)
	}

	var result map[string]any = make(map[string]any)
	result["total"] = trxList.Total
	result["transaction"] = response

	return helper.ResponseSuccess("ok", nil, result, 200)
}

func (usecase *TransactionUsecaseImpl) UpdateTransaction(transactionID string, body dto.UpdateTransactionPending) dto.Response {
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

	log.Print(body)

	player, _ := usecase.PlayerRepository.GetPlayerByID(body.PlayerID)
	if player {
		return helper.ResponseError("failed", "player not found", 404)
	}

	bankPlayer, err := usecase.PlayerRepository.GetPlayerBankByID(body.BankPlayerID)
	if err != nil {
		return helper.ResponseError("failed", "bank player not found", 404)
	}

	trx, err := usecase.TrxRepository.GetDetailTransaction(transactionID)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return helper.ResponseError("failed", err.Error(), 404)
	}

	if trx.Status != "PENDING" {
		return helper.ResponseError("failed", "Update Transaction Only PENDING STATUS", 400)
	}

	coin, err := usecase.CoinRepository.GetDetailCoin()
	if err != nil {
		return helper.ResponseError("failed", "coin not found", 404)
	}

	balanceCoin := coin.Balance - trx.Ammount

	trxUpdate, err := usecase.TrxRepository.UpdateTransaction(transactionID, body.PlayerID, bankPlayer.BankPlayerID, body.Status, balanceCoin)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return helper.ResponseError("failed", err.Error(), 404)
	}

	if len(trxUpdate) == 0 {
		return helper.ResponseError("failed", "Gagal update transaksi", 404)
	}

	return helper.ResponseSuccess("ok", nil, map[string]interface{}{"id": trxUpdate}, 200)
}
