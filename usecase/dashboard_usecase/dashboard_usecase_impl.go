package dashboard_usecase

import (
	"erp-service/helper"
	"erp-service/model/dto"
	"erp-service/repository/coin_repository"
	"erp-service/repository/dashboard_repository"
)

type DashboardUsecaseImpl struct {
	DashboardRepository dashboard_repository.DashboardRepository
	CoinRepository      coin_repository.CoinRepository
}

func NewDashboardUsecase(dashboardRepository dashboard_repository.DashboardRepository, coinRepository coin_repository.CoinRepository) DashboardUsecase {
	return &DashboardUsecaseImpl{
		DashboardRepository: dashboardRepository,
		CoinRepository:      coinRepository,
	}
}

func (usecase *DashboardUsecaseImpl) GetDashboard(dateFrom string, dateTo string) dto.Response {
	transactionValue, _ := usecase.DashboardRepository.GetTransactionValue(dateFrom, dateTo)

	topDepo, _ := usecase.DashboardRepository.GetTopPlayerDeposit()

	topWD, _ := usecase.DashboardRepository.GetTopPlayerWithdraw()

	getCoinGame, err := usecase.CoinRepository.GetDetailCoin()
	helper.PanicIfError(err)

	response := dto.Dashboard{
		TransactionValue:  transactionValue,
		TopPlayerDeposit:  topDepo,
		TopPlayerWithdraw: topWD,
		CoinGame:          getCoinGame,
	}

	var result map[string]any = make(map[string]any)
	result["dashboard"] = response

	return helper.ResponseSuccess("ok", nil, result, 200)
}
