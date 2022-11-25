package player_usecase

import "erp-service/model/dto"

type PlayerUsecase interface {
	GetAllPlayer() dto.Response
	AddPlayer(dto.AddPlayer) dto.Response
	AddPlayerBank(dto.AddBankPlayer) dto.Response
}
