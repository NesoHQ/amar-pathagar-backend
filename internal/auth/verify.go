package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yourusername/online-library/internal/domain"
	"go.uber.org/zap"
)

func (s *service) generateToken(userID, role string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *service) ValidateToken(ctx context.Context, tokenString string) (*TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		s.log.Warn("failed to parse token", zap.Error(err))
		return nil, domain.ErrInvalidToken
	}

	if !token.Valid {
		s.log.Warn("invalid token")
		return nil, domain.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		s.log.Warn("invalid token claims")
		return nil, domain.ErrInvalidToken
	}

	userID, _ := claims["user_id"].(string)
	role, _ := claims["role"].(string)

	return &TokenClaims{
		UserID: userID,
		Role:   role,
	}, nil
}
