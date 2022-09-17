package payment

import (
	"github.com/sallandpioneers/go-stripe-eboekhouden/internal"
	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/config"
	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/model"
	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/payment/stripe/service"
)

type Service interface {
	CreateReport(model.Report) error
}

func New(c *config.Payment) (p Service, err error) {
	switch c.Current {
	case "stripe":
		if p, err = service.New(c.Stripe); err != nil {
			return nil, err
		}
	default:
		return nil, internal.ModeUnknown("payment", c.Current, "stripe")
	}
	return p, nil
}
