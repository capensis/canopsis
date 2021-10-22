package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
)

const EnvURL = "CPS_POSTGRES_URL"

func NewPool(ctx context.Context) (*pgxpool.Pool, error) {
	connStr := os.Getenv(EnvURL)
	if connStr == "" {
		return nil, fmt.Errorf("environment variable %s empty", EnvURL)
	}

	return pgxpool.Connect(ctx, connStr)
}
