package main

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/engine"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/errt"
	"git.canopsis.net/canopsis/go-engines/lib/utils"
	"github.com/streadway/amqp"
	"runtime/trace"
)

type EngineChe struct {
	canopsis.DefaultEngine
	Options    Options
	References References
}

func (e *EngineChe) ConsumerChan() (<-chan amqp.Delivery, error) {
	if e.Options.Purge {
		e.Logger().Info().Msg("Purging queue")
		_, err := e.References.ChannelSub.QueuePurge(e.Options.ConsumeFromQueue, false)
		if err != nil {
			return nil, fmt.Errorf("error while purging queue: %v", err)
		}
	}
	msgs, err := e.References.ChannelSub.Consume(
		e.Options.ConsumeFromQueue, // queue
		canopsis.CheConsumerName,   // consumer
		false,                      // auto-ack
		false,                      // exclusive
		false,                      // no-local
		false,                      // no-wait
		nil,                        // args
	)

	if err != nil {
		return nil, fmt.Errorf("consume error: %v", err)
	}

	return msgs, err
}

func (e *EngineChe) Initialize() error {
	e.Logger().Debug().Msg("Loading event filter data sources")
	err := e.References.EventFilterService.LoadDataSourceFactories(
		e.References.EnrichmentCenter,
		e.References.EnrichFields,
		e.Options.DataSourceDirectory)
	if err != nil {
		return fmt.Errorf("unable to load data sources: %v", err)
	}

	e.Logger().Debug().Msg("Loading event filter rules")
	err = e.References.EventFilterService.LoadRules()
	if err != nil {
		return fmt.Errorf("unable to load rules: %v", err)
	}

	if e.References.Config.Alarm.OutputLength < 1 {
		e.Logger().Warn().Msg("OutputLength value is not set or less than 1: the event's output and long_output won't be truncated")
	}

	return nil
}

// Stop force bulk writes before stopping
func (e *EngineChe) Stop() int {
	return e.DefaultEngine.Stop()
}

func (e *EngineChe) PeriodicalProcess() {
	e.Logger().Debug().Msg("Loading event filter rules")
	err := e.References.EventFilterService.LoadRules()
	if err != nil {
		e.Logger().Error().Err(err).Msg("unable to load rules")
	}

	e.Logger().Debug().Msg("Loading watchers")
	err = e.References.EnrichmentCenter.LoadWatchers()
	if err != nil {
		e.Logger().Error().Err(err).Msg("unable to load watchers")
	}
}

// processWorkerError read errors, ack amqp messages and stop the engine if needed
func (e *EngineChe) processWorkerError(err error, msg amqp.Delivery) {
	_, isDropError := err.(eventfilter.DropError)
	if isDropError {
		e.Logger().Info().Str("event", string(msg.Body)).Msg("event dropped by event filter")
	} else if !e.Options.PrintEventOnError {
		e.Logger().Error().Err(err).Msg("event processing error")
	} else {
		e.Logger().Error().Err(err).Str("event", string(msg.Body)).Msg("event processing error")
	}

	e.DefaultEngine.ProcessWorkerError(err, msg)
}

func (e *EngineChe) WorkerProcess(msg amqp.Delivery) {
	fifoAck := true

	defer func() {
		if fifoAck {
			err := e.unlockAlarm(msg.Body)
			if err != nil {
				e.Logger().Error().Err(err).Msg("Failed to send an fifoAck message to the engine_fifo, the alarm will be unlocked automatically by the engine_fifo")
			}
		}
	}()

	ctx, task := trace.NewTask(context.Background(), "che.WorkerProcess")
	defer task.End()

	trace.Logf(ctx, "event_size", "%d", len(msg.Body))

	var event types.Event
	err := e.References.JSONDecoder.Decode(msg.Body, &event)
	if err != nil {
		e.processWorkerError(err, msg)
		return
	}

	err = event.InjectExtraInfos(msg.Body)
	if err != nil {
		e.processWorkerError(err, msg)
		return
	}

	e.Logger().Debug().Msgf("valid input event: %v", string(msg.Body))
	trace.Log(ctx, "event.event_type", event.EventType)
	trace.Log(ctx, "event.timestamp", event.Timestamp.String())
	trace.Log(ctx, "event.source_type", event.SourceType)
	trace.Log(ctx, "event.connector", event.Connector)
	trace.Log(ctx, "event.connector_name", event.ConnectorName)
	trace.Log(ctx, "event.component", event.Component)
	trace.Log(ctx, "event.resource", event.Resource)

	err = event.IsValid()
	if err != nil {
		e.processWorkerError(err, msg)
		return
	}

	event.Output = utils.TruncateString(event.Output, e.References.Config.Alarm.OutputLength)
	event.LongOutput = utils.TruncateString(event.LongOutput, e.References.Config.Alarm.OutputLength)

	if e.Options.FeatureEventProcessing {
		var report eventfilter.Report
		event, report, err = e.References.EventFilterService.ProcessEvent(event)
		if err != nil {
			e.processWorkerError(err, msg)
			return
		}

		if report.EntityUpdated && event.Entity != nil {
			*event.Entity = e.References.EnrichmentCenter.Update(*event.Entity)
		}
	}

	// Enrich the event with the entity and create the context if this has not
	// been done by the event filter.
	// FIXME: this could cause issues if the event's entity is defined but is
	// not added to the context.
	if event.Entity == nil && (e.Options.FeatureContextCreation || e.Options.FeatureContextEnrich) {
		event, err = e.EnrichContextFromEvent(event)
		if err != nil {
			e.processWorkerError(err, msg)
			return
		}
	}

	if event.Entity != nil && event.IsContextable() {
		err := e.References.EnrichmentCenter.EnrichResourceInfoWithComponentInfo(&event, event.Entity)
		if err == nil {
			err = e.References.EnrichmentCenter.Flush()
		}

		if err != nil {
			e.Logger().Warn().Err(err).Msg("can't update resource's info")
		}
	}

	if e.Options.FeatureEventProcessing {
		event.Format()

		fifoAck, err = e.enrichAndSend(event)
		if err != nil {
			e.processWorkerError(err, msg)
			return
		}
	}

	e.DefaultEngine.AckMessage(msg)
}

// EnrichContextFromEvent creates the entity corresponding to an event, and
// enriches it with the event's extra infos.
func (e *EngineChe) EnrichContextFromEvent(event types.Event) (types.Event, error) {
	if e.Options.FeatureContextCreation && event.IsContextable() {
		event.Entity = e.References.EnrichmentCenter.Handle(event, e.References.EnrichFields)
	}
	return event, nil
}

// enrichAndSend enriches an event with its entity and alarm, and pushes it to
// the new queue
func (e *EngineChe) enrichAndSend(event types.Event) (bool, error) {
	var err error

	// The entity should be set by the enrichment center, as long as the
	// FeatureContextCreation flag is true.
	if event.Entity == nil {
		event.Entity, err = e.References.EnrichmentCenter.Get(event)
		if err != nil {
			return true, err
		}
	}

	bevent, err := e.References.JSONEncoder.Encode(event)
	if err != nil {
		return true, err
	}

	e.Logger().Debug().Msgf("output event about to be published: %+v", string(bevent))
	return e.publishToNext(bevent)
}

// publishToNext publish a received amqp message to another queue
func (e *EngineChe) publishToNext(event []byte) (bool, error) {
	if e.Options.PublishToQueue == "void" {
		return true, nil
	}

	err := errt.NewIOError(e.References.ChannelPub.Publish(
		"",
		e.Options.PublishToQueue,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json", // this type is mandatory to avoid bad conversions into Python.
			Body:         event,
			DeliveryMode: amqp.Persistent,
		},
	))

	if err != nil {
		return true, err
	} else {
		return false, err
	}
}

// send ack message to the fifo ack in order to engine_fifo unlock an alarm
func (e *EngineChe) unlockAlarm(event []byte) error {
	return errt.NewIOError(e.References.ChannelPub.Publish(
		"",
		canopsis.FIFOAckQueueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json", // this type is mandatory to avoid bad conversions into Python.
			Body:         event,
			DeliveryMode: amqp.Persistent,
		},
	))
}

func (e *EngineChe) GetRunInfo() engine.RunInfo {
	publishQueue := e.Options.PublishToQueue
	if publishQueue == "void" {
		publishQueue = ""
	}

	return engine.RunInfo{
		Name:         "engine-che",
		ConsumeQueue: e.Options.ConsumeFromQueue,
		PublishQueue: publishQueue,
	}
}
