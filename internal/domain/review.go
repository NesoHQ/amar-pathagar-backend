package domain

import "time"

type UserReview struct {
	ID                  string    `json:"id"`
	ReviewerID          string    `json:"reviewer_id"`
	Reviewer            *User     `json:"reviewer,omitempty"`
	RevieweeID          string    `json:"reviewee_id"`
	Reviewee            *User     `json:"reviewee,omitempty"`
	BookID              *string   `json:"book_id,omitempty"`
	BehaviorRating      *int      `json:"behavior_rating,omitempty"`
	BookConditionRating *int      `json:"book_condition_rating,omitempty"`
	CommunicationRating *int      `json:"communication_rating,omitempty"`
	Comment             string    `json:"comment,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
}
