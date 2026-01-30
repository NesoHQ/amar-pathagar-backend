package admin

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/online-library/internal/domain"
	"github.com/yourusername/online-library/internal/notification"
	"github.com/yourusername/online-library/internal/successscore"
	"go.uber.org/zap"
)

type service struct {
	adminRepo       AdminRepo
	successScoreSvc successscore.Service
	notificationSvc notification.Service
	log             *zap.Logger
}

func NewService(adminRepo AdminRepo, successScoreSvc successscore.Service, notificationSvc notification.Service, log *zap.Logger) Service {
	return &service{
		adminRepo:       adminRepo,
		successScoreSvc: successScoreSvc,
		notificationSvc: notificationSvc,
		log:             log,
	}
}

func (s *service) GetPendingRequests(ctx context.Context, limit, offset int) ([]*domain.BookRequest, error) {
	return s.adminRepo.GetPendingRequests(ctx, limit, offset)
}

func (s *service) ApproveBookRequest(ctx context.Context, requestID string, dueDate string) error {
	processedAt := time.Now().Format(time.RFC3339)

	// Get the request details first
	requests, err := s.adminRepo.GetPendingRequests(ctx, 1000, 0)
	if err != nil {
		return err
	}

	var targetRequest *domain.BookRequest
	for _, req := range requests {
		if req.ID == requestID {
			targetRequest = req
			break
		}
	}

	if targetRequest == nil {
		return fmt.Errorf("request not found")
	}

	// Update request status
	if err := s.adminRepo.UpdateRequestStatus(ctx, requestID, "approved", processedAt, &dueDate); err != nil {
		s.log.Error("failed to approve request", zap.String("request_id", requestID), zap.Error(err))
		return err
	}

	// Update book status and assign to user
	if err := s.adminRepo.AssignBookToUser(ctx, targetRequest.BookID, targetRequest.UserID); err != nil {
		s.log.Error("failed to assign book to user", zap.Error(err))
		return err
	}

	// Create reading history
	if err := s.adminRepo.CreateReadingHistory(ctx, targetRequest.BookID, targetRequest.UserID, dueDate); err != nil {
		s.log.Error("failed to create reading history", zap.Error(err))
		return err
	}

	// Increment user's books_received counter
	if err := s.adminRepo.IncrementUserBooksReceived(ctx, targetRequest.UserID); err != nil {
		s.log.Error("failed to increment books received", zap.Error(err))
	}

	// Send notification to user
	if err := s.notificationSvc.NotifyRequestApproved(ctx, targetRequest.UserID, targetRequest.BookID, targetRequest.Book.Title); err != nil {
		s.log.Error("failed to send notification", zap.Error(err))
	}

	s.log.Info("book request approved", zap.String("request_id", requestID))
	return nil
}

func (s *service) RejectBookRequest(ctx context.Context, requestID string, reason string) error {
	processedAt := time.Now().Format(time.RFC3339)

	if err := s.adminRepo.UpdateRequestStatus(ctx, requestID, "rejected", processedAt, nil); err != nil {
		s.log.Error("failed to reject request", zap.String("request_id", requestID), zap.Error(err))
		return err
	}

	s.log.Info("book request rejected", zap.String("request_id", requestID), zap.String("reason", reason))
	return nil
}

func (s *service) GetRequestsByBook(ctx context.Context, bookID string) ([]*domain.BookRequest, error) {
	return s.adminRepo.GetRequestsByBook(ctx, bookID)
}

func (s *service) GetAllUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	return s.adminRepo.GetAllUsers(ctx, limit, offset)
}

func (s *service) AdjustSuccessScore(ctx context.Context, userID string, amount int, reason string) error {
	if err := s.successScoreSvc.AdjustScore(ctx, userID, amount, reason, "", nil); err != nil {
		s.log.Error("failed to adjust success score", zap.String("user_id", userID), zap.Error(err))
		return err
	}

	s.log.Info("success score adjusted by admin", zap.String("user_id", userID), zap.Int("amount", amount))
	return nil
}

func (s *service) UpdateUserRole(ctx context.Context, userID string, role domain.UserRole) error {
	if role != domain.RoleAdmin && role != domain.RoleMember {
		return fmt.Errorf("invalid role: %s", role)
	}

	if err := s.adminRepo.UpdateUserRole(ctx, userID, string(role)); err != nil {
		s.log.Error("failed to update user role", zap.String("user_id", userID), zap.Error(err))
		return err
	}

	s.log.Info("user role updated", zap.String("user_id", userID), zap.String("role", string(role)))
	return nil
}

func (s *service) GetSystemStats(ctx context.Context) (*SystemStats, error) {
	return s.adminRepo.GetSystemStats(ctx)
}

func (s *service) GetAuditLogs(ctx context.Context, limit, offset int) ([]*AuditLog, error) {
	return s.adminRepo.GetAuditLogs(ctx, limit, offset)
}

func (s *service) GetAllBooks(ctx context.Context, limit, offset int, filters BookFilters) ([]*domain.Book, error) {
	return s.adminRepo.GetAllBooks(ctx, limit, offset, filters)
}

func (s *service) UpdateBookStatus(ctx context.Context, bookID string, status domain.BookStatus) error {
	return s.adminRepo.UpdateBookStatus(ctx, bookID, string(status))
}
