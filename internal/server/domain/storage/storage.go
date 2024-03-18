package storage

import "github.com/sallandpioneers/go-stripe-eboekhouden/internal/domain/customer"

type Storage struct {
	Customer customer.Storager
}
