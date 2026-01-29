package donation

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/online-library/internal/domain"
	"go.uber.org/zap"
)

type service struct {
	donationRepo    DonationRepo
	successScoreSvc SuccessScoreSvc
	log             *zap.Logger
}

func NewService(donationRepo DonationRepo, successScoreSvc SuccessScoreSvc, log *zap.Logger) Service {
	return &service{
		donationRepo:    donationRepo,
		successScoreSvc: successScoreSvc,
		log:             log,
	}
}

func (s *service) Create(ctx context.Context, donation *domain.Donation) (*domain.Donation, error) {
	donation.ID = uuid.New().String()
	donation.CreatedAt = time.Now()

	if err := s.donationRepo.Create(ctx, donation); err != nil {
		s.log.Error("failed to create donation", zap.Error(err))
		return nil, err
	}

	// Update success score based on donation type
	if donation.DonationType == domain.DonationTypeBook {
		if err := s.successScoreSvc.ProcessBookDonation(ctx, donation.DonorID, donation.ID); err != nil {
			s.log.Warn("failed to update success score for book donation", zap.Error(err))
		}
	} else if donation.DonationType == domain.DonationTypeMoney {
		if err := s.successScoreSvc.ProcessMoneyDonation(ctx, donation.DonorID, donation.ID); err != nil {
			s.log.Warn("failed to update success score for money donation", zap.Error(err))
		}
	}

	s.log.Info("donation created successfully", zap.String("donation_id", donation.ID))
	return donation, nil
}

func (s *service) List(ctx context.Context, limit, offset int) ([]*domain.Donation, error) {
	donations, err := s.donationRepo.List(ctx, limit, offset)
	if err != nil {
		s.log.Error("failed to list donations", zap.Error(err))
		return nil, err
	}
	return donations, nil
}
