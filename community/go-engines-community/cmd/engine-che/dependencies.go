package main

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/che"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
)

func NewEngine(ctx context.Context, opts che.Options, logger zerolog.Logger) engine.Engine {
	mongoClient, err := mongo.NewClient(ctx, 0, 0, logger)
	if err != nil {
		panic(fmt.Errorf("cannot connect to MongoDb: %w", err))
	}

	return che.NewEngine(ctx, opts, mongoClient, nil, metrics.NewNullMetaUpdater(), eventfilter.NewExternalDataGetterContainer(), logger)
}
