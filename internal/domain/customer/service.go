package customer

import "context"

type Servicer interface {
	Create(context.Context, *Service) error
}
