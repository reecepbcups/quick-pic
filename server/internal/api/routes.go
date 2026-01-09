package api

import (
	"github.com/gin-gonic/gin"
	"github.com/quickpic/server/internal/api/handlers"
	"github.com/quickpic/server/internal/api/middleware"
	"github.com/quickpic/server/internal/services"
	"github.com/quickpic/server/internal/storage"
)

func SetupRoutes(
	router *gin.Engine,
	authService *services.AuthService,
	userService *services.UserService,
	friendService *services.FriendService,
	messageService *services.MessageService,
	userRepo storage.UserRepo,
) {
	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	friendHandler := handlers.NewFriendHandler(friendService)
	messageHandler := handlers.NewMessageHandler(messageService, userRepo)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Auth routes (public)
	auth := router.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)
		auth.POST("/logout", authHandler.Logout)
	}

	// Protected routes
	protected := router.Group("")
	protected.Use(middleware.AuthMiddleware(authService))
	{
		// User routes
		protected.GET("/users/:username", userHandler.GetByUsername)

		// Friend routes
		friends := protected.Group("/friends")
		{
			friends.POST("/request", friendHandler.SendRequest)
			friends.GET("/requests", friendHandler.GetPendingRequests)
			friends.POST("/accept", friendHandler.AcceptRequest)
			friends.POST("/reject", friendHandler.RejectRequest)
			friends.GET("", friendHandler.GetFriends)
		}

		// Message routes
		messages := protected.Group("/messages")
		{
			messages.POST("", messageHandler.SendMessage)
			messages.GET("", messageHandler.GetMessages)
		}
	}
}
