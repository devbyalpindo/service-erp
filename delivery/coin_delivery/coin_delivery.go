package coin_delivery

import "github.com/gin-gonic/gin"

type CoinDelivery interface {
	GetCoin(*gin.Context)
	GetDetailCoin(*gin.Context)
	UpdateCoinBalance(*gin.Context)
	GetMutation(*gin.Context)
}
