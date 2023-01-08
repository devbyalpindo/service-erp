package bank_repository

import (
	"erp-service/model/entity"
)

type BankRepository interface {
	AddBank(*entity.Bank) (*string, error)
	GetAllBank() ([]entity.Bank, error)
	GetDetailBank(id string) (*entity.Bank, error)
	UpdateBank(id string, bank *entity.Bank) (*string, error)
	UpdateBalanceBank(bank *entity.Bank, types string) (*string, float64, error)
	TransferToBank(idFrom string, balanceBankFrom float64, idBankTo string, balanceBankTo float64, amount float64, adminFee float64, nameBankfrom string, nameBankTo string) (*string, error)
	TransactionBank(bank entity.MutationBank) (*string, error)
	GetMutation(bankID string, types string, IsTransactionBank bool, limit int, offset int, dateFrom string, dateTo string) (entity.BankithTotal, error)
}
