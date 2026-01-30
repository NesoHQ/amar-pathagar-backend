package domain

import "time"

type HandoverThread struct {
	ID               string     `json:"id"`
	BookID           string     `json:"book_id"`
	CurrentHolderID  string     `json:"current_holder_id"`
	NextHolderID     string     `json:"next_holder_id"`
	ReadingHistoryID *string    `json:"reading_history_id,omitempty"`
	Status           string     `json:"status"` // active, completed, cancelled
	HandoverDueDate  time.Time  `json:"handover_due_date"`
	IsPublic         bool       `json:"is_public"`
	CreatedAt        time.Time  `json:"created_at"`
	CompletedAt      *time.Time `json:"completed_at,omitempty"`
	UpdatedAt        time.Time  `json:"updated_at"`

	// Populated fields
	Book          *Book             `json:"book,omitempty"`
	CurrentHolder *User             `json:"current_holder,omitempty"`
	NextHolder    *User             `json:"next_holder,omitempty"`
	Messages      []HandoverMessage `json:"messages,omitempty"`
}

type HandoverMessage struct {
	ID              string    `json:"id"`
	ThreadID        string    `json:"thread_id"`
	UserID          string    `json:"user_id"`
	Message         string    `json:"message"`
	IsSystemMessage bool      `json:"is_system_message"`
	CreatedAt       time.Time `json:"created_at"`

	// Populated fields
	User *User `json:"user,omitempty"`
}

type ReadingHistoryExtended struct {
	*ReadingHistory
	DueDate           *time.Time `json:"due_date,omitempty"`
	IsCompleted       bool       `json:"is_completed"`
	CompletedAt       *time.Time `json:"completed_at,omitempty"`
	NextReaderID      *string    `json:"next_reader_id,omitempty"`
	DeliveryStatus    string     `json:"delivery_status"` // not_started, in_transit, delivered
	MarkedDeliveredAt *time.Time `json:"marked_delivered_at,omitempty"`
	NextReader        *User      `json:"next_reader,omitempty"`
}

type DeliveryStatus string

const (
	DeliveryNotStarted DeliveryStatus = "not_started"
	DeliveryInTransit  DeliveryStatus = "in_transit"
	DeliveryDelivered  DeliveryStatus = "delivered"
)

type HandoverThreadStatus string

const (
	HandoverActive    HandoverThreadStatus = "active"
	HandoverCompleted HandoverThreadStatus = "completed"
	HandoverCancelled HandoverThreadStatus = "cancelled"
)
