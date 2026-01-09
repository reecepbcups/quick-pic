package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/quickpic/server/internal/models"
	"github.com/quickpic/server/internal/storage"
)

type FriendService struct {
	friendRepo storage.FriendRepo
	userRepo   storage.UserRepo
}

func NewFriendService(friendRepo storage.FriendRepo, userRepo storage.UserRepo) *FriendService {
	return &FriendService{
		friendRepo: friendRepo,
		userRepo:   userRepo,
	}
}

func (s *FriendService) SendFriendRequest(ctx context.Context, fromUserID uuid.UUID, toUsername string) (*models.FriendRequest, error) {
	toUser, err := s.userRepo.GetByUsername(ctx, toUsername)
	if err != nil {
		return nil, err
	}

	if fromUserID == toUser.ID {
		return nil, models.ErrCannotAddSelf
	}

	return s.friendRepo.CreateFriendRequest(ctx, fromUserID, toUser.ID)
}

func (s *FriendService) GetPendingRequests(ctx context.Context, userID uuid.UUID) ([]models.FriendRequestWithUser, error) {
	return s.friendRepo.GetPendingRequests(ctx, userID)
}

func (s *FriendService) AcceptFriendRequest(ctx context.Context, userID uuid.UUID, requestID uuid.UUID) error {
	request, err := s.friendRepo.GetFriendRequest(ctx, requestID)
	if err != nil {
		return err
	}

	// Verify the current user is the recipient
	if request.ToUserID != userID {
		return models.ErrUnauthorized
	}

	if request.Status != models.FriendRequestPending {
		return models.ErrFriendRequestNotFound
	}

	// Update request status
	if err := s.friendRepo.UpdateFriendRequestStatus(ctx, requestID, models.FriendRequestAccepted); err != nil {
		return err
	}

	// Create friendship
	return s.friendRepo.CreateFriendship(ctx, request.FromUserID, request.ToUserID)
}

func (s *FriendService) RejectFriendRequest(ctx context.Context, userID uuid.UUID, requestID uuid.UUID) error {
	request, err := s.friendRepo.GetFriendRequest(ctx, requestID)
	if err != nil {
		return err
	}

	// Verify the current user is the recipient
	if request.ToUserID != userID {
		return models.ErrUnauthorized
	}

	if request.Status != models.FriendRequestPending {
		return models.ErrFriendRequestNotFound
	}

	return s.friendRepo.UpdateFriendRequestStatus(ctx, requestID, models.FriendRequestRejected)
}

func (s *FriendService) GetFriends(ctx context.Context, userID uuid.UUID) ([]models.Friend, error) {
	return s.friendRepo.GetFriends(ctx, userID)
}
