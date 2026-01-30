package bookmark

import (
	"context"

	"github.com/yourusername/online-library/internal/domain"
)

type Service interface {
	Create(ctx context.Context, bookmark *domain.UserBookmark) (*domain.UserBookmark, error)
	Delete(ctx context.Context, userID, bookID, bookmarkType string) error
	GetByUser(ctx context.Context, userID string) ([]*domain.UserBookmark, error)
}

type BookmarkRepo interface {
	Create(ctx context.Context, bookmark *domain.UserBookmark) error
	Delete(ctx context.Context, userID, bookID, bookmarkType string) error
	FindByUserID(ctx context.Context, userID string) ([]*domain.UserBookmark, error)
}
