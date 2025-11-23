package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secretKey  []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewJWTManager(secret string, accessTTL, refreshTTL time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:  []byte(secret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (j *JWTManager) GenerateAccessToken(userID string) (string, time.Time, error) {
	exp := time.Now().Add(j.accessTTL)
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": exp.Unix(),
		"iat": time.Now().Unix(),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := tok.SignedString(j.secretKey)
	return s, exp, err
}
