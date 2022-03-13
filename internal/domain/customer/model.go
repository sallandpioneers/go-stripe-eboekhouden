package customer

import "github.com/oklog/ulid"

type Gender string
type Type string

const Male Gender = "m"
const Female Gender = "v"

const Business Type = "b"
const Private Type = "p"

type Service struct {
	ID           ulid.ULID
	StripeID     string
	BoekhoudenID string
	Name         string

	RelationID int64
	Contact    string
	Addresses  struct {
		Business struct {
			Address string
			ZipCode string
			City    string
			Country string
		}
		Mailing struct {
			Address string
			ZipCode string
			City    string
			Country string
		}
	}
	Phone      string
	Email      string
	Website    string
	Notition   string
	VAT        string
	COC        string
	Salutation string
	IBAN       string
	BIC        string
	Type       Type
}
