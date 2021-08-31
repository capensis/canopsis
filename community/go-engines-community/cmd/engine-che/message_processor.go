package main

import (
	"context"
	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	libcontext "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/errt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
	"runtime/trace"
	"time"
)

type messageProcessor struct {
	FeaturePrintEventOnError bool
	FeatureEventProcessing   bool
	FeatureContextCreation   bool
	AlarmConfigProvider      config.AlarmConfigProvider
	EventFilterService       eventfilter.Service
	EnrichmentCenter         libcontext.EnrichmentCenter
	EnrichFields             libcontext.EnrichFields
	AmqpPublisher            libamqp.Publisher
	AlarmAdapter             alarm.Adapter
	Encoder                  encoding.Encoder
	Decoder                  encoding.Decoder
	Logger                   zerolog.Logger
}

func (p *messageProcessor) Process(parentCtx context.Context, d amqp.Delivery) ([]byte, error) {
	ctx, task := trace.NewTask(parentCtx, "che.WorkerProcess")
	defer task.End()

	trace.Logf(ctx, "event_size", "%d", len(d.Body))
	var event types.Event
	err := p.Decoder.Decode(d.Body, &event)
	if err != nil {
		p.logError(err, "cannot decode event", d.Body)
		return nil, nil
	}

	err = event.InjectExtraInfos(d.Body)
	if err != nil {
		p.logError(err, "cannot inject extra infos", d.Body)
		return nil, nil
	}

	p.Logger.Debug().Msgf("valid input event: %v", string(d.Body))
	trace.Log(ctx, "event.event_type", event.EventType)
	trace.Log(ctx, "event.timestamp", event.Timestamp.String())
	trace.Log(ctx, "event.source_type", event.SourceType)
	trace.Log(ctx, "event.connector", event.Connector)
	trace.Log(ctx, "event.connector_name", event.ConnectorName)
	trace.Log(ctx, "event.component", event.Component)
	trace.Log(ctx, "event.resource", event.Resource)

	err = event.IsValid()
	if err != nil {
		p.logError(err, "invalid event", d.Body)
		return nil, nil
	}

	alarmConfig := p.AlarmConfigProvider.Get()
	event.Output = utils.TruncateString(event.Output, alarmConfig.OutputLength)
	event.LongOutput = utils.TruncateString(event.LongOutput, alarmConfig.OutputLength)
	updatedEntityServices := libcontext.UpdatedEntityServices{}

	// Process event by event filters.
	if p.FeatureEventProcessing {
		var report eventfilter.Report
		event, report, err = p.EventFilterService.ProcessEvent(ctx, event)
		if err != nil {
			if engine.IsConnectionError(err) {
				return nil, err
			}

			p.logError(err, "cannot apply event filters on event", d.Body)
			return nil, nil
		}

		updatedEntityServices = updatedEntityServices.Add(report.UpdatedEntityServices)

		if report.EntityUpdated && event.Entity != nil {
			updated, err := p.EnrichmentCenter.UpdateEntityInfos(ctx, event.Entity)
			if err != nil {
				if engine.IsConnectionError(err) {
					return nil, err
				}

				p.logError(err, "cannot update entity infos", d.Body)
				return nil, nil
			}

			updatedEntityServices = updatedEntityServices.Add(updated)
		}
	}

	// Enrich the event with the entity and create the context if this has not
	// been done by the event filter.
	if event.Entity == nil && p.FeatureContextCreation && event.IsContextable() {
		entity, updated, err := p.EnrichmentCenter.Handle(ctx, event, p.EnrichFields)
		if err != nil {
			if engine.IsConnectionError(err) {
				return nil, err
			}

			p.logError(err, "cannot update context graph", d.Body)
			return nil, nil
		}
		event.Entity = entity
		updatedEntityServices = updatedEntityServices.Add(updated)
	}

	// Update context graph for entity service.
	if event.EventType == types.EventTypeUpdateEntityService {
		err = p.EnrichmentCenter.ReloadService(ctx, event.GetEID())
		if err != nil {
			if engine.IsConnectionError(err) {
				return nil, err
			}

			p.logError(err, "cannot update service", d.Body)
			return nil, nil
		}
	} else if event.EventType == types.EventTypeRecomputeEntityService {
		updated, err := p.EnrichmentCenter.HandleEntityServiceUpdate(ctx, event.GetEID())
		if err != nil {
			if engine.IsConnectionError(err) {
				return nil, err
			}

			p.logError(err, "cannot update entity service", d.Body)
			return nil, nil
		}

		if updated != nil {
			updatedEntityServices = updatedEntityServices.Add(*updated)
		} else {
			// If context graph is not updated do not recompute service state.
			event.EventType = types.EventTypeUpdateEntityService
		}
	}

	if !p.FeatureEventProcessing {
		return nil, nil
	}

	event.Format()

	// Find entity if still empty.
	if event.Entity == nil {
		event.Entity, err = p.EnrichmentCenter.Get(ctx, event)
		if err != nil {
			if engine.IsConnectionError(err) {
				return nil, err
			}

			p.logError(err, "cannot find entity", d.Body)
			return nil, nil
		}
	}

	event.AddedToServices = append(event.AddedToServices, updatedEntityServices.AddedTo...)
	event.RemovedFromServices = append(event.RemovedFromServices, updatedEntityServices.RemovedFrom...)

	err = p.publishComponentInfosUpdatedEvents(ctx, updatedEntityServices.UpdatedComponentInfosResources)
	if err != nil {
		return nil, err
	}

	body, err := p.Encoder.Encode(event)
	if err != nil {
		p.logError(err, "cannot encode event", d.Body)
		return nil, nil
	}

	p.Logger.Debug().Msgf("output event about to be published: %+v", string(body))

	return body, nil
}

func (p *messageProcessor) logError(err error, errMsg string, msg []byte) {
	if p.FeaturePrintEventOnError {
		p.Logger.Err(err).Str("event", string(msg)).Msg(errMsg)
	} else {
		p.Logger.Err(err).Msg(errMsg)
	}
}

// publishComponentInfosUpdatedEvents sends event to update context graph and state of
// entity services if there are entity services which depend on component infos and
// component infos of resources have been updated on component event.
// It's not possible to immediately process such resources  since only component entity
// is locked by engine fifo and resource entity can be updated by another event in parallel.
func (p *messageProcessor) publishComponentInfosUpdatedEvents(ctx context.Context, resources []string) error {
	if len(resources) == 0 {
		return nil
	}

	alarms := make([]types.Alarm, 0)
	err := p.AlarmAdapter.GetOpenedAlarmsByIDs(ctx, resources, &alarms)
	if err != nil {
		if engine.IsConnectionError(err) {
			return err
		}

		p.Logger.Err(err).Msg("cannot find alarms")
		return nil
	}
	sentForResources := make(map[string]bool, len(alarms))

	for _, a := range alarms {
		sentForResources[a.EntityID] = true
		e := types.Event{
			EventType:     types.EventTypeEntityUpdated,
			Connector:     a.Value.Connector,
			ConnectorName: a.Value.ConnectorName,
			Component:     a.Value.Component,
			Resource:      a.Value.Resource,
			Timestamp:     types.CpsTime{Time: time.Now()},
			SourceType:    types.SourceTypeResource,
			Output:        "updated component infos",
			Initiator:     types.InitiatorSystem,
		}

		err := p.publishToEngineFIFO(e)
		if err != nil {
			return err
		}
	}

	for _, resource := range resources {
		if !sentForResources[resource] {
			p.Logger.Warn().Str("entity", resource).
				Msg("resource doesn't have opened alarm so no event was fired on component_infos update and context graph won't be updated until new alarm")
		}
	}

	return nil
}

func (p *messageProcessor) publishToEngineFIFO(event types.Event) error {
	body, err := p.Encoder.Encode(event)
	if err != nil {
		p.Logger.Err(err).Msg("cannot encode event")
		return nil
	}
	return errt.NewIOError(p.AmqpPublisher.Publish(
		"",
		canopsis.FIFOQueueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json", // this type is mandatory to avoid bad conversions into Python.
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	))
}
