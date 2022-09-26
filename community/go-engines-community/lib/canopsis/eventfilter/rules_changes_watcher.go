package eventfilter

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

const loadTimerDuration = 200 * time.Millisecond

type RuleChangesWatcher interface {
	Watch(ctx context.Context, types []string) error
}

func NewRulesChangesWatcher(client mongo.DbClient, service Service) RuleChangesWatcher {
	return &rulesChangesWatcher{
		collection:        client.Collection(mongo.EventFilterRulesMongoCollection),
		service:           service,
		loadSleepDuration: loadTimerDuration,
	}
}

type rulesChangesWatcher struct {
	collection        mongo.DbCollection
	service           Service
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

	// load in case of something was changed in the collection when it wasn't watched
	err = w.service.LoadRules(ctx, types)
	if err != nil {
		return err
	}

	go func() {
		defer close(loadChan)

		for stream.Next(ctx) {
			select {
			case <-ctx.Done():
				return
			case loadChan <- struct{}{}:
			default:
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case _, ok := <-loadChan:
			if !ok {
				return nil
			}

			select {
			case <-ctx.Done():
				return nil
			case <-time.After(w.loadSleepDuration):
			}

			err := w.service.LoadRules(ctx, types)
			if err != nil {
				return err
			}
		}
	}
}
