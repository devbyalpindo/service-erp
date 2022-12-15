package player_repository

import (
	"erp-service/model/entity"
)

type PlayerRepository interface {
	GetAllPlayer() ([]entity.Player, error)
	GetPlayerByID(string) (bool, error)
	GetPlayerBankByID(string) (*entity.BankPlayer, error)
	AddPlayer(*entity.Player) (*string, error)
	AddBankPlayer(*entity.BankPlayer) (*string, error)
	UpdatePlayer(*entity.Player) (*string, error)
	UpdateBankPlayer(*entity.BankPlayer) (*string, error)
}
