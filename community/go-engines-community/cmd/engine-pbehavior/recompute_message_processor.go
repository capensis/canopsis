package main

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
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
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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
	InheritedServiceResolver libpbehavior.InheritedServicePbhResolver
	Exchange, Queue          string
	TimezoneConfigProvider   config.TimezoneConfigProvider
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
		resolver, err = p.PbhService.Recompute(ctx, p.TimezoneConfigProvider.Get().Location)
	} else {
		resolver, err = p.PbhService.RecomputeByIds(ctx, ids, p.TimezoneConfigProvider.Get().Location)
	}
	if err != nil {
		return fmt.Errorf("failed to recompute pbehaviors: %w", err)
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

	if len(ids) != 0 {
		var inheritedServicePbhResult libpbehavior.InheritedServicesPbhResolveResult
		var excludeIDs []string

		var servicePbhEventsData libpbehavior.ServiceEventsData

		if event.RecomputeInherited {
			inheritedServicePbhResult, servicePbhEventsData, err = p.InheritedServiceResolver.ComputeAndResolveInheritedServicePbh(ctx, resolver)
			if err != nil {
				return fmt.Errorf("failed to resolve inherited service pbehaviors: %w", err)
			}

			err = p.sendServiceEvents(ctx, servicePbhEventsData, event, ids)
			if err != nil {
				return fmt.Errorf("failed to send service events: %w", err)
			}

			excludeIDs = make([]string, len(inheritedServicePbhResult.IDs))
			copy(excludeIDs, inheritedServicePbhResult.IDs)
		} else {
			inheritedServicePbhResult, err = p.InheritedServiceResolver.GetResolvedInheritedServicePbh(ctx)
			if err != nil {
				return fmt.Errorf("failed to get resolved inherited pbh for services: %w", err)
			}
		}

		for _, id := range ids {
			excludeIDs, err = p.updateAlarms(ctx, id, excludeIDs, inheritedServicePbhResult, resolver, ids, event)
			if err != nil {
				return fmt.Errorf("failed to update alarms: %w", err)
			}
		}
	}

	return nil
}

func (p *recomputeMessageProcessor) sendServiceEvents(
	ctx context.Context,
	servicePbhEventsData libpbehavior.ServiceEventsData,
	rpcEvent rpc.PbehaviorRecomputeEvent,
	updatedPbhIds []string,
) error {
	serviceEvents := servicePbhEventsData.ServiceEvents
	serviceMap := servicePbhEventsData.ServicesMap

	for idx := range serviceEvents {
		p.resolveEventAuthor(rpcEvent, updatedPbhIds, serviceMap[serviceEvents[idx].Component].PbehaviorInfo.ID, &serviceEvents[idx])

		body, err := p.Encoder.Encode(serviceEvents[idx])
		if err != nil {
			return fmt.Errorf("cannot encode event: %w", err)
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
			return fmt.Errorf("cannot send event: %w", err)
		}

		p.Logger.Debug().
			Str("pbehavior", serviceEvents[idx].PbehaviorInfo.ID).
			Str("entity", serviceEvents[idx].Component).
			Msgf("send %s event", serviceEvents[idx].EventType)
	}

	return nil
}

func (p *recomputeMessageProcessor) updateAlarms(
	ctx context.Context,
	pbhId string,
	excludeIds []string,
	inheritedServicePbhResult libpbehavior.InheritedServicesPbhResolveResult,
	resolver libpbehavior.ComputedEntityTypeResolver,
	updatedPbhIds []string,
	event rpc.PbehaviorRecomputeEvent,
) ([]string, error) {
	matchByPbehaviorId := make(bson.M)
	if len(excludeIds) > 0 {
		matchByPbehaviorId["_id"] = bson.M{"$nin": excludeIds}
	}

	if len(inheritedServicePbhResult.IDs) > 0 {
		matchByPbehaviorId["$or"] = []bson.M{
			{"pbehavior_info.id": pbhId},
			{"services": bson.M{"$in": inheritedServicePbhResult.IDs}},
		}
	} else {
		matchByPbehaviorId["pbehavior_info.id"] = pbhId
	}

	cursor, err := p.EntityCollection.Find(ctx, matchByPbehaviorId)
	if err != nil {
		return excludeIds, err
	}

	idsByPbhInfo, err := p.sendAlarmEvents(ctx, cursor, pbhId, inheritedServicePbhResult, resolver, updatedPbhIds, event)
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

	idsByPattern, err := p.sendAlarmEvents(ctx, cursor, pbhId, inheritedServicePbhResult, resolver, updatedPbhIds, event)
	if err != nil {
		return excludeIds, err
	}

	return append(excludeIds, idsByPattern...), nil
}

func (p *recomputeMessageProcessor) sendAlarmEvents(
	ctx context.Context,
	cursor mongo.Cursor,
	pbhId string,
	inheritedServicePbhResult libpbehavior.InheritedServicesPbhResolveResult,
	resolver libpbehavior.ComputedEntityTypeResolver,
	updatedPbhIds []string,
	rpcEvent rpc.PbehaviorRecomputeEvent,
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

		resolveResult = inheritedServicePbhResult.ResolveForEntity(resolveResult, entity.Services)

		event, err := p.EventManager.GetEvent(resolveResult, entity, datetime.CpsTime{Time: now})
		if err != nil {
			p.Logger.Err(err).Str("entity", entity.ID).Msg("cannot generate event")
			continue
		}

		if event.EventType == "" {
			continue
		}

		p.resolveEventAuthor(rpcEvent, updatedPbhIds, entity.PbehaviorInfo.ID, &event)

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

func (p *recomputeMessageProcessor) resolveEventAuthor(
	rpcEvent rpc.PbehaviorRecomputeEvent,
	updatedPbhIds []string,
	prevPbhId string,
	event *types.Event,
) {
	if rpcEvent.Author != "" {
		newPbhId := event.PbehaviorInfo.ID

		if newPbhId != "" && slices.Contains(updatedPbhIds, newPbhId) ||
			prevPbhId != "" && slices.Contains(updatedPbhIds, prevPbhId) {
			event.Author = rpcEvent.Author
			if !event.PbehaviorInfo.IsDefaultActive() {
				event.PbehaviorInfo.Author = rpcEvent.Author
			}

			if rpcEvent.UserID != "" {
				event.UserID = rpcEvent.UserID
			}

			if rpcEvent.Initiator != "" {
				event.Initiator = rpcEvent.Initiator
			}
		}
	}
}
