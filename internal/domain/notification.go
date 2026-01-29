package domain

import "time"

type Notification struct {
	ID        string           `json:"id"`
	UserID    string           `json:"user_id"`
	Type      NotificationType `json:"type"`
	Title     string           `json:"title"`
	Message   string           `json:"message"`
	Link      string           `json:"link,omitempty"`
	IsRead    bool             `json:"is_read"`
	CreatedAt time.Time        `json:"created_at"`
}

type NotificationType string

const (
	NotificationTypeBookAvailable NotificationType = "book_available"
	NotificationTypeBookDue       NotificationType = "book_due"
	NotificationTypeReview        NotificationType = "review"
	NotificationTypeIdea          NotificationType = "idea"
	NotificationTypeDonation      NotificationType = "donation"
)
