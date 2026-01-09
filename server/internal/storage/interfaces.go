package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/quickpic/server/internal/models"
)

// UserRepo defines the interface for user data operations
type UserRepo interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	StoreRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error
	ValidateRefreshToken(ctx context.Context, tokenHash string) (uuid.UUID, error)
	DeleteRefreshToken(ctx context.Context, tokenHash string) error
	DeleteAllRefreshTokens(ctx context.Context, userID uuid.UUID) error
}

// FriendRepo defines the interface for friend-related operations
type FriendRepo interface {
	CreateFriendRequest(ctx context.Context, fromUserID, toUserID uuid.UUID) (*models.FriendRequest, error)
	GetPendingRequests(ctx context.Context, userID uuid.UUID) ([]models.FriendRequestWithUser, error)
	GetFriendRequest(ctx context.Context, requestID uuid.UUID) (*models.FriendRequest, error)
	UpdateFriendRequestStatus(ctx context.Context, requestID uuid.UUID, status models.FriendRequestStatus) error
	CreateFriendship(ctx context.Context, userAID, userBID uuid.UUID) error
	AreFriends(ctx context.Context, userAID, userBID uuid.UUID) (bool, error)
	GetFriends(ctx context.Context, userID uuid.UUID) ([]models.Friend, error)
}

// MessageRepo defines the interface for message operations
type MessageRepo interface {
	Create(ctx context.Context, msg *models.Message) error
	GetPendingMessages(ctx context.Context, userID uuid.UUID) ([]models.MessageWithSender, error)
	GetByID(ctx context.Context, messageID uuid.UUID) (*models.Message, error)
	Delete(ctx context.Context, messageID uuid.UUID) error
	DeleteOldMessages(ctx context.Context, olderThan time.Duration) (int64, error)
}

// Backend represents a complete storage backend with all repositories
// Use type assertions to get the concrete repository types from backends
type Backend interface {
	// Close closes any underlying connections
	Close() error
	// Name returns the backend name (e.g., "sqlite", "blockchain")
	Name() string
}

// SQLiteBackend is implemented by the SQLite backend
type SQLiteBackend interface {
	Backend
	Users() UserRepo
	Friends() FriendRepo
	Messages() MessageRepo
}

// BlockchainBackend is implemented by the blockchain backend
type BlockchainBackend interface {
	Backend
	Users() UserRepo
	Friends() FriendRepo
	Messages() MessageRepo
}

// Repositories holds all repository interfaces for a backend
type Repositories struct {
	Users    UserRepo
	Friends  FriendRepo
	Messages MessageRepo
}
