package invoice

import (
	"time"

	"github.com/oklog/ulid"
)

type DirectDebitType string

const (
	SingleUse  DirectDebitType = "E"
	Continuous DirectDebitType = "D"
)

type Service struct {
	ID               ulid.ULID
	StripeID         string
	CustomerID       ulid.ULID
	StripeCustomerID string
	MutationID       int64
	Number           string
	DueDate          time.Time
	CollectionMethod string
	Items            []ItemService
	Subtotal         int64
	Tax              int64
	Total            int64
	AmountDue        int64
	AmountPaid       int64
	AmountRemaining  int64
	CreatedAt        time.Time
}

type ItemService struct {
	ID              ulid.ULID
	StripeID        string
	StripePlanID    string
	StripeProductID string
	Quantity        int64
	Description     string
	Amount          int64 // In Cents
	BTWCode         string
	// Unit              string
	// Code              string
}
