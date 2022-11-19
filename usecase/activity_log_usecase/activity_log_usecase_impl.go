package activity_log_usecase

import (
	"erp-service/helper"
	"erp-service/model/dto"
	"erp-service/model/entity"
	"erp-service/repository/activity_log_repository"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ActivityLogUsecaseImpl struct {
	ActivityLogRepository activity_log_repository.ActivityLogRepository
}

func NewActivityLogUsecase(activityLogRepository activity_log_repository.ActivityLogRepository) ActivityLogUsecase {
	return &ActivityLogUsecaseImpl{
		ActivityLogRepository: activityLogRepository,
	}
}

func (usecase *ActivityLogUsecaseImpl) AddActivity(body dto.ActivityLog) (*string, error) {

	createID := uuid.New().String()

	payloadLog := &entity.ActivityLog{
		ActivityID:    createID,
		UserID:        body.UserID,
		IsTransaction: body.IsTransaction,
		TransactionID: body.TransactionID,
		Description:   body.Description,
		CreatedAt:     time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt:     time.Now().Format("2006-01-02 15:04:05"),
	}

	log, err := usecase.ActivityLogRepository.AddActivity(payloadLog)
	helper.PanicIfError(err)

	return log, nil

}

func (usecase *ActivityLogUsecaseImpl) GetActivity(limit int, offset int) dto.Response {
	logList, err := usecase.ActivityLogRepository.GetActivity(limit, offset)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return helper.ResponseError("failed", "Data not found", 404)
	} else if err != nil {
		return helper.ResponseError("failed", err, 500)
	}
	helper.PanicIfError(err)
	response := []dto.ActivityLog{}
	for _, logs := range logList.ActivityLog {
		responseData := dto.ActivityLog{
			UserID:        logs.UserID,
			IsTransaction: logs.IsTransaction,
			TransactionID: logs.TransactionID,
			Description:   logs.Description,
			CreatedAt:     logs.CreatedAt,
		}
		response = append(response, responseData)
	}

	var result map[string]any = make(map[string]any)
	result["total"] = logList.Total
	result["log"] = response

	return helper.ResponseSuccess("ok", nil, result, 200)
}
