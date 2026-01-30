package notification

import (
	"context"
	"fmt"

	"github.com/yourusername/online-library/internal/domain"
	"go.uber.org/zap"
)

type service struct {
	notificationRepo NotificationRepo
	log              *zap.Logger
}

func NewService(notificationRepo NotificationRepo, log *zap.Logger) Service {
	return &service{
		notificationRepo: notificationRepo,
		log:              log,
	}
}

func (s *service) NotifyIdeaVote(ctx context.Context, userID, voterName, ideaTitle string, isUpvote bool) error {
	voteType := "upvoted"
	if !isUpvote {
		voteType = "downvoted"
	}
	return s.notificationRepo.Create(
		ctx,
		userID,
		"idea_vote",
		"Idea Vote",
		fmt.Sprintf("%s %s your idea '%s'", voterName, voteType, ideaTitle),
		"/ideas",
	)
}

func (s *service) NotifyReviewReceived(ctx context.Context, userID, reviewerName string) error {
	return s.notificationRepo.Create(
		ctx,
		userID,
		"review_received",
		"New Review",
		fmt.Sprintf("%s left you a review!", reviewerName),
		"/profile/reviews",
	)
}

func (s *service) NotifyBookAvailable(ctx context.Context, userID, bookID, bookTitle string) error {
	return s.notificationRepo.Create(
		ctx,
		userID,
		"book_available",
		"Book Available",
		fmt.Sprintf("The book '%s' is now available for you!", bookTitle),
		fmt.Sprintf("/books/%s", bookID),
	)
}

func (s *service) NotifyRequestApproved(ctx context.Context, userID, bookID, bookTitle string) error {
	return s.notificationRepo.Create(
		ctx,
		userID,
		"request_approved",
		"Request Approved",
		fmt.Sprintf("Your request for '%s' has been approved!", bookTitle),
		fmt.Sprintf("/books/%s", bookID),
	)
}

func (s *service) NotifyReturnDue(ctx context.Context, userID, bookID, bookTitle string, daysLeft int) error {
	return s.notificationRepo.Create(
		ctx,
		userID,
		"return_due",
		"Book Return Reminder",
		fmt.Sprintf("Please return '%s' in %d days.", bookTitle, daysLeft),
		fmt.Sprintf("/books/%s", bookID),
	)
}

func (s *service) GetUserNotifications(ctx context.Context, userID string, limit int) ([]*domain.Notification, error) {
	return s.notificationRepo.GetByUserID(ctx, userID, limit)
}

func (s *service) MarkAsRead(ctx context.Context, notificationID string) error {
	return s.notificationRepo.MarkAsRead(ctx, notificationID)
}

func (s *service) MarkAllAsRead(ctx context.Context, userID string) error {
	return s.notificationRepo.MarkAllAsRead(ctx, userID)
}
