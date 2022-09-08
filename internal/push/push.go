package push

import (
	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/config"
	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/push/soap"
	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/server/domain/push"
)

func New(em *push.Push, c *config.EBoekHouden, isDevelopment bool, sendToExternal bool) error {
	if err := soap.New(em.Soap, c); err != nil {
		return err
	}
	return nil
}
