package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // Never expose in JSON
	PublicKey    string    `json:"public_key"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserPublic struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	PublicKey string    `json:"public_key"`
}

func (u *User) ToPublic() UserPublic {
	return UserPublic{
		ID:        u.ID,
		Username:  u.Username,
		PublicKey: u.PublicKey,
	}
}

type RegisterRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=32"`
	Password  string `json:"password" binding:"required,min=8"`
	PublicKey string `json:"public_key" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	ExpiresIn    int64      `json:"expires_in"`
	User         UserPublic `json:"user"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
