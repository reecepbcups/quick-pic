package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/quickpic/server/internal/api"
	"github.com/quickpic/server/internal/repository"
	"github.com/quickpic/server/internal/services"
)

var (
	testServer *httptest.Server
	testRouter *gin.Engine
	testDB     *repository.DB
)

// TestMain sets up and tears down the test environment
func TestMain(m *testing.M) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	jwtSecret := "test-secret-key-for-testing"

	// Initialize in-memory SQLite database
	var err error
	testDB, err = repository.NewDB(":memory:")
	if err != nil {
		fmt.Printf("Failed to create test database: %v\n", err)
		os.Exit(1)
	}

	// Run migrations
	if err := testDB.Migrate(); err != nil {
		fmt.Printf("Failed to run migrations: %v\n", err)
		os.Exit(1)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(testDB)
	friendRepo := repository.NewFriendRepository(testDB)
	messageRepo := repository.NewMessageRepository(testDB)

	// Initialize services
	authService := services.NewAuthService(userRepo, jwtSecret)
	userService := services.NewUserService(userRepo)
	friendService := services.NewFriendService(friendRepo, userRepo)
	messageService := services.NewMessageService(messageRepo, friendRepo)

	// Initialize router
	testRouter = gin.New()
	testRouter.Use(gin.Recovery())

	// Setup routes
	api.SetupRoutesWithRepo(testRouter, authService, userService, friendService, messageService, userRepo)

	// Create test server
	testServer = httptest.NewServer(testRouter)

	// Run tests
	code := m.Run()

	// Cleanup
	testServer.Close()
	testDB.Close()

	os.Exit(code)
}

// TestClient provides helper methods for making HTTP requests
type TestClient struct {
	t           *testing.T
	baseURL     string
	accessToken string
}

func NewTestClient(t *testing.T) *TestClient {
	// Reset database before each test
	testDB.Reset()

	return &TestClient{
		t:       t,
		baseURL: testServer.URL,
	}
}

func (c *TestClient) SetAccessToken(token string) {
	c.accessToken = token
}

func (c *TestClient) ClearAccessToken() {
	c.accessToken = ""
}

// HTTP helper methods

func (c *TestClient) Post(path string, body interface{}) *http.Response {
	return c.request("POST", path, body)
}

func (c *TestClient) Get(path string) *http.Response {
	return c.request("GET", path, nil)
}

func (c *TestClient) request(method, path string, body interface{}) *http.Response {
	var req *http.Request
	var err error

	url := c.baseURL + path

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			c.t.Fatalf("Failed to marshal request body: %v", err)
		}
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
		if err != nil {
			c.t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			c.t.Fatalf("Failed to create request: %v", err)
		}
	}

	if c.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.accessToken)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.t.Fatalf("Failed to make request: %v", err)
	}

	return resp
}

// Response parsing helpers

func (c *TestClient) ParseJSON(resp *http.Response, v interface{}) {
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		c.t.Fatalf("Failed to parse JSON response: %v", err)
	}
}

func (c *TestClient) ExpectStatus(resp *http.Response, expected int) {
	if resp.StatusCode != expected {
		c.t.Errorf("Expected status %d, got %d", expected, resp.StatusCode)
	}
}

// Test data generators

type RegisterRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	PublicKey string `json:"public_key"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	User         struct {
		ID        string `json:"id"`
		Username  string `json:"username"`
		PublicKey string `json:"public_key"`
	} `json:"user"`
}

type SendFriendRequest struct {
	Username string `json:"username"`
}

type FriendRequestAction struct {
	RequestID string `json:"request_id"`
}

type FriendRequest struct {
	ID         string `json:"id"`
	FromUserID string `json:"from_user_id"`
	ToUserID   string `json:"to_user_id"`
	Status     string `json:"status"`
	FromUser   struct {
		ID        string `json:"id"`
		Username  string `json:"username"`
		PublicKey string `json:"public_key"`
	} `json:"from_user"`
}

type Friend struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	PublicKey string `json:"public_key"`
	Since     string `json:"since"`
}

type SendMessageRequest struct {
	ToUsername       string `json:"to_username"`
	EncryptedContent string `json:"encrypted_content"`
	ContentType      string `json:"content_type"`
	Signature        string `json:"signature"`
}

type MessageResponse struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
}

type Message struct {
	ID               string `json:"id"`
	FromUserID       string `json:"from_user_id"`
	ToUserID         string `json:"to_user_id"`
	EncryptedContent string `json:"encrypted_content"`
	ContentType      string `json:"content_type"`
	Signature        string `json:"signature"`
	CreatedAt        string `json:"created_at"`
	FromUsername     string `json:"from_username"`
	FromPublicKey    string `json:"from_public_key"`
}

// Helper to generate fake public key (base64 encoded 32 bytes)
func fakePublicKey() string {
	return "dGVzdC1wdWJsaWMta2V5LWZvci10ZXN0aW5nLXB1cnBvc2Vz"
}

// Helper to generate unique username
var userCounter = 0

func uniqueUsername() string {
	userCounter++
	return fmt.Sprintf("testuser%d", userCounter)
}
