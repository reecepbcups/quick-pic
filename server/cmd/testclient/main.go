package main

import (
	"bytes"
	"compress/flate"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/hkdf"
)

const (
	baseURL        = "http://localhost:8080"
	targetUsername = "reecepbcups" // The user to friend and message
)

type TestClient struct {
	httpClient  *http.Client
	accessToken string
	privateKey  [32]byte
	publicKey   [32]byte
	username    string
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         struct {
		ID        string `json:"id"`
		Username  string `json:"username"`
		PublicKey string `json:"public_key"`
	} `json:"user"`
}

type Friend struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	PublicKey string `json:"public_key"`
}

type FriendRequest struct {
	ID           string `json:"id"`
	FromUserID   string `json:"from_user_id"`
	FromUsername string `json:"from_username"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: testclient <command>")
		fmt.Println("Commands:")
		fmt.Println("  register <username>                        - Create account only (no friend request)")
		fmt.Println("  friend <username> <target>                 - Send friend request to target user")
		fmt.Println("  setup <username>                           - Create account and send friend request to", targetUsername)
		fmt.Println("  message --from <user> --to <user> <msg>    - Send a message (requires friendship)")
		fmt.Println("  status <username>                          - Check friend request status")
		fmt.Println("  receive <username>                         - Receive and decrypt messages")
		fmt.Println("  accept <username> <request_id>             - Accept a friend request")
		fmt.Println("  pending <username>                         - List pending friend requests")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "debug":
		// Test encryption/decryption locally
		if err := runDebug(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case "register":
		if len(os.Args) < 3 {
			fmt.Println("Usage: testclient register <username>")
			os.Exit(1)
		}
		username := os.Args[2]
		if err := runRegisterOnly(username); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case "friend":
		if len(os.Args) < 4 {
			fmt.Println("Usage: testclient friend <username> <target_username>")
			os.Exit(1)
		}
		username := os.Args[2]
		target := os.Args[3]
		if err := runFriendOnly(username, target); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case "setup":
		if len(os.Args) < 3 {
			fmt.Println("Usage: testclient setup <username>")
			os.Exit(1)
		}
		username := os.Args[2]
		if err := runSetup(username); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case "message":
		msgFlags := flag.NewFlagSet("message", flag.ExitOnError)
		fromUser := msgFlags.String("from", "", "sender username (required)")
		toUser := msgFlags.String("to", "", "recipient username (required)")
		if err := msgFlags.Parse(os.Args[2:]); err != nil {
			fmt.Printf("Failed to parse message flags: %v\n", err)
			os.Exit(1)
		}

		if *fromUser == "" || *toUser == "" {
			fmt.Println("Usage: testclient message --from <sender> --to <recipient> <message>")
			fmt.Println("  --from  sender username (required)")
			fmt.Println("  --to    recipient username (required)")
			os.Exit(1)
		}

		args := msgFlags.Args()
		if len(args) < 1 {
			fmt.Println("Usage: testclient message --from <sender> --to <recipient> <message>")
			os.Exit(1)
		}
		message := args[0]

		if err := runMessage(*fromUser, *toUser, message); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case "status":
		if len(os.Args) < 3 {
			fmt.Println("Usage: testclient status <username>")
			os.Exit(1)
		}
		username := os.Args[2]
		if err := runStatus(username); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case "receive":
		if len(os.Args) < 3 {
			fmt.Println("Usage: testclient receive <username>")
			os.Exit(1)
		}
		username := os.Args[2]
		if err := runReceive(username); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case "pending":
		if len(os.Args) < 3 {
			fmt.Println("Usage: testclient pending <username>")
			os.Exit(1)
		}
		username := os.Args[2]
		if err := runPending(username); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case "accept":
		if len(os.Args) < 4 {
			fmt.Println("Usage: testclient accept <username> <request_id>")
			os.Exit(1)
		}
		username := os.Args[2]
		requestID := os.Args[3]
		if err := runAccept(username, requestID); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func runRegisterOnly(username string) error {
	client := &TestClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		username:   username,
	}

	// Generate keypair
	fmt.Println("Generating keypair...")
	if err := client.generateKeyPair(); err != nil {
		return fmt.Errorf("failed to generate keypair: %w", err)
	}

	// Register
	fmt.Printf("Registering user '%s'...\n", username)
	if err := client.register(); err != nil {
		return fmt.Errorf("failed to register: %w", err)
	}
	fmt.Println("Registration successful!")

	// Save credentials
	if err := client.saveCredentials(); err != nil {
		return fmt.Errorf("failed to save credentials: %w", err)
	}
	fmt.Printf("Credentials saved to %s.json\n", username)
	fmt.Printf("\nTo send a friend request, run:\n")
	fmt.Printf("  go run ./cmd/testclient friend %s <target_username>\n", username)

	return nil
}

func runFriendOnly(username, target string) error {
	client := &TestClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		username:   username,
	}

	// Load credentials
	if err := client.loadCredentials(); err != nil {
		return fmt.Errorf("failed to load credentials (did you run register first?): %w", err)
	}

	// Send friend request
	fmt.Printf("Sending friend request to '%s'...\n", target)
	if err := client.sendFriendRequest(target); err != nil {
		return fmt.Errorf("failed to send friend request: %w", err)
	}
	fmt.Println("Friend request sent!")

	return nil
}

func runPending(username string) error {
	client := &TestClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		username:   username,
	}

	// Load credentials
	if err := client.loadCredentials(); err != nil {
		return fmt.Errorf("failed to load credentials (did you run register first?): %w", err)
	}

	// Get pending requests
	req, _ := http.NewRequest("GET", baseURL+"/friends/requests", nil)
	req.Header.Set("Authorization", "Bearer "+client.accessToken)

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("get pending requests failed (%d): %s", resp.StatusCode, string(body))
	}

	var requests []struct {
		ID       string `json:"id"`
		FromUser struct {
			ID       string `json:"id"`
			Username string `json:"username"`
		} `json:"from_user"`
		CreatedAt string `json:"created_at"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&requests); err != nil {
		return fmt.Errorf("decode failed: %w", err)
	}

	fmt.Printf("Pending friend requests (%d):\n", len(requests))
	for _, r := range requests {
		fmt.Printf("  ID: %s | From: %s | At: %s\n", r.ID, r.FromUser.Username, r.CreatedAt)
	}

	if len(requests) > 0 {
		fmt.Printf("\nTo accept a request, run:\n")
		fmt.Printf("  go run ./cmd/testclient accept %s <request_id>\n", username)
	}

	return nil
}

func runAccept(username, requestID string) error {
	client := &TestClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		username:   username,
	}

	// Load credentials
	if err := client.loadCredentials(); err != nil {
		return fmt.Errorf("failed to load credentials (did you run register first?): %w", err)
	}

	// Accept friend request
	req, _ := http.NewRequest("POST", baseURL+"/friends/requests/"+requestID+"/accept", nil)
	req.Header.Set("Authorization", "Bearer "+client.accessToken)

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("accept request failed (%d): %s", resp.StatusCode, string(body))
	}

	fmt.Println("Friend request accepted!")
	return nil
}

func runSetup(username string) error {
	client := &TestClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		username:   username,
	}

	// Generate keypair
	fmt.Println("Generating keypair...")
	if err := client.generateKeyPair(); err != nil {
		return fmt.Errorf("failed to generate keypair: %w", err)
	}

	// Register
	fmt.Printf("Registering user '%s'...\n", username)
	if err := client.register(); err != nil {
		return fmt.Errorf("failed to register: %w", err)
	}
	fmt.Println("Registration successful!")

	// Save credentials
	if err := client.saveCredentials(); err != nil {
		return fmt.Errorf("failed to save credentials: %w", err)
	}
	fmt.Printf("Credentials saved to %s.json\n", username)

	// Send friend request
	fmt.Printf("Sending friend request to '%s'...\n", targetUsername)
	if err := client.sendFriendRequest(targetUsername); err != nil {
		return fmt.Errorf("failed to send friend request: %w", err)
	}
	fmt.Println("Friend request sent!")
	fmt.Printf("\nNow accept the friend request in your iOS app, then run:\n")
	fmt.Printf("  go run ./cmd/testclient message %s \"Hello from test client!\"\n", username)

	return nil
}

func runMessage(fromUsername, toUsername, message string) error {
	client := &TestClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		username:   fromUsername,
	}

	// Load credentials
	if err := client.loadCredentials(); err != nil {
		return fmt.Errorf("failed to load credentials for '%s' (did you run register first?): %w", fromUsername, err)
	}

	// Check if we're friends
	fmt.Printf("Checking if %s is friends with %s...\n", fromUsername, toUsername)
	friends, err := client.getFriends()
	if err != nil {
		return fmt.Errorf("failed to get friends: %w", err)
	}

	var targetFriend *Friend
	for _, f := range friends {
		if f.Username == toUsername {
			targetFriend = &f
			break
		}
	}

	if targetFriend == nil {
		return fmt.Errorf("not friends with %s yet - send/accept a friend request first", toUsername)
	}

	fmt.Printf("Friends with %s! Sending encrypted message...\n", toUsername)

	// Encrypt and send message
	if err := client.sendMessage(targetFriend, message); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	fmt.Println("Message sent successfully!")
	return nil
}

type ReceivedMessage struct {
	ID               string `json:"id"`
	FromUserID       string `json:"from_user_id"`
	ToUserID         string `json:"to_user_id"`
	EncryptedContent []byte `json:"encrypted_content"` // base64 decoded by Go json
	ContentType      string `json:"content_type"`
	Signature        string `json:"signature"`
	CreatedAt        string `json:"created_at"`
	FromUsername     string `json:"from_username"`
	FromPublicKey    string `json:"from_public_key"`
}

func runReceive(username string) error {
	client := &TestClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		username:   username,
	}

	// Load credentials
	if err := client.loadCredentials(); err != nil {
		return fmt.Errorf("failed to load credentials (did you run setup?): %w", err)
	}

	// Fetch messages
	fmt.Println("Fetching messages...")
	req, _ := http.NewRequest("GET", baseURL+"/messages", nil)
	req.Header.Set("Authorization", "Bearer "+client.accessToken)

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("fetch messages failed (%d): %s", resp.StatusCode, string(body))
	}

	var messages []ReceivedMessage
	if err := json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		return fmt.Errorf("decode messages failed: %w", err)
	}

	fmt.Printf("Found %d message(s)\n\n", len(messages))

	for i, msg := range messages {
		fmt.Printf("=== Message %d ===\n", i+1)
		fmt.Printf("From: %s\n", msg.FromUsername)
		fmt.Printf("Content-Type: %s\n", msg.ContentType)
		fmt.Printf("Encrypted length: %d bytes\n", len(msg.EncryptedContent))
		fmt.Printf("Encrypted (hex, first 64 bytes): %x\n", msg.EncryptedContent[:min(64, len(msg.EncryptedContent))])

		// Decode sender's public key
		senderPubKeyBytes, err := base64.StdEncoding.DecodeString(msg.FromPublicKey)
		if err != nil {
			fmt.Printf("ERROR: Failed to decode sender public key: %v\n\n", err)
			continue
		}
		var senderPubKey [32]byte
		copy(senderPubKey[:], senderPubKeyBytes)

		// Try to decrypt
		fmt.Println("\nAttempting decryption...")
		decrypted, err := client.decrypt(msg.EncryptedContent, senderPubKey)
		if err != nil {
			fmt.Printf("DECRYPT FAILED: %v\n\n", err)
		} else {
			fmt.Printf("DECRYPT SUCCESS!\n")
			if msg.ContentType == "text" {
				fmt.Printf("Message: %s\n\n", string(decrypted))
			} else {
				fmt.Printf("Image data: %d bytes\n\n", len(decrypted))
			}
		}
	}

	return nil
}

func runStatus(username string) error {
	client := &TestClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		username:   username,
	}

	// Load credentials
	if err := client.loadCredentials(); err != nil {
		return fmt.Errorf("failed to load credentials (did you run setup?): %w", err)
	}

	// Check friends
	friends, err := client.getFriends()
	if err != nil {
		return fmt.Errorf("failed to get friends: %w", err)
	}

	fmt.Printf("Friends (%d):\n", len(friends))
	for _, f := range friends {
		fmt.Printf("  - %s\n", f.Username)
	}

	isFriend := false
	for _, f := range friends {
		if f.Username == targetUsername {
			isFriend = true
			break
		}
	}

	if isFriend {
		fmt.Printf("\nYou are friends with %s! You can send messages.\n", targetUsername)
	} else {
		fmt.Printf("\nNot yet friends with %s. Accept the request in the iOS app.\n", targetUsername)
	}

	return nil
}

func (c *TestClient) generateKeyPair() error {
	// Generate random private key
	if _, err := rand.Read(c.privateKey[:]); err != nil {
		return err
	}

	// Derive public key
	curve25519.ScalarBaseMult(&c.publicKey, &c.privateKey)
	return nil
}

func (c *TestClient) register() error {
	publicKeyB64 := base64.StdEncoding.EncodeToString(c.publicKey[:])

	reqBody := map[string]string{
		"username":   c.username,
		"password":   "testpassword123",
		"public_key": publicKeyB64,
	}
	body, _ := json.Marshal(reqBody)

	resp, err := c.httpClient.Post(baseURL+"/auth/register", "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("registration failed (%d): %s", resp.StatusCode, string(respBody))
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return err
	}

	c.accessToken = authResp.AccessToken
	return nil
}

func (c *TestClient) login() error {
	reqBody := map[string]string{
		"username": c.username,
		"password": "testpassword123",
	}
	body, _ := json.Marshal(reqBody)

	resp, err := c.httpClient.Post(baseURL+"/auth/login", "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("login failed (%d): %s", resp.StatusCode, string(respBody))
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return err
	}

	c.accessToken = authResp.AccessToken
	return nil
}

func (c *TestClient) sendFriendRequest(username string) error {
	reqBody := map[string]string{"username": username}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", baseURL+"/friends/request", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("friend request failed (%d): %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func (c *TestClient) getFriends() ([]Friend, error) {
	req, _ := http.NewRequest("GET", baseURL+"/friends", nil)
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get friends failed (%d): %s", resp.StatusCode, string(respBody))
	}

	var friends []Friend
	if err := json.NewDecoder(resp.Body).Decode(&friends); err != nil {
		return nil, err
	}

	return friends, nil
}

func (c *TestClient) sendMessage(recipient *Friend, content string) error {
	// Decode recipient's public key
	recipientPubKeyBytes, err := base64.StdEncoding.DecodeString(recipient.PublicKey)
	if err != nil {
		return fmt.Errorf("invalid recipient public key: %w", err)
	}

	var recipientPubKey [32]byte
	copy(recipientPubKey[:], recipientPubKeyBytes)

	// Encrypt the message
	encryptedData, signature, err := c.encrypt([]byte(content), recipientPubKey)
	if err != nil {
		return fmt.Errorf("encryption failed: %w", err)
	}

	// Send the message
	reqBody := map[string]string{
		"to_username":       recipient.Username,
		"encrypted_content": base64.StdEncoding.EncodeToString(encryptedData),
		"content_type":      "text",
		"signature":         signature,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", baseURL+"/messages", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("send message failed (%d): %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// Encryption matching iOS CryptoService
func (c *TestClient) encrypt(content []byte, recipientPubKey [32]byte) ([]byte, string, error) {
	// 1. Compress content with raw deflate (iOS Compression framework uses raw deflate, not zlib)
	var compressed bytes.Buffer
	w, err := flate.NewWriter(&compressed, flate.DefaultCompression)
	if err != nil {
		return nil, "", fmt.Errorf("flate writer failed: %w", err)
	}
	if _, err := w.Write(content); err != nil {
		return nil, "", fmt.Errorf("compression failed: %w", err)
	}
	if err := w.Close(); err != nil {
		return nil, "", fmt.Errorf("compression close failed: %w", err)
	}
	compressedData := compressed.Bytes()

	// 2. Generate ephemeral symmetric key (256 bits)
	symmetricKey := make([]byte, 32)
	if _, err := rand.Read(symmetricKey); err != nil {
		return nil, "", err
	}

	// 3. Encrypt content with ChaCha20-Poly1305
	contentCipher, err := chacha20poly1305.New(symmetricKey)
	if err != nil {
		return nil, "", err
	}

	contentNonce := make([]byte, contentCipher.NonceSize())
	if _, err := rand.Read(contentNonce); err != nil {
		return nil, "", err
	}

	// ChaChaPoly format: nonce || ciphertext || tag (tag is appended by Seal)
	encryptedContent := contentCipher.Seal(contentNonce, contentNonce, compressedData, nil)

	// 4. Derive shared secret using ECDH
	sharedSecret, err := curve25519.X25519(c.privateKey[:], recipientPubKey[:])
	if err != nil {
		return nil, "", fmt.Errorf("ECDH failed: %w", err)
	}

	// 5. Derive key encryption key using HKDF-SHA256
	hkdfReader := hkdf.New(sha256.New, sharedSecret, nil, []byte("QuickPic-Key-Encryption"))
	derivedKey := make([]byte, 32)
	if _, err := io.ReadFull(hkdfReader, derivedKey); err != nil {
		return nil, "", fmt.Errorf("HKDF failed: %w", err)
	}

	// 6. Encrypt symmetric key with derived key
	keyCipher, err := chacha20poly1305.New(derivedKey)
	if err != nil {
		return nil, "", err
	}

	keyNonce := make([]byte, keyCipher.NonceSize())
	if _, err := rand.Read(keyNonce); err != nil {
		return nil, "", err
	}

	encryptedKey := keyCipher.Seal(keyNonce, keyNonce, symmetricKey, nil)

	// 7. Combine: [4-byte key length (big-endian)][encrypted key][encrypted content]
	result := make([]byte, 4+len(encryptedKey)+len(encryptedContent))
	binary.BigEndian.PutUint32(result[0:4], uint32(len(encryptedKey)))
	copy(result[4:], encryptedKey)
	copy(result[4+len(encryptedKey):], encryptedContent)

	// 8. Sign with Ed25519 (derived from Curve25519 private key)
	// Note: Ed25519 private key is the seed, same as X25519 private key bytes
	edPrivateKey := ed25519.NewKeyFromSeed(c.privateKey[:])
	signature := ed25519.Sign(edPrivateKey, result)

	return result, base64.StdEncoding.EncodeToString(signature), nil
}

// Decrypt - for testing that our encryption format is correct
func (c *TestClient) decrypt(encryptedData []byte, senderPubKey [32]byte) ([]byte, error) {
	// 1. Extract key length
	if len(encryptedData) < 4 {
		return nil, fmt.Errorf("data too short for key length")
	}
	keyLength := binary.BigEndian.Uint32(encryptedData[0:4])
	fmt.Printf("DEBUG: keyLength = %d\n", keyLength)

	if len(encryptedData) < int(4+keyLength) {
		return nil, fmt.Errorf("data too short for encrypted key")
	}

	// 2. Extract encrypted key and content
	encryptedKey := encryptedData[4 : 4+keyLength]
	encryptedContent := encryptedData[4+keyLength:]
	fmt.Printf("DEBUG: encryptedKey len = %d, encryptedContent len = %d\n", len(encryptedKey), len(encryptedContent))

	// 3. Derive shared secret (recipient private + sender public)
	sharedSecret, err := curve25519.X25519(c.privateKey[:], senderPubKey[:])
	if err != nil {
		return nil, fmt.Errorf("ECDH failed: %w", err)
	}

	// 4. Derive KEK
	hkdfReader := hkdf.New(sha256.New, sharedSecret, nil, []byte("QuickPic-Key-Encryption"))
	derivedKey := make([]byte, 32)
	if _, err := io.ReadFull(hkdfReader, derivedKey); err != nil {
		return nil, fmt.Errorf("HKDF failed: %w", err)
	}

	// 5. Decrypt symmetric key
	keyCipher, err := chacha20poly1305.New(derivedKey)
	if err != nil {
		return nil, err
	}
	if len(encryptedKey) < keyCipher.NonceSize() {
		return nil, fmt.Errorf("encrypted key too short for nonce")
	}
	keyNonce := encryptedKey[:keyCipher.NonceSize()]
	keyCiphertext := encryptedKey[keyCipher.NonceSize():]
	symmetricKey, err := keyCipher.Open(nil, keyNonce, keyCiphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt symmetric key: %w", err)
	}
	fmt.Printf("DEBUG: decrypted symmetricKey len = %d\n", len(symmetricKey))

	// 6. Decrypt content
	contentCipher, err := chacha20poly1305.New(symmetricKey)
	if err != nil {
		return nil, err
	}
	if len(encryptedContent) < contentCipher.NonceSize() {
		return nil, fmt.Errorf("encrypted content too short for nonce")
	}
	contentNonce := encryptedContent[:contentCipher.NonceSize()]
	contentCiphertext := encryptedContent[contentCipher.NonceSize():]
	compressedData, err := contentCipher.Open(nil, contentNonce, contentCiphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt content: %w", err)
	}
	fmt.Printf("DEBUG: decrypted compressed len = %d\n", len(compressedData))

	// 7. Decompress (raw deflate, not zlib)
	r := flate.NewReader(bytes.NewReader(compressedData))
	defer func() { _ = r.Close() }()
	plaintext, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("decompress failed: %w", err)
	}

	return plaintext, nil
}

func runDebug() error {
	fmt.Println("=== Crypto Debug ===")

	// Generate sender keypair
	var senderPrivate, senderPublic [32]byte
	if _, err := rand.Read(senderPrivate[:]); err != nil {
		return err
	}
	curve25519.ScalarBaseMult(&senderPublic, &senderPrivate)

	// Generate recipient keypair
	var recipientPrivate, recipientPublic [32]byte
	if _, err := rand.Read(recipientPrivate[:]); err != nil {
		return err
	}
	curve25519.ScalarBaseMult(&recipientPublic, &recipientPrivate)

	fmt.Printf("Sender X25519 private:  %s\n", base64.StdEncoding.EncodeToString(senderPrivate[:]))
	fmt.Printf("Sender X25519 public:   %s\n", base64.StdEncoding.EncodeToString(senderPublic[:]))
	fmt.Printf("Recipient X25519 public: %s\n", base64.StdEncoding.EncodeToString(recipientPublic[:]))

	// Ed25519 from sender private key
	edPrivateKey := ed25519.NewKeyFromSeed(senderPrivate[:])
	edPublicKey := edPrivateKey.Public().(ed25519.PublicKey)
	fmt.Printf("Sender Ed25519 public:  %s\n", base64.StdEncoding.EncodeToString(edPublicKey))
	fmt.Printf("\nNOTE: Ed25519 public != X25519 public (this is the problem!)\n")

	// Test ECDH
	sharedSecret1, _ := curve25519.X25519(senderPrivate[:], recipientPublic[:])
	sharedSecret2, _ := curve25519.X25519(recipientPrivate[:], senderPublic[:])
	fmt.Printf("\nECDH shared (sender):    %s\n", base64.StdEncoding.EncodeToString(sharedSecret1))
	fmt.Printf("ECDH shared (recipient): %s\n", base64.StdEncoding.EncodeToString(sharedSecret2))
	fmt.Printf("ECDH matches: %v\n", bytes.Equal(sharedSecret1, sharedSecret2))

	// Test encrypt/decrypt
	message := []byte("Hello, World!")
	fmt.Printf("\nOriginal message: %s\n", message)

	client := &TestClient{privateKey: senderPrivate, publicKey: senderPublic}
	encrypted, signature, err := client.encrypt(message, recipientPublic)
	if err != nil {
		return fmt.Errorf("encrypt failed: %w", err)
	}

	fmt.Printf("Encrypted length: %d bytes\n", len(encrypted))
	fmt.Printf("Signature: %s\n", signature)

	// Verify signature
	sigBytes, _ := base64.StdEncoding.DecodeString(signature)
	// This uses the Ed25519 public key, but iOS uses X25519 public key bytes!
	if ed25519.Verify(edPublicKey, encrypted, sigBytes) {
		fmt.Println("Signature valid (using Ed25519 public key)")
	} else {
		fmt.Println("Signature INVALID")
	}

	// Try verifying with X25519 public key bytes (what iOS does)
	// This will fail because X25519 public != Ed25519 public
	if ed25519.Verify(senderPublic[:], encrypted, sigBytes) {
		fmt.Println("Signature valid (using X25519 public key)")
	} else {
		fmt.Println("Signature INVALID with X25519 public key (expected - this is the bug)")
	}

	// Test decryption round-trip
	fmt.Println("\n=== Testing Decrypt ===")
	recipient := &TestClient{privateKey: recipientPrivate, publicKey: recipientPublic}
	decrypted, err := recipient.decrypt(encrypted, senderPublic)
	if err != nil {
		fmt.Printf("Decrypt FAILED: %v\n", err)
	} else {
		fmt.Printf("Decrypt SUCCESS: %s\n", string(decrypted))
	}

	return nil
}

// Credentials file handling
type savedCredentials struct {
	Username   string `json:"username"`
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
}

func (c *TestClient) saveCredentials() error {
	creds := savedCredentials{
		Username:   c.username,
		PrivateKey: base64.StdEncoding.EncodeToString(c.privateKey[:]),
		PublicKey:  base64.StdEncoding.EncodeToString(c.publicKey[:]),
	}

	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(c.username+".json", data, 0600)
}

func (c *TestClient) loadCredentials() error {
	data, err := os.ReadFile(c.username + ".json")
	if err != nil {
		return err
	}

	var creds savedCredentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return err
	}

	privKey, err := base64.StdEncoding.DecodeString(creds.PrivateKey)
	if err != nil {
		return err
	}
	copy(c.privateKey[:], privKey)

	pubKey, err := base64.StdEncoding.DecodeString(creds.PublicKey)
	if err != nil {
		return err
	}
	copy(c.publicKey[:], pubKey)

	// Login to get fresh token
	return c.login()
}
