package bank_repository

import (
	"erp-service/helper"
	"erp-service/model/entity"
	"errors"

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
