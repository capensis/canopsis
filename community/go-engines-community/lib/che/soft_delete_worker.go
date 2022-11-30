package che

import (
	"context"
	"fmt"
	"strings"
	"time"

	libmongo "go.mongodb.org/mongo-driver/mongo"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/importcontextgraph"

	"go.mongodb.org/mongo-driver/bson"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"

	"github.com/rs/zerolog"
)

type softDeletePeriodicalWorker struct {
	collection         mongo.DbCollection
	PeriodicalInterval time.Duration
	EventPublisher     importcontextgraph.EventPublisher
	Logger             zerolog.Logger
}

func (w *softDeletePeriodicalWorker) GetInterval() time.Duration {
	return w.PeriodicalInterval
}

func (w *softDeletePeriodicalWorker) Work(ctx context.Context) {
	now := types.CpsTime{Time: time.Now()}

	cursor, err := w.collection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"soft_deleted": true,
				},
			},
			{
				"$lookup": bson.M{
					"from": mongo.AlarmMongoCollection,
					"let":  bson.M{"id": "$_id"},
					"pipeline": []bson.M{
						{
							"$match": bson.M{"$and": []bson.M{
								{"$expr": bson.M{"$eq": bson.A{"$d", "$$id"}}},
								{"v.resolved": nil},
							}},
						},
						{"$limit": 1},
					},
					"as": "alarm",
				},
			},
			{
				"$addFields": bson.M{
					"alarm_exist": bson.M{
						"$cond": bson.M{
							"if":   bson.M{"$eq": bson.A{bson.M{"$size": "$alarm"}, 0}},
							"then": false,
							"else": true,
						},
					},
				},
			},
		},
	)
	if err != nil {
		w.Logger.Error().Err(err).Msg("unable to load soft deleted entities")
	}

	defer cursor.Close(ctx)

	writeModels := make([]libmongo.WriteModel, 0, canopsis.DefaultBulkSize)
	events := make([]types.Event, 0, canopsis.DefaultBulkSize)
	bulkBytesSize := 0

	for cursor.Next(ctx) {
		sendEvent := false

		var ent struct {
			ID                      string         `bson:"_id"`
			Type                    string         `bson:"type"`
			ResolveDeletedEventSend *types.CpsTime `bson:"resolve_deleted_event_sent,omitempty"`
			AlarmExist              bool           `bson:"alarm_exist"`
		}

		err = cursor.Decode(&ent)
		if err != nil {
			w.Logger.Error().Err(err).Msg("unable to decode an entity")
			continue
		}

		var newModel libmongo.WriteModel

		if !ent.AlarmExist {
			newModel = libmongo.
				NewDeleteOneModel().
				SetFilter(bson.M{"_id": ent.ID, "soft_deleted": true})
		} else if ent.ResolveDeletedEventSend == nil || ent.ResolveDeletedEventSend.Add(time.Hour).Before(now.Time) {
			sendEvent = true

			newModel = libmongo.
				NewUpdateOneModel().
				SetFilter(bson.M{"_id": ent.ID, "soft_deleted": true}).
				SetUpdate(bson.M{"$set": bson.M{"resolve_deleted_event_sent": now}})
		} else {
			continue
		}

		b, err := bson.Marshal(newModel)
		if err != nil {
			w.Logger.Error().Err(err).Msg("unable to marshal eventfilter update model")
			continue
		}

		newModelLen := len(b)
		if bulkBytesSize+newModelLen > canopsis.DefaultBulkBytesSize {
			_, err := w.collection.BulkWrite(ctx, writeModels)
			if err != nil {
				w.Logger.Error().Err(err).Msg("unable to bulk write soft deletable entities")
			}

			for _, event := range events {
				err = w.EventPublisher.SendEvent(ctx, event)
				if err != nil {
					w.Logger.Error().Err(err).Msg("failed to send event")
					continue
				}
			}

			writeModels = writeModels[:0]
			events = events[:0]
			bulkBytesSize = 0
		}

		if sendEvent {
			event, err := w.createEvent(types.EventTypeResolveDeleted, ent.Type, ent.ID, now)
			if err != nil {
				w.Logger.Error().Err(err).Msg("failed to create event")
				continue
			}

			events = append(events, event)
		}

		writeModels = append(writeModels, newModel)
		bulkBytesSize += newModelLen

		if len(writeModels) == canopsis.DefaultBulkSize {
			_, err := w.collection.BulkWrite(ctx, writeModels)
			if err != nil {
				w.Logger.Error().Err(err).Msg("unable to bulk write soft deletable entities")
			}

			for _, event := range events {
				err = w.EventPublisher.SendEvent(ctx, event)
				if err != nil {
					w.Logger.Error().Err(err).Msg("failed to send event")
					continue
				}
			}

			writeModels = writeModels[:0]
			events = events[:0]
			bulkBytesSize = 0
		}
	}

	if len(writeModels) > 0 {
		_, err := w.collection.BulkWrite(ctx, writeModels)
		if err != nil {
			w.Logger.Error().Err(err).Msg("unable to bulk write soft deletable entities")
		}

		for _, event := range events {
			err = w.EventPublisher.SendEvent(ctx, event)
			if err != nil {
				w.Logger.Error().Err(err).Msg("failed to send event")
				return
			}
		}
	}
}

func (w *softDeletePeriodicalWorker) createEvent(eventType string, t, id string, now types.CpsTime) (types.Event, error) {
	event := types.Event{
		Connector:     "engine",
		ConnectorName: "engine-che",
		EventType:     eventType,
		Timestamp:     now,
		Author:        canopsis.DefaultEventAuthor,
	}

	switch t {
	case types.EntityTypeComponent:
		event.Component = id
		event.SourceType = types.SourceTypeComponent
	case types.EntityTypeResource:
		idParts := strings.Split(id, "/")
		if len(idParts) != 2 {
			return types.Event{}, fmt.Errorf("invalid resource id = %s", id)
		}
		event.Resource = idParts[0]
		event.Component = idParts[1]
		event.SourceType = types.SourceTypeResource
	}

	return event, nil
}
