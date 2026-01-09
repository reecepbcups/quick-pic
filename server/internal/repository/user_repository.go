package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/quickpic/server/internal/models"
)

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	user.ID = uuid.New()
	user.Username = strings.ToLower(user.Username)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := `
		INSERT INTO users (id, username, password_hash, public_key, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID.String(), user.Username, user.PasswordHash, user.PublicKey, user.CreatedAt, user.UpdatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint") || strings.Contains(err.Error(), "unique constraint") {
			return models.ErrUsernameExists
		}
		return err
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, username, password_hash, public_key, created_at, updated_at
		FROM users WHERE id = ?
	`

	var user models.User
	var idStr string
	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&idStr, &user.Username, &user.PasswordHash, &user.PublicKey, &user.CreatedAt, &user.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	user.ID, _ = uuid.Parse(idStr)
	return &user, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, username, password_hash, public_key, created_at, updated_at
		FROM users WHERE username = ?
	`

	var user models.User
	var idStr string
	err := r.db.QueryRowContext(ctx, query, strings.ToLower(username)).Scan(
		&idStr, &user.Username, &user.PasswordHash, &user.PublicKey, &user.CreatedAt, &user.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	user.ID, _ = uuid.Parse(idStr)
	return &user, nil
}

func (r *UserRepository) StoreRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at)
		VALUES (?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query, uuid.New().String(), userID.String(), tokenHash, expiresAt)
	return err
}

func (r *UserRepository) ValidateRefreshToken(ctx context.Context, tokenHash string) (uuid.UUID, error) {
	query := `
		SELECT user_id FROM refresh_tokens
		WHERE token_hash = ? AND expires_at > datetime('now')
	`

	var userIDStr string
	err := r.db.QueryRowContext(ctx, query, tokenHash).Scan(&userIDStr)

	if errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, models.ErrInvalidToken
	}
	if err != nil {
		return uuid.Nil, err
	}

	userID, _ := uuid.Parse(userIDStr)
	return userID, nil
}

func (r *UserRepository) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	query := `DELETE FROM refresh_tokens WHERE token_hash = ?`
	_, err := r.db.ExecContext(ctx, query, tokenHash)
	return err
}

func (r *UserRepository) DeleteAllRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = ?`
	_, err := r.db.ExecContext(ctx, query, userID.String())
	return err
}
