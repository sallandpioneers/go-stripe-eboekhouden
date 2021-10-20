package invoice

import (
	"time"

	"github.com/oklog/ulid"
)

type DirectDebitType string

const (
	SingleUse  DirectDebitType = "E"
	Continuous DirectDebitType = "D"
)

type Service struct {
	ID              ulid.ULID
	Number          string
	CustomerCode    string
	Date            time.Time
	DaysUntilDue    int
	InvoiceTemplate string
	Email           struct {
		Send    bool
		Subject string
		Message string
		Address string
		Name    string
	}
	AutomaticDebtCollection bool
	DirectDebit             struct {
		IBAN                string
		Type                DirectDebitType
		ID                  string
		Date                time.Time
		First               bool
		InNameOf            string
		Location            string
		Description         string
		Add                 bool
		MutationDescription string
	}
	Items []ItemService
}

type ItemService struct {
	Amount            int
	Unit              string
	Code              string
	Description       string
	PricePerUnit      string
	BTWCode           string
	ContraAccountCode string // !?
	KostenplaatsID    int    // !?
}
