package service

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/aceworksdev/go-stripe-eboekhouden/internal/config"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/invoice"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/mutation"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/utils/id"
)

type invoiceService struct {
	mutation mutation.Servicer
	cfg      *config.EBoekHouden
}

func NewInvoice(c *config.EBoekHouden) (invoice.Servicer, error) {
	return &invoiceService{
		cfg: c,
	}, nil
}

func (service *invoiceService) AddMutation(item mutation.Servicer) {
	service.mutation = item
}

func (service *invoiceService) Finalize(ctx context.Context, item *invoice.Service) error {
	var ledgerAccountCode string
	var description string
	if len(item.Items) > 0 {
		ledgerAccountCode = service.getLedgerAccountCode(item.Items[0].StripeProductID, item.Items[0].StripePlanID)
		description = fmt.Sprintf("%s %s - %s\n%s", item.Items[0].Description, "28 feb.", "28 mrt.", item.InvoiceURL)
	}

	var err error
	itemMutation := &mutation.Service{
		BoekhoudenCustomerID: item.BoekhoudenCustomerID,
		Type:                 mutation.InvoiceSend,
		Date:                 item.CreatedAt,
		LedgerAccountCode:    "1300",
		InvoiceNumber:        item.Number,
		InvoiceURL:           item.InvoiceURL,
		PaymentTerm:          strconv.Itoa(int(item.DueDate.Sub(time.Now()).Hours() / 24)),
		PaymentFeature:       "", // !?
		Items:                make([]mutation.ServiceItem, 0),
		Description:          description,
	}

	if itemMutation.ID, err = id.New(); err != nil {
		return err
	}

	var amountVAT int64 = 0
	var amount int64 = 0
	for _, v := range item.Items {
		for _, v := range v.TaxAmounts {
			amountVAT += v.Amount
		}

		amount += v.Amount * v.Quantity
	}

	amountF := float64(amount) / 100
	amountVATF := float64(amountVAT) / 100
	vatPercentage := amountVATF / amountF * 100
	if math.IsNaN(vatPercentage) {
		vatPercentage = 0
	}

	itemMutation.Items = append(itemMutation.Items, mutation.ServiceItem{
		Amount:            amountF,
		AmountExVAT:       amountF,
		AmountVAT:         amountVATF,
		AmountInVAT:       amountF + amountVATF,
		VATCode:           mutation.VATHighSales21,
		VATPercentage:     vatPercentage,
		LedgerAccountCode: ledgerAccountCode,
	})

	if err := service.mutation.Create(ctx, itemMutation); err != nil {
		return err
	}
	return nil
}

func (service *invoiceService) Paid(ctx context.Context, item *invoice.Service) error {
	var err error
	itemMutation := &mutation.Service{
		BoekhoudenCustomerID: item.BoekhoudenCustomerID,
		Type:                 mutation.InvoicePaymentReceived,
		Date:                 item.CreatedAt,
		LedgerAccountCode:    "1010",
		InvoiceNumber:        item.Number,
		InvoiceURL:           item.InvoiceURL,
		PaymentTerm:          strconv.Itoa(int(item.DueDate.Sub(time.Now()).Hours() / 24)),
		PaymentFeature:       "", // !?
		Items:                make([]mutation.ServiceItem, 0),
		Description:          fmt.Sprintf("Professional 28 feb. = 28 mrt. 2022\n%s", item.InvoiceURL),
	}

	if itemMutation.ID, err = id.New(); err != nil {
		return err
	}

	var amountVAT int64 = 0
	var amount int64 = 0
	for _, v := range item.Items {
		for _, v := range v.TaxAmounts {
			amountVAT += v.Amount
		}

		amount += v.Amount * v.Quantity
	}

	amountF := float64(amount) / 100
	amountVATF := float64(amountVAT) / 100
	vatPercentage := amountVATF / amountF * 100
	if math.IsNaN(vatPercentage) {
		vatPercentage = 0
	}

	itemMutation.Items = append(itemMutation.Items, mutation.ServiceItem{
		Amount:            amountF,
		AmountExVAT:       amountF,
		AmountVAT:         amountVATF,
		AmountInVAT:       amountF + amountVATF,
		VATCode:           mutation.VATHighSales21,
		VATPercentage:     vatPercentage,
		LedgerAccountCode: "1300",
	})

	fmt.Printf("itemMutation: %+v\n", itemMutation)

	if err := service.mutation.Create(ctx, itemMutation); err != nil {
		return err
	}
	return nil
}

func (service *invoiceService) getLedgerAccountCode(stripeProductID string, stripePlanID string) string {
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
