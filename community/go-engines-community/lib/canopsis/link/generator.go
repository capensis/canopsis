package link

import (
	"context"
)

type Generator interface {
	GenerateForAlarms(ctx context.Context, ids []string) (map[string]LinksByCategory, error)
	GenerateForEntities(ctx context.Context, ids []string) (map[string]LinksByCategory, error)
	GenerateForAllAlarms(ctx context.Context, ids []string) ([]Link, error)
	Load(ctx context.Context) error
}
