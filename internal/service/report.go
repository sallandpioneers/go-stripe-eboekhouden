package service

import (
	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/domain/report"
	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/model"
	"github.com/sallandpioneers/go-stripe-eboekhouden/internal/payment"
)

type reportService struct {
	payment payment.Service
}

func NewReport(p payment.Service) (report.Servicer, error) {
	return &reportService{
		payment: p,
	}, nil
}

func (s *reportService) Create() error {
	return s.payment.CreateReport(model.Report{})
}
