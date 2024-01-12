package eventfilter

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/sync/errgroup"
)

// loadTimerDuration is needed to avoid flood from mongo change stream when a lot of documents are updated in the same time.
const loadTimerDuration = 200 * time.Millisecond

type RuleChangesWatcher interface {
	Watch(ctx context.Context, types []string) error
}

func NewRulesChangesWatcher(client mongo.DbClient, service Service) RuleChangesWatcher {
	return &rulesChangesWatcher{
		collection:        client.Collection(mongo.EventFilterRuleCollection),
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
	eg, ctx := errgroup.WithContext(ctx)

	// buf = 1, in order not to lose stream event when loader was slept
	loadChan := make(chan struct{}, 1)

	// load in case of something was changed in the collection when it wasn't watched
	err := w.service.LoadRules(ctx, types)
	if err != nil {
		return err
	}

	eg.Go(func() error {
		defer close(loadChan)

		stream, err := w.collection.Watch(ctx, []bson.M{})
		if err != nil {
			return err
		}

		defer stream.Close(ctx)

		for stream.Next(ctx) {
			select {
			case <-ctx.Done():
				return nil
			case loadChan <- struct{}{}:
			default:
			}
		}

		return nil
	})

	eg.Go(func() error {
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
	})

	return eg.Wait()
}
