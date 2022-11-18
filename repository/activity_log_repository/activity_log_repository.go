package activity_log_repository

import (
	"erp-service/model/dto"
	"erp-service/model/entity"
)

type ActivityLogRepository interface {
	AddActivity(*entity.ActivityLog) (*string, error)
	GetActivity(limit int, offset int) (dto.ActivityLogWithTotal, error)
}
