package json

import (
	"github.com/aceworksdev/go-stripe-eboekhouden/internal"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/config"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/server/domain/handler"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/server/domain/service"

	jsoniter "github.com/json-iterator/go"
)

func New(h *handler.Handler, s *service.Service, c *config.Config, sa *internal.ServicesAvailable) error {
	jip, jsp := NewPool()
	var err error
	if h.Hooks, err = NewHooks(s.Hooks, c.Stripe.Secret, jip, jsp); err != nil {
		return err
	}
	return nil
}

func NewPool() (jsonIteratorPool jsoniter.IteratorPool, jsonStreamPool jsoniter.StreamPool) {
	jsonIteratorPool = jsoniter.NewIterator(jsoniter.ConfigFastest).Pool()

	// TODO not sure about the 100 buffer size, will lookup in the future
	jsonStreamPool = jsoniter.NewStream(jsoniter.ConfigFastest, nil, 100).Pool()
	return jsonIteratorPool, jsonStreamPool
}
