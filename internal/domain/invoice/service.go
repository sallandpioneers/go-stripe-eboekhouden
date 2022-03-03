package invoice

import (
	"context"

	"github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/mutation"
)

type Servicer interface {
	AddMutation(item mutation.Servicer)

	Create(ctx context.Context, item *Service, customerCode string) error
	Update(ctx context.Context, item *Service, customerCode string) error
}
