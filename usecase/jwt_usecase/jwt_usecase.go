package jwt_usecase

import "github.com/golang-jwt/jwt/v4"

type JwtUsecase interface {
	GenerateToken(string, string, string) (string, error)
	ValidateToken(string) (*jwt.Token, error)
	ValidateTokenAndGetUserId(string) (string, error)
	ValidateTokenAndGetUser(string) (string, string, string, error)
}
