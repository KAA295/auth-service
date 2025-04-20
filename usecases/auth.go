package usecases

import (
	"auth_service/domain"
)

type AuthService interface {
	GenerateTokens(userID string, ip string) (domain.Tokens, error)
	RefreshTokens(userID string, ip string, accessToken string, refreshToken string) (domain.Tokens, error)
}
