package hooks

import (
	"context"

	"github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/customer"
)

type Servicer interface {
	CustomerCreate(context.Context, *customer.Service) error
}
