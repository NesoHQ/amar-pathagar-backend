package book

import (
	"context"
	"strings"
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

	// Check if user already has a pending request for this book
	// For now, we'll just create the request and let the unique constraint handle it
	// In production, you'd want to check first and return a friendly error

	// Create book request
	request := &domain.BookRequest{
		ID:          uuid.New().String(),
		BookID:      bookID,
		UserID:      userID,
		Status:      "pending",
		RequestedAt: time.Now(),
	}

	if err := s.bookRepo.CreateRequest(ctx, request); err != nil {
		// Check if it's a duplicate key error
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			s.log.Warn("user already has a pending request for this book", zap.String("book_id", bookID), zap.String("user_id", userID))
			return nil, domain.ErrBookAlreadyRequested
		}
		s.log.Error("failed to create book request", zap.String("error", err.Error()))
		return nil, err
	}

	s.log.Info("book request created", zap.String("request_id", request.ID), zap.String("book_id", bookID))
	return request, nil
}

func (s *service) GetUserRequests(ctx context.Context, userID string) ([]*domain.BookRequest, error) {
	requests, err := s.bookRepo.FindRequestsByUserID(ctx, userID)
	if err != nil {
		s.log.Error("failed to get user requests", zap.String("user_id", userID), zap.String("error", err.Error()))
		return nil, err
	}
	return requests, nil
}

func (s *service) CheckBookRequested(ctx context.Context, bookID, userID string) (bool, error) {
	request, err := s.bookRepo.FindRequestByBookAndUser(ctx, bookID, userID)
	if err != nil {
		s.log.Error("failed to check book request", zap.String("book_id", bookID), zap.String("user_id", userID), zap.String("error", err.Error()))
		return false, err
	}
	return request != nil, nil
}
