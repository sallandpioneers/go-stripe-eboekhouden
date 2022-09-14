package payout

import (
	"context"

	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/domain/mutation"
)

type Servicer interface {
	AddMutation(item mutation.Servicer)
	Paid(context.Context, *Service) error
}
