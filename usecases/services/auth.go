package services

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"auth_service/domain"
	"auth_service/repository"
	"auth_service/usecases"
)

type AuthService struct {
	authRepository      repository.AuthRepository
	emailService        usecases.EmailService
	AccessTokenExpTime  time.Duration
	RefreshTokenExpTime time.Duration
}

func NewAuthService(authRepository repository.AuthRepository, emailService usecases.EmailService, accessTokenExpTime int, refreshTokenExpTime int) usecases.AuthService {
	return &AuthService{authRepository: authRepository, emailService: emailService, AccessTokenExpTime: time.Hour * 24 * time.Duration(accessTokenExpTime), RefreshTokenExpTime: time.Hour * 24 * time.Duration(refreshTokenExpTime)}
}

func (s *AuthService) GenerateTokens(userID string, ip string) (domain.Tokens, error) {
	pairID := uuid.New()
	accessToken, err := s.generateAccessToken(userID, ip, pairID)
	if err != nil {
		return domain.Tokens{}, err
	}
	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		return domain.Tokens{}, err
	}
	err = s.postRefreshToken(refreshToken, userID, pairID)
	if err != nil {
		return domain.Tokens{}, err
	}
	return domain.Tokens{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *AuthService) RefreshTokens(userID string, ip string, accessToken string, refreshToken string) (domain.Tokens, error) {
	claims := &domain.CustomClaims{}

	_, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})

	if userID != claims.UserID {
		return domain.Tokens{}, domain.ErrUnauthorized
	}

	token, err := s.authRepository.GetToken(claims.PairID)
	if err != nil {
		return domain.Tokens{}, err
	}

	if time.Now().After(token.Expires) {
		err := s.authRepository.DeleteToken(claims.PairID)
		if err != nil {
			return domain.Tokens{}, err
		}
		return domain.Tokens{}, domain.ErrUnauthorized
	}

	if token.UserID != userID {
		return domain.Tokens{}, domain.ErrUnauthorized
	}

	err = bcrypt.CompareHashAndPassword([]byte(token.Token), []byte(refreshToken))
	if err != nil {
		return domain.Tokens{}, domain.ErrUnauthorized
	}

	if claims.Ip != ip {
		s.emailService.Send("Warning, ip changed")
	}

	err = s.authRepository.DeleteToken(claims.PairID)
	if err != nil {
		return domain.Tokens{}, err
	}

	return s.GenerateTokens(userID, ip)
}

func (s *AuthService) generateAccessToken(userID string, ip string, pairID uuid.UUID) (domain.Token, error) {
	expTime := time.Now().Add(s.AccessTokenExpTime)
	claims := domain.CustomClaims{
		UserID: userID,
		Ip:     ip,
		PairID: pairID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	signedToken, err := token.SignedString([]byte(os.Getenv("SECRET"))) //?
	if err != nil {
		return "", err
	}

	return domain.Token(signedToken), nil
}

func (s *AuthService) generateRefreshToken() (domain.Token, error) {
	data := make([]byte, 32)
	_, err := rand.Read(data)
	if err != nil {
		return "", err
	}
	token := base64.URLEncoding.EncodeToString(data)
	return domain.Token(token), nil
}

func (s *AuthService) postRefreshToken(refreshToken domain.Token, userID string, pairID uuid.UUID) error {
	encryptedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	expTime := time.Now().Add(s.RefreshTokenExpTime)
	err = s.authRepository.AddToken(domain.RefreshToken{
		UserID:    userID,
		Token:     domain.Token(encryptedToken),
		Expires:   expTime,
		TokenPair: pairID,
	})
	if err != nil {
		return err
	}
	return nil
}
