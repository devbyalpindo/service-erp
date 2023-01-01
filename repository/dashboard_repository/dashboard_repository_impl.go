package dashboard_repository

import (
	"erp-service/model/dto"

	"gorm.io/gorm"
)

type DashboardRepositoryImpl struct {
	DB *gorm.DB
}

func NewDashboardRepository(DB *gorm.DB) DashboardRepository {
	return &DashboardRepositoryImpl{DB: DB}
}

func (repository *DashboardRepositoryImpl) GetTransactionValue(dateFrom string, dateTo string) ([]dto.TransactionValue, error) {
	trxValue := []dto.TransactionValue{}
	err := repository.DB.Table("transactions").Select("count(transactions.ammount) as total, type_transactions.type_transaction as value, sum(transactions.ammount) as total_ammount").Joins("inner join type_transactions on type_transactions.type_id = transactions.type_id").Where("DATE(transactions.created_at) >= ? AND DATE(transactions.created_at) <= ? AND transactions.status != ?", dateFrom, dateTo, "CANCELED").Group("type_transactions.type_transaction").Limit(10).Find(&trxValue).Error

	if err != nil {
		return trxValue, nil
	}

	return trxValue, nil
}

func (repository *DashboardRepositoryImpl) GetTopPlayerDeposit() ([]dto.TopPlayerDeposit, error) {
	topDepo := []dto.TopPlayerDeposit{}
	err := repository.DB.Table("transactions").Select("transactions.player_id, players.player_name, SUM(transactions.ammount) as total_deposit").Joins("inner join type_transactions on type_transactions.type_id = transactions.type_id").Joins("inner join players on transactions.player_id = players.player_id").Where("type_transactions.type_transaction = ?", "DEPOSIT").Group("transactions.player_id, players.player_name").Limit(10).Order("total_deposit DESC").Find(&topDepo).Error

	if err != nil {
		return topDepo, nil
	}

	return topDepo, nil
}

func (repository *DashboardRepositoryImpl) GetTopPlayerWithdraw() ([]dto.TopPlayerWithdraw, error) {
	topWD := []dto.TopPlayerWithdraw{}
	err := repository.DB.Table("transactions").Select("transactions.player_id, players.player_name, SUM(transactions.ammount) as total_withdraw").Joins("inner join type_transactions on type_transactions.type_id = transactions.type_id").Joins("inner join players on transactions.player_id = players.player_id").Where("type_transactions.type_transaction = ?", "WITHDRAW").Group("transactions.player_id, players.player_name").Limit(10).Order("total_withdraw DESC").Find(&topWD).Error

	if err != nil {
		return topWD, nil
	}

	return topWD, nil
}
