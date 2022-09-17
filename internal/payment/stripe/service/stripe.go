package service

import (
	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/config"
	"github.com/stripe/stripe-go/v73"
)

type Service struct {
}

func New(cs *config.Stripe) (*Service, error) {
	stripe.Key = cs.Key
	return &Service{}, nil
}
