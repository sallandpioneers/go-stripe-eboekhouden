package service

import (
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/customer"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/hooks"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/invoice"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/mutation"
)

type Service struct {
	Hooks    hooks.Servicer
	Customer customer.Servicer
	Invoice  invoice.Servicer
	Mutation mutation.Servicer
}
