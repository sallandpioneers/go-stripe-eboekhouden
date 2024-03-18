package service

import (
	"context"

	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/config"
	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/domain/mutation"
	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/domain/payout"
	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/utils/id"
)

type payoutService struct {
	mutation mutation.Servicer
	cfg      *config.EBoekHouden
}

func NewPayout(cfg *config.EBoekHouden) (payout.Servicer, error) {
	return &payoutService{
		cfg: cfg,
	}, nil
}

func (service *payoutService) AddMutation(item mutation.Servicer) {
	service.mutation = item
}

func (service *payoutService) Paid(ctx context.Context, item *payout.Service) error {
	var err error
	itemMutation := &mutation.Service{
		BoekhoudenCustomerID: item.BoekhoudenCustomerID,
		Type:                 mutation.MoneyReceived,
		Date:                 item.CreatedAt,
		LedgerAccountCode:    service.cfg.LedgerAccountCode.Bank,
		PaymentFeature:       "", // !?
		Items: []mutation.ServiceItem{
			{
				Amount:            float64(item.Amount) / 100,
				AmountExVAT:       float64(item.Amount) / 100,
				AmountVAT:         0,
				AmountInVAT:       float64(item.Amount) / 100,
				VATCode:           mutation.VATNo,
				VATPercentage:     0,
				LedgerAccountCode: "1301",
			},
		},
		Description: "STRIPE PAYOUT",
	}

	if itemMutation.ID, err = id.New(); err != nil {
		return err
	}

	if err := service.mutation.Create(ctx, itemMutation); err != nil {
		return err
	}
	return nil
}
