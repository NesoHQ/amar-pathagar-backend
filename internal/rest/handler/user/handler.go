package userhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/online-library/internal/rest/response"
	"github.com/yourusername/online-library/internal/user"
	"go.uber.org/zap"
)

type Handler struct {
	userSvc user.Service
	log     *zap.Logger
}

func NewHandler(userSvc user.Service, log *zap.Logger) *Handler {
	return &Handler{userSvc: userSvc, log: log}
}

func (h *Handler) GetProfile(c *gin.Context) {
	userID := c.Param("id")
	user, err := h.userSvc.GetProfile(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, user)
}

func (h *Handler) GetLeaderboard(c *gin.Context) {
	users, err := h.userSvc.GetLeaderboard(c.Request.Context(), 10)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, users)
}

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	r.GET("/users/:id/profile", h.GetProfile)
	r.GET("/leaderboard", h.GetLeaderboard)
}
