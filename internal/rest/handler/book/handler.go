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
	Title        string   `json:"title" binding:"required"`
	Author       string   `json:"author" binding:"required"`
	ISBN         string   `json:"isbn"`
	CoverURL     string   `json:"cover_url"`
	Description  string   `json:"description"`
	Category     string   `json:"category"`
	Tags         []string `json:"tags"`
	Topics       []string `json:"topics"`
	PhysicalCode string   `json:"physical_code"`
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
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
		CreatedBy:    &userID,
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
	books, err := h.bookSvc.List(c.Request.Context(), 50, 0)
	if err != nil {
		response.Error(c, err)
		return
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

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	books := r.Group("/books")
	{
		books.GET("", h.List)
		books.GET("/:id", h.GetByID)
		books.POST("", h.Create)
		books.PATCH("/:id", h.Update)
		books.DELETE("/:id", h.Delete)
		books.POST("/:id/request", h.RequestBook)
	}
}
