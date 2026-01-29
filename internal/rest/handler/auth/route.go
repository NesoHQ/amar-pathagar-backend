package authhandler

import "github.com/gin-gonic/gin"

func RegisterPublicRoutes(r *gin.RouterGroup, h *Handler) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
	}
}

func RegisterProtectedRoutes(r *gin.RouterGroup, h *Handler) {
	r.GET("/me", h.Me)
}
