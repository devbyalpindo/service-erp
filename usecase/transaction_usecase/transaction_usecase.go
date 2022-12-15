package transaction_usecase

import "erp-service/model/dto"

type TransactionUsecase interface {
	AddTransaction(userID string, body dto.AddTransaction) dto.Response
	GetAllTransaction(roleName string, limit int, offset int, dateFrom string, dateTo string, types string, status string) dto.Response
	UpdateTransaction(string, dto.UpdateTransactionPending) dto.Response
}
