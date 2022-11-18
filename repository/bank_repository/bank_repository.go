package bank_repository

import (
	"erp-service/model/entity"
)

type BankRepository interface {
	AddBank(*entity.Bank) (*string, error)
	GetAllBank() ([]entity.Bank, error)
	GetDetailBank(id string) (*entity.Bank, error)
	UpdateBank(id string, bank *entity.Bank) (*string, error)
	UpdateBalanceBank(bank *entity.Bank, types string) (*string, error)
}
