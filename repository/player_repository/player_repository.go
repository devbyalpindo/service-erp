package player_repository

import (
	"erp-service/model/dto"
	"erp-service/model/entity"
)

type PlayerRepository interface {
	GetAllPlayer(playerID string, limit int, offset int) (dto.PlayerWithTotal, error)
	GetPlayerByID(string) (bool, error)
	GetPlayerBankByID(string) (*entity.BankPlayer, error)
	AddPlayer(*entity.Player) (*string, error)
	AddBankPlayer(*entity.BankPlayer) (*string, error)
	UpdatePlayer(*entity.Player) (*string, error)
	UpdateBankPlayer(*entity.BankPlayer) (*string, error)
	BulkInsertPlayer([]entity.Player) (string, error)
	BulkInsertBankPlayer([]entity.BankPlayer) (string, error)
}
