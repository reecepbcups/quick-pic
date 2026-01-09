package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/quickpic/server/internal/models"
)

type MessageRepository struct {
	db *DB
}

func NewMessageRepository(db *DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(ctx context.Context, msg *models.Message) error {
	msg.ID = uuid.New()
	msg.CreatedAt = time.Now()

	query := `
		INSERT INTO messages (id, from_user_id, to_user_id, encrypted_content, content_type, signature, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		msg.ID.String(), msg.FromUserID.String(), msg.ToUserID.String(),
		msg.EncryptedContent, msg.ContentType, msg.Signature, msg.CreatedAt)

	return err
}

func (r *MessageRepository) GetPendingMessages(ctx context.Context, userID uuid.UUID) ([]models.MessageWithSender, error) {
	query := `
		SELECT m.id, m.from_user_id, m.to_user_id, m.encrypted_content, m.content_type, m.signature, m.created_at,
		       u.username, u.public_key
		FROM messages m
		JOIN users u ON u.id = m.from_user_id
		WHERE m.to_user_id = ?
		ORDER BY m.created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.MessageWithSender
	for rows.Next() {
		var msg models.MessageWithSender
		var idStr, fromUserIDStr, toUserIDStr string
		err := rows.Scan(
			&idStr, &fromUserIDStr, &toUserIDStr,
			&msg.EncryptedContent, &msg.ContentType, &msg.Signature, &msg.CreatedAt,
			&msg.FromUsername, &msg.FromPublicKey)
		if err != nil {
			return nil, err
		}
		msg.ID, _ = uuid.Parse(idStr)
		msg.FromUserID, _ = uuid.Parse(fromUserIDStr)
		msg.ToUserID, _ = uuid.Parse(toUserIDStr)
		messages = append(messages, msg)
	}

	return messages, rows.Err()
}

func (r *MessageRepository) GetByID(ctx context.Context, messageID uuid.UUID) (*models.Message, error) {
	query := `
		SELECT id, from_user_id, to_user_id, encrypted_content, content_type, signature, created_at
		FROM messages WHERE id = ?
	`

	var msg models.Message
	var idStr, fromUserIDStr, toUserIDStr string
	err := r.db.QueryRowContext(ctx, query, messageID.String()).Scan(
		&idStr, &fromUserIDStr, &toUserIDStr,
		&msg.EncryptedContent, &msg.ContentType, &msg.Signature, &msg.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrMessageNotFound
	}
	if err != nil {
		return nil, err
	}

	msg.ID, _ = uuid.Parse(idStr)
	msg.FromUserID, _ = uuid.Parse(fromUserIDStr)
	msg.ToUserID, _ = uuid.Parse(toUserIDStr)
	return &msg, nil
}

func (r *MessageRepository) Delete(ctx context.Context, messageID uuid.UUID) error {
	query := `DELETE FROM messages WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, messageID.String())
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return models.ErrMessageNotFound
	}

	return nil
}

func (r *MessageRepository) DeleteOldMessages(ctx context.Context, olderThan time.Duration) (int64, error) {
	query := `DELETE FROM messages WHERE created_at < ?`
	result, err := r.db.ExecContext(ctx, query, time.Now().Add(-olderThan))
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
