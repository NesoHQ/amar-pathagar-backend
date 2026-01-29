package review

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/online-library/internal/domain"
	"go.uber.org/zap"
)

type service struct {
	reviewRepo      ReviewRepo
	successScoreSvc SuccessScoreSvc
	notificationSvc NotificationSvc
	log             *zap.Logger
}

func NewService(reviewRepo ReviewRepo, successScoreSvc SuccessScoreSvc, notificationSvc NotificationSvc, log *zap.Logger) Service {
	return &service{
		reviewRepo:      reviewRepo,
		successScoreSvc: successScoreSvc,
		notificationSvc: notificationSvc,
		log:             log,
	}
}

func (s *service) Create(ctx context.Context, review *domain.UserReview) (*domain.UserReview, error) {
	review.ID = uuid.New().String()
	review.CreatedAt = time.Now()

	if err := s.reviewRepo.Create(ctx, review); err != nil {
		s.log.Error("failed to create review", zap.Error(err))
		return nil, err
	}

	// Calculate average rating and update success score
	avgRating := 0
	count := 0
	if review.BehaviorRating != nil {
		avgRating += *review.BehaviorRating
		count++
	}
	if review.BookConditionRating != nil {
		avgRating += *review.BookConditionRating
		count++
	}
	if review.CommunicationRating != nil {
		avgRating += *review.CommunicationRating
		count++
	}

	if count > 0 {
		avgRating = avgRating / count
		if avgRating >= 4 {
			if err := s.successScoreSvc.ProcessPositiveReview(ctx, review.RevieweeID, review.ID); err != nil {
				s.log.Warn("failed to update success score for positive review", zap.Error(err))
			}
		} else if avgRating < 3 {
			if err := s.successScoreSvc.ProcessNegativeReview(ctx, review.RevieweeID, review.ID); err != nil {
				s.log.Warn("failed to update success score for negative review", zap.Error(err))
			}
		}
	}

	s.log.Info("review created successfully", zap.String("review_id", review.ID))
	return review, nil
}

func (s *service) GetByUser(ctx context.Context, userID string) ([]*domain.UserReview, error) {
	reviews, err := s.reviewRepo.FindByUserID(ctx, userID)
	if err != nil {
		s.log.Error("failed to get reviews", zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}
	return reviews, nil
}
