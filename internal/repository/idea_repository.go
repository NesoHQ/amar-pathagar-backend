package repository

import (
	"context"
	"database/sql"

	"github.com/yourusername/online-library/internal/domain"
	"github.com/yourusername/online-library/internal/idea"
	"go.uber.org/zap"
)

type IdeaRepository struct {
	db  *sql.DB
	log *zap.Logger
}

var _ idea.IdeaRepo = (*IdeaRepository)(nil)

func NewIdeaRepository(db *sql.DB, log *zap.Logger) *IdeaRepository {
	return &IdeaRepository{db: db, log: log}
}

func (r *IdeaRepository) Create(ctx context.Context, i *domain.ReadingIdea) error {
	query := `INSERT INTO reading_ideas (id, book_id, user_id, title, content, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query, i.ID, i.BookID, i.UserID, i.Title, i.Content, i.CreatedAt, i.UpdatedAt)
	return err
}

func (r *IdeaRepository) FindByBookID(ctx context.Context, bookID string) ([]*domain.ReadingIdea, error) {
	query := `SELECT id, book_id, user_id, title, content, COALESCE(upvotes, 0), COALESCE(downvotes, 0), created_at, updated_at
	          FROM reading_ideas WHERE book_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ideas []*domain.ReadingIdea
	for rows.Next() {
		i := &domain.ReadingIdea{}
		err := rows.Scan(&i.ID, &i.BookID, &i.UserID, &i.Title, &i.Content, &i.Upvotes, &i.Downvotes, &i.CreatedAt, &i.UpdatedAt)
		if err != nil {
			return nil, err
		}
		ideas = append(ideas, i)
	}
	return ideas, nil
}

func (r *IdeaRepository) AddVote(ctx context.Context, vote *domain.IdeaVote) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check if user already voted
	var existingVoteType string
	err = tx.QueryRowContext(ctx, `SELECT vote_type FROM idea_votes WHERE idea_id = $1 AND user_id = $2`,
		vote.IdeaID, vote.UserID).Scan(&existingVoteType)

	if err == sql.ErrNoRows {
		// No existing vote, insert new one
		_, err = tx.ExecContext(ctx, `INSERT INTO idea_votes (id, idea_id, user_id, vote_type, created_at) VALUES ($1, $2, $3, $4, $5)`,
			vote.ID, vote.IdeaID, vote.UserID, vote.VoteType, vote.CreatedAt)
		if err != nil {
			return err
		}

		// Update vote counts
		if vote.VoteType == domain.VoteTypeUp {
			_, err = tx.ExecContext(ctx, `UPDATE reading_ideas SET upvotes = upvotes + 1 WHERE id = $1`, vote.IdeaID)
		} else {
			_, err = tx.ExecContext(ctx, `UPDATE reading_ideas SET downvotes = downvotes + 1 WHERE id = $1`, vote.IdeaID)
		}
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		// User already voted
		if existingVoteType == string(vote.VoteType) {
			// Same vote type, remove the vote (toggle off)
			_, err = tx.ExecContext(ctx, `DELETE FROM idea_votes WHERE idea_id = $1 AND user_id = $2`, vote.IdeaID, vote.UserID)
			if err != nil {
				return err
			}

			// Decrement vote count
			if vote.VoteType == domain.VoteTypeUp {
				_, err = tx.ExecContext(ctx, `UPDATE reading_ideas SET upvotes = GREATEST(upvotes - 1, 0) WHERE id = $1`, vote.IdeaID)
			} else {
				_, err = tx.ExecContext(ctx, `UPDATE reading_ideas SET downvotes = GREATEST(downvotes - 1, 0) WHERE id = $1`, vote.IdeaID)
			}
			if err != nil {
				return err
			}
		} else {
			// Different vote type, update the vote
			_, err = tx.ExecContext(ctx, `UPDATE idea_votes SET vote_type = $1 WHERE idea_id = $2 AND user_id = $3`,
				vote.VoteType, vote.IdeaID, vote.UserID)
			if err != nil {
				return err
			}

			// Update vote counts (decrement old, increment new)
			if vote.VoteType == domain.VoteTypeUp {
				_, err = tx.ExecContext(ctx, `UPDATE reading_ideas SET upvotes = upvotes + 1, downvotes = GREATEST(downvotes - 1, 0) WHERE id = $1`, vote.IdeaID)
			} else {
				_, err = tx.ExecContext(ctx, `UPDATE reading_ideas SET downvotes = downvotes + 1, upvotes = GREATEST(upvotes - 1, 0) WHERE id = $1`, vote.IdeaID)
			}
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *IdeaRepository) RemoveVote(ctx context.Context, ideaID, userID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM idea_votes WHERE idea_id = $1 AND user_id = $2`, ideaID, userID)
	return err
}

func (r *IdeaRepository) UpdateVoteCounts(ctx context.Context, ideaID string, upvotes, downvotes int) error {
	_, err := r.db.ExecContext(ctx, `UPDATE reading_ideas SET upvotes = $1, downvotes = $2 WHERE id = $3`, upvotes, downvotes, ideaID)
	return err
}
