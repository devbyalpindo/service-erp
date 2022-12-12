package bonus_repository

import (
	"erp-service/helper"
	"erp-service/model/dto"
	"erp-service/model/entity"

	"gorm.io/gorm"
)

type BonusRepositoryImpl struct {
	DB *gorm.DB
}

func NewBonusRepository(DB *gorm.DB) BonusRepository {
	return &BonusRepositoryImpl{DB: DB}
}

func (repository *BonusRepositoryImpl) GetAllBonus(types string, dateFrom string, dateTo string, limit int, offset int) (dto.BonusWithTotal, error) {
	bonus := []entity.Bonus{}
	var totalData int64
	var err error
	if len(types) > 0 {
		err = repository.DB.Model(&entity.Bonus{}).Count(&totalData).Limit(limit).Offset(offset).Where("DATE(created_at) >= ? AND DATE(created_at) <= ? AND type = ?", dateFrom, dateTo, types).Order("created_at DESC").Scan(&bonus).Error
	} else {
		err = repository.DB.Model(&entity.Bonus{}).Count(&totalData).Limit(limit).Offset(offset).Where("DATE(created_at) >= ? AND DATE(created_at) <= ?", dateFrom, dateTo).Order("created_at DESC").Scan(&bonus).Error
	}
	helper.PanicIfError(err)

	if len(bonus) <= 0 {
		resultError := dto.BonusWithTotal{
			Total:      0,
			TotalBonus: 0,
			Bonus:      nil,
		}
		return resultError, gorm.ErrRecordNotFound
	}
	var totalBonus float64

	for _, item := range bonus {
		totalBonus += item.Ammount
	}

	result := dto.BonusWithTotal{
		Total:      totalData,
		Bonus:      bonus,
		TotalBonus: totalBonus,
	}

	return result, nil
}

func (repository *BonusRepositoryImpl) AddBonus(bonus *entity.Bonus) (*string, error) {
	if err := repository.DB.Create(&bonus).Error; err != nil {
		return nil, err
	}

	return &bonus.BonusID, nil
}
