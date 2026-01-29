package domain

import "errors"

var (
	// Common errors
	ErrNotFound       = errors.New("resource not found")
	ErrAlreadyExists  = errors.New("resource already exists")
	ErrInvalidInput   = errors.New("invalid input")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrInternalServer = errors.New("internal server error")

	// Auth errors
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
	ErrEmailExists        = errors.New("email already exists")
	ErrUsernameExists     = errors.New("username already exists")

	// Book errors
	ErrBookNotAvailable    = errors.New("book not available")
	ErrBookAlreadyBorrowed = errors.New("book already borrowed")
	ErrInvalidBookStatus   = errors.New("invalid book status")

	// User errors
	ErrUserNotFound      = errors.New("user not found")
	ErrInsufficientScore = errors.New("insufficient success score")
)
