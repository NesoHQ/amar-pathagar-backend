package idea

import (
	"context"

	"github.com/yourusername/online-library/internal/domain"
)

type Service interface {
	Create(ctx context.Context, idea *domain.ReadingIdea) (*domain.ReadingIdea, error)
	GetByBook(ctx context.Context, bookID string) ([]*domain.ReadingIdea, error)
	Vote(ctx context.Context, ideaID, userID string, voteType domain.VoteType) error
}

type IdeaRepo interface {
	Create(ctx context.Context, idea *domain.ReadingIdea) error
	FindByBookID(ctx context.Context, bookID string) ([]*domain.ReadingIdea, error)
	AddVote(ctx context.Context, vote *domain.IdeaVote) error
	RemoveVote(ctx context.Context, ideaID, userID string) error
	UpdateVoteCounts(ctx context.Context, ideaID string, upvotes, downvotes int) error
}

type SuccessScoreSvc interface {
	ProcessIdeaPosted(ctx context.Context, userID, ideaID string) error
	ProcessIdeaUpvote(ctx context.Context, userID, ideaID string) error
	ProcessIdeaDownvote(ctx context.Context, userID, ideaID string) error
}

type NotificationSvc interface {
	NotifyIdeaVote(ctx context.Context, userID, voterName, ideaTitle string, isUpvote bool) error
}
