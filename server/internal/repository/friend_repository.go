package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/quickpic/server/internal/models"
)

type FriendRepository struct {
	db *DB
}

func NewFriendRepository(db *DB) *FriendRepository {
	return &FriendRepository{db: db}
}

func (r *FriendRepository) CreateFriendRequest(ctx context.Context, fromUserID, toUserID uuid.UUID) (*models.FriendRequest, error) {
	// Check if request already exists
	existingQuery := `
		SELECT id FROM friend_requests
		WHERE from_user_id = ? AND to_user_id = ? AND status = 'pending'
	`
	var existingID string
	err := r.db.QueryRowContext(ctx, existingQuery, fromUserID.String(), toUserID.String()).Scan(&existingID)
	if err == nil {
		return nil, models.ErrFriendRequestExists
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// Check if reverse request exists (they already sent us a request)
	err = r.db.QueryRowContext(ctx, existingQuery, toUserID.String(), fromUserID.String()).Scan(&existingID)
	if err == nil {
		return nil, models.ErrFriendRequestExists
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// Check if already friends
	if areFriends, _ := r.AreFriends(ctx, fromUserID, toUserID); areFriends {
		return nil, models.ErrAlreadyFriends
	}

	request := &models.FriendRequest{
		ID:         uuid.New(),
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Status:     models.FriendRequestPending,
		CreatedAt:  time.Now(),
	}

	query := `
		INSERT INTO friend_requests (id, from_user_id, to_user_id, status, created_at)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err = r.db.ExecContext(ctx, query,
		request.ID.String(), request.FromUserID.String(), request.ToUserID.String(), request.Status, request.CreatedAt)
	if err != nil {
		return nil, err
	}

	return request, nil
}

func (r *FriendRepository) GetPendingRequests(ctx context.Context, userID uuid.UUID) ([]models.FriendRequestWithUser, error) {
	query := `
		SELECT fr.id, fr.from_user_id, fr.to_user_id, fr.status, fr.created_at,
		       u.id, u.username, u.public_key
		FROM friend_requests fr
		JOIN users u ON u.id = fr.from_user_id
		WHERE fr.to_user_id = ? AND fr.status = 'pending'
		ORDER BY fr.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []models.FriendRequestWithUser
	for rows.Next() {
		var req models.FriendRequestWithUser
		var idStr, fromUserIDStr, toUserIDStr, fromUserIDStr2 string
		err := rows.Scan(
			&idStr, &fromUserIDStr, &toUserIDStr, &req.Status, &req.CreatedAt,
			&fromUserIDStr2, &req.FromUser.Username, &req.FromUser.PublicKey)
		if err != nil {
			return nil, err
		}
		req.ID, _ = uuid.Parse(idStr)
		req.FromUserID, _ = uuid.Parse(fromUserIDStr)
		req.ToUserID, _ = uuid.Parse(toUserIDStr)
		req.FromUser.ID, _ = uuid.Parse(fromUserIDStr2)
		requests = append(requests, req)
	}

	return requests, rows.Err()
}

func (r *FriendRepository) GetFriendRequest(ctx context.Context, requestID uuid.UUID) (*models.FriendRequest, error) {
	query := `
		SELECT id, from_user_id, to_user_id, status, created_at
		FROM friend_requests WHERE id = ?
	`

	var req models.FriendRequest
	var idStr, fromUserIDStr, toUserIDStr string
	err := r.db.QueryRowContext(ctx, query, requestID.String()).Scan(
		&idStr, &fromUserIDStr, &toUserIDStr, &req.Status, &req.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrFriendRequestNotFound
	}
	if err != nil {
		return nil, err
	}

	req.ID, _ = uuid.Parse(idStr)
	req.FromUserID, _ = uuid.Parse(fromUserIDStr)
	req.ToUserID, _ = uuid.Parse(toUserIDStr)
	return &req, nil
}

func (r *FriendRepository) UpdateFriendRequestStatus(ctx context.Context, requestID uuid.UUID, status models.FriendRequestStatus) error {
	query := `UPDATE friend_requests SET status = ? WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, status, requestID.String())
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return models.ErrFriendRequestNotFound
	}

	return nil
}

func (r *FriendRepository) CreateFriendship(ctx context.Context, userAID, userBID uuid.UUID) error {
	// Ensure user_a_id < user_b_id for the unique constraint
	if userAID.String() > userBID.String() {
		userAID, userBID = userBID, userAID
	}

	query := `
		INSERT OR IGNORE INTO friendships (id, user_a_id, user_b_id, created_at)
		VALUES (?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query, uuid.New().String(), userAID.String(), userBID.String(), time.Now())
	return err
}

func (r *FriendRepository) AreFriends(ctx context.Context, userAID, userBID uuid.UUID) (bool, error) {
	// Ensure user_a_id < user_b_id
	if userAID.String() > userBID.String() {
		userAID, userBID = userBID, userAID
	}

	query := `SELECT 1 FROM friendships WHERE user_a_id = ? AND user_b_id = ?`
	var exists int
	err := r.db.QueryRowContext(ctx, query, userAID.String(), userBID.String()).Scan(&exists)

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *FriendRepository) GetFriends(ctx context.Context, userID uuid.UUID) ([]models.Friend, error) {
	query := `
		SELECT u.id, u.username, u.public_key, f.created_at
		FROM friendships f
		JOIN users u ON (
			(f.user_a_id = ? AND u.id = f.user_b_id) OR
			(f.user_b_id = ? AND u.id = f.user_a_id)
		)
		WHERE f.user_a_id = ? OR f.user_b_id = ?
		ORDER BY u.username
	`

	userIDStr := userID.String()
	rows, err := r.db.QueryContext(ctx, query, userIDStr, userIDStr, userIDStr, userIDStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []models.Friend
	for rows.Next() {
		var f models.Friend
		var userIDStr string
		if err := rows.Scan(&userIDStr, &f.Username, &f.PublicKey, &f.Since); err != nil {
			return nil, err
		}
		f.UserID, _ = uuid.Parse(userIDStr)
		friends = append(friends, f)
	}

	return friends, rows.Err()
}
