package player_repository

import (
	"erp-service/helper"
	"erp-service/model/entity"
	"time"

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

func (repository *PlayerRepositoryImpl) UpdatePlayer(body *entity.Player) (*string, error) {

	result := repository.DB.Model(&body).Where("player_id = ?", body.PlayerID).Updates(entity.Player{PlayerName: body.PlayerName, UpdatedAt: time.Now().Format("2006-01-02 15:04:05")})
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &body.PlayerID, nil
}

func (repository *PlayerRepositoryImpl) UpdateBankPlayer(body *entity.BankPlayer) (*string, error) {

	result := repository.DB.Model(&body).Where("bank_player_id = ?", body.BankPlayerID).Updates(entity.BankPlayer{BankName: body.BankName, AccountName: body.AccountName, AccountNumber: body.AccountNumber, Category: body.Category, UpdatedAt: time.Now().Format("2006-01-02 15:04:05")})
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &body.BankPlayerID, nil
}

func (repository *PlayerRepositoryImpl) BulkInsertPlayer(player []entity.Player) (string, error) {
	if err := repository.DB.Create(&player).Error; err != nil {
		return "", err
	}

	return "ok", nil
}

func (repository *PlayerRepositoryImpl) BulkInsertBankPlayer(bank []entity.BankPlayer) (string, error) {
	if err := repository.DB.Create(&bank).Error; err != nil {
		return "", err
	}

	return "ok", nil
}
