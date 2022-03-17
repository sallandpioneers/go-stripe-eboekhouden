package service

import (
	"context"
	"strings"

	"github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/customer"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/utils/id"
	"github.com/oklog/ulid"
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

func (service *customerService) Get(ctx context.Context, id ulid.ULID) (*customer.Service, error) {
	return service.storage.Get(ctx, id)
}

func (service *customerService) GetBasedOnStripeID(ctx context.Context, id string) (*customer.Service, error) {
	return service.storage.GetBasedOnStripeID(ctx, id)
}

func (service *customerService) GetBasedOnBoekhoudenID(ctx context.Context, id int64) (*customer.Service, error) {
	return service.storage.GetBasedOnBoekhoudenID(ctx, id)
}

func (service *customerService) Create(ctx context.Context, item *customer.Service) error {
	var err error
	item.ID, err = id.New()
	if err != nil {
		return err
	}

	item.BoekhoudenID = strings.TrimPrefix(item.StripeID, "cus_")
	if len(item.BoekhoudenID) > 15 {
		item.BoekhoudenID = item.BoekhoudenID[len(item.BoekhoudenID)-15:]
	}
	if item.Name == "" {
		item.Name = item.Email
	}
	item.Type = customer.Business

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
	item.BoekhoudenID = itemOld.BoekhoudenID

	if item.Name == "" {
		item.Name = item.Email
	}
	item.Type = customer.Business

	if err := service.push.Update(ctx, item); err != nil {
		return err
	}
	return nil
}
