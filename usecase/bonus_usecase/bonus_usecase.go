package bonus_usecase

import "erp-service/model/dto"

type BonusUsecase interface {
	GetAllBonus(types string, dateFrom string, dateTo string, limit int, offset int) dto.Response
	AddBonus(dto.BonusAdd) dto.Response
	UpdateBonus(id string, body dto.BonusAdd) dto.Response
}
