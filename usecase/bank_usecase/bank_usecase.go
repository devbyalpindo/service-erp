package bank_usecase

import "erp-service/model/dto"

type BankUsecase interface {
	AddBank(dto.BankAdd) dto.Response
	GetAllBank() dto.Response
	UpdateBank(string, dto.BankUpdate) dto.Response
	UpdateBankBalance(dto.BankUpdateBalance) dto.Response
}
