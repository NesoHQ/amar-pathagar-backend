package book

import (
	"context"

	"github.com/yourusername/online-library/internal/domain"
	"go.uber.org/zap"
)

func (s *service) GetByID(ctx context.Context, id string) (*domain.Book, error) {
	book, err := s.bookRepo.FindByID(ctx, id)
	if err != nil {
		s.log.Error("failed to get book", zap.String("book_id", id), zap.Error(err))
		return nil, err
	}

	return book, nil
}

func (s *service) List(ctx context.Context, limit, offset int) ([]*domain.Book, error) {
	books, err := s.bookRepo.List(ctx, limit, offset)
	if err != nil {
		s.log.Error("failed to list books", zap.Error(err))
		return nil, err
	}

	return books, nil
}
