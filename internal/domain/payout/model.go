package payout

import (
	"time"

	"github.com/oklog/ulid"
)

type Service struct {
	ID                   ulid.ULID
	StripeID             string
	BoekhoudenCustomerID string
	Amount               int64
	Description          string
	CreatedAt            time.Time
}
