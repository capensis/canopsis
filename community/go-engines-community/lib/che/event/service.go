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

type serviceProcessor struct {
	dbClient            libmongo.DbClient
	dbCollection        libmongo.DbCollection
	contextGraphManager contextgraph.Manager
	eventFilterService  eventfilter.Service
}

func NewServiceProcessor(
	dbClient libmongo.DbClient,
	contextGraphManager contextgraph.Manager,
	eventFilterService eventfilter.Service,
) Processor {
	return &serviceProcessor{
		dbClient:            dbClient,
		dbCollection:        dbClient.Collection(libmongo.EntityMongoCollection),
		contextGraphManager: contextGraphManager,
		eventFilterService:  eventFilterService,
	}
}

func (p *serviceProcessor) Process(ctx context.Context, event *types.Event) (
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

	if event.EventType == types.EventTypeRecomputeEntityService {
		var eventEntity types.Entity

		err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
			commRegister.Clear()

			eventEntity = types.Entity{}
			var err error

			eventEntity, err = p.contextGraphManager.RecomputeService(ctx, event.GetEID(), commRegister)
			if err != nil {
				return fmt.Errorf("cannot recompute service %s: %w", event.Component, err)
			}

			return commRegister.Commit(ctx)
		})
		if err != nil {
			return nil, nil, eventMetric, err
		}

		event.Entity = &eventEntity
		eventMetric.EntityType = eventEntity.Type

		return nil, nil, eventMetric, nil
	}

	err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		commRegister.Clear()

		var err error
		_, err = p.contextGraphManager.HandleService(ctx, event, commRegister)
		if err != nil {
			return fmt.Errorf("cannot update context graph: %w", err)
		}

		return commRegister.Commit(ctx)
	})
	if err != nil {
		return nil, nil, eventMetric, err
	}

	if event.Entity == nil {
		return nil, nil, eventMetric, errors.New("unexpected empty entity")
	}

	eventMetric.EntityType = event.Entity.Type
	var checkServices bool

	// Process event by event filters.
	if event.Entity.Enabled && !event.Healthcheck {
		isInfosUpdated, err := p.eventFilterService.ProcessEvent(ctx, event)
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

	if event.Healthcheck || !checkServices {
		return nil, nil, eventMetric, nil
	}

	entityIdsToMetrics := []string{event.Entity.ID}

	err = p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		commRegister.Clear()
		eventMetric.IsServicesUpdated = false

		var service types.Entity
		err := p.dbCollection.FindOne(ctx, bson.M{"_id": event.Entity.ID}).Decode(&service)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return errors.New("service was deleted during event processing")
			}

			return err
		}

		// todo: should be called to get fresh services from the db, should be removed when we do something with cache
		err = p.contextGraphManager.LoadServices(ctx)
		if err != nil {
			return fmt.Errorf("cannot refresh services: %w", err)
		}

		p.contextGraphManager.AssignServices(&service, commRegister)

		_, err = p.contextGraphManager.AssignStateSetting(ctx, &service, commRegister)
		if err != nil {
			return fmt.Errorf("cannot inherit component fields: %w", err)
		}

		err = commRegister.Commit(ctx)
		if err != nil {
			return err
		}

		event.Entity = &service
		eventMetric.IsServicesUpdated = len(event.Entity.ServicesToAdd) > 0 || len(event.Entity.ServicesToRemove) > 0

		return nil
	})
	if err != nil {
		return nil, nil, eventMetric, err
	}

	return nil, entityIdsToMetrics, eventMetric, nil
}
