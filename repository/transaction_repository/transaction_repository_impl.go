package transaction_repository

import (
	"erp-service/helper"
	"erp-service/model/dto"
	"erp-service/model/entity"
	"time"

	"gorm.io/gorm"
)

type TransactionRepositoryImpl struct {
	DB *gorm.DB
}

func NewTransactionRepository(DB *gorm.DB) TransactionRepository {
	return &TransactionRepositoryImpl{DB}
}

func (repository *TransactionRepositoryImpl) AddTransaction(trx *entity.Transaction, balanceCoin float32, balanceBank float32) (string, error) {
	bank := entity.Bank{}
	coin := entity.Coin{}

	tx := repository.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&trx).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	if err := tx.Model(&bank).Where("bank_id = ?", trx.BankID).Updates(map[string]interface{}{"balance": balanceBank, "updated_at": time.Now().Format("2006-01-02 15:04:05")}).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	if err := tx.Model(&coin).Where("coin_name = ?", "SALJU 88").Updates(map[string]interface{}{"balance": balanceCoin, "updated_at": time.Now().Format("2006-01-02 15:04:05")}).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	tx.Commit()

	return trx.TransactionID, nil
}

func (repository *TransactionRepositoryImpl) GetAllTransaction(roleName string, limit int, offset int, dateFrom string, dateTo string) (dto.TransactionWithTotal, error) {

	trx := []entity.TransactionJoin{}
	var totalData int64
	var err error

	if roleName == "ADMIN" {
		err = repository.DB.Table("transactions").Select("transactions.transaction_id, transactions.user_id, transactions.player_name, transactions.player_id, transactions.bank_player, transactions.account_number, transactions.bank_id, banks.bank_name, banks.account_number as account_number_bank, transactions.type_id, type_transactions.type_transaction, transactions.ammount, transactions.admin_fee, transactions.last_balance_coin, transactions.last_balance_bank, users.username as created_by, transactions.created_at, transactions.updated_at").Joins("inner join banks on banks.bank_id = transactions.bank_id").Joins("inner join type_transactions on type_transactions.type_id = transactions.type_id").Joins("inner join users on users.user_id = transactions.user_id").Where("DATE(transactions.created_at) >= ? AND DATE(transactions.created_at) <= ?", dateFrom, dateTo).Order("transactions.created_at DESC").Count(&totalData).Limit(limit).Offset(offset).Find(&trx).Error
	}

	if roleName == "DEPOSITOR" {
		err = repository.DB.Table("transactions").Select("transactions.transaction_id, transactions.user_id, transactions.player_name, transactions.player_id, transactions.bank_player, transactions.account_number, transactions.bank_id, banks.bank_name, banks.account_number as account_number_bank, transactions.type_id, type_transactions.type_transaction, transactions.ammount, transactions.admin_fee, transactions.last_balance_coin, transactions.last_balance_bank, users.username as created_by, transactions.created_at, transactions.updated_at").Joins("inner join banks on banks.bank_id = transactions.bank_id").Joins("inner join type_transactions on type_transactions.type_id = transactions.type_id").Joins("inner join users on users.user_id = transactions.user_id").Where("DATE(transactions.created_at) >= ? AND DATE(transactions.created_at) <= ? AND type_transaction IN ('DEPOSIT', 'BONUS')", dateFrom, dateTo).Order("transactions.created_at DESC").Count(&totalData).Limit(limit).Offset(offset).Find(&trx).Error
	}

	if roleName == "WITHDRAWER" {
		err = repository.DB.Table("transactions").Select("transactions.transaction_id, transactions.user_id, transactions.player_name, transactions.player_id, transactions.bank_player, transactions.account_number, transactions.bank_id, banks.bank_name, banks.account_number as account_number_bank, transactions.type_id, type_transactions.type_transaction, transactions.ammount, transactions.admin_fee, transactions.last_balance_coin, transactions.last_balance_bank, users.username as created_by, transactions.created_at, transactions.updated_at").Joins("inner join banks on banks.bank_id = transactions.bank_id").Joins("inner join type_transactions on type_transactions.type_id = transactions.type_id").Joins("inner join users on users.user_id = transactions.user_id").Where("DATE(transactions.created_at) >= ? AND DATE(transactions.created_at) <= ? AND type_transaction = ?", dateFrom, dateTo, "WITHDRAW").Order("transactions.created_at DESC").Count(&totalData).Limit(limit).Offset(offset).Find(&trx).Error
	}

	helper.PanicIfError(err)

	if len(trx) <= 0 {
		resultError := dto.TransactionWithTotal{
			Total:       0,
			Transaction: nil,
		}
		return resultError, gorm.ErrRecordNotFound
	}

	result := dto.TransactionWithTotal{
		Total:       totalData,
		Transaction: trx,
	}

	return result, nil
}
