package soap

import (
	"context"

	"github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/mutation"
	eboekhouden "github.com/aceworksdev/go-stripe-eboekhouden/internal/push/soap/generated"
	"github.com/hooklift/gowsdl/soap"
)

type mutationPush struct {
	client eboekhouden.SoapAppSoap
}

func NewMutation(client eboekhouden.SoapAppSoap) (mutation.Pusher, error) {
	return &mutationPush{
		client: client,
	}, nil
}

func (service *mutationPush) Create(ctx context.Context, item *mutation.Service) error {
	session, err := GetSession()
	if err != nil {
		return err
	}

	mutationType := eboekhouden.EnMutatieSoorten(item.Type)

	addMutation := &eboekhouden.AddMutatie{
		SessionID:     session.SessionID,
		SecurityCode2: session.SecurityCode2,
		OMut: &eboekhouden.CMutatie{
			MutatieNr:        0,
			Soort:            &mutationType,
			Datum:            soap.CreateXsdDateTime(item.Date, false),
			Rekening:         item.LedgerAccountCode,
			RelatieCode:      item.BoekhoudenCustomerID,
			Factuurnummer:    item.InvoiceNumber,
			Boekstuk:         item.LedgerAccountCode,
			Omschrijving:     item.Description,
			Betalingstermijn: item.PaymentTerm,
			Betalingskenmerk: item.PaymentFeature,
			InExBTW:          string(mutation.Exclusive),
			MutatieRegels: &eboekhouden.ArrayOfCMutatieRegel{
				CMutatieRegel: make([]*eboekhouden.CMutatieRegel, len(item.Items)),
			},
		},
	}

	for k, v := range item.Items {
		addMutation.OMut.MutatieRegels.CMutatieRegel[k] = &eboekhouden.CMutatieRegel{
			BedragInvoer:      v.Amount,
			BedragExclBTW:     v.AmountExVAT,
			BedragBTW:         v.AmountVAT,
			BedragInclBTW:     v.AmountInVAT,
			BTWCode:           string(v.VATCode),
			BTWPercentage:     v.VATPercentage,
			TegenrekeningCode: v.LedgerAccountCode,
			KostenplaatsID:    v.KostenplaatsID,
		}
	}

	resp, err := service.client.AddMutatieContext(ctx, addMutation)
	if err != nil {
		return err
	}
	if resp.AddMutatieResult != nil {
		if err := handleError(resp.AddMutatieResult.ErrorMsg); err != nil {
			return err
		}
	}

	return nil
}
