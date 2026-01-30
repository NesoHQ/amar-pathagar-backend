package notificationhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/online-library/internal/notification"
	"github.com/yourusername/online-library/internal/rest/middleware"
	"github.com/yourusername/online-library/internal/rest/response"
	"go.uber.org/zap"
)

type Handler struct {
	notificationSvc notification.Service
	log             *zap.Logger
}

func NewHandler(notificationSvc notification.Service, log *zap.Logger) *Handler {
	return &Handler{notificationSvc: notificationSvc, log: log}
}

func (h *Handler) GetUserNotifications(c *gin.Context) {
	userID := middleware.GetUserID(c)

	notifications, err := h.notificationSvc.GetUserNotifications(c.Request.Context(), userID, 50)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, notifications)
}

func (h *Handler) MarkAsRead(c *gin.Context) {
	notificationID := c.Param("id")

	if err := h.notificationSvc.MarkAsRead(c.Request.Context(), notificationID); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "notification marked as read"})
}

func (h *Handler) MarkAllAsRead(c *gin.Context) {
	userID := middleware.GetUserID(c)

	if err := h.notificationSvc.MarkAllAsRead(c.Request.Context(), userID); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "all notifications marked as read"})
}

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	r.GET("/notifications", h.GetUserNotifications)
	r.PUT("/notifications/:id/read", h.MarkAsRead)
	r.PUT("/notifications/read-all", h.MarkAllAsRead)
}
