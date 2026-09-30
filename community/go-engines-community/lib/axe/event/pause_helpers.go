package event

import (
	"context"
	"errors"
	"fmt"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	libevent "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// pauseChange is the pause or the resume applied to an entity and its alarm outside of a pause
// or resume event, e.g. when the entity has been added to or removed from services.
type pauseChange struct {
	// alarmChange is types.AlarmChangeTypePause when the entity has been paused,
	// types.AlarmChangeTypeResume when it has been resumed, and it's empty when the entity stays
	// paused because its pause has only been replaced by another one, or when its pause hasn't changed.
	alarmChange types.AlarmChangeType
	// pause is the applied pause, nil on a resume or if the pause hasn't changed. It's propagated
	// to the dependencies of the entity even if it only replaces its previous pause.
	pause *types.LatestPause
}

// resolveInheritedPause returns the pause which the entity inherits from its parents:
// the newest inheritable pause of its services and of its component, nil if none of them is paused.
func resolveInheritedPause(ctx context.Context, entityCollection mongo.DbCollection, entity types.Entity) (*types.LatestPause, error) {
	parentIDs := make([]string, len(entity.Services))
	copy(parentIDs, entity.Services)
	if entity.Type == types.EntityTypeResource {
		parentIDs = append(parentIDs, entity.Component)
	}

	if len(parentIDs) == 0 {
		return nil, nil
	}

	parent := types.Entity{}
	err := entityCollection.
		FindOne(ctx,
			// A pause which is not inherited is applied only to the parent itself.
			bson.M{"_id": bson.M{"$in": parentIDs}, "pause.inherited": true},
			// pause.source breaks ties so all entities of the parents inherit the same pause.
			options.FindOne().
				SetSort(bson.D{{Key: "pause.timestamp", Value: -1}, {Key: "pause.source", Value: -1}}).
				SetProjection(bson.M{"pause": 1, "type": 1, "name": 1})).
		Decode(&parent)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	return parent.InheritedPause(), nil
}

// updatePauseOnServicesChange updates the pause of the entity and of its alarm when the entity has been
// added to or removed from services: the entity inherits the newest inheritable pause of its parents,
// and its inherited pause is ended when none of its parents has it anymore. Its own pause is kept.
// The entity and the alarm are updated in the caller transaction and in place, the returned alarm change
// is a pause, a resume or empty when the pause hasn't changed.
func updatePauseOnServicesChange(
	ctx context.Context,
	entityCollection, alarmCollection mongo.DbCollection,
	entity *types.Entity,
	alarm *types.Alarm,
	now datetime.CpsTime,
) (pauseChange, error) {
	if len(entity.ServicesToAdd) == 0 && len(entity.ServicesToRemove) == 0 {
		return pauseChange{}, nil
	}

	inherited, err := resolveInheritedPause(ctx, entityCollection, *entity)
	if err != nil {
		return pauseChange{}, err
	}

	if inherited == nil {
		// The entity is not paused, or it's paused by itself, so the change doesn't affect its pause.
		if entity.Pause == nil || entity.Pause.Source == "" {
			return pauseChange{}, nil
		}

		// None of the parents has the entity pause anymore.
		err = resumeInheritedPause(ctx, entityCollection, alarmCollection, entity, alarm, now)
		if err != nil {
			return pauseChange{}, err
		}

		return pauseChange{alarmChange: types.AlarmChangeTypeResume}, nil
	}

	// The entity already has the pause of its parents.
	if entity.Pause != nil && entity.Pause.Source == inherited.Source && entity.Pause.Timestamp.Equal(inherited.Timestamp) {
		return pauseChange{}, nil
	}

	// The entity has already received a pause which the inherited one doesn't replace.
	if entity.LatestPause != nil && !inherited.CanApplyOver(entity.LatestPause, entity.ID) {
		return pauseChange{}, nil
	}

	wasPaused := entity.Pause != nil

	err = applyInheritedPause(ctx, entityCollection, alarmCollection, entity, alarm, *inherited, now)
	if err != nil {
		return pauseChange{}, err
	}

	change := pauseChange{pause: inherited}
	if !wasPaused {
		change.alarmChange = types.AlarmChangeTypePause
	}

	return change, nil
}

func applyInheritedPause(
	ctx context.Context,
	entityCollection, alarmCollection mongo.DbCollection,
	entity *types.Entity,
	alarm *types.Alarm,
	pause types.LatestPause,
	now datetime.CpsTime,
) error {
	wasPaused := entity.Pause != nil

	newEntity, err := updateEntityByID(ctx, entity.ID, bson.M{"$set": bson.M{
		"pause":        pause.Pause,
		"latest_pause": pause,
	}}, entityCollection)
	if err != nil {
		return err
	}

	*entity = newEntity

	if alarm.ID == "" {
		return nil
	}

	newStepQuery := stepUpdateQueryWithInPbhIntervalByStep(types.NewAlarmStep(types.AlarmStepPause, pause.Timestamp,
		pause.Author, pause.Message, pause.UserID, "", pause.Initiator, false))
	set := bson.M{
		"v.pause":          newStepQuery,
		"v.steps":          addStepUpdateQuery(newStepQuery),
		"v.inactive_start": pause.Timestamp,
		"v.inactive_duration": bson.M{"$sum": bson.A{
			"$v.inactive_duration",
			bson.M{"$cond": bson.M{
				"if": bson.M{"$gt": bson.A{"$v.inactive_start", 0}},
				"then": bson.M{"$subtract": bson.A{
					pause.Timestamp,
					"$v.inactive_start",
				}},
				"else": 0,
			}},
		}},
		"v.last_update_date": now,
	}
	if wasPaused {
		// The replaced pause is not resumed, so its duration is kept.
		set["v.pause_duration"] = bson.M{"$sum": bson.A{
			"$v.pause_duration",
			bson.M{"$subtract": bson.A{pause.Timestamp, "$v.pause.t"}},
		}}
	}

	update := []bson.M{{"$set": set}}

	err = alarmCollection.
		FindOneAndUpdate(ctx, bson.M{"_id": alarm.ID}, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).
		Decode(alarm)
	if err != nil {
		return fmt.Errorf("cannot pause alarm: %w", err)
	}

	return nil
}

func resumeInheritedPause(
	ctx context.Context,
	entityCollection, alarmCollection mongo.DbCollection,
	entity *types.Entity,
	alarm *types.Alarm,
	now datetime.CpsTime,
) error {
	ended := types.LatestPause{Pause: *entity.Pause}
	if entity.LatestPause != nil {
		ended = *entity.LatestPause
	}

	ended.End = &now

	newEntity, err := updateEntityByID(ctx, entity.ID, bson.M{
		"$set":   bson.M{"latest_pause": ended},
		"$unset": bson.M{"pause": ""},
	}, entityCollection)
	if err != nil {
		return err
	}

	*entity = newEntity

	if alarm.ID == "" {
		return nil
	}

	newStepQuery := stepUpdateQueryWithInPbhIntervalByStep(types.NewAlarmStep(types.AlarmStepResume, now,
		"", ended.Message, "", "", types.InitiatorSystem, false))
	update := []bson.M{
		{"$set": bson.M{
			"v.steps": addStepUpdateQuery(newStepQuery),
			"v.pause_duration": bson.M{"$sum": bson.A{
				"$v.pause_duration",
				bson.M{"$subtract": bson.A{now, "$v.pause.t"}},
			}},
			"v.inactive_duration": bson.M{"$sum": bson.A{
				"$v.inactive_duration",
				bson.M{"$subtract": bson.A{now, "$v.inactive_start"}},
			}},
			"v.inactive_start":   updateInactiveStart(now, true, true, true, true, false),
			"v.last_update_date": now,
		}},
		{"$unset": "v.pause"},
	}

	err = alarmCollection.
		FindOneAndUpdate(ctx, bson.M{"_id": alarm.ID}, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).
		Decode(alarm)
	if err != nil {
		return fmt.Errorf("cannot resume alarm: %w", err)
	}

	return nil
}

// publishPauseChange sends the trigger event of the entity alarm and, if the entity is a service
// or a component, a pause or resume event to each of its dependencies, so they apply the change
// and propagate it to their own dependencies.
func publishPauseChange(
	ctx context.Context,
	entityCollection mongo.DbCollection,
	publisher libamqp.Publisher,
	encoder encoding.Encoder,
	eventGenerator libevent.Generator,
	entity types.Entity,
	alarmID string,
	change pauseChange,
	now datetime.CpsTime,
) error {
	if change.alarmChange == "" && change.pause == nil {
		return nil
	}

	// A pause which only replaces the previous pause of the entity is not an alarm change,
	// but it's still propagated to its dependencies.
	if change.alarmChange != "" {
		triggerEvent, err := eventGenerator.Generate(entity)
		if err != nil {
			return fmt.Errorf("cannot generate event: %w", err)
		}

		triggerEvent.Timestamp = now
		triggerEvent.AlarmChange = &types.AlarmChange{Type: change.alarmChange}
		triggerEvent.EventType = types.EventTypeTrigger
		triggerEvent.AlarmID = alarmID
		err = publishEvent(ctx, publisher, encoder, triggerEvent)
		if err != nil {
			return err
		}
	}

	var filter bson.M
	switch entity.Type {
	case types.EntityTypeService:
		filter = bson.M{"services": entity.ID, "enabled": true}
	case types.EntityTypeComponent:
		filter = bson.M{"component": entity.ID, "type": types.EntityTypeResource, "enabled": true}
	default:
		return nil
	}

	cursor, err := entityCollection.Find(ctx, filter)
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		dep := types.Entity{}
		err = cursor.Decode(&dep)
		if err != nil {
			return err
		}

		depEvent, err := eventGenerator.Generate(dep)
		if err != nil {
			return fmt.Errorf("cannot generate event: %w", err)
		}

		depEvent.Inherited = true
		if change.pause != nil {
			depEvent.EventType = types.EventTypePause
			depEvent.Timestamp = change.pause.Timestamp
			depEvent.PauseSource = change.pause.Source
			depEvent.PauseSourceType = change.pause.SourceType
			depEvent.PauseType = change.pause.Type
			depEvent.Output = change.pause.Message
			depEvent.Author = change.pause.Author
			depEvent.UserID = change.pause.UserID
			depEvent.Initiator = change.pause.Initiator
		} else {
			depEvent.EventType = types.EventTypeResume
			depEvent.Timestamp = now
			depEvent.PauseSource = entity.LatestPause.Source
			depEvent.PauseSourceType = entity.LatestPause.SourceType
			depEvent.Initiator = types.InitiatorSystem
		}

		err = publishEvent(ctx, publisher, encoder, depEvent)
		if err != nil {
			return err
		}
	}

	return cursor.Err()
}

func publishEvent(ctx context.Context, publisher libamqp.Publisher, encoder encoding.Encoder, event types.Event) error {
	body, err := encoder.Encode(event)
	if err != nil {
		return fmt.Errorf("cannot encode event: %w", err)
	}

	return publisher.PublishWithContext(
		ctx,
		"",
		canopsis.FIFOQueueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  canopsis.JsonContentType,
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
}
