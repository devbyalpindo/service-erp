package coin_repository

import (
	"erp-service/helper"
	"erp-service/model/entity"

	"gorm.io/gorm"
)

type CoinRepositoryImpl struct {
	DB *gorm.DB
}

func NewCoinRepository(DB *gorm.DB) CoinRepository {
	return &CoinRepositoryImpl{DB: DB}
}

func (repository *CoinRepositoryImpl) GetCoin() ([]entity.Coin, error) {
	coin := []entity.Coin{}
	err := repository.DB.Model(&entity.Coin{}).Scan(&coin).Error
	helper.PanicIfError(err)
	if len(coin) <= 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return coin, nil
}

func (repository *CoinRepositoryImpl) GetDetailCoin() (*entity.Coin, error) {
	coin := entity.Coin{}
	result := repository.DB.First(&coin)

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &coin, nil
}

func (repository *CoinRepositoryImpl) UpdateBalanceCoin(coin *entity.Coin, types string) (*string, error) {
	coinResult := entity.Coin{}
	result := repository.DB.Where("coin_id = ?", coin.CoinID).Find(&coinResult)
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	balance := coinResult.Balance

	switch types {
	case "MINUS":
		balance = balance - coin.Balance
	case "PLUS":
		balance = balance + coin.Balance
	default:
		return nil, gorm.ErrRecordNotFound
	}

	results := repository.DB.Model(&coin).Where("coin_id = ?", coin.CoinID).Updates(map[string]interface{}{"balance": balance, "updated_at": coin.UpdatedAt})

	if results.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &coin.CoinID, nil
}
