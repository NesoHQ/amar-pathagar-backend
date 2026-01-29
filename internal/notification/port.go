package notification

import "context"

type Service interface {
	NotifyIdeaVote(ctx context.Context, userID, voterName, ideaTitle string, isUpvote bool) error
	NotifyReviewReceived(ctx context.Context, userID, reviewerName string) error
	NotifyBookAvailable(ctx context.Context, userID, bookID, bookTitle string) error
	NotifyRequestApproved(ctx context.Context, userID, bookID, bookTitle string) error
	NotifyReturnDue(ctx context.Context, userID, bookID, bookTitle string, daysLeft int) error
}

type NotificationRepo interface {
	Create(ctx context.Context, userID, notifType, title, message, link string) error
}
