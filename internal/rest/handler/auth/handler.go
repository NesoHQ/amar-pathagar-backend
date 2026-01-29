package authhandler

import (
	"github.com/yourusername/online-library/internal/auth"
	"go.uber.org/zap"
)

type Handler struct {
	authSvc auth.Service
	log     *zap.Logger
}

func NewHandler(authSvc auth.Service, log *zap.Logger) *Handler {
	return &Handler{authSvc: authSvc, log: log}
}
