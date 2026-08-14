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

type connectorProcessor struct {
	dbClient                libmongo.DbClient
	dbCollection            libmongo.DbCollection
	contextGraphManager     contextgraph.Manager
	eventFilterService      eventfilter.Service
	entityInfosUpdateSender metrics.EntityInfosUpdateSender
	encoder                 encoding.Encoder
	decoder                 encoding.Decoder
}

func NewConnectorProcessor(
	dbClient libmongo.DbClient,
	contextGraphManager contextgraph.Manager,
	eventFilterService eventfilter.Service,
	entityInfosUpdateSender metrics.EntityInfosUpdateSender,
	encoder encoding.Encoder,
	decoder encoding.Decoder,
) Processor {
	return &connectorProcessor{
		dbClient:                dbClient,
		dbCollection:            dbClient.Collection(libmongo.EntityMongoCollection),
		contextGraphManager:     contextGraphManager,
		eventFilterService:      eventFilterService,
		entityInfosUpdateSender: entityInfosUpdateSender,
		encoder:                 encoder,
		decoder:                 decoder,
	}
}

func (p *connectorProcessor) Process(ctx context.Context, event *types.Event, partialRes *ProcessorResult) (ProcessorResult, error) {
	res := ProcessorResult{}
	var report contextgraph.Report
	commRegister := libmongo.NewCommandsRegister(p.dbCollection, canopsis.DefaultBulkSize)
	if partialRes == nil {
		res.EventMetric = techmetrics.CheEventMetric{
			EventMetric: techmetrics.EventMetric{
				EventType: event.EventType,
			},
		}

		err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
			commRegister.Clear()

			var err error
			report, err = p.contextGraphManager.HandleConnector(ctx, event, commRegister)
			if err != nil {
				return fmt.Errorf("cannot update context graph: %w", err)
			}

			return commRegister.Commit(ctx)
		})
		if err != nil {
			return res, err
		}

		res.EventMetric.EntityType = event.Entity.Type
	} else {
		res.EventMetric = partialRes.EventMetric
		report = partialRes.ContextGraphReport
	}

	if event.Entity == nil {
		return res, errors.New("unexpected empty entity")
	}

	if event.Healthcheck {
		return res, nil
	}

	var updatedInfosNames []string

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
			_, err = p.dbCollection.UpdateOne(
				ctx,
				bson.M{"_id": event.Entity.ID},
				bson.M{"$set": bson.M{"infos": event.Entity.Infos}},
			)
			if err != nil {
				return res, fmt.Errorf("cannot update entities: %w", err)
			}

			res.EventMetric.IsInfosUpdated = true
			report.CheckInfoChanged = true
			logInfosUpdate(p.entityInfosUpdateSender, event.Entity.ID, efr.UpdatedEntityInfos)

			updatedInfosNames = make([]string, 0, len(efr.UpdatedEntityInfos))
			for k := range efr.UpdatedEntityInfos {
				updatedInfosNames = append(updatedInfosNames, k)
			}
		}
	}

	if !report.CheckConnector && !report.CheckInfoChanged {
		return res, nil
	}

	res.UpdatedEntityIdsForMetrics = []string{event.Entity.ID}

	err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		commRegister.Clear()

		res.EventMetric.IsServicesUpdated = false

		var connector types.Entity
		err := p.dbCollection.FindOne(ctx, bson.M{"_id": event.Entity.ID}).Decode(&connector)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return errors.New("connector was deleted during event processing")
			}

			return err
		}

		// todo: should be called to get fresh services from the db, should be removed when we do something with cache
		err = p.contextGraphManager.LoadServices(ctx)
		if err != nil {
			return fmt.Errorf("cannot refresh services: %w", err)
		}

		if report.CheckConnector {
			p.contextGraphManager.AssignServices(&connector, commRegister)
		} else if report.CheckInfoChanged {
			p.contextGraphManager.AssignServicesByInfoNames(&connector, updatedInfosNames, commRegister)
		}

		err = commRegister.Commit(ctx)
		if err != nil {
			return err
		}

		event.Entity = &connector
		res.EventMetric.IsServicesUpdated = len(connector.ServicesToAdd) > 0 || len(connector.ServicesToRemove) > 0

		return nil
	})
	if err != nil {
		return res, err
	}

	return res, nil
}
