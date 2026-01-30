package repository

import (
	"context"
	"database/sql"

	"github.com/yourusername/online-library/internal/domain"
	"github.com/yourusername/online-library/internal/notification"
	"go.uber.org/zap"
)

type NotificationRepository struct {
	db  *sql.DB
	log *zap.Logger
}

var _ notification.NotificationRepo = (*NotificationRepository)(nil)

func NewNotificationRepository(db *sql.DB, log *zap.Logger) *NotificationRepository {
	return &NotificationRepository{db: db, log: log}
}

func (r *NotificationRepository) Create(ctx context.Context, userID, notifType, title, message, link string) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO notifications (user_id, type, title, message, link, created_at)
	                                  VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)`,
		userID, notifType, title, message, link)
	return err
}

func (r *NotificationRepository) GetByUserID(ctx context.Context, userID string, limit int) ([]*domain.Notification, error) {
	query := `
		SELECT id, user_id, type, title, message, COALESCE(link, ''), is_read, created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`
	rows, err := r.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*domain.Notification
	for rows.Next() {
		n := &domain.Notification{}
		err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Message, &n.Link, &n.IsRead, &n.CreatedAt)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	return notifications, nil
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, notificationID string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE notifications SET is_read = true WHERE id = $1`, notificationID)
	return err
}

func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE notifications SET is_read = true WHERE user_id = $1`, userID)
	return err
}
