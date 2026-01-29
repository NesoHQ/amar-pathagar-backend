package donationhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/online-library/internal/domain"
	"github.com/yourusername/online-library/internal/donation"
	"github.com/yourusername/online-library/internal/rest/middleware"
	"github.com/yourusername/online-library/internal/rest/response"
	"go.uber.org/zap"
)

type Handler struct {
	donationSvc donation.Service
	log         *zap.Logger
}

func NewHandler(donationSvc donation.Service, log *zap.Logger) *Handler {
	return &Handler{donationSvc: donationSvc, log: log}
}

type CreateDonationRequest struct {
	DonationType string   `json:"donation_type" binding:"required"`
	BookID       string   `json:"book_id"`
	Amount       *float64 `json:"amount"`
	Currency     string   `json:"currency"`
	Message      string   `json:"message"`
	IsPublic     bool     `json:"is_public"`
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateDonationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	var bookID *string
	if req.BookID != "" {
		bookID = &req.BookID
	}

	donation := &domain.Donation{
		DonorID:      userID,
		DonationType: domain.DonationType(req.DonationType),
		BookID:       bookID,
		Amount:       req.Amount,
		Currency:     req.Currency,
		Message:      req.Message,
		IsPublic:     req.IsPublic,
	}

	created, err := h.donationSvc.Create(c.Request.Context(), donation)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, created)
}

func (h *Handler) List(c *gin.Context) {
	donations, err := h.donationSvc.List(c.Request.Context(), 50, 0)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, donations)
}

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	r.POST("/donations", h.Create)
	r.GET("/donations", h.List)
}
