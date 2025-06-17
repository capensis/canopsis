package che

import (
	"context"
	"fmt"
	"strings"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/importcontextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const ResolveDeletedEventWaitTime = time.Hour

type entityData struct {
	ID                           string            `bson:"_id"`
	Name                         string            `bson:"name"`
	Component                    string            `bson:"component"`
	Type                         string            `bson:"type"`
	ResolveDeletedEventSend      *datetime.CpsTime `bson:"resolve_deleted_event_sent,omitempty"`
	ResolveDeletedEventProcessed *datetime.CpsTime `bson:"resolve_deleted_event_processed,omitempty"`
	SoftDeleted                  datetime.CpsTime  `bson:"soft_deleted"`
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
	now := datetime.NewCpsTime()
	softDeleteTime := now.Add(-w.softDeleteWaitTime)
	resolveDeletedEventAfterTime := now.Add(-ResolveDeletedEventWaitTime)
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
		w.logger.Err(err).Msg("unable to load soft deleted entities")

		return
	}

	defer cursor.Close(ctx)

	idsToRemove := make([]string, 0, canopsis.DefaultBulkSize)
	idsToUpdateEventDate := make([]string, 0, canopsis.DefaultBulkSize)
	serviceIDs := make([]string, 0, canopsis.DefaultBulkSize)
	connectorIDs := make([]string, 0, canopsis.DefaultBulkSize)
	events := make([]types.Event, 0, canopsis.DefaultBulkSize)

	for cursor.Next(ctx) {
		var ent entityData
		err = cursor.Decode(&ent)
		if err != nil {
			w.logger.Err(err).Msg("unable to decode an entity")
			continue
		}

		if ent.ResolveDeletedEventProcessed != nil {
			if softDeleteTime.Before(ent.SoftDeleted.Time) {
				continue
			}

			idsToRemove = append(idsToRemove, ent.ID)
			switch ent.Type {
			case types.EntityTypeConnector:
				connectorIDs = append(connectorIDs, ent.ID)
			case types.EntityTypeService:
				serviceIDs = append(serviceIDs, ent.ID)
			}
		} else if ent.Type != types.EntityTypeService && (ent.ResolveDeletedEventSend == nil || ent.ResolveDeletedEventSend.Time.Before(resolveDeletedEventAfterTime)) {
			event, err := w.createEvent(types.EventTypeResolveDeleted, ent, now)
			if err != nil {
				w.logger.Err(err).Msg("failed to create event")
				continue
			}

			events = append(events, event)
			idsToUpdateEventDate = append(idsToUpdateEventDate, ent.ID)
		} else {
			continue
		}

		if len(idsToRemove) == canopsis.DefaultBulkSize {
			_, err = w.entityCollection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": idsToRemove}, "soft_deleted": bson.M{"$exists": true}})
			if err != nil {
				w.logger.Err(err).Msg("unable to delete soft deletable entities")
			}

			err = w.removeLinksToServices(ctx, serviceIDs)
			if err != nil {
				w.logger.Err(err).Msg("unable to remove links to services")
			}

			err = w.removeLinksToConnectors(ctx, connectorIDs)
			if err != nil {
				w.logger.Err(err).Msg("unable to remove links to services")
			}

			idsToRemove = idsToRemove[:0]
			serviceIDs = serviceIDs[:0]
			connectorIDs = connectorIDs[:0]
		}

		if len(idsToUpdateEventDate) == canopsis.DefaultBulkSize {
			_, err = w.entityCollection.UpdateMany(ctx,
				bson.M{"_id": bson.M{"$in": idsToUpdateEventDate}, "soft_deleted": bson.M{"$exists": true}},
				bson.M{"$set": bson.M{"resolve_deleted_event_sent": now}},
			)
			if err != nil {
				w.logger.Err(err).Msg("unable to update soft deletable entities")
			}

			for _, event := range events {
				err = w.eventPublisher.SendEvent(ctx, event)
				if err != nil {
					w.logger.Warn().Err(err).Msg("failed to send event")
				}
			}

			idsToUpdateEventDate = idsToUpdateEventDate[:0]
			events = events[:0]
		}
	}

	if err = cursor.Err(); err != nil {
		w.logger.Err(err).Msg("unable to fetch soft deleted entities")

		return
	}

	if len(idsToRemove) > 0 {
		_, err = w.entityCollection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": idsToRemove}, "soft_deleted": bson.M{"$exists": true}})
		if err != nil {
			w.logger.Err(err).Msg("unable to delete soft deletable entities")
		}

		err = w.removeLinksToServices(ctx, serviceIDs)
		if err != nil {
			w.logger.Err(err).Msg("unable to remove links to services")
		}

		err = w.removeLinksToConnectors(ctx, connectorIDs)
		if err != nil {
			w.logger.Err(err).Msg("unable to remove links to services")
		}
	}

	if len(idsToUpdateEventDate) > 0 {
		_, err = w.entityCollection.UpdateMany(ctx,
			bson.M{"_id": bson.M{"$in": idsToUpdateEventDate}, "soft_deleted": bson.M{"$exists": true}},
			bson.M{"$set": bson.M{"resolve_deleted_event_sent": now}},
		)
		if err != nil {
			w.logger.Err(err).Msg("unable to update soft deletable entities")
		}

		for _, event := range events {
			err = w.eventPublisher.SendEvent(ctx, event)
			if err != nil {
				w.logger.Warn().Err(err).Msg("failed to send event")
			}
		}
	}
}

func (w *softDeletePeriodicalWorker) createEvent(eventType string, ent entityData, now datetime.CpsTime) (types.Event, error) {
	event := types.Event{
		Connector:     canopsis.CheConnector,
		ConnectorName: canopsis.CheConnector,
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

func (w *softDeletePeriodicalWorker) removeLinksToServices(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	_, err := w.serviceCountersCollection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return fmt.Errorf("unable to delete entity service counters: %w", err)
	}

	for _, id := range ids {
		err = w.removeLinksToService(ctx, id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *softDeletePeriodicalWorker) removeLinksToService(ctx context.Context, id string) error {
	cursor, err := w.entityCollection.Find(ctx, bson.M{"services": id}, options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		return fmt.Errorf("unable to find linked entities: %w", err)
	}

	defer cursor.Close(ctx)
	linkedIDs := make([]string, 0, canopsis.DefaultBulkSize)
	for cursor.Next(ctx) {
		ent := struct {
			ID string `bson:"_id"`
		}{}
		err = cursor.Decode(&ent)
		if err != nil {
			return fmt.Errorf("unable to decode linked entity: %w", err)
		}

		linkedIDs = append(linkedIDs, ent.ID)
		if len(linkedIDs) == canopsis.DefaultBulkSize {
			_, err := w.entityCollection.UpdateMany(ctx, bson.M{"_id": bson.M{"$in": linkedIDs}}, bson.M{"$pull": bson.M{"services": id}})
			if err != nil {
				return fmt.Errorf("unable to update linked entities: %w", err)
			}

			linkedIDs = linkedIDs[:0]
		}
	}

	if len(linkedIDs) > 0 {
		_, err := w.entityCollection.UpdateMany(ctx, bson.M{"_id": bson.M{"$in": linkedIDs}}, bson.M{"$pull": bson.M{"services": id}})
		if err != nil {
			return fmt.Errorf("unable to update linked entities: %w", err)
		}
	}

	if err = cursor.Err(); err != nil {
		return fmt.Errorf("unable to fetch linked entities: %w", err)
	}

	return nil
}

func (w *softDeletePeriodicalWorker) removeLinksToConnectors(ctx context.Context, ids []string) error {
	for _, id := range ids {
		err := w.removeLinksToConnector(ctx, id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *softDeletePeriodicalWorker) removeLinksToConnector(ctx context.Context, id string) error {
	cursor, err := w.entityCollection.Find(ctx, bson.M{"connector": id}, options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		return fmt.Errorf("unable to find linked entities: %w", err)
	}

	defer cursor.Close(ctx)
	linkedIDs := make([]string, 0, canopsis.DefaultBulkSize)
	for cursor.Next(ctx) {
		ent := struct {
			ID string `bson:"_id"`
		}{}
		err = cursor.Decode(&ent)
		if err != nil {
			return fmt.Errorf("unable to decode linked entity: %w", err)
		}

		linkedIDs = append(linkedIDs, ent.ID)
		if len(linkedIDs) == canopsis.DefaultBulkSize {
			_, err := w.entityCollection.UpdateMany(ctx, bson.M{"_id": bson.M{"$in": linkedIDs}}, bson.M{"$unset": bson.M{"connector": ""}})
			if err != nil {
				return fmt.Errorf("unable to update linked entities: %w", err)
			}

			linkedIDs = linkedIDs[:0]
		}
	}

	if len(linkedIDs) > 0 {
		_, err := w.entityCollection.UpdateMany(ctx, bson.M{"_id": bson.M{"$in": linkedIDs}}, bson.M{"$unset": bson.M{"connector": ""}})
		if err != nil {
			return fmt.Errorf("unable to update linked entities: %w", err)
		}
	}

	if err = cursor.Err(); err != nil {
		return fmt.Errorf("unable to fetch linked entities: %w", err)
	}

	return nil
}
