package customer

import "context"

type Servicer interface {
	Create(context.Context, *Service) error
	Update(ctx context.Context, item *Service) error
}
