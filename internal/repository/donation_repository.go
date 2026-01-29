package repository

import (
	"context"
	"database/sql"

	"github.com/yourusername/online-library/internal/domain"
	"github.com/yourusername/online-library/internal/donation"
	"go.uber.org/zap"
)

type DonationRepository struct {
	db  *sql.DB
	log *zap.Logger
}

var _ donation.DonationRepo = (*DonationRepository)(nil)

func NewDonationRepository(db *sql.DB, log *zap.Logger) *DonationRepository {
	return &DonationRepository{db: db, log: log}
}

func (r *DonationRepository) Create(ctx context.Context, d *domain.Donation) error {
	query := `INSERT INTO donations (id, donor_id, donation_type, book_id, amount, currency, message, is_public, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.ExecContext(ctx, query, d.ID, d.DonorID, d.DonationType,
		nullString(d.BookID), nullFloat64(d.Amount), d.Currency, d.Message, d.IsPublic, d.CreatedAt)
	return err
}

func (r *DonationRepository) List(ctx context.Context, limit, offset int) ([]*domain.Donation, error) {
	query := `SELECT id, donor_id, donation_type, book_id, amount, currency, message, is_public, created_at
	          FROM donations WHERE is_public = true ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var donations []*domain.Donation
	for rows.Next() {
		d := &domain.Donation{}
		var bookID sql.NullString
		var amount sql.NullFloat64
		err := rows.Scan(&d.ID, &d.DonorID, &d.DonationType, &bookID, &amount, &d.Currency, &d.Message, &d.IsPublic, &d.CreatedAt)
		if err != nil {
			return nil, err
		}
		d.BookID = stringPtr(bookID)
		d.Amount = float64Ptr(amount)
		donations = append(donations, d)
	}
	return donations, nil
}
