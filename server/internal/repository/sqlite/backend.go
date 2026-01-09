package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/quickpic/server/internal/models"
)

// Backend implements repository.Backend for SQLite storage
type Backend struct {
	db       *sql.DB
	users    *UserRepository
	friends  *FriendRepository
	messages *MessageRepository
}

// NewBackend creates a new SQLite backend
// Use ":memory:" for in-memory database (testing)
// Use a file path like "./quickpic.db" for persistent storage
func NewBackend(dataSourceName string) (*Backend, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	backend := &Backend{
		db: db,
	}
	backend.users = &UserRepository{db: db}
	backend.friends = &FriendRepository{db: db}
	backend.messages = &MessageRepository{db: db}

	// Run migrations
	if err := backend.migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return backend, nil
}

func (b *Backend) migrate() error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			user_number INTEGER UNIQUE,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			public_key TEXT NOT NULL,
			created_at DATETIME DEFAULT (datetime('now')),
			updated_at DATETIME DEFAULT (datetime('now'))
		)`,
		`CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)`,
		`CREATE TABLE IF NOT EXISTS refresh_tokens (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			token_hash TEXT UNIQUE NOT NULL,
			expires_at DATETIME NOT NULL,
			created_at DATETIME DEFAULT (datetime('now'))
		)`,
		`CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id)`,
		`CREATE TABLE IF NOT EXISTS friend_requests (
			id TEXT PRIMARY KEY,
			from_user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			to_user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			status TEXT NOT NULL DEFAULT 'pending',
			created_at DATETIME DEFAULT (datetime('now')),
			UNIQUE(from_user_id, to_user_id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_friend_requests_to_user ON friend_requests(to_user_id, status)`,
		`CREATE TABLE IF NOT EXISTS friendships (
			id TEXT PRIMARY KEY,
			user_a_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			user_b_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at DATETIME DEFAULT (datetime('now')),
			UNIQUE(user_a_id, user_b_id),
			CHECK(user_a_id < user_b_id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_friendships_user_a ON friendships(user_a_id)`,
		`CREATE INDEX IF NOT EXISTS idx_friendships_user_b ON friendships(user_b_id)`,
		`CREATE TABLE IF NOT EXISTS messages (
			id TEXT PRIMARY KEY,
			from_user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			to_user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			encrypted_content BLOB NOT NULL,
			content_type TEXT NOT NULL,
			signature TEXT NOT NULL,
			created_at DATETIME DEFAULT (datetime('now'))
		)`,
		`CREATE INDEX IF NOT EXISTS idx_messages_to_user ON messages(to_user_id, created_at)`,
	}

	for _, migration := range migrations {
		if _, err := b.db.Exec(migration); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	return nil
}

func (b *Backend) Users() *UserRepository {
	return b.users
}

func (b *Backend) Friends() *FriendRepository {
	return b.friends
}

func (b *Backend) Messages() *MessageRepository {
	return b.messages
}

func (b *Backend) Close() error {
	return b.db.Close()
}

func (b *Backend) Name() string {
	return "sqlite"
}

// Reset clears all data (for testing)
func (b *Backend) Reset() error {
	tables := []string{"messages", "friendships", "friend_requests", "refresh_tokens", "users"}
	for _, table := range tables {
		if _, err := b.db.Exec("DELETE FROM " + table); err != nil {
			return err
		}
	}
	return nil
}

// ============ UserRepository ============

type UserRepository struct {
	db *sql.DB
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	user.ID = uuid.New()
	user.Username = strings.ToLower(user.Username)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Get next user number
	var maxUserNumber sql.NullInt64
	err := r.db.QueryRowContext(ctx, "SELECT MAX(user_number) FROM users").Scan(&maxUserNumber)
	if err != nil {
		return err
	}
	if maxUserNumber.Valid {
		user.UserNumber = maxUserNumber.Int64 + 1
	} else {
		user.UserNumber = 1
	}

	query := `
		INSERT INTO users (id, user_number, username, password_hash, public_key, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err = r.db.ExecContext(ctx, query,
		user.ID.String(), user.UserNumber, user.Username, user.PasswordHash, user.PublicKey, user.CreatedAt, user.UpdatedAt)

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
		SELECT id, user_number, username, password_hash, public_key, created_at, updated_at
		FROM users WHERE id = ?
	`

	var user models.User
	var idStr string
	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&idStr, &user.UserNumber, &user.Username, &user.PasswordHash, &user.PublicKey, &user.CreatedAt, &user.UpdatedAt)

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
		SELECT id, user_number, username, password_hash, public_key, created_at, updated_at
		FROM users WHERE username = ?
	`

	var user models.User
	var idStr string
	err := r.db.QueryRowContext(ctx, query, strings.ToLower(username)).Scan(
		&idStr, &user.UserNumber, &user.Username, &user.PasswordHash, &user.PublicKey, &user.CreatedAt, &user.UpdatedAt)

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

// ============ FriendRepository ============

type FriendRepository struct {
	db *sql.DB
}

func (r *FriendRepository) CreateFriendRequest(ctx context.Context, fromUserID, toUserID uuid.UUID) (*models.FriendRequest, error) {
	// Check if request already exists
	existingQuery := `
		SELECT id FROM friend_requests
		WHERE from_user_id = ? AND to_user_id = ? AND status = 'pending'
	`
	var existingID string
	err := r.db.QueryRowContext(ctx, existingQuery, fromUserID.String(), toUserID.String()).Scan(&existingID)
	if err == nil {
		return nil, models.ErrFriendRequestExists
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// Check if reverse request exists (they already sent us a request)
	err = r.db.QueryRowContext(ctx, existingQuery, toUserID.String(), fromUserID.String()).Scan(&existingID)
	if err == nil {
		return nil, models.ErrFriendRequestExists
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// Check if already friends
	if areFriends, _ := r.AreFriends(ctx, fromUserID, toUserID); areFriends {
		return nil, models.ErrAlreadyFriends
	}

	request := &models.FriendRequest{
		ID:         uuid.New(),
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Status:     models.FriendRequestPending,
		CreatedAt:  time.Now(),
	}

	query := `
		INSERT INTO friend_requests (id, from_user_id, to_user_id, status, created_at)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err = r.db.ExecContext(ctx, query,
		request.ID.String(), request.FromUserID.String(), request.ToUserID.String(), request.Status, request.CreatedAt)
	if err != nil {
		return nil, err
	}

	return request, nil
}

func (r *FriendRepository) GetPendingRequests(ctx context.Context, userID uuid.UUID) ([]models.FriendRequestWithUser, error) {
	query := `
		SELECT fr.id, fr.from_user_id, fr.to_user_id, fr.status, fr.created_at,
		       u.id, u.username, u.public_key
		FROM friend_requests fr
		JOIN users u ON u.id = fr.from_user_id
		WHERE fr.to_user_id = ? AND fr.status = 'pending'
		ORDER BY fr.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []models.FriendRequestWithUser
	for rows.Next() {
		var req models.FriendRequestWithUser
		var idStr, fromUserIDStr, toUserIDStr, fromUserIDStr2 string
		err := rows.Scan(
			&idStr, &fromUserIDStr, &toUserIDStr, &req.Status, &req.CreatedAt,
			&fromUserIDStr2, &req.FromUser.Username, &req.FromUser.PublicKey)
		if err != nil {
			return nil, err
		}
		req.ID, _ = uuid.Parse(idStr)
		req.FromUserID, _ = uuid.Parse(fromUserIDStr)
		req.ToUserID, _ = uuid.Parse(toUserIDStr)
		req.FromUser.ID, _ = uuid.Parse(fromUserIDStr2)
		requests = append(requests, req)
	}

	return requests, rows.Err()
}

func (r *FriendRepository) GetFriendRequest(ctx context.Context, requestID uuid.UUID) (*models.FriendRequest, error) {
	query := `
		SELECT id, from_user_id, to_user_id, status, created_at
		FROM friend_requests WHERE id = ?
	`

	var req models.FriendRequest
	var idStr, fromUserIDStr, toUserIDStr string
	err := r.db.QueryRowContext(ctx, query, requestID.String()).Scan(
		&idStr, &fromUserIDStr, &toUserIDStr, &req.Status, &req.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrFriendRequestNotFound
	}
	if err != nil {
		return nil, err
	}

	req.ID, _ = uuid.Parse(idStr)
	req.FromUserID, _ = uuid.Parse(fromUserIDStr)
	req.ToUserID, _ = uuid.Parse(toUserIDStr)
	return &req, nil
}

func (r *FriendRepository) UpdateFriendRequestStatus(ctx context.Context, requestID uuid.UUID, status models.FriendRequestStatus) error {
	query := `UPDATE friend_requests SET status = ? WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, status, requestID.String())
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return models.ErrFriendRequestNotFound
	}

	return nil
}

func (r *FriendRepository) CreateFriendship(ctx context.Context, userAID, userBID uuid.UUID) error {
	// Ensure user_a_id < user_b_id for the unique constraint
	if userAID.String() > userBID.String() {
		userAID, userBID = userBID, userAID
	}

	query := `
		INSERT OR IGNORE INTO friendships (id, user_a_id, user_b_id, created_at)
		VALUES (?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query, uuid.New().String(), userAID.String(), userBID.String(), time.Now())
	return err
}

func (r *FriendRepository) AreFriends(ctx context.Context, userAID, userBID uuid.UUID) (bool, error) {
	// Ensure user_a_id < user_b_id
	if userAID.String() > userBID.String() {
		userAID, userBID = userBID, userAID
	}

	query := `SELECT 1 FROM friendships WHERE user_a_id = ? AND user_b_id = ?`
	var exists int
	err := r.db.QueryRowContext(ctx, query, userAID.String(), userBID.String()).Scan(&exists)

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *FriendRepository) GetFriends(ctx context.Context, userID uuid.UUID) ([]models.Friend, error) {
	query := `
		SELECT u.id, u.username, u.public_key, f.created_at
		FROM friendships f
		JOIN users u ON (
			(f.user_a_id = ? AND u.id = f.user_b_id) OR
			(f.user_b_id = ? AND u.id = f.user_a_id)
		)
		WHERE f.user_a_id = ? OR f.user_b_id = ?
		ORDER BY u.username
	`

	userIDStr := userID.String()
	rows, err := r.db.QueryContext(ctx, query, userIDStr, userIDStr, userIDStr, userIDStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []models.Friend
	for rows.Next() {
		var f models.Friend
		var userIDStr string
		if err := rows.Scan(&userIDStr, &f.Username, &f.PublicKey, &f.Since); err != nil {
			return nil, err
		}
		f.UserID, _ = uuid.Parse(userIDStr)
		friends = append(friends, f)
	}

	return friends, rows.Err()
}

// ============ MessageRepository ============

type MessageRepository struct {
	db *sql.DB
}

func (r *MessageRepository) Create(ctx context.Context, msg *models.Message) error {
	msg.ID = uuid.New()
	msg.CreatedAt = time.Now()

	query := `
		INSERT INTO messages (id, from_user_id, to_user_id, encrypted_content, content_type, signature, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		msg.ID.String(), msg.FromUserID.String(), msg.ToUserID.String(),
		msg.EncryptedContent, msg.ContentType, msg.Signature, msg.CreatedAt)

	return err
}

func (r *MessageRepository) GetPendingMessages(ctx context.Context, userID uuid.UUID) ([]models.MessageWithSender, error) {
	query := `
		SELECT m.id, m.from_user_id, m.to_user_id, m.encrypted_content, m.content_type, m.signature, m.created_at,
		       u.username, u.public_key
		FROM messages m
		JOIN users u ON u.id = m.from_user_id
		WHERE m.to_user_id = ?
		ORDER BY m.created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.MessageWithSender
	for rows.Next() {
		var msg models.MessageWithSender
		var idStr, fromUserIDStr, toUserIDStr string
		err := rows.Scan(
			&idStr, &fromUserIDStr, &toUserIDStr,
			&msg.EncryptedContent, &msg.ContentType, &msg.Signature, &msg.CreatedAt,
			&msg.FromUsername, &msg.FromPublicKey)
		if err != nil {
			return nil, err
		}
		msg.ID, _ = uuid.Parse(idStr)
		msg.FromUserID, _ = uuid.Parse(fromUserIDStr)
		msg.ToUserID, _ = uuid.Parse(toUserIDStr)
		messages = append(messages, msg)
	}

	return messages, rows.Err()
}

func (r *MessageRepository) GetByID(ctx context.Context, messageID uuid.UUID) (*models.Message, error) {
	query := `
		SELECT id, from_user_id, to_user_id, encrypted_content, content_type, signature, created_at
		FROM messages WHERE id = ?
	`

	var msg models.Message
	var idStr, fromUserIDStr, toUserIDStr string
	err := r.db.QueryRowContext(ctx, query, messageID.String()).Scan(
		&idStr, &fromUserIDStr, &toUserIDStr,
		&msg.EncryptedContent, &msg.ContentType, &msg.Signature, &msg.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrMessageNotFound
	}
	if err != nil {
		return nil, err
	}

	msg.ID, _ = uuid.Parse(idStr)
	msg.FromUserID, _ = uuid.Parse(fromUserIDStr)
	msg.ToUserID, _ = uuid.Parse(toUserIDStr)
	return &msg, nil
}

func (r *MessageRepository) Delete(ctx context.Context, messageID uuid.UUID) error {
	query := `DELETE FROM messages WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, messageID.String())
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return models.ErrMessageNotFound
	}

	return nil
}

func (r *MessageRepository) DeleteOldMessages(ctx context.Context, olderThan time.Duration) (int64, error) {
	query := `DELETE FROM messages WHERE created_at < ?`
	result, err := r.db.ExecContext(ctx, query, time.Now().Add(-olderThan))
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
