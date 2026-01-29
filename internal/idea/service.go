package idea

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/online-library/internal/domain"
	"go.uber.org/zap"
)

type service struct {
	ideaRepo        IdeaRepo
	successScoreSvc SuccessScoreSvc
	notificationSvc NotificationSvc
	log             *zap.Logger
}

func NewService(ideaRepo IdeaRepo, successScoreSvc SuccessScoreSvc, notificationSvc NotificationSvc, log *zap.Logger) Service {
	return &service{
		ideaRepo:        ideaRepo,
		successScoreSvc: successScoreSvc,
		notificationSvc: notificationSvc,
		log:             log,
	}
}

func (s *service) Create(ctx context.Context, idea *domain.ReadingIdea) (*domain.ReadingIdea, error) {
	idea.ID = uuid.New().String()
	idea.CreatedAt = time.Now()
	idea.UpdatedAt = time.Now()

	if err := s.ideaRepo.Create(ctx, idea); err != nil {
		s.log.Error("failed to create idea", zap.Error(err))
		return nil, err
	}

	// Update success score
	if err := s.successScoreSvc.ProcessIdeaPosted(ctx, idea.UserID, idea.ID); err != nil {
		s.log.Warn("failed to update success score for idea", zap.Error(err))
	}

	s.log.Info("idea created successfully", zap.String("idea_id", idea.ID))
	return idea, nil
}

func (s *service) GetByBook(ctx context.Context, bookID string) ([]*domain.ReadingIdea, error) {
	ideas, err := s.ideaRepo.FindByBookID(ctx, bookID)
	if err != nil {
		s.log.Error("failed to get ideas", zap.String("book_id", bookID), zap.Error(err))
		return nil, err
	}
	return ideas, nil
}

func (s *service) Vote(ctx context.Context, ideaID, userID string, voteType domain.VoteType) error {
	vote := &domain.IdeaVote{
		ID:        uuid.New().String(),
		IdeaID:    ideaID,
		UserID:    userID,
		VoteType:  voteType,
		CreatedAt: time.Now(),
	}

	if err := s.ideaRepo.AddVote(ctx, vote); err != nil {
		s.log.Error("failed to add vote", zap.Error(err))
		return err
	}

	// Update success score
	if voteType == domain.VoteTypeUp {
		if err := s.successScoreSvc.ProcessIdeaUpvote(ctx, userID, ideaID); err != nil {
			s.log.Warn("failed to update success score for upvote", zap.Error(err))
		}
	} else {
		if err := s.successScoreSvc.ProcessIdeaDownvote(ctx, userID, ideaID); err != nil {
			s.log.Warn("failed to update success score for downvote", zap.Error(err))
		}
	}

	s.log.Info("vote added successfully", zap.String("idea_id", ideaID))
	return nil
}
