package repository

import (
	"context"
	"database/sql"

	"github.com/yourusername/online-library/internal/domain"
	"github.com/yourusername/online-library/internal/review"
	"go.uber.org/zap"
)

type ReviewRepository struct {
	db  *sql.DB
	log *zap.Logger
}

var _ review.ReviewRepo = (*ReviewRepository)(nil)

func NewReviewRepository(db *sql.DB, log *zap.Logger) *ReviewRepository {
	return &ReviewRepository{db: db, log: log}
}

func (r *ReviewRepository) Create(ctx context.Context, rev *domain.UserReview) error {
	query := `INSERT INTO user_reviews (id, reviewer_id, reviewee_id, book_id, behavior_rating, 
	                                     book_condition_rating, communication_rating, comment, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.ExecContext(ctx, query, rev.ID, rev.ReviewerID, rev.RevieweeID,
		nullString(rev.BookID), nullInt64(rev.BehaviorRating), nullInt64(rev.BookConditionRating),
		nullInt64(rev.CommunicationRating), rev.Comment, rev.CreatedAt)
	return err
}

func (r *ReviewRepository) FindByUserID(ctx context.Context, userID string) ([]*domain.UserReview, error) {
	query := `SELECT id, reviewer_id, reviewee_id, book_id, behavior_rating, 
	                 book_condition_rating, communication_rating, comment, created_at
	          FROM user_reviews WHERE reviewee_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*domain.UserReview
	for rows.Next() {
		rev := &domain.UserReview{}
		var bookID sql.NullString
		var behaviorRating, bookConditionRating, communicationRating sql.NullInt64
		err := rows.Scan(&rev.ID, &rev.ReviewerID, &rev.RevieweeID, &bookID,
			&behaviorRating, &bookConditionRating, &communicationRating, &rev.Comment, &rev.CreatedAt)
		if err != nil {
			return nil, err
		}
		rev.BookID = stringPtr(bookID)
		rev.BehaviorRating = intPtr(behaviorRating)
		rev.BookConditionRating = intPtr(bookConditionRating)
		rev.CommunicationRating = intPtr(communicationRating)
		reviews = append(reviews, rev)
	}
	return reviews, nil
}
