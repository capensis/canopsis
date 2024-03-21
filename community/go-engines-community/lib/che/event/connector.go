package event

import (
	"context"
	"errors"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/contextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type connectorProcessor struct {
	dbClient            libmongo.DbClient
	dbCollection        libmongo.DbCollection
	contextGraphManager contextgraph.Manager
	eventFilterService  eventfilter.Service
}

func NewConnectorProcessor(
	dbClient libmongo.DbClient,
	contextGraphManager contextgraph.Manager,
	eventFilterService eventfilter.Service,
) Processor {
	return &connectorProcessor{
		dbClient:            dbClient,
		dbCollection:        dbClient.Collection(libmongo.EntityMongoCollection),
		contextGraphManager: contextGraphManager,
		eventFilterService:  eventFilterService,
	}
}

func (p *connectorProcessor) Process(ctx context.Context, event *types.Event) (
	[]types.Entity,
	[]string,
	techmetrics.CheEventMetric,
	error,
) {
	eventMetric := techmetrics.CheEventMetric{
		EventMetric: techmetrics.EventMetric{
			EventType: event.EventType,
		},
	}

	commRegister := libmongo.NewCommandsRegister(p.dbCollection, canopsis.DefaultBulkSize)
	err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		commRegister.Clear()

		_, err := p.contextGraphManager.HandleConnector(ctx, event, commRegister)
		if err != nil {
			return fmt.Errorf("cannot update context graph: %w", err)
		}

		return commRegister.Commit(ctx)
	})
	if err != nil {
		return nil, nil, eventMetric, err
	}

	if event.Entity == nil {
		return nil, nil, eventMetric, fmt.Errorf("unexpected empty entity")
	}

	eventMetric.EntityType = event.Entity.Type
	var checkServices bool

	if event.Healthcheck {
		return nil, nil, eventMetric, nil
	}

	// Process event by event filters.
	if event.Entity.Enabled {
		var isInfosUpdated bool
		isInfosUpdated, eventMetric.ExecutedEnrichRules, eventMetric.ExternalRequests, err = p.eventFilterService.ProcessEvent(ctx, event)
		if err != nil {
			return nil, nil, eventMetric, err
		}

		if isInfosUpdated {
			_, err = p.dbCollection.UpdateOne(
				ctx,
				bson.M{"_id": event.Entity.ID},
				bson.M{"$set": bson.M{"infos": event.Entity.Infos}},
			)
			if err != nil {
				return nil, nil, eventMetric, fmt.Errorf("cannot update entities: %w", err)
			}

			eventMetric.IsInfosUpdated = true
			checkServices = true
		}
	}

	if !checkServices {
		return nil, nil, eventMetric, nil
	}

	entityIdsToMetrics := []string{event.Entity.ID}

	err = p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		commRegister.Clear()

		eventMetric.IsServicesUpdated = false

		var connector types.Entity
		err := p.dbCollection.FindOne(ctx, bson.M{"_id": event.Entity.ID}).Decode(&connector)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return fmt.Errorf("connector was deleted during event processing")
			}

			return err
		}

		// todo: should be called to get fresh services from the db, should be removed when we do something with cache
		err = p.contextGraphManager.LoadServices(ctx)
		if err != nil {
			return fmt.Errorf("cannot refresh services: %w", err)
		}

		p.contextGraphManager.AssignServices(&connector, commRegister)

		err = commRegister.Commit(ctx)
		if err != nil {
			return err
		}

		event.Entity = &connector
		eventMetric.IsServicesUpdated = len(event.Entity.ServicesToAdd) > 0 || len(event.Entity.ServicesToRemove) > 0

		return nil
	})
	if err != nil {
		return nil, nil, eventMetric, err
	}

	return nil, entityIdsToMetrics, eventMetric, nil
}
