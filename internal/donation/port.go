package donation

import (
	"context"

	"github.com/yourusername/online-library/internal/domain"
)

type Service interface {
	Create(ctx context.Context, donation *domain.Donation) (*domain.Donation, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Donation, error)
}

type DonationRepo interface {
	Create(ctx context.Context, donation *domain.Donation) error
	List(ctx context.Context, limit, offset int) ([]*domain.Donation, error)
}

type SuccessScoreSvc interface {
	ProcessBookDonation(ctx context.Context, userID, donationID string) error
	ProcessMoneyDonation(ctx context.Context, userID, donationID string) error
}
