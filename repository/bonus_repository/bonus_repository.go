package bonus_repository

import (
	"erp-service/model/dto"
	"erp-service/model/entity"
)

type BonusRepository interface {
	GetAllBonus(types string, dateFrom string, dateTo string, limit int, offset int) (dto.BonusWithTotal, error)
	GetBonusByID(id string) (*entity.Bonus, error)
	AddBonus(bonus *entity.Bonus, balanceCoin float64) (*string, error)
	UpdateBonus(id string, types string, ammount float64, balanceCoin float64) (*string, error)
}
