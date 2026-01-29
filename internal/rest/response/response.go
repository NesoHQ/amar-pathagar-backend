package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/online-library/internal/domain"
)

// Success sends a successful response
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

// Created sends a created response
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    data,
	})
}

// Error sends an error response
func Error(c *gin.Context, err error) {
	statusCode := http.StatusInternalServerError
	message := "internal server error"

	// Map domain errors to HTTP status codes
	switch err {
	case domain.ErrNotFound, domain.ErrUserNotFound:
		statusCode = http.StatusNotFound
		message = err.Error()
	case domain.ErrInvalidCredentials, domain.ErrInvalidToken, domain.ErrTokenExpired:
		statusCode = http.StatusUnauthorized
		message = err.Error()
	case domain.ErrEmailExists, domain.ErrUsernameExists, domain.ErrAlreadyExists:
		statusCode = http.StatusConflict
		message = err.Error()
	case domain.ErrInvalidInput:
		statusCode = http.StatusBadRequest
		message = err.Error()
	case domain.ErrForbidden:
		statusCode = http.StatusForbidden
		message = err.Error()
	case domain.ErrBookNotAvailable, domain.ErrBookAlreadyBorrowed:
		statusCode = http.StatusBadRequest
		message = err.Error()
	default:
		if err != nil {
			message = err.Error()
		}
	}

	c.JSON(statusCode, gin.H{
		"success": false,
		"error":   message,
	})
}

// BadRequest sends a bad request response
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"error":   message,
	})
}

// Unauthorized sends an unauthorized response
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"success": false,
		"error":   message,
	})
}

// NotFound sends a not found response
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, gin.H{
		"success": false,
		"error":   message,
	})
}
