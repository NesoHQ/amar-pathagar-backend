package domain

import "time"

type ReadingIdea struct {
	ID        string    `json:"id"`
	BookID    string    `json:"book_id"`
	Book      *Book     `json:"book,omitempty"`
	UserID    string    `json:"user_id"`
	User      *User     `json:"user,omitempty"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Upvotes   int       `json:"upvotes"`
	Downvotes int       `json:"downvotes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type IdeaVote struct {
	ID        string    `json:"id"`
	IdeaID    string    `json:"idea_id"`
	UserID    string    `json:"user_id"`
	VoteType  VoteType  `json:"vote_type"`
	CreatedAt time.Time `json:"created_at"`
}

type VoteType string

const (
	VoteTypeUp   VoteType = "upvote"
	VoteTypeDown VoteType = "downvote"
)
