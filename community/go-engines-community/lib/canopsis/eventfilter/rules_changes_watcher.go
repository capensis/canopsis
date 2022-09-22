package eventfilter

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
)

const loadTimerDuration = 200 * time.Millisecond

type RuleChangesWatcher interface {
	Watch(ctx context.Context, types []string) error
}

func NewRulesChangesWatcher(client mongo.DbClient, service Service, logger zerolog.Logger) RuleChangesWatcher {
	return &rulesChangesWatcher{
		collection:        client.Collection(mongo.EventFilterRulesMongoCollection),
		service:           service,
		logger:            logger,
		loadTimerDuration: loadTimerDuration,
	}
}

type rulesChangesWatcher struct {
	collection        mongo.DbCollection
	service           Service
	logger            zerolog.Logger
	loadTimerDuration time.Duration
}

func (w *rulesChangesWatcher) Watch(ctx context.Context, types []string) error {
	stream, err := w.collection.Watch(ctx, []bson.M{})
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	loadTimer := time.NewTimer(w.loadTimerDuration)
	loadTimer.Stop()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-loadTimer.C:
				err := w.service.LoadRules(ctx, types)
				if err != nil {
					w.logger.Error().Err(err).Msg("unable to load rules")
				}
			}
		}
	}()

	defer stream.Close(ctx)

	for stream.Next(ctx) {
		loadTimer.Reset(w.loadTimerDuration)
	}

	return nil
}
