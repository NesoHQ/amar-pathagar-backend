package domain

import "time"

type Book struct {
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	Author          string     `json:"author"`
	ISBN            string     `json:"isbn,omitempty"`
	CoverURL        string     `json:"cover_url,omitempty"`
	Description     string     `json:"description,omitempty"`
	Category        string     `json:"category,omitempty"`
	Tags            []string   `json:"tags,omitempty"`
	Topics          []string   `json:"topics,omitempty"`
	PhysicalCode    string     `json:"physical_code,omitempty"`
	Status          BookStatus `json:"status"`
	MaxReadingDays  int        `json:"max_reading_days"`
	CurrentHolderID *string    `json:"current_holder_id,omitempty"`
	CurrentHolder   *User      `json:"current_holder,omitempty"`
	CreatedBy       *string    `json:"created_by,omitempty"`
	DonatedBy       *string    `json:"donated_by,omitempty"`
	IsDonated       bool       `json:"is_donated"`
	DonationDate    *time.Time `json:"donation_date,omitempty"`
	TotalReads      int        `json:"total_reads"`
	AverageRating   float64    `json:"average_rating"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type BookStatus string

const (
	StatusAvailable BookStatus = "available"
	StatusReading   BookStatus = "reading"
	StatusReserved  BookStatus = "reserved"
	StatusRequested BookStatus = "requested"
)

type ReadingHistory struct {
	ID           string     `json:"id"`
	BookID       string     `json:"book_id"`
	Book         *Book      `json:"book,omitempty"`
	ReaderID     string     `json:"reader_id"`
	Reader       *User      `json:"reader,omitempty"`
	StartDate    time.Time  `json:"start_date"`
	EndDate      *time.Time `json:"end_date,omitempty"`
	DurationDays *int       `json:"duration_days,omitempty"`
	Notes        string     `json:"notes,omitempty"`
	Rating       *int       `json:"rating,omitempty"`
	Review       string     `json:"review,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type WaitingQueue struct {
	ID       string    `json:"id"`
	BookID   string    `json:"book_id"`
	Book     *Book     `json:"book,omitempty"`
	UserID   string    `json:"user_id"`
	User     *User     `json:"user,omitempty"`
	Position int       `json:"position"`
	JoinedAt time.Time `json:"joined_at"`
	Notified bool      `json:"notified"`
}

type BookRequest struct {
	ID                 string     `json:"id"`
	BookID             string     `json:"book_id"`
	Book               *Book      `json:"book,omitempty"`
	UserID             string     `json:"user_id"`
	User               *User      `json:"user,omitempty"`
	Status             string     `json:"status"`
	PriorityScore      float64    `json:"priority_score"`
	InterestMatchScore float64    `json:"interest_match_score"`
	DistanceKm         *float64   `json:"distance_km,omitempty"`
	RequestedAt        time.Time  `json:"requested_at"`
	ProcessedAt        *time.Time `json:"processed_at,omitempty"`
	DueDate            *time.Time `json:"due_date,omitempty"`
}
