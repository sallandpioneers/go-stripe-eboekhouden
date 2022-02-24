package service

import (
	"context"

	"github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/customer"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/hooks"
)

type hooksService struct {
	customer customer.Servicer
}

func NewHooks() (hooks.Servicer, error) {
	return &hooksService{}, nil
}

func (service *hooksService) AddCustomer(customer customer.Servicer) {
	service.customer = customer
}

func (service *hooksService) CustomerCreate(ctx context.Context, item *customer.Service) error {
	if err := service.customer.Create(ctx, item); err != nil {
		return err
	}
	return nil
}
func (service *hooksService) CustomerUpdate(ctx context.Context, item *customer.Service) error {
	if err := service.customer.Update(ctx, item); err != nil {
		return err
	}
	return nil
}
