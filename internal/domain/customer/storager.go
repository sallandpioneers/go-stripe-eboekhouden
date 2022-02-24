package customer

import (
	"context"

	"github.com/oklog/ulid"
)

type Storager interface {
	Create(context.Context, *Service) error
	Get(context.Context, ulid.ULID) (*Service, error)
	GetBasedOnStripeID(context.Context, string) (*Service, error)
	GetBasedOnBoekhoudenID(context.Context, int64) (*Service, error)
}
