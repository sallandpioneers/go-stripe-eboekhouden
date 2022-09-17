package adapter

import (
	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/model"
	stripeModel "github.com/sallandpioneers/go-stripe-eboekhouden/internal/payment/stripe/model"
)

func ToReport(i model.Report) stripeModel.Report {
	return stripeModel.Report{}
}
