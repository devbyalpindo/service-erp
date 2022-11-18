package activity_log_repository

import (
	"erp-service/helper"
	"erp-service/model/dto"
	"erp-service/model/entity"

	"gorm.io/gorm"
)

type ActivityLogRepositoryImpl struct {
	DB *gorm.DB
}

func NewActivityLogRepository(DB *gorm.DB) ActivityLogRepository {
	return &ActivityLogRepositoryImpl{DB: DB}
}

func (repository *ActivityLogRepositoryImpl) AddActivity(log *entity.ActivityLog) (*string, error) {
	if err := repository.DB.Create(&log).Error; err != nil {
		return nil, err
	}

	return &log.ActivityID, nil
}

func (repository *ActivityLogRepositoryImpl) GetActivity(limit int, offset int) (dto.ActivityLogWithTotal, error) {
	log := []entity.ActivityLog{}
	var totalData int64
	err := repository.DB.Model(&entity.ActivityLog{}).Count(&totalData).Limit(limit).Offset(offset).Order("created_at DESC").Scan(&log).Error
	helper.PanicIfError(err)
	if len(log) <= 0 {
		resultError := dto.ActivityLogWithTotal{
			Total:       0,
			ActivityLog: nil,
		}
		return resultError, gorm.ErrRecordNotFound
	}

	result := dto.ActivityLogWithTotal{
		Total:       totalData,
		ActivityLog: log,
	}

	return result, nil
}
