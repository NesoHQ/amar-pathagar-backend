package book

import (
	"context"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/online-library/internal/domain"
	"go.uber.org/zap"
)

func (s *service) RequestBook(ctx context.Context, bookID, userID string) (*domain.BookRequest, error) {
	// Check if book exists and is available
	book, err := s.bookRepo.FindByID(ctx, bookID)
	if err != nil {
		s.log.Error("book not found", zap.String("book_id", bookID), zap.Error(err))
		return nil, domain.ErrNotFound
	}

	// Check if user already has a pending request for this book
	existing, err := s.bookRepo.FindRequestByBookAndUser(ctx, bookID, userID)
	if err != nil {
		s.log.Error("failed to check existing request", zap.Error(err))
		return nil, err
	}
	if existing != nil {
		return nil, domain.ErrAlreadyRequested
	}

	// Calculate priority score (success score + interest match + distance)
	priorityScore := calculatePriorityScore(book, userID)

	request := &domain.BookRequest{
		ID:                 uuid.New().String(),
		BookID:             bookID,
		UserID:             userID,
		Status:             "pending",
		PriorityScore:      priorityScore,
		InterestMatchScore: 0.0, // TODO: Calculate based on user interests
		RequestedAt:        time.Now(),
	}

	if err := s.bookRepo.CreateRequest(ctx, request); err != nil {
		s.log.Error("failed to create book request", zap.Error(err))
		return nil, err
	}

	s.log.Info("book request created", zap.String("request_id", request.ID), zap.String("book_id", bookID), zap.String("user_id", userID))
	return request, nil
}

func (s *service) GetUserRequests(ctx context.Context, userID string) ([]*domain.BookRequest, error) {
	requests, err := s.bookRepo.FindRequestsByUserID(ctx, userID)
	if err != nil {
		s.log.Error("failed to get user requests", zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}
	return requests, nil
}

func (s *service) CheckBookRequested(ctx context.Context, bookID, userID string) (bool, error) {
	request, err := s.bookRepo.FindRequestByBookAndUser(ctx, bookID, userID)
	if err != nil {
		return false, err
	}
	return request != nil, nil
}

func (s *service) CancelRequest(ctx context.Context, bookID, userID string) error {
	if err := s.bookRepo.CancelRequest(ctx, bookID, userID); err != nil {
		s.log.Error("failed to cancel request", zap.String("book_id", bookID), zap.String("user_id", userID), zap.Error(err))
		return err
	}
	s.log.Info("book request cancelled", zap.String("book_id", bookID), zap.String("user_id", userID))
	return nil
}

// calculatePriorityScore calculates the priority score for a book request
// Based on: success score (weight: 0.5) + interest match (weight: 0.3) + distance (weight: 0.2)
func calculatePriorityScore(book *domain.Book, userID string) float64 {
	// Base score from success score (normalized to 0-100)
	// This will be enhanced with actual user data
	baseScore := 50.0 // Default middle score

	// Interest match score (0-100)
	interestScore := 0.0

	// Distance score (0-100, closer is better)
	distanceScore := 50.0 // Default middle score

	// Weighted calculation
	priorityScore := (baseScore * 0.5) + (interestScore * 0.3) + (distanceScore * 0.2)

	return priorityScore
}

// calculateDistance calculates distance between two coordinates in kilometers
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371 // Earth's radius in kilometers

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}
