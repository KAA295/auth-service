package repository

import (
	"auth_service/domain"

	"github.com/google/uuid"
)

type AuthRepository interface {
	GetToken(pairID uuid.UUID) (domain.RefreshToken, error)
	AddToken(refreshToken domain.RefreshToken) error
	DeleteToken(pairID uuid.UUID) error
}
