package player_usecase

import "erp-service/model/dto"

type PlayerUsecase interface {
	GetAllPlayer() dto.Response
	AddPlayer(dto.AddPlayer) dto.Response
	AddPlayerBank(dto.AddBankPlayer) dto.Response
	UpdatePlayer(dto.UpdatePlayer) dto.Response
	UpdateBankPlayer(dto.UpdateBankPlayer) dto.Response
	BulkInsertPlayer([]dto.BulkInsertPlayer) dto.Response
	BulkInsertBankPlayer([]dto.BulkInsertBankPlayer) dto.Response
}
