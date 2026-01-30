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
		                   tags, topics, physical_code, status, max_reading_days, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`
	_, err := r.db.ExecContext(ctx, query,
		b.ID, b.Title, b.Author, b.ISBN, b.CoverURL, b.Description, b.Category,
		pq.Array(b.Tags), pq.Array(b.Topics), b.PhysicalCode, b.Status, b.MaxReadingDays,
		b.CreatedBy, b.CreatedAt, b.UpdatedAt)
	return err
}

func (r *BookRepository) FindByID(ctx context.Context, id string) (*domain.Book, error) {
	b := &domain.Book{}
	query := `
		SELECT id, title, author, COALESCE(isbn, ''), COALESCE(cover_url, ''),
		       COALESCE(description, ''), COALESCE(category, ''),
		       COALESCE(tags, '{}'), COALESCE(topics, '{}'),
		       COALESCE(physical_code, ''), status, COALESCE(max_reading_days, 14), current_holder_id,
		       COALESCE(is_donated, false), COALESCE(total_reads, 0),
		       COALESCE(average_rating, 0), created_at, updated_at
		FROM books WHERE id = $1
	`
	var currentHolderID sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&b.ID, &b.Title, &b.Author, &b.ISBN, &b.CoverURL, &b.Description, &b.Category,
		pq.Array(&b.Tags), pq.Array(&b.Topics), &b.PhysicalCode, &b.Status, &b.MaxReadingDays,
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

func (r *BookRepository) FindRequestsByUserID(ctx context.Context, userID string) ([]*domain.BookRequest, error) {
	query := `
		SELECT 
			br.id, br.book_id, br.user_id, br.status, br.priority_score,
			br.interest_match_score, br.distance_km, br.requested_at, br.processed_at, br.due_date,
			b.id, b.title, b.author, COALESCE(b.cover_url, ''), COALESCE(b.category, ''),
			b.status, COALESCE(b.average_rating, 0)
		FROM book_requests br
		LEFT JOIN books b ON br.book_id = b.id
		WHERE br.user_id = $1 AND br.status = 'pending'
		ORDER BY br.requested_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*domain.BookRequest
	for rows.Next() {
		req := &domain.BookRequest{}
		book := &domain.Book{}
		var distanceKm sql.NullFloat64
		var processedAt, dueDate sql.NullTime

		err := rows.Scan(
			&req.ID, &req.BookID, &req.UserID, &req.Status, &req.PriorityScore,
			&req.InterestMatchScore, &distanceKm, &req.RequestedAt, &processedAt, &dueDate,
			&book.ID, &book.Title, &book.Author, &book.CoverURL, &book.Category,
			&book.Status, &book.AverageRating,
		)
		if err != nil {
			return nil, err
		}

		if distanceKm.Valid {
			req.DistanceKm = &distanceKm.Float64
		}
		if processedAt.Valid {
			req.ProcessedAt = &processedAt.Time
		}
		if dueDate.Valid {
			req.DueDate = &dueDate.Time
		}
		req.Book = book
		requests = append(requests, req)
	}
	return requests, nil
}

func (r *BookRepository) FindRequestByBookAndUser(ctx context.Context, bookID, userID string) (*domain.BookRequest, error) {
	query := `
		SELECT id, book_id, user_id, status, priority_score,
		       interest_match_score, distance_km, requested_at, processed_at, due_date
		FROM book_requests
		WHERE book_id = $1 AND user_id = $2 AND status = 'pending'
		LIMIT 1
	`
	req := &domain.BookRequest{}
	var distanceKm sql.NullFloat64
	var processedAt, dueDate sql.NullTime

	err := r.db.QueryRowContext(ctx, query, bookID, userID).Scan(
		&req.ID, &req.BookID, &req.UserID, &req.Status, &req.PriorityScore,
		&req.InterestMatchScore, &distanceKm, &req.RequestedAt, &processedAt, &dueDate,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if distanceKm.Valid {
		req.DistanceKm = &distanceKm.Float64
	}
	if processedAt.Valid {
		req.ProcessedAt = &processedAt.Time
	}
	if dueDate.Valid {
		req.DueDate = &dueDate.Time
	}

	return req, nil
}

func (r *BookRepository) CancelRequest(ctx context.Context, bookID, userID string) error {
	query := `DELETE FROM book_requests WHERE book_id = $1 AND user_id = $2 AND status = 'pending'`
	result, err := r.db.ExecContext(ctx, query, bookID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *BookRepository) Search(ctx context.Context, query string, limit, offset int) ([]*domain.Book, error) {
	searchQuery := `
		SELECT id, title, author, COALESCE(cover_url, ''), COALESCE(category, ''),
		       status, COALESCE(average_rating, 0), created_at
		FROM books
		WHERE title ILIKE $1 OR author ILIKE $1 OR category ILIKE $1
		   OR $1 = ANY(tags) OR $1 = ANY(topics)
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, searchQuery, "%"+query+"%", limit, offset)
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

func (r *BookRepository) FilterByStatus(ctx context.Context, status string, limit, offset int) ([]*domain.Book, error) {
	query := `
		SELECT id, title, author, COALESCE(cover_url, ''), COALESCE(category, ''),
		       status, COALESCE(average_rating, 0), created_at
		FROM books
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, status, limit, offset)
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

func (r *BookRepository) ReturnBook(ctx context.Context, bookID string) error {
	query := `UPDATE books SET status = 'available', current_holder_id = NULL, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, bookID)
	return err
}

func (r *BookRepository) CompleteReadingHistory(ctx context.Context, bookID, userID string) error {
	query := `
		UPDATE reading_history 
		SET end_date = NOW(), 
		    duration_days = EXTRACT(DAY FROM (NOW() - start_date))::INTEGER,
		    updated_at = NOW()
		WHERE book_id = $1 AND reader_id = $2 AND end_date IS NULL
	`
	_, err := r.db.ExecContext(ctx, query, bookID, userID)
	return err
}

func (r *BookRepository) GetReadingHistoryByUser(ctx context.Context, userID string) ([]*domain.ReadingHistory, error) {
	query := `
		SELECT rh.id, rh.book_id, rh.reader_id, rh.start_date, rh.end_date, 
		       rh.duration_days, COALESCE(rh.notes, ''), rh.rating, COALESCE(rh.review, ''),
		       rh.due_date, rh.is_completed, rh.delivery_status,
		       b.title, b.author, COALESCE(b.cover_url, '')
		FROM reading_history rh
		LEFT JOIN books b ON rh.book_id = b.id
		WHERE rh.reader_id = $1
		ORDER BY rh.start_date DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []*domain.ReadingHistory
	for rows.Next() {
		h := &domain.ReadingHistory{Book: &domain.Book{}}
		var endDate sql.NullTime
		var durationDays sql.NullInt64
		var rating sql.NullInt64
		var dueDate sql.NullTime
		var isCompleted sql.NullBool
		var deliveryStatus sql.NullString

		err := rows.Scan(
			&h.ID, &h.BookID, &h.ReaderID, &h.StartDate, &endDate,
			&durationDays, &h.Notes, &rating, &h.Review,
			&dueDate, &isCompleted, &deliveryStatus,
			&h.Book.Title, &h.Book.Author, &h.Book.CoverURL,
		)
		if err != nil {
			return nil, err
		}

		if endDate.Valid {
			h.EndDate = &endDate.Time
		}
		if durationDays.Valid {
			days := int(durationDays.Int64)
			h.DurationDays = &days
		}
		if rating.Valid {
			r := int(rating.Int64)
			h.Rating = &r
		}

		history = append(history, h)
	}
	return history, nil
}
