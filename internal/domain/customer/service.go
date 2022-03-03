package customer

import (
	"context"

	"github.com/oklog/ulid"
)

type Servicer interface {
	Get(context.Context, ulid.ULID) (*Service, error)
	GetBasedOnStripeID(context.Context, string) (*Service, error)
	GetBasedOnBoekhoudenID(context.Context, int64) (*Service, error)

	Create(context.Context, *Service) error
	Update(ctx context.Context, item *Service) error
}
