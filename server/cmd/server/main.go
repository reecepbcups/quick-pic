package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/quickpic/server/internal/api"
	"github.com/quickpic/server/internal/repository"
	"github.com/quickpic/server/internal/services"
)

func main() {
	// Load configuration
	dbPath := getEnv("DATABASE_PATH", "./quickpic.db")
	jwtSecret := getEnv("JWT_SECRET", "development-secret-change-in-production")
	port := getEnv("PORT", "8080")

	// Initialize database
	db, err := repository.NewDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	friendRepo := repository.NewFriendRepository(db)
	messageRepo := repository.NewMessageRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, jwtSecret)
	userService := services.NewUserService(userRepo)
	friendService := services.NewFriendService(friendRepo, userRepo)
	messageService := services.NewMessageService(messageRepo, friendRepo)

	// Initialize router
	router := gin.Default()

	// Setup routes
	api.SetupRoutesWithRepo(router, authService, userService, friendService, messageService, userRepo)

	// Start server
	log.Printf("Starting QuickPic server on port %s", port)
	log.Printf("Using SQLite database: %s", dbPath)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
