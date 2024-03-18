package hooks

import "github.com/valyala/fasthttp"

type Handler interface {
	AllHooks(ctx *fasthttp.RequestCtx)
}
