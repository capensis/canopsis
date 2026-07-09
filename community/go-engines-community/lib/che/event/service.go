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
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type serviceProcessor struct {
	dbClient                libmongo.DbClient
	dbCollection            libmongo.DbCollection
	contextGraphManager     contextgraph.Manager
	eventFilterService      eventfilter.Service
	entityInfosUpdateSender metrics.EntityInfosUpdateSender
	encoder                 encoding.Encoder
	decoder                 encoding.Decoder
}

func NewServiceProcessor(
	dbClient libmongo.DbClient,
	contextGraphManager contextgraph.Manager,
	eventFilterService eventfilter.Service,
	entityInfosUpdateSender metrics.EntityInfosUpdateSender,
	encoder encoding.Encoder,
	decoder encoding.Decoder,
) Processor {
	return &serviceProcessor{
		dbClient:                dbClient,
		dbCollection:            dbClient.Collection(libmongo.EntityMongoCollection),
		contextGraphManager:     contextGraphManager,
		eventFilterService:      eventFilterService,
		entityInfosUpdateSender: entityInfosUpdateSender,
		encoder:                 encoder,
		decoder:                 decoder,
	}
}

func (p *serviceProcessor) Process(ctx context.Context, event *types.Event, partialRes *ProcessorResult) (ProcessorResult, error) {
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

	commRegister := libmongo.NewCommandsRegister(p.dbCollection, canopsis.DefaultBulkSize)
	if partialRes == nil {
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
				return res, err
			}

			event.Entity = &eventEntity
			res.EventMetric.EntityType = eventEntity.Type

			return res, nil
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
			return res, err
		}

		res.EventMetric.EntityType = event.Entity.Type
	}

	if event.Entity == nil {
		return res, errors.New("unexpected empty entity")
	}

	var checkServices bool

	// Process event by event filters.
	if event.Entity.Enabled && !event.Healthcheck {
		efr, suspended, err := runEventFilters(ctx, p.eventFilterService, p.encoder, p.decoder, event, &res, partialRes)
		if err != nil {
			return res, err
		}

		if suspended {
			return res, nil
		}

		if len(efr.UpdatedEntityInfos) > 0 {
			_, err = p.dbCollection.UpdateOne(
				ctx,
				bson.M{"_id": event.Entity.ID},
				bson.M{"$set": bson.M{"infos": event.Entity.Infos}},
			)
			if err != nil {
				return res, fmt.Errorf("cannot update entities: %w", err)
			}

			res.EventMetric.IsInfosUpdated = true
			checkServices = true
			logInfosUpdate(p.entityInfosUpdateSender, event.Entity.ID, efr.UpdatedEntityInfos)
		}
	}

	if event.Healthcheck || !checkServices {
		return res, nil
	}

	res.UpdatedEntityIdsForMetrics = []string{event.Entity.ID}

	err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		commRegister.Clear()
		res.EventMetric.IsServicesUpdated = false
		res.EventMetric.IsStateSettingUpdated = false

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

		res.EventMetric.IsStateSettingUpdated, err = p.contextGraphManager.AssignStateSetting(ctx, &service, commRegister)
		if err != nil {
			return fmt.Errorf("cannot inherit component fields: %w", err)
		}

		err = commRegister.Commit(ctx)
		if err != nil {
			return err
		}

		event.Entity = &service
		res.EventMetric.IsServicesUpdated = len(event.Entity.ServicesToAdd) > 0 || len(event.Entity.ServicesToRemove) > 0

		return nil
	})
	if err != nil {
		return res, err
	}

	return res, nil
}
