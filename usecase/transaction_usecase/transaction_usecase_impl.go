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
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
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

	if types.TypeTransaction == "BONUS" && body.Status == "PENDING" {
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

	var lastBalanceCoin float64
	var lastBalanceBank float64

	if types.TypeTransaction == "WITHDRAW" {
		lastBalanceCoin = coin.Balance + body.Ammount
		lastBalanceBank = bank.Balance - (body.Ammount + body.AdminFee)
	} else if types.TypeTransaction == "DEPOSIT" {
		lastBalanceCoin = coin.Balance - body.Ammount
		lastBalanceBank = bank.Balance + (body.Ammount + body.AdminFee)
	} else if types.TypeTransaction == "BONUS" {
		lastBalanceCoin = coin.Balance - body.Ammount
		lastBalanceBank = bank.Balance
	} else {
		lastBalanceCoin = coin.Balance
		lastBalanceBank = bank.Balance
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
		LastBalanceCoin: lastBalanceCoin,
		LastBalanceBank: lastBalanceBank,
		Status:          body.Status,
		Note:            body.Note,
		CreatedAt:       time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt:       time.Now().Format("2006-01-02 15:04:05"),
	}

	if types.TypeTransaction == "WITHDRAW" {
		if bank.Balance < (body.Ammount + body.AdminFee) {
			return helper.ResponseError("failed", "saldo bank tidak mencukupi", 400)
		}
	}
	var typeMutation string

	switch types.TypeTransaction {
	case "WITHDRAW":
		typeMutation = "DEBET"
	case "DEPOSIT":
		typeMutation = "CREDIT"
	case "BONUS":
		typeMutation = "BONUS"
	default:
		return helper.ResponseError("failed", "transaksi gagal", 400)
	}

	trx, err := usecase.TrxRepository.AddTransaction(payloadTrx, lastBalanceCoin, lastBalanceBank, typeMutation)
	helper.PanicIfError(err)

	return helper.ResponseSuccess("ok", nil, map[string]interface{}{"id": trx}, 201)
}

func (usecase *TransactionUsecaseImpl) GetAllTransaction(roleName string, limit int, offset int, dateFrom string, dateTo string, types string, status string, keyword string, filter string) dto.Response {
	if len(filter) > 0 {
		filterValid := []string{"note", "player_id"}

		if !slices.Contains(filterValid, filter) {
			return helper.ResponseError("bad request", "Filter Only note and player_id", 400)
		}
	}

	trxList, err := usecase.TrxRepository.GetAllTransaction(roleName, limit, offset, dateFrom, dateTo, types, status, keyword, filter)
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
			Note:                trx.Note,
			CreatedBy:           trx.CreatedBy,
			CreatedAt:           timeCreated.Format("2006-01-02 15:04:05"),
		}
		response = append(response, responseData)
	}

	var result map[string]any = make(map[string]any)
	result["total"] = trxList.Total
	result["total_deposit"] = trxList.TotalDeposit
	result["total_withdraw"] = trxList.TotalWithdraw
	result["total_bonus"] = trxList.TotalBonus
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

	typeMutation := "DEBET"
	desc := "Mengurangi saldo coin dari update transaksi pending dari ID Transaksi " + trxUpdate

	payloadMutation := entity.MutationCoin{
		MutationCoinID:    uuid.New().String(),
		Type:              typeMutation,
		Ammount:           trx.Ammount,
		LastBalance:       balanceCoin,
		Description:       desc,
		IsTransactionBank: false,
		CreatedAt:         time.Now().Format("2006-01-02 15:04:05"),
	}

	mutationID, err := usecase.CoinRepository.TransactionCoin(payloadMutation)
	if err != nil {
		return helper.ResponseError("failed", err.Error(), 400)
	}

	return helper.ResponseSuccess("ok", nil, map[string]interface{}{"id": trxUpdate, "mutation_coin_id": mutationID}, 200)
}

func (usecase *TransactionUsecaseImpl) CanceledTransaction(transactionID string) dto.Response {

	if len(transactionID) == 0 {
		return helper.ResponseError("failed", "transaction id required", 400)
	}

	transaction, err := usecase.TrxRepository.GetDetailTransaction(transactionID)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return helper.ResponseError("failed", err.Error(), 404)
	}

	bank, err := usecase.BankRepository.GetDetailBank(transaction.BankID)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return helper.ResponseError("failed", err.Error(), 404)
	}

	typeTrx, err := usecase.TypeRepository.GetDetailType(transaction.TypeID)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return helper.ResponseError("failed", err.Error(), 404)
	}

	coin, _ := usecase.CoinRepository.GetDetailCoin()

	var balanceBank float64
	var balanceCoin float64
	var typeMutationBank string
	var typeMutationCoin string
	var descBank string
	var descCoin string

	if typeTrx.TypeTransaction == "DEPOSIT" {
		if bank.Balance <= transaction.Ammount {
			return helper.ResponseError("failed", "bank balance is not enough, please top up the bank balance", 400)
		}

		if transaction.Status == "PENDING" {
			balanceCoin = coin.Balance
		} else {
			balanceCoin = coin.Balance + transaction.Ammount
			typeMutationCoin = "CREDIT"
			descCoin = "Saldo coin ditambahkan sebesar " + fmt.Sprintf("%.2f", transaction.Ammount) + " dari cancel transaksi ID " + transaction.TransactionID
		}
		balanceBank = bank.Balance - (transaction.Ammount + transaction.AdminFee)
		typeMutationBank = "DEBET"
		descBank = "Saldo bank dikurangi sebesar " + fmt.Sprintf("%.2f", transaction.Ammount+transaction.AdminFee) + " dari cancel transaksi ID " + transaction.TransactionID
	} else if typeTrx.TypeTransaction == "WITHDRAW" {
		balanceBank = bank.Balance + (transaction.Ammount + transaction.AdminFee)
		balanceCoin = coin.Balance - transaction.Ammount
		typeMutationBank = "CREDIT"
		typeMutationCoin = "DEBET"
		descBank = "Saldo bank ditambahkan sebesar " + fmt.Sprintf("%.2f", transaction.Ammount+transaction.AdminFee) + " dari cancel transaksi ID " + transaction.TransactionID
		descCoin = "Saldo coin dikurangi sebesar " + fmt.Sprintf("%.2f", transaction.Ammount) + " dari cancel transaksi ID " + transaction.TransactionID
	} else if typeTrx.TypeTransaction == "BONUS" {
		balanceBank = bank.Balance
		balanceCoin = coin.Balance + transaction.Ammount
		typeMutationCoin = "CREDIT"
		descCoin = "Saldo coin ditambahkan sebesar " + fmt.Sprintf("%.2f", transaction.Ammount) + " dari cancel transaksi ID " + transaction.TransactionID
	} else {
		return helper.ResponseError("failed", "Cancel transaction failed, Transaction not found", 404)
	}

	cancelTrx, err := usecase.TrxRepository.CanceledTransaction(transaction.TransactionID, bank.BankID, balanceBank, balanceCoin)

	if err != nil {
		return helper.ResponseError("failed", err.Error(), 404)
	}

	payloadMutationBank := entity.MutationBank{
		MutationBankID:    uuid.New().String(),
		BankID:            bank.BankID,
		Type:              typeMutationBank,
		Ammount:           (transaction.Ammount + transaction.AdminFee),
		LastBalance:       balanceBank,
		Description:       descBank,
		IsTransactionBank: false,
		CreatedAt:         time.Now().Format("2006-01-02 15:04:05"),
	}

	mutationID, err := usecase.BankRepository.TransactionBank(payloadMutationBank)
	if err != nil {
		return helper.ResponseError("failed", err.Error(), 400)
	}

	if transaction.Status != "PENDING" {
		payloadMutationCoin := entity.MutationCoin{
			MutationCoinID:    uuid.New().String(),
			Type:              typeMutationCoin,
			Ammount:           transaction.Ammount,
			LastBalance:       balanceCoin,
			Description:       descCoin,
			IsTransactionBank: false,
			CreatedAt:         time.Now().Format("2006-01-02 15:04:05"),
		}

		_, err := usecase.CoinRepository.TransactionCoin(payloadMutationCoin)
		if err != nil {
			return helper.ResponseError("failed", err.Error(), 400)
		}
	}

	return helper.ResponseSuccess("ok", nil, map[string]interface{}{"id": cancelTrx, "mutation_bank_id": mutationID}, 200)
}
