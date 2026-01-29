package user

import (
	"context"

	"github.com/yourusername/online-library/internal/domain"
)

type Service interface {
	GetProfile(ctx context.Context, userID string) (*domain.User, error)
	UpdateProfile(ctx context.Context, userID string, user *domain.User) (*domain.User, error)
	AddInterests(ctx context.Context, userID string, interests []string) error
	GetLeaderboard(ctx context.Context, limit int) ([]*domain.User, error)
}

type UserRepo interface {
	FindByID(ctx context.Context, id string) (*domain.User, error)
	Update(ctx context.Context, id string, user *domain.User) error
	AddInterests(ctx context.Context, userID string, interests []string) error
	GetTopUsers(ctx context.Context, limit int) ([]*domain.User, error)
}
