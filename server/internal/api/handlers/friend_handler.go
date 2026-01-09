package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/quickpic/server/internal/models"
	"github.com/quickpic/server/internal/services"
)

type FriendHandler struct {
	friendService *services.FriendService
}

func NewFriendHandler(friendService *services.FriendService) *FriendHandler {
	return &FriendHandler{friendService: friendService}
}

func (h *FriendHandler) SendRequest(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var req models.SendFriendRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	request, err := h.friendService.SendFriendRequest(c.Request.Context(), userID, req.Username)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		case errors.Is(err, models.ErrCannotAddSelf):
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot add yourself as a friend"})
		case errors.Is(err, models.ErrFriendRequestExists):
			c.JSON(http.StatusConflict, gin.H{"error": "friend request already exists"})
		case errors.Is(err, models.ErrAlreadyFriends):
			c.JSON(http.StatusConflict, gin.H{"error": "already friends"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send friend request"})
		}
		return
	}

	c.JSON(http.StatusCreated, request)
}

func (h *FriendHandler) GetPendingRequests(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	requests, err := h.friendService.GetPendingRequests(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get friend requests"})
		return
	}

	if requests == nil {
		requests = []models.FriendRequestWithUser{}
	}

	c.JSON(http.StatusOK, requests)
}

func (h *FriendHandler) AcceptRequest(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var req models.FriendRequestActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.friendService.AcceptFriendRequest(c.Request.Context(), userID, req.RequestID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrFriendRequestNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "friend request not found"})
		case errors.Is(err, models.ErrUnauthorized):
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized to accept this request"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to accept friend request"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "friend request accepted"})
}

func (h *FriendHandler) RejectRequest(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var req models.FriendRequestActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.friendService.RejectFriendRequest(c.Request.Context(), userID, req.RequestID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrFriendRequestNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "friend request not found"})
		case errors.Is(err, models.ErrUnauthorized):
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized to reject this request"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reject friend request"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "friend request rejected"})
}

func (h *FriendHandler) GetFriends(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	friends, err := h.friendService.GetFriends(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get friends"})
		return
	}

	if friends == nil {
		friends = []models.Friend{}
	}

	c.JSON(http.StatusOK, friends)
}
