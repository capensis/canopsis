package che

import (
	"context"
	"fmt"
	"strings"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/importcontextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	libmongo "go.mongodb.org/mongo-driver/mongo"
)

const ResolveDeletedEventWaitTime = time.Hour

type entityData struct {
	ID                           string         `bson:"_id"`
	Name                         string         `bson:"name"`
	Component                    string         `bson:"component"`
	Type                         string         `bson:"type"`
	ResolveDeletedEventSend      *types.CpsTime `bson:"resolve_deleted_event_sent,omitempty"`
	ResolveDeletedEventProcessed *types.CpsTime `bson:"resolve_deleted_event_processed,omitempty"`
	SoftDeleted                  types.CpsTime  `bson:"soft_deleted"`
}

type softDeletePeriodicalWorker struct {
	entityCollection          mongo.DbCollection
	serviceCountersCollection mongo.DbCollection
	periodicalInterval        time.Duration
	eventPublisher            importcontextgraph.EventPublisher
	softDeleteWaitTime        time.Duration
	logger                    zerolog.Logger
}

func (w *softDeletePeriodicalWorker) GetInterval() time.Duration {
	return w.periodicalInterval
}

// Work checks all soft deleted entities.
// If service counters are recomputed it deletes entity. If not it sends another event (but not for services).
func (w *softDeletePeriodicalWorker) Work(ctx context.Context) {
	now := types.CpsTime{Time: time.Now()}

	cursor, err := w.entityCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"soft_deleted": bson.M{"$exists": true},
				},
			},
			{
				"$project": bson.M{
					"_id":                             1,
					"name":                            1,
					"component":                       1,
					"type":                            1,
					"resolve_deleted_event_sent":      1,
					"resolve_deleted_event_processed": 1,
					"soft_deleted":                    1,
				},
			},
		},
	)
	if err != nil {
		w.logger.Error().Err(err).Msg("unable to load soft deleted entities")
		return
	}

	defer cursor.Close(ctx)

	writeModels := make([]libmongo.WriteModel, 0, canopsis.DefaultBulkSize)
	serviceCountersIDs := make([]string, 0, canopsis.DefaultBulkSize)
	events := make([]types.Event, 0, canopsis.DefaultBulkSize)
	bulkBytesSize := 0

	for cursor.Next(ctx) {
		sendEvent := false
		var ent entityData
		err = cursor.Decode(&ent)
		if err != nil {
			w.logger.Error().Err(err).Msg("unable to decode an entity")
			continue
		}

		var newModels []libmongo.WriteModel
		if ent.Type == types.EntityTypeService {
			serviceCountersIDs = append(serviceCountersIDs, ent.ID)
			if len(serviceCountersIDs) == canopsis.DefaultBulkSize {
				_, err = w.serviceCountersCollection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": serviceCountersIDs}})
				if err != nil {
					w.logger.Error().Err(err).Msg("unable to delete entity service counters")
				}

				serviceCountersIDs = serviceCountersIDs[:0]
			}
		}

		if ent.ResolveDeletedEventProcessed != nil {
			if now.Add(-w.softDeleteWaitTime).Before(ent.SoftDeleted.Time) {
				continue
			}

			switch ent.Type {
			case types.EntityTypeConnector:
				newModels = []libmongo.WriteModel{
					libmongo.NewUpdateManyModel().
						SetFilter(bson.M{"connector": ent.ID}).
						SetUpdate(bson.M{"$unset": bson.M{"connector": ""}}),
				}
			case types.EntityTypeService:
				newModels = []libmongo.WriteModel{
					libmongo.NewUpdateManyModel().
						SetFilter(bson.M{"services": ent.ID}).
						SetUpdate(bson.M{"$pull": bson.M{"services": ent.ID}}),
				}
			}

			newModels = append(newModels,
				libmongo.
					NewDeleteOneModel().
					SetFilter(bson.M{"_id": ent.ID, "soft_deleted": bson.M{"$exists": true}}),
			)
		} else if ent.Type != types.EntityTypeService && (ent.ResolveDeletedEventSend == nil || ent.ResolveDeletedEventSend.Add(ResolveDeletedEventWaitTime).Before(now.Time)) {
			sendEvent = true
			newModels = []libmongo.WriteModel{
				libmongo.
					NewUpdateOneModel().
					SetFilter(bson.M{"_id": ent.ID, "soft_deleted": bson.M{"$exists": true}}).
					SetUpdate(bson.M{"$set": bson.M{"resolve_deleted_event_sent": now}}),
			}
		} else {
			continue
		}

		b, err := bson.Marshal(struct {
			Arr []libmongo.WriteModel
		}{Arr: writeModels})
		if err != nil {
			w.logger.Error().Err(err).Msg("unable to marshal eventfilter update model")
			continue
		}

		newModelLen := len(b)
		if bulkBytesSize+newModelLen > canopsis.DefaultBulkBytesSize {
			_, err := w.entityCollection.BulkWrite(ctx, writeModels)
			if err != nil {
				w.logger.Error().Err(err).Msg("unable to bulk write soft deletable entities")
			}

			for _, event := range events {
				err = w.eventPublisher.SendEvent(ctx, event)
				if err != nil {
					w.logger.Error().Err(err).Msg("failed to send event")
					continue
				}
			}

			writeModels = writeModels[:0]
			events = events[:0]
			bulkBytesSize = 0
		}

		if sendEvent {
			event, err := w.createEvent(types.EventTypeResolveDeleted, ent, now)
			if err != nil {
				w.logger.Error().Err(err).Msg("failed to create event")
				continue
			}

			events = append(events, event)
		}

		writeModels = append(writeModels, newModels...)
		bulkBytesSize += newModelLen

		if len(writeModels) == canopsis.DefaultBulkSize {
			_, err := w.entityCollection.BulkWrite(ctx, writeModels)
			if err != nil {
				w.logger.Error().Err(err).Msg("unable to bulk write soft deletable entities")
			}

			for _, event := range events {
				err = w.eventPublisher.SendEvent(ctx, event)
				if err != nil {
					w.logger.Error().Err(err).Msg("failed to send event")
					continue
				}
			}

			writeModels = writeModels[:0]
			events = events[:0]
			bulkBytesSize = 0
		}
	}

	if len(writeModels) > 0 {
		_, err := w.entityCollection.BulkWrite(ctx, writeModels)
		if err != nil {
			w.logger.Error().Err(err).Msg("unable to bulk write soft deletable entities")
		}

		for _, event := range events {
			err = w.eventPublisher.SendEvent(ctx, event)
			if err != nil {
				w.logger.Error().Err(err).Msg("failed to send event")
				return
			}
		}
	}

	if len(serviceCountersIDs) > 0 {
		_, err = w.serviceCountersCollection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": serviceCountersIDs}})
		if err != nil {
			w.logger.Error().Err(err).Msg("unable to delete entity service counters")
		}
	}
}

func (w *softDeletePeriodicalWorker) createEvent(eventType string, ent entityData, now types.CpsTime) (types.Event, error) {
	event := types.Event{
		Connector:     canopsis.EngineConnector,
		ConnectorName: canopsis.CheEngineName,
		EventType:     eventType,
		Timestamp:     now,
		Author:        canopsis.DefaultEventAuthor,
		Initiator:     types.InitiatorSystem,
	}

	switch ent.Type {
	case types.EntityTypeConnector:
		event.SourceType = types.SourceTypeConnector
		event.Connector = strings.TrimSuffix(ent.ID, "/"+ent.Name)
		event.ConnectorName = ent.Name
	case types.EntityTypeComponent:
		event.SourceType = types.SourceTypeComponent
		event.Component = ent.ID
	case types.EntityTypeResource:
		event.SourceType = types.SourceTypeResource

		event.Resource = ent.Name
		event.Component = ent.Component
		if event.Component == "" {
			idParts := strings.Split(ent.ID, "/")
			if len(idParts) != 2 {
				return types.Event{}, fmt.Errorf("invalid resource id = %s", ent.ID)
			}

			event.Resource = idParts[0]
			event.Component = idParts[1]
		}
	}

	return event, nil
}
