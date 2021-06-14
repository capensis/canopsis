package main

import (
	"context"
	"errors"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	libpbehavior "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/bsm/redislock"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"time"
)

type rpcServerMessageProcessor struct {
	FeaturePrintEventOnError bool
	Processor                createPbehaviorMessageProcessor
	Decoder                  encoding.Decoder
	Encoder                  encoding.Encoder
	Logger                   zerolog.Logger
}

func (p *rpcServerMessageProcessor) Process(ctx context.Context, d amqp.Delivery) ([]byte, error) {
	msg := d.Body
	var event types.RPCPBehaviorEvent
	err := p.Decoder.Decode(msg, &event)
	if err != nil || event.Alarm == nil {
		p.logError(err, "Pbehavior RPC server: invalid event", msg)

		return p.getErrRpcEvent(errors.New("invalid event")), nil
	}

	pbhEvent, err := p.Processor.Process(
		ctx,
		event.Alarm,
		event.Entity,
		event.Params,
		msg,
		"Pbehavior RPC server",
	)
	if err != nil {
		if engine.IsConnectionError(err) {
			return nil, err
		}

		p.logError(err, "cannot process event", msg)
		return p.getErrRpcEvent(err), nil
	}

	if pbhEvent == nil {
		pbhEvent = &types.Event{}
	}

	return p.getRpcEvent(types.RPCPBehaviorResultEvent{
		Alarm:    event.Alarm,
		PbhEvent: *pbhEvent,
		Error:    nil,
	})
}

type createPbehaviorMessageProcessor struct {
	FeaturePrintEventOnError bool
	DbClient                 mongo.DbClient
	LockerClient             redis.LockClient
	Store                    redis.Store
	PbhService               libpbehavior.Service
	EventManager             libpbehavior.EventManager
	AlarmAdapter             alarm.Adapter
	TimezoneConfigProvider   config.TimezoneConfigProvider
	Logger                   zerolog.Logger
}

func (p *createPbehaviorMessageProcessor) Process(
	ctx context.Context,
	alarm *types.Alarm,
	entity *types.Entity,
	params types.ActionPBehaviorParameters,
	msg []byte,
	logPrefix string,
) (*types.Event, error) {
	pbehavior, err := p.createPbehavior(ctx, params, entity)
	if err != nil {
		p.logError(err, fmt.Sprintf("%s: failed to create pbehavior", logPrefix), msg)

		return nil, err
	}

	if pbehavior == nil {
		p.logError(err, fmt.Sprintf("%s: createPbehavior returned err = nil, but pbehavior is empty", logPrefix), msg)

		return nil, errors.New("pbehavior is empty")
	}

	err = p.recomputePbehavior(ctx, pbehavior.ID)
	if err != nil {
		p.logError(err, fmt.Sprintf("%s: failed to recompute pbehaviors", logPrefix), msg)

		return nil, err
	}

	resolveResult, err := p.getResolveResult(ctx, entity)
	if err != nil {
		p.logError(err, fmt.Sprintf("%s: failed to resolve pbehavior for an entity", logPrefix), msg)

		return nil, err
	}

	if alarm == nil {
		alarms := make([]types.Alarm, 0)
		err := p.AlarmAdapter.GetOpenedAlarmsByIDs([]string{entity.ID}, &alarms)
		if err != nil {
			p.logError(err, fmt.Sprintf("%s: failed to find alarm", logPrefix), msg)

			return nil, err
		}

		if len(alarms) == 0 {
			return nil, nil
		}

		alarm = &alarms[0]
	}

	pbhEvent := p.EventManager.GetEvent(resolveResult, *alarm, time.Now())
	if pbhEvent.EventType == "" {
		return nil, nil
	}

	pbhEvent.Entity = entity

	return &pbhEvent, nil
}

func (p *createPbehaviorMessageProcessor) createPbehavior(
	ctx context.Context,
	params types.ActionPBehaviorParameters,
	entity *types.Entity,
) (*libpbehavior.PBehavior, error) {
	typeCollection := p.DbClient.Collection(mongo.PbehaviorTypeMongoCollection)
	res := typeCollection.FindOne(ctx, bson.M{"_id": params.Type})
	if err := res.Err(); err != nil {
		if err == mongodriver.ErrNoDocuments {
			p.Logger.Error().Str("type", params.Type).Msg("Pbehavior RPC server: pbehavior type not exist")
			return nil, nil
		} else {
			p.Logger.Err(err).Msg("Pbehavior RPC server: cannot get pbehavior type")
			return nil, err
		}
	}

	reasonCollection := p.DbClient.Collection(mongo.PbehaviorReasonMongoCollection)
	res = reasonCollection.FindOne(ctx, bson.M{"_id": params.Reason})
	if err := res.Err(); err != nil {
		if err == mongodriver.ErrNoDocuments {
			p.Logger.Error().Str("reason", params.Reason).Msg("Worker process: pbehavior reason not exist")
			return nil, nil
		} else {
			p.Logger.Err(err).Msg("Pbehavior RPC server: cannot get pbehavior reason")
			return nil, err
		}
	}

	now := types.NewCpsTime(time.Now().Unix())
	var start, stop *types.CpsTime
	if params.Tstart != nil && params.Tstop != nil {
		start = &types.CpsTime{Time: time.Unix(*params.Tstart, 0)}
		stop = &types.CpsTime{Time: time.Unix(*params.Tstop, 0)}
	} else if params.StartOnTrigger != nil && *params.StartOnTrigger &&
		params.Duration != nil && params.Duration.Seconds > 0 {
		now := time.Now()
		start = &types.CpsTime{Time: now}
		stop = &types.CpsTime{Time: now.Add(time.Duration(params.Duration.Seconds) * time.Second)}
	}

	if start == nil {
		err := fmt.Errorf("invalid action parameters: %+v", params)
		p.Logger.Err(err).Msg("Pbehavior RPC server: cannot create pbehavior")
		return nil, nil
	}

	pbehavior := libpbehavior.PBehavior{
		ID:      utils.NewID(),
		Author:  params.Author,
		Enabled: true,
		Filter:  fmt.Sprintf(`{"_id": "%s"}`, entity.ID),
		Name:    params.Name,
		Reason:  params.Reason,
		RRule:   params.RRule,
		Start:   start,
		Stop:    stop,
		Type:    params.Type,
		Created: now,
		Updated: now,
	}

	collection := p.DbClient.Collection(mongo.PbehaviorMongoCollection)
	_, err := collection.InsertOne(ctx, pbehavior)
	if err != nil {
		p.Logger.Err(err).Msg("Pbehavior RPC server: create new pbehavior failed!")
		return nil, err
	}

	p.Logger.Info().Str("pbehavior", pbehavior.ID).Msg("Pbehavior RPC server: create pbehavior")
	return &pbehavior, nil
}

func (p *createPbehaviorMessageProcessor) recomputePbehavior(ctx context.Context, pbehaviorID string) error {
	ok, err := p.Store.Restore(ctx, p.PbhService)
	if err != nil || !ok {
		if err == nil {
			err = fmt.Errorf("pbehavior intervals are not computed, cache is empty")
		}
		p.Logger.Err(err).Msg("Pbehavior RPC server: get pbehavior's frames from redis failed! Skip periodical process")
		return err
	}

	computeLock, err := p.LockerClient.Obtain(
		ctx,
		redis.RecomputeLockKey,
		redis.RecomputeLockDuration,
		&redislock.Options{
			RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(100*time.Millisecond), 10),
		},
	)
	if err != nil {
		p.Logger.Err(err).Msg("Pbehavior RPC server: obtain redlock failed! Skip recompute")
		return err
	}

	defer func() {
		if computeLock != nil {
			err := computeLock.Release(ctx)
			if err != nil && err != redislock.ErrLockNotHeld {
				p.Logger.Warn().Err(err).Msg("Pbehavior RPC server: failed to manually release compute-lock, the lock will be released by ttl")
			}
		}
	}()

	err = p.PbhService.Recompute(ctx, pbehaviorID)

	if err != nil {
		p.Logger.Err(err).Msgf("Pbehavior RPC server: pbehavior recompute failed!")
		return err
	}

	err = p.Store.Save(ctx, p.PbhService)
	if err != nil {
		p.Logger.Err(err).Msg("Pbehavior RPC server: save pbehavior's frames to redis failed! The data might be inconsistent")
		return err
	}

	p.Logger.Debug().Str("pbehavior", pbehaviorID).Msg("Pbehavior RPC server: pbehavior recomputed")

	return nil
}

func (p *createPbehaviorMessageProcessor) getResolveResult(ctx context.Context, entity *types.Entity) (libpbehavior.ResolveResult, error) {
	location := p.TimezoneConfigProvider.Get().Location
	now := time.Now().In(location)
	resolveResult, err := p.PbhService.Resolve(ctx, entity, now)
	if err != nil {
		p.Logger.Err(err).Str("entity_id", entity.ID).Msg("Pbehavior RPC server: resolve an entity failed")
		return libpbehavior.ResolveResult{}, err
	}

	return resolveResult, nil
}

func (p *createPbehaviorMessageProcessor) logError(err error, errMsg string, msg []byte) {
	if p.FeaturePrintEventOnError {
		p.Logger.Err(err).Str("event", string(msg)).Msg(errMsg)
	} else {
		p.Logger.Err(err).Msg(errMsg)
	}
}

func (p *rpcServerMessageProcessor) getErrRpcEvent(err error) []byte {
	msg, _ := p.getRpcEvent(types.RPCPBehaviorResultEvent{Error: &types.RPCError{Error: err}})
	return msg
}

func (p *rpcServerMessageProcessor) getRpcEvent(event types.RPCPBehaviorResultEvent) ([]byte, error) {
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
