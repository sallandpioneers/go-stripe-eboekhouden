package push

import "github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/customer"

type Push struct {
	Soap *Soap
}

type Soap struct {
	Customer customer.Pusher
}
