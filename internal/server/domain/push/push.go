package push

import (
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/customer"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/mutation"
)

type Push struct {
	Soap *Soap
}

type Soap struct {
	Customer customer.Pusher
	Mutation mutation.Pusher
}
