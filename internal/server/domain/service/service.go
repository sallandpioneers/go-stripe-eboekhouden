package service

import "github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/customer"

type Service struct {
	Customer customer.Servicer
}
