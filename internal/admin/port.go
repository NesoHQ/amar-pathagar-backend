package admin

import (
	"context"

	"github.com/yourusername/online-library/internal/domain"
)

type Service interface {
	// Book Request Management
	GetPendingRequests(ctx context.Context, limit, offset int) ([]*domain.BookRequest, error)
	ApproveBookRequest(ctx context.Context, requestID string, dueDate string) error
	RejectBookRequest(ctx context.Context, requestID string, reason string) error
	GetRequestsByBook(ctx context.Context, bookID string) ([]*domain.BookRequest, error)

	// User Management
	GetAllUsers(ctx context.Context, limit, offset int) ([]*domain.User, error)
	AdjustSuccessScore(ctx context.Context, userID string, amount int, reason string) error
	UpdateUserRole(ctx context.Context, userID string, role domain.UserRole) error

	// Statistics
	GetSystemStats(ctx context.Context) (*SystemStats, error)
	GetAuditLogs(ctx context.Context, limit, offset int) ([]*AuditLog, error)

	// Book Management
	GetAllBooks(ctx context.Context, limit, offset int, filters BookFilters) ([]*domain.Book, error)
	UpdateBookStatus(ctx context.Context, bookID string, status domain.BookStatus) error
}

type AdminRepo interface {
	GetPendingRequests(ctx context.Context, limit, offset int) ([]*domain.BookRequest, error)
	GetRequestsByBook(ctx context.Context, bookID string) ([]*domain.BookRequest, error)
	UpdateRequestStatus(ctx context.Context, requestID string, status string, processedAt string, dueDate *string) error
	GetAllUsers(ctx context.Context, limit, offset int) ([]*domain.User, error)
	UpdateUserRole(ctx context.Context, userID string, role string) error
	GetSystemStats(ctx context.Context) (*SystemStats, error)
	GetAuditLogs(ctx context.Context, limit, offset int) ([]*AuditLog, error)
	CreateAuditLog(ctx context.Context, log *AuditLog) error
	GetAllBooks(ctx context.Context, limit, offset int, filters BookFilters) ([]*domain.Book, error)
	UpdateBookStatus(ctx context.Context, bookID string, status string) error
}

type SystemStats struct {
	TotalUsers         int     `json:"total_users"`
	TotalBooks         int     `json:"total_books"`
	AvailableBooks     int     `json:"available_books"`
	BooksInCirculation int     `json:"books_in_circulation"`
	PendingRequests    int     `json:"pending_requests"`
	TotalDonations     int     `json:"total_donations"`
	TotalIdeas         int     `json:"total_ideas"`
	TotalReviews       int     `json:"total_reviews"`
	AvgSuccessScore    float64 `json:"avg_success_score"`
}

type AuditLog struct {
	ID           string                 `json:"id"`
	UserID       *string                `json:"user_id,omitempty"`
	Action       string                 `json:"action"`
	ResourceType string                 `json:"resource_type"`
	ResourceID   *string                `json:"resource_id,omitempty"`
	Details      map[string]interface{} `json:"details,omitempty"`
	IPAddress    string                 `json:"ip_address,omitempty"`
	UserAgent    string                 `json:"user_agent,omitempty"`
	CreatedAt    string                 `json:"created_at"`
}

type BookFilters struct {
	Search   string
	Category string
	Status   string
}
