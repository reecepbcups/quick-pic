package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/quickpic/server/internal/models"
	"github.com/quickpic/server/internal/storage"
	"golang.org/x/crypto/argon2"
)

const (
	accessTokenDuration  = 15 * time.Minute
	refreshTokenDuration = 7 * 24 * time.Hour
)

type AuthService struct {
	userRepo  storage.UserRepo
	jwtSecret []byte
}

func NewAuthService(userRepo storage.UserRepo, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *AuthService) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	// Hash password with Argon2id
	passwordHash := s.hashPassword(req.Password)

	user := &models.User{
		Username:     req.Username,
		PasswordHash: passwordHash,
		PublicKey:    req.PublicKey,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return s.generateTokens(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, models.ErrInvalidCredentials
	}

	if !s.verifyPassword(req.Password, user.PasswordHash) {
		return nil, models.ErrInvalidCredentials
	}

	return s.generateTokens(ctx, user)
}

func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken string) (*models.AuthResponse, error) {
	tokenHash := s.hashToken(refreshToken)

	userID, err := s.userRepo.ValidateRefreshToken(ctx, tokenHash)
	if err != nil {
		return nil, err
	}

	// Delete the old refresh token (rotation)
	_ = s.userRepo.DeleteRefreshToken(ctx, tokenHash)

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.generateTokens(ctx, user)
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	tokenHash := s.hashToken(refreshToken)
	return s.userRepo.DeleteRefreshToken(ctx, tokenHash)
}

func (s *AuthService) ValidateAccessToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, models.ErrInvalidToken
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return uuid.Nil, models.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return uuid.Nil, models.ErrInvalidToken
	}

	// Check expiration
	if exp, ok := claims["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			return uuid.Nil, models.ErrTokenExpired
		}
	}

	userIDStr, ok := claims["sub"].(string)
	if !ok {
		return uuid.Nil, models.ErrInvalidToken
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, models.ErrInvalidToken
	}

	return userID, nil
}

func (s *AuthService) generateTokens(ctx context.Context, user *models.User) (*models.AuthResponse, error) {
	// Generate access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID.String(),
		"exp": time.Now().Add(accessTokenDuration).Unix(),
		"iat": time.Now().Unix(),
	})

	accessTokenString, err := accessToken.SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshTokenBytes := make([]byte, 32)
	if _, err := rand.Read(refreshTokenBytes); err != nil {
		return nil, err
	}
	refreshToken := base64.URLEncoding.EncodeToString(refreshTokenBytes)

	// Store refresh token hash
	tokenHash := s.hashToken(refreshToken)
	expiresAt := time.Now().Add(refreshTokenDuration)
	if err := s.userRepo.StoreRefreshToken(ctx, user.ID, tokenHash, expiresAt); err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(accessTokenDuration.Seconds()),
		User:         user.ToPublic(),
	}, nil
}

func (s *AuthService) hashPassword(password string) string {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		panic("failed to generate random salt: " + err.Error())
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	// Encode as base64: salt + hash
	combined := append(salt, hash...)
	return base64.StdEncoding.EncodeToString(combined)
}

func (s *AuthService) verifyPassword(password, encodedHash string) bool {
	combined, err := base64.StdEncoding.DecodeString(encodedHash)
	if err != nil || len(combined) < 48 {
		return false
	}

	salt := combined[:16]
	storedHash := combined[16:]

	computedHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	// Constant-time comparison
	if len(storedHash) != len(computedHash) {
		return false
	}
	var result byte
	for i := range storedHash {
		result |= storedHash[i] ^ computedHash[i]
	}
	return result == 0
}

func (s *AuthService) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.StdEncoding.EncodeToString(hash[:])
}
