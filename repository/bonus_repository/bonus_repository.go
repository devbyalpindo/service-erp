package bonus_repository

import (
	"erp-service/model/dto"
	"erp-service/model/entity"
)

type BonusRepository interface {
	GetAllBonus(types string, dateFrom string, dateTo string, limit int, offset int) (dto.BonusWithTotal, error)
	AddBonus(*entity.Bonus) (*string, error)
}
