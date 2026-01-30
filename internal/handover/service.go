package handover

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/online-library/internal/domain"
	"github.com/yourusername/online-library/internal/notification"
	"go.uber.org/zap"
)

type service struct {
	handoverRepo    HandoverRepo
	notificationSvc notification.Service
	log             *zap.Logger
}

func NewService(handoverRepo HandoverRepo, notificationSvc notification.Service, log *zap.Logger) Service {
	return &service{
		handoverRepo:    handoverRepo,
		notificationSvc: notificationSvc,
		log:             log,
	}
}

func (s *service) MarkBookCompleted(ctx context.Context, userID, bookID string) error {
	// Get active reading history
	history, err := s.handoverRepo.GetActiveReadingHistory(ctx, bookID)
	if err != nil {
		return fmt.Errorf("failed to get reading history: %w", err)
	}

	if history == nil {
		// Check if there's already a completed reading history for this user
		lastHistory, err := s.handoverRepo.GetLastCompletedReadingHistory(ctx, bookID)
		if err == nil && lastHistory != nil && lastHistory.ReaderID == userID {
			// User already completed this book
			return fmt.Errorf("you have already marked this book as completed")
		}

		// No active reading history found - this can happen for old books
		// Create a reading history entry for this book
		s.log.Warn("no active reading history found, creating one", zap.String("book_id", bookID), zap.String("user_id", userID))

		// Create reading history starting from now
		if err := s.handoverRepo.StartNewReadingHistory(ctx, bookID, userID); err != nil {
			return fmt.Errorf("failed to create reading history: %w", err)
		}

		// Get the newly created history
		history, err = s.handoverRepo.GetActiveReadingHistory(ctx, bookID)
		if err != nil || history == nil {
			return fmt.Errorf("failed to get newly created reading history: %w", err)
		}
	}

	if history.ReaderID != userID {
		return fmt.Errorf("you are not the current holder of this book")
	}

	if history.IsCompleted {
		return fmt.Errorf("book already marked as completed")
	}

	// Mark as completed
	completedAt := time.Now()
	if err := s.handoverRepo.UpdateReadingHistoryCompleted(ctx, history.ID, completedAt); err != nil {
		return fmt.Errorf("failed to mark book as completed: %w", err)
	}

	// Check if there's a next reader
	hasNextReader := history.NextReaderID != nil && *history.NextReaderID != ""

	if hasNextReader {
		// Update delivery status to in_transit
		if err := s.handoverRepo.UpdateReadingHistoryDeliveryStatus(ctx, history.ID, domain.DeliveryInTransit, nil); err != nil {
			s.log.Error("failed to update delivery status", zap.Error(err))
		}

		// Notify next reader
		if err := s.notificationSvc.NotifyBookInTransit(ctx, *history.NextReaderID, bookID, history.Book.Title); err != nil {
			s.log.Error("failed to send notification", zap.Error(err))
		}
	} else {
		// No next reader, close reading history and mark book as on_hold
		// Book stays with the reader (no penalty) until someone requests it
		if err := s.handoverRepo.CloseReadingHistory(ctx, history.ID, completedAt); err != nil {
			s.log.Error("failed to close reading history", zap.Error(err))
			return fmt.Errorf("failed to close reading history: %w", err)
		}

		if err := s.handoverRepo.UpdateBookStatus(ctx, bookID, domain.StatusOnHold); err != nil {
			s.log.Error("failed to update book status to on_hold", zap.String("book_id", bookID), zap.Error(err))
			return fmt.Errorf("failed to update book status: %w", err)
		}

		s.log.Info("book status updated to on_hold", zap.String("book_id", bookID))

		// Keep the book assigned to the current reader
		// They will hold it until someone requests it

		// Close any active handover thread since there's no next reader
		activeThread, err := s.handoverRepo.GetActiveHandoverThreadByBook(ctx, bookID)
		if err == nil && activeThread != nil {
			if err := s.handoverRepo.UpdateHandoverThreadStatus(ctx, activeThread.ID, domain.HandoverCompleted, &completedAt); err != nil {
				s.log.Error("failed to complete handover thread", zap.Error(err))
			} else {
				// Post system message
				systemMsg := &domain.HandoverMessage{
					ThreadID:        activeThread.ID,
					UserID:          userID,
					Message:         "📚 Book reading completed. Book is on hold until next request.",
					IsSystemMessage: true,
					CreatedAt:       time.Now(),
				}
				if err := s.handoverRepo.CreateHandoverMessage(ctx, systemMsg); err != nil {
					s.log.Error("failed to create system message", zap.Error(err))
				}
			}
		}
	}

	s.log.Info("book marked as completed",
		zap.String("book_id", bookID),
		zap.String("user_id", userID),
		zap.Bool("has_next_reader", hasNextReader))
	return nil
}

func (s *service) MarkBookDelivered(ctx context.Context, userID, bookID string) error {
	// Check if there's an active handover thread where user is the next holder
	thread, err := s.handoverRepo.GetActiveHandoverThreadByBook(ctx, bookID)
	if err != nil {
		return fmt.Errorf("failed to get handover thread: %w", err)
	}

	if thread == nil {
		return fmt.Errorf("no active handover thread found")
	}

	// Check if user is the next holder in the thread
	if thread.NextHolderID != userID {
		return fmt.Errorf("you are not the next holder for this book")
	}

	// Get active reading history (may not exist for initial handover)
	history, err := s.handoverRepo.GetActiveReadingHistory(ctx, bookID)

	if history == nil {
		// Initial handover case: No reading history exists yet
		// Create new reading history for the first reader
		if err := s.handoverRepo.StartNewReadingHistory(ctx, bookID, userID); err != nil {
			return fmt.Errorf("failed to start reading history: %w", err)
		}

		// Update book status to reading and assign to user
		if err := s.handoverRepo.UpdateBookStatus(ctx, bookID, domain.StatusReading); err != nil {
			s.log.Error("failed to update book status", zap.Error(err))
		}

		// Assign book to user
		if err := s.handoverRepo.AssignBookToUser(ctx, bookID, userID); err != nil {
			s.log.Error("failed to assign book to user", zap.Error(err))
		}
	} else {
		// Reader-to-reader handover case
		if history.DeliveryStatus == string(domain.DeliveryDelivered) {
			return fmt.Errorf("book already marked as delivered")
		}

		// Mark as delivered
		deliveredAt := time.Now()
		if err := s.handoverRepo.UpdateReadingHistoryDeliveryStatus(ctx, history.ID, domain.DeliveryDelivered, &deliveredAt); err != nil {
			return fmt.Errorf("failed to mark book as delivered: %w", err)
		}

		// Close the old reading history
		if err := s.handoverRepo.CloseReadingHistory(ctx, history.ID, deliveredAt); err != nil {
			s.log.Error("failed to close reading history", zap.Error(err))
		}

		// Start new reading history for the next reader
		if err := s.handoverRepo.StartNewReadingHistory(ctx, bookID, userID); err != nil {
			s.log.Error("failed to start new reading history", zap.Error(err))
		}

		// Notify previous holder
		if err := s.notificationSvc.NotifyBookDelivered(ctx, history.ReaderID, bookID, history.Book.Title); err != nil {
			s.log.Error("failed to send notification", zap.Error(err))
		}
	}

	// Complete the handover thread
	completedAt := time.Now()
	if err := s.handoverRepo.UpdateHandoverThreadStatus(ctx, thread.ID, domain.HandoverCompleted, &completedAt); err != nil {
		s.log.Error("failed to complete handover thread", zap.Error(err))
	}

	// Post system message
	systemMsg := &domain.HandoverMessage{
		ThreadID:        thread.ID,
		UserID:          userID,
		Message:         "📦 Book has been delivered and received successfully!",
		IsSystemMessage: true,
		CreatedAt:       time.Now(),
	}
	if err := s.handoverRepo.CreateHandoverMessage(ctx, systemMsg); err != nil {
		s.log.Error("failed to create system message", zap.Error(err))
	}

	s.log.Info("book marked as delivered", zap.String("book_id", bookID), zap.String("user_id", userID))
	return nil
}

func (s *service) GetActiveHandoverThread(ctx context.Context, bookID string) (*domain.HandoverThread, error) {
	return s.handoverRepo.GetActiveHandoverThreadByBook(ctx, bookID)
}

func (s *service) GetUserHandoverThreads(ctx context.Context, userID string) ([]*domain.HandoverThread, error) {
	return s.handoverRepo.GetHandoverThreadsByUser(ctx, userID)
}

func (s *service) PostHandoverMessage(ctx context.Context, threadID, userID, message string) error {
	// Verify thread exists and user is participant
	thread, err := s.handoverRepo.GetHandoverThreadByID(ctx, threadID)
	if err != nil {
		return fmt.Errorf("failed to get thread: %w", err)
	}

	if thread == nil {
		return fmt.Errorf("thread not found")
	}

	if thread.CurrentHolderID != userID && thread.NextHolderID != userID {
		return fmt.Errorf("you are not a participant in this handover")
	}

	msg := &domain.HandoverMessage{
		ThreadID:        threadID,
		UserID:          userID,
		Message:         message,
		IsSystemMessage: false,
		CreatedAt:       time.Now(),
	}

	if err := s.handoverRepo.CreateHandoverMessage(ctx, msg); err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	// Notify the other participant
	otherUserID := thread.NextHolderID
	if userID == thread.NextHolderID {
		otherUserID = thread.CurrentHolderID
	}

	if err := s.notificationSvc.NotifyHandoverMessage(ctx, otherUserID, thread.BookID, thread.Book.Title); err != nil {
		s.log.Error("failed to send notification", zap.Error(err))
	}

	return nil
}

func (s *service) GetHandoverMessages(ctx context.Context, threadID string) ([]domain.HandoverMessage, error) {
	return s.handoverRepo.GetHandoverMessagesByThread(ctx, threadID)
}

func (s *service) CheckAndCreateHandoverThreads(ctx context.Context) error {
	// Get reading histories due within 7 days
	histories, err := s.handoverRepo.GetReadingHistoriesDueSoon(ctx, 7)
	if err != nil {
		return fmt.Errorf("failed to get due histories: %w", err)
	}

	for _, history := range histories {
		// Skip if already has next reader assigned
		if history.NextReaderID != nil && *history.NextReaderID != "" {
			continue
		}

		// Get next approved request
		nextRequest, err := s.handoverRepo.GetNextApprovedRequest(ctx, history.BookID)
		if err != nil {
			s.log.Error("failed to get next request", zap.String("book_id", history.BookID), zap.Error(err))
			continue
		}

		if nextRequest == nil {
			// No next reader, skip
			continue
		}

		// Update reading history with next reader
		if err := s.handoverRepo.UpdateReadingHistoryNextReader(ctx, history.ID, nextRequest.UserID); err != nil {
			s.log.Error("failed to update next reader", zap.Error(err))
			continue
		}

		// Check if handover thread already exists
		existingThread, _ := s.handoverRepo.GetActiveHandoverThreadByBook(ctx, history.BookID)
		if existingThread != nil {
			continue
		}

		// Create handover thread
		thread := &domain.HandoverThread{
			BookID:           history.BookID,
			CurrentHolderID:  history.ReaderID,
			NextHolderID:     nextRequest.UserID,
			ReadingHistoryID: &history.ID,
			Status:           string(domain.HandoverActive),
			HandoverDueDate:  *history.DueDate,
			IsPublic:         true,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		if err := s.handoverRepo.CreateHandoverThread(ctx, thread); err != nil {
			s.log.Error("failed to create handover thread", zap.Error(err))
			continue
		}

		// Create system message
		systemMsg := &domain.HandoverMessage{
			ThreadID:        thread.ID,
			UserID:          history.ReaderID,
			Message:         fmt.Sprintf("📚 Handover thread created. Book is due on %s. Please coordinate the handover.", history.DueDate.Format("Jan 2, 2006")),
			IsSystemMessage: true,
			CreatedAt:       time.Now(),
		}

		if err := s.handoverRepo.CreateHandoverMessage(ctx, systemMsg); err != nil {
			s.log.Error("failed to create system message", zap.Error(err))
		}

		// Notify both users
		if err := s.notificationSvc.NotifyHandoverThreadCreated(ctx, history.ReaderID, nextRequest.UserID, history.BookID, history.Book.Title); err != nil {
			s.log.Error("failed to send notifications", zap.Error(err))
		}

		s.log.Info("handover thread created",
			zap.String("book_id", history.BookID),
			zap.String("current_holder", history.ReaderID),
			zap.String("next_holder", nextRequest.UserID))
	}

	return nil
}

func (s *service) GetReadingHistoryExtended(ctx context.Context, bookID, userID string) (*domain.ReadingHistoryExtended, error) {
	history, err := s.handoverRepo.GetActiveReadingHistory(ctx, bookID)
	if err != nil {
		return nil, err
	}

	if history == nil {
		return nil, fmt.Errorf("no active reading history found")
	}

	if history.ReaderID != userID {
		return nil, fmt.Errorf("you are not the current holder")
	}

	return history, nil
}
