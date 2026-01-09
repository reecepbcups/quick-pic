package services

import (
	"context"

	"github.com/quickpic/server/internal/models"
	"github.com/quickpic/server/internal/storage"
)

type UserService struct {
	userRepo storage.UserRepo
}

func NewUserService(userRepo storage.UserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (*models.UserPublic, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	public := user.ToPublic()
	return &public, nil
}
