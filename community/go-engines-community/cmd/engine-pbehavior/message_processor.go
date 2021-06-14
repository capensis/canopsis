package main

import (
	"context"
	"fmt"
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
	"runtime/trace"
	"time"
)

type messageProcessor struct {
	Store                    redis.Store
	PbhService               libpbehavior.Service
	TimezoneConfigProvider   config.TimezoneConfigProvider
	FeaturePrintEventOnError bool
	Encoder                  encoding.Encoder
	Decoder                  encoding.Decoder
	CreatePbehaviroProcessor createPbehaviorMessageProcessor
	ChannelPub               libamqp.Channel
	Logger                   zerolog.Logger
}

func (p *messageProcessor) Process(ctx context.Context, d amqp.Delivery) ([]byte, error) {
	p.Logger.Debug().Msg("Process message")

	msg := d.Body
	ctx, task := trace.NewTask(context.Background(), "pbehavior.MessageProcessor")
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

	if event.Entity != nil && event.EventType == types.EventTypePbhCreate {
		params := types.ActionPBehaviorParameters{}
		err := p.Decoder.Decode([]byte(event.PbhParameters), &params)
		if err != nil {
			p.logError(err, "Message processor: invalid params for create pbehavior", msg)
			return nil, nil
		}

		pbhEvent, err := p.CreatePbehaviroProcessor.Process(ctx, event.Alarm, event.Entity,
			params, msg, "Message processor")
		if err != nil {
			if redis.IsConnectionError(err) {
				return nil, err
			}

			p.logError(err, "Message processor: cannot process event", msg)
			return nil, nil
		}

		if pbhEvent != nil {
			err := p.publishTo(*pbhEvent, canopsis.FIFOQueueName)
			if err != nil {
				p.logError(err, "Message processor: cannot publish pbh event", msg)
				return nil, err
			}
		}
	} else if event.Entity != nil && !event.IsPbehaviorEvent() {
		ok, err := p.Store.Restore(ctx, p.PbhService)
		if err != nil || !ok {
			if err == nil {
				err = fmt.Errorf("pbehavior intervals are not computed, cache is empty")
			}
			p.logError(err, "Message processor: get pbehavior's frames from redis failed! Skip periodical process", msg)
			return nil, err
		}

		now := time.Now().In(p.TimezoneConfigProvider.Get().Location)
		resolveResult, err := p.PbhService.Resolve(ctx, event.Entity, now)
		if err != nil {
			if redis.IsConnectionError(err) {
				return nil, err
			}

			p.Logger.Err(err).Str("entity_id", event.Entity.ID).Msg("Message processor: resolve an entity failed!")
			return nil, nil
		}

		event.PbehaviorInfo = libpbehavior.NewPBehaviorInfo(resolveResult)
	}

	// Encode and publish the event to the next engine
	var bevent []byte
	trace.WithRegion(ctx, "encode-event", func() {
		bevent, err = p.Encoder.Encode(event)
	})
	if err != nil {
		p.logError(err, "cannot encode event", msg)
		return nil, nil
	}

	return bevent, nil
}

func (p *messageProcessor) publishTo(event types.Event, queue string) error {
	bevent, err := p.Encoder.Encode(event)
	if err != nil {
		return fmt.Errorf("publishTo(): error while encoding event %+v", err)
	}

	return errt.NewIOError(p.ChannelPub.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json", // this type is mandatory to avoid bad conversions into Python.
			Body:         bevent,
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
