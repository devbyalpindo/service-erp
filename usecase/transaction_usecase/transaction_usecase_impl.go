package transaction_usecase

import (
	"erp-service/helper"
	"erp-service/model/dto"
	"erp-service/model/entity"
	"erp-service/repository/bank_repository"
	"erp-service/repository/coin_repository"
	"erp-service/repository/transaction_repository"
	"erp-service/repository/type_repository"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionUsecaseImpl struct {
	TrxRepository  transaction_repository.TransactionRepository
	CoinRepository coin_repository.CoinRepository
	BankRepository bank_repository.BankRepository
	TypeRepository type_repository.TypeRepository
	Validate       *validator.Validate
}

func NewTransactionUsecase(trxRepository transaction_repository.TransactionRepository, coinRepository coin_repository.CoinRepository, bankRepository bank_repository.BankRepository, typeRepository type_repository.TypeRepository, validate *validator.Validate) TransactionUsecase {
	return &TransactionUsecaseImpl{
		TrxRepository:  trxRepository,
		CoinRepository: coinRepository,
		BankRepository: bankRepository,
		TypeRepository: typeRepository,
		Validate:       validate,
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

	coin, err := usecase.CoinRepository.GetDetailCoin()
	if err != nil {
		return helper.ResponseError("failed", "coin not found", 404)
	}

	bank, err := usecase.BankRepository.GetDetailBank(body.BankID)
	if err != nil {
		return helper.ResponseError("failed", "bank not found", 404)
	}

	types, err := usecase.TypeRepository.GetDetailType(body.TypeID)
	if err != nil {
		return helper.ResponseError("failed", "type not found", 404)
	}

	payloadTrx := &entity.Transaction{
		TransactionID:   createID,
		UserID:          userID,
		PlayerName:      body.PlayerName,
		PlayerID:        body.PlayerID,
		BankPlayer:      body.BankPlayer,
		AccountNumber:   body.AccountNumber,
		BankID:          body.BankID,
		TypeID:          body.TypeID,
		Ammount:         body.Ammount,
		AdminFee:        body.AdminFee,
		LastBalanceCoin: coin.Balance,
		LastBalanceBank: bank.Balance,
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

	switch types.TypeTransaction {
	case "WITHDRAW":
		bankBalance = bank.Balance - (body.Ammount + body.AdminFee)
		coinBalance = coin.Balance - body.Ammount
	case "DEPOSIT":
		bankBalance = bank.Balance + (body.Ammount + body.AdminFee)
		coinBalance = coin.Balance + body.Ammount
	case "BONUS":
		bankBalance = bank.Balance - (body.Ammount + body.AdminFee)
		coinBalance = coin.Balance - body.Ammount
	default:
		return helper.ResponseError("failed", "transaksi gagal", 400)
	}

	trx, err := usecase.TrxRepository.AddTransaction(payloadTrx, coinBalance, bankBalance)
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
		responseData := dto.TransactionJoin{
			TransactionID:     trx.TransactionID,
			UserID:            trx.UserID,
			PlayerName:        trx.PlayerName,
			PlayerID:          trx.PlayerID,
			BankPlayer:        trx.BankPlayer,
			AccountNumber:     trx.AccountNumber,
			BankID:            trx.BankID,
			BankName:          trx.BankName,
			AccountNumberBank: trx.AccountNumberBank,
			TypeID:            trx.TypeID,
			TypeTransaction:   trx.TypeTransaction,
			Ammount:           trx.Ammount,
			AdminFee:          trx.AdminFee,
			LastBalanceCoin:   trx.LastBalanceCoin,
			LastBalanceBank:   trx.LastBalanceBank,
			CreatedBy:         trx.CreatedBy,
			CreatedAt:         trx.CreatedAt,
		}
		response = append(response, responseData)
	}

	var result map[string]any = make(map[string]any)
	result["total"] = trxList.Total
	result["transaction"] = response

	return helper.ResponseSuccess("ok", nil, result, 200)
}
