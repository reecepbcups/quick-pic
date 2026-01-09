package services

import (
	"context"

	"github.com/quickpic/server/internal/models"
	"github.com/quickpic/server/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
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
