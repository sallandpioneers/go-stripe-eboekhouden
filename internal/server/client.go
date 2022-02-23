package server

import (
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/config"
	"github.com/valyala/fasthttp"
)

func NewClient(c *config.Client) *fasthttp.Client {
	return &fasthttp.Client{
		ReadTimeout:     c.ReadTimeout,
		ReadBufferSize:  c.ReadBufferSize,
		WriteTimeout:    c.WriteTimeout,
		WriteBufferSize: c.WriteBufferSize,
	}
}
