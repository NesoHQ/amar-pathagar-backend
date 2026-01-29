package user

import (
	"context"
	"time"

	"github.com/yourusername/online-library/internal/domain"
	"go.uber.org/zap"
)

type service struct {
	userRepo UserRepo
	log      *zap.Logger
}

func NewService(userRepo UserRepo, log *zap.Logger) Service {
	return &service{
		userRepo: userRepo,
		log:      log,
	}
}

func (s *service) GetProfile(ctx context.Context, userID string) (*domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		s.log.Error("failed to get user profile", zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}
	return user, nil
}

func (s *service) UpdateProfile(ctx context.Context, userID string, user *domain.User) (*domain.User, error) {
	user.UpdatedAt = time.Now()
	if err := s.userRepo.Update(ctx, userID, user); err != nil {
		s.log.Error("failed to update user profile", zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}
	s.log.Info("user profile updated", zap.String("user_id", userID))
	return user, nil
}

func (s *service) AddInterests(ctx context.Context, userID string, interests []string) error {
	if err := s.userRepo.AddInterests(ctx, userID, interests); err != nil {
		s.log.Error("failed to add interests", zap.String("user_id", userID), zap.Error(err))
		return err
	}
	s.log.Info("interests added", zap.String("user_id", userID))
	return nil
}

func (s *service) GetLeaderboard(ctx context.Context, limit int) ([]*domain.User, error) {
	users, err := s.userRepo.GetTopUsers(ctx, limit)
	if err != nil {
		s.log.Error("failed to get leaderboard", zap.Error(err))
		return nil, err
	}
	return users, nil
}
