package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/quickpic/server/internal/api"
	"github.com/quickpic/server/internal/backend"
	"github.com/quickpic/server/internal/services"
)

func main() {
	// Load configuration
	backendType := getEnv("BACKEND_TYPE", "sqlite")
	jwtSecret := getEnv("JWT_SECRET", "development-secret-change-in-production")
	port := getEnv("PORT", "8080")

	// Build backend configuration based on type
	var cfg backend.Config
	switch strings.ToLower(backendType) {
	case "sqlite":
		dbPath := getEnv("DATABASE_PATH", "./quickpic.db")
		cfg = backend.DefaultSQLiteConfig(dbPath)
		log.Printf("Using SQLite backend: %s", dbPath)

	case "blockchain":
		rpcURL := getEnv("BLOCKCHAIN_RPC_URL", "http://localhost:8545")
		privateKey := getEnv("BLOCKCHAIN_PRIVATE_KEY", "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
		contractAddress := getEnv("BLOCKCHAIN_CONTRACT_ADDRESS", "")

		if contractAddress == "" {
			log.Fatal("BLOCKCHAIN_CONTRACT_ADDRESS is required for blockchain backend")
		}

		cfg = backend.Config{
			Type:                      backend.TypeBlockchain,
			BlockchainRPCURL:          rpcURL,
			BlockchainPrivateKey:      privateKey,
			BlockchainContractAddress: contractAddress,
		}
		log.Printf("Using blockchain backend: %s (contract: %s)", rpcURL, contractAddress)

	default:
		log.Fatalf("Unknown backend type: %s (use 'sqlite' or 'blockchain')", backendType)
	}

	// Initialize backend
	result, err := backend.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize backend: %v", err)
	}
	defer result.Backend.Close()

	// Initialize services
	authService := services.NewAuthService(result.Repos.Users, jwtSecret)
	userService := services.NewUserService(result.Repos.Users)
	friendService := services.NewFriendService(result.Repos.Friends, result.Repos.Users)
	messageService := services.NewMessageService(result.Repos.Messages, result.Repos.Friends)

	// Initialize router
	router := gin.Default()

	// Setup routes
	api.SetupRoutes(router, authService, userService, friendService, messageService, result.Repos.Users)

	// Start server
	log.Printf("Starting QuickPic server on port %s", port)
	log.Printf("Backend: %s", result.Backend.Name())
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
