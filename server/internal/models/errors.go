package models

import "errors"

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrUsernameExists        = errors.New("username already exists")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrInvalidToken          = errors.New("invalid token")
	ErrTokenExpired          = errors.New("token expired")
	ErrFriendRequestNotFound = errors.New("friend request not found")
	ErrFriendRequestExists   = errors.New("friend request already exists")
	ErrAlreadyFriends        = errors.New("already friends")
	ErrNotFriends            = errors.New("not friends with this user")
	ErrMessageNotFound       = errors.New("message not found")
	ErrUnauthorized          = errors.New("unauthorized")
	ErrCannotAddSelf         = errors.New("cannot add yourself as a friend")
)
