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

func (repository *CoinRepositoryImpl) UpdateBalanceCoin(coin *entity.Coin, types string) (*string, float64, error) {
	coinResult := entity.Coin{}
	result := repository.DB.Where("coin_id = ?", coin.CoinID).Find(&coinResult)
	if result.RowsAffected == 0 {
		return nil, 0, gorm.ErrRecordNotFound
	}

	balance := coinResult.Balance

	switch types {
	case "MINUS":
		balance = balance - coin.Balance
	case "PLUS":
		balance = balance + coin.Balance
	default:
		return nil, 0, gorm.ErrRecordNotFound
	}

	results := repository.DB.Model(&coin).Where("coin_id = ?", coin.CoinID).Updates(map[string]interface{}{"balance": balance, "updated_at": coin.UpdatedAt})

	if results.RowsAffected == 0 {
		return nil, 0, gorm.ErrRecordNotFound
	}

	return &coin.CoinID, balance, nil
}

func (repository *CoinRepositoryImpl) TransactionCoin(coin entity.MutationCoin) (*string, error) {
	if err := repository.DB.Create(&coin).Error; err != nil {
		return nil, err
	}

	return &coin.MutationCoinID, nil
}

func (repository *CoinRepositoryImpl) GetMutation(types string, isTransactionBank bool, limit int, offset int, dateFrom string, dateTo string) (entity.CoinWithTotal, error) {

	mutation := []entity.MutationCoin{}
	var err error
	var totalData int64

	if len(types) > 0 {
		if isTransactionBank {
			err = repository.DB.Model(&entity.MutationCoin{}).Where("type = ? AND DATE(created_at) >= ? AND DATE(created_at) <= ? AND is_transaction_bank = ?", types, dateFrom, dateTo, isTransactionBank).Order("created_at DESC").Count(&totalData).Limit(limit).Offset(offset).Find(&mutation).Error
		} else {
			err = repository.DB.Model(&entity.MutationCoin{}).Where("type = ? AND DATE(created_at) >= ? AND DATE(created_at) <= ?", types, dateFrom, dateTo).Order("created_at DESC").Count(&totalData).Limit(limit).Offset(offset).Find(&mutation).Error
		}

	}

	if len(types) == 0 {
		if isTransactionBank {
			err = repository.DB.Model(&entity.MutationCoin{}).Where("DATE(created_at) >= ? AND DATE(created_at) <= ? AND is_transaction_bank = ?", dateFrom, dateTo, isTransactionBank).Order("created_at DESC").Count(&totalData).Limit(limit).Offset(offset).Find(&mutation).Error
		} else {
			err = repository.DB.Model(&entity.MutationCoin{}).Where("DATE(created_at) >= ? AND DATE(created_at) <= ?", dateFrom, dateTo).Order("created_at DESC").Count(&totalData).Limit(limit).Offset(offset).Find(&mutation).Error
		}

	}

	helper.PanicIfError(err)

	if len(mutation) <= 0 {
		resultError := entity.CoinWithTotal{
			Total:    0,
			Mutation: nil,
		}
		return resultError, gorm.ErrRecordNotFound
	}

	result := entity.CoinWithTotal{
		Total:    totalData,
		Mutation: mutation,
	}

	return result, nil
}
