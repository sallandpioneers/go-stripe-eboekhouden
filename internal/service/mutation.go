package service

import (
	"context"

	"github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/mutation"
)

type mutationService struct {
	push mutation.Pusher
}

func NewMutation(p mutation.Pusher) (mutation.Servicer, error) {
	return &mutationService{
		push: p,
	}, nil
}

func (service *mutationService) Create(ctx context.Context, item *mutation.Service) error {
	if err := service.push.Create(ctx, item); err != nil {
		return err
	}
	return nil
}
