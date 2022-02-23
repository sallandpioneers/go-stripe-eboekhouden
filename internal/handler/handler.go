package handler

import (
	"github.com/aceworksdev/go-stripe-eboekhouden/internal"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/config"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/handler/json"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/server/domain/handler"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/server/domain/service"
)

func New(mode string, h *handler.Handler, s *service.Service, sa *internal.ServicesAvailable, c *config.Config) error {
	switch mode {
	case "json":
		if err := json.New(h, s, c, sa); err != nil {
			return err
		}
	default:
		return internal.ModeUnknown("handler", mode, "json")
	}

	return nil
}
