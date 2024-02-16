package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	libevent "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/db"
	libpbehavior "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const maxLogIds = 20

type recomputeMessageProcessor struct {
	FeaturePrintEventOnError bool
	PbhService               libpbehavior.Service
	PbehaviorCollection      mongo.DbCollection
	EntityCollection         mongo.DbCollection
	EventGenerator           libevent.Generator
	EventManager             libpbehavior.EventManager
	Encoder                  encoding.Encoder
	Decoder                  encoding.Decoder
	Publisher                libamqp.Publisher
	Exchange, Queue          string
	Logger                   zerolog.Logger
}

func (p *recomputeMessageProcessor) Process(ctx context.Context, d amqp.Delivery) ([]byte, error) {
	msg := d.Body
	var event rpc.PbehaviorRecomputeEvent
	err := p.Decoder.Decode(msg, &event)
	if err != nil {
		p.logError(err, "invalid event", msg)

		return nil, nil
	}

	err = p.computePbehaviors(ctx, event.Ids)
	if err != nil {
		if engine.IsConnectionError(err) {
			return nil, err
		}

		p.logError(err, "cannot process event", msg)
		return nil, nil
	}

	return nil, nil
}

func (p *recomputeMessageProcessor) computePbehaviors(ctx context.Context, ids []string) error {
	var resolver libpbehavior.ComputedEntityTypeResolver
	var err error
	if len(ids) == 0 {
		resolver, err = p.PbhService.Recompute(ctx)
	} else {
		resolver, err = p.PbhService.RecomputeByIds(ctx, ids)
	}
	if err != nil {
		return err
	}

	if len(ids) == 0 {
		p.Logger.Info().Msg("all pbehaviors recomputed")
	} else if len(ids) <= maxLogIds {
		p.Logger.Info().Strs("pbehaviors", ids).Msg("pbehaviors recomputed")
	} else {
		p.Logger.Info().
			Strs("first_pbehaviors", ids[:maxLogIds]).
			Int("pbehaviors", len(ids)).
			Msg("pbehaviors recomputed")
	}

	excludeIds := make([]string, 0)
	for _, id := range ids {
		excludeIds, err = p.updateAlarms(ctx, id, excludeIds, resolver)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *recomputeMessageProcessor) updateAlarms(
	ctx context.Context,
	id string,
	excludeIds []string,
	resolver libpbehavior.ComputedEntityTypeResolver,
) ([]string, error) {
	pbehavior := libpbehavior.PBehavior{}
	err := p.PbehaviorCollection.FindOne(ctx, bson.M{"_id": id},
		options.FindOne().SetProjection(bson.M{
			"entity_pattern": 1,
		})).Decode(&pbehavior)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return excludeIds, nil
		}

		return excludeIds, err
	}

	matchByPattern, err := db.EntityPatternToMongoQuery(pbehavior.EntityPattern, "")
	if err != nil || len(matchByPattern) == 0 {
		return excludeIds, err
	}

	if len(excludeIds) > 0 {
		matchByPattern = bson.M{"$and": bson.A{
			bson.M{"_id": bson.M{"$nin": excludeIds}},
			matchByPattern,
		}}
	}

	cursor, err := p.EntityCollection.Find(ctx, matchByPattern)
	if err != nil {
		return excludeIds, err
	}

	idsByPattern, err := p.sendAlarmEvents(ctx, cursor, id, resolver)
	if err != nil {
		return excludeIds, err
	}

	excludeIds = append(excludeIds, idsByPattern...)
	matchByPbehaviorId := bson.M{"pbehavior_info.id": id}
	if len(excludeIds) > 0 {
		matchByPbehaviorId["_id"] = bson.M{"$nin": excludeIds}
	}

	cursor, err = p.EntityCollection.Find(ctx, matchByPbehaviorId)
	if err != nil {
		return excludeIds, err
	}

	idsByPbhInfo, err := p.sendAlarmEvents(ctx, cursor, id, resolver)
	if err != nil {
		return excludeIds, err
	}

	excludeIds = append(excludeIds, idsByPbhInfo...)

	return excludeIds, nil
}

func (p *recomputeMessageProcessor) sendAlarmEvents(
	ctx context.Context,
	cursor mongo.Cursor,
	id string,
	resolver libpbehavior.ComputedEntityTypeResolver,
) ([]string, error) {
	if cursor == nil {
		return nil, nil
	}

	defer cursor.Close(ctx)

	ids := make([]string, 0)
	now := time.Now()
	for cursor.Next(ctx) {
		entity := types.Entity{}
		err := cursor.Decode(&entity)
		if err != nil {
			p.Logger.Err(err).Msg("cannot decode alarm")
			continue
		}

		ids = append(ids, entity.ID)
		resolveResult, err := resolver.Resolve(ctx, entity, now)
		if err != nil {
			return nil, fmt.Errorf("cannot resolve pbehavior for entity: %w", err)
		}

		eventType, output := p.EventManager.GetEventType(resolveResult, entity.PbehaviorInfo)
		if eventType == "" {
			continue
		}

		event, err := p.EventGenerator.Generate(entity)
		if err != nil {
			return nil, fmt.Errorf("cannot generate event: %w", err)
		}

		event.EventType = eventType
		event.Output = output
		event.PbehaviorInfo = libpbehavior.NewPBehaviorInfo(datetime.CpsTime{Time: now}, resolveResult)
		body, err := p.Encoder.Encode(event)
		if err != nil {
			return nil, fmt.Errorf("cannot encode event: %w", err)
		}

		err = p.Publisher.PublishWithContext(
			ctx,
			p.Exchange,
			p.Queue,
			false,
			false,
			amqp.Publishing{
				ContentType:  canopsis.JsonContentType,
				Body:         body,
				DeliveryMode: amqp.Persistent,
			},
		)

		if err != nil {
			return nil, fmt.Errorf("cannot send event: %w", err)
		}

		p.Logger.Debug().Str("pbehavior", id).Str("entity", entity.ID).Msgf("send %s event", event.EventType)
	}

	return ids, nil
}

func (p *recomputeMessageProcessor) logError(err error, errMsg string, msg []byte) {
	if p.FeaturePrintEventOnError {
		p.Logger.Err(err).Str("event", string(msg)).Msg(errMsg)
	} else {
		p.Logger.Err(err).Msg(errMsg)
	}
}
