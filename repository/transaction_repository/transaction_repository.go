package transaction_repository

import (
	"erp-service/model/dto"
	"erp-service/model/entity"
)

type TransactionRepository interface {
	AddTransaction(trx *entity.Transaction, balanceCoin float32, balanceBank float32, typeTrx string) (string, error)
	GetAllTransaction(roleID string, limit int, offset int, dateFrom string, dateTo string, types string) (dto.TransactionWithTotal, error)
	GetDetailTransaction(string) (*entity.Transaction, error)
	UpdateTransaction(transactionID string, playerID string, bankPlayerID string, status string, balanceCoin float32) (string, error)
}
