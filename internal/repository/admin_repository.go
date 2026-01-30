package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/lib/pq"
	"github.com/yourusername/online-library/internal/admin"
	"github.com/yourusername/online-library/internal/domain"
	"go.uber.org/zap"
)

type AdminRepository struct {
	db  *sql.DB
	log *zap.Logger
}

var _ admin.AdminRepo = (*AdminRepository)(nil)

func NewAdminRepository(db *sql.DB, log *zap.Logger) *AdminRepository {
	return &AdminRepository{db: db, log: log}
}

func (r *AdminRepository) GetPendingRequests(ctx context.Context, limit, offset int) ([]*domain.BookRequest, error) {
	query := `
		SELECT 
			br.id, br.book_id, br.user_id, br.status, br.priority_score,
			br.interest_match_score, br.distance_km, br.requested_at, br.processed_at, br.due_date,
			b.title, b.author, b.cover_url, b.status,
			u.username, u.full_name, u.success_score
		FROM book_requests br
		LEFT JOIN books b ON br.book_id = b.id
		LEFT JOIN users u ON br.user_id = u.id
		WHERE br.status = 'pending'
		ORDER BY br.priority_score DESC, br.requested_at ASC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*domain.BookRequest
	for rows.Next() {
		req := &domain.BookRequest{Book: &domain.Book{}, User: &domain.User{}}
		var distanceKm sql.NullFloat64
		var processedAt, dueDate sql.NullTime
		var coverURL sql.NullString

		err := rows.Scan(
			&req.ID, &req.BookID, &req.UserID, &req.Status, &req.PriorityScore,
			&req.InterestMatchScore, &distanceKm, &req.RequestedAt, &processedAt, &dueDate,
			&req.Book.Title, &req.Book.Author, &coverURL, &req.Book.Status,
			&req.User.Username, &req.User.FullName, &req.User.SuccessScore,
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
		if coverURL.Valid {
			req.Book.CoverURL = coverURL.String
		}
		req.Book.ID = req.BookID
		req.User.ID = req.UserID

		requests = append(requests, req)
	}
	return requests, nil
}

func (r *AdminRepository) GetRequestsByBook(ctx context.Context, bookID string) ([]*domain.BookRequest, error) {
	query := `
		SELECT 
			br.id, br.book_id, br.user_id, br.status, br.priority_score,
			br.interest_match_score, br.distance_km, br.requested_at,
			u.username, u.full_name, u.success_score, u.location_address
		FROM book_requests br
		LEFT JOIN users u ON br.user_id = u.id
		WHERE br.book_id = $1 AND br.status = 'pending'
		ORDER BY br.priority_score DESC, br.requested_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*domain.BookRequest
	for rows.Next() {
		req := &domain.BookRequest{User: &domain.User{}}
		var distanceKm sql.NullFloat64
		var locationAddress sql.NullString

		err := rows.Scan(
			&req.ID, &req.BookID, &req.UserID, &req.Status, &req.PriorityScore,
			&req.InterestMatchScore, &distanceKm, &req.RequestedAt,
			&req.User.Username, &req.User.FullName, &req.User.SuccessScore, &locationAddress,
		)
		if err != nil {
			return nil, err
		}

		if distanceKm.Valid {
			req.DistanceKm = &distanceKm.Float64
		}
		if locationAddress.Valid {
			req.User.LocationAddress = locationAddress.String
		}
		req.User.ID = req.UserID

		requests = append(requests, req)
	}
	return requests, nil
}

func (r *AdminRepository) UpdateRequestStatus(ctx context.Context, requestID string, status string, processedAt string, dueDate *string) error {
	query := `UPDATE book_requests SET status = $1, processed_at = $2, due_date = $3 WHERE id = $4`
	_, err := r.db.ExecContext(ctx, query, status, processedAt, dueDate, requestID)
	return err
}

func (r *AdminRepository) GetAllUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	query := `
		SELECT id, username, email, full_name, role, 
		       COALESCE(avatar_url, ''), COALESCE(success_score, 100),
		       COALESCE(books_shared, 0), COALESCE(books_received, 0),
		       COALESCE(reviews_received, 0), COALESCE(ideas_posted, 0),
		       COALESCE(is_donor, false), created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		u := &domain.User{}
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.Role,
			&u.AvatarURL, &u.SuccessScore, &u.BooksShared, &u.BooksReceived,
			&u.ReviewsReceived, &u.IdeasPosted, &u.IsDonor, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *AdminRepository) UpdateUserRole(ctx context.Context, userID string, role string) error {
	query := `UPDATE users SET role = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, role, userID)
	return err
}

func (r *AdminRepository) GetSystemStats(ctx context.Context) (*admin.SystemStats, error) {
	stats := &admin.SystemStats{}

	// Total users
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&stats.TotalUsers)
	if err != nil {
		return nil, err
	}

	// Total books
	err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM books").Scan(&stats.TotalBooks)
	if err != nil {
		return nil, err
	}

	// Available books
	err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM books WHERE status = 'available'").Scan(&stats.AvailableBooks)
	if err != nil {
		return nil, err
	}

	// Books in circulation
	err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM books WHERE status = 'reading'").Scan(&stats.BooksInCirculation)
	if err != nil {
		return nil, err
	}

	// Pending requests
	err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM book_requests WHERE status = 'pending'").Scan(&stats.PendingRequests)
	if err != nil {
		return nil, err
	}

	// Total donations
	err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM donations").Scan(&stats.TotalDonations)
	if err != nil {
		return nil, err
	}

	// Total ideas
	err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM reading_ideas").Scan(&stats.TotalIdeas)
	if err != nil {
		return nil, err
	}

	// Total reviews
	err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM user_reviews").Scan(&stats.TotalReviews)
	if err != nil {
		return nil, err
	}

	// Average success score
	err = r.db.QueryRowContext(ctx, "SELECT COALESCE(AVG(success_score), 0) FROM users").Scan(&stats.AvgSuccessScore)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (r *AdminRepository) GetAuditLogs(ctx context.Context, limit, offset int) ([]*admin.AuditLog, error) {
	query := `
		SELECT id, user_id, action, resource_type, resource_id, details, 
		       COALESCE(ip_address, ''), COALESCE(user_agent, ''), created_at
		FROM audit_logs
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*admin.AuditLog
	for rows.Next() {
		log := &admin.AuditLog{}
		var userID, resourceID sql.NullString
		var detailsJSON []byte

		err := rows.Scan(&log.ID, &userID, &log.Action, &log.ResourceType, &resourceID,
			&detailsJSON, &log.IPAddress, &log.UserAgent, &log.CreatedAt)
		if err != nil {
			return nil, err
		}

		if userID.Valid {
			log.UserID = &userID.String
		}
		if resourceID.Valid {
			log.ResourceID = &resourceID.String
		}
		if len(detailsJSON) > 0 {
			json.Unmarshal(detailsJSON, &log.Details)
		}

		logs = append(logs, log)
	}
	return logs, nil
}

func (r *AdminRepository) CreateAuditLog(ctx context.Context, log *admin.AuditLog) error {
	detailsJSON, _ := json.Marshal(log.Details)
	query := `
		INSERT INTO audit_logs (id, user_id, action, resource_type, resource_id, details, ip_address, user_agent, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
	`
	_, err := r.db.ExecContext(ctx, query, log.ID, log.UserID, log.Action, log.ResourceType,
		log.ResourceID, detailsJSON, log.IPAddress, log.UserAgent)
	return err
}

func (r *AdminRepository) GetAllBooks(ctx context.Context, limit, offset int, filters admin.BookFilters) ([]*domain.Book, error) {
	query := `
		SELECT id, title, author, COALESCE(cover_url, ''), COALESCE(category, ''),
		       COALESCE(tags, '{}'), COALESCE(topics, '{}'), status, 
		       COALESCE(total_reads, 0), COALESCE(average_rating, 0), created_at
		FROM books
		WHERE 1=1
	`
	args := []interface{}{}
	argPos := 1

	if filters.Search != "" {
		query += fmt.Sprintf(" AND (title ILIKE $%d OR author ILIKE $%d)", argPos, argPos)
		args = append(args, "%"+filters.Search+"%")
		argPos++
	}

	if filters.Category != "" {
		query += fmt.Sprintf(" AND category = $%d", argPos)
		args = append(args, filters.Category)
		argPos++
	}

	if filters.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argPos)
		args = append(args, filters.Status)
		argPos++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*domain.Book
	for rows.Next() {
		b := &domain.Book{}
		err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.CoverURL, &b.Category,
			pq.Array(&b.Tags), pq.Array(&b.Topics), &b.Status,
			&b.TotalReads, &b.AverageRating, &b.CreatedAt)
		if err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}

func (r *AdminRepository) UpdateBookStatus(ctx context.Context, bookID string, status string) error {
	query := `UPDATE books SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, bookID)
	return err
}
