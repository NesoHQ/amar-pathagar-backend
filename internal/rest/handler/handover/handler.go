package handover

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/online-library/internal/handover"
	"github.com/yourusername/online-library/internal/rest/response"
	"go.uber.org/zap"
)

type Handler struct {
	handoverSvc handover.Service
	log         *zap.Logger
}

func NewHandler(handoverSvc handover.Service, log *zap.Logger) *Handler {
	return &Handler{
		handoverSvc: handoverSvc,
		log:         log,
	}
}

// MarkBookCompleted marks a book as reading completed by current holder
// POST /api/v1/books/:id/complete
func (h *Handler) MarkBookCompleted(c *gin.Context) {
	userID := c.GetString("user_id")
	bookID := c.Param("id")

	if err := h.handoverSvc.MarkBookCompleted(c.Request.Context(), userID, bookID); err != nil {
		h.log.Error("failed to mark book completed", zap.Error(err))
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "Book marked as completed"})
}

// MarkBookDelivered marks a book as delivered by the receiver
// POST /api/v1/books/:id/delivered
func (h *Handler) MarkBookDelivered(c *gin.Context) {
	userID := c.GetString("user_id")
	bookID := c.Param("id")

	if err := h.handoverSvc.MarkBookDelivered(c.Request.Context(), userID, bookID); err != nil {
		h.log.Error("failed to mark book delivered", zap.Error(err))
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "Book marked as delivered"})
}

// GetActiveHandoverThread gets the active handover thread for a book
// GET /api/v1/books/:id/handover
func (h *Handler) GetActiveHandoverThread(c *gin.Context) {
	bookID := c.Param("id")

	thread, err := h.handoverSvc.GetActiveHandoverThread(c.Request.Context(), bookID)
	if err != nil {
		h.log.Error("failed to get handover thread", zap.Error(err))
		response.Error(c, fmt.Errorf("failed to get handover thread"))
		return
	}

	if thread == nil {
		response.Success(c, nil)
		return
	}

	response.Success(c, thread)
}

// GetUserHandoverThreads gets all handover threads for a user
// GET /api/v1/handover/threads
func (h *Handler) GetUserHandoverThreads(c *gin.Context) {
	userID := c.GetString("user_id")

	threads, err := h.handoverSvc.GetUserHandoverThreads(c.Request.Context(), userID)
	if err != nil {
		h.log.Error("failed to get user handover threads", zap.Error(err))
		response.Error(c, fmt.Errorf("failed to get handover threads"))
		return
	}

	response.Success(c, threads)
}

// PostHandoverMessage posts a message to a handover thread
// POST /api/v1/handover/threads/:id/messages
func (h *Handler) PostHandoverMessage(c *gin.Context) {
	userID := c.GetString("user_id")
	threadID := c.Param("id")

	var req struct {
		Message string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request")
		return
	}

	if err := h.handoverSvc.PostHandoverMessage(c.Request.Context(), threadID, userID, req.Message); err != nil {
		h.log.Error("failed to post handover message", zap.Error(err))
		response.Error(c, err)
		return
	}

	response.Created(c, gin.H{"message": "Message posted"})
}

// GetHandoverMessages gets all messages for a handover thread
// GET /api/v1/handover/threads/:id/messages
func (h *Handler) GetHandoverMessages(c *gin.Context) {
	threadID := c.Param("id")

	messages, err := h.handoverSvc.GetHandoverMessages(c.Request.Context(), threadID)
	if err != nil {
		h.log.Error("failed to get handover messages", zap.Error(err))
		response.Error(c, fmt.Errorf("failed to get messages"))
		return
	}

	response.Success(c, messages)
}

// GetReadingHistoryExtended gets extended reading history for current holder
// GET /api/v1/books/:id/reading-status
func (h *Handler) GetReadingHistoryExtended(c *gin.Context) {
	userID := c.GetString("user_id")
	bookID := c.Param("id")

	history, err := h.handoverSvc.GetReadingHistoryExtended(c.Request.Context(), bookID, userID)
	if err != nil {
		h.log.Error("failed to get reading history", zap.Error(err))
		response.Error(c, err)
		return
	}

	response.Success(c, history)
}

// RegisterRoutes registers handover routes
func RegisterRoutes(router *gin.RouterGroup, h *Handler) {
	books := router.Group("/books")
	{
		books.POST("/:id/complete", h.MarkBookCompleted)
		books.POST("/:id/delivered", h.MarkBookDelivered)
		books.GET("/:id/handover", h.GetActiveHandoverThread)
		books.GET("/:id/reading-status", h.GetReadingHistoryExtended)
	}

	handover := router.Group("/handover")
	{
		handover.GET("/threads", h.GetUserHandoverThreads)
		handover.POST("/threads/:id/messages", h.PostHandoverMessage)
		handover.GET("/threads/:id/messages", h.GetHandoverMessages)
	}
}
