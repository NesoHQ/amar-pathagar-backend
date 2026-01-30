package book

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/online-library/internal/domain"
	"go.uber.org/zap"
)

func (s *service) RequestBook(ctx context.Context, bookID, userID string) (*domain.BookRequest, error) {
	// Check if book exists and is available
	book, err := s.bookRepo.FindByID(ctx, bookID)
	if err != nil {
		s.log.Error("failed to find book", zap.String("book_id", bookID), zap.String("error", err.Error()))
		return nil, domain.ErrBookNotFound
	}

	if book.Status != "available" {
		s.log.Warn("book not available for request", zap.String("book_id", bookID), zap.String("status", string(book.Status)))
		return nil, domain.ErrBookNotAvailable
	}

	// Create book request
	request := &domain.BookRequest{
		ID:          uuid.New().String(),
		BookID:      bookID,
		UserID:      userID,
		Status:      "pending",
		RequestedAt: time.Now(),
	}

	if err := s.bookRepo.CreateRequest(ctx, request); err != nil {
		s.log.Error("failed to create book request", zap.String("error", err.Error()))
		return nil, err
	}

	s.log.Info("book request created", zap.String("request_id", request.ID), zap.String("book_id", bookID))
	return request, nil
}
