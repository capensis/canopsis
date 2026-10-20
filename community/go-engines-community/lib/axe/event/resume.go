package event

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"slices"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters/calculator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewResumeProcessor(
	client mongo.DbClient,
	autoInstructionMatcher AutoInstructionMatcher,
	remediationRpcClient engine.RPCClient,
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator,
	componentCountersCalculator calculator.ComponentCountersCalculator,
	eventsSender entitycounters.EventsSender,
	eventGenerator event.Generator,
	publisher amqp.Publisher,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) Processor {
	return &resumeProcessor{
		client:                 client,
		alarmCollection:        client.Collection(mongo.AlarmMongoCollection),
		entityCollection:       client.Collection(mongo.EntityMongoCollection),
		autoInstructionMatcher: autoInstructionMatcher,
		remediationRpcClient:   remediationRpcClient,
		eventGenerator:         eventGenerator,
		publisher:              publisher,
		encoder:                encoder,
		logger:                 logger,
		countersHelper:         newCountersHelper(entityServiceCountersCalculator, componentCountersCalculator, eventsSender, logger),
	}
}

type resumeProcessor struct {
	client                 mongo.DbClient
	alarmCollection        mongo.DbCollection
	entityCollection       mongo.DbCollection
	autoInstructionMatcher AutoInstructionMatcher
	remediationRpcClient   engine.RPCClient
	eventGenerator         event.Generator
	publisher              amqp.Publisher
	encoder                encoding.Encoder
	logger                 zerolog.Logger
	countersHelper         *countersHelper
}

func (p *resumeProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil || !event.Entity.Enabled {
		return result, nil
	}

	result, countersRes, ended, err := p.processEvent(ctx, event, nil)
	if err != nil {
		return result, err
	}

	if result.AlarmChange.Type != types.AlarmChangeTypeNone {
		p.postProcess(context.WithoutCancel(ctx), event, result, countersRes)
	}

	// An inherited pause is always resumed on dependents, whatever the resume event is,
	// otherwise they would stay paused forever.
	if ended != nil && ended.Inherited {
		p.processDepends(ctx, event, result, ended.Timestamp)
	}

	return result, nil
}

// processEvent ends the pause which started at pauseStart.
// pauseStart is nil for a direct resume, then the entity's own current pause is ended.
// It returns the ended pause if the resume is applied, nil if it is ignored.
func (p *resumeProcessor) processEvent(ctx context.Context, event rpc.AxeEvent, pauseStart *datetime.CpsTime) (_ Result, _ countersResult, _ *types.LatestPause, resErr error) {
	// A panic is returned as an error, so a panic on a dependency doesn't stop the processing of the other ones.
	defer func() {
		if r := recover(); r != nil {
			resErr = fmt.Errorf("panic on resume of %s: %v\n%s", event.Entity.ID, r, debug.Stack())
		}
	}()

	match := getOpenAlarmMatch(event)
	now := datetime.NewCpsTime()
	newStepQuery := stepUpdateQueryWithInPbhIntervalByStep(types.NewAlarmStep(types.AlarmStepResume, now,
		event.Parameters.Author, event.Parameters.Output, event.Parameters.User, event.Parameters.Role,
		event.Parameters.Initiator, false))
	update := []bson.M{
		{"$set": bson.M{
			"v.steps": addStepUpdateQuery(newStepQuery),
			"v.pause_duration": bson.M{"$sum": bson.A{
				"$v.pause_duration",
				bson.M{"$subtract": bson.A{
					now,
					"$v.pause.t",
				}},
			}},
			"v.inactive_duration": bson.M{"$sum": bson.A{
				"$v.inactive_duration",
				bson.M{"$subtract": bson.A{
					now,
					"$v.inactive_start",
				}},
			}},
			"v.inactive_start":   updateInactiveStart(now, true, true, true, true, false),
			"v.last_update_date": now,
		}},
		{"$unset": "v.pause"},
	}

	result := Result{}
	countersRes := countersResult{}
	var endedPause *types.LatestPause
	err := p.client.WithTransaction(ctx, func(ctx context.Context) error {
		result = Result{}
		countersRes = countersResult{}
		endedPause = nil

		entity := types.Entity{}
		err := p.entityCollection.FindOne(ctx, bson.M{"_id": event.Entity.ID}).Decode(&entity)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		source := event.Parameters.PauseSource
		sourceType := event.Parameters.PauseSourceType
		if source == "" {
			source = entity.ID
			sourceType = entity.Type
		}

		start := pauseStart
		if start == nil {
			// Direct resume: only the entity's active pause from this source can be ended.
			if entity.LatestPause == nil || entity.LatestPause.End != nil || entity.LatestPause.Source != source {
				return nil
			}

			start = &entity.LatestPause.Timestamp
		}

		end := event.Parameters.Timestamp
		if end.Before(*start) {
			p.logger.Warn().
				Str("entity", entity.ID).
				Str("connector", event.Parameters.Connector+"/"+event.Parameters.ConnectorName).
				Str("pause", start.String()).
				Str("resume", end.String()).
				Msg("the resume is older than the pause it ends, it is ignored")

			return nil
		}

		ended := types.LatestPause{
			// A pause of another entity is always inherited.
			Timestamp:  *start,
			Source:     source,
			SourceType: sourceType,
			Inherited:  source != entity.ID,
			Path:       event.Parameters.PausePath,
			End:        new(end),
		}
		switch {
		case entity.LatestPause != nil && entity.LatestPause.Source == source && entity.LatestPause.Timestamp.Equal(*start):
			if entity.LatestPause.End != nil {
				return nil // duplicate
			}

			ended.Pause = entity.LatestPause.Pause
		case ended.CanApplyOver(entity.LatestPause, entity.ID):
			// Resume arrived before its pause, keep the ended pause so the late pause and older ones are ignored.
		default:
			return nil // resume of a replaced pause
		}

		endedPause = &ended
		wasPaused := entity.Pause != nil
		entity.Pause = nil
		entity.LatestPause = &ended
		result.Entity = entity
		_, err = p.entityCollection.UpdateOne(ctx, bson.M{"_id": entity.ID}, bson.M{
			"$set":   bson.M{"latest_pause": ended},
			"$unset": bson.M{"pause": ""},
		})
		if err != nil {
			return fmt.Errorf("cannot update entity: %w", err)
		}

		if !wasPaused {
			return nil
		}

		alarm := types.Alarm{}
		err = p.alarmCollection.
			FindOneAndUpdate(ctx, match, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).
			Decode(&alarm)
		if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
			return err
		}

		alarmChange := types.NewAlarmChange()
		alarmChange.Type = types.AlarmChangeTypeResume
		alarmChange.PreviousPbehaviorTypeID = entity.PbehaviorInfo.TypeID
		alarmChange.PreviousPbehaviorCannonicalType = entity.PbehaviorInfo.CanonicalType
		result.Forward = true
		result.Alarm = alarm
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
		return result, countersRes, endedPause, err
	}

	if result.Alarm.ID != "" {
		result.IsInstructionMatched = isInstructionMatched(event, result, p.autoInstructionMatcher, p.logger)
	}

	return result, countersRes, endedPause, err
}

func (p *resumeProcessor) postProcess(
	ctx context.Context,
	event rpc.AxeEvent,
	result Result,
	countersRes countersResult,
) {
	// A panic is logged, so a panic on a dependency doesn't stop the processing of the other ones.
	defer func() {
		if r := recover(); r != nil {
			p.logger.Error().Str("entity", event.Entity.ID).Interface("panic", r).Bytes("stack", debug.Stack()).
				Msg("panic on post process of resume")
		}
	}()

	p.countersHelper.UpdateStates(ctx, countersRes)

	if result.Alarm.ID != "" {
		err := sendRemediationEvent(ctx, event, result, p.remediationRpcClient, p.encoder)
		if err != nil {
			p.logger.Err(err).Msg("cannot send event to engine-remediation")
		}
	}
}

func (p *resumeProcessor) processDepends(ctx context.Context, event rpc.AxeEvent, result Result, pauseStart datetime.CpsTime) {
	entity := *event.Entity
	if result.Entity.ID != "" {
		entity = result.Entity
	}

	var filter bson.M
	// Dependents are not filtered by pause.source: they may not have received the pause yet
	// and must keep the ended pause to ignore it when it arrives.
	switch entity.Type {
	case types.EntityTypeComponent:
		filter = bson.M{
			"component": entity.ID,
			"type":      types.EntityTypeResource,
			"enabled":   true,
		}
	case types.EntityTypeService:
		filter = bson.M{
			"services": entity.ID,
			"enabled":  true,
		}
	default:
		return
	}

	output := event.Parameters.Output
	pauseSource := event.Parameters.PauseSource
	pauseSourceType := event.Parameters.PauseSourceType
	// The resumed pause of another entity passes through this one, its own pause starts from it,
	// so it's not in the path.
	pausePath := event.Parameters.PausePath
	if pauseSource == "" {
		pauseSource = entity.ID
		pauseSourceType = entity.Type
		if entity.LatestPause != nil {
			output = entity.LatestPause.GetInheritedMessage(entity.Type, entity.Name)
		}
	} else {
		if len(pausePath) >= types.MaxPausePathLength {
			p.logger.Error().Str("entity", entity.ID).Int("path_length", len(pausePath)).
				Msg("the pause path is too long, the resume is not propagated to the dependencies")

			return
		}

		pausePath = append(slices.Clone(pausePath), entity.ID)
	}

	cursor, err := p.entityCollection.Find(ctx, filter)
	if err != nil {
		p.logger.Err(err).Str("entity", entity.ID).Msg("cannot find dependencies")

		return
	}

	defer cursor.Close(ctx)

	depEvent := rpc.AxeEvent{
		EventType: event.EventType,
		Parameters: rpc.AxeParameters{
			Output:          output,
			Author:          event.Parameters.Author,
			User:            event.Parameters.User,
			Role:            event.Parameters.Role,
			Initiator:       event.Parameters.Initiator,
			Timestamp:       event.Parameters.Timestamp,
			Inherited:       true,
			Connector:       event.Parameters.Connector,
			ConnectorName:   event.Parameters.ConnectorName,
			PauseSource:     pauseSource,
			PauseSourceType: pauseSourceType,
			PausePath:       pausePath,
		},
	}

	for cursor.Next(ctx) {
		dep := types.Entity{}
		err = cursor.Decode(&dep)
		if err != nil {
			p.logger.Err(err).Msg("cannot decode dependency")

			continue
		}

		depEvent.Entity = &dep
		depRes, countersRes, depEnded, err := p.processEvent(ctx, depEvent, &pauseStart)
		if err != nil {
			p.logger.Err(err).Str("entity", entity.ID).Str("dependency", dep.ID).
				Msg("cannot process resume for dependency")

			continue
		}

		if depRes.AlarmChange.Type != types.AlarmChangeTypeNone {
			te, err := p.eventGenerator.Generate(dep)
			if err != nil {
				p.logger.Err(err).Str("entity", entity.ID).Str("dependency", dep.ID).
					Msg("cannot generate event for dependency")

				continue
			}

			te.Timestamp = depEvent.Parameters.Timestamp
			te.AlarmChange = &depRes.AlarmChange
			te.EventType = types.EventTypeTrigger
			te.AlarmID = depRes.Alarm.ID
			err = p.sendEvent(ctx, te)
			if err != nil {
				p.logger.Err(err).Str("entity", entity.ID).Str("dependency", dep.ID).
					Msg("cannot send trigger event for dependency")

				continue
			}

			p.postProcess(context.WithoutCancel(ctx), depEvent, depRes, countersRes)
		}

		if depEnded != nil {
			p.processDepends(ctx, depEvent, depRes, pauseStart)
		}
	}

	if err = cursor.Err(); err != nil {
		p.logger.Err(err).Str("entity", entity.ID).Msg("cannot fetch dependencies")
	}
}

func (p *resumeProcessor) sendEvent(ctx context.Context, event types.Event) error {
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
