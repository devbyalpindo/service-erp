package transaction_repository

import (
	"erp-service/model/dto"
	"erp-service/model/entity"
)

type TransactionRepository interface {
	AddTransaction(trx *entity.Transaction, balanceCoin float32, balanceBank float32) (string, error)
	GetAllTransaction(roleID string, limit int, offset int, dateFrom string, dateTo string) (dto.TransactionWithTotal, error)
}
