package successscore

import "context"

type Service interface {
	ProcessIdeaPosted(ctx context.Context, userID, ideaID string) error
	ProcessIdeaUpvote(ctx context.Context, userID, ideaID string) error
	ProcessIdeaDownvote(ctx context.Context, userID, ideaID string) error
	ProcessPositiveReview(ctx context.Context, userID, reviewID string) error
	ProcessNegativeReview(ctx context.Context, userID, reviewID string) error
	ProcessBookDonation(ctx context.Context, userID, donationID string) error
	ProcessMoneyDonation(ctx context.Context, userID, donationID string) error
	ProcessReturnOnTime(ctx context.Context, userID, bookID string) error
	ProcessReturnLate(ctx context.Context, userID, bookID string) error
	ProcessLostBook(ctx context.Context, userID, bookID string) error
}

type ScoreRepo interface {
	UpdateScore(ctx context.Context, userID string, change int, reason, refType string, refID *string) error
}
