package json

import (
	"context"
	"log"
	"time"

	"github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/customer"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/hooks"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/invoice"

	jsoniter "github.com/json-iterator/go"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/plan"
	"github.com/stripe/stripe-go/v72/transfer"
	"github.com/stripe/stripe-go/v72/webhook"
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
		err = h.InvoiceFinalized(ev.Data.Raw)
	case "invoice.marked_uncollectible":
		// End the subscription, this will be done by stripe, so dont do anything
		// Other data will be send through invoice.updated, so we dont have to worry about that
		err = nil
	case "invoice.paid":
		// Send signal to RPI with noise that we got money
		// Other data will be send through invoice.updated, so we dont have to worry about that
		err = h.invoicePaid(ev.Data.Raw)
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
		err = h.customerCreate(ev.Data.Raw)
	case "customer.updated":
		err = h.customerUpdate(ev.Data.Raw)
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
		err = h.payoutPaid(ev.Data.Raw)
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

func (h *hooksHandler) payoutPaid(data []byte) error {
	var c stripe.Payout
	jsonIterator := h.jsonIteratorPool.BorrowIterator(data)
	defer h.jsonIteratorPool.ReturnIterator(jsonIterator)
	jsonIterator.ReadVal(&c)
	if jsonIterator.Error != nil {
		return jsonIterator.Error
	}

	t, _ := transfer.Get(c.ID, nil)
	log.Println(t)

	t, _ = transfer.Get("po_1Ke7iCBrr18AunNBFtUELFcS/transactions?from_stmt=true&count=20&expand%5B%5D=data.charge", &stripe.TransferParams{
		Params: stripe.Params{},
	})
	log.Println(t)

	return nil
}

func (h *hooksHandler) InvoiceFinalized(data []byte) error {
	item, err := h.getInvoice(data)
	if err != nil {
		return err
	}

	if err := h.service.InvoiceFinalized(context.Background(), item); err != nil {
		return err
	}

	return nil
}

func (h *hooksHandler) invoicePaid(data []byte) error {
	item, err := h.getInvoice(data)
	if err != nil {
		return err
	}

	if err := h.service.InvoicePaid(context.Background(), item); err != nil {
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
		CollectionMethod: string(*c.CollectionMethod),
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

		for k2, v2 := range v.TaxAmounts {
			item.Items[k].TaxAmounts[k2].Amount = v2.Amount
			item.Items[k].TaxAmounts[k2].Inclusive = v2.Inclusive
		}
	}
	return item, nil
}

func (h *hooksHandler) customerCreate(data []byte) error {
	item, err := h.getCustomer(data)
	if err != nil {
		return err
	}
	if err := h.service.CustomerCreate(context.Background(), item); err != nil {
		return err
	}
	return nil
}

func (h *hooksHandler) customerUpdate(data []byte) error {
	item, err := h.getCustomer(data)
	if err != nil {
		return err
	}
	if err := h.service.CustomerUpdate(context.Background(), item); err != nil {
		return err
	}
	return nil
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
	}

	item.Addresses.Business.Address = c.Address.Line1
	item.Addresses.Business.ZipCode = c.Address.PostalCode
	item.Addresses.Business.City = c.Address.City
	item.Addresses.Business.Country = c.Address.Country

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
