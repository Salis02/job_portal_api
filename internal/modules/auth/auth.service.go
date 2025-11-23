package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo             *Repo
	jwt              *JWTManager
	accessTTLSeconds int64
}

func NewService(r *Repo, jwt *JWTManager) *Service {
	return &Service{
		repo:             r,
		jwt:              jwt,
		accessTTLSeconds: int64(jwt.accessTTL.Seconds()),
	}
}

func (s *Service) Register(ctx context.Context, req *RegisterRequest) (string, error) {
	// Normalize email
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// Check valid email
	if !strings.Contains(email, "@") {
		return "", errors.New("invalid email")
	}

	// Rules minimal password
	if len(req.Password) < 6 {
		return "", errors.New("password too short")
	}

	//Hash Password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return s.repo.CreateUser(ctx, req.Name, email, string(hash))
}

func (s *Service) Login(ctx context.Context, req *LoginRequest, userAgent, ip string) (*AuthResponse, error) {
	id, _, _, passwordHash, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate access token
	accessToken, exp, err := s.jwt.GenerateAccessToken(id)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshTokenPlain, refreshExp, err := s.jwt.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// Hash refresh token before storing to db
	h := sha256.Sum256([]byte(refreshTokenPlain))
	tokenHash := hex.EncodeToString(h[:])

	// Store refresh token hashed
	if err := s.repo.StoreRefreshToken(ctx, id, tokenHash, userAgent, ip, refreshExp); err != nil {
		return nil, err
	}

	resp := &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenPlain,
		TokenType:    "Bearer",
		ExpiresIn:    int64(exp.Sub(time.Now()).Seconds()),
	}
	return resp, nil
}

func (s *Service) Refresh(ctx context.Context, refreshTokenPlain string, userAgent, ip string) (*AuthResponse, error) {
	// hash input token
	h := sha256.Sum256([]byte(refreshTokenPlain))
	tokenHash := hex.EncodeToString(h[:])

	userID, err := s.repo.ValidateRefreshToken(ctx, tokenHash)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// revoke old token (rotate)
	_ = s.repo.RevokeRefreshToken(ctx, tokenHash)

	// issue new access + refresh
	accessToken, exp, err := s.jwt.GenerateAccessToken(userID)
	if err != nil {
		return nil, err
	}
	refreshTokenPlainNew, refreshExp, err := s.jwt.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}
	h2 := sha256.Sum256([]byte(refreshTokenPlainNew))
	tokenHash2 := hex.EncodeToString(h2[:])
	if err := s.repo.StoreRefreshToken(ctx, userID, tokenHash2, userAgent, ip, refreshExp); err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenPlainNew,
		TokenType:    "Bearer",
		ExpiresIn:    int64(exp.Sub(time.Now()).Seconds()),
	}, nil
}

func (s *Service) Logout(ctx context.Context, refreshTokenPlain string) error {
	h := sha256.Sum256([]byte(refreshTokenPlain))
	tokenHash := hex.EncodeToString(h[:])
	return s.repo.RevokeRefreshToken(ctx, tokenHash)
}
