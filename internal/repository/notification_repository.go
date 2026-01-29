package repository

import (
	"context"
	"database/sql"

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
