package book

import (
	"context"

	"github.com/yourusername/online-library/internal/domain"
)

// Service defines the book service interface
type Service interface {
	Create(ctx context.Context, book *domain.Book) (*domain.Book, error)
	GetByID(ctx context.Context, id string) (*domain.Book, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Book, error)
	Search(ctx context.Context, query string, limit, offset int) ([]*domain.Book, error)
	Update(ctx context.Context, id string, book *domain.Book) (*domain.Book, error)
	Delete(ctx context.Context, id string) error
	RequestBook(ctx context.Context, bookID, userID string) (*domain.BookRequest, error)
	GetUserRequests(ctx context.Context, userID string) ([]*domain.BookRequest, error)
	CheckBookRequested(ctx context.Context, bookID, userID string) (bool, error)
	CancelRequest(ctx context.Context, bookID, userID string) error
	ReturnBook(ctx context.Context, bookID, userID string) error
	GetReadingHistory(ctx context.Context, userID string) ([]*domain.ReadingHistory, error)
}

// BookRepo defines the book repository interface
type BookRepo interface {
	Create(ctx context.Context, book *domain.Book) error
	FindByID(ctx context.Context, id string) (*domain.Book, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Book, error)
	Search(ctx context.Context, query string, limit, offset int) ([]*domain.Book, error)
	FilterByStatus(ctx context.Context, status string, limit, offset int) ([]*domain.Book, error)
	Update(ctx context.Context, id string, book *domain.Book) error
	Delete(ctx context.Context, id string) error
	CreateRequest(ctx context.Context, request *domain.BookRequest) error
	FindRequestsByUserID(ctx context.Context, userID string) ([]*domain.BookRequest, error)
	FindRequestByBookAndUser(ctx context.Context, bookID, userID string) (*domain.BookRequest, error)
	CancelRequest(ctx context.Context, bookID, userID string) error
	ReturnBook(ctx context.Context, bookID string) error
	CompleteReadingHistory(ctx context.Context, bookID, userID string) error
	GetReadingHistoryByUser(ctx context.Context, userID string) ([]*domain.ReadingHistory, error)
}
