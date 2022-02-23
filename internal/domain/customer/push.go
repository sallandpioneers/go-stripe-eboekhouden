package customer

import "context"

type Pusher interface {
	Create(ctx context.Context, item *Service) error
}