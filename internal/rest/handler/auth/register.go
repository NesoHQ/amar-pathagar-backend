package authhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/online-library/internal/rest/response"
	"go.uber.org/zap"
)

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.authSvc.Register(c.Request.Context(), req.Username, req.Email, req.Password, req.FullName)
	if err != nil {
		h.log.Error("registration failed", zap.Error(err))
		response.Error(c, err)
		return
	}

	response.Created(c, AuthResponse{
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
