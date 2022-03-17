package mysql

import (
	"context"

	"github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/customer"
	queries "github.com/aceworksdev/go-stripe-eboekhouden/internal/storage/mysql/queries/generated"
	"github.com/oklog/ulid"
)

type customerStorage struct {
	dbs     *DBs
	queries *queries.Queries
}

func NewCustomer(dbs *DBs, querier *queries.Queries) (customer.Storager, error) {
	return &customerStorage{
		dbs:     dbs,
		queries: querier,
	}, nil
}

func (storage *customerStorage) Create(ctx context.Context, item *customer.Service) error {
	arg := &queries.CreateCustomerParams{
		ID:             item.ID,
		StripeID:       item.StripeID,
		BoekhoudenID:   item.RelationID,
		BoekhoudenCode: item.BoekhoudenID,
	}

	if err := storage.queries.CreateCustomer(ctx, arg); err != nil {
		return err
	}
	return nil
}

func (storage *customerStorage) Get(ctx context.Context, id ulid.ULID) (*customer.Service, error) {
	dbRecord, err := storage.queries.GetCustomer(ctx, id)
	if err != nil {
		return nil, err
	}
	return &customer.Service{
		ID:           dbRecord.ID,
		StripeID:     dbRecord.StripeID,
		RelationID:   dbRecord.BoekhoudenID,
		BoekhoudenID: dbRecord.BoekhoudenCode,
	}, nil
}

func (storage *customerStorage) GetBasedOnStripeID(ctx context.Context, id string) (*customer.Service, error) {
	dbRecord, err := storage.queries.GetCustomerBasedOnStripeID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &customer.Service{
		ID:           dbRecord.ID,
		StripeID:     dbRecord.StripeID,
		RelationID:   dbRecord.BoekhoudenID,
		BoekhoudenID: dbRecord.BoekhoudenCode,
	}, nil
}

func (storage *customerStorage) GetBasedOnBoekhoudenID(ctx context.Context, id int64) (*customer.Service, error) {
	dbRecord, err := storage.queries.GetCustomerBasedOnBoekhoudenID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &customer.Service{
		ID:           dbRecord.ID,
		StripeID:     dbRecord.StripeID,
		RelationID:   dbRecord.BoekhoudenID,
		BoekhoudenID: dbRecord.BoekhoudenCode,
	}, nil
}
