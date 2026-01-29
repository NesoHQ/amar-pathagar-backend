package review

import (
	"context"

	"github.com/yourusername/online-library/internal/domain"
)

type Service interface {
	Create(ctx context.Context, review *domain.UserReview) (*domain.UserReview, error)
	GetByUser(ctx context.Context, userID string) ([]*domain.UserReview, error)
}

type ReviewRepo interface {
	Create(ctx context.Context, review *domain.UserReview) error
	FindByUserID(ctx context.Context, userID string) ([]*domain.UserReview, error)
}

type SuccessScoreSvc interface {
	ProcessPositiveReview(ctx context.Context, userID, reviewID string) error
	ProcessNegativeReview(ctx context.Context, userID, reviewID string) error
}

type NotificationSvc interface {
	NotifyReviewReceived(ctx context.Context, userID, reviewerName string) error
}
