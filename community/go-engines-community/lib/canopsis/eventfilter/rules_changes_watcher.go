package eventfilter

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type RuleChangesWatcherInterface interface {
	Watch(ctx context.Context, types []string) error
}

func NewRulesChangesWatcher(client mongo.DbClient, service Service) RuleChangesWatcherInterface {
	return &rulesChangesWatcher{
		collection: client.Collection(mongo.EventFilterRulesMongoCollection),
		service:    service,
	}
}

type rulesChangesWatcher struct {
	collection mongo.DbCollection
	service    Service
}

func (w *rulesChangesWatcher) Watch(ctx context.Context, types []string) error {
	stream, err := w.collection.Watch(ctx, []bson.M{})
	if err != nil {
		return err
	}

	defer stream.Close(ctx)

	for stream.Next(ctx) {
		err := w.service.LoadRules(ctx, types)
		if err != nil {
			return fmt.Errorf("unable to load rules: %w", err)
		}
	}

	return nil
}
