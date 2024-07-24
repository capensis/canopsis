package main

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
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

	err = p.computePbehaviors(ctx, event)
	if err != nil {
		if engine.IsConnectionError(err) {
			return nil, err
		}

		p.logError(err, "cannot process event", msg)
		return nil, nil
	}

	return nil, nil
}

func (p *recomputeMessageProcessor) computePbehaviors(ctx context.Context, event rpc.PbehaviorRecomputeEvent) error {
	ids := event.Ids
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
		excludeIds, err = p.updateAlarms(ctx, id, excludeIds, resolver, ids, event.Author, event.UserID, event.Initiator)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *recomputeMessageProcessor) updateAlarms(
	ctx context.Context,
	pbhId string,
	excludeIds []string,
	resolver libpbehavior.ComputedEntityTypeResolver,
	updatedPbhIds []string,
	author, userID, initiator string,
) ([]string, error) {
	matchByPbehaviorId := bson.M{"pbehavior_info.id": pbhId}
	if len(excludeIds) > 0 {
		matchByPbehaviorId["_id"] = bson.M{"$nin": excludeIds}
	}

	cursor, err := p.EntityCollection.Find(ctx, matchByPbehaviorId)
	if err != nil {
		return excludeIds, err
	}

	idsByPbhInfo, err := p.sendAlarmEvents(ctx, cursor, pbhId, resolver, updatedPbhIds, author, userID, initiator)
	if err != nil {
		return excludeIds, err
	}

	excludeIds = append(excludeIds, idsByPbhInfo...)
	pbehavior := libpbehavior.PBehavior{}
	err = p.PbehaviorCollection.FindOne(ctx, bson.M{"_id": pbhId},
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

	cursor, err = p.EntityCollection.Find(ctx, matchByPattern)
	if err != nil {
		return excludeIds, err
	}

	idsByPattern, err := p.sendAlarmEvents(ctx, cursor, pbhId, resolver, updatedPbhIds, author, userID, initiator)
	if err != nil {
		return excludeIds, err
	}

	excludeIds = append(excludeIds, idsByPattern...)

	return excludeIds, nil
}

func (p *recomputeMessageProcessor) sendAlarmEvents(
	ctx context.Context,
	cursor mongo.Cursor,
	pbhId string,
	resolver libpbehavior.ComputedEntityTypeResolver,
	updatedPbhIds []string,
	author, userID, initiator string,
) ([]string, error) {
	if cursor == nil {
		return nil, nil
	}

	defer cursor.Close(ctx)

	entityIds := make([]string, 0)
	now := time.Now()
	for cursor.Next(ctx) {
		entity := types.Entity{}
		err := cursor.Decode(&entity)
		if err != nil {
			p.Logger.Err(err).Msg("cannot decode alarm")
			continue
		}

		entityIds = append(entityIds, entity.ID)
		resolveResult, err := resolver.Resolve(ctx, entity, now)
		if err != nil {
			return nil, fmt.Errorf("cannot resolve pbehavior for entity: %w", err)
		}

		event, err := p.EventManager.GetEvent(resolveResult, entity, datetime.CpsTime{Time: now})
		if err != nil {
			p.Logger.Err(err).Str("entity", entity.ID).Msg("cannot generate event")
			continue
		}

		if event.EventType == "" {
			continue
		}

		if author != "" {
			newPbhId := event.PbehaviorInfo.ID
			prevPbhId := entity.PbehaviorInfo.ID
			if newPbhId != "" && slices.Contains(updatedPbhIds, newPbhId) ||
				prevPbhId != "" && slices.Contains(updatedPbhIds, prevPbhId) {
				event.Author = author
				if !event.PbehaviorInfo.IsDefaultActive() {
					event.PbehaviorInfo.Author = author
				}

				if userID != "" {
					event.UserID = userID
				}

				if initiator != "" {
					event.Initiator = initiator
				}
			}
		}

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

		p.Logger.Debug().Str("pbehavior", pbhId).Str("entity", entity.ID).Msgf("send %s event", event.EventType)
	}

	return entityIds, nil
}

func (p *recomputeMessageProcessor) logError(err error, errMsg string, msg []byte) {
	if p.FeaturePrintEventOnError {
		p.Logger.Err(err).Str("event", string(msg)).Msg(errMsg)
	} else {
		p.Logger.Err(err).Msg(errMsg)
	}
}
