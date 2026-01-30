package handover

import (
	"context"
	"time"

	"github.com/yourusername/online-library/internal/domain"
)

type Service interface {
	// Mark book as reading completed
	MarkBookCompleted(ctx context.Context, userID, bookID string) error

	// Mark book as delivered by the receiver
	MarkBookDelivered(ctx context.Context, userID, bookID string) error

	// Get active handover thread for a book
	GetActiveHandoverThread(ctx context.Context, bookID string) (*domain.HandoverThread, error)

	// Get handover threads for a user (as sender or receiver)
	GetUserHandoverThreads(ctx context.Context, userID string) ([]*domain.HandoverThread, error)

	// Post message to handover thread
	PostHandoverMessage(ctx context.Context, threadID, userID, message string) error

	// Get messages for a handover thread
	GetHandoverMessages(ctx context.Context, threadID string) ([]domain.HandoverMessage, error)

	// Check and create handover threads for books nearing due date (cron job)
	CheckAndCreateHandoverThreads(ctx context.Context) error

	// Get reading history with extended fields
	GetReadingHistoryExtended(ctx context.Context, bookID, userID string) (*domain.ReadingHistoryExtended, error)
}

type HandoverRepo interface {
	// Reading history operations
	GetActiveReadingHistory(ctx context.Context, bookID string) (*domain.ReadingHistoryExtended, error)
	GetLastCompletedReadingHistory(ctx context.Context, bookID string) (*domain.ReadingHistoryExtended, error)
	UpdateReadingHistoryCompleted(ctx context.Context, historyID string, completedAt time.Time) error
	UpdateReadingHistoryDeliveryStatus(ctx context.Context, historyID string, status domain.DeliveryStatus, deliveredAt *time.Time) error
	GetReadingHistoriesDueSoon(ctx context.Context, daysThreshold int) ([]*domain.ReadingHistoryExtended, error)
	CloseReadingHistory(ctx context.Context, historyID string, endDate time.Time) error
	StartNewReadingHistory(ctx context.Context, bookID, userID string) error

	// Handover thread operations
	CreateHandoverThread(ctx context.Context, thread *domain.HandoverThread) error
	GetHandoverThreadByID(ctx context.Context, threadID string) (*domain.HandoverThread, error)
	GetActiveHandoverThreadByBook(ctx context.Context, bookID string) (*domain.HandoverThread, error)
	GetHandoverThreadsByUser(ctx context.Context, userID string) ([]*domain.HandoverThread, error)
	UpdateHandoverThreadStatus(ctx context.Context, threadID string, status domain.HandoverThreadStatus, completedAt *time.Time) error

	// Handover message operations
	CreateHandoverMessage(ctx context.Context, message *domain.HandoverMessage) error
	GetHandoverMessagesByThread(ctx context.Context, threadID string) ([]domain.HandoverMessage, error)

	// Book operations
	GetNextApprovedRequest(ctx context.Context, bookID string) (*domain.BookRequest, error)
	UpdateReadingHistoryNextReader(ctx context.Context, historyID, nextReaderID string) error
	UpdateBookStatus(ctx context.Context, bookID string, status domain.BookStatus) error
	AssignBookToUser(ctx context.Context, bookID, userID string) error
}
