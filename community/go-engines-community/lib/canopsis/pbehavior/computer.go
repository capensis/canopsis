package pbehavior

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
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

// Computer is used to implement pbehavior intervals recompute on singnal from channel.
type Computer interface {
	// Compute recomputes all pbehaviors on empty signal
	// or recomputes only one pbehavior on signal which conatins pbehavior id.
	Compute(ctx context.Context, ch <-chan ComputeTask)
}

type ComputeTask struct {
	// PbehaviorIds defines for which pbehavior intervals should be recomputed.
	// If empty all intervals are recomputed.
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
	encoder             encoding.Encoder
	publisher           libamqp.Publisher
	queue               string

	tasksMx sync.Mutex
	tasks   []string

	isTest bool
}

// NewCancelableComputer creates new computer.
func NewCancelableComputer(
	service Service,
	dbClient mongo.DbClient,
	publisher libamqp.Publisher,
	eventManager EventManager,
	encoder encoding.Encoder,
	queue string,
	logger zerolog.Logger,
	isTest bool,
) Computer {
	return &cancelableComputer{
		logger:              logger,
		service:             service,
		dbClient:            dbClient,
		alarmCollection:     dbClient.Collection(mongo.AlarmMongoCollection),
		entityCollection:    dbClient.Collection(mongo.EntityMongoCollection),
		pbehaviorCollection: dbClient.Collection(mongo.PbehaviorMongoCollection),
		eventManager:        eventManager,
		encoder:             encoder,
		publisher:           publisher,
		queue:               queue,
		isTest:              isTest,
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

	var err error
	if all {
		err = c.service.Recompute(ctx)
	} else {
		err = c.service.RecomputeByIds(ctx, pbehaviorIds)
	}

	if err != nil {
		c.logger.Err(err).Msgf("API pbehavior recompute failed")
		return
	}

	if c.isTest {
		return
	}

	excludeIds := make([]string, 0)

	for _, pbehaviorID := range pbehaviorIds {
		if pbehaviorID != calculateAll {
			c.logger.Info().Str("pbehavior", pbehaviorID).Msg("API pbehavior recompute: pbehavior recomputed")

			excludeIds, err = c.updateAlarms(ctx, pbehaviorID, excludeIds)
			if err != nil {
				c.logger.Err(err).Msgf("API pbehavior update alarms failed")
				return
			}
		}
	}
}

func (c *cancelableComputer) updateAlarms(
	ctx context.Context,
	pbehaviorID string,
	excludeIds []string,
) ([]string, error) {
	eventGenerator := libevent.NewGenerator(entity.NewAdapter(c.dbClient))

	cursor, err := c.findEntitiesMatchPbhFilter(ctx, pbehaviorID, excludeIds)
	if err != nil {
		return nil, err
	}

	idsByFilter, err := c.sendAlarmEvents(ctx, cursor, pbehaviorID, eventGenerator)
	if err != nil {
		return nil, err
	}

	excludeIds = append(excludeIds, idsByFilter...)
	cursor, err = c.findEntitiesMatchPbhID(ctx, pbehaviorID, excludeIds)
	if err != nil {
		return nil, err
	}

	idsByPbhInfo, err := c.sendAlarmEvents(ctx, cursor, pbehaviorID, eventGenerator)
	if err != nil {
		return nil, err
	}

	excludeIds = append(excludeIds, idsByPbhInfo...)

	return excludeIds, nil
}

func (c *cancelableComputer) sendAlarmEvents(ctx context.Context, cursor mongo.Cursor, pbehaviorID string, eventGenerator libevent.Generator) ([]string, error) {
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
		resolveResult, err := c.service.Resolve(ctx, entity.ID, now)
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
			event, err = eventGenerator.Generate(ctx, entity)
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

		err = c.publisher.Publish(
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

func (c *cancelableComputer) findEntitiesMatchPbhFilter(
	ctx context.Context,
	pbehaviorID string,
	excludeIds []string,
) (mongo.Cursor, error) {
	pbehavior := PBehavior{}
	err := c.pbehaviorCollection.FindOne(ctx, bson.M{"_id": pbehaviorID},
		options.FindOne().SetProjection(bson.M{
			"filter": 1,
		})).Decode(&pbehavior)

	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			return nil, nil
		} else {
			return nil, err
		}
	}

	var filter interface{}
	err = json.Unmarshal([]byte(pbehavior.Filter), &filter)
	if err != nil {
		return nil, err
	}

	if len(excludeIds) > 0 {
		filter = bson.M{"$and": bson.A{
			bson.M{"_id": bson.M{"$nin": excludeIds}},
			filter,
		}}
	}

	return c.entityCollection.Find(ctx, filter)
}

func (c *cancelableComputer) findEntitiesMatchPbhID(
	ctx context.Context,
	pbehaviorID string,
	excludeIds []string,
) (mongo.Cursor, error) {
	filter := bson.M{"pbehavior_info.id": pbehaviorID}

	if len(excludeIds) > 0 {
		filter["_id"] = bson.M{"$nin": excludeIds}
	}

	return c.entityCollection.Find(ctx, filter)
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
