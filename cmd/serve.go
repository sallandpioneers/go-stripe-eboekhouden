package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/aceworksdev/go-stripe-eboekhouden/internal"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/config"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/customer"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/handler"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/push"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/server"
	handlerStruct "github.com/aceworksdev/go-stripe-eboekhouden/internal/server/domain/handler"
	pushStruct "github.com/aceworksdev/go-stripe-eboekhouden/internal/server/domain/push"
	serviceStruct "github.com/aceworksdev/go-stripe-eboekhouden/internal/server/domain/service"
	storageStruct "github.com/aceworksdev/go-stripe-eboekhouden/internal/server/domain/storage"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/service"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/storage"

	"github.com/spf13/cobra"
	"github.com/stripe/stripe-go/v72"
)

//nolint:gochecknoinits // this is the recommended way to use cobra
func init() {
}

//nolint:gochecknoglobals // this is the recommended way to use cobra
var serveCommand = &cobra.Command{
	Use:   "serve",
	Short: "Start the webhook listener and will wait until stripe events are thrown",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("STARTING stripe eboekhouding coupling")
		log.SetFlags(log.Lshortfile | log.Ltime)

		// Load config
		c := config.New()
		if err := c.Validate(); err != nil {
			log.Fatal("[CONFIG] ", err)
		}

		// Set stripe ID
		stripe.Key = c.Stripe.Key

		// Init Server
		serve := server.New(c.Server)
		client := server.NewClient(c.Client)

		sa := &internal.ServicesAvailable{}
		// Build storage layer. DB and FS

		s := &storageStruct.Storage{}
		if err := storage.New(s, c.DB, sa); err != nil {
			log.Fatal(err)
		}

		p := &pushStruct.Push{
			Soap: &pushStruct.Soap{},
		}
		if err := push.New(p, c.EBoekHouden, true, true); err != nil {
			log.Fatal(err)
		}

		// Build service layer
		serv := &serviceStruct.Service{}
		if err := service.New(serv, s, p, client, c); err != nil {
			log.Fatal(err)
		}

		// Build hand layer. JSON
		hand := &handlerStruct.Handler{}
		err := handler.New("json", hand, serv, sa, c)
		if err != nil {
			log.Fatal(err)
		}

		server.NewRouter(serve, hand, c.Router)

		if err := p.Soap.Customer.Create(context.TODO(), &customer.Service{
			Company: "Test 1234",
		}); err != nil {
			log.Fatal(err)
		}

		log.Printf("Starting up go-stripe-eboekhouden back-end, listening on port: %d\n", c.Server.Port)
		log.Fatal(serve.ListenAndServe(fmt.Sprintf(":%d", c.Server.Port)))
	},
}
