package repository

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/yourusername/online-library/internal/book"
	"github.com/yourusername/online-library/internal/domain"
	"go.uber.org/zap"
)

type BookRepository struct {
	db  *sql.DB
	log *zap.Logger
}

var _ book.BookRepo = (*BookRepository)(nil)

func NewBookRepository(db *sql.DB, log *zap.Logger) *BookRepository {
	return &BookRepository{db: db, log: log}
}

func (r *BookRepository) Create(ctx context.Context, b *domain.Book) error {
	query := `
		INSERT INTO books (id, title, author, isbn, cover_url, description, category, 
		                   tags, topics, physical_code, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err := r.db.ExecContext(ctx, query,
		b.ID, b.Title, b.Author, b.ISBN, b.CoverURL, b.Description, b.Category,
		pq.Array(b.Tags), pq.Array(b.Topics), b.PhysicalCode, b.Status,
		b.CreatedAt, b.UpdatedAt)
	return err
}

func (r *BookRepository) FindByID(ctx context.Context, id string) (*domain.Book, error) {
	b := &domain.Book{}
	query := `
		SELECT id, title, author, COALESCE(isbn, ''), COALESCE(cover_url, ''),
		       COALESCE(description, ''), COALESCE(category, ''),
		       COALESCE(tags, '{}'), COALESCE(topics, '{}'),
		       COALESCE(physical_code, ''), status, current_holder_id,
		       COALESCE(is_donated, false), COALESCE(total_reads, 0),
		       COALESCE(average_rating, 0), created_at, updated_at
		FROM books WHERE id = $1
	`
	var currentHolderID sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&b.ID, &b.Title, &b.Author, &b.ISBN, &b.CoverURL, &b.Description, &b.Category,
		pq.Array(&b.Tags), pq.Array(&b.Topics), &b.PhysicalCode, &b.Status,
		&currentHolderID, &b.IsDonated, &b.TotalReads, &b.AverageRating,
		&b.CreatedAt, &b.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	b.CurrentHolderID = stringPtr(currentHolderID)
	return b, err
}

func (r *BookRepository) List(ctx context.Context, limit, offset int) ([]*domain.Book, error) {
	query := `
		SELECT id, title, author, COALESCE(cover_url, ''), COALESCE(category, ''),
		       status, COALESCE(average_rating, 0), created_at
		FROM books
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*domain.Book
	for rows.Next() {
		b := &domain.Book{}
		err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.CoverURL, &b.Category,
			&b.Status, &b.AverageRating, &b.CreatedAt)
		if err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}

func (r *BookRepository) Update(ctx context.Context, id string, b *domain.Book) error {
	query := `
		UPDATE books SET title = $1, author = $2, isbn = $3, cover_url = $4,
		       description = $5, category = $6, tags = $7, topics = $8,
		       status = $9, updated_at = $10
		WHERE id = $11
	`
	_, err := r.db.ExecContext(ctx, query,
		b.Title, b.Author, b.ISBN, b.CoverURL, b.Description, b.Category,
		pq.Array(b.Tags), pq.Array(b.Topics), b.Status, b.UpdatedAt, id)
	return err
}

func (r *BookRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM books WHERE id = $1", id)
	return err
}

func (r *BookRepository) CreateRequest(ctx context.Context, req *domain.BookRequest) error {
	query := `
		INSERT INTO book_requests (id, book_id, user_id, status, priority_score, 
		                          interest_match_score, distance_km, requested_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		req.ID, req.BookID, req.UserID, req.Status, req.PriorityScore,
		req.InterestMatchScore, req.DistanceKm, req.RequestedAt)
	return err
}
