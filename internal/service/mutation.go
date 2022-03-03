package service

import (
	"context"
	"log"

	"github.com/aceworksdev/go-stripe-eboekhouden/internal/config"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/mutation"
)

type mutationService struct {
	push mutation.Pusher
	cfg  *config.EBoekHouden
}

func NewMutation(p mutation.Pusher, c *config.EBoekHouden) (mutation.Servicer, error) {
	return &mutationService{
		push: p,
		cfg:  c,
	}, nil
}

func (service *mutationService) Create(ctx context.Context, item *mutation.Service, customerCode string) error {
	item.LedgerAccountCode = service.GetLedgerAccountCode("", "")

	if err := service.push.Create(ctx, item, customerCode); err != nil {
		return err
	}
	return nil
}

func (service *mutationService) Update(ctx context.Context, item *mutation.Service, customerCode string) error {
	if err := service.push.Update(ctx, item, customerCode); err != nil {
		return err
	}
	return nil
}

func (service *mutationService) GetLedgerAccountCode(stripeProductID string, stripePlanID string) string {
	ledgerAccountCode := service.cfg.LedgerAccountCodeDefault

	if service.cfg.UseLedgerAccountCodeForAll {
		ledgerAccountCode = service.cfg.LedgerAccountCodeDefault
	} else if service.cfg.UseLedgerAccountCodePerProduct {
		for ledgerCode, productID := range service.cfg.LedgerAccountCodeProducts {
			if productID == stripeProductID {
				ledgerAccountCode = ledgerCode
			}
		}
	} else if service.cfg.UseLedgerAccountCodePerPlan {
		for ledgerCode, planID := range service.cfg.LedgerAccountCodePlans {
			if planID == stripePlanID {
				ledgerAccountCode = ledgerCode
			}
		}
	}

	if ledgerAccountCode == "" {
		log.Fatal("EBoekhouden config not setup correctly, missing default ledger account code")
	}

	return ledgerAccountCode
}
