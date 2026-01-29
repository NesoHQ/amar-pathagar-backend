package bookmarkhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/online-library/internal/bookmark"
	"github.com/yourusername/online-library/internal/domain"
	"github.com/yourusername/online-library/internal/rest/middleware"
	"github.com/yourusername/online-library/internal/rest/response"
	"go.uber.org/zap"
)

type Handler struct {
	bookmarkSvc bookmark.Service
	log         *zap.Logger
}

func NewHandler(bookmarkSvc bookmark.Service, log *zap.Logger) *Handler {
	return &Handler{bookmarkSvc: bookmarkSvc, log: log}
}

type CreateBookmarkRequest struct {
	BookID        string `json:"book_id" binding:"required"`
	BookmarkType  string `json:"bookmark_type" binding:"required"`
	PriorityLevel int    `json:"priority_level"`
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateBookmarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	bookmark := &domain.UserBookmark{
		UserID:        userID,
		BookID:        req.BookID,
		BookmarkType:  domain.BookmarkType(req.BookmarkType),
		PriorityLevel: req.PriorityLevel,
	}

	created, err := h.bookmarkSvc.Create(c.Request.Context(), bookmark)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, created)
}

func (h *Handler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	bookID := c.Param("bookId")

	if err := h.bookmarkSvc.Delete(c.Request.Context(), userID, bookID); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "bookmark deleted"})
}

func (h *Handler) GetByUser(c *gin.Context) {
	userID := middleware.GetUserID(c)
	bookmarks, err := h.bookmarkSvc.GetByUser(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, bookmarks)
}

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	r.POST("/bookmarks", h.Create)
	r.DELETE("/bookmarks/:bookId", h.Delete)
	r.GET("/bookmarks", h.GetByUser)
}
