package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/quickpic/server/internal/models"
	"github.com/quickpic/server/internal/storage"
)

type MessageService struct {
	messageRepo storage.MessageRepo
	friendRepo  storage.FriendRepo
}

func NewMessageService(messageRepo storage.MessageRepo, friendRepo storage.FriendRepo) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
		friendRepo:  friendRepo,
	}
}

func (s *MessageService) SendMessage(ctx context.Context, fromUserID uuid.UUID, toUserID uuid.UUID, req *models.SendMessageRequest) (*models.Message, error) {
	// Verify users are friends
	areFriends, err := s.friendRepo.AreFriends(ctx, fromUserID, toUserID)
	if err != nil {
		return nil, err
	}
	if !areFriends {
		return nil, models.ErrNotFriends
	}

	msg := &models.Message{
		FromUserID:       fromUserID,
		ToUserID:         toUserID,
		EncryptedContent: req.EncryptedContent,
		ContentType:      req.ContentType,
		Signature:        req.Signature,
	}

	if err := s.messageRepo.Create(ctx, msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func (s *MessageService) GetPendingMessages(ctx context.Context, userID uuid.UUID) ([]models.MessageWithSender, error) {
	return s.messageRepo.GetPendingMessages(ctx, userID)
}
