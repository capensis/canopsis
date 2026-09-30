package event

import (
	"context"
	"errors"
	"fmt"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/contextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	libevent "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type serviceProcessor struct {
	dbClient                libmongo.DbClient
	dbCollection            libmongo.DbCollection
	contextGraphManager     contextgraph.Manager
	eventFilterService      eventfilter.Service
	entityInfosUpdateSender metrics.EntityInfosUpdateSender
	eventGenerator          libevent.Generator
	publisher               libamqp.Publisher
	encoder                 encoding.Encoder
	decoder                 encoding.Decoder
}

func NewServiceProcessor(
	dbClient libmongo.DbClient,
	contextGraphManager contextgraph.Manager,
	eventFilterService eventfilter.Service,
	entityInfosUpdateSender metrics.EntityInfosUpdateSender,
	eventGenerator libevent.Generator,
	publisher libamqp.Publisher,
	encoder encoding.Encoder,
	decoder encoding.Decoder,
) Processor {
	return &serviceProcessor{
		dbClient:                dbClient,
		dbCollection:            dbClient.Collection(libmongo.EntityMongoCollection),
		contextGraphManager:     contextGraphManager,
		eventFilterService:      eventFilterService,
		entityInfosUpdateSender: entityInfosUpdateSender,
		eventGenerator:          eventGenerator,
		publisher:               publisher,
		encoder:                 encoder,
		decoder:                 decoder,
	}
}

func (p *serviceProcessor) Process(ctx context.Context, event *types.Event, partialRes *ProcessorResult) (ProcessorResult, error) {
	res := ProcessorResult{}
	var report contextgraph.Report
	commRegister := libmongo.NewCommandsRegister(p.dbCollection, canopsis.DefaultBulkSize)
	if partialRes == nil {
		res.EventMetric = techmetrics.CheEventMetric{
			EventType: event.EventType,
		}

		if event.EventType == types.EventTypeRecomputeEntityService {
			var eventEntity types.Entity
			var pauseChanges contextgraph.ServicePauseChanges

			err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
				commRegister.Clear()

				eventEntity = types.Entity{}
				pauseChanges = contextgraph.ServicePauseChanges{}
				var err error

				eventEntity, pauseChanges, err = p.contextGraphManager.RecomputeService(ctx, event.GetEID(), commRegister)
				if err != nil {
					return fmt.Errorf("cannot recompute service %s: %w", event.Component, err)
				}

				return commRegister.Commit(ctx)
			})
			if err != nil {
				return res, err
			}

			// The recompute doesn't update the pause of the dependencies, so pause and resume events
			// are sent for the ones which have been added to or removed from the service.
			err = p.sendPauseChanges(ctx, event, pauseChanges)
			if err != nil {
				return res, fmt.Errorf("cannot send pause changes of service %s: %w", event.Component, err)
			}

			event.Entity = &eventEntity
			res.EventMetric.EntityType = eventEntity.Type

			return res, nil
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
	if event.Entity.Enabled && !event.Healthcheck {
		efr, suspended, err := runEventFilters(ctx, p.eventFilterService, p.encoder, p.decoder, event, &res, partialRes)
		if err != nil {
			return res, err
		}

		if suspended {
			res.ContextGraphReport = report

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

	if !report.CheckService && !report.CheckInfoChanged {
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

		if report.CheckService {
			p.contextGraphManager.AssignServices(&service, commRegister)
		} else if report.CheckInfoChanged {
			p.contextGraphManager.AssignServicesByInfoNames(&service, updatedInfosNames, commRegister)
		}

		res.EventMetric.IsStateSettingUpdated, err = p.contextGraphManager.AssignStateSetting(ctx, &service, commRegister)
		if err != nil {
			return fmt.Errorf("cannot inherit component fields: %w", err)
		}

		err = commRegister.Commit(ctx)
		if err != nil {
			return err
		}

		event.Entity = &service
		res.EventMetric.IsServicesUpdated = len(service.ServicesToAdd) > 0 || len(service.ServicesToRemove) > 0

		return nil
	})
	if err != nil {
		return res, err
	}

	return res, nil
}

// sendPauseChanges sends a pause event for each new dependency of the paused service and a resume event
// for each removed dependency which the service has paused, so their pause, their alarm
// and their own dependencies are updated by engine-axe.
func (p *serviceProcessor) sendPauseChanges(ctx context.Context, event *types.Event, changes contextgraph.ServicePauseChanges) error {
	for _, entity := range changes.ToPause {
		pauseEvent, err := p.eventGenerator.Generate(entity)
		if err != nil {
			return fmt.Errorf("cannot generate pause event of %s: %w", entity.ID, err)
		}

		pause := changes.ServicePause
		pauseEvent.EventType = types.EventTypePause
		pauseEvent.Timestamp = pause.Timestamp
		pauseEvent.Output = pause.Message
		pauseEvent.Author = pause.Author
		pauseEvent.UserID = pause.UserID
		pauseEvent.Initiator = pause.Initiator
		pauseEvent.PauseType = pause.Type
		pauseEvent.PauseSource = pause.Source
		pauseEvent.PauseSourceType = pause.SourceType
		pauseEvent.PausePath = pause.Path
		pauseEvent.Inherited = true

		err = p.sendEvent(ctx, pauseEvent)
		if err != nil {
			return err
		}
	}

	for _, entity := range changes.ToResume {
		resumeEvent, err := p.eventGenerator.Generate(entity)
		if err != nil {
			return fmt.Errorf("cannot generate resume event of %s: %w", entity.ID, err)
		}

		resumeEvent.EventType = types.EventTypeResume
		resumeEvent.Timestamp = event.Timestamp
		resumeEvent.Initiator = types.InitiatorSystem
		resumeEvent.PauseSource = entity.Pause.Source
		resumeEvent.PauseSourceType = entity.Pause.SourceType
		resumeEvent.PausePath = entity.Pause.Path

		err = p.sendEvent(ctx, resumeEvent)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *serviceProcessor) sendEvent(ctx context.Context, event types.Event) error {
	body, err := p.encoder.Encode(event)
	if err != nil {
		return fmt.Errorf("cannot encode event: %w", err)
	}

	return p.publisher.PublishWithContext(
		ctx,
		"",
		canopsis.FIFOQueueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  canopsis.JsonContentType,
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
}
