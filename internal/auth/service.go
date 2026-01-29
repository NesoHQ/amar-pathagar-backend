package auth

import (
	"go.uber.org/zap"
)

type service struct {
	userRepo  UserRepo
	jwtSecret string
	log       *zap.Logger
}

// NewService creates a new auth service
func NewService(userRepo UserRepo, jwtSecret string, log *zap.Logger) Service {
	return &service{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		log:       log,
	}
}
