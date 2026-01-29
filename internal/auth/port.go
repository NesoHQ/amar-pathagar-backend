package auth

import (
	"context"

	"github.com/yourusername/online-library/internal/domain"
)

// Service defines the auth service interface
type Service interface {
	Register(ctx context.Context, username, email, password, fullName string) (*AuthResult, error)
	Login(ctx context.Context, emailOrUsername, password string) (*AuthResult, error)
	ValidateToken(ctx context.Context, token string) (*TokenClaims, error)
}

// UserRepo defines the user repository interface
type UserRepo interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
}

// AuthResult represents the authentication result
type AuthResult struct {
	AccessToken  string
	RefreshToken string
	User         *domain.User
}

// TokenClaims represents JWT token claims
type TokenClaims struct {
	UserID string
	Role   string
}
