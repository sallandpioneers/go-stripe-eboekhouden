package mutation

import "context"

type Servicer interface {
	Create(ctx context.Context, item *Service) error
}
