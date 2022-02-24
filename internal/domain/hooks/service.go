package hooks

import (
	"context"

	"github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/customer"
)

type Servicer interface {
	AddCustomer(service customer.Servicer)

	CustomerCreate(context.Context, *customer.Service) error
	CustomerUpdate(context.Context, *customer.Service) error
}
