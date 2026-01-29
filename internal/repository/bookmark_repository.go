package repository

import (
	"context"
	"database/sql"

	"github.com/yourusername/online-library/internal/bookmark"
	"github.com/yourusername/online-library/internal/domain"
	"go.uber.org/zap"
)

type BookmarkRepository struct {
	db  *sql.DB
	log *zap.Logger
}

var _ bookmark.BookmarkRepo = (*BookmarkRepository)(nil)

func NewBookmarkRepository(db *sql.DB, log *zap.Logger) *BookmarkRepository {
	return &BookmarkRepository{db: db, log: log}
}

func (r *BookmarkRepository) Create(ctx context.Context, b *domain.UserBookmark) error {
	query := `INSERT INTO user_bookmarks (id, user_id, book_id, bookmark_type, priority_level, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, b.ID, b.UserID, b.BookID, b.BookmarkType, b.PriorityLevel, b.CreatedAt)
	return err
}

func (r *BookmarkRepository) Delete(ctx context.Context, userID, bookID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM user_bookmarks WHERE user_id = $1 AND book_id = $2`, userID, bookID)
	return err
}

func (r *BookmarkRepository) FindByUserID(ctx context.Context, userID string) ([]*domain.UserBookmark, error) {
	query := `SELECT id, user_id, book_id, bookmark_type, priority_level, created_at
	          FROM user_bookmarks WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookmarks []*domain.UserBookmark
	for rows.Next() {
		b := &domain.UserBookmark{}
		err := rows.Scan(&b.ID, &b.UserID, &b.BookID, &b.BookmarkType, &b.PriorityLevel, &b.CreatedAt)
		if err != nil {
			return nil, err
		}
		bookmarks = append(bookmarks, b)
	}
	return bookmarks, nil
}
