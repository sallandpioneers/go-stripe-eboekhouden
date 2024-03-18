package invoice

import "context"

type Pusher interface {
	Create(ctx context.Context, item *Service, customerCode string) error
	Update(ctx context.Context, item *Service, customerCode string) error
}
