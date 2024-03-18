package invoice

import (
	"context"

	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/domain/mutation"
)

type Servicer interface {
	AddMutation(item mutation.Servicer)

	Finalize(ctx context.Context, item *Service) error
	Paid(ctx context.Context, item *Service) error
}
