package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/quickpic/server/internal/models"
	"github.com/quickpic/server/internal/repository"
)

type MessageService struct {
	messageRepo *repository.MessageRepository
	friendRepo  *repository.FriendRepository
}

func NewMessageService(messageRepo *repository.MessageRepository, friendRepo *repository.FriendRepository) *MessageService {
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

func (s *MessageService) AcknowledgeMessage(ctx context.Context, userID uuid.UUID, messageID uuid.UUID) error {
	msg, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return err
	}

	// Verify the current user is the recipient
	if msg.ToUserID != userID {
		return models.ErrUnauthorized
	}

	// Delete the message
	return s.messageRepo.Delete(ctx, messageID)
}
