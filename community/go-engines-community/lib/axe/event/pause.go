package event

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters/calculator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewPauseProcessor(
	client mongo.DbClient,
	autoInstructionMatcher AutoInstructionMatcher,
	remediationRpcClient engine.RPCClient,
	metricsSender metrics.Sender,
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator,
	componentCountersCalculator calculator.ComponentCountersCalculator,
	eventsSender entitycounters.EventsSender,
	eventGenerator event.Generator,
	publisher amqp.Publisher,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) Processor {
	return &pauseProcessor{
		client:                 client,
		alarmCollection:        client.Collection(mongo.AlarmMongoCollection),
		entityCollection:       client.Collection(mongo.EntityMongoCollection),
		autoInstructionMatcher: autoInstructionMatcher,
		remediationRpcClient:   remediationRpcClient,
		metricsSender:          metricsSender,
		eventGenerator:         eventGenerator,
		publisher:              publisher,
		encoder:                encoder,
		logger:                 logger,
		countersHelper:         newCountersHelper(entityServiceCountersCalculator, componentCountersCalculator, eventsSender, logger),
	}
}

type pauseProcessor struct {
	client                 mongo.DbClient
	alarmCollection        mongo.DbCollection
	entityCollection       mongo.DbCollection
	autoInstructionMatcher AutoInstructionMatcher
	remediationRpcClient   engine.RPCClient
	metricsSender          metrics.Sender
	eventGenerator         event.Generator
	publisher              amqp.Publisher
	encoder                encoding.Encoder
	logger                 zerolog.Logger
	countersHelper         *countersHelper
}

func (p *pauseProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil || !event.Entity.Enabled {
		return result, nil
	}

	result, countersRes, processDepends, err := p.processEvent(ctx, event)
	if err != nil {
		return result, err
	}

	if result.AlarmChange.Type != types.AlarmChangeTypeNone {
		p.postProcess(context.WithoutCancel(ctx), event, result, countersRes)
	}

	if processDepends {
		err = p.processDepends(ctx, event, result)
	}

	return result, err
}

func (p *pauseProcessor) processEvent(ctx context.Context, event rpc.AxeEvent) (Result, countersResult, bool, error) {
	newPause := types.Pause{
		Timestamp: event.Parameters.Timestamp,
		Message:   event.Parameters.Output,
		Author:    event.Parameters.Author,
		UserID:    event.Parameters.User,
		Initiator: event.Parameters.Initiator,
		Type:      event.Parameters.PauseType,
		Source:    event.Parameters.PauseSource,
		Inherited: event.Parameters.Inherited,
		Path:      event.Parameters.PausePath,
	}
	match := getOpenAlarmMatch(event)
	newStepQuery := stepUpdateQueryWithInPbhInterval(types.AlarmStepPause, event.Parameters.Output, event.Parameters)
	// wasPaused is true if the pause of the alarm is replaced by this one.
	alarmUpdate := func(wasPaused bool) []bson.M {
		set := bson.M{
			"v.pause":          newStepQuery,
			"v.steps":          addStepUpdateQuery(newStepQuery),
			"v.inactive_start": event.Parameters.Timestamp,
			"v.inactive_duration": bson.M{"$sum": bson.A{
				"$v.inactive_duration",
				bson.M{"$cond": bson.M{
					"if": bson.M{"$gt": bson.A{"$v.inactive_start", 0}},
					"then": bson.M{"$subtract": bson.A{
						event.Parameters.Timestamp,
						"$v.inactive_start",
					}},
					"else": 0,
				}},
			}},
			"v.last_update_date": event.Parameters.Timestamp,
		}
		if wasPaused {
			// The replaced pause is not resumed, so its duration is kept.
			set["v.pause_duration"] = bson.M{"$sum": bson.A{
				"$v.pause_duration",
				bson.M{"$subtract": bson.A{event.Parameters.Timestamp, "$v.pause.t"}},
			}}
		}

		return []bson.M{{"$set": set}}
	}

	result := Result{}
	countersRes := countersResult{}
	processDepends := false
	err := p.client.WithTransaction(ctx, func(ctx context.Context) error {
		result = Result{}
		countersRes = countersResult{}
		processDepends = false

		entity := types.Entity{}
		err := p.entityCollection.FindOne(ctx, bson.M{"_id": event.Entity.ID}).Decode(&entity)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		latestPause := types.LatestPause{Pause: newPause}
		if latestPause.Source == "" {
			latestPause.Source = entity.ID
		}

		// An inherited pause replaces an older one and the entity's own pause, the entity's own pause
		// is ignored while it's paused; see types.LatestPause.CanApplyOver.
		if !latestPause.CanApplyOver(entity.LatestPause, entity.ID) {
			return nil
		}

		processDepends = true
		wasPaused := entity.Pause != nil
		entity.Pause = &newPause
		entity.LatestPause = &latestPause
		result.Entity = entity
		_, err = p.entityCollection.UpdateOne(ctx, bson.M{"_id": entity.ID}, bson.M{"$set": bson.M{
			"pause":        newPause,
			"latest_pause": latestPause,
		}})
		if err != nil {
			return fmt.Errorf("cannot update entity: %w", err)
		}

		alarm := types.Alarm{}
		err = p.alarmCollection.
			FindOneAndUpdate(ctx, match, alarmUpdate(wasPaused), options.FindOneAndUpdate().SetReturnDocument(options.After)).
			Decode(&alarm)
		if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
			return err
		}

		result.Forward = true
		result.Alarm = alarm

		// Only the first pause of the alarm is an alarm change.
		if wasPaused {
			return nil
		}

		alarmChange := types.NewAlarmChange()
		alarmChange.Type = types.AlarmChangeTypePause
		result.AlarmChange = alarmChange

		result.IsCountersUpdated, countersRes, err = p.countersHelper.CalculateCounters(
			ctx,
			&result.Alarm,
			&result.Entity,
			result.AlarmChange,
		)

		return err
	})

	if err != nil {
		return result, countersRes, processDepends, err
	}

	if result.Alarm.ID != "" {
		result.IsInstructionMatched = isInstructionMatched(event, result, p.autoInstructionMatcher, p.logger)
	}

	return result, countersRes, processDepends, err
}

func (p *pauseProcessor) postProcess(
	ctx context.Context,
	event rpc.AxeEvent,
	result Result,
	countersRes countersResult,
) {
	entity := *event.Entity
	if result.Entity.ID != "" {
		entity = result.Entity
	}

	p.metricsSender.SendEventMetrics(
		result.Alarm,
		entity,
		result.AlarmChange,
		event.Parameters.Timestamp.Time,
		event.Parameters.Initiator,
		event.Parameters.User,
		event.Parameters.Instruction,
		"",
	)

	p.countersHelper.UpdateStates(ctx, countersRes)

	if result.Alarm.ID != "" {
		err := sendRemediationEvent(ctx, event, result, p.remediationRpcClient, p.encoder)
		if err != nil {
			p.logger.Err(err).Msg("cannot send event to engine-remediation")
		}
	}
}

func (p *pauseProcessor) processDepends(ctx context.Context, event rpc.AxeEvent, result Result) error {
	if !event.Parameters.Inherited {
		return nil
	}

	entity := *event.Entity
	if result.Entity.ID != "" {
		entity = result.Entity
	}

	label := ""
	var filter bson.M
	switch entity.Type {
	case types.EntityTypeComponent:
		filter = bson.M{
			"component": entity.ID,
			"type":      types.EntityTypeResource,
			"enabled":   true,
		}
		label = "Component: "
	case types.EntityTypeService:
		filter = bson.M{
			"services": entity.ID,
			"enabled":  true,
		}
		label = "Service: "
	default:
		return nil
	}

	output := event.Parameters.Output
	pauseSource := event.Parameters.PauseSource
	// The pause of another entity passes through this one, its own pause starts from it,
	// so it's not in the path.
	pausePath := event.Parameters.PausePath
	if pauseSource == "" {
		pauseSource = entity.ID
		output += "\n" + label + entity.Name
	} else {
		pausePath = append(slices.Clone(pausePath), entity.ID)
	}

	cursor, err := p.entityCollection.Find(ctx, filter)
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)

	depEvent := rpc.AxeEvent{
		EventType: event.EventType,
		Parameters: rpc.AxeParameters{
			Output:        output,
			Author:        event.Parameters.Author,
			User:          event.Parameters.User,
			Role:          event.Parameters.Role,
			Initiator:     event.Parameters.Initiator,
			Timestamp:     event.Parameters.Timestamp,
			Inherited:     true,
			Connector:     event.Parameters.Connector,
			ConnectorName: event.Parameters.ConnectorName,
			PauseType:     event.Parameters.PauseType,
			PauseSource:   pauseSource,
			PausePath:     pausePath,
		},
	}

	for cursor.Next(ctx) {
		dep := types.Entity{}
		err = cursor.Decode(&dep)
		if err != nil {
			return err
		}

		depEvent.Entity = &dep
		depRes, countersRes, processDepends, err := p.processEvent(ctx, depEvent)
		if err != nil {
			return err
		}

		if depRes.AlarmChange.Type == types.AlarmChangeTypePause {
			te, err := p.eventGenerator.Generate(dep)
			if err != nil {
				return err
			}

			te.Timestamp = depEvent.Parameters.Timestamp
			te.AlarmChange = &depRes.AlarmChange
			te.EventType = types.EventTypeTrigger
			te.AlarmID = depRes.Alarm.ID
			err = p.sendEvent(ctx, te)
			if err != nil {
				return err
			}

			p.postProcess(context.WithoutCancel(ctx), depEvent, depRes, countersRes)
		}

		if processDepends {
			err = p.processDepends(ctx, depEvent, depRes)
			if err != nil {
				return err
			}
		}
	}

	if err = cursor.Err(); err != nil {
		return err
	}

	return err
}

func (p *pauseProcessor) sendEvent(ctx context.Context, event types.Event) error {
	b, err := p.encoder.Encode(event)
	if err != nil {
		return fmt.Errorf("failed to encode event %w", err)
	}

	return p.publisher.PublishWithContext(
		ctx,
		"",
		canopsis.FIFOQueueName,
		false,
		false,
		amqp091.Publishing{
			ContentType:  canopsis.JsonContentType,
			Body:         b,
			DeliveryMode: amqp091.Persistent,
		},
	)
}
