package hooks

import (
	"context"

	"github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/customer"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/invoice"
)

type Servicer interface {
	AddCustomer(service customer.Servicer)
	AddInvoice(item invoice.Servicer)

	InvoiceFinalized(context.Context, *invoice.Service) error
	InvoicePaid(context.Context, *invoice.Service) error

	CustomerCreate(context.Context, *customer.Service) error
	CustomerUpdate(context.Context, *customer.Service) error
}
