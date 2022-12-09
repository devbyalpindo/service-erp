package bonus_usecase

import "erp-service/model/dto"

type BonusUsecase interface {
	GetAllBonus(types string, dateFrom string, dateTo string, limit int, offset int) dto.Response
	AddBonus(dto.BonusAdd) dto.Response
}
