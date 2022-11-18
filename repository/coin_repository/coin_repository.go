package coin_repository

import (
	"erp-service/model/entity"
)

type CoinRepository interface {
	GetCoin() ([]entity.Coin, error)
	UpdateBalanceCoin(bank *entity.Coin, types string) (*string, error)
	GetDetailCoin() (*entity.Coin, error)
}
