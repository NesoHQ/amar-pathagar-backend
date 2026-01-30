package bookhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/online-library/internal/book"
	"github.com/yourusername/online-library/internal/domain"
	"github.com/yourusername/online-library/internal/rest/middleware"
	"github.com/yourusername/online-library/internal/rest/response"
	"go.uber.org/zap"
)

type Handler struct {
	bookSvc book.Service
	log     *zap.Logger
}

func NewHandler(bookSvc book.Service, log *zap.Logger) *Handler {
	return &Handler{bookSvc: bookSvc, log: log}
}

type CreateBookRequest struct {
	Title          string   `json:"title" binding:"required"`
	Author         string   `json:"author" binding:"required"`
	ISBN           string   `json:"isbn"`
	CoverURL       string   `json:"cover_url"`
	Description    string   `json:"description"`
	Category       string   `json:"category"`
	Tags           []string `json:"tags"`
	Topics         []string `json:"topics"`
	PhysicalCode   string   `json:"physical_code"`
	MaxReadingDays int      `json:"max_reading_days"`
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)

	// Default to 14 days if not provided
	maxReadingDays := req.MaxReadingDays
	if maxReadingDays <= 0 {
		maxReadingDays = 14
	}

	book := &domain.Book{
		Title:          req.Title,
		Author:         req.Author,
		ISBN:           req.ISBN,
		CoverURL:       req.CoverURL,
		Description:    req.Description,
		Category:       req.Category,
		Tags:           req.Tags,
		Topics:         req.Topics,
		PhysicalCode:   req.PhysicalCode,
		MaxReadingDays: maxReadingDays,
		CreatedBy:      &userID,
	}

	created, err := h.bookSvc.Create(c.Request.Context(), book)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, created)
}

func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	book, err := h.bookSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, book)
}

func (h *Handler) List(c *gin.Context) {
	search := c.Query("search")
	status := c.Query("status")

	var books []*domain.Book
	var err error

	if search != "" {
		books, err = h.bookSvc.Search(c.Request.Context(), search, 50, 0)
	} else if status != "" {
		books, err = h.bookSvc.List(c.Request.Context(), 50, 0) // Will filter in frontend for now
	} else {
		books, err = h.bookSvc.List(c.Request.Context(), 50, 0)
	}

	if err != nil {
		response.Error(c, err)
		return
	}

	// Filter by status if provided
	if status != "" && search == "" {
		filtered := []*domain.Book{}
		for _, book := range books {
			if string(book.Status) == status {
				filtered = append(filtered, book)
			}
		}
		books = filtered
	}

	response.Success(c, books)
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var req CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	book := &domain.Book{
		Title:        req.Title,
		Author:       req.Author,
		ISBN:         req.ISBN,
		CoverURL:     req.CoverURL,
		Description:  req.Description,
		Category:     req.Category,
		Tags:         req.Tags,
		Topics:       req.Topics,
		PhysicalCode: req.PhysicalCode,
	}

	updated, err := h.bookSvc.Update(c.Request.Context(), id, book)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, updated)
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.bookSvc.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, gin.H{"message": "book deleted"})
}

func (h *Handler) RequestBook(c *gin.Context) {
	id := c.Param("id")
	userID := middleware.GetUserID(c)

	request, err := h.bookSvc.RequestBook(c.Request.Context(), id, userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, request)
}

func (h *Handler) GetUserRequests(c *gin.Context) {
	userID := middleware.GetUserID(c)

	requests, err := h.bookSvc.GetUserRequests(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, requests)
}

func (h *Handler) CheckBookRequested(c *gin.Context) {
	id := c.Param("id")
	userID := middleware.GetUserID(c)

	requested, err := h.bookSvc.CheckBookRequested(c.Request.Context(), id, userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"requested": requested})
}

func (h *Handler) CancelRequest(c *gin.Context) {
	id := c.Param("id")
	userID := middleware.GetUserID(c)

	if err := h.bookSvc.CancelRequest(c.Request.Context(), id, userID); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "request cancelled"})
}

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	books := r.Group("/books")
	{
		books.GET("", h.List)
		books.GET("/:id", h.GetByID)
		books.POST("", h.Create)
		books.PATCH("/:id", h.Update)
		books.DELETE("/:id", h.Delete)
		books.POST("/:id/request", h.RequestBook)
		books.DELETE("/:id/request", h.CancelRequest)
		books.GET("/:id/requested", h.CheckBookRequested)
		books.POST("/:id/return", h.ReturnBook)
	}

	// User's book requests and history
	r.GET("/my-requests", h.GetUserRequests)
	r.GET("/my-reading-history", h.GetReadingHistory)
	r.GET("/my-books-on-hold", h.GetBooksOnHold)
}

func (h *Handler) ReturnBook(c *gin.Context) {
	id := c.Param("id")
	userID := middleware.GetUserID(c)

	if err := h.bookSvc.ReturnBook(c.Request.Context(), id, userID); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "book returned successfully"})
}

func (h *Handler) GetReadingHistory(c *gin.Context) {
	userID := middleware.GetUserID(c)

	history, err := h.bookSvc.GetReadingHistory(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, history)
}

func (h *Handler) GetBooksOnHold(c *gin.Context) {
	userID := middleware.GetUserID(c)

	books, err := h.bookSvc.GetBooksOnHold(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, books)
}
