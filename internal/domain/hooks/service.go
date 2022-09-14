package hooks

import (
	"context"

	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/domain/customer"
	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/domain/invoice"
	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/domain/payout"
)

type Servicer interface {
	AddCustomer(service customer.Servicer)
	AddInvoice(item invoice.Servicer)
	AddPayout(item payout.Servicer)

	InvoiceFinalized(context.Context, *invoice.Service) error
	InvoicePaid(context.Context, *invoice.Service) error

	CustomerCreate(context.Context, *customer.Service) error
	CustomerUpdate(context.Context, *customer.Service, *customer.BalanceUpdate) error

	PayoutPaid(context.Context, *payout.Service) error
}
