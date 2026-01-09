package e2e

import (
	"net/http"
	"testing"
)

// =============================================================================
// AUTH TESTS
// =============================================================================

func TestHealthCheck(t *testing.T) {
	client := NewTestClient(t)

	resp := client.Get("/health")
	defer func() { _ = resp.Body.Close() }()

	client.ExpectStatus(resp, http.StatusOK)

	var result map[string]string
	client.ParseJSON(resp, &result)

	if result["status"] != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", result["status"])
	}
}

func TestRegister_Success(t *testing.T) {
	client := NewTestClient(t)

	req := RegisterRequest{
		Username:  uniqueUsername(),
		Password:  "password123",
		PublicKey: fakePublicKey(),
	}

	resp := client.Post("/auth/register", req)
	client.ExpectStatus(resp, http.StatusCreated)

	var authResp AuthResponse
	client.ParseJSON(resp, &authResp)

	if authResp.AccessToken == "" {
		t.Error("Expected access token to be set")
	}
	if authResp.RefreshToken == "" {
		t.Error("Expected refresh token to be set")
	}
	if authResp.User.Username != req.Username {
		t.Errorf("Expected username '%s', got '%s'", req.Username, authResp.User.Username)
	}
	if authResp.User.PublicKey != req.PublicKey {
		t.Errorf("Expected public key '%s', got '%s'", req.PublicKey, authResp.User.PublicKey)
	}
}

func TestRegister_DuplicateUsername(t *testing.T) {
	client := NewTestClient(t)
	username := uniqueUsername()

	// First registration should succeed
	req := RegisterRequest{
		Username:  username,
		Password:  "password123",
		PublicKey: fakePublicKey(),
	}

	resp := client.Post("/auth/register", req)
	client.ExpectStatus(resp, http.StatusCreated)
	_ = resp.Body.Close()

	// Second registration with same username should fail
	resp = client.Post("/auth/register", req)
	client.ExpectStatus(resp, http.StatusConflict)
	_ = resp.Body.Close()
}

func TestRegister_InvalidInput(t *testing.T) {
	client := NewTestClient(t)

	tests := []struct {
		name string
		req  RegisterRequest
	}{
		{
			name: "short username",
			req: RegisterRequest{
				Username:  "ab",
				Password:  "password123",
				PublicKey: fakePublicKey(),
			},
		},
		{
			name: "short password",
			req: RegisterRequest{
				Username:  uniqueUsername(),
				Password:  "short",
				PublicKey: fakePublicKey(),
			},
		},
		{
			name: "missing public key",
			req: RegisterRequest{
				Username: uniqueUsername(),
				Password: "password123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := client.Post("/auth/register", tt.req)
			client.ExpectStatus(resp, http.StatusBadRequest)
			_ = resp.Body.Close()
		})
	}
}

func TestLogin_Success(t *testing.T) {
	client := NewTestClient(t)
	username := uniqueUsername()
	password := "password123"

	// Register first
	regReq := RegisterRequest{
		Username:  username,
		Password:  password,
		PublicKey: fakePublicKey(),
	}
	resp := client.Post("/auth/register", regReq)
	_ = resp.Body.Close()

	// Login
	loginReq := LoginRequest{
		Username: username,
		Password: password,
	}

	resp = client.Post("/auth/login", loginReq)
	client.ExpectStatus(resp, http.StatusOK)

	var authResp AuthResponse
	client.ParseJSON(resp, &authResp)

	if authResp.AccessToken == "" {
		t.Error("Expected access token to be set")
	}
	if authResp.User.Username != username {
		t.Errorf("Expected username '%s', got '%s'", username, authResp.User.Username)
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	client := NewTestClient(t)
	username := uniqueUsername()

	// Register first
	regReq := RegisterRequest{
		Username:  username,
		Password:  "password123",
		PublicKey: fakePublicKey(),
	}
	resp := client.Post("/auth/register", regReq)
	_ = resp.Body.Close()

	// Login with wrong password
	loginReq := LoginRequest{
		Username: username,
		Password: "wrongpassword",
	}

	resp = client.Post("/auth/login", loginReq)
	client.ExpectStatus(resp, http.StatusUnauthorized)
	_ = resp.Body.Close()
}

func TestLogin_NonexistentUser(t *testing.T) {
	client := NewTestClient(t)

	loginReq := LoginRequest{
		Username: "nonexistentuser",
		Password: "password123",
	}

	resp := client.Post("/auth/login", loginReq)
	client.ExpectStatus(resp, http.StatusUnauthorized)
	_ = resp.Body.Close()
}

func TestRefreshToken_Success(t *testing.T) {
	client := NewTestClient(t)

	// Register and get tokens
	regReq := RegisterRequest{
		Username:  uniqueUsername(),
		Password:  "password123",
		PublicKey: fakePublicKey(),
	}

	resp := client.Post("/auth/register", regReq)
	var authResp AuthResponse
	client.ParseJSON(resp, &authResp)

	// Refresh token
	refreshReq := map[string]string{
		"refresh_token": authResp.RefreshToken,
	}

	resp = client.Post("/auth/refresh", refreshReq)
	client.ExpectStatus(resp, http.StatusOK)

	var newAuthResp AuthResponse
	client.ParseJSON(resp, &newAuthResp)

	if newAuthResp.AccessToken == "" {
		t.Error("Expected new access token")
	}
	if newAuthResp.RefreshToken == "" {
		t.Error("Expected new refresh token")
	}
	if newAuthResp.User.Username != regReq.Username {
		t.Errorf("Expected user '%s', got '%s'", regReq.Username, newAuthResp.User.Username)
	}
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	client := NewTestClient(t)

	refreshReq := map[string]string{
		"refresh_token": "invalid-token",
	}

	resp := client.Post("/auth/refresh", refreshReq)
	client.ExpectStatus(resp, http.StatusUnauthorized)
	_ = resp.Body.Close()
}

func TestLogout_Success(t *testing.T) {
	client := NewTestClient(t)

	// Register and get tokens
	regReq := RegisterRequest{
		Username:  uniqueUsername(),
		Password:  "password123",
		PublicKey: fakePublicKey(),
	}

	resp := client.Post("/auth/register", regReq)
	var authResp AuthResponse
	client.ParseJSON(resp, &authResp)

	// Logout
	logoutReq := map[string]string{
		"refresh_token": authResp.RefreshToken,
	}

	resp = client.Post("/auth/logout", logoutReq)
	client.ExpectStatus(resp, http.StatusOK)
	_ = resp.Body.Close()

	// Try to use the refresh token after logout - should fail
	refreshReq := map[string]string{
		"refresh_token": authResp.RefreshToken,
	}

	resp = client.Post("/auth/refresh", refreshReq)
	client.ExpectStatus(resp, http.StatusUnauthorized)
	_ = resp.Body.Close()
}

// =============================================================================
// FRIEND TESTS
// =============================================================================

func TestFriends_SendRequest_Success(t *testing.T) {
	client := NewTestClient(t)

	// Create two users
	user1 := createAuthenticatedUser(t, client)
	user2 := createAuthenticatedUser(t, client)

	// User1 sends friend request to User2
	client.SetAccessToken(user1.AccessToken)

	friendReq := SendFriendRequest{
		Username: user2.User.Username,
	}

	resp := client.Post("/friends/request", friendReq)
	client.ExpectStatus(resp, http.StatusCreated)

	var friendRequest FriendRequest
	client.ParseJSON(resp, &friendRequest)

	if friendRequest.Status != "pending" {
		t.Errorf("Expected status 'pending', got '%s'", friendRequest.Status)
	}
}

func TestFriends_SendRequest_ToSelf(t *testing.T) {
	client := NewTestClient(t)

	user := createAuthenticatedUser(t, client)
	client.SetAccessToken(user.AccessToken)

	friendReq := SendFriendRequest{
		Username: user.User.Username,
	}

	resp := client.Post("/friends/request", friendReq)
	client.ExpectStatus(resp, http.StatusBadRequest)
	_ = resp.Body.Close()
}

func TestFriends_SendRequest_UserNotFound(t *testing.T) {
	client := NewTestClient(t)

	user := createAuthenticatedUser(t, client)
	client.SetAccessToken(user.AccessToken)

	friendReq := SendFriendRequest{
		Username: "nonexistentuser",
	}

	resp := client.Post("/friends/request", friendReq)
	client.ExpectStatus(resp, http.StatusNotFound)
	_ = resp.Body.Close()
}

func TestFriends_SendRequest_Duplicate(t *testing.T) {
	client := NewTestClient(t)

	user1 := createAuthenticatedUser(t, client)
	user2 := createAuthenticatedUser(t, client)

	client.SetAccessToken(user1.AccessToken)

	friendReq := SendFriendRequest{
		Username: user2.User.Username,
	}

	// First request
	resp := client.Post("/friends/request", friendReq)
	client.ExpectStatus(resp, http.StatusCreated)
	_ = resp.Body.Close()

	// Duplicate request
	resp = client.Post("/friends/request", friendReq)
	client.ExpectStatus(resp, http.StatusConflict)
	_ = resp.Body.Close()
}

func TestFriends_GetPendingRequests(t *testing.T) {
	client := NewTestClient(t)

	user1 := createAuthenticatedUser(t, client)
	user2 := createAuthenticatedUser(t, client)

	// User1 sends request to User2
	client.SetAccessToken(user1.AccessToken)
	resp := client.Post("/friends/request", SendFriendRequest{Username: user2.User.Username})
	_ = resp.Body.Close()

	// User2 checks pending requests
	client.SetAccessToken(user2.AccessToken)
	resp = client.Get("/friends/requests")
	client.ExpectStatus(resp, http.StatusOK)

	var requests []FriendRequest
	client.ParseJSON(resp, &requests)

	if len(requests) != 1 {
		t.Errorf("Expected 1 pending request, got %d", len(requests))
	}
	if len(requests) > 0 && requests[0].FromUser.Username != user1.User.Username {
		t.Errorf("Expected request from '%s', got '%s'", user1.User.Username, requests[0].FromUser.Username)
	}
}

func TestFriends_AcceptRequest(t *testing.T) {
	client := NewTestClient(t)

	user1 := createAuthenticatedUser(t, client)
	user2 := createAuthenticatedUser(t, client)

	// User1 sends request to User2
	client.SetAccessToken(user1.AccessToken)
	resp := client.Post("/friends/request", SendFriendRequest{Username: user2.User.Username})
	var friendRequest FriendRequest
	client.ParseJSON(resp, &friendRequest)

	// User2 accepts the request
	client.SetAccessToken(user2.AccessToken)
	acceptReq := FriendRequestAction{
		RequestID: friendRequest.ID,
	}

	resp = client.Post("/friends/accept", acceptReq)
	client.ExpectStatus(resp, http.StatusOK)
	_ = resp.Body.Close()

	// Verify they are now friends
	resp = client.Get("/friends")
	client.ExpectStatus(resp, http.StatusOK)

	var friends []Friend
	client.ParseJSON(resp, &friends)

	if len(friends) != 1 {
		t.Errorf("Expected 1 friend, got %d", len(friends))
	}
	if len(friends) > 0 && friends[0].Username != user1.User.Username {
		t.Errorf("Expected friend '%s', got '%s'", user1.User.Username, friends[0].Username)
	}
}

func TestFriends_RejectRequest(t *testing.T) {
	client := NewTestClient(t)

	user1 := createAuthenticatedUser(t, client)
	user2 := createAuthenticatedUser(t, client)

	// User1 sends request to User2
	client.SetAccessToken(user1.AccessToken)
	resp := client.Post("/friends/request", SendFriendRequest{Username: user2.User.Username})
	var friendRequest FriendRequest
	client.ParseJSON(resp, &friendRequest)

	// User2 rejects the request
	client.SetAccessToken(user2.AccessToken)
	rejectReq := FriendRequestAction{
		RequestID: friendRequest.ID,
	}

	resp = client.Post("/friends/reject", rejectReq)
	client.ExpectStatus(resp, http.StatusOK)
	_ = resp.Body.Close()

	// Verify they are not friends
	resp = client.Get("/friends")
	client.ExpectStatus(resp, http.StatusOK)

	var friends []Friend
	client.ParseJSON(resp, &friends)

	if len(friends) != 0 {
		t.Errorf("Expected 0 friends, got %d", len(friends))
	}
}

func TestFriends_GetFriends(t *testing.T) {
	client := NewTestClient(t)

	user1 := createAuthenticatedUser(t, client)
	user2 := createAuthenticatedUser(t, client)

	// Make them friends
	makeFriends(t, client, user1, user2)

	// User1 checks friends list
	client.SetAccessToken(user1.AccessToken)
	resp := client.Get("/friends")
	client.ExpectStatus(resp, http.StatusOK)

	var friends []Friend
	client.ParseJSON(resp, &friends)

	if len(friends) != 1 {
		t.Errorf("Expected 1 friend, got %d", len(friends))
	}
	if len(friends) > 0 {
		if friends[0].Username != user2.User.Username {
			t.Errorf("Expected friend '%s', got '%s'", user2.User.Username, friends[0].Username)
		}
		if friends[0].PublicKey == "" {
			t.Error("Expected friend's public key to be included")
		}
	}
}

// =============================================================================
// MESSAGE TESTS
// =============================================================================

func TestMessages_Send_Success(t *testing.T) {
	client := NewTestClient(t)

	user1 := createAuthenticatedUser(t, client)
	user2 := createAuthenticatedUser(t, client)

	// Make them friends first
	makeFriends(t, client, user1, user2)

	// User1 sends message to User2
	client.SetAccessToken(user1.AccessToken)

	msgReq := SendMessageRequest{
		ToUsername:       user2.User.Username,
		EncryptedContent: "ZW5jcnlwdGVkLWNvbnRlbnQtaGVyZQ==", // base64 encoded
		ContentType:      "text",
		Signature:        "c2lnbmF0dXJlLWhlcmU=",
	}

	resp := client.Post("/messages", msgReq)
	client.ExpectStatus(resp, http.StatusCreated)

	var msgResp MessageResponse
	client.ParseJSON(resp, &msgResp)

	if msgResp.ID == "" {
		t.Error("Expected message ID to be set")
	}
}

func TestMessages_Send_NotFriends(t *testing.T) {
	client := NewTestClient(t)

	user1 := createAuthenticatedUser(t, client)
	user2 := createAuthenticatedUser(t, client)

	// Don't make them friends - try to send message
	client.SetAccessToken(user1.AccessToken)

	msgReq := SendMessageRequest{
		ToUsername:       user2.User.Username,
		EncryptedContent: "ZW5jcnlwdGVkLWNvbnRlbnQ=",
		ContentType:      "text",
		Signature:        "c2lnbmF0dXJl",
	}

	resp := client.Post("/messages", msgReq)
	client.ExpectStatus(resp, http.StatusForbidden)
	_ = resp.Body.Close()
}

func TestMessages_Send_UserNotFound(t *testing.T) {
	client := NewTestClient(t)

	user := createAuthenticatedUser(t, client)
	client.SetAccessToken(user.AccessToken)

	msgReq := SendMessageRequest{
		ToUsername:       "nonexistentuser",
		EncryptedContent: "ZW5jcnlwdGVkLWNvbnRlbnQ=",
		ContentType:      "text",
		Signature:        "c2lnbmF0dXJl",
	}

	resp := client.Post("/messages", msgReq)
	client.ExpectStatus(resp, http.StatusNotFound)
	_ = resp.Body.Close()
}

func TestMessages_GetPending(t *testing.T) {
	client := NewTestClient(t)

	user1 := createAuthenticatedUser(t, client)
	user2 := createAuthenticatedUser(t, client)

	makeFriends(t, client, user1, user2)

	// User1 sends message to User2
	client.SetAccessToken(user1.AccessToken)
	msgReq := SendMessageRequest{
		ToUsername:       user2.User.Username,
		EncryptedContent: "ZW5jcnlwdGVkLW1lc3NhZ2U=",
		ContentType:      "text",
		Signature:        "c2lnbmF0dXJl",
	}
	resp := client.Post("/messages", msgReq)
	_ = resp.Body.Close()

	// User2 fetches pending messages
	client.SetAccessToken(user2.AccessToken)
	resp = client.Get("/messages")
	client.ExpectStatus(resp, http.StatusOK)

	var messages []Message
	client.ParseJSON(resp, &messages)

	if len(messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(messages))
	}
	if len(messages) > 0 {
		if messages[0].FromUsername != user1.User.Username {
			t.Errorf("Expected message from '%s', got '%s'", user1.User.Username, messages[0].FromUsername)
		}
		if messages[0].ContentType != "text" {
			t.Errorf("Expected content type 'text', got '%s'", messages[0].ContentType)
		}
	}
}


// =============================================================================
// PROTECTED ROUTES TESTS
// =============================================================================

func TestProtectedRoutes_RequireAuth(t *testing.T) {
	client := NewTestClient(t)
	client.ClearAccessToken()

	routes := []struct {
		method string
		path   string
	}{
		{"GET", "/friends"},
		{"GET", "/friends/requests"},
		{"POST", "/friends/request"},
		{"POST", "/friends/accept"},
		{"POST", "/friends/reject"},
		{"GET", "/messages"},
		{"POST", "/messages"},
		{"GET", "/users/someuser"},
	}

	for _, route := range routes {
		t.Run(route.method+" "+route.path, func(t *testing.T) {
			var resp *http.Response
			if route.method == "GET" {
				resp = client.Get(route.path)
			} else {
				resp = client.Post(route.path, map[string]string{})
			}
			client.ExpectStatus(resp, http.StatusUnauthorized)
			_ = resp.Body.Close()
		})
	}
}

func TestProtectedRoutes_InvalidToken(t *testing.T) {
	client := NewTestClient(t)
	client.SetAccessToken("invalid-token")

	resp := client.Get("/friends")
	client.ExpectStatus(resp, http.StatusUnauthorized)
	_ = resp.Body.Close()
}

// =============================================================================
// FULL USER JOURNEY TEST
// =============================================================================

func TestFullUserJourney(t *testing.T) {
	client := NewTestClient(t)

	// 1. Alice registers
	aliceUsername := uniqueUsername()
	aliceReq := RegisterRequest{
		Username:  aliceUsername,
		Password:  "alicepassword123",
		PublicKey: fakePublicKey(),
	}
	resp := client.Post("/auth/register", aliceReq)
	client.ExpectStatus(resp, http.StatusCreated)
	var aliceAuth AuthResponse
	client.ParseJSON(resp, &aliceAuth)
	t.Logf("Alice registered: %s", aliceAuth.User.Username)

	// 2. Bob registers
	bobUsername := uniqueUsername()
	bobReq := RegisterRequest{
		Username:  bobUsername,
		Password:  "bobpassword123",
		PublicKey: fakePublicKey(),
	}
	resp = client.Post("/auth/register", bobReq)
	client.ExpectStatus(resp, http.StatusCreated)
	var bobAuth AuthResponse
	client.ParseJSON(resp, &bobAuth)
	t.Logf("Bob registered: %s", bobAuth.User.Username)

	// 3. Alice sends friend request to Bob
	client.SetAccessToken(aliceAuth.AccessToken)
	resp = client.Post("/friends/request", SendFriendRequest{Username: bobUsername})
	client.ExpectStatus(resp, http.StatusCreated)
	var friendReq FriendRequest
	client.ParseJSON(resp, &friendReq)
	t.Logf("Alice sent friend request to Bob")

	// 4. Bob checks pending requests
	client.SetAccessToken(bobAuth.AccessToken)
	resp = client.Get("/friends/requests")
	client.ExpectStatus(resp, http.StatusOK)
	var requests []FriendRequest
	client.ParseJSON(resp, &requests)
	if len(requests) != 1 {
		t.Fatalf("Bob should have 1 pending request, got %d", len(requests))
	}
	t.Logf("Bob has 1 pending request from Alice")

	// 5. Bob accepts the request
	resp = client.Post("/friends/accept", FriendRequestAction{RequestID: friendReq.ID})
	client.ExpectStatus(resp, http.StatusOK)
	_ = resp.Body.Close()
	t.Logf("Bob accepted Alice's friend request")

	// 6. Alice sends a message to Bob
	client.SetAccessToken(aliceAuth.AccessToken)
	msgReq := SendMessageRequest{
		ToUsername:       bobUsername,
		EncryptedContent: "SGVsbG8gQm9iIQ==", // "Hello Bob!" base64
		ContentType:      "text",
		Signature:        "YWxpY2Utc2lnbmF0dXJl",
	}
	resp = client.Post("/messages", msgReq)
	client.ExpectStatus(resp, http.StatusCreated)
	var msgResp MessageResponse
	client.ParseJSON(resp, &msgResp)
	t.Logf("Alice sent message to Bob: %s", msgResp.ID)

	// 7. Bob fetches his messages
	client.SetAccessToken(bobAuth.AccessToken)
	resp = client.Get("/messages")
	client.ExpectStatus(resp, http.StatusOK)
	var messages []Message
	client.ParseJSON(resp, &messages)
	if len(messages) != 1 {
		t.Fatalf("Bob should have 1 message, got %d", len(messages))
	}
	if messages[0].FromUsername != aliceUsername {
		t.Errorf("Message should be from Alice, got %s", messages[0].FromUsername)
	}
	t.Logf("Bob received message from Alice")

	// 8. Alice logs out
	client.SetAccessToken(aliceAuth.AccessToken)
	resp = client.Post("/auth/logout", map[string]string{"refresh_token": aliceAuth.RefreshToken})
	client.ExpectStatus(resp, http.StatusOK)
	_ = resp.Body.Close()
	t.Logf("Alice logged out")

	// 9. Alice's refresh token should no longer work
	resp = client.Post("/auth/refresh", map[string]string{"refresh_token": aliceAuth.RefreshToken})
	client.ExpectStatus(resp, http.StatusUnauthorized)
	_ = resp.Body.Close()
	t.Logf("Alice's refresh token invalidated after logout")

	t.Log("Full user journey completed successfully!")
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

func createAuthenticatedUser(t *testing.T, client *TestClient) AuthResponse {
	req := RegisterRequest{
		Username:  uniqueUsername(),
		Password:  "password123",
		PublicKey: fakePublicKey(),
	}

	resp := client.Post("/auth/register", req)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to create user: status %d", resp.StatusCode)
	}

	var authResp AuthResponse
	client.ParseJSON(resp, &authResp)
	return authResp
}

func makeFriends(t *testing.T, client *TestClient, user1, user2 AuthResponse) {
	// User1 sends request
	client.SetAccessToken(user1.AccessToken)
	resp := client.Post("/friends/request", SendFriendRequest{Username: user2.User.Username})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to send friend request: status %d", resp.StatusCode)
	}
	var friendReq FriendRequest
	client.ParseJSON(resp, &friendReq)

	// User2 accepts
	client.SetAccessToken(user2.AccessToken)
	resp = client.Post("/friends/accept", FriendRequestAction{RequestID: friendReq.ID})
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Failed to accept friend request: status %d", resp.StatusCode)
	}
	_ = resp.Body.Close()
}
