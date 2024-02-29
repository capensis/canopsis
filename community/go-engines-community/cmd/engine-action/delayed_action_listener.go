package main

import (
	"context"
	"fmt"
	"time"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/action"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type delayedScenarioListener struct {
	PeriodicalInterval     time.Duration
	DelayedScenarioManager action.DelayedScenarioManager
	AmqpChannel            libamqp.Channel
	Queue                  string
	Encoder                encoding.Encoder
	Logger                 zerolog.Logger
}

func (l *delayedScenarioListener) Listen(ctx context.Context, ch <-chan action.DelayedScenarioTask) {
	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-ch:
			if !ok {
				return
			}

			err := l.publishRunDelayedScenarioEvent(ctx, task)
			if err != nil {
				l.Logger.Err(err).Msg("cannot send run delayed scenario event")
				continue
			}

			l.Logger.Debug().Str("scenario", task.Scenario.ID).Str("alarm", task.Alarm.ID).Msg("send run delayed scenario event")
		}
	}
}

func (l *delayedScenarioListener) publishRunDelayedScenarioEvent(
	ctx context.Context,
	task action.DelayedScenarioTask,
) error {
	b, err := l.Encoder.Encode(task.AdditionalData)
	if err != nil {
		return fmt.Errorf("cannot encode event: %w", err)
	}
	event := types.Event{
		EventType:     types.EventTypeRunDelayedScenario,
		Connector:     canopsis.ActionConnector,
		ConnectorName: canopsis.ActionConnector,
		Component:     task.Alarm.Value.Component,
		Resource:      task.Alarm.Value.Resource,
		SourceType:    task.Entity.Type,
		Timestamp:     datetime.NewCpsTime(),
		Output:        "run delayed scenario",
		Author:        canopsis.DefaultEventAuthor,
		Initiator:     types.InitiatorSystem,

		DelayedScenarioID:   task.Scenario.ID,
		DelayedScenarioData: string(b),
	}

	body, err := l.Encoder.Encode(event)
	if err != nil {
		return fmt.Errorf("cannot encode event: %w", err)
	}

	err = l.AmqpChannel.PublishWithContext(
		ctx,
		"",
		l.Queue,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		l.Logger.Err(err).Msg("cannot publish message to amqp channel")
		return err
	}

	return nil
}
