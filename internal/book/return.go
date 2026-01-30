package book

import (
	"context"

	"github.com/yourusername/online-library/internal/domain"
	"go.uber.org/zap"
)

func (s *service) ReturnBook(ctx context.Context, bookID, userID string) error {
	// Get book details
	book, err := s.bookRepo.FindByID(ctx, bookID)
	if err != nil {
		s.log.Error("book not found", zap.String("book_id", bookID), zap.Error(err))
		return domain.ErrNotFound
	}

	// Verify user is the current holder
	if book.CurrentHolderID == nil || *book.CurrentHolderID != userID {
		return domain.ErrUnauthorized
	}

	// Update book status
	if err := s.bookRepo.ReturnBook(ctx, bookID); err != nil {
		s.log.Error("failed to return book", zap.String("book_id", bookID), zap.Error(err))
		return err
	}

	// Update reading history
	if err := s.bookRepo.CompleteReadingHistory(ctx, bookID, userID); err != nil {
		s.log.Error("failed to complete reading history", zap.Error(err))
	}

	s.log.Info("book returned successfully", zap.String("book_id", bookID), zap.String("user_id", userID))
	return nil
}

func (s *service) GetReadingHistory(ctx context.Context, userID string) ([]*domain.ReadingHistory, error) {
	history, err := s.bookRepo.GetReadingHistoryByUser(ctx, userID)
	if err != nil {
		s.log.Error("failed to get reading history", zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}
	return history, nil
}

func (s *service) GetBooksOnHold(ctx context.Context, userID string) ([]*domain.Book, error) {
	books, err := s.bookRepo.GetBooksOnHoldByUser(ctx, userID)
	if err != nil {
		s.log.Error("failed to get books on hold", zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}
	return books, nil
}
