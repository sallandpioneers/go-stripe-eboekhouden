package service

import (
	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/config"
	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/payment"
	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/server/domain/push"
	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/server/domain/service"
	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/server/domain/storage"
	"github.com/valyala/fasthttp"
)

func New(s *service.Service, db *storage.Storage, paym payment.Service, p *push.Push, c *fasthttp.Client, cfg *config.Config) error {
	var err error
	if s.Customer, err = NewCustomer(db.Customer, p.Soap.Customer); err != nil {
		return err
	}
	if s.Invoice, err = NewInvoice(cfg.EBoekHouden); err != nil {
		return err
	}
	if s.Mutation, err = NewMutation(p.Soap.Mutation); err != nil {
		return err
	}
	if s.Payout, err = NewPayout(cfg.EBoekHouden); err != nil {
		return err
	}
	if s.Hooks, err = NewHooks(); err != nil {
		return err
	}
	if s.Report, err = NewReport(paym); err != nil {
		return err
	}

	s.Hooks.AddCustomer(s.Customer)
	s.Hooks.AddInvoice(s.Invoice)
	s.Hooks.AddPayout(s.Payout)

	s.Customer.AddMutation(s.Mutation)
	s.Invoice.AddMutation(s.Mutation)
	s.Payout.AddMutation(s.Mutation)

	return nil
}
