package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/online-library/internal/domain"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) Register(ctx context.Context, username, email, password, fullName string) (*AuthResult, error) {
	// Check if user exists
	if _, err := s.userRepo.FindByUsername(ctx, username); err == nil {
		s.log.Warn("username already exists", zap.String("username", username))
		return nil, domain.ErrUsernameExists
	}
	if _, err := s.userRepo.FindByEmail(ctx, email); err == nil {
		s.log.Warn("email already exists", zap.String("email", email))
		return nil, domain.ErrEmailExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("failed to hash password", zap.Error(err))
		return nil, domain.ErrInternalServer
	}

	// Create user
	user := &domain.User{
		ID:           uuid.New().String(),
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		FullName:     fullName,
		Role:         domain.RoleMember,
		SuccessScore: 100, // Starting score
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.log.Error("failed to create user", zap.Error(err))
		return nil, err
	}

	// Generate tokens
	accessToken, err := s.generateToken(user.ID, string(user.Role), 24*time.Hour)
	if err != nil {
		s.log.Error("failed to generate access token", zap.Error(err))
		return nil, err
	}

	refreshToken, err := s.generateToken(user.ID, string(user.Role), 168*time.Hour)
	if err != nil {
		s.log.Error("failed to generate refresh token", zap.Error(err))
		return nil, err
	}

	s.log.Info("user registered successfully", zap.String("user_id", user.ID))

	return &AuthResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}
