package bonus_repository

import (
	"erp-service/helper"
	"erp-service/model/dto"
	"erp-service/model/entity"
	"time"

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

func (repository *BonusRepositoryImpl) AddBonus(bonus *entity.Bonus, balanceCoin float64) (*string, error) {
	tx := repository.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&bonus).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Model(&entity.Coin{}).Where("coin_name = ?", "SALJU 88").Updates(map[string]interface{}{"balance": balanceCoin, "updated_at": time.Now().Format("2006-01-02 15:04:05")}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &bonus.BonusID, nil
}

func (repository *BonusRepositoryImpl) GetBonusByID(id string) (*entity.Bonus, error) {
	bonus := entity.Bonus{}
	result := repository.DB.Where("bonus_id = ?", id).Find(&bonus)

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &bonus, nil
}

func (repository *BonusRepositoryImpl) UpdateBonus(id string, types string, ammount float64, balanceCoin float64) (*string, error) {

	tx := repository.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&entity.Coin{}).Where("coin_name = ?", "SALJU 88").Updates(map[string]interface{}{"balance": balanceCoin, "updated_at": time.Now().Format("2006-01-02 15:04:05")}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Model(&entity.Bonus{}).Where("bonus_id = ?", id).Updates(map[string]interface{}{"type": types, "ammount": ammount, "updated_at": time.Now().Format("2006-01-02 15:04:05")}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &id, nil
}
