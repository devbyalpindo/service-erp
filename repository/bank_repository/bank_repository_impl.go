package bank_repository

import (
	"erp-service/helper"
	"erp-service/model/entity"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BankRepositoryImpl struct {
	DB *gorm.DB
}

func NewBankRepository(DB *gorm.DB) BankRepository {
	return &BankRepositoryImpl{DB: DB}
}

func (repository *BankRepositoryImpl) AddBank(bank *entity.Bank) (*string, error) {
	if err := repository.DB.Create(&bank).Error; err != nil {
		return nil, err
	}

	return &bank.BankID, nil
}

func (repository *BankRepositoryImpl) GetAllBank() ([]entity.Bank, error) {
	bank := []entity.Bank{}
	err := repository.DB.Model(&entity.Bank{}).Scan(&bank).Error
	helper.PanicIfError(err)
	if len(bank) <= 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return bank, nil
}

func (repository *BankRepositoryImpl) GetDetailBank(id string) (*entity.Bank, error) {
	bank := entity.Bank{}
	result := repository.DB.Where("bank_id = ?", id).Find(&bank)

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &bank, nil
}

func (repository *BankRepositoryImpl) UpdateBank(id string, bank *entity.Bank) (*string, error) {
	result := repository.DB.Model(&bank).Where("bank_id = ?", id).Updates(entity.Bank{BankName: bank.BankName, AccountName: bank.AccountName, Category: bank.Category, AccountNumber: bank.AccountNumber, Active: bank.Active, Ibanking: bank.Ibanking, CodeAccess: bank.CodeAccess, Pin: bank.Pin, UpdatedAt: bank.UpdatedAt})
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &id, nil
}

func (repository *BankRepositoryImpl) UpdateBalanceBank(bank *entity.Bank, types string) (*string, error) {
	bankResult := entity.Bank{}
	result := repository.DB.Where("bank_id = ?", bank.BankID).Find(&bankResult)
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	balance := bankResult.Balance

	if types == "MINUS" {
		if balance < bank.Balance {
			return nil, errors.New("saldo bank tidak mencukupi")
		}
		balance = balance - bank.Balance
	}

	if types == "PLUS" {
		balance = balance + bank.Balance
	}

	results := repository.DB.Model(&bank).Where("bank_id = ?", bank.BankID).Updates(map[string]interface{}{"balance": balance, "updated_at": bank.UpdatedAt})

	if results.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &bank.BankID, nil
}

func (repository *BankRepositoryImpl) TransferToBank(idFrom string, balanceBankFrom float32, idBankTo string, balanceBankTo float32, ammount float32) (*string, error) {
	bank := entity.Bank{}
	mutationBank := []entity.MutationBank{
		{
			MutationBankID: uuid.New().String(),
			BankID:         idFrom,
			Type:           "DEBET",
			Ammount:        ammount,
			Description:    "Transfer ke bank lain",
			CreatedAt:      time.Now().Format("2006-01-02 15:04:05"),
		},
		{
			MutationBankID: uuid.New().String(),
			BankID:         idBankTo,
			Type:           "CREDIT",
			Ammount:        ammount,
			Description:    "Menerima transfer dari bank lain",
			CreatedAt:      time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	tx := repository.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&mutationBank).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Model(&bank).Where("bank_id = ?", idFrom).Updates(map[string]interface{}{"balance": balanceBankFrom, "updated_at": time.Now().Format("2006-01-02 15:04:05")}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Model(&bank).Where("bank_id = ?", idBankTo).Updates(map[string]interface{}{"balance": balanceBankTo, "updated_at": time.Now().Format("2006-01-02 15:04:05")}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &idBankTo, nil
}

func (repository *BankRepositoryImpl) TransactionBank(bank entity.MutationBank) (*string, error) {
	if err := repository.DB.Create(&bank).Error; err != nil {
		return nil, err
	}

	return &bank.MutationBankID, nil
}

func (repository *BankRepositoryImpl) GetMutation(bankID string, types string, limit int, offset int, dateFrom string, dateTo string) (entity.BankithTotal, error) {

	mutation := []entity.BankJoinMutation{}
	var err error
	var totalData int64

	if len(bankID) > 0 && len(types) > 0 {
		err = repository.DB.Table("mutation_banks").Select("mutation_banks.mutation_bank_id, mutation_banks.bank_id, banks.bank_name, banks.account_number, mutation_banks.type, mutation_banks.ammount, banks.balance as last_balance, mutation_banks.description").Joins("inner join banks on mutation_banks.bank_id = banks.bank_id").Where("mutation_banks.bank_id = ? AND mutation_banks.type = ? AND DATE(mutation_banks.created_at) >= ? AND DATE(mutation_banks.created_at) <= ?", bankID, types, dateFrom, dateTo).Order("mutation_banks.created_at DESC").Count(&totalData).Limit(limit).Offset(offset).Find(&mutation).Error
	}

	if len(bankID) == 0 && len(types) > 0 {
		err = repository.DB.Table("mutation_banks").Select("mutation_banks.mutation_bank_id, mutation_banks.bank_id, banks.bank_name, banks.account_number, mutation_banks.type, mutation_banks.ammount, banks.balance as last_balance, mutation_banks.description").Joins("inner join banks on mutation_banks.bank_id = banks.bank_id").Where("mutation_banks.type = ? AND DATE(mutation_banks.created_at) >= ? AND DATE(mutation_banks.created_at) <= ?", types, dateFrom, dateTo).Order("mutation_banks.created_at DESC").Count(&totalData).Limit(limit).Offset(offset).Find(&mutation).Error
	}

	if len(types) == 0 && len(bankID) > 0 {
		err = repository.DB.Table("mutation_banks").Select("mutation_banks.mutation_bank_id, mutation_banks.bank_id, banks.bank_name, banks.account_number, mutation_banks.type, mutation_banks.ammount, banks.balance as last_balance, mutation_banks.description").Joins("inner join banks on mutation_banks.bank_id = banks.bank_id").Where("mutation_banks.bank_id = ? AND DATE(mutation_banks.created_at) >= ? AND DATE(mutation_banks.created_at) <= ?", bankID, dateFrom, dateTo).Order("mutation_banks.created_at DESC").Count(&totalData).Limit(limit).Offset(offset).Find(&mutation).Error
	}

	if len(bankID) == 0 && len(types) == 0 {
		err = repository.DB.Table("mutation_banks").Select("mutation_banks.mutation_bank_id, mutation_banks.bank_id, banks.bank_name, banks.account_number, mutation_banks.type, mutation_banks.ammount, banks.balance as last_balance, mutation_banks.description").Joins("inner join banks on mutation_banks.bank_id = banks.bank_id").Where("DATE(mutation_banks.created_at) >= ? AND DATE(mutation_banks.created_at) <= ?", dateFrom, dateTo).Order("mutation_banks.created_at DESC").Count(&totalData).Limit(limit).Offset(offset).Find(&mutation).Error
	}

	helper.PanicIfError(err)

	if len(mutation) <= 0 {
		resultError := entity.BankithTotal{
			Total:    0,
			Mutation: nil,
		}
		return resultError, gorm.ErrRecordNotFound
	}

	result := entity.BankithTotal{
		Total:    totalData,
		Mutation: mutation,
	}

	return result, nil
}
