package book

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/online-library/internal/domain"
	"go.uber.org/zap"
)

func (s *service) Create(ctx context.Context, book *domain.Book) (*domain.Book, error) {
	book.ID = uuid.New().String()
	book.Status = domain.StatusAvailable
	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()

	if err := s.bookRepo.Create(ctx, book); err != nil {
		s.log.Error("failed to create book", zap.Error(err))
		return nil, err
	}

	s.log.Info("book created successfully", zap.String("book_id", book.ID))
	return book, nil
}
