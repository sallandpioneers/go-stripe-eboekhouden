package storage

import "github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/customer"

type Storage struct {
	Customer customer.Storager
}
