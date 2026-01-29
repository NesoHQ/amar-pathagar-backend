package logger

import (
	"go.uber.org/zap"
)

// NewLogger creates a new zap logger based on environment
func NewLogger(env string) (*zap.Logger, error) {
	if env == "production" {
		return zap.NewProduction()
	}
	return zap.NewDevelopment()
}
