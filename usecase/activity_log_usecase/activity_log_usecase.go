package activity_log_usecase

import "erp-service/model/dto"

type ActivityLogUsecase interface {
	AddActivity(dto.ActivityLog) (*string, error)
	GetActivity(limit int, offset int) dto.Response
}
