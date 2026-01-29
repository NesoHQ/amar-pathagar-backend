package reviewhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/online-library/internal/domain"
	"github.com/yourusername/online-library/internal/rest/middleware"
	"github.com/yourusername/online-library/internal/rest/response"
	"github.com/yourusername/online-library/internal/review"
	"go.uber.org/zap"
)

type Handler struct {
	reviewSvc review.Service
	log       *zap.Logger
}

func NewHandler(reviewSvc review.Service, log *zap.Logger) *Handler {
	return &Handler{reviewSvc: reviewSvc, log: log}
}

type CreateReviewRequest struct {
	RevieweeID          string `json:"reviewee_id" binding:"required"`
	BookID              string `json:"book_id"`
	BehaviorRating      *int   `json:"behavior_rating"`
	BookConditionRating *int   `json:"book_condition_rating"`
	CommunicationRating *int   `json:"communication_rating"`
	Comment             string `json:"comment"`
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	var bookID *string
	if req.BookID != "" {
		bookID = &req.BookID
	}

	review := &domain.UserReview{
		ReviewerID:          userID,
		RevieweeID:          req.RevieweeID,
		BookID:              bookID,
		BehaviorRating:      req.BehaviorRating,
		BookConditionRating: req.BookConditionRating,
		CommunicationRating: req.CommunicationRating,
		Comment:             req.Comment,
	}

	created, err := h.reviewSvc.Create(c.Request.Context(), review)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, created)
}

func (h *Handler) GetByUser(c *gin.Context) {
	userID := c.Param("id")
	reviews, err := h.reviewSvc.GetByUser(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, reviews)
}

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	r.POST("/reviews", h.Create)
	r.GET("/users/:id/reviews", h.GetByUser)
}
