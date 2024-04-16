package calculator

import (
	"context"
	"errors"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters/strategy/component"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/db"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type componentCountersCalculator struct {
	entityCountersCollection mongo.DbCollection
	entityCollection         mongo.DbCollection
	eventsSender             entitycounters.EventsSender

	options *options.FindOneOptions
}

func NewComponentCountersCalculator(dbClient mongo.DbClient, eventsSender entitycounters.EventsSender) ComponentCountersCalculator {
	return &componentCountersCalculator{
		entityCountersCollection: dbClient.Collection(mongo.EntityCountersCollection),
		entityCollection:         dbClient.Collection(mongo.EntityMongoCollection),
		eventsSender:             eventsSender,
		options: options.FindOne().SetProjection(
			bson.M{
				"component_state_settings":           1,
				"component_state_settings_to_add":    1,
				"component_state_settings_to_remove": 1,
			},
		),
	}
}

func (s *componentCountersCalculator) RecomputeCounters(ctx context.Context, component *types.Entity) (int, error) {
	if component == nil {
		return 0, nil
	}

	var counters entitycounters.EntityCounters
	err := s.entityCountersCollection.FindOne(ctx, bson.M{"_id": component.ID}).Decode(&counters)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return 0, nil
		}

		return 0, err
	}

	counters.Reset()

	match := []bson.M{
		{"$match": bson.M{"component": component.ID, "type": types.EntityTypeResource}},
	}

	if counters.Rule.InheritedEntityPattern != nil {
		patternMongoQuery, err := db.EntityPatternToMongoQuery(*counters.Rule.InheritedEntityPattern, "")
		if err != nil {
			return 0, err
		}

		match = append(match, bson.M{"$match": patternMongoQuery})
	}

	cursor, err := s.entityCollection.Aggregate(ctx, append(match, []bson.M{
		{
			"$project": bson.M{
				"entity": "$$ROOT",
				"_id":    0,
			},
		},
		{
			"$lookup": bson.M{
				"from": mongo.AlarmMongoCollection,
				"let":  bson.M{"id": "$entity._id"},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$and": []bson.M{
								{"$expr": bson.M{"$eq": bson.A{"$d", "$$id"}}},
								{"v.resolved": nil},
							},
						},
					},
					{
						"$limit": 1,
					},
				},
				"as": "alarm",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$alarm",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$project": bson.M{
				"entity.pbehavior_info": 1,
				"alarm._id":             1,
				"alarm.v.state":         1,
			},
		},
	}...))
	if err != nil {
		return 0, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var depEnt types.AlarmWithEntity
		err := cursor.Decode(&depEnt)
		if err != nil {
			return 0, err
		}

		curActive := depEnt.Entity.PbehaviorInfo.IsActive()
		curState := types.AlarmStateOK

		if depEnt.Alarm.ID != "" && curActive {
			curState = int(depEnt.Alarm.CurrentState())
		}

		counters.IncrementState(curState, false)
	}

	_, err = s.entityCountersCollection.UpdateOne(
		ctx,
		bson.M{"_id": component.ID},
		bson.M{"$set": counters},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return 0, err
	}

	return counters.GetWorstState(), nil
}

func (s *componentCountersCalculator) RecomputeAll(ctx context.Context) error {
	cursor, err := s.entityCollection.Find(ctx, bson.M{
		"enabled":    true,
		"type":       types.EntityTypeComponent,
		"state_info": bson.M{"$nin": bson.A{"", nil}},
	}, options.Find().SetProjection(bson.M{"_id": 1, "connector": 1}))
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var comp types.Entity
		err = cursor.Decode(&comp)
		if err != nil {
			return fmt.Errorf("unable to decode entity service document: %w", err)
		}

		err = s.eventsSender.RecomputeComponent(ctx, comp.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *componentCountersCalculator) CalculateCounters(
	ctx context.Context,
	entity *types.Entity,
	alarm *types.Alarm,
	alarmChange types.AlarmChange,
) (bool, bool, int, error) {
	if entity == nil || entity.Type != types.EntityTypeResource {
		return false, false, 0, nil
	}

	var strategy ComponentCountersStrategy

	switch alarmChange.Type {
	case types.AlarmChangeTypeNone:
		strategy = component.NoChangeStrategy{}
	case types.AlarmChangeTypeCreate:
		if alarm == nil || alarm.ID == "" {
			return false, false, 0, nil
		}

		strategy = component.CreateStrategy{}
	case types.AlarmChangeTypeCreateAndPbhEnter:
		if alarm == nil || alarm.ID == "" {
			return false, false, 0, nil
		}

		strategy = component.CreateAndPbhEnterStrategy{}
	case types.AlarmChangeTypeStateDecrease, types.AlarmChangeTypeStateIncrease, types.AlarmChangeTypeChangeState:
		strategy = component.ChangeStateStrategy{}
	case types.AlarmChangeTypePbhEnter:
		strategy = component.PbhEnterStrategy{}
	case types.AlarmChangeTypePbhLeave:
		strategy = component.PbhLeaveStrategy{}
	case types.AlarmChangeTypePbhLeaveAndEnter:
		strategy = component.PbhLeaveAndEnterStrategy{}
	case types.AlarmChangeTypeResolve:
		if alarm == nil || alarm.ID == "" {
			return false, false, 0, nil
		}

		strategy = component.ResolveStrategy{}
	default:
		return false, false, 0, nil
	}

	return s.calculateCounters(ctx, entity, alarm, alarmChange, strategy)
}

func (s *componentCountersCalculator) calculateCounters(
	ctx context.Context,
	entity *types.Entity,
	alarm *types.Alarm,
	alarmChange types.AlarmChange,
	strategy ComponentCountersStrategy,
) (bool, bool, int, error) {
	var calcData entitycounters.ComponentCountersCalcData
	var err error

	// Some services data may be changed if another event for the same entity was processed or because of transaction retry.
	// Always get the fresh data.
	calcData.Info, err = s.getStateSettingsInfo(ctx, entity.ID)
	if err != nil {
		return false, false, 0, err
	}

	if !calcData.Info.ComponentStateSettings && !calcData.Info.ComponentStateSettingsToRemove {
		return false, false, 0, nil
	}

	calcData.PrevActive = alarmChange.PreviousPbehaviorCannonicalType == types.PbhCanonicalTypeActive || alarmChange.PreviousPbehaviorCannonicalType == ""
	calcData.CurActive = entity.PbehaviorInfo.IsActive()
	calcData.AlarmExists = alarm != nil && alarm.ID != ""

	// todo: garbage condition, but it works, is it possible to simplify?
	if calcData.AlarmExists && calcData.PrevActive {
		if alarmChange.Type == types.AlarmChangeTypeStateDecrease ||
			alarmChange.Type == types.AlarmChangeTypeStateIncrease ||
			alarmChange.Type == types.AlarmChangeTypeChangeState {
			if calcData.CurActive {
				calcData.PrevState = int(alarmChange.PreviousState)
			}
		} else if alarmChange.Type == types.AlarmChangeTypeResolve {
			if calcData.CurActive {
				calcData.PrevState = int(alarm.CurrentState())
			}
		} else {
			calcData.PrevState = int(alarm.CurrentState())
		}
	}

	if calcData.AlarmExists && calcData.CurActive {
		calcData.CurState = int(alarm.CurrentState())
	}

	// each alarm change have some conditions where we can be 100% sure that counters won't be changed, so we can
	// avoid extra calculation and updates and skip
	if strategy.CanSkip(calcData) {
		return false, false, 0, nil
	}

	var counters entitycounters.EntityCounters
	err = s.entityCountersCollection.FindOne(ctx, bson.M{"_id": entity.Component}).Decode(&counters)
	if err != nil {
		return false, false, 0, err
	}

	calcData.Counters = counters.Copy()
	newCounters := strategy.Calculate(calcData)
	diff := newCounters.Sub(counters)
	if len(diff) == 0 {
		return false, false, 0, nil
	}

	_, err = s.entityCountersCollection.UpdateOne(
		ctx,
		bson.M{"_id": entity.Component},
		bson.M{"$inc": diff},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return false, false, 0, err
	}

	if calcData.Info.ComponentStateSettingsToAdd || calcData.Info.ComponentStateSettingsToRemove {
		_, err := s.entityCollection.UpdateOne(
			ctx,
			bson.M{"_id": entity.ID},
			bson.M{"$unset": bson.M{"component_state_settings_to_add": 1, "component_state_settings_to_remove": 1}})
		if err != nil {
			return false, false, 0, err
		}
	}

	newWorstState := newCounters.GetWorstState()
	if counters.GetWorstState() != newWorstState {
		return true, true, newWorstState, nil
	}

	return true, false, 0, nil
}

func (s *componentCountersCalculator) getStateSettingsInfo(ctx context.Context, entityID string) (entitycounters.StateSettingsInfo, error) {
	var info entitycounters.StateSettingsInfo
	err := s.entityCollection.FindOne(ctx, bson.M{"_id": entityID}, s.options).Decode(&info)
	if err != nil {
		return entitycounters.StateSettingsInfo{}, err
	}

	return info, nil
}
