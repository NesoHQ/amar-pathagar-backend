package domain

import "time"

type User struct {
	ID              string    `json:"id"`
	Username        string    `json:"username"`
	Email           string    `json:"email"`
	PasswordHash    string    `json:"-"`
	FullName        string    `json:"full_name"`
	Role            UserRole  `json:"role"`
	AvatarURL       string    `json:"avatar_url,omitempty"`
	Bio             string    `json:"bio,omitempty"`
	LocationLat     *float64  `json:"location_lat,omitempty"`
	LocationLng     *float64  `json:"location_lng,omitempty"`
	LocationAddress string    `json:"location_address,omitempty"`
	SuccessScore    int       `json:"success_score"`
	BooksShared     int       `json:"books_shared"`
	BooksReceived   int       `json:"books_received"`
	ReviewsReceived int       `json:"reviews_received"`
	IdeasPosted     int       `json:"ideas_posted"`
	TotalUpvotes    int       `json:"total_upvotes"`
	TotalDownvotes  int       `json:"total_downvotes"`
	IsDonor         bool      `json:"is_donor"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type UserRole string

const (
	RoleAdmin  UserRole = "admin"
	RoleMember UserRole = "member"
)

type UserInterest struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Interest  string    `json:"interest"`
	Weight    float64   `json:"weight"`
	CreatedAt time.Time `json:"created_at"`
}

type SuccessScoreHistory struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	ChangeAmount  int       `json:"change_amount"`
	Reason        string    `json:"reason"`
	ReferenceType string    `json:"reference_type,omitempty"`
	ReferenceID   *string   `json:"reference_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}
