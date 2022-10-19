package main

import (
	"context"
	"fmt"
	"runtime/trace"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/streadway/amqp"
)

type EngineFIFO struct {
	canopsis.DefaultEngine
	Options    Options
	References References
	ctlChan    chan struct{}
}

func (e *EngineFIFO) ConsumerChan() (<-chan amqp.Delivery, error) {
	msgs, err := e.References.ChannelSub.Consume(
		e.Options.ConsumeFromQueue, // queue
		canopsis.FIFOConsumerName,  // consumer
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

func (e *EngineFIFO) AckConsumerChan() (<-chan amqp.Delivery, error) {
	msgs, err := e.References.AckChanSub.Consume(
		canopsis.FIFOAckQueueName, // queue
		"fifo_ack",                // consumer
		false,                     // auto-ack
		false,                     // exclusive
		false,                     // no-local
		false,                     // no-wait
		nil,                       // args
	)

	if err != nil {
		return nil, fmt.Errorf("consume error: %v", err)
	}

	return msgs, err
}

func (e *EngineFIFO) Initialize(ctx context.Context) error {
	e.References.Scheduler.Start(ctx)
	e.Logger().Info().Msg("FIFO engine started")

	go e.ackManager(ctx)

	return nil
}

func (e *EngineFIFO) PeriodicalProcess(ctx context.Context) {
	select {
	case <-e.ctlChan:
		e.Logger().Error().Msg("exit engine signal received")
		e.AskStop(canopsis.ExitEngine)
	default:
	}

}

func (e *EngineFIFO) WorkerProcess(parentCtx context.Context, msg amqp.Delivery) {
	ctx, task := trace.NewTask(parentCtx, "fifo.WorkerProcess")
	defer task.End()

	trace.Logf(ctx, "event_size", "%d", len(msg.Body))

	var event types.Event
	err := e.References.JSONDecoder.Decode(msg.Body, &event)
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

	event.Format()
	e.References.StatsSender.Add(time.Now().Unix(), true)

	err = event.InjectExtraInfos(msg.Body)
	if err != nil {
		e.processWorkerError(err, msg)
		return
	}

	e.Logger().Debug().Str("event", fmt.Sprintf("%+v", event)).Msg("sent to scheduler")
	err = e.References.Scheduler.ProcessEvent(ctx, event)

	if err != nil {
		e.processWorkerError(err, msg)
		return
	}

	e.DefaultEngine.AckMessage(msg)
}

// processWorkerError read errors, ack amqp messages and stop the engine if needed
func (e *EngineFIFO) processWorkerError(err error, msg amqp.Delivery) {
	const msgError = "event processing error"

	if !e.Options.PrintEventOnError {
		e.Logger().Error().Err(err).Msg(msgError)
	} else {
		e.Logger().Error().Err(err).Str("event", string(msg.Body)).Msg(msgError)
	}

	e.DefaultEngine.ProcessWorkerError(err, msg)
}

func (e *EngineFIFO) ackManager(ctx context.Context) {
	exitFIFO := struct{}{}
	ackChan, err := e.AckConsumerChan()
	if err != nil {
		e.Logger().Debug().Err(err).Msg("ack consumer chan error")
		e.ctlChan <- exitFIFO
		return
	}
	e.Logger().Debug().Msg("Ack consumer initialized")

	for msg := range ackChan {
		e.Logger().Debug().Msgf("Ack received %s", string(msg.Body))
		var event types.Event
		err := e.References.JSONDecoder.Decode(msg.Body, &event)

		if err != nil {
			e.Logger().Err(err).
				Str("lockID", event.GetLockID()).
				Msg("JSON decode")
			err = e.confirmAck(msg)
			if err != nil {
				e.Logger().Err(err).
					Str("lockID", event.GetLockID()).
					Msg("confirm ack")
			}
			continue
		}

		if err := e.References.Scheduler.AckEvent(ctx, event); err != nil {
			e.Logger().Err(err).
				Str("lockID", event.GetLockID()).
				Msg("Error on acking message")
		}

		err = e.confirmAck(msg)
		if err != nil {
			e.Logger().Err(err).
				Str("lockID", event.GetLockID()).
				Msg("AckChan ack")
		}
	}
}

func (e *EngineFIFO) confirmAck(msg amqp.Delivery) error {
	return e.References.AckChanSub.Ack(msg.DeliveryTag, false)
}

func (e *EngineFIFO) GetRunInfo() engine.RunInfo {
	return engine.RunInfo{
		Name:         "engine-fifo",
		ConsumeQueue: e.Options.ConsumeFromQueue,
		PublishQueue: e.Options.PublishToQueue,
	}
}
