package player_repository

import (
	"erp-service/helper"
	"erp-service/model/entity"

	"gorm.io/gorm"
)

type PlayerRepositoryImpl struct {
	DB *gorm.DB
}

func NewPlayerRepository(DB *gorm.DB) PlayerRepository {
	return &PlayerRepositoryImpl{DB}
}

func (repository *PlayerRepositoryImpl) GetAllPlayer() ([]entity.Player, error) {
	player := []entity.Player{}
	err := repository.DB.Model(&entity.Player{}).Preload("BankPlayer").Find(&player).Error
	helper.PanicIfError(err)
	if len(player) <= 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return player, nil
}

func (repository *PlayerRepositoryImpl) GetPlayerByID(id string) (bool, error) {
	player := entity.Player{}
	result := repository.DB.Where("player_id = ?", id).Find(&player)

	if result.RowsAffected == 0 {
		return true, gorm.ErrRecordNotFound
	}

	return false, nil
}

func (repository *PlayerRepositoryImpl) GetPlayerBankByID(id string) (*entity.BankPlayer, error) {
	bankPlayer := entity.BankPlayer{}
	result := repository.DB.Where("bank_player_id = ?", id).Find(&bankPlayer)

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &bankPlayer, nil
}

func (repository *PlayerRepositoryImpl) AddPlayer(player *entity.Player) (*string, error) {
	if err := repository.DB.Create(&player).Error; err != nil {
		return nil, err
	}

	return &player.PlayerID, nil
}

func (repository *PlayerRepositoryImpl) AddBankPlayer(bankPlayer *entity.BankPlayer) (*string, error) {
	if err := repository.DB.Create(&bankPlayer).Error; err != nil {
		return nil, err
	}

	return &bankPlayer.BankPlayerID, nil
}
