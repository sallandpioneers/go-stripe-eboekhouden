package json

import (
	"context"
	"log"
	"time"

	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/domain/customer"
	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/domain/hooks"
	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/domain/invoice"
	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/domain/payout"

	jsoniter "github.com/json-iterator/go"
	"github.com/stripe/stripe-go/v73"
	"github.com/stripe/stripe-go/v73/plan"
	"github.com/stripe/stripe-go/v73/webhook"
	"github.com/valyala/fasthttp"
)

type hooksHandler struct {
	service          hooks.Servicer
	jsonIteratorPool jsoniter.IteratorPool
	jsonStreamPool   jsoniter.StreamPool
	secret           string
}

func NewHooks(service hooks.Servicer, s string, jip jsoniter.IteratorPool, jsp jsoniter.StreamPool) (hooks.Handler, error) {
	return &hooksHandler{
		service:          service,
		jsonIteratorPool: jip,
		jsonStreamPool:   jsp,
		secret:           s,
	}, nil
}

func (h *hooksHandler) AllHooks(ctx *fasthttp.RequestCtx) {
	ev, err := validateSignature(ctx.PostBody(), string(ctx.Request.Header.Peek("Stripe-Signature")), h.secret)
	if err != nil {
		log.Printf("stripe webhook failure\t %v", err)
		return
	}

	switch ev.Type {
	case "invoice.created":
		err = nil
	case "invoice.deleted":
		// This will happen when a draft invoice has been deleted.
		err = nil
	case "invoice.finalized":
		// invoice is set to open, maybe send email to customer?
		// Other data will be send through invoice.updated, so we dont have to worry about that
		err = h.InvoiceFinalized(ctx, ev.Data.Raw)
	case "invoice.marked_uncollectible":
		// End the subscription, this will be done by stripe, so dont do anything
		// Other data will be send through invoice.updated, so we dont have to worry about that
		err = nil
	case "invoice.paid":
		// Send signal to RPI with noise that we got money
		// Other data will be send through invoice.updated, so we dont have to worry about that
		err = h.invoicePaid(ctx, ev.Data.Raw)
	case "invoice.payment_action_required":
		// User is suppose to do some shit
		// TODO figure out when this can happen
		// TODO create system that will let the user know about this
		err = nil
	case "invoice.payment_failed":
		// Payment failed, dont do anything
		// TODO make logging system that can track this
		err = nil
	case "invoice.payment_succeeded":
		// Payment succeeded
		// Send signal to RPI with noise that we got money
		// Other data will be send through invoice.updated, so we dont have to worry about that
		err = nil
	case "invoice.sent":
		// dont do anything
		err = nil
	case "invoice.upcoming":
		// upcoming invoice, maybe send an email to the customer? upcoming means x amount of days before invoice is due.
		// This can be changed in dashboard: https://dashboard.stripe.com/settings/billing/automatic
		err = nil
	case "invoice.updated":
		// err = h.invoiceUpdate(ev.Data.Raw)
	case "invoice.voided":
		// Invoice cannot be used anymore
		err = nil
	case "customer.subscription.created":
		err = nil
	case "customer.subscription.updated":
		err = nil
	case "customer.subscription.deleted":
		err = nil
	case "customer.created":
		err = h.customerCreate(ctx, ev.Data.Raw)
	case "customer.updated":
		err = h.customerUpdate(ctx, ev.Data.Raw, ev.Data.PreviousAttributes)
	case "customer.tax_id.updated":
		err = nil
	case "payment_intent.succeeded":
		err = nil
	case "setup_intent.requires_action":
		err = nil
	case "setup_intent.succeeded":
		err = nil
	case "setup_intent.setup_failed":
		err = nil
	case "payout.paid":
		err = h.payoutPaid(ctx, ev.Data.Raw)
	default:
		log.Printf("webhook not supported\t\t\t%s", ev.Type)
		err = nil
	}

	if err != nil {
		log.Println(err)
		ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	ctx.Response.SetStatusCode(fasthttp.StatusNoContent)
}

func (h *hooksHandler) payoutPaid(ctx context.Context, data []byte) error {
	item, err := h.getPayout(data)
	if err != nil {
		return err
	}

	if err := h.service.PayoutPaid(ctx, item); err != nil {
		return err
	}
	return nil
}

func (h *hooksHandler) getPayout(data []byte) (*payout.Service, error) {
	var c stripe.Payout
	jsonIterator := h.jsonIteratorPool.BorrowIterator(data)
	defer h.jsonIteratorPool.ReturnIterator(jsonIterator)
	jsonIterator.ReadVal(&c)
	if jsonIterator.Error != nil {
		return nil, jsonIterator.Error
	}

	return &payout.Service{
		StripeID:    c.ID,
		Amount:      c.Amount,
		Description: c.Description,
		CreatedAt:   time.Unix(c.Created, 0),
	}, nil
}

func (h *hooksHandler) InvoiceFinalized(ctx context.Context, data []byte) error {
	item, err := h.getInvoice(data)
	if err != nil {
		return err
	}

	if err := h.service.InvoiceFinalized(ctx, item); err != nil {
		return err
	}

	return nil
}

func (h *hooksHandler) invoicePaid(ctx context.Context, data []byte) error {
	item, err := h.getInvoice(data)
	if err != nil {
		return err
	}

	if err := h.service.InvoicePaid(ctx, item); err != nil {
		return err
	}

	return nil
}

func (h *hooksHandler) getInvoice(data []byte) (*invoice.Service, error) {
	var c stripe.Invoice
	jsonIterator := h.jsonIteratorPool.BorrowIterator(data)
	defer h.jsonIteratorPool.ReturnIterator(jsonIterator)
	jsonIterator.ReadVal(&c)
	if jsonIterator.Error != nil {
		return nil, jsonIterator.Error
	}

	item := &invoice.Service{
		StripeID:         c.ID,
		StripeCustomerID: c.Customer.ID,
		Number:           c.Number,
		DueDate:          time.Unix(c.DueDate, 0),
		CollectionMethod: string(c.CollectionMethod),
		Items:            make([]invoice.ItemService, c.Lines.TotalCount),
		Subtotal:         c.Subtotal,
		Tax:              c.Tax,
		Total:            c.Total,
		AmountDue:        c.AmountDue,
		AmountPaid:       c.AmountPaid,
		AmountRemaining:  c.AmountRemaining,
		InvoiceURL:       c.InvoicePDF,
		CreatedAt:        time.Unix(c.Created, 0),
	}

	for k, v := range c.Lines.Data {
		var planID string
		var productID string
		if v.Plan == nil {
			productID = v.Price.Product.ID
		} else {
			plan, err := plan.Get(v.Plan.ID, nil)
			if err != nil {
				return nil, err
			}

			productID = plan.Product.ID
		}

		item.Items[k].StripeID = v.ID
		if v.Plan != nil {
			item.Items[k].StripePlanID = planID
		}
		item.Items[k].StripeProductID = productID
		item.Items[k].Quantity = v.Quantity
		item.Items[k].Description = v.Description
		item.Items[k].Amount = v.Amount
		item.Items[k].TaxAmounts = make([]invoice.InvoiceTaxAmountService, len(v.TaxAmounts))
		item.Items[k].PeriodStart = time.Unix(v.Period.Start, 0)
		item.Items[k].PeriodEnd = time.Unix(v.Period.End, 0)

		for k2, v2 := range v.TaxAmounts {
			item.Items[k].TaxAmounts[k2].Amount = v2.Amount
			item.Items[k].TaxAmounts[k2].Inclusive = v2.Inclusive
		}
	}
	return item, nil
}

func (h *hooksHandler) customerCreate(ctx context.Context, data []byte) error {
	item, err := h.getCustomer(data)
	if err != nil {
		return err
	}
	if err := h.service.CustomerCreate(ctx, item); err != nil {
		return err
	}
	return nil
}

func (h *hooksHandler) customerUpdate(ctx context.Context, data []byte, previousAttributes map[string]interface{}) error {
	item, err := h.getCustomer(data)
	if err != nil {
		return err
	}
	balance, err := h.getCustomerPreviousAttributes(ctx, previousAttributes)
	if err != nil {
		return err
	}
	balance.NewBalance = item.Balance

	if err := h.service.CustomerUpdate(ctx, item, balance); err != nil {
		return err
	}
	return nil
}

func (h *hooksHandler) getCustomerPreviousAttributes(ctx context.Context, data map[string]interface{}) (*customer.BalanceUpdate, error) {
	item := &customer.BalanceUpdate{}
	if balance, ok := data["balance"]; ok {
		item.OldBalance = int64(balance.(float64))
		item.HasChanged = true
	}
	return item, nil
}

func (h *hooksHandler) getCustomer(data []byte) (*customer.Service, error) {
	var c stripe.Customer
	jsonIterator := h.jsonIteratorPool.BorrowIterator(data)
	defer h.jsonIteratorPool.ReturnIterator(jsonIterator)
	jsonIterator.ReadVal(&c)
	if jsonIterator.Error != nil {
		return nil, jsonIterator.Error
	}

	item := &customer.Service{
		StripeID: c.ID,
		Name:     c.Name,
		Email:    c.Email,
		Balance:  c.Balance,
	}

	if c.Address != nil {
		item.Addresses.Business.Address = c.Address.Line1
		item.Addresses.Business.ZipCode = c.Address.PostalCode
		item.Addresses.Business.City = c.Address.City
		item.Addresses.Business.Country = c.Address.Country
	}

	if c.Shipping != nil {
		item.Addresses.Mailing.Address = c.Shipping.Address.Line1
		item.Addresses.Mailing.ZipCode = c.Shipping.Address.PostalCode
		item.Addresses.Mailing.City = c.Shipping.Address.City
		item.Addresses.Mailing.Country = c.Shipping.Address.Country
	}
	return item, nil
}

func validateSignature(payload []byte, header, secret string) (stripe.Event, error) {
	return webhook.ConstructEvent(payload, header, secret)
}
