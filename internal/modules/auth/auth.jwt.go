package auth

import (
	"errors"
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

func (j *JWTManager) ParseAccessToken(tokenStr string) (string, error) {
	tok, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.secretKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := tok.Claims.(jwt.MapClaims); ok && tok.Valid {
		sub, ok := claims["sub"].(string)
		if !ok {
			return "", errors.New("invalid subject")
		}
		return sub, nil
	}
	return "", errors.New("invalid token")
}
func (j *JWTManager) GenerateRefreshToken() (string, time.Time, error) {
	// Generate a secure random string (here using jwt with random id)
	exp := time.Now().Add(j.refreshTTL)
	claims := jwt.MapClaims{
		"jti":  jwt.NewNumericDate(time.Now()).String(), // Not important just unique
		"exp":  exp.Unix(),
		"liat": time.Now().Unix(),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := tok.SignedString(j.secretKey)
	return s, exp, err
}
