package main

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"time"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	libpbehavior "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type rpcServerMessageProcessor struct {
	FeaturePrintEventOnError bool
	DbClient                 mongo.DbClient
	PbhService               libpbehavior.Service
	EventManager             libpbehavior.EventManager
	TimezoneConfigProvider   config.TimezoneConfigProvider
	PubChannel               libamqp.Channel
	Decoder                  encoding.Decoder
	Encoder                  encoding.Encoder
	Logger                   zerolog.Logger
}

func (p *rpcServerMessageProcessor) Process(ctx context.Context, d amqp.Delivery) ([]byte, error) {
	msg := d.Body
	var event rpc.PbehaviorEvent
	err := p.Decoder.Decode(msg, &event)
	if err != nil || event.Alarm == nil || event.Entity == nil {
		p.logError(err, "invalid event", msg)

		return p.getErrRpcEvent(errors.New("invalid event")), nil
	}

	if event.Healthcheck {
		return p.getRpcEvent(rpc.PbehaviorResultEvent{
			Alarm:  event.Alarm,
			Entity: event.Entity,
		})
	}

	var pbhEvent *types.Event
	switch event.Type {
	case rpc.PbehaviorEventTypeCreate:
		pbhEvent, err = p.processCreatePbhEvent(ctx, *event.Entity, event.Params)
	case rpc.PbehaviorEventTypeDelete:
		pbhEvent, err = p.processDeletePbhEvent(ctx, *event.Entity, event.Params)
	default:
		p.logError(nil, "invalid event type: "+event.Type, msg)

		return p.getErrRpcEvent(errors.New("invalid event type: " + event.Type)), nil
	}

	if err != nil {
		if engine.IsConnectionError(err) {
			return nil, err
		}

		p.logError(err, "cannot process event", msg)

		return p.getErrRpcEvent(err), nil
	}

	if pbhEvent != nil && d.ReplyTo == "" {
		b, err := p.Encoder.Encode(pbhEvent)
		if err != nil {
			p.logError(err, "cannot encode event", msg)

			return nil, nil
		}

		err = p.PubChannel.PublishWithContext(
			ctx,
			"",
			canopsis.FIFOQueueName,
			false,
			false,
			amqp.Publishing{
				ContentType:  canopsis.JsonContentType,
				Body:         b,
				DeliveryMode: amqp.Persistent,
			},
		)
		if err != nil {
			if engine.IsConnectionError(err) {
				return nil, err
			}

			p.logError(err, "cannot publish event", msg)
		}

		return nil, nil
	}

	resEvent := rpc.PbehaviorResultEvent{
		Alarm:  event.Alarm,
		Entity: event.Entity,
	}
	if pbhEvent != nil {
		resEvent.PbhEvent = *pbhEvent
		resEvent.PbhEvent.Entity = event.Entity
	}

	return p.getRpcEvent(resEvent)
}

func (p *rpcServerMessageProcessor) processCreatePbhEvent(
	ctx context.Context,
	entity types.Entity,
	params rpc.PbehaviorParameters,
) (*types.Event, error) {
	pbehavior, oldPbhIds, err := p.createPbehavior(ctx, params, entity)
	if err != nil {
		return nil, err
	}

	oldPbhIds = append(oldPbhIds, pbehavior.ID)
	resolver, err := p.PbhService.RecomputeByIds(ctx, oldPbhIds)
	if err != nil {
		return nil, fmt.Errorf("pbehavior recompute failed: %w", err)
	}

	p.Logger.Debug().Str("pbehavior", pbehavior.ID).Msg("pbehavior recomputed")

	resolveResult, err := p.getResolveResult(ctx, entity, resolver)
	if err != nil {
		return nil, err
	}

	pbhEvent, err := p.EventManager.GetEvent(resolveResult, entity, datetime.NewCpsTime())
	if err != nil {
		return nil, err
	}

	if pbhEvent.EventType == "" {
		return nil, nil
	}

	if pbhEvent.PbehaviorInfo.ID == pbehavior.ID {
		if params.RuleName != "" {
			pbhEvent.PbehaviorInfo.RuleName = params.RuleName
			pbhEvent.Output = pbhEvent.PbehaviorInfo.GetStepMessage()
		}

		if params.Author != "" {
			pbhEvent.Author = params.Author
			pbhEvent.PbehaviorInfo.Author = params.Author
		}
	}

	return &pbhEvent, nil
}

func (p *rpcServerMessageProcessor) processDeletePbhEvent(
	ctx context.Context,
	entity types.Entity,
	params rpc.PbehaviorParameters,
) (*types.Event, error) {
	pbhId, err := p.deletePbehavior(ctx, params, entity)
	if pbhId == "" || err != nil {
		return nil, err
	}

	resolver, err := p.PbhService.RecomputeByIds(ctx, []string{pbhId})
	if err != nil {
		return nil, fmt.Errorf("pbehavior recompute failed: %w", err)
	}

	p.Logger.Debug().Str("pbehavior", pbhId).Msg("pbehavior removed")

	resolveResult, err := p.getResolveResult(ctx, entity, resolver)
	if err != nil {
		return nil, err
	}

	pbhEvent, err := p.EventManager.GetEvent(resolveResult, entity, datetime.NewCpsTime())
	if err != nil {
		return nil, err
	}

	if pbhEvent.EventType == "" {
		return nil, nil
	}

	if pbhEvent.PbehaviorInfo.ID == "" {
		if params.RuleName != "" {
			prevPbhInfo := entity.PbehaviorInfo
			prevPbhInfo.RuleName = params.RuleName
			pbhEvent.Output = prevPbhInfo.GetStepMessage()
		}

		if params.Author != "" {
			pbhEvent.Author = params.Author
		}
	}

	return &pbhEvent, nil
}

func (p *rpcServerMessageProcessor) createPbehavior(
	ctx context.Context,
	params rpc.PbehaviorParameters,
	entity types.Entity,
) (*libpbehavior.PBehavior, []string, error) {
	typeCollection := p.DbClient.Collection(mongo.PbehaviorTypeMongoCollection)
	err := typeCollection.FindOne(ctx, bson.M{"_id": params.Type}).Err()
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil, fmt.Errorf("pbehavior type not exist: %q", params.Type)
		}

		return nil, nil, fmt.Errorf("cannot get pbehavior type: %w", err)
	}

	reasonCollection := p.DbClient.Collection(mongo.PbehaviorReasonMongoCollection)
	err = reasonCollection.FindOne(ctx, bson.M{"_id": params.Reason}).Err()
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil, fmt.Errorf("pbehavior reason not exist: %q", params.Reason)
		}

		return nil, nil, fmt.Errorf("cannot get pbehavior reason: %w", err)
	}

	now := datetime.NewCpsTime()
	var start, stop datetime.CpsTime
	if params.Tstart != nil && params.Tstop != nil {
		start = *params.Tstart
		stop = *params.Tstop
	} else if params.StartOnTrigger != nil && *params.StartOnTrigger &&
		params.Duration != nil && params.Duration.Value > 0 {
		start = now
		stop = params.Duration.AddTo(now)
	}

	if start.IsZero() {
		return nil, nil, fmt.Errorf("invalid action parameters, tstart with tstop or start_on_trigger with duration must be defined: %+v", params)
	}

	pbehavior := libpbehavior.PBehavior{
		ID:       utils.NewID(),
		Enabled:  true,
		Comments: make([]libpbehavior.Comment, 0),
		Name:     params.Name,
		Reason:   params.Reason,
		RRule:    params.RRule,
		Start:    &start,
		Stop:     &stop,
		Type:     params.Type,
		Color:    params.Color,
		Created:  &now,
		Updated:  &now,
		Origin:   params.Origin,
		Entity:   entity.ID,
		EntityPatternFields: savedpattern.EntityPatternFields{
			EntityPattern: pattern.Entity{
				{
					{
						Field: "_id",
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: entity.ID,
						},
					},
				},
			},
		},
	}

	if params.Comment != "" {
		pbehavior.Comments = append(pbehavior.Comments, libpbehavior.Comment{
			ID:        utils.NewID(),
			Origin:    cmp.Or(params.Author, canopsis.DefaultEventAuthor),
			Timestamp: &now,
			Message:   params.Comment,
		})
	}

	oldPbhIds := make([]string, 0)
	collection := p.DbClient.Collection(mongo.PbehaviorMongoCollection)
	err = p.DbClient.WithTransaction(ctx, func(ctx context.Context) error {
		oldPbhIds = oldPbhIds[:0]
		cursor, err := collection.Find(ctx, bson.M{
			"entity": entity.ID,
			"origin": params.Origin,
		}, options.Find().SetProjection(bson.M{"_id": 1}))
		if err != nil {
			return fmt.Errorf("cannot find old pbehavior: %w", err)
		}

		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			pbh := libpbehavior.PBehavior{}
			err = cursor.Decode(&pbh)
			if err != nil {
				return err
			}

			oldPbhIds = append(oldPbhIds, pbh.ID)
		}

		_, err = collection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": oldPbhIds}})
		if err != nil {
			return fmt.Errorf("cannot delete old pbehavior: %w", err)
		}

		_, err = collection.InsertOne(ctx, pbehavior)
		if err != nil {
			return fmt.Errorf("create new pbehavior failed: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	p.Logger.Info().Str("pbehavior", pbehavior.ID).Msg("create pbehavior")

	return &pbehavior, oldPbhIds, nil
}

func (p *rpcServerMessageProcessor) deletePbehavior(
	ctx context.Context,
	params rpc.PbehaviorParameters,
	entity types.Entity,
) (string, error) {
	collection := p.DbClient.Collection(mongo.PbehaviorMongoCollection)
	var pbehavior libpbehavior.PBehavior
	err := collection.FindOneAndDelete(ctx, bson.M{
		"entity": entity.ID,
		"origin": params.Origin,
	}, options.FindOneAndDelete().SetProjection(bson.M{"_id": 1})).Decode(&pbehavior)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return "", nil
		}

		return "", fmt.Errorf("delete pbehavior failed: %w", err)
	}

	p.Logger.Info().Str("pbehavior", pbehavior.ID).Msg("delete pbehavior")

	return pbehavior.ID, nil
}

func (p *rpcServerMessageProcessor) getResolveResult(
	ctx context.Context,
	entity types.Entity,
	resolver libpbehavior.ComputedEntityTypeResolver,
) (libpbehavior.ResolveResult, error) {
	location := p.TimezoneConfigProvider.Get().Location
	now := time.Now().In(location)
	resolveResult, err := resolver.Resolve(ctx, entity, now)
	if err != nil {
		return libpbehavior.ResolveResult{}, fmt.Errorf("resolve an entity failed: %w", err)
	}

	return resolveResult, nil
}

func (p *rpcServerMessageProcessor) getErrRpcEvent(err error) []byte {
	msg, _ := p.getRpcEvent(rpc.PbehaviorResultEvent{Error: &rpc.Error{Error: err}})
	return msg
}

func (p *rpcServerMessageProcessor) getRpcEvent(event rpc.PbehaviorResultEvent) ([]byte, error) {
	msg, err := p.Encoder.Encode(event)
	if err != nil {
		p.Logger.Err(err).Msg("cannot encode event")
		return nil, nil
	}

	return msg, nil
}

func (p *rpcServerMessageProcessor) logError(err error, errMsg string, msg []byte) {
	if p.FeaturePrintEventOnError {
		p.Logger.Err(err).Str("event", string(msg)).Msg(errMsg)
	} else {
		p.Logger.Err(err).Msg(errMsg)
	}
}
