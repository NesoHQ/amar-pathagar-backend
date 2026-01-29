package successscore

import (
	"context"

	"go.uber.org/zap"
)

const (
	ScoreReturnOnTime   = 10
	ScoreReturnLate     = -15
	ScorePositiveReview = 5
	ScoreNegativeReview = -10
	ScoreIdeaPosted     = 3
	ScoreIdeaUpvote     = 1
	ScoreIdeaDownvote   = -1
	ScoreLostBook       = -50
	ScoreBookDonated    = 20
	ScoreMoneyDonated   = 10
)

type service struct {
	scoreRepo ScoreRepo
	log       *zap.Logger
}

func NewService(scoreRepo ScoreRepo, log *zap.Logger) Service {
	return &service{
		scoreRepo: scoreRepo,
		log:       log,
	}
}

func (s *service) ProcessIdeaPosted(ctx context.Context, userID, ideaID string) error {
	return s.scoreRepo.UpdateScore(ctx, userID, ScoreIdeaPosted, "Posted reading idea", "idea", &ideaID)
}

func (s *service) ProcessIdeaUpvote(ctx context.Context, userID, ideaID string) error {
	return s.scoreRepo.UpdateScore(ctx, userID, ScoreIdeaUpvote, "Idea received upvote", "idea", &ideaID)
}

func (s *service) ProcessIdeaDownvote(ctx context.Context, userID, ideaID string) error {
	return s.scoreRepo.UpdateScore(ctx, userID, ScoreIdeaDownvote, "Idea received downvote", "idea", &ideaID)
}

func (s *service) ProcessPositiveReview(ctx context.Context, userID, reviewID string) error {
	return s.scoreRepo.UpdateScore(ctx, userID, ScorePositiveReview, "Received positive review", "review", &reviewID)
}

func (s *service) ProcessNegativeReview(ctx context.Context, userID, reviewID string) error {
	return s.scoreRepo.UpdateScore(ctx, userID, ScoreNegativeReview, "Received negative review", "review", &reviewID)
}

func (s *service) ProcessBookDonation(ctx context.Context, userID, donationID string) error {
	return s.scoreRepo.UpdateScore(ctx, userID, ScoreBookDonated, "Donated book", "donation", &donationID)
}

func (s *service) ProcessMoneyDonation(ctx context.Context, userID, donationID string) error {
	return s.scoreRepo.UpdateScore(ctx, userID, ScoreMoneyDonated, "Made financial contribution", "donation", &donationID)
}

func (s *service) ProcessReturnOnTime(ctx context.Context, userID, bookID string) error {
	return s.scoreRepo.UpdateScore(ctx, userID, ScoreReturnOnTime, "Returned book on time", "book", &bookID)
}

func (s *service) ProcessReturnLate(ctx context.Context, userID, bookID string) error {
	return s.scoreRepo.UpdateScore(ctx, userID, ScoreReturnLate, "Returned book late", "book", &bookID)
}

func (s *service) ProcessLostBook(ctx context.Context, userID, bookID string) error {
	return s.scoreRepo.UpdateScore(ctx, userID, ScoreLostBook, "Lost book", "book", &bookID)
}
