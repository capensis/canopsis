package main

import (
	"context"
	"errors"
	"fmt"
	"runtime/trace"
	"time"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	libpbehavior "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/errt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
)

const (
	resolveTimeout                     = time.Millisecond * 100
	resolveDeadlineExceededErrInterval = time.Second * 30
)

type messageProcessor struct {
	PbhService               libpbehavior.EntityTypeResolver
	TimezoneConfigProvider   config.TimezoneConfigProvider
	FeaturePrintEventOnError bool
	Encoder                  encoding.Encoder
	Decoder                  encoding.Decoder
	CreatePbehaviorProcessor createPbehaviorMessageProcessor
	ChannelPub               libamqp.Channel
	Logger                   zerolog.Logger
	// resolveDeadlineExceededAt contains time of last logged deadline exceeded error.
	// The error is logged only once in resolveDeadlineExceededErrInterval.
	resolveDeadlineExceededAt time.Time
}

func (p *messageProcessor) Process(ctx context.Context, d amqp.Delivery) ([]byte, error) {
	msg := d.Body
	ctx, task := trace.NewTask(ctx, "pbehavior.MessageProcessor")
	defer task.End()
	trace.Logf(ctx, "event_size", "%d", len(msg))
	var err error
	var event types.Event
	trace.WithRegion(ctx, "decode-event", func() {
		err = p.Decoder.Decode(msg, &event)
	})
	if err != nil {
		p.logError(err, "cannot decode event", msg)
		return nil, nil
	}

	p.Logger.Debug().Msgf("unmarshaled: %+v", event)
	trace.Log(ctx, "event.event_type", event.EventType)
	trace.Log(ctx, "event.timestamp", event.Timestamp.String())
	trace.Log(ctx, "event.source_type", event.SourceType)
	trace.Log(ctx, "event.connector", event.Connector)
	trace.Log(ctx, "event.connector_name", event.ConnectorName)
	trace.Log(ctx, "event.component", event.Component)
	trace.Log(ctx, "event.resource", event.Resource)

	if event.Entity != nil {
		if event.EventType == types.EventTypePbhCreate {
			err := p.processPbhCreateEvent(ctx, event, msg)
			if err != nil {
				return nil, err
			}
		} else if !event.IsPbehaviorEvent() {
			pbhInfo, err := p.processEvent(ctx, event, msg)
			if err != nil {
				return nil, err
			}

			event.PbehaviorInfo = pbhInfo
		}
	}

	var body []byte
	trace.WithRegion(ctx, "encode-event", func() {
		body, err = p.Encoder.Encode(event)
	})
	if err != nil {
		p.logError(err, "cannot encode event", msg)
		return nil, nil
	}

	return body, nil
}

// processEvent tries to resolve pbehavior type for entity.
// It logs error and sends event to next engine on fail and
// pbehavior type will be resolved in periodical worker.
func (p *messageProcessor) processEvent(ctx context.Context, event types.Event, msg []byte) (types.PbehaviorInfo, error) {
	// Skip resolve if the entity is already in pbehavior.
	if !event.Entity.PbehaviorInfo.IsDefaultActive() {
		return event.Entity.PbehaviorInfo, nil
	}
	// Resolve type in case if the entity is new.
	ctx, cancel := context.WithTimeout(ctx, resolveTimeout)
	defer cancel()
	now := time.Now().In(p.TimezoneConfigProvider.Get().Location)

	resolveResult, err := p.PbhService.Resolve(ctx, *event.Entity, now)
	if err == nil {
		if !p.resolveDeadlineExceededAt.IsZero() {
			p.resolveDeadlineExceededAt = time.Time{}
			p.Logger.Info().Msg("entity resolving has been fixed")
		}

		return libpbehavior.NewPBehaviorInfo(types.CpsTime{Time: now}, resolveResult), nil
	}

	if errors.Is(err, context.DeadlineExceeded) {
		if p.resolveDeadlineExceededAt.IsZero() || time.Since(p.resolveDeadlineExceededAt) > resolveDeadlineExceededErrInterval {
			p.resolveDeadlineExceededAt = now
			p.Logger.Err(err).
				Str("entity", event.Entity.ID).
				Msg("resolve an entity too long")
		}

		return types.PbehaviorInfo{}, nil
	}

	if redis.IsConnectionError(err) {
		return types.PbehaviorInfo{}, err
	}

	p.logError(err, "resolve an entity failed", msg)
	return types.PbehaviorInfo{}, nil
}

func (p *messageProcessor) processPbhCreateEvent(ctx context.Context, event types.Event, msg []byte) error {
	params := types.ActionPBehaviorParameters{}
	err := p.Decoder.Decode([]byte(event.PbhParameters), &params)
	if err != nil {
		p.logError(err, "invalid params for create pbehavior", msg)
		return nil
	}

	pbhEvent, err := p.CreatePbehaviorProcessor.Process(ctx, event.Alarm, event.Entity,
		params, msg)
	if err != nil {
		if redis.IsConnectionError(err) {
			return err
		}

		p.logError(err, "cannot process event", msg)
		return nil
	}

	if pbhEvent != nil {
		err := p.publishTo(*pbhEvent, canopsis.FIFOQueueName)
		if err != nil {
			return fmt.Errorf("cannot publish pbh event: %w", err)
		}
	}

	return nil
}

func (p *messageProcessor) publishTo(event types.Event, queue string) error {
	body, err := p.Encoder.Encode(event)
	if err != nil {
		return fmt.Errorf("error while encoding event: %w", err)
	}

	return errt.NewIOError(p.ChannelPub.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType:  canopsis.JsonContentType,
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	))
}

func (p *messageProcessor) logError(err error, errMsg string, msg []byte) {
	if p.FeaturePrintEventOnError {
		p.Logger.Err(err).Str("event", string(msg)).Msg(errMsg)
	} else {
		p.Logger.Err(err).Msg(errMsg)
	}
}
