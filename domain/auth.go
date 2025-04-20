package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Token string

type AccessToken struct {
	Token Token
}

type RefreshToken struct {
	Token     Token
	UserID    string
	Expires   time.Time
	TokenPair uuid.UUID
}

type Tokens struct {
	AccessToken  Token
	RefreshToken Token
}

type CustomClaims struct {
	UserID string
	Ip     string
	PairID uuid.UUID
	jwt.RegisteredClaims
}
