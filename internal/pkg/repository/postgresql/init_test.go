package postgresql

import (
	"context"
	"homework-7/internal/tests/postgres"
)

var (
	ctx context.Context
	Db  *postgres.TDB
)

func init() {
	Db = postgres.NewFromEnv()
	ctx = context.Background()
}
