package postgres

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"auth_service/domain"
)

type AuthRepository struct {
	DB *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{DB: db}
}

func (repo *AuthRepository) GetToken(pairID uuid.UUID) (domain.RefreshToken, error) {
	query := "SELECT token, user_id, expiration_time FROM tokens WHERE pair_id = $1"
	var refreshToken domain.RefreshToken
	err := repo.DB.QueryRow(query, pairID).Scan(&refreshToken.Token, &refreshToken.UserID, &refreshToken.Expires)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.RefreshToken{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.RefreshToken{}, err
	}

	return refreshToken, nil
}

func (repo *AuthRepository) AddToken(refreshToken domain.RefreshToken) error {
	query := "INSERT INTO tokens (token, user_id, expiration_time, pair_id) VALUES ($1, $2, $3, $4)"
	_, err := repo.DB.Exec(query, refreshToken.Token, refreshToken.UserID, refreshToken.Expires, refreshToken.TokenPair)
	if err != nil {
		return err
	}
	return nil
}

func (repo *AuthRepository) DeleteToken(pairID uuid.UUID) error {
	query := "DELETE FROM tokens WHERE pair_id = $1"
	_, err := repo.DB.Exec(query, pairID)
	if err != nil {
		return err
	}
	return nil
}
