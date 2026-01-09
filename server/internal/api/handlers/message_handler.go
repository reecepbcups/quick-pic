package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/quickpic/server/internal/models"
	"github.com/quickpic/server/internal/repository"
	"github.com/quickpic/server/internal/services"
)

type MessageHandler struct {
	messageService *services.MessageService
	userRepo       *repository.UserRepository
}

func NewMessageHandler(messageService *services.MessageService, userRepo *repository.UserRepository) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
		userRepo:       userRepo,
	}
}

func (h *MessageHandler) SendMessage(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var req models.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get recipient user ID
	toUser, err := h.userRepo.GetByUsername(c.Request.Context(), req.ToUsername)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "recipient not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find recipient"})
		return
	}

	msg, err := h.messageService.SendMessage(c.Request.Context(), userID, toUser.ID, &req)
	if err != nil {
		if errors.Is(err, models.ErrNotFriends) {
			c.JSON(http.StatusForbidden, gin.H{"error": "you can only send messages to friends"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send message"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         msg.ID,
		"created_at": msg.CreatedAt,
	})
}

func (h *MessageHandler) GetMessages(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	messages, err := h.messageService.GetPendingMessages(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get messages"})
		return
	}

	if messages == nil {
		messages = []models.MessageWithSender{}
	}

	c.JSON(http.StatusOK, messages)
}

func (h *MessageHandler) AcknowledgeMessage(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	messageIDStr := c.Param("id")
	messageID, err := uuid.Parse(messageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid message ID"})
		return
	}

	err = h.messageService.AcknowledgeMessage(c.Request.Context(), userID, messageID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrMessageNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "message not found"})
		case errors.Is(err, models.ErrUnauthorized):
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized to acknowledge this message"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to acknowledge message"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "message acknowledged and deleted"})
}
