package domain

import "time"

type UserBookmark struct {
	ID            string       `json:"id"`
	UserID        string       `json:"user_id"`
	BookID        string       `json:"book_id"`
	Book          *Book        `json:"book,omitempty"`
	BookmarkType  BookmarkType `json:"bookmark_type"`
	PriorityLevel int          `json:"priority_level"`
	CreatedAt     time.Time    `json:"created_at"`
}

type BookmarkType string

const (
	BookmarkTypeWishlist BookmarkType = "wishlist"
	BookmarkTypeFavorite BookmarkType = "favorite"
	BookmarkTypeReading  BookmarkType = "reading"
)
