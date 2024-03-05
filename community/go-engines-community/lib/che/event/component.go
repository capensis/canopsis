package event

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/contextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type componentProcessor struct {
	dbClient            libmongo.DbClient
	dbCollection        libmongo.DbCollection
	contextGraphManager contextgraph.Manager
	eventFilterService  eventfilter.Service
}

func NewComponentProcessor(
	dbClient libmongo.DbClient,
	contextGraphManager contextgraph.Manager,
	eventFilterService eventfilter.Service,
) Processor {
	return &componentProcessor{
		dbClient:            dbClient,
		dbCollection:        dbClient.Collection(libmongo.EntityMongoCollection),
		contextGraphManager: contextGraphManager,
		eventFilterService:  eventFilterService,
	}
}

func (p *componentProcessor) Process(ctx context.Context, event *types.Event) (
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

	var report contextgraph.Report
	commRegister := libmongo.NewCommandsRegister(p.dbCollection, canopsis.DefaultBulkSize)

	err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		commRegister.Clear()

		var err error
		report, err = p.contextGraphManager.HandleComponent(ctx, event, commRegister)
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
	eventMetric.IsNewEntity = report.IsNew

	if event.Healthcheck {
		return nil, nil, eventMetric, nil
	}

	// Process event by event filters.
	if event.Entity.Enabled {
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
			report.CheckComponent = true
		}
	}

	// cap = 2: component and connector.
	entityIdsToCheck := make([]string, 0, 2)

	if report.CheckComponent {
		entityIdsToCheck = append(entityIdsToCheck, event.Entity.ID)
	}

	if report.CheckConnector {
		entityIdsToCheck = append(entityIdsToCheck, event.Entity.Connector)
	}

	if len(entityIdsToCheck) == 0 {
		return nil, nil, eventMetric, nil
	}

	// cap = 1 for a potential connector counter update.
	toCountersUpdate := make([]types.Entity, 0, 1)
	resourceIDsToUpdateMetrics := make([]string, 0)

	err = p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		commRegister.Clear()

		toCountersUpdate = toCountersUpdate[:0]
		resourceIDsToUpdateMetrics = resourceIDsToUpdateMetrics[:0]

		eventMetric.IsServicesUpdated = false

		var component types.Entity
		var connector types.Entity

		cursor, err := p.dbCollection.Find(ctx, bson.M{"_id": bson.M{"$in": entityIdsToCheck}})
		if err != nil {
			return err
		}

		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var ent types.Entity

			err = cursor.Decode(&ent)
			if err != nil {
				return err
			}

			if ent.Type == types.EntityTypeComponent {
				component = ent
			} else {
				connector = ent
			}
		}

		if component.ID == "" {
			return fmt.Errorf("component was deleted during event processing")
		}

		// todo: should be called to get fresh services from the db, should be removed when we do something with cache
		err = p.contextGraphManager.LoadServices(ctx)
		if err != nil {
			return fmt.Errorf("cannot refresh services: %w", err)
		}

		p.contextGraphManager.AssignServices(&component, commRegister)

		if connector.ID != "" && report.CheckConnector {
			p.contextGraphManager.AssignServices(&connector, commRegister)
			if len(connector.ServicesToAdd) > 0 || len(connector.ServicesToRemove) > 0 {
				toCountersUpdate = append(toCountersUpdate, connector)
			}
		}

		stateSettingUpdated, err := p.contextGraphManager.AssignStateSetting(ctx, &component, commRegister)
		if err != nil {
			return fmt.Errorf("cannot assign state setting: %w", err)
		}

		resourceIDs, err := p.contextGraphManager.ProcessComponentDependencies(ctx, &component, commRegister)
		if err != nil {
			return fmt.Errorf("cannot process resources: %w", err)
		}

		resourceIDsToUpdateMetrics = append(resourceIDsToUpdateMetrics, resourceIDs...)

		err = commRegister.Commit(ctx)
		if err != nil {
			return err
		}

		event.Entity = &component
		event.StateSettingUpdated = stateSettingUpdated
		eventMetric.IsServicesUpdated = len(event.Entity.ServicesToAdd) > 0 || len(event.Entity.ServicesToRemove) > 0

		return nil
	})
	if err != nil {
		return nil, nil, eventMetric, err
	}

	return toCountersUpdate, append(resourceIDsToUpdateMetrics, entityIdsToCheck...), eventMetric, nil
}
