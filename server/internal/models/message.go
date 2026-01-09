package models

import (
	"time"

	"github.com/google/uuid"
)

type ContentType string

const (
	ContentTypeText  ContentType = "text"
	ContentTypeImage ContentType = "image"
)

type Message struct {
	ID               uuid.UUID   `json:"id"`
	FromUserID       uuid.UUID   `json:"from_user_id"`
	ToUserID         uuid.UUID   `json:"to_user_id"`
	EncryptedContent []byte      `json:"encrypted_content"`
	ContentType      ContentType `json:"content_type"`
	Signature        string      `json:"signature"`
	CreatedAt        time.Time   `json:"created_at"`
}

type MessageWithSender struct {
	Message
	FromUsername  string `json:"from_username"`
	FromPublicKey string `json:"from_public_key"`
}

type SendMessageRequest struct {
	ToUsername       string      `json:"to_username" binding:"required"`
	EncryptedContent []byte      `json:"encrypted_content" binding:"required"`
	ContentType      ContentType `json:"content_type" binding:"required"`
	Signature        string      `json:"signature" binding:"required"`
}

type AcknowledgeMessageRequest struct {
	MessageID uuid.UUID `json:"message_id" binding:"required"`
}
