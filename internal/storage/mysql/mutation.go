package mysql

import (
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/domain/mutation"
	queries "github.com/aceworksdev/go-stripe-eboekhouden/internal/storage/mysql/queries/generated"
)

type mutationStorager struct {
	dbs     *DBs
	queries *queries.Queries
}

func NewMutation(dbs *DBs, querier *queries.Queries) (mutation.Storager, error) {
	return &customerStorage{
		dbs:     dbs,
		queries: querier,
	}, nil
}
