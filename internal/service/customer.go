package service

import (
	"context"
	"strings"
	"time"

	"github.com/oklog/ulid"
	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/domain/customer"
	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/domain/mutation"
	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/utils/id"
)

type customerService struct {
	storage customer.Storager
	push    customer.Pusher

	mutation mutation.Servicer
}

func NewCustomer(s customer.Storager, p customer.Pusher) (customer.Servicer, error) {
	return &customerService{
		storage: s,
		push:    p,
	}, nil
}

func (service *customerService) AddMutation(item mutation.Servicer) {
	service.mutation = item
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

func (service *customerService) Update(ctx context.Context, item *customer.Service, balance *customer.BalanceUpdate) error {
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

	// If the balance is negative we own the customer money
	// If the balance is positive the customer owns use money
	if balance.HasChanged {
		balanceChange := balance.NewBalance - balance.OldBalance
		balanceChange *= -1 // Make the balance positive

		var err error
		itemMutation := &mutation.Service{
			BoekhoudenCustomerID: item.BoekhoudenID,
			Type:                 mutation.Memorial,
			Date:                 time.Now(),
			LedgerAccountCode:    "1301",
			PaymentFeature:       item.BoekhoudenID, // !?
			Items: []mutation.ServiceItem{
				{
					Amount:            float64(balanceChange) / 100,
					AmountExVAT:       float64(balanceChange) / 100,
					AmountVAT:         0,
					AmountInVAT:       float64(balanceChange) / 100,
					VATCode:           mutation.VATNo,
					VATPercentage:     0,
					LedgerAccountCode: "1302",
				},
			},
			Description: "Customer Balance Update",
		}

		if itemMutation.ID, err = id.New(); err != nil {
			return err
		}

		if err := service.mutation.Create(ctx, itemMutation); err != nil {
			return err
		}

	}

	return nil
}
