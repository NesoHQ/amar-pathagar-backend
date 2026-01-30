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
	handoverRepo    HandoverRepo
	log             *zap.Logger
}

func NewService(adminRepo AdminRepo, successScoreSvc successscore.Service, notificationSvc notification.Service, handoverRepo HandoverRepo, log *zap.Logger) Service {
	return &service{
		adminRepo:       adminRepo,
		successScoreSvc: successScoreSvc,
		notificationSvc: notificationSvc,
		handoverRepo:    handoverRepo,
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

	// Update request status to approved
	if err := s.adminRepo.UpdateRequestStatus(ctx, requestID, "approved", processedAt, &dueDate); err != nil {
		s.log.Error("failed to approve request", zap.String("request_id", requestID), zap.Error(err))
		return err
	}

	// Reject all other pending requests for this book
	otherRequests, err := s.adminRepo.GetRequestsByBook(ctx, targetRequest.BookID)
	if err != nil {
		s.log.Error("failed to get other requests", zap.Error(err))
	} else {
		for _, req := range otherRequests {
			if req.ID != requestID && req.Status == "pending" {
				if err := s.adminRepo.UpdateRequestStatus(ctx, req.ID, "rejected", processedAt, nil); err != nil {
					s.log.Error("failed to reject other request", zap.String("request_id", req.ID), zap.Error(err))
				}
			}
		}
	}

	// Get book details
	book, err := s.adminRepo.GetBookByID(ctx, targetRequest.BookID)
	if err != nil {
		s.log.Error("failed to get book details", zap.Error(err))
		return err
	}

	// Update book status to "requested"
	if err := s.adminRepo.UpdateBookStatus(ctx, targetRequest.BookID, "requested"); err != nil {
		s.log.Error("failed to update book status", zap.Error(err))
		return err
	}

	// Determine who is the current holder
	var currentHolderID string

	if book.Status == domain.StatusAvailable || book.Status == domain.StatusOnHold {
		// Book is available or on_hold - check if someone read it before
		lastHistory, err := s.handoverRepo.GetLastCompletedReadingHistory(ctx, targetRequest.BookID)
		if err != nil || lastHistory == nil {
			// No one has read it yet, use book creator
			if book.CreatedBy != nil && *book.CreatedBy != "" {
				currentHolderID = *book.CreatedBy
			} else {
				return fmt.Errorf("cannot determine current holder")
			}
		} else {
			// Someone read it before, they still have the physical book
			currentHolderID = lastHistory.ReaderID
		}
	} else {
		return fmt.Errorf("book is not available for request (status: %s)", book.Status)
	}

	// Create handover thread between current holder and new requester
	parsedDueDate, parseErr := time.Parse(time.RFC3339, dueDate)
	if parseErr != nil {
		s.log.Error("failed to parse due date", zap.Error(parseErr))
		parsedDueDate = time.Now().AddDate(0, 0, 14) // Default 14 days
	}

	thread := &domain.HandoverThread{
		BookID:          targetRequest.BookID,
		CurrentHolderID: currentHolderID,
		NextHolderID:    targetRequest.UserID,
		Status:          string(domain.HandoverActive),
		HandoverDueDate: parsedDueDate,
		IsPublic:        true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.handoverRepo.CreateHandoverThread(ctx, thread); err != nil {
		s.log.Error("failed to create handover thread", zap.Error(err))
	} else {
		// Create system message
		systemMsg := &domain.HandoverMessage{
			ThreadID:        thread.ID,
			UserID:          currentHolderID,
			Message:         fmt.Sprintf("📚 Book handover thread created. Please coordinate delivery of \"%s\" to the reader. Due date: %s", book.Title, parsedDueDate.Format("Jan 2, 2006")),
			IsSystemMessage: true,
			CreatedAt:       time.Now(),
		}
		if err := s.handoverRepo.CreateHandoverMessage(ctx, systemMsg); err != nil {
			s.log.Error("failed to create system message", zap.Error(err))
		}

		// Notify both users
		if err := s.notificationSvc.NotifyHandoverThreadCreated(ctx, currentHolderID, targetRequest.UserID, targetRequest.BookID, book.Title); err != nil {
			s.log.Error("failed to send notifications", zap.Error(err))
		}
	}

	// Increment user's books_received counter
	if err := s.adminRepo.IncrementUserBooksReceived(ctx, targetRequest.UserID); err != nil {
		s.log.Error("failed to increment books received", zap.Error(err))
	}

	// Send notification to user
	if err := s.notificationSvc.NotifyRequestApproved(ctx, targetRequest.UserID, targetRequest.BookID, targetRequest.Book.Title); err != nil {
		s.log.Error("failed to send notification", zap.Error(err))
	}

	s.log.Info("book request approved", zap.String("request_id", requestID), zap.String("book_status", string(book.Status)))
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
