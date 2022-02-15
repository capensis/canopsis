package main

import (
	"context"
	libengine "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/depmake"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/fifo"
	"github.com/rs/zerolog"
)

func NewEngine(ctx context.Context, options fifo.Options, logger zerolog.Logger) libengine.Engine {
	defer depmake.Catch(logger)

	var m depmake.DependencyMaker

	return fifo.Default(ctx, options, m.DepMongoClient(ctx, logger), eventfilter.NewExternalDataGetterContainer(), logger)
}
