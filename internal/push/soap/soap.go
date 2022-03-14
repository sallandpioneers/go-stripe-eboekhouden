package soap

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aceworksdev/go-stripe-eboekhouden/internal/config"
	eboekhouden "github.com/aceworksdev/go-stripe-eboekhouden/internal/push/soap/generated"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/server/domain/push"
	"github.com/hooklift/gowsdl/soap"
)

func New(p *push.Soap, c *config.EBoekHouden) error {
	client := soap.NewClient("https://soap.e-boekhouden.nl/soap.asmx?WSDL")
	service = eboekhouden.NewSoapAppSoap(client)
	cfg = *c

	// Add the securityCode to the session
	session.SecurityCode2 = c.SecurityCode2

	var err error
	if p.Customer, err = NewCustomer(service); err != nil {
		return err
	}
	if p.Mutation, err = NewMutation(service); err != nil {
		return err
	}
	return nil
}

var session Session
var service eboekhouden.SoapAppSoap
var cfg config.EBoekHouden

func GetSession() (*Session, error) {
	if session.createdAt.Before(time.Now().Add(-5 * time.Minute)) {
		if err := createNewSession(); err != nil {
			return nil, err
		}
	}
	return &session, nil
}

func createNewSession() error {
	resp, err := service.OpenSessionContext(context.TODO(), &eboekhouden.OpenSession{
		Username:      cfg.Username,
		SecurityCode1: cfg.SecurityCode1,
		SecurityCode2: cfg.SecurityCode2,
	})
	if err != nil {
		return err
	}

	if resp.OpenSessionResult != nil {
		if err := handleError(resp.OpenSessionResult.ErrorMsg); err != nil {
			return err
		}
		session.SessionID = resp.OpenSessionResult.SessionID
		return nil
	}
	return errors.New("Something went wrong while creating a new session")
}

type Session struct {
	SessionID     string
	SecurityCode2 string
	createdAt     time.Time
}

func handleError(err *eboekhouden.CError) error {
	if err == nil {
		return nil
	}

	switch err.LastErrorCode {
	case "":
		return nil
	case "E0023":
		createNewSession()
		fallthrough
	default:
		return fmt.Errorf("An eboekhouding error happend:\ncode\t\t%s\ndescription:\t%s", err.LastErrorCode, err.LastErrorDescription)
	}
}
