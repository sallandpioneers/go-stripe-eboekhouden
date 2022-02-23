package server

import (
	"log"

	"github.com/aceworksdev/go-stripe-eboekhouden/internal/config"
	"github.com/valyala/fasthttp"
)

func New(config *config.Server) (server *fasthttp.Server) {
	log.Println("Init server")
	server = &fasthttp.Server{
		Name:               config.Name,
		MaxRequestBodySize: config.MaxRequestBodySize,
	}
	log.Println("Init server done")
	return
}
