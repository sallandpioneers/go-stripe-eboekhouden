package service

import (
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/config"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/server/domain/push"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/server/domain/service"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/server/domain/storage"
	"github.com/valyala/fasthttp"
)

func New(s *service.Service, db *storage.Storage, p *push.Push, c *fasthttp.Client, config *config.Config) error {
	var err error
	if s.Customer, err = NewCustomer(db.Customer, p.Soap.Customer); err != nil {
		return err
	}

	return nil
}
