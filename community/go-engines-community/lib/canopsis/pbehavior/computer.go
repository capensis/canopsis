package pbehavior

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	libevent "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const calculateAll = "all"

// Computer is used to implement pbehavior intervals recompute on signal from channel.
type Computer interface {
	// Compute recomputes all pbehaviors on empty signal
	// or recomputes only one pbehavior on signal which contains pbehavior id.
	Compute(ctx context.Context, ch <-chan ComputeTask)
}

type ComputeTask struct {
	PbehaviorIds []string
}

// cancelableComputer recomputes pbehaviors.
type cancelableComputer struct {
	logger              zerolog.Logger
	service             Service
	dbClient            mongo.DbClient
	alarmCollection     mongo.DbCollection
	entityCollection    mongo.DbCollection
	pbehaviorCollection mongo.DbCollection
	eventManager        EventManager
	decoder             encoding.Decoder
	encoder             encoding.Encoder
	publisher           libamqp.Publisher
	queue               string

	tasksMx sync.Mutex
	tasks   []string
}

// NewCancelableComputer creates new computer.
func NewCancelableComputer(
	service Service,
	dbClient mongo.DbClient,
	publisher libamqp.Publisher,
	eventManager EventManager,
	decoder encoding.Decoder,
	encoder encoding.Encoder,
	queue string,
	logger zerolog.Logger,
) Computer {
	return &cancelableComputer{
		logger:              logger,
		service:             service,
		dbClient:            dbClient,
		alarmCollection:     dbClient.Collection(mongo.AlarmMongoCollection),
		entityCollection:    dbClient.Collection(mongo.EntityMongoCollection),
		pbehaviorCollection: dbClient.Collection(mongo.PbehaviorMongoCollection),
		eventManager:        eventManager,
		decoder:             decoder,
		encoder:             encoder,
		publisher:           publisher,
		queue:               queue,
	}
}

func (c *cancelableComputer) Compute(ctx context.Context, ch <-chan ComputeTask) {
	startComputeCh := make(chan bool, 1)

	go func() {
		defer close(startComputeCh)

		for {
			select {
			case <-ctx.Done():
				return
			case task, ok := <-ch:
				if !ok {
					return
				}

				if len(task.PbehaviorIds) == 0 {
					c.addTasks([]string{calculateAll})
				} else {
					c.addTasks(task.PbehaviorIds)
				}

				select {
				case startComputeCh <- true:
				default:
				}
			}
		}
	}()

	for range startComputeCh {
		for {
			tasks := c.getTasks()
			if len(tasks) == 0 {
				break
			}

			c.computePbehaviors(ctx, tasks)
		}
	}
}

func (c *cancelableComputer) addTasks(tasks []string) {
	c.tasksMx.Lock()
	defer c.tasksMx.Unlock()

	c.tasks = append(c.tasks, tasks...)
}

func (c *cancelableComputer) getTasks() []string {
	c.tasksMx.Lock()
	defer c.tasksMx.Unlock()

	tasks := c.tasks
	c.tasks = make([]string, 0)

	return tasks
}

// computePbehaviors obtains lock and calls computer.
func (c *cancelableComputer) computePbehaviors(
	ctx context.Context,
	pbehaviorIds []string,
) {
	all := false
	for _, id := range pbehaviorIds {
		if id == calculateAll {
			all = true
			break
		}
	}

	var resolver ComputedEntityTypeResolver
	var err error
	if all {
		resolver, err = c.service.Recompute(ctx)
	} else {
		resolver, err = c.service.RecomputeByIds(ctx, pbehaviorIds)
	}

	if err != nil {
		c.logger.Err(err).Msgf("API pbehavior recompute failed")
		return
	}

	excludeIds := make([]string, 0)

	for _, pbehaviorID := range pbehaviorIds {
		if pbehaviorID != calculateAll {
			c.logger.Info().Str("pbehavior", pbehaviorID).Msg("pbehavior recomputed")

			excludeIds, err = c.updateAlarms(ctx, pbehaviorID, excludeIds, resolver)
			if err != nil {
				c.logger.Err(err).Str("pbehavior", pbehaviorID).Msg("API pbehavior update alarms failed")
				return
			}
		}
	}
}

func (c *cancelableComputer) updateAlarms(
	ctx context.Context,
	pbehaviorID string,
	excludeIds []string,
	resolver ComputedEntityTypeResolver,
) ([]string, error) {
	eventGenerator := libevent.NewGenerator("api", "api")
	pbehavior := PBehavior{}
	err := c.pbehaviorCollection.FindOne(ctx, bson.M{"_id": pbehaviorID},
		options.FindOne().SetProjection(bson.M{
			"entity_pattern":  1,
			"old_mongo_query": 1,
		})).Decode(&pbehavior)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	var query interface{}
	if len(pbehavior.EntityPattern) > 0 {
		query, err = pbehavior.EntityPattern.ToMongoQuery("")
		if err != nil {
			return nil, err
		}
	} else {
		var oldMongoQuery map[string]interface{}
		err = c.decoder.Decode([]byte(pbehavior.OldMongoQuery), &oldMongoQuery)
		if err != nil {
			return nil, err
		}

		query = oldMongoQuery
	}

	matchByPattern := query
	if len(excludeIds) > 0 {
		matchByPattern = bson.M{"$and": bson.A{
			bson.M{"_id": bson.M{"$nin": excludeIds}},
			matchByPattern,
		}}
	}

	cursor, err := c.entityCollection.Find(ctx, matchByPattern)
	if err != nil {
		return nil, err
	}

	idsByPattern, err := c.sendAlarmEvents(ctx, cursor, pbehaviorID, eventGenerator, resolver)
	if err != nil {
		return nil, err
	}

	excludeIds = append(excludeIds, idsByPattern...)

	matchByPbehaviorId := bson.M{"pbehavior_info.id": pbehaviorID}
	if len(excludeIds) > 0 {
		matchByPbehaviorId["_id"] = bson.M{"$nin": excludeIds}
	}

	cursor, err = c.entityCollection.Find(ctx, matchByPbehaviorId)
	if err != nil {
		return nil, err
	}

	idsByPbhInfo, err := c.sendAlarmEvents(ctx, cursor, pbehaviorID, eventGenerator, resolver)
	if err != nil {
		return nil, err
	}

	excludeIds = append(excludeIds, idsByPbhInfo...)

	return excludeIds, nil
}

func (c *cancelableComputer) sendAlarmEvents(
	ctx context.Context,
	cursor mongo.Cursor,
	pbehaviorID string,
	eventGenerator libevent.Generator,
	resolver ComputedEntityTypeResolver,
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
			c.logger.Err(err).Msg("cannot decode alarm")
			continue
		}

		ids = append(ids, entity.ID)
		resolveResult, err := resolver.Resolve(ctx, entity, now)
		if err != nil {
			return nil, fmt.Errorf("cannot resolve pbehavior for entity: %w", err)
		}

		eventType, output := c.eventManager.GetEventType(resolveResult, entity.PbehaviorInfo)
		if eventType == "" {
			continue
		}

		alarm, err := c.findLastAlarm(ctx, entity.ID)
		if err != nil {
			return nil, fmt.Errorf("cannot find alarm: %w", err)
		}

		event := types.Event{
			Initiator: types.InitiatorSystem,
		}
		if alarm == nil {
			event, err = eventGenerator.Generate(entity)
			if err != nil {
				return nil, fmt.Errorf("cannot generate event: %w", err)
			}
		} else {
			event.Connector = alarm.Value.Connector
			event.ConnectorName = alarm.Value.ConnectorName
			event.Component = alarm.Value.Component
			event.Resource = alarm.Value.Resource
			event.SourceType = event.DetectSourceType()
		}

		event.EventType = eventType
		event.Output = output
		event.PbehaviorInfo = NewPBehaviorInfo(types.CpsTime{Time: now}, resolveResult)
		body, err := c.encoder.Encode(event)
		if err != nil {
			return nil, fmt.Errorf("cannot encode event: %w", err)
		}

		err = c.publisher.PublishWithContext(
			ctx,
			"",
			c.queue,
			false,
			false,
			amqp.Publishing{
				ContentType:  "application/json",
				Body:         body,
				DeliveryMode: amqp.Persistent,
			},
		)

		if err != nil {
			return nil, fmt.Errorf("cannot send event: %w", err)
		}

		c.logger.Info().Str("pbehavior", pbehaviorID).Str("entity", entity.ID).Msgf("send %s event", event.EventType)
	}

	return ids, nil
}

func (c *cancelableComputer) findLastAlarm(
	ctx context.Context,
	entityID string,
) (*types.Alarm, error) {
	alarm := types.Alarm{}
	err := c.alarmCollection.
		FindOne(ctx, bson.M{"d": entityID}, options.FindOne().SetSort(bson.M{"t": -1})).
		Decode(&alarm)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	return &alarm, nil
}
