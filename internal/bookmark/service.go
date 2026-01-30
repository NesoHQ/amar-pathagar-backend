package bookmark

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/online-library/internal/domain"
	"go.uber.org/zap"
)

type service struct {
	bookmarkRepo BookmarkRepo
	log          *zap.Logger
}

func NewService(bookmarkRepo BookmarkRepo, log *zap.Logger) Service {
	return &service{
		bookmarkRepo: bookmarkRepo,
		log:          log,
	}
}

func (s *service) Create(ctx context.Context, bookmark *domain.UserBookmark) (*domain.UserBookmark, error) {
	bookmark.ID = uuid.New().String()
	bookmark.CreatedAt = time.Now()

	if err := s.bookmarkRepo.Create(ctx, bookmark); err != nil {
		s.log.Error("failed to create bookmark", zap.Error(err))
		return nil, err
	}

	s.log.Info("bookmark created successfully", zap.String("bookmark_id", bookmark.ID))
	return bookmark, nil
}

func (s *service) Delete(ctx context.Context, userID, bookID, bookmarkType string) error {
	if err := s.bookmarkRepo.Delete(ctx, userID, bookID, bookmarkType); err != nil {
		s.log.Error("failed to delete bookmark", zap.String("user_id", userID), zap.String("book_id", bookID), zap.String("bookmark_type", bookmarkType), zap.String("error", err.Error()))
		return err
	}

	s.log.Info("bookmark deleted successfully", zap.String("user_id", userID), zap.String("book_id", bookID), zap.String("bookmark_type", bookmarkType))
	return nil
}

func (s *service) GetByUser(ctx context.Context, userID string) ([]*domain.UserBookmark, error) {
	bookmarks, err := s.bookmarkRepo.FindByUserID(ctx, userID)
	if err != nil {
		s.log.Error("failed to get bookmarks", zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}
	return bookmarks, nil
}
