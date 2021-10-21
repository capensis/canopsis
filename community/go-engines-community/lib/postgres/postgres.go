package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
)

const envURL = "CPS_POSTGRES_URL"

func NewPool(ctx context.Context) (*pgxpool.Pool, error) {
	connStr := os.Getenv(envURL)
	if connStr == "" {
		return nil, fmt.Errorf("environment variable %s empty", envURL)
	}

	return pgxpool.Connect(ctx, connStr)
}
