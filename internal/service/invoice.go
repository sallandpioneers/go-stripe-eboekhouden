package service

import (
	"context"

	"github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/invoice"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/mutation"
)

type invoiceService struct {
	push     invoice.Pusher
	mutation mutation.Servicer
}

func NewInvoice(p invoice.Pusher) (invoice.Servicer, error) {
	return &invoiceService{
		push: p,
	}, nil
}

func (service *invoiceService) AddMutation(item mutation.Servicer) {
	service.mutation = item
}

func (service *invoiceService) Create(ctx context.Context, item *invoice.Service, customerCode string) error {
	itemMutation := &mutation.Service{}
	if err := service.mutation.Create(ctx, itemMutation, customerCode); err != nil {
		return err
	}
	return nil
}

func (service *invoiceService) Update(ctx context.Context, item *invoice.Service, customerCode string) error {
	itemMutation := &mutation.Service{}
	if err := service.mutation.Update(ctx, itemMutation, customerCode); err != nil {
		return err
	}
	return nil
}
