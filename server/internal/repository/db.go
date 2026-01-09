package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

// NewDB creates a new SQLite database connection
// Use ":memory:" for in-memory database (testing)
// Use a file path like "./quickpic.db" for persistent storage
func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db}, nil
}

func (db *DB) Migrate() error {
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
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	return nil
}

// Reset clears all data (for testing)
func (db *DB) Reset() error {
	tables := []string{"messages", "friendships", "friend_requests", "refresh_tokens", "users"}
	for _, table := range tables {
		if _, err := db.Exec("DELETE FROM " + table); err != nil {
			return err
		}
	}
	return nil
}
