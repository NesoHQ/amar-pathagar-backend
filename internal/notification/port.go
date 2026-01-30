package notification

import (
	"context"

	"github.com/yourusername/online-library/internal/domain"
)

type Service interface {
	NotifyIdeaVote(ctx context.Context, userID, voterName, ideaTitle string, isUpvote bool) error
	NotifyReviewReceived(ctx context.Context, userID, reviewerName string) error
	NotifyBookAvailable(ctx context.Context, userID, bookID, bookTitle string) error
	NotifyRequestApproved(ctx context.Context, userID, bookID, bookTitle string) error
	NotifyReturnDue(ctx context.Context, userID, bookID, bookTitle string, daysLeft int) error
	GetUserNotifications(ctx context.Context, userID string, limit int) ([]*domain.Notification, error)
	MarkAsRead(ctx context.Context, notificationID string) error
	MarkAllAsRead(ctx context.Context, userID string) error
}

type NotificationRepo interface {
	Create(ctx context.Context, userID, notifType, title, message, link string) error
	GetByUserID(ctx context.Context, userID string, limit int) ([]*domain.Notification, error)
	MarkAsRead(ctx context.Context, notificationID string) error
	MarkAllAsRead(ctx context.Context, userID string) error
}
