package datastorage

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

var ErrConfNotFound = errors.New("cannot find configuration _id=" + ID)

type adapter struct {
	collection mongo.DbCollection
}

func NewAdapter(client mongo.DbClient) Adapter {
	return &adapter{
		collection: client.Collection(mongo.ConfigurationMongoCollection),
	}
}

func (a *adapter) Get(ctx context.Context) (DataStorage, error) {
	data := DataStorage{}
	err := a.collection.FindOne(ctx, bson.M{"_id": ID}).Decode(&data)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return data, nil
		}

		return data, err
	}

	return data, nil
}

func (a *adapter) updateOne(ctx context.Context, upd bson.M) error {
	res, err := a.collection.UpdateOne(ctx, bson.M{"_id": ID}, bson.M{"$set": upd})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrConfNotFound
	}

	return nil
}

func (a *adapter) updateDatetime(ctx context.Context, key string, t datetime.CpsTime) error {
	return a.updateOne(ctx, bson.M{key: t})
}

func (a *adapter) updateHistoryWithCount(ctx context.Context, key string, history HistoryWithCount) error {
	return a.updateOne(ctx, bson.M{key: history})
}

func (a *adapter) UpdateHistoryJunit(ctx context.Context, t datetime.CpsTime) error {
	return a.updateDatetime(ctx, "history.junit", t)
}

func (a *adapter) UpdateHistoryRemediation(ctx context.Context, t datetime.CpsTime) error {
	return a.updateDatetime(ctx, "history.remediation", t)
}

func (a *adapter) UpdateHistoryAlarm(ctx context.Context, history HistoryWithCount) error {
	return a.updateHistoryWithCount(ctx, "history.alarm", history)
}

func (a *adapter) UpdateHistoryAlarmExternalTag(ctx context.Context, history HistoryWithCount) error {
	return a.updateHistoryWithCount(ctx, "history.alarm_external_tag", history)
}

func (a *adapter) UpdateHistoryEntityDisabled(ctx context.Context, history HistoryWithCount) error {
	return a.updateHistoryWithCount(ctx, "history.entity_disabled", history)
}

func (a *adapter) UpdateHistoryEntityUnlinked(ctx context.Context, history HistoryWithCount) error {
	return a.updateHistoryWithCount(ctx, "history.entity_unlinked", history)
}

func (a *adapter) UpdateHistoryEntityCleaned(ctx context.Context, history HistoryWithCount) error {
	return a.updateHistoryWithCount(ctx, "history.entity_cleaned", history)
}

func (a *adapter) UpdateHistoryPbehavior(ctx context.Context, t datetime.CpsTime) error {
	return a.updateDatetime(ctx, "history.pbehavior", t)
}

func (a *adapter) UpdateHistoryHealthCheck(ctx context.Context, t datetime.CpsTime) error {
	return a.updateDatetime(ctx, "history.health_check", t)
}

func (a *adapter) UpdateHistoryWebhook(ctx context.Context, t datetime.CpsTime) error {
	return a.updateDatetime(ctx, "history.webhook", t)
}

func (a *adapter) UpdateHistoryEventFilterFailure(ctx context.Context, t datetime.CpsTime) error {
	return a.updateDatetime(ctx, "history.event_filter_failure", t)
}

func (a *adapter) UpdateHistoryEventRecords(ctx context.Context, t datetime.CpsTime) error {
	return a.updateDatetime(ctx, "history.event_records", t)
}
