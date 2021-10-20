package mutation

import (
	"time"

	"github.com/oklog/ulid"
)

type Type string
type TaxType string

const (
	InvoiceReceived        Type = "FactuurOntvangen"
	InvoiceSend            Type = "FactuurVerstuurd"
	InvoicePaymentReceived Type = "FactuurbetalingOntvangen"
	InvoicePaymentSend     Type = "FactuurbetalingVerstuurd"
	MoneyReceived          Type = "GeldOntvangen"
	MoneySpend             Type = "GeldUitgegeven"
	Memorial               Type = "Memoriaal"

	Inclusive TaxType = "IN"
	Exclusice TaxType = "EX"
)

type Service struct {
	ID                ulid.ULID
	MutationNR        int
	Type              Type
	Date              time.Time
	LedgerAccountCode string
	CustomerCode      string
	InvoiceNumber     string
	InvoiceURL        string // Boekstuk, will be url to pdf hosted by stripe
	Description       string
	PaymentTerm       string
	PaymentFeature    string
	Items             []ServiceItem
}

type ServiceItem struct {
	VAT               int
	Amount            int
	VATCode           string
	VATAmount         int
	LedgerAccountCode string
	KostenplaatsID    int
}
