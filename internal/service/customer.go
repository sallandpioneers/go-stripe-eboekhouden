package service

import (
	"context"

	"github.com/aceworksdev/go-stripe-eboekhouden/internal"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/customer"
)

type customerService struct {
	storage customer.Storager
	push    customer.Pusher
}

func NewCustomer(s customer.Storager, p customer.Pusher) (customer.Servicer, error) {
	return &customerService{
		storage: s,
		push:    p,
	}, nil
}

func (service *customerService) Create(ctx context.Context, item *customer.Service) error {
	var err error
	item.Code, err = internal.GetRandomChars(15)
	if err != nil {
		return err
	}

	if err := service.push.Create(ctx, item); err != nil {
		return err
	}
	return nil
}
