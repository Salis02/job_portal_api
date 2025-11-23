package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{db: db}
}

// CreateUser returns created user's id
func (r *Repo) CreateUser(ctx context.Context, name, email, passwordHash string) (string, error) {
	id := uuid.New().String()
	_, err := r.db.Exec(ctx, `
        INSERT INTO users (id, name, email, password_hash, created_at, updated_at)
        VALUES ($1,$2,$3,$4,NOW(),NOW())
    `, id, name, email, passwordHash)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *Repo) GetUserByEmail(ctx context.Context, email string) (id, name, emailOut, passwordHash string, err error) {
	row := r.db.QueryRow(ctx, `SELECT id, name, email, password_hash FROM users WHERE email=$1`, email)
	err = row.Scan(&id, &name, &emailOut, &passwordHash)
	return
}

func (r *Repo) GetUserById(ctx context.Context, id string) (name, email string, err error) {
	row := r.db.QueryRow(ctx, `SELECT name, email FROM users WHERE id=$1`, id)
	err = row.Scan(&name, &email)
	return
}

// Refresh token
func (r *Repo) StoreRefreshToken(ctx context.Context, userID, tokenHash, userAgent, ip string, expiresAt time.Time) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO refresh_tokens (id, user_id, token_hash, user_agent, ip_address, expires_at, created_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, NOW())
	`, userID, tokenHash, userAgent, ip, expiresAt)
	return err
}

func (r *Repo) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE refresh_tokens SET revoked = true WHERE token_hash = $1
	`, tokenHash)
	return err
}

// Validate refresh token exists, not revoked and not expired
func (r *Repo) ValidateRefreshToken(ctx context.Context, tokenHash string) (userID string, err error) {
	row := r.db.QueryRow(ctx, `
		SELECT user_id FROM refresh_tokens
		WHERE token_hash = $1 AND revoked = false AND expires_at > NOW()
	`, tokenHash)
	err = row.Scan(&userID)
	return
}
