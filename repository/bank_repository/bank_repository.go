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
	TransferToBank(idFrom string, balanceBankFrom float32, idBankTo string, balanceBankTo float32, amount float32) (*string, error)
	TransactionBank(bank entity.MutationBank) (*string, error)
	GetMutation(bankID string, types string, limit int, offset int, dateFrom string, dateTo string) (entity.BankithTotal, error)
}
