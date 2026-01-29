package book

import (
	"go.uber.org/zap"
)

type service struct {
	bookRepo BookRepo
	log      *zap.Logger
}

// NewService creates a new book service
func NewService(bookRepo BookRepo, log *zap.Logger) Service {
	return &service{
		bookRepo: bookRepo,
		log:      log,
	}
}
