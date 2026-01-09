package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/quickpic/server/internal/models"
)

// Backend implements repository.Backend for blockchain storage
type Backend struct {
	client          *ethclient.Client
	contract        *QuickPicStorage
	contractAddress common.Address
	privateKey      *ecdsa.PrivateKey
	chainID         *big.Int

	users    *UserRepository
	friends  *FriendRepository
	messages *MessageRepository

	// In-memory refresh token storage (not stored on-chain)
	refreshTokens     map[string]refreshTokenEntry
	refreshTokenMutex sync.RWMutex
}

type refreshTokenEntry struct {
	userID    uuid.UUID
	expiresAt time.Time
}

// Config holds the configuration for the blockchain backend
type Config struct {
	RPCURL          string // e.g., "http://localhost:8545" for Anvil
	PrivateKey      string // hex-encoded private key (without 0x prefix)
	ContractAddress string // deployed contract address
}

// NewBackend creates a new blockchain backend
func NewBackend(cfg Config) (*Backend, error) {
	client, err := ethclient.Dial(cfg.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(cfg.PrivateKey, "0x"))
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	contractAddress := common.HexToAddress(cfg.ContractAddress)
	contract, err := NewQuickPicStorage(contractAddress, client)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to bind contract: %w", err)
	}

	backend := &Backend{
		client:          client,
		contract:        contract,
		contractAddress: contractAddress,
		privateKey:      privateKey,
		chainID:         chainID,
		refreshTokens:   make(map[string]refreshTokenEntry),
	}

	backend.users = &UserRepository{backend: backend}
	backend.friends = &FriendRepository{backend: backend}
	backend.messages = &MessageRepository{backend: backend}

	return backend, nil
}

func (b *Backend) Users() *UserRepository {
	return b.users
}

func (b *Backend) Friends() *FriendRepository {
	return b.friends
}

func (b *Backend) Messages() *MessageRepository {
	return b.messages
}

func (b *Backend) Close() error {
	b.client.Close()
	return nil
}

func (b *Backend) Name() string {
	return "blockchain"
}

func (b *Backend) getTransactOpts(ctx context.Context) (*bind.TransactOpts, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(b.privateKey, b.chainID)
	if err != nil {
		return nil, err
	}
	auth.Context = ctx
	return auth, nil
}

// uuidToBytes32 converts a UUID to a bytes32 value
func uuidToBytes32(id uuid.UUID) [32]byte {
	var result [32]byte
	copy(result[:], id[:])
	return result
}

// bytes32ToUUID converts a bytes32 value back to a UUID
func bytes32ToUUID(b [32]byte) uuid.UUID {
	var id uuid.UUID
	copy(id[:], b[:16])
	return id
}

// ============ UserRepository ============

type UserRepository struct {
	backend *Backend
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	user.ID = uuid.New()
	user.Username = strings.ToLower(user.Username)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	auth, err := r.backend.getTransactOpts(ctx)
	if err != nil {
		return err
	}

	tx, err := r.backend.contract.CreateUser(
		auth,
		uuidToBytes32(user.ID),
		user.Username,
		user.PasswordHash,
		user.PublicKey,
	)
	if err != nil {
		if strings.Contains(err.Error(), "Username already taken") {
			return models.ErrUsernameExists
		}
		return err
	}

	// Wait for transaction to be mined
	receipt, err := bind.WaitMined(ctx, r.backend.client, tx)
	if err != nil {
		return err
	}
	if receipt.Status == 0 {
		return fmt.Errorf("transaction failed")
	}

	// Get the user number from the contract
	result, err := r.backend.contract.GetUser(&bind.CallOpts{Context: ctx}, uuidToBytes32(user.ID))
	if err != nil {
		return err
	}
	user.UserNumber = result.UserNumber.Int64()

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	result, err := r.backend.contract.GetUser(&bind.CallOpts{Context: ctx}, uuidToBytes32(id))
	if err != nil {
		if strings.Contains(err.Error(), "User not found") {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	return &models.User{
		ID:           bytes32ToUUID(result.UserId),
		UserNumber:   result.UserNumber.Int64(),
		Username:     result.Username,
		PasswordHash: result.PasswordHash,
		PublicKey:    result.PublicKey,
		CreatedAt:    time.Unix(result.CreatedAt.Int64(), 0),
		UpdatedAt:    time.Unix(result.UpdatedAt.Int64(), 0),
	}, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	result, err := r.backend.contract.GetUserByUsername(&bind.CallOpts{Context: ctx}, strings.ToLower(username))
	if err != nil {
		if strings.Contains(err.Error(), "User not found") {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	return &models.User{
		ID:           bytes32ToUUID(result.UserId),
		UserNumber:   result.UserNumber.Int64(),
		Username:     result.UsernameOut,
		PasswordHash: result.PasswordHash,
		PublicKey:    result.PublicKey,
		CreatedAt:    time.Unix(result.CreatedAt.Int64(), 0),
		UpdatedAt:    time.Unix(result.UpdatedAt.Int64(), 0),
	}, nil
}

// Refresh tokens are stored in-memory (not on-chain) for performance
func (r *UserRepository) StoreRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	r.backend.refreshTokenMutex.Lock()
	defer r.backend.refreshTokenMutex.Unlock()

	r.backend.refreshTokens[tokenHash] = refreshTokenEntry{
		userID:    userID,
		expiresAt: expiresAt,
	}
	return nil
}

func (r *UserRepository) ValidateRefreshToken(ctx context.Context, tokenHash string) (uuid.UUID, error) {
	r.backend.refreshTokenMutex.RLock()
	defer r.backend.refreshTokenMutex.RUnlock()

	entry, ok := r.backend.refreshTokens[tokenHash]
	if !ok || time.Now().After(entry.expiresAt) {
		return uuid.Nil, models.ErrInvalidToken
	}
	return entry.userID, nil
}

func (r *UserRepository) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	r.backend.refreshTokenMutex.Lock()
	defer r.backend.refreshTokenMutex.Unlock()

	delete(r.backend.refreshTokens, tokenHash)
	return nil
}

func (r *UserRepository) DeleteAllRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	r.backend.refreshTokenMutex.Lock()
	defer r.backend.refreshTokenMutex.Unlock()

	for hash, entry := range r.backend.refreshTokens {
		if entry.userID == userID {
			delete(r.backend.refreshTokens, hash)
		}
	}
	return nil
}

// ============ FriendRepository ============

type FriendRepository struct {
	backend *Backend
}

func (r *FriendRepository) CreateFriendRequest(ctx context.Context, fromUserID, toUserID uuid.UUID) (*models.FriendRequest, error) {
	// Check if already friends
	areFriends, err := r.AreFriends(ctx, fromUserID, toUserID)
	if err != nil {
		return nil, err
	}
	if areFriends {
		return nil, models.ErrAlreadyFriends
	}

	requestID := uuid.New()

	auth, err := r.backend.getTransactOpts(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := r.backend.contract.CreateFriendRequest(
		auth,
		uuidToBytes32(requestID),
		uuidToBytes32(fromUserID),
		uuidToBytes32(toUserID),
	)
	if err != nil {
		errStr := strings.ToLower(err.Error())
		if strings.Contains(errStr, "request already exists") {
			return nil, models.ErrFriendRequestExists
		}
		if strings.Contains(errStr, "already friends") {
			return nil, models.ErrAlreadyFriends
		}
		if strings.Contains(errStr, "cannot send request to self") {
			return nil, models.ErrCannotAddSelf
		}
		return nil, fmt.Errorf("blockchain error: %w", err)
	}

	receipt, err := bind.WaitMined(ctx, r.backend.client, tx)
	if err != nil {
		return nil, err
	}
	if receipt.Status == 0 {
		return nil, fmt.Errorf("transaction failed")
	}

	return &models.FriendRequest{
		ID:         requestID,
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Status:     models.FriendRequestPending,
		CreatedAt:  time.Now(),
	}, nil
}

func (r *FriendRepository) GetPendingRequests(ctx context.Context, userID uuid.UUID) ([]models.FriendRequestWithUser, error) {
	requestIds, err := r.backend.contract.GetPendingRequestsForUser(&bind.CallOpts{Context: ctx}, uuidToBytes32(userID))
	if err != nil {
		return nil, err
	}

	var requests []models.FriendRequestWithUser
	for _, reqID := range requestIds {
		result, err := r.backend.contract.GetFriendRequest(&bind.CallOpts{Context: ctx}, reqID)
		if err != nil {
			continue
		}

		fromUserID := bytes32ToUUID(result.FromUserId)
		fromUser, err := r.backend.users.GetByID(ctx, fromUserID)
		if err != nil {
			continue
		}

		req := models.FriendRequestWithUser{
			FriendRequest: models.FriendRequest{
				ID:         bytes32ToUUID(result.RequestId),
				FromUserID: fromUserID,
				ToUserID:   bytes32ToUUID(result.ToUserId),
				Status:     models.FriendRequestStatus([]string{"pending", "accepted", "rejected"}[result.Status]),
				CreatedAt:  time.Unix(result.CreatedAt.Int64(), 0),
			},
			FromUser: models.UserPublic{
				ID:        fromUser.ID,
				Username:  fromUser.Username,
				PublicKey: fromUser.PublicKey,
			},
		}
		requests = append(requests, req)
	}

	return requests, nil
}

func (r *FriendRepository) GetFriendRequest(ctx context.Context, requestID uuid.UUID) (*models.FriendRequest, error) {
	result, err := r.backend.contract.GetFriendRequest(&bind.CallOpts{Context: ctx}, uuidToBytes32(requestID))
	if err != nil {
		if strings.Contains(err.Error(), "Friend request not found") {
			return nil, models.ErrFriendRequestNotFound
		}
		return nil, err
	}

	return &models.FriendRequest{
		ID:         bytes32ToUUID(result.RequestId),
		FromUserID: bytes32ToUUID(result.FromUserId),
		ToUserID:   bytes32ToUUID(result.ToUserId),
		Status:     models.FriendRequestStatus([]string{"pending", "accepted", "rejected"}[result.Status]),
		CreatedAt:  time.Unix(result.CreatedAt.Int64(), 0),
	}, nil
}

func (r *FriendRepository) UpdateFriendRequestStatus(ctx context.Context, requestID uuid.UUID, status models.FriendRequestStatus) error {
	auth, err := r.backend.getTransactOpts(ctx)
	if err != nil {
		return err
	}

	var statusVal uint8
	switch status {
	case models.FriendRequestPending:
		statusVal = 0
	case models.FriendRequestAccepted:
		statusVal = 1
	case models.FriendRequestRejected:
		statusVal = 2
	}

	tx, err := r.backend.contract.UpdateFriendRequestStatus(auth, uuidToBytes32(requestID), statusVal)
	if err != nil {
		if strings.Contains(err.Error(), "Friend request not found") {
			return models.ErrFriendRequestNotFound
		}
		return err
	}

	receipt, err := bind.WaitMined(ctx, r.backend.client, tx)
	if err != nil {
		return err
	}
	if receipt.Status == 0 {
		return fmt.Errorf("transaction failed")
	}

	return nil
}

func (r *FriendRepository) CreateFriendship(ctx context.Context, userAID, userBID uuid.UUID) error {
	friendshipID := uuid.New()

	auth, err := r.backend.getTransactOpts(ctx)
	if err != nil {
		return err
	}

	tx, err := r.backend.contract.CreateFriendship(
		auth,
		uuidToBytes32(friendshipID),
		uuidToBytes32(userAID),
		uuidToBytes32(userBID),
	)
	if err != nil {
		return err
	}

	receipt, err := bind.WaitMined(ctx, r.backend.client, tx)
	if err != nil {
		return err
	}
	if receipt.Status == 0 {
		return fmt.Errorf("transaction failed")
	}

	return nil
}

func (r *FriendRepository) AreFriends(ctx context.Context, userAID, userBID uuid.UUID) (bool, error) {
	return r.backend.contract.AreFriends(&bind.CallOpts{Context: ctx}, uuidToBytes32(userAID), uuidToBytes32(userBID))
}

func (r *FriendRepository) GetFriends(ctx context.Context, userID uuid.UUID) ([]models.Friend, error) {
	friendIds, err := r.backend.contract.GetFriendsOfUser(&bind.CallOpts{Context: ctx}, uuidToBytes32(userID))
	if err != nil {
		return nil, err
	}

	var friends []models.Friend
	for _, friendID := range friendIds {
		friendUser, err := r.backend.users.GetByID(ctx, bytes32ToUUID(friendID))
		if err != nil {
			continue
		}

		// Get friendship creation time
		friendships, err := r.backend.contract.GetUserFriendships(&bind.CallOpts{Context: ctx}, uuidToBytes32(userID))
		var since time.Time
		for _, fshipID := range friendships {
			fship, err := r.backend.contract.GetFriendship(&bind.CallOpts{Context: ctx}, fshipID)
			if err != nil {
				continue
			}
			if bytes32ToUUID(fship.UserAId) == friendUser.ID || bytes32ToUUID(fship.UserBId) == friendUser.ID {
				since = time.Unix(fship.CreatedAt.Int64(), 0)
				break
			}
		}

		friends = append(friends, models.Friend{
			UserID:    friendUser.ID,
			Username:  friendUser.Username,
			PublicKey: friendUser.PublicKey,
			Since:     since,
		})
	}

	return friends, nil
}

// ============ MessageRepository ============

type MessageRepository struct {
	backend *Backend
}

func (r *MessageRepository) Create(ctx context.Context, msg *models.Message) error {
	msg.ID = uuid.New()
	msg.CreatedAt = time.Now()

	auth, err := r.backend.getTransactOpts(ctx)
	if err != nil {
		return err
	}

	var contentType uint8
	switch msg.ContentType {
	case models.ContentTypeText:
		contentType = 0
	case models.ContentTypeImage:
		contentType = 1
	}

	tx, err := r.backend.contract.CreateMessage(
		auth,
		uuidToBytes32(msg.ID),
		uuidToBytes32(msg.FromUserID),
		uuidToBytes32(msg.ToUserID),
		msg.EncryptedContent,
		contentType,
		msg.Signature,
	)
	if err != nil {
		return err
	}

	receipt, err := bind.WaitMined(ctx, r.backend.client, tx)
	if err != nil {
		return err
	}
	if receipt.Status == 0 {
		return fmt.Errorf("transaction failed")
	}

	return nil
}

func (r *MessageRepository) GetPendingMessages(ctx context.Context, userID uuid.UUID) ([]models.MessageWithSender, error) {
	messageIds, err := r.backend.contract.GetMessagesForUser(&bind.CallOpts{Context: ctx}, uuidToBytes32(userID))
	if err != nil {
		return nil, err
	}

	var messages []models.MessageWithSender
	for _, msgID := range messageIds {
		result, err := r.backend.contract.GetMessage(&bind.CallOpts{Context: ctx}, msgID)
		if err != nil {
			continue
		}

		fromUser, err := r.backend.users.GetByID(ctx, bytes32ToUUID(result.FromUserId))
		if err != nil {
			continue
		}

		msg := models.MessageWithSender{
			Message: models.Message{
				ID:               bytes32ToUUID(result.MessageId),
				FromUserID:       bytes32ToUUID(result.FromUserId),
				ToUserID:         bytes32ToUUID(result.ToUserId),
				EncryptedContent: result.EncryptedContent,
				ContentType:      models.ContentType([]string{"text", "image"}[result.ContentType]),
				Signature:        result.Signature,
				CreatedAt:        time.Unix(result.CreatedAt.Int64(), 0),
			},
			FromUsername:  fromUser.Username,
			FromPublicKey: fromUser.PublicKey,
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func (r *MessageRepository) GetByID(ctx context.Context, messageID uuid.UUID) (*models.Message, error) {
	result, err := r.backend.contract.GetMessage(&bind.CallOpts{Context: ctx}, uuidToBytes32(messageID))
	if err != nil {
		if strings.Contains(err.Error(), "Message not found") {
			return nil, models.ErrMessageNotFound
		}
		return nil, err
	}

	return &models.Message{
		ID:               bytes32ToUUID(result.MessageId),
		FromUserID:       bytes32ToUUID(result.FromUserId),
		ToUserID:         bytes32ToUUID(result.ToUserId),
		EncryptedContent: result.EncryptedContent,
		ContentType:      models.ContentType([]string{"text", "image"}[result.ContentType]),
		Signature:        result.Signature,
		CreatedAt:        time.Unix(result.CreatedAt.Int64(), 0),
	}, nil
}

func (r *MessageRepository) Delete(ctx context.Context, messageID uuid.UUID) error {
	auth, err := r.backend.getTransactOpts(ctx)
	if err != nil {
		return err
	}

	tx, err := r.backend.contract.DeleteMessage(auth, uuidToBytes32(messageID))
	if err != nil {
		if strings.Contains(err.Error(), "Message not found") {
			return models.ErrMessageNotFound
		}
		return err
	}

	receipt, err := bind.WaitMined(ctx, r.backend.client, tx)
	if err != nil {
		return err
	}
	if receipt.Status == 0 {
		return fmt.Errorf("transaction failed")
	}

	return nil
}

func (r *MessageRepository) DeleteOldMessages(ctx context.Context, olderThan time.Duration) (int64, error) {
	// Note: This would require iterating through all messages on-chain,
	// which is expensive. For now, we return 0 as this operation is
	// better handled by off-chain indexing or a scheduled task.
	return 0, nil
}
