package models

import (
	"time"

	"github.com/google/uuid"
)

type FriendRequestStatus string

const (
	FriendRequestPending  FriendRequestStatus = "pending"
	FriendRequestAccepted FriendRequestStatus = "accepted"
	FriendRequestRejected FriendRequestStatus = "rejected"
)

type FriendRequest struct {
	ID         uuid.UUID           `json:"id"`
	FromUserID uuid.UUID           `json:"from_user_id"`
	ToUserID   uuid.UUID           `json:"to_user_id"`
	Status     FriendRequestStatus `json:"status"`
	CreatedAt  time.Time           `json:"created_at"`
}

type FriendRequestWithUser struct {
	FriendRequest
	FromUser UserPublic `json:"from_user"`
}

type Friendship struct {
	ID        uuid.UUID `json:"id"`
	UserAID   uuid.UUID `json:"user_a_id"`
	UserBID   uuid.UUID `json:"user_b_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Friend struct {
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	PublicKey string    `json:"public_key"`
	Since     time.Time `json:"since"`
}

type SendFriendRequestRequest struct {
	Username string `json:"username" binding:"required"`
}

type FriendRequestActionRequest struct {
	RequestID uuid.UUID `json:"request_id" binding:"required"`
}
