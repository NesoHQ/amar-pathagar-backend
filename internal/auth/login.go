package auth

import (
	"context"
	"time"

	"github.com/yourusername/online-library/internal/domain"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) Login(ctx context.Context, emailOrUsername, password string) (*AuthResult, error) {
	var user *domain.User
	var err error

	// Try to find user by email first, then username
	user, err = s.userRepo.FindByEmail(ctx, emailOrUsername)
	if err != nil {
		user, err = s.userRepo.FindByUsername(ctx, emailOrUsername)
		if err != nil {
			s.log.Warn("user not found", zap.String("identifier", emailOrUsername))
			return nil, domain.ErrInvalidCredentials
		}
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		s.log.Warn("invalid password", zap.String("user_id", user.ID))
		return nil, domain.ErrInvalidCredentials
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

	s.log.Info("user logged in successfully", zap.String("user_id", user.ID))

	return &AuthResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}
