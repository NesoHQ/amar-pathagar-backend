package book

import (
	"context"
	"time"

	"github.com/yourusername/online-library/internal/domain"
	"go.uber.org/zap"
)

func (s *service) Update(ctx context.Context, id string, book *domain.Book) (*domain.Book, error) {
	// Check if book exists
	existing, err := s.bookRepo.FindByID(ctx, id)
	if err != nil {
		s.log.Error("book not found", zap.String("book_id", id), zap.Error(err))
		return nil, domain.ErrNotFound
	}

	// Update fields
	book.ID = existing.ID
	book.CreatedAt = existing.CreatedAt
	book.UpdatedAt = time.Now()

	if err := s.bookRepo.Update(ctx, id, book); err != nil {
		s.log.Error("failed to update book", zap.String("book_id", id), zap.Error(err))
		return nil, err
	}

	s.log.Info("book updated successfully", zap.String("book_id", id))
	return book, nil
}
