package domain

import "time"

type Donation struct {
	ID           string       `json:"id"`
	DonorID      string       `json:"donor_id"`
	Donor        *User        `json:"donor,omitempty"`
	DonationType DonationType `json:"donation_type"`
	BookID       *string      `json:"book_id,omitempty"`
	Book         *Book        `json:"book,omitempty"`
	Amount       *float64     `json:"amount,omitempty"`
	Currency     string       `json:"currency,omitempty"`
	Message      string       `json:"message,omitempty"`
	IsPublic     bool         `json:"is_public"`
	CreatedAt    time.Time    `json:"created_at"`
}

type DonationType string

const (
	DonationTypeBook  DonationType = "book"
	DonationTypeMoney DonationType = "money"
)
