package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/yourusername/online-library/internal/domain"
	"github.com/yourusername/online-library/internal/handover"
	"go.uber.org/zap"
)

type HandoverRepository struct {
	db  *sql.DB
	log *zap.Logger
}

var _ handover.HandoverRepo = (*HandoverRepository)(nil)

func NewHandoverRepository(db *sql.DB, log *zap.Logger) *HandoverRepository {
	return &HandoverRepository{db: db, log: log}
}

// Reading history operations
func (r *HandoverRepository) GetActiveReadingHistory(ctx context.Context, bookID string) (*domain.ReadingHistoryExtended, error) {
	query := `
		SELECT 
			rh.id, rh.book_id, rh.reader_id, rh.start_date, rh.end_date,
			rh.due_date, rh.is_completed, rh.completed_at, rh.next_reader_id,
			rh.delivery_status, rh.marked_delivered_at,
			b.title, b.author, b.cover_url, b.status,
			u.username, u.full_name, u.success_score,
			COALESCE(nu.username, '') as next_username,
			COALESCE(nu.full_name, '') as next_full_name
		FROM reading_history rh
		LEFT JOIN books b ON rh.book_id = b.id
		LEFT JOIN users u ON rh.reader_id = u.id
		LEFT JOIN users nu ON rh.next_reader_id = nu.id
		WHERE rh.book_id = $1 AND rh.end_date IS NULL
		ORDER BY rh.start_date DESC
		LIMIT 1
	`

	row := r.db.QueryRowContext(ctx, query, bookID)

	history := &domain.ReadingHistoryExtended{
		ReadingHistory: &domain.ReadingHistory{
			Book:   &domain.Book{},
			Reader: &domain.User{},
		},
	}

	var dueDate, completedAt, markedDeliveredAt sql.NullTime
	var nextReaderID sql.NullString
	var coverURL sql.NullString
	var nextUsername, nextFullName sql.NullString

	err := row.Scan(
		&history.ID, &history.BookID, &history.ReaderID, &history.StartDate, &history.EndDate,
		&dueDate, &history.IsCompleted, &completedAt, &nextReaderID,
		&history.DeliveryStatus, &markedDeliveredAt,
		&history.Book.Title, &history.Book.Author, &coverURL, &history.Book.Status,
		&history.Reader.Username, &history.Reader.FullName, &history.Reader.SuccessScore,
		&nextUsername, &nextFullName,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if dueDate.Valid {
		history.DueDate = &dueDate.Time
	}
	if completedAt.Valid {
		history.CompletedAt = &completedAt.Time
	}
	if markedDeliveredAt.Valid {
		history.MarkedDeliveredAt = &markedDeliveredAt.Time
	}
	if nextReaderID.Valid {
		history.NextReaderID = &nextReaderID.String
		if nextUsername.Valid {
			history.NextReader = &domain.User{
				ID:       nextReaderID.String,
				Username: nextUsername.String,
				FullName: nextFullName.String,
			}
		}
	}
	if coverURL.Valid {
		history.Book.CoverURL = coverURL.String
	}

	history.Book.ID = history.BookID
	history.Reader.ID = history.ReaderID

	return history, nil
}

func (r *HandoverRepository) GetLastCompletedReadingHistory(ctx context.Context, bookID string) (*domain.ReadingHistoryExtended, error) {
	query := `
		SELECT 
			rh.id, rh.book_id, rh.reader_id, rh.start_date, rh.end_date,
			rh.due_date, rh.is_completed, rh.completed_at,
			b.title, b.author, b.cover_url, b.status,
			u.username, u.full_name, u.success_score
		FROM reading_history rh
		LEFT JOIN books b ON rh.book_id = b.id
		LEFT JOIN users u ON rh.reader_id = u.id
		WHERE rh.book_id = $1 AND rh.is_completed = true AND rh.end_date IS NOT NULL
		ORDER BY rh.completed_at DESC
		LIMIT 1
	`

	row := r.db.QueryRowContext(ctx, query, bookID)

	history := &domain.ReadingHistoryExtended{
		ReadingHistory: &domain.ReadingHistory{
			Book:   &domain.Book{},
			Reader: &domain.User{},
		},
	}

	var dueDate, endDate, completedAt sql.NullTime
	var coverURL sql.NullString

	err := row.Scan(
		&history.ID, &history.BookID, &history.ReaderID, &history.StartDate, &endDate,
		&dueDate, &history.IsCompleted, &completedAt,
		&history.Book.Title, &history.Book.Author, &coverURL, &history.Book.Status,
		&history.Reader.Username, &history.Reader.FullName, &history.Reader.SuccessScore,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if dueDate.Valid {
		history.DueDate = &dueDate.Time
	}
	if endDate.Valid {
		history.EndDate = &endDate.Time
	}
	if completedAt.Valid {
		history.CompletedAt = &completedAt.Time
	}
	if coverURL.Valid {
		history.Book.CoverURL = coverURL.String
	}

	history.Book.ID = history.BookID
	history.Reader.ID = history.ReaderID

	return history, nil
}

func (r *HandoverRepository) UpdateReadingHistoryCompleted(ctx context.Context, historyID string, completedAt time.Time) error {
	query := `
		UPDATE reading_history 
		SET is_completed = true, completed_at = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, completedAt, historyID)
	return err
}

func (r *HandoverRepository) UpdateReadingHistoryDeliveryStatus(ctx context.Context, historyID string, status domain.DeliveryStatus, deliveredAt *time.Time) error {
	query := `
		UPDATE reading_history 
		SET delivery_status = $1, marked_delivered_at = $2, updated_at = NOW()
		WHERE id = $3
	`
	_, err := r.db.ExecContext(ctx, query, status, deliveredAt, historyID)
	return err
}

func (r *HandoverRepository) GetReadingHistoriesDueSoon(ctx context.Context, daysThreshold int) ([]*domain.ReadingHistoryExtended, error) {
	query := `
		SELECT 
			rh.id, rh.book_id, rh.reader_id, rh.start_date, rh.due_date,
			rh.is_completed, rh.next_reader_id, rh.delivery_status,
			b.title, b.author
		FROM reading_history rh
		LEFT JOIN books b ON rh.book_id = b.id
		WHERE rh.end_date IS NULL 
		  AND rh.due_date IS NOT NULL
		  AND rh.due_date <= NOW() + INTERVAL '1 day' * $1
		  AND rh.due_date > NOW()
		  AND (rh.next_reader_id IS NULL OR rh.next_reader_id = '')
	`

	rows, err := r.db.QueryContext(ctx, query, daysThreshold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []*domain.ReadingHistoryExtended
	for rows.Next() {
		history := &domain.ReadingHistoryExtended{
			ReadingHistory: &domain.ReadingHistory{
				Book: &domain.Book{},
			},
		}

		var dueDate sql.NullTime
		var nextReaderID sql.NullString

		err := rows.Scan(
			&history.ID, &history.BookID, &history.ReaderID, &history.StartDate, &dueDate,
			&history.IsCompleted, &nextReaderID, &history.DeliveryStatus,
			&history.Book.Title, &history.Book.Author,
		)
		if err != nil {
			return nil, err
		}

		if dueDate.Valid {
			history.DueDate = &dueDate.Time
		}
		if nextReaderID.Valid {
			history.NextReaderID = &nextReaderID.String
		}

		history.Book.ID = history.BookID
		histories = append(histories, history)
	}

	return histories, nil
}

func (r *HandoverRepository) UpdateReadingHistoryNextReader(ctx context.Context, historyID, nextReaderID string) error {
	query := `
		UPDATE reading_history 
		SET next_reader_id = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, nextReaderID, historyID)
	return err
}

// Handover thread operations
func (r *HandoverRepository) CreateHandoverThread(ctx context.Context, thread *domain.HandoverThread) error {
	query := `
		INSERT INTO handover_threads (
			id, book_id, current_holder_id, next_holder_id, reading_history_id,
			status, handover_due_date, is_public, created_at, updated_at
		) VALUES (
			gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, $9
		) RETURNING id
	`

	return r.db.QueryRowContext(ctx, query,
		thread.BookID, thread.CurrentHolderID, thread.NextHolderID, thread.ReadingHistoryID,
		thread.Status, thread.HandoverDueDate, thread.IsPublic, thread.CreatedAt, thread.UpdatedAt,
	).Scan(&thread.ID)
}

func (r *HandoverRepository) GetHandoverThreadByID(ctx context.Context, threadID string) (*domain.HandoverThread, error) {
	query := `
		SELECT 
			ht.id, ht.book_id, ht.current_holder_id, ht.next_holder_id,
			ht.reading_history_id, ht.status, ht.handover_due_date, ht.is_public,
			ht.created_at, ht.completed_at, ht.updated_at,
			b.title, b.author, b.cover_url,
			u1.username as current_username, u1.full_name as current_full_name,
			u2.username as next_username, u2.full_name as next_full_name
		FROM handover_threads ht
		LEFT JOIN books b ON ht.book_id = b.id
		LEFT JOIN users u1 ON ht.current_holder_id = u1.id
		LEFT JOIN users u2 ON ht.next_holder_id = u2.id
		WHERE ht.id = $1
	`

	thread := &domain.HandoverThread{
		Book:          &domain.Book{},
		CurrentHolder: &domain.User{},
		NextHolder:    &domain.User{},
	}

	var readingHistoryID sql.NullString
	var completedAt sql.NullTime
	var coverURL sql.NullString

	err := r.db.QueryRowContext(ctx, query, threadID).Scan(
		&thread.ID, &thread.BookID, &thread.CurrentHolderID, &thread.NextHolderID,
		&readingHistoryID, &thread.Status, &thread.HandoverDueDate, &thread.IsPublic,
		&thread.CreatedAt, &completedAt, &thread.UpdatedAt,
		&thread.Book.Title, &thread.Book.Author, &coverURL,
		&thread.CurrentHolder.Username, &thread.CurrentHolder.FullName,
		&thread.NextHolder.Username, &thread.NextHolder.FullName,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if readingHistoryID.Valid {
		thread.ReadingHistoryID = &readingHistoryID.String
	}
	if completedAt.Valid {
		thread.CompletedAt = &completedAt.Time
	}
	if coverURL.Valid {
		thread.Book.CoverURL = coverURL.String
	}

	thread.Book.ID = thread.BookID
	thread.CurrentHolder.ID = thread.CurrentHolderID
	thread.NextHolder.ID = thread.NextHolderID

	return thread, nil
}

func (r *HandoverRepository) GetActiveHandoverThreadByBook(ctx context.Context, bookID string) (*domain.HandoverThread, error) {
	query := `
		SELECT 
			ht.id, ht.book_id, ht.current_holder_id, ht.next_holder_id,
			ht.reading_history_id, ht.status, ht.handover_due_date, ht.is_public,
			ht.created_at, ht.completed_at, ht.updated_at,
			b.title, b.author, b.cover_url,
			u1.username as current_username, u1.full_name as current_full_name,
			u2.username as next_username, u2.full_name as next_full_name
		FROM handover_threads ht
		LEFT JOIN books b ON ht.book_id = b.id
		LEFT JOIN users u1 ON ht.current_holder_id = u1.id
		LEFT JOIN users u2 ON ht.next_holder_id = u2.id
		WHERE ht.book_id = $1 AND ht.status = 'active'
		ORDER BY ht.created_at DESC
		LIMIT 1
	`

	thread := &domain.HandoverThread{
		Book:          &domain.Book{},
		CurrentHolder: &domain.User{},
		NextHolder:    &domain.User{},
	}

	var readingHistoryID sql.NullString
	var completedAt sql.NullTime
	var coverURL sql.NullString

	err := r.db.QueryRowContext(ctx, query, bookID).Scan(
		&thread.ID, &thread.BookID, &thread.CurrentHolderID, &thread.NextHolderID,
		&readingHistoryID, &thread.Status, &thread.HandoverDueDate, &thread.IsPublic,
		&thread.CreatedAt, &completedAt, &thread.UpdatedAt,
		&thread.Book.Title, &thread.Book.Author, &coverURL,
		&thread.CurrentHolder.Username, &thread.CurrentHolder.FullName,
		&thread.NextHolder.Username, &thread.NextHolder.FullName,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if readingHistoryID.Valid {
		thread.ReadingHistoryID = &readingHistoryID.String
	}
	if completedAt.Valid {
		thread.CompletedAt = &completedAt.Time
	}
	if coverURL.Valid {
		thread.Book.CoverURL = coverURL.String
	}

	thread.Book.ID = thread.BookID
	thread.CurrentHolder.ID = thread.CurrentHolderID
	thread.NextHolder.ID = thread.NextHolderID

	return thread, nil
}

func (r *HandoverRepository) GetHandoverThreadsByUser(ctx context.Context, userID string) ([]*domain.HandoverThread, error) {
	query := `
		SELECT 
			ht.id, ht.book_id, ht.current_holder_id, ht.next_holder_id,
			ht.status, ht.handover_due_date, ht.is_public,
			ht.created_at, ht.completed_at,
			b.title, b.author, b.cover_url,
			u1.username as current_username, u1.full_name as current_full_name,
			u2.username as next_username, u2.full_name as next_full_name
		FROM handover_threads ht
		LEFT JOIN books b ON ht.book_id = b.id
		LEFT JOIN users u1 ON ht.current_holder_id = u1.id
		LEFT JOIN users u2 ON ht.next_holder_id = u2.id
		WHERE ht.current_holder_id = $1 OR ht.next_holder_id = $1
		ORDER BY ht.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []*domain.HandoverThread
	for rows.Next() {
		thread := &domain.HandoverThread{
			Book:          &domain.Book{},
			CurrentHolder: &domain.User{},
			NextHolder:    &domain.User{},
		}

		var completedAt sql.NullTime
		var coverURL sql.NullString

		err := rows.Scan(
			&thread.ID, &thread.BookID, &thread.CurrentHolderID, &thread.NextHolderID,
			&thread.Status, &thread.HandoverDueDate, &thread.IsPublic,
			&thread.CreatedAt, &completedAt,
			&thread.Book.Title, &thread.Book.Author, &coverURL,
			&thread.CurrentHolder.Username, &thread.CurrentHolder.FullName,
			&thread.NextHolder.Username, &thread.NextHolder.FullName,
		)
		if err != nil {
			return nil, err
		}

		if completedAt.Valid {
			thread.CompletedAt = &completedAt.Time
		}
		if coverURL.Valid {
			thread.Book.CoverURL = coverURL.String
		}

		thread.Book.ID = thread.BookID
		thread.CurrentHolder.ID = thread.CurrentHolderID
		thread.NextHolder.ID = thread.NextHolderID

		threads = append(threads, thread)
	}

	return threads, nil
}

func (r *HandoverRepository) UpdateHandoverThreadStatus(ctx context.Context, threadID string, status domain.HandoverThreadStatus, completedAt *time.Time) error {
	query := `
		UPDATE handover_threads 
		SET status = $1, completed_at = $2, updated_at = NOW()
		WHERE id = $3
	`
	_, err := r.db.ExecContext(ctx, query, status, completedAt, threadID)
	return err
}

// Handover message operations
func (r *HandoverRepository) CreateHandoverMessage(ctx context.Context, message *domain.HandoverMessage) error {
	query := `
		INSERT INTO handover_messages (
			id, thread_id, user_id, message, is_system_message, created_at
		) VALUES (
			gen_random_uuid(), $1, $2, $3, $4, $5
		) RETURNING id
	`

	return r.db.QueryRowContext(ctx, query,
		message.ThreadID, message.UserID, message.Message, message.IsSystemMessage, message.CreatedAt,
	).Scan(&message.ID)
}

func (r *HandoverRepository) GetHandoverMessagesByThread(ctx context.Context, threadID string) ([]domain.HandoverMessage, error) {
	query := `
		SELECT 
			hm.id, hm.thread_id, hm.user_id, hm.message, hm.is_system_message, hm.created_at,
			u.username, u.full_name, u.avatar_url
		FROM handover_messages hm
		LEFT JOIN users u ON hm.user_id = u.id
		WHERE hm.thread_id = $1
		ORDER BY hm.created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, threadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []domain.HandoverMessage
	for rows.Next() {
		msg := domain.HandoverMessage{
			User: &domain.User{},
		}

		var avatarURL sql.NullString

		err := rows.Scan(
			&msg.ID, &msg.ThreadID, &msg.UserID, &msg.Message, &msg.IsSystemMessage, &msg.CreatedAt,
			&msg.User.Username, &msg.User.FullName, &avatarURL,
		)
		if err != nil {
			return nil, err
		}

		if avatarURL.Valid {
			msg.User.AvatarURL = avatarURL.String
		}
		msg.User.ID = msg.UserID

		messages = append(messages, msg)
	}

	return messages, nil
}

// Book operations
func (r *HandoverRepository) GetNextApprovedRequest(ctx context.Context, bookID string) (*domain.BookRequest, error) {
	query := `
		SELECT 
			br.id, br.book_id, br.user_id, br.status, br.priority_score,
			br.requested_at, br.processed_at, br.due_date,
			u.username, u.full_name, u.success_score
		FROM book_requests br
		LEFT JOIN users u ON br.user_id = u.id
		WHERE br.book_id = $1 AND br.status = 'approved'
		ORDER BY br.priority_score DESC, br.requested_at ASC
		LIMIT 1
	`

	req := &domain.BookRequest{
		User: &domain.User{},
	}

	var processedAt, dueDate sql.NullTime

	err := r.db.QueryRowContext(ctx, query, bookID).Scan(
		&req.ID, &req.BookID, &req.UserID, &req.Status, &req.PriorityScore,
		&req.RequestedAt, &processedAt, &dueDate,
		&req.User.Username, &req.User.FullName, &req.User.SuccessScore,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if processedAt.Valid {
		req.ProcessedAt = &processedAt.Time
	}
	if dueDate.Valid {
		req.DueDate = &dueDate.Time
	}

	req.User.ID = req.UserID

	return req, nil
}

func (r *HandoverRepository) UpdateBookStatus(ctx context.Context, bookID string, status domain.BookStatus) error {
	query := `UPDATE books SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, bookID)
	return err
}

func (r *HandoverRepository) CloseReadingHistory(ctx context.Context, historyID string, endDate time.Time) error {
	query := `
		UPDATE reading_history 
		SET end_date = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, endDate, historyID)
	return err
}

func (r *HandoverRepository) StartNewReadingHistory(ctx context.Context, bookID, userID string) error {
	query := `
		INSERT INTO reading_history (
			id, book_id, reader_id, start_date, created_at, updated_at
		) VALUES (
			gen_random_uuid(), $1, $2, NOW(), NOW(), NOW()
		)
	`
	_, err := r.db.ExecContext(ctx, query, bookID, userID)
	return err
}

func (r *HandoverRepository) AssignBookToUser(ctx context.Context, bookID, userID string) error {
	query := `UPDATE books SET current_holder_id = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, userID, bookID)
	return err
}
