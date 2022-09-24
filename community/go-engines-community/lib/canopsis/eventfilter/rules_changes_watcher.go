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
		loadSleepDuration: loadTimerDuration,
	}
}

type rulesChangesWatcher struct {
	collection        mongo.DbCollection
	service           Service
	logger            zerolog.Logger
	loadSleepDuration time.Duration
}

func (w *rulesChangesWatcher) Watch(ctx context.Context, types []string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	stream, err := w.collection.Watch(ctx, []bson.M{})
	if err != nil {
		return err
	}
	defer stream.Close(ctx)

	// buf = 1, in order not to lose stream event when loader was slept
	loadChan := make(chan struct{}, 1)
	defer close(loadChan)

	// load in case of something was changed in the collection when it wasn't watched
	err = w.service.LoadRules(ctx, types)
	if err != nil {
		w.logger.Error().Err(err).Msg("unable to load rules")
		return nil
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case _, ok := <-loadChan:
				if !ok {
					return
				}
				time.Sleep(w.loadSleepDuration)

				err := w.service.LoadRules(ctx, types)
				if err != nil {
					w.logger.Error().Err(err).Msg("unable to load rules")
				}
			}
		}
	}()

	for stream.Next(ctx) {
		select {
		case <-ctx.Done():
			return nil
		case loadChan <- struct{}{}:
		default:
		}
	}

	return nil
}
