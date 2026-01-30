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
	Update(ctx context.Context, id string, book *domain.Book) (*domain.Book, error)
	Delete(ctx context.Context, id string) error
	RequestBook(ctx context.Context, bookID, userID string) (*domain.BookRequest, error)
}

// BookRepo defines the book repository interface
type BookRepo interface {
	Create(ctx context.Context, book *domain.Book) error
	FindByID(ctx context.Context, id string) (*domain.Book, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Book, error)
	Update(ctx context.Context, id string, book *domain.Book) error
	Delete(ctx context.Context, id string) error
	CreateRequest(ctx context.Context, request *domain.BookRequest) error
}
