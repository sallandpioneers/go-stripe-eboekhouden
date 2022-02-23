package json

import (
	"log"

	"github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/hooks"

	jsoniter "github.com/json-iterator/go"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/webhook"
	"github.com/valyala/fasthttp"
)

type hooksHandler struct {
	service          hooks.Servicer
	jsonIteratorPool jsoniter.IteratorPool
	jsonStreamPool   jsoniter.StreamPool
	secret           string
}

func NewHooks(service hooks.Servicer, s string, jip jsoniter.IteratorPool, jsp jsoniter.StreamPool) (hooks.Handler, error) {
	return &hooksHandler{
		service:          service,
		jsonIteratorPool: jip,
		jsonStreamPool:   jsp,
		secret:           s,
	}, nil
}

func (h *hooksHandler) AllHooks(ctx *fasthttp.RequestCtx) {
	ev, err := validateSignature(ctx.PostBody(), string(ctx.Request.Header.Peek("Stripe-Signature")), h.secret)
	if err != nil {
		log.Printf("stripe webhook failure\t %v", err)
		return
	}

	switch ev.Type {
	case "invoice.created":
	default:
		log.Printf("webhook not supported\t\t\t%s", ev.Type)
		err = nil
	}

	if err != nil {
		ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	ctx.Response.SetStatusCode(fasthttp.StatusNoContent)

}

func validateSignature(payload []byte, header, secret string) (stripe.Event, error) {
	return webhook.ConstructEvent(payload, header, secret)
}
