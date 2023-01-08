package coin_repository

import (
	"erp-service/model/entity"
)

type CoinRepository interface {
	GetCoin() ([]entity.Coin, error)
	UpdateBalanceCoin(bank *entity.Coin, types string) (*string, float64, error)
	GetDetailCoin() (*entity.Coin, error)
	TransactionCoin(coin entity.MutationCoin) (*string, error)
	GetMutation(types string, IsTransactionBank bool, limit int, offset int, dateFrom string, dateTo string) (entity.CoinWithTotal, error)
}
