package mutation

import (
	"time"

	"github.com/oklog/ulid"
)

type Type string
type TaxType string

const (
	InvoiceReceived        Type = "InvoiceReceived"
	InvoiceSend            Type = "InvoiceSend"
	InvoicePaymentReceived Type = "InvoicePaidReceived"
	InvoicePaymentSend     Type = "InvoicePaidSend"
	MoneyReceived          Type = "MoneyReceived"
	MoneySpend             Type = "MoneySpend"
	Memorial               Type = "Memorial"

	Inclusive TaxType = "IN"
	Exclusice TaxType = "EX"
)

type Service struct {
	ID                ulid.ULID
	Number            int
	Type              Type
	Date              time.Time
	LedgerAccountCode string
	CustomerCode      string
	InvoiceNumber     string
	InvoiceURL        string
	Description       string
	PaymentTerm       string
	PaymentFeature    string
	Items             []ServiceItem
}

type ServiceItem struct {
	VAT               float64
	Amount            float64
	VATCode           string
	VATAmount         int
	LedgerAccountCode string
	KostenplaatsID    int64
}
