package transaction_repository

import (
	"erp-service/helper"
	"erp-service/model/dto"
	"erp-service/model/entity"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepositoryImpl struct {
	DB *gorm.DB
}

func NewTransactionRepository(DB *gorm.DB) TransactionRepository {
	return &TransactionRepositoryImpl{DB}
}

func (repository *TransactionRepositoryImpl) AddTransaction(trx *entity.Transaction, balanceCoin float32, balanceBank float32, typeTrx string) (string, error) {
	bank := entity.Bank{}
	coin := entity.Coin{}
	var desc string
	if typeTrx == "DEBET" {
		desc = "Mengurangi saldo bank dari transaksi WITHDRAW"
	}

	if typeTrx == "CREDIT" {
		desc = "Menambah saldo bank dari transaksi DEPOSIT"
	}

	mutation := entity.MutationBank{
		MutationBankID: uuid.New().String(),
		BankID:         trx.BankID,
		Type:           typeTrx,
		Ammount:        trx.Ammount,
		LastBalance:    balanceBank,
		Description:    desc,
		CreatedAt:      time.Now().Format("2006-01-02 15:04:05"),
	}

	tx := repository.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if typeTrx != "NOT" {
		if err := tx.Create(&mutation).Error; err != nil {
			tx.Rollback()
			return "", err
		}
	}

	if err := tx.Create(&trx).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	if err := tx.Model(&bank).Where("bank_id = ?", trx.BankID).Updates(map[string]interface{}{"balance": balanceBank, "updated_at": time.Now().Format("2006-01-02 15:04:05")}).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	if trx.Status == "COMPLETED" {
		if err := tx.Model(&coin).Where("coin_name = ?", "SALJU 88").Updates(map[string]interface{}{"balance": balanceCoin, "updated_at": time.Now().Format("2006-01-02 15:04:05")}).Error; err != nil {
			tx.Rollback()
			return "", err
		}
	}

	tx.Commit()

	return trx.TransactionID, nil
}

func (repository *TransactionRepositoryImpl) GetAllTransaction(roleName string, limit int, offset int, dateFrom string, dateTo string, types string) (dto.TransactionWithTotal, error) {

	trx := []entity.TransactionJoin{}
	var totalData int64
	var err error

	if roleName == "ADMIN" {
		if len(types) > 0 {
			err = repository.DB.Table("transactions").Select("transactions.transaction_id, transactions.user_id, players.player_name, players.player_id, bank_players.bank_name as bank_player_name, bank_players.account_number as account_number_player, transactions.bank_id, banks.bank_name, banks.account_number as account_number_bank, transactions.type_id, type_transactions.type_transaction, transactions.ammount, transactions.admin_fee, transactions.last_balance_coin, transactions.last_balance_bank, transactions.status, users.username as created_by, transactions.created_at, transactions.updated_at").Joins("inner join banks on banks.bank_id = transactions.bank_id").Joins("inner join type_transactions on type_transactions.type_id = transactions.type_id").Joins("inner join users on users.user_id = transactions.user_id").Joins("left join players on players.player_id = transactions.player_id").Joins("left join bank_players on bank_players.bank_player_id = transactions.bank_player_id").Where("DATE(transactions.created_at) >= ? AND DATE(transactions.created_at) <= ? AND type_transactions.type_transaction = ?", dateFrom, dateTo, strings.ToUpper(types)).Order("transactions.created_at DESC").Count(&totalData).Limit(limit).Offset(offset).Find(&trx).Error
		} else {
			err = repository.DB.Table("transactions").Select("transactions.transaction_id, transactions.user_id, players.player_name, players.player_id, bank_players.bank_name as bank_player_name, bank_players.account_number as account_number_player, transactions.bank_id, banks.bank_name, banks.account_number as account_number_bank, transactions.type_id, type_transactions.type_transaction, transactions.ammount, transactions.admin_fee, transactions.last_balance_coin, transactions.last_balance_bank, transactions.status, users.username as created_by, transactions.created_at, transactions.updated_at").Joins("inner join banks on banks.bank_id = transactions.bank_id").Joins("inner join type_transactions on type_transactions.type_id = transactions.type_id").Joins("inner join users on users.user_id = transactions.user_id").Joins("left join players on players.player_id = transactions.player_id").Joins("left join bank_players on bank_players.bank_player_id = transactions.bank_player_id").Where("DATE(transactions.created_at) >= ? AND DATE(transactions.created_at) <= ?", dateFrom, dateTo).Order("transactions.created_at DESC").Count(&totalData).Limit(limit).Offset(offset).Find(&trx).Error
		}
	}

	if roleName == "DEPOSITOR" {
		if len(types) > 0 {
			err = repository.DB.Table("transactions").Select("transactions.transaction_id, transactions.user_id, players.player_name, players.player_id, bank_players.bank_name as bank_player_name, bank_players.account_number as account_number_player, transactions.bank_id, banks.bank_name, banks.account_number as account_number_bank, transactions.type_id, type_transactions.type_transaction, transactions.ammount, transactions.admin_fee, transactions.last_balance_coin, transactions.last_balance_bank,  transactions.status, users.username as created_by, transactions.created_at, transactions.updated_at").Joins("inner join banks on banks.bank_id = transactions.bank_id").Joins("inner join type_transactions on type_transactions.type_id = transactions.type_id").Joins("inner join users on users.user_id = transactions.user_id").Joins("left join players on players.player_id = transactions.player_id").Joins("left join bank_players on bank_players.bank_player_id = transactions.bank_player_id").Where("DATE(transactions.created_at) >= ? AND DATE(transactions.created_at) <= ? AND type_transaction IN ('DEPOSIT', 'BONUS') AND type_transactions.type_transaction = ?", dateFrom, dateTo, strings.ToUpper(types)).Order("transactions.created_at DESC").Count(&totalData).Limit(limit).Offset(offset).Find(&trx).Error
		} else {
			err = repository.DB.Table("transactions").Select("transactions.transaction_id, transactions.user_id, players.player_name, players.player_id, bank_players.bank_name as bank_player_name, bank_players.account_number as account_number_player, transactions.bank_id, banks.bank_name, banks.account_number as account_number_bank, transactions.type_id, type_transactions.type_transaction, transactions.ammount, transactions.admin_fee, transactions.last_balance_coin, transactions.last_balance_bank,  transactions.status, users.username as created_by, transactions.created_at, transactions.updated_at").Joins("inner join banks on banks.bank_id = transactions.bank_id").Joins("inner join type_transactions on type_transactions.type_id = transactions.type_id").Joins("inner join users on users.user_id = transactions.user_id").Joins("left join players on players.player_id = transactions.player_id").Joins("left join bank_players on bank_players.bank_player_id = transactions.bank_player_id").Where("DATE(transactions.created_at) >= ? AND DATE(transactions.created_at) <= ? AND type_transaction IN ('DEPOSIT', 'BONUS')", dateFrom, dateTo).Order("transactions.created_at DESC").Count(&totalData).Limit(limit).Offset(offset).Find(&trx).Error
		}

	}

	if roleName == "WITHDRAWER" {
		if len(types) > 0 {
			err = repository.DB.Table("transactions").Select("transactions.transaction_id, transactions.user_id, players.player_name, players.player_id, bank_players.bank_name as bank_player_name, bank_players.account_number as account_number_player, transactions.bank_id, banks.bank_name, banks.account_number as account_number_bank, transactions.type_id, type_transactions.type_transaction, transactions.ammount, transactions.admin_fee, transactions.last_balance_coin, transactions.last_balance_bank, transactions.status, users.username as created_by, transactions.created_at, transactions.updated_at").Joins("inner join banks on banks.bank_id = transactions.bank_id").Joins("inner join type_transactions on type_transactions.type_id = transactions.type_id").Joins("inner join users on users.user_id = transactions.user_id").Joins("left join players on players.player_id = transactions.player_id").Joins("left join bank_players on bank_players.bank_player_id = transactions.bank_player_id").Where("DATE(transactions.created_at) >= ? AND DATE(transactions.created_at) <= ? AND type_transaction = ? AND type_transactions.type_transaction = ?", dateFrom, dateTo, "WITHDRAW", strings.ToUpper(types)).Order("transactions.created_at DESC").Count(&totalData).Limit(limit).Offset(offset).Find(&trx).Error
		} else {
			err = repository.DB.Table("transactions").Select("transactions.transaction_id, transactions.user_id, players.player_name, players.player_id, bank_players.bank_name as bank_player_name, bank_players.account_number as account_number_player, transactions.bank_id, banks.bank_name, banks.account_number as account_number_bank, transactions.type_id, type_transactions.type_transaction, transactions.ammount, transactions.admin_fee, transactions.last_balance_coin, transactions.last_balance_bank, transactions.status, users.username as created_by, transactions.created_at, transactions.updated_at").Joins("inner join banks on banks.bank_id = transactions.bank_id").Joins("inner join type_transactions on type_transactions.type_id = transactions.type_id").Joins("inner join users on users.user_id = transactions.user_id").Joins("left join players on players.player_id = transactions.player_id").Joins("left join bank_players on bank_players.bank_player_id = transactions.bank_player_id").Where("DATE(transactions.created_at) >= ? AND DATE(transactions.created_at) <= ? AND type_transaction = ?", dateFrom, dateTo, "WITHDRAW").Order("transactions.created_at DESC").Count(&totalData).Limit(limit).Offset(offset).Find(&trx).Error
		}

	}

	helper.PanicIfError(err)

	var totalDeposit float32
	var totalWithdraw float32
	var totalBonus float32

	for _, item := range trx {
		switch item.TypeTransaction {
		case "DEPOSIT":
			totalDeposit += item.Ammount
		case "WITHDRAW":
			totalWithdraw += item.Ammount
		case "BONUS":
			totalBonus += item.Ammount
		default:
			totalDeposit += 0
			totalWithdraw += 0
			totalBonus += 0
		}
	}

	if len(trx) <= 0 {
		resultError := dto.TransactionWithTotal{
			Total:         0,
			TotalDeposit:  0,
			TotalWithdraw: 0,
			Transaction:   nil,
		}
		return resultError, gorm.ErrRecordNotFound
	}

	result := dto.TransactionWithTotal{
		Total:         totalData,
		TotalDeposit:  totalDeposit,
		TotalWithdraw: totalWithdraw,
		TotalBonus:    totalBonus,
		Transaction:   trx,
	}

	return result, nil
}

func (repository *TransactionRepositoryImpl) GetDetailTransaction(id string) (*entity.Transaction, error) {
	trx := entity.Transaction{}
	result := repository.DB.Where("transaction_id = ?", id).Find(&trx)

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &trx, nil
}

func (repository *TransactionRepositoryImpl) UpdateTransaction(transactionID string, playerID string, bankPlayerID string, status string, balanceCoin float32) (string, error) {
	trx := entity.Transaction{}
	coin := entity.Coin{}

	tx := repository.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&coin).Where("coin_name = ?", "SALJU 88").Updates(map[string]interface{}{"balance": balanceCoin, "updated_at": time.Now().Format("2006-01-02 15:04:05")}).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	if err := tx.Model(&trx).Where("transaction_id = ?", transactionID).Updates(entity.Transaction{PlayerID: playerID, BankPlayerID: bankPlayerID, Status: status, UpdatedAt: time.Now().Format("2006-01-02 15:04:05")}).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	tx.Commit()

	return transactionID, nil
}
