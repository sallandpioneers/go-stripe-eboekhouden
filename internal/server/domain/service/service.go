package service

import (
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/customer"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/hooks"
)

type Service struct {
	Hooks    hooks.Servicer
	Customer customer.Servicer
}
