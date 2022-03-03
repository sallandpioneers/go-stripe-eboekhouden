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

func (service *mutationPush) Create(ctx context.Context, item *mutation.Service, customerCode string) error {
	session, err := GetSession()
	if err != nil {
		return err
	}

	addMutation := &eboekhouden.AddMutatie{
		SessionID:     session.SessionID,
		SecurityCode2: session.SecurityCode2,
		OMut: &eboekhouden.CMutatie{
			MutatieNr:        0,
			Soort:            nil,
			Datum:            soap.CreateXsdDateTime(item.Date, false),
			Rekening:         "",
			RelatieCode:      customerCode,
			Factuurnummer:    item.InvoiceNumber,
			Boekstuk:         "",
			Omschrijving:     "",
			Betalingstermijn: "",
			Betalingskenmerk: "",
			InExBTW:          "",
			MutatieRegels: &eboekhouden.ArrayOfCMutatieRegel{
				CMutatieRegel: make([]*eboekhouden.CMutatieRegel, len(item.Items)),
			},
		},
	}

	for k, v := range item.Items {
		addMutation.OMut.MutatieRegels.CMutatieRegel[k].BedragInvoer = v.Amount
		addMutation.OMut.MutatieRegels.CMutatieRegel[k].BedragExclBTW = 0
		addMutation.OMut.MutatieRegels.CMutatieRegel[k].BedragBTW = 0
		addMutation.OMut.MutatieRegels.CMutatieRegel[k].BedragInclBTW = 0
		addMutation.OMut.MutatieRegels.CMutatieRegel[k].BTWCode = ""
		addMutation.OMut.MutatieRegels.CMutatieRegel[k].BTWPercentage = 0
		addMutation.OMut.MutatieRegels.CMutatieRegel[k].TegenrekeningCode = v.LedgerAccountCode
		addMutation.OMut.MutatieRegels.CMutatieRegel[k].KostenplaatsID = v.KostenplaatsID
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

func (service *mutationPush) Update(ctx context.Context, item *mutation.Service, customerCode string) error {
	return nil
}
