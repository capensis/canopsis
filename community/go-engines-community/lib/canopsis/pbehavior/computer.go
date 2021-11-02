package pbehavior

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	libevent "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

const (
	OperationCreate = iota
	OperationUpdate
	OperationDelete
)

const (
	CalculateAll  = ""
	CalculateNone = "none"
)

// Computer is used to implement pbehavior intervals recompute on singnal from channel.
type Computer interface {
	// Compute recomputes all pbehaviors on empty signal
	// or recomputes only one pbehavior on signal which conatins pbehavior id.
	Compute(ctx context.Context, ch <-chan ComputeTask)
}

type ComputeTask struct {
	// PbehaviorID defines for which pbehavior intervals should be recomputed.
	// If empty all intervals are recomputed.
	PbehaviorID string
	// UpdateAlarms prescribes to update alarms by pbehavior. PbehaviorID shouldn't be empty.
	OperationType int
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
	}
}

func (c *cancelableComputer) Compute(parentCtx context.Context, ch <-chan ComputeTask) {
	var ctx context.Context
	var cancel context.CancelFunc
	c.logger.Debug().Msg("compute started")
	defer func() {
		if cancel != nil {
			cancel()
		}
		c.logger.Debug().Msg("compute ended")
	}()

	ongoingID := CalculateNone
	ongoingMx := sync.Mutex{}

	wg := sync.WaitGroup{}

	for {
		select {
		case <-parentCtx.Done():
			return
		case task, ok := <-ch:
			if !ok {
				return
			}

			pbehaviorID := task.PbehaviorID

			ongoingMx.Lock()
			currentOngoing := ongoingID
			ongoingMx.Unlock()

			// Cancel ongoing calculation if
			// - ongoing calculated pbehavior is the same as in the signal
			// - ongoing calculation calculates all pbehavior
			// - new calculation calculates all pbehavior
			if currentOngoing != CalculateNone && (pbehaviorID == CalculateAll || currentOngoing == CalculateAll || currentOngoing == pbehaviorID) {
				c.logger.Info().Str("pbehaviorID", pbehaviorID).Str("ongoingID", currentOngoing).Msg("Test cancel log")
				c.logger.Debug().Msg("API pbehavior recompute: STOP SIGNAL!!!")

				cancel()
				wg.Wait()
				// If canceled calculation calculates all pbehavior keep it
				if currentOngoing != CalculateAll {
					ongoingID = pbehaviorID
				} else {
					ongoingID = CalculateAll
				}
			} else {
				// Wait ongoing calculation to avoid conflicts
				wg.Wait()
				ongoingID = pbehaviorID
			}

			ctx, cancel = context.WithCancel(parentCtx)
			wg.Add(1)

			go func(ctx context.Context, pbhId string, opType int) {
				defer func() {
					ongoingMx.Lock()
					ongoingID = CalculateNone
					ongoingMx.Unlock()

					wg.Done()
				}()

				c.computePbehavior(ctx, pbhId, opType)
			}(ctx, ongoingID, task.OperationType)
		}
	}
}

// computePbehavior obtains lock and calls computer.
func (c *cancelableComputer) computePbehavior(
	ctx context.Context,
	pbehaviorID string,
	operationType int,
) {
	var err error
	if pbehaviorID == "" {
		err = c.service.Recompute(ctx)
	} else {
		err = c.service.RecomputeByID(ctx, pbehaviorID)
	}

	if err != nil {
		c.logger.Err(err).Msgf("API pbehavior recompute failed")
		return
	}

	c.logger.Info().Str("pbehavior", pbehaviorID).Msg("API pbehavior recompute: pbehavior recomputed")

	if pbehaviorID != "" {
		err := c.updateAlarms(ctx, pbehaviorID, operationType)
		if err != nil {
			c.logger.Err(err).Msgf("API pbehavior update alarms failed")
			return
		}
	}
}

func (c *cancelableComputer) updateAlarms(
	ctx context.Context,
	pbehaviorID string,
	operationType int,
) error {
	eventGenerator := libevent.NewGenerator(entity.NewAdapter(c.dbClient))

	switch operationType {
	case OperationCreate:
		cursor, err := c.findEntitiesMatchPbhFilter(ctx, pbehaviorID)
		if err != nil {
			return fmt.Errorf("cannot find alarms: %w", err)
		}

		err = c.sendAlarmEvents(ctx, cursor, pbehaviorID, eventGenerator)
		if err != nil {
			return err
		}
	case OperationUpdate:
		cursor, err := c.findEntitiesMatchPbhFilter(ctx, pbehaviorID)
		if err != nil {
			return err
		}

		err = c.sendAlarmEvents(ctx, cursor, pbehaviorID, eventGenerator)
		if err != nil {
			return err
		}

		cursor, err = c.findEntitiesMatchPbhID(ctx, pbehaviorID)
		if err != nil {
			return err
		}

		err = c.sendAlarmEvents(ctx, cursor, pbehaviorID, eventGenerator)
		if err != nil {
			return err
		}
	case OperationDelete:
		cursor, err := c.findEntitiesMatchPbhID(ctx, pbehaviorID)
		if err != nil {
			return fmt.Errorf("cannot find alarms: %w", err)
		}

		err = c.sendAlarmEvents(ctx, cursor, pbehaviorID, eventGenerator)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *cancelableComputer) sendAlarmEvents(ctx context.Context, cursor mongo.Cursor, pbehaviorID string, eventGenerator libevent.Generator) error {
	now := time.Now()
	for cursor.Next(ctx) {
		entity := types.Entity{}
		err := cursor.Decode(&entity)
		if err != nil {
			c.logger.Err(err).Msg("cannot decode alarm")
			continue
		}

		resolveResult, err := c.service.Resolve(ctx, entity.ID, now)
		if err != nil {
			return fmt.Errorf("cannot resolve pbehavior for entity: %w", err)
		}

		eventType, output := c.eventManager.GetEventType(resolveResult, entity.PbehaviorInfo)
		if eventType == "" {
			continue
		}

		alarm, err := c.findLastAlarm(ctx, entity.ID)
		if err != nil {
			return fmt.Errorf("cannot find alarm: %w", err)
		}

		event := types.Event{}
		if alarm == nil {
			event, err = eventGenerator.Generate(ctx, entity)
			if err != nil {
				return fmt.Errorf("cannot generate event: %w", err)
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
			return fmt.Errorf("cannot encode event: %w", err)
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
			return fmt.Errorf("cannot send event: %w", err)
		}

		c.logger.Info().Str("pbehavior", pbehaviorID).Str("entity", entity.ID).Msgf("send %s event", event.EventType)
	}

	return nil
}

func (c *cancelableComputer) findEntitiesMatchPbhFilter(
	ctx context.Context,
	pbehaviorID string,
) (mongo.Cursor, error) {
	pbehavior := PBehavior{}
	err := c.pbehaviorCollection.FindOne(ctx, bson.M{"_id": pbehaviorID},
		options.FindOne().SetProjection(bson.M{
			"filter": 1,
		})).Decode(&pbehavior)

	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			return nil, fmt.Errorf("cannot find pbehavior id=%s", pbehaviorID)
		} else {
			return nil, err
		}
	}

	var filter interface{}
	err = json.Unmarshal([]byte(pbehavior.Filter), &filter)
	if err != nil {
		return nil, err
	}

	return c.entityCollection.Find(ctx, filter)
}

func (c *cancelableComputer) findEntitiesMatchPbhID(
	ctx context.Context,
	pbehaviorID string,
) (mongo.Cursor, error) {
	return c.entityCollection.Find(ctx, bson.M{"pbehavior_info.id": pbehaviorID})
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
