package event

import (
	"context"
	"errors"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/contextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type componentProcessor struct {
	dbClient                libmongo.DbClient
	dbEntityCollection      libmongo.DbCollection
	dbAlarmCollection       libmongo.DbCollection
	contextGraphManager     contextgraph.Manager
	eventFilterService      eventfilter.Service
	entityInfosUpdateSender metrics.EntityInfosUpdateSender
	encoder                 encoding.Encoder
	decoder                 encoding.Decoder
	logger                  zerolog.Logger
}

func NewComponentProcessor(
	dbClient libmongo.DbClient,
	contextGraphManager contextgraph.Manager,
	eventFilterService eventfilter.Service,
	entityInfosUpdateSender metrics.EntityInfosUpdateSender,
	encoder encoding.Encoder,
	decoder encoding.Decoder,
	logger zerolog.Logger,
) Processor {
	return &componentProcessor{
		dbClient:                dbClient,
		dbEntityCollection:      dbClient.Collection(libmongo.EntityMongoCollection),
		dbAlarmCollection:       dbClient.Collection(libmongo.AlarmMongoCollection),
		contextGraphManager:     contextGraphManager,
		eventFilterService:      eventFilterService,
		entityInfosUpdateSender: entityInfosUpdateSender,
		encoder:                 encoder,
		decoder:                 decoder,
		logger:                  logger,
	}
}

func (p *componentProcessor) Process(ctx context.Context, event *types.Event, partialRes *ProcessorResult) (ProcessorResult, error) {
	res := ProcessorResult{}
	if partialRes == nil {
		res.EventMetric = techmetrics.CheEventMetric{
			EventMetric: techmetrics.EventMetric{
				EventType: event.EventType,
			},
		}
	} else {
		res.EventMetric = partialRes.EventMetric
	}

	var report contextgraph.Report
	commRegister := libmongo.NewCommandsRegister(p.dbEntityCollection, canopsis.DefaultBulkSize)
	if partialRes == nil {
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
			return res, err
		}

		res.EventMetric.EntityType = event.Entity.Type
		res.EventMetric.IsNewEntity = report.IsNew
	} else {
		report = partialRes.ContextGraphReport
	}

	if event.Entity == nil {
		return res, errors.New("unexpected empty entity")
	}

	if event.Healthcheck {
		return res, nil
	}

	// Process event by event filters.
	if event.Entity.Enabled {
		efr, suspended, err := runEventFilters(ctx, p.eventFilterService, p.encoder, p.decoder, event, &res, partialRes)
		if err != nil {
			return res, err
		}

		if suspended {
			return res, nil
		}

		if len(efr.UpdatedEntityInfos) > 0 {
			_, err = p.dbEntityCollection.UpdateOne(
				ctx,
				bson.M{"_id": event.Entity.ID},
				bson.M{"$set": bson.M{"infos": event.Entity.Infos}},
			)
			if err != nil {
				return res, fmt.Errorf("cannot update entities: %w", err)
			}

			res.EventMetric.IsInfosUpdated = true
			report.CheckComponent = true
			logInfosUpdate(p.entityInfosUpdateSender, event.Entity.ID, efr.UpdatedEntityInfos)
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
		return res, nil
	}

	// cap = 1 for a potential connector counter update.
	toCountersUpdate := make([]types.Entity, 0, 1)
	resourceIDsToUpdateMetrics := make([]string, 0)

	err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		commRegister.Clear()

		toCountersUpdate = toCountersUpdate[:0]
		resourceIDsToUpdateMetrics = resourceIDsToUpdateMetrics[:0]

		res.EventMetric.IsServicesUpdated = false
		res.EventMetric.IsStateSettingUpdated = false

		var component types.Entity
		var connector types.Entity

		cursor, err := p.dbEntityCollection.Find(ctx, bson.M{"_id": bson.M{"$in": entityIdsToCheck}})
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
			return errors.New("component was deleted during event processing")
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
		res.EventMetric.IsServicesUpdated = len(event.Entity.ServicesToAdd) > 0 || len(event.Entity.ServicesToRemove) > 0
		res.EventMetric.IsStateSettingUpdated = stateSettingUpdated

		return nil
	})
	if err != nil {
		return res, err
	}

	if res.EventMetric.IsInfosUpdated {
		go func() {
			err = updateMetaAlarmInfos(ctx, event.Entity.ID, event.Entity.Infos, p.dbAlarmCollection, p.dbEntityCollection)
			if err != nil {
				p.logger.Err(err).Str("entity", event.Entity.ID).Msg("cannot update meta alarm infos")
			}
		}()
	}

	res.UpdatedEntityIdsForMetrics = append(resourceIDsToUpdateMetrics, entityIdsToCheck...)

	return res, nil
}
