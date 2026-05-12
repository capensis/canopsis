package event

import (
	"context"
	"errors"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/contextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type serviceProcessor struct {
	dbClient            libmongo.DbClient
	dbCollection        libmongo.DbCollection
	contextGraphManager contextgraph.Manager
	eventFilterService  eventfilter.Service
	metricsSender       metrics.Sender
}

func NewServiceProcessor(
	dbClient libmongo.DbClient,
	contextGraphManager contextgraph.Manager,
	eventFilterService eventfilter.Service,
	metricsSender metrics.Sender,
) Processor {
	return &serviceProcessor{
		dbClient:            dbClient,
		dbCollection:        dbClient.Collection(libmongo.EntityMongoCollection),
		contextGraphManager: contextGraphManager,
		eventFilterService:  eventFilterService,
		metricsSender:       metricsSender,
	}
}

func (p *serviceProcessor) Process(ctx context.Context, event *types.Event) (ProcessorResult, error) {
	result := ProcessorResult{
		EventMetric: techmetrics.CheEventMetric{
			EventMetric: techmetrics.EventMetric{
				EventType: event.EventType,
			},
		},
	}

	var report contextgraph.Report
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
			return result, err
		}

		event.Entity = &eventEntity
		result.EventMetric.EntityType = eventEntity.Type

		return result, nil
	}

	err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		commRegister.Clear()

		var err error
		report, err = p.contextGraphManager.HandleService(ctx, event, commRegister)
		if err != nil {
			return fmt.Errorf("cannot update context graph: %w", err)
		}

		return commRegister.Commit(ctx)
	})
	if err != nil {
		return result, err
	}

	if event.Entity == nil {
		return result, errors.New("unexpected empty entity")
	}

	result.EventMetric.EntityType = event.Entity.Type

	if event.Healthcheck {
		return result, nil
	}

	var updatedInfos map[string]eventfilter.UpdatedValue
	var updatedInfosNames []string

	// Process event by event filters.
	if event.Entity.Enabled {
		updatedInfos, result.EventMetric.ExecutedEnrichRules, result.EventMetric.ExternalRequests, err = p.eventFilterService.ProcessEvent(ctx, event)
		if err != nil {
			return result, err
		}

		if len(updatedInfos) > 0 {
			_, err = p.dbCollection.UpdateOne(
				ctx,
				bson.M{"_id": event.Entity.ID},
				bson.M{"$set": bson.M{"infos": event.Entity.Infos}},
			)
			if err != nil {
				return result, fmt.Errorf("cannot update entities: %w", err)
			}

			result.EventMetric.IsInfosUpdated = true
			report.CheckInfoChanged = true
			logInfosUpdate(p.metricsSender, event.Entity.ID, updatedInfos)

			updatedInfosNames = make([]string, 0, len(updatedInfos))
			for k := range updatedInfos {
				updatedInfosNames = append(updatedInfosNames, k)
			}
		}
	}

	if !report.CheckService && !report.CheckInfoChanged {
		return result, nil
	}

	entityIdsToMetrics := []string{event.Entity.ID}

	err = p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		commRegister.Clear()
		result.EventMetric.IsServicesUpdated = false
		result.EventMetric.IsStateSettingUpdated = false

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

		if report.CheckService {
			p.contextGraphManager.AssignServices(&service, commRegister)
		} else if report.CheckInfoChanged {
			p.contextGraphManager.AssignServicesByInfoNames(&service, updatedInfosNames, commRegister)
		}

		result.EventMetric.IsStateSettingUpdated, err = p.contextGraphManager.AssignStateSetting(ctx, &service, commRegister)
		if err != nil {
			return fmt.Errorf("cannot inherit component fields: %w", err)
		}

		err = commRegister.Commit(ctx)
		if err != nil {
			return err
		}

		event.Entity = &service
		result.EventMetric.IsServicesUpdated = len(service.ServicesToAdd) > 0 || len(service.ServicesToRemove) > 0

		return nil
	})
	if err != nil {
		return result, err
	}

	result.UpdatedEntityIDsForMetrics = entityIdsToMetrics

	return result, nil
}
