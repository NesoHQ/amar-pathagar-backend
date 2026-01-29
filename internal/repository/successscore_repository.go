package repository

import (
	"context"
	"database/sql"

	"github.com/yourusername/online-library/internal/successscore"
	"go.uber.org/zap"
)

type SuccessScoreRepository struct {
	db  *sql.DB
	log *zap.Logger
}

var _ successscore.ScoreRepo = (*SuccessScoreRepository)(nil)

func NewSuccessScoreRepository(db *sql.DB, log *zap.Logger) *SuccessScoreRepository {
	return &SuccessScoreRepository{db: db, log: log}
}

func (r *SuccessScoreRepository) UpdateScore(ctx context.Context, userID string, change int, reason, refType string, refID *string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update user success score
	_, err = tx.ExecContext(ctx, `UPDATE users SET success_score = success_score + $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`, change, userID)
	if err != nil {
		return err
	}

	// Record history
	var refIDVal interface{} = nil
	if refID != nil {
		refIDVal = *refID
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO success_score_history (user_id, change_amount, reason, reference_type, reference_id, created_at)
	                               VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)`,
		userID, change, reason, refType, refIDVal)
	if err != nil {
		return err
	}

	return tx.Commit()
}
