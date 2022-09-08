package push

import (
	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/domain/customer"
	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/domain/mutation"
)

type Push struct {
	Soap *Soap
}

type Soap struct {
	Customer customer.Pusher
	Mutation mutation.Pusher
}
