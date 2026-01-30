package userhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/online-library/internal/domain"
	"github.com/yourusername/online-library/internal/rest/middleware"
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

type UpdateProfileRequest struct {
	FullName        string   `json:"full_name"`
	Bio             string   `json:"bio"`
	AvatarURL       string   `json:"avatar_url"`
	LocationAddress string   `json:"location_address"`
	LocationLat     *float64 `json:"location_lat"`
	LocationLng     *float64 `json:"location_lng"`
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user := &domain.User{
		FullName:        req.FullName,
		Bio:             req.Bio,
		AvatarURL:       req.AvatarURL,
		LocationAddress: req.LocationAddress,
		LocationLat:     req.LocationLat,
		LocationLng:     req.LocationLng,
	}

	updated, err := h.userSvc.UpdateProfile(c.Request.Context(), userID, user)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, updated)
}

type AddInterestsRequest struct {
	Interests []string `json:"interests" binding:"required"`
}

func (h *Handler) AddInterests(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req AddInterestsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.userSvc.AddInterests(c.Request.Context(), userID, req.Interests); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, gin.H{"message": "interests added"})
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
	r.PUT("/users/profile", h.UpdateProfile)
	r.POST("/users/interests", h.AddInterests)
	r.GET("/leaderboard", h.GetLeaderboard)
}
