package service

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/config"
	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/domain/invoice"
	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/domain/mutation"
	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/utils/id"
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

		for index := range item.Items {
			description += fmt.Sprintf("%s\t%s - %s\n", item.Items[index].Description, item.Items[index].PeriodStart.Format("02 Jan. '06"), item.Items[index].PeriodEnd.Format("02 Jan. '06"))
		}

		// Because the description length has a maximum of 200 we need to check if this is allowed
		// IDs are 27 long, to be sure we use 35
		eBoekhoudenDescriptionMaxLength := 200
		idLength := 35
		maxDescriptionBeforeID := eBoekhoudenDescriptionMaxLength - idLength

		if len(description) >= maxDescriptionBeforeID {
			description = description[:maxDescriptionBeforeID-3]
			description += "...\n"
		}

		// Append the Stripe ID so we can find the invoice in stripe
		description = fmt.Sprintf("%sID: %s", description, item.StripeID)
	}

	var err error
	itemMutation := &mutation.Service{
		BoekhoudenCustomerID: item.BoekhoudenCustomerID,
		Type:                 mutation.InvoiceSend,
		Date:                 item.CreatedAt,
		LedgerAccountCode:    service.cfg.LedgerAccountCode.Debtors,
		InvoiceNumber:        item.Number,
		InvoiceURL:           item.InvoiceURL,
		PaymentTerm:          strconv.Itoa(int(time.Until(item.DueDate).Hours() / 24)),
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
		Type:                 mutation.Memorial,
		Date:                 item.CreatedAt,
		LedgerAccountCode:    "1301",
		InvoiceNumber:        item.Number,
		InvoiceURL:           item.InvoiceURL,
		Description:          fmt.Sprintln("Customer paid, lets move the money to the stripe debtor until stripe paid us"),
		PaymentTerm:          strconv.Itoa(int(time.Until(item.DueDate).Hours() / 24)),
		PaymentFeature:       "", // !?
		Items:                make([]mutation.ServiceItem, 0),
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

	amountF := float64(amount+amountVAT) / 100

	itemMutation.Items = append(itemMutation.Items, mutation.ServiceItem{
		Amount:            amountF,
		LedgerAccountCode: "1300",
	})

	if err := service.mutation.Create(ctx, itemMutation); err != nil {
		return err
	}
	return nil
}

func (service *invoiceService) getLedgerAccountCode(stripeProductID string, stripePlanID string) string {
	ledgerAccountCode := service.cfg.LedgerAccountCode.Default

	if service.cfg.UseLedgerAccountCode.ForAll {
		ledgerAccountCode = service.cfg.LedgerAccountCode.Default
	} else if service.cfg.UseLedgerAccountCode.PerProduct {
		for ledgerCode, productID := range service.cfg.LedgerAccountCode.Products {
			if productID == stripeProductID {
				ledgerAccountCode = ledgerCode
			}
		}
	} else if service.cfg.UseLedgerAccountCode.PerPlan {
		for ledgerCode, planID := range service.cfg.LedgerAccountCode.Plans {
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
