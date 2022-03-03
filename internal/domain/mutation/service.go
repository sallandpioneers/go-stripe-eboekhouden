package mutation

import "context"

type Servicer interface {
	Create(ctx context.Context, item *Service, customerCode string) error
	Update(ctx context.Context, item *Service, customerCode string) error
}
