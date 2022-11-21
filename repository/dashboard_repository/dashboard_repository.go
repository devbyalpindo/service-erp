package dashboard_repository

import "erp-service/model/dto"

type DashboardRepository interface {
	GetTransactionValue(dateFrom string, dateTo string) ([]dto.TransactionValue, error)
	GetTopPlayerDeposit() ([]dto.TopPlayerDeposit, error)
	GetTopPlayerWithdraw() ([]dto.TopPlayerWithdraw, error)
}
