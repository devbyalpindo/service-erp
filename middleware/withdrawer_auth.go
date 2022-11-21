package middleware

import (
	"erp-service/helper"
	"erp-service/usecase/jwt_usecase"

	"github.com/gin-gonic/gin"
)

func UserAuth(jwtUsecase jwt_usecase.JwtUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			resp := helper.ResponseError("You are unathorized", "Invalid token", 401)
			c.AbortWithStatusJSON(resp.StatusCode, resp)
			return
		}

		userId, username, role, err := jwtUsecase.ValidateTokenAndGetUser(authHeader)
		if err != nil {
			resp := helper.ResponseError("You are unathorized", "Invalid token", 401)
			c.AbortWithStatusJSON(resp.StatusCode, resp)
			return
		}

		if !(role == "ADMIN" || role == "WITHDRAWER" || role == "DEPOSITOR") {
			resp := helper.ResponseError("Forbidden Access", "You have no access to do this action", 403)
			c.AbortWithStatusJSON(resp.StatusCode, resp)
			return
		}

		c.Set("user_id", userId)
		c.Set("username", username)
		c.Set("role", role)
	}
}
