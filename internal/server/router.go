package server

import (
	"log"

	"github.com/aceworksdev/go-stripe-eboekhouden/internal/config"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/server/domain/handler"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func NewRouter(s *fasthttp.Server, h *handler.Handler, config *config.Router) *router.Router {
	log.Println("Init router")

	r := router.New()

	r.RedirectTrailingSlash = config.RedirectTrailingSlash
	r.RedirectFixedPath = config.RedirectFixedPath
	r.HandleMethodNotAllowed = config.HandleMethodNotAllowed
	r.HandleOPTIONS = config.HandleOPTIONS

	// Set servers router handler to the newly made router
	s.Handler = r.Handler

	initRoutes(r, h)
	log.Println("Init router done")

	return r
}

//nolint function to long
func initRoutes(r *router.Router, h *handler.Handler) {
	r.POST("/hooks", h.Hooks.AllHooks)
}
