package coin_usecase

import "erp-service/model/dto"

type CoinUsecase interface {
	GetCoin() dto.Response
	GetDetailCoin() dto.Response
	UpdateCoinBalance(dto.CoinUpdateBalance) dto.Response
	GetMutation(dto.GetMutationCoin) dto.Response
}
