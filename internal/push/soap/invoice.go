package soap

import (
	"context"
	"errors"

	"github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/invoice"
	eboekhouden "github.com/aceworksdev/go-stripe-eboekhouden/internal/push/soap/generated"
)

type invoicePush struct {
	client eboekhouden.SoapAppSoap
}

func NewInvoice(client eboekhouden.SoapAppSoap) (invoice.Pusher, error) {
	return &invoicePush{
		client: client,
	}, nil
}

func (service *invoicePush) Create(ctx context.Context, item *invoice.Service, customerCode string) error {
	return errors.New("not_implemented")
}

func (service *invoicePush) Update(ctx context.Context, item *invoice.Service, customerCode string) error {
	return errors.New("not_implemented")
}
