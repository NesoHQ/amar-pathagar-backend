package authhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/online-library/internal/rest/response"
	"go.uber.org/zap"
)

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	emailOrUsername := req.Email
	if emailOrUsername == "" {
		emailOrUsername = req.Username
	}

	result, err := h.authSvc.Login(c.Request.Context(), emailOrUsername, req.Password)
	if err != nil {
		h.log.Error("login failed", zap.Error(err))
		response.Error(c, err)
		return
	}

	response.Success(c, AuthResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		User: &UserDTO{
			ID:           result.User.ID,
			Username:     result.User.Username,
			Email:        result.User.Email,
			FullName:     result.User.FullName,
			Role:         string(result.User.Role),
			SuccessScore: result.User.SuccessScore,
		},
	})
}

func (h *Handler) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	response.Success(c, gin.H{"user_id": userID})
}
