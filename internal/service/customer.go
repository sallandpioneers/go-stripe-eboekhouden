package service

import (
	"context"
	"strings"

	"github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/customer"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/utils/id"
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
	item.ID, err = id.New()
	if err != nil {
		return err
	}

	item.Code = strings.TrimPrefix(item.StripeID, "cus_")
	if len(item.Code) > 15 {
		item.Code = item.Code[len(item.Code)-15:]
	}
	if item.Name != "" {
		item.Company = item.Name
	} else if item.Email != "" {
		item.Company = item.Email
	}

	if err := service.push.Create(ctx, item); err != nil {
		return err
	}

	if err := service.storage.Create(ctx, item); err != nil {
		return err
	}

	return nil
}

func (service *customerService) Update(ctx context.Context, item *customer.Service) error {
	itemOld, err := service.storage.GetBasedOnStripeID(ctx, item.StripeID)
	if err != nil {
		return err
	}
	item.Code = itemOld.Code
	if err := service.push.Update(ctx, item); err != nil {
		return err
	}
	return nil
}
