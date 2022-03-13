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
			AddDatum:  soap.CreateXsdDateTime(time.Now(), true),
			Code:      item.BoekhoudenID,
			Bedrijf:   item.Name,
			BP:        string(item.Type),
			Adres:     item.Addresses.Business.Address,
			Postcode:  item.Addresses.Business.ZipCode,
			Plaats:    item.Addresses.Business.City,
			Land:      item.Addresses.Business.Country,
			Adres2:    item.Addresses.Mailing.Address,
			Postcode2: item.Addresses.Mailing.ZipCode,
			Plaats2:   item.Addresses.Mailing.City,
			Land2:     item.Addresses.Mailing.Country,
			Telefoon:  item.Phone,
			Email:     item.Email,
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

func (service *customerPush) Update(ctx context.Context, item *customer.Service) error {
	session, err := GetSession()
	if err != nil {
		return err
	}

	resp, err := service.client.UpdateRelatieContext(ctx, &eboekhouden.UpdateRelatie{
		SessionID:     session.SessionID,
		SecurityCode2: session.SecurityCode2,
		ORel: &eboekhouden.CRelatie{
			AddDatum:  soap.CreateXsdDateTime(time.Now(), true),
			Code:      item.BoekhoudenID,
			Bedrijf:   item.Name,
			BP:        string(item.Type),
			Adres:     item.Addresses.Business.Address,
			Postcode:  item.Addresses.Business.ZipCode,
			Plaats:    item.Addresses.Business.City,
			Land:      item.Addresses.Business.Country,
			Adres2:    item.Addresses.Mailing.Address,
			Postcode2: item.Addresses.Mailing.ZipCode,
			Plaats2:   item.Addresses.Mailing.City,
			Land2:     item.Addresses.Mailing.Country,
			Telefoon:  item.Phone,
			Email:     item.Email,
		},
	})
	if err != nil {
		return err
	}
	if resp.UpdateRelatieResult != nil {
		if err := handleError(resp.UpdateRelatieResult); err != nil {
			return err
		}
	}

	return nil
}
