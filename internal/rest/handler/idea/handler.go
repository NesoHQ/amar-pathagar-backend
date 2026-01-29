package ideahandler

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/online-library/internal/domain"
	"github.com/yourusername/online-library/internal/idea"
	"github.com/yourusername/online-library/internal/rest/middleware"
	"github.com/yourusername/online-library/internal/rest/response"
	"go.uber.org/zap"
)

type Handler struct {
	ideaSvc idea.Service
	log     *zap.Logger
}

func NewHandler(ideaSvc idea.Service, log *zap.Logger) *Handler {
	return &Handler{ideaSvc: ideaSvc, log: log}
}

type CreateIdeaRequest struct {
	BookID  string `json:"book_id" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateIdeaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	idea := &domain.ReadingIdea{
		BookID:  req.BookID,
		UserID:  userID,
		Title:   req.Title,
		Content: req.Content,
	}

	created, err := h.ideaSvc.Create(c.Request.Context(), idea)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, created)
}

func (h *Handler) GetByBook(c *gin.Context) {
	bookID := c.Param("bookId")
	ideas, err := h.ideaSvc.GetByBook(c.Request.Context(), bookID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, ideas)
}

func (h *Handler) Vote(c *gin.Context) {
	ideaID := c.Param("id")
	userID := middleware.GetUserID(c)
	voteType := domain.VoteTypeUp
	if c.Query("type") == "down" {
		voteType = domain.VoteTypeDown
	}

	if err := h.ideaSvc.Vote(c.Request.Context(), ideaID, userID, voteType); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "vote recorded"})
}

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	r.POST("/ideas", h.Create)
	r.GET("/books/:bookId/ideas", h.GetByBook)
	r.POST("/ideas/:id/vote", h.Vote)
}
