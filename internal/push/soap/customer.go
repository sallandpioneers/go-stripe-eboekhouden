package soap

import (
	"context"
	"time"

	"github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/customer"
	eboekhouden "github.com/aceworksdev/go-stripe-eboekhouden/internal/push/soap/generated"
	"github.com/hooklift/gowsdl/soap"
)

type customerPush struct {
	isDevelopment bool
	client        eboekhouden.SoapAppSoap
}

func NewCustomer(client eboekhouden.SoapAppSoap) (customer.Pusher, error) {
	return &customerPush{
		client: client,
	}, nil
}

func (service *customerPush) Create(ctx context.Context, item *customer.Service) error {
	session, err := GetSession()
	if err != nil {
		return err
	}

	resp, err := service.client.AddRelatieContext(ctx, &eboekhouden.AddRelatie{
		SessionID:     session.SessionID,
		SecurityCode2: session.SecurityCode2,
		ORel: &eboekhouden.CRelatie{
			AddDatum: soap.CreateXsdDateTime(time.Now(), true),
			Code:     item.Code, // Generate random code
			Bedrijf:  item.Company,
			BP:       string(item.Type),
		},
	})
	if err != nil {
		return err
	}
	if resp.AddRelatieResult != nil {
		if err := handleError(resp.AddRelatieResult.ErrorMsg); err != nil {
			return err
		}

		item.RelationID = resp.AddRelatieResult.Rel_ID
	}
	return nil
}
