package adminhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/online-library/internal/admin"
	"github.com/yourusername/online-library/internal/domain"
	"github.com/yourusername/online-library/internal/rest/response"
	"go.uber.org/zap"
)

type Handler struct {
	adminSvc admin.Service
	log      *zap.Logger
}

func NewHandler(adminSvc admin.Service, log *zap.Logger) *Handler {
	return &Handler{adminSvc: adminSvc, log: log}
}

// GetPendingRequests returns all pending book requests
func (h *Handler) GetPendingRequests(c *gin.Context) {
	requests, err := h.adminSvc.GetPendingRequests(c.Request.Context(), 100, 0)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, requests)
}

// ApproveBookRequest approves a book request
type ApproveRequestReq struct {
	DueDate string `json:"due_date" binding:"required"`
}

func (h *Handler) ApproveBookRequest(c *gin.Context) {
	requestID := c.Param("id")
	var req ApproveRequestReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.adminSvc.ApproveBookRequest(c.Request.Context(), requestID, req.DueDate); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "request approved"})
}

// RejectBookRequest rejects a book request
type RejectRequestReq struct {
	Reason string `json:"reason"`
}

func (h *Handler) RejectBookRequest(c *gin.Context) {
	requestID := c.Param("id")
	var req RejectRequestReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.adminSvc.RejectBookRequest(c.Request.Context(), requestID, req.Reason); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "request rejected"})
}

// GetRequestsByBook returns all requests for a specific book
func (h *Handler) GetRequestsByBook(c *gin.Context) {
	bookID := c.Param("bookId")
	requests, err := h.adminSvc.GetRequestsByBook(c.Request.Context(), bookID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, requests)
}

// GetAllUsers returns all users
func (h *Handler) GetAllUsers(c *gin.Context) {
	users, err := h.adminSvc.GetAllUsers(c.Request.Context(), 100, 0)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, users)
}

// AdjustSuccessScore adjusts a user's success score
type AdjustScoreReq struct {
	Amount int    `json:"amount" binding:"required"`
	Reason string `json:"reason" binding:"required"`
}

func (h *Handler) AdjustSuccessScore(c *gin.Context) {
	userID := c.Param("userId")
	var req AdjustScoreReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.adminSvc.AdjustSuccessScore(c.Request.Context(), userID, req.Amount, req.Reason); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "success score adjusted"})
}

// UpdateUserRole updates a user's role
type UpdateRoleReq struct {
	Role string `json:"role" binding:"required"`
}

func (h *Handler) UpdateUserRole(c *gin.Context) {
	userID := c.Param("userId")
	var req UpdateRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	role := domain.UserRole(req.Role)
	if err := h.adminSvc.UpdateUserRole(c.Request.Context(), userID, role); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "user role updated"})
}

// GetSystemStats returns system statistics
func (h *Handler) GetSystemStats(c *gin.Context) {
	stats, err := h.adminSvc.GetSystemStats(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, stats)
}

// GetAuditLogs returns audit logs
func (h *Handler) GetAuditLogs(c *gin.Context) {
	logs, err := h.adminSvc.GetAuditLogs(c.Request.Context(), 100, 0)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, logs)
}

// GetAllBooks returns all books with filters
func (h *Handler) GetAllBooks(c *gin.Context) {
	filters := admin.BookFilters{
		Search:   c.Query("search"),
		Category: c.Query("category"),
		Status:   c.Query("status"),
	}

	books, err := h.adminSvc.GetAllBooks(c.Request.Context(), 100, 0, filters)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, books)
}

// UpdateBookStatus updates a book's status
type UpdateBookStatusReq struct {
	Status string `json:"status" binding:"required"`
}

func (h *Handler) UpdateBookStatus(c *gin.Context) {
	bookID := c.Param("bookId")
	var req UpdateBookStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	status := domain.BookStatus(req.Status)
	if err := h.adminSvc.UpdateBookStatus(c.Request.Context(), bookID, status); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "book status updated"})
}

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	admin := r.Group("/admin")
	{
		// Statistics
		admin.GET("/stats", h.GetSystemStats)
		admin.GET("/audit-logs", h.GetAuditLogs)

		// Book Request Management
		admin.GET("/requests/pending", h.GetPendingRequests)
		admin.POST("/requests/:id/approve", h.ApproveBookRequest)
		admin.POST("/requests/:id/reject", h.RejectBookRequest)
		admin.GET("/books/:bookId/requests", h.GetRequestsByBook)

		// User Management
		admin.GET("/users", h.GetAllUsers)
		admin.POST("/users/:userId/score", h.AdjustSuccessScore)
		admin.PUT("/users/:userId/role", h.UpdateUserRole)

		// Book Management
		admin.GET("/books", h.GetAllBooks)
		admin.PUT("/books/:bookId/status", h.UpdateBookStatus)
	}
}
