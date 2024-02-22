package statesetting

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

func NewRulesChangesWatcher(client mongo.DbClient, service Assigner) RuleChangesWatcher {
	return &rulesChangesWatcher{
		collection: client.Collection(mongo.EngineNotificationCollection),
		service:    service,
	}
}

type rulesChangesWatcher struct {
	collection mongo.DbCollection
	service    Assigner
}

func (w *rulesChangesWatcher) Watch(ctx context.Context) error {
	// load in case of something was changed in the collection when it wasn't watched
	err := w.service.LoadRules(ctx)
	if err != nil {
		return err
	}

	stream, err := w.collection.Watch(
		ctx,
		[]bson.M{{"$match": bson.M{"documentKey._id": StateSettingsNotificationID}}},
	)
	if err != nil {
		return err
	}

	defer stream.Close(ctx)

	for stream.Next(ctx) {
		select {
		case <-ctx.Done():
			return nil
		default:
			err = w.service.LoadRules(ctx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
