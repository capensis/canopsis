package calculator

import (
	"context"
	"errors"
	"fmt"
	"math"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters/strategy/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type entityServiceCountersCalculator struct {
	serviceCountersCollection mongo.DbCollection
	entityCollection          mongo.DbCollection
	templateExecutor          template.Executor
	eventsSender              entitycounters.EventsSender

	options *options.FindOneOptions
}

func NewEntityServiceCountersCalculator(dbClient mongo.DbClient, executor template.Executor, eventsSender entitycounters.EventsSender) EntityServiceCountersCalculator {
	return &entityServiceCountersCalculator{
		serviceCountersCollection: dbClient.Collection(mongo.EntityCountersCollection),
		entityCollection:          dbClient.Collection(mongo.EntityMongoCollection),
		templateExecutor:          executor,
		eventsSender:              eventsSender,

		options: options.FindOne().SetProjection(
			bson.M{
				"services":           1,
				"services_to_add":    1,
				"services_to_remove": 1,
			},
		),
	}
}

func (s *entityServiceCountersCalculator) getServicesInfo(ctx context.Context, entityID string) (entitycounters.ServicesInfo, error) {
	var info entitycounters.ServicesInfo
	err := s.entityCollection.FindOne(ctx, bson.M{"_id": entityID}, s.options).Decode(&info)
	if err != nil {
		return entitycounters.ServicesInfo{}, err
	}

	return info, nil
}

func (s *entityServiceCountersCalculator) RecomputeCounters(ctx context.Context, service *types.Entity) (map[string]entitycounters.UpdatedServicesInfo, error) {
	var counters entitycounters.EntityCounters
	err := s.serviceCountersCollection.FindOne(ctx, bson.M{"_id": service.ID}).Decode(&counters)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	counters.Reset()

	cursor, err := s.entityCollection.Aggregate(ctx, []bson.M{
		{
			"$match": bson.M{"services": service.ID},
		},
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
				"entity":        1,
				"alarm._id":     1,
				"alarm.v.state": 1,
				"alarm.v.ack":   1,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var depEnt types.AlarmWithEntity
		err := cursor.Decode(&depEnt)
		if err != nil {
			return nil, err
		}

		inherited := false
		if counters.Rule != nil && counters.Rule.InheritedEntityPattern != nil {
			inherited, err = match.MatchEntityPattern(*counters.Rule.InheritedEntityPattern, &depEnt.Entity)
			if err != nil {
				return nil, err
			}
		}

		curActive := depEnt.Entity.PbehaviorInfo.IsActive()
		curPbhTypeID := depEnt.Entity.PbehaviorInfo.TypeID

		curState := types.AlarmStateOK
		if depEnt.Alarm.ID != "" && curActive {
			curState = int(depEnt.Alarm.CurrentState())
		}

		isAcked := depEnt.Alarm.ID != "" && depEnt.Alarm.IsAck()

		counters.Depends++
		counters.IncrementState(curState, inherited)
		if !curActive {
			counters.IncrementPbhCounters(curPbhTypeID)
		}

		if depEnt.Alarm.ID != "" {
			counters.IncrementAlarmCounters(isAcked, curActive)
		}
	}

	counters.Output, err = s.templateExecutor.Execute(counters.OutputTemplate, counters)
	if err != nil {
		return nil, err
	}

	_, err = s.serviceCountersCollection.UpdateOne(ctx, bson.M{"_id": service.ID}, bson.M{"$set": counters}, options.Update().SetUpsert(true))
	if err != nil {
		return nil, err
	}

	updatedServiceStates := make(map[string]entitycounters.UpdatedServicesInfo)
	updatedServiceStates[service.ID] = entitycounters.UpdatedServicesInfo{
		State:  counters.GetWorstState(),
		Output: counters.Output,
	}

	return updatedServiceStates, nil
}

func (s *entityServiceCountersCalculator) RecomputeAll(ctx context.Context) error {
	cursor, err := s.entityCollection.Find(ctx, bson.M{"enabled": true, "type": types.EntityTypeService})
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var serv types.Entity
		err = cursor.Decode(&serv)
		if err != nil {
			return fmt.Errorf("unable to decode entity service document: %w", err)
		}

		err = s.eventsSender.RecomputeService(ctx, serv.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *entityServiceCountersCalculator) CalculateCounters(
	ctx context.Context,
	entity *types.Entity,
	alarm *types.Alarm,
	alarmChange types.AlarmChange,
) (bool, map[string]entitycounters.UpdatedServicesInfo, error) {
	if entity == nil {
		return false, nil, nil
	}

	var strategy EntityServiceCountersStrategy

	switch alarmChange.Type {
	case types.AlarmChangeTypeNone:
		strategy = entityservice.NoChangeStrategy{}
	case types.AlarmChangeTypeCreate:
		if alarm == nil || alarm.ID == "" {
			return false, nil, nil
		}

		strategy = entityservice.CreateStrategy{}
	case types.AlarmChangeTypeCreateAndPbhEnter:
		if alarm == nil || alarm.ID == "" {
			return false, nil, nil
		}

		strategy = entityservice.CreateAndPbhEnterStrategy{}
	case types.AlarmChangeTypePbhEnter:
		strategy = entityservice.PbhEnterStrategy{}
	case types.AlarmChangeTypePbhLeave:
		strategy = entityservice.PbhLeaveStrategy{}
	case types.AlarmChangeTypePbhLeaveAndEnter:
		strategy = entityservice.PbhLeaveAndEnterStrategy{}
	case types.AlarmChangeTypeStateDecrease, types.AlarmChangeTypeStateIncrease, types.AlarmChangeTypeChangeState:
		if alarm == nil || alarm.ID == "" {
			return false, nil, nil
		}

		strategy = entityservice.ChangeStateStrategy{}
	case types.AlarmChangeTypeResolve:
		if alarm == nil || alarm.ID == "" {
			return false, nil, nil
		}

		strategy = entityservice.ResolveStrategy{}
	case types.AlarmChangeTypeAck:
		if alarm == nil || alarm.ID == "" {
			return false, nil, nil
		}

		strategy = entityservice.AckStrategy{}
	case types.AlarmChangeTypeAckremove:
		if alarm == nil || alarm.ID == "" {
			return false, nil, nil
		}

		strategy = entityservice.AckRemoveStrategy{}
	case types.AlarmChangeTypeEnabled:
		// on toggle event when entity becomes enabled, nothing actually happens with an alarm, so it is basically the same as noChange strategy.
		strategy = entityservice.NoChangeStrategy{}
	default:
		return false, nil, nil
	}

	return s.calculateCounters(ctx, entity, alarm, alarmChange, strategy)
}

func (s *entityServiceCountersCalculator) calculateCounters(
	ctx context.Context,
	entity *types.Entity,
	alarm *types.Alarm,
	alarmChange types.AlarmChange,
	strategy EntityServiceCountersStrategy,
) (isAnyServiceCountersUpdated bool, updatedServicesInfos map[string]entitycounters.UpdatedServicesInfo, _ error) {
	// Some services data may be changed if another event for the same entity was processed or because of transaction retry.
	// Always get the fresh data.
	info, err := s.getServicesInfo(ctx, entity.ID)
	if err != nil {
		return false, nil, err
	}

	if len(info.Services) == 0 && len(info.ServicesToRemove) == 0 {
		return false, nil, err
	}

	var calcData entitycounters.EntityServiceCountersCalcData

	calcData.PrevActive = alarmChange.PreviousPbehaviorCannonicalType == types.PbhCanonicalTypeActive || alarmChange.PreviousPbehaviorCannonicalType == ""
	calcData.CurActive = entity.PbehaviorInfo.IsActive()
	calcData.CurPbhTypeID = entity.PbehaviorInfo.TypeID
	calcData.PrevPbhTypeID = alarmChange.PreviousPbehaviorTypeID
	calcData.EntityEnabled = entity.Enabled
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

	if calcData.AlarmExists {
		calcData.IsAcked = alarm.IsAck()
	}

	calcData.ServicesToRemove = make(map[string]bool, len(info.ServicesToRemove))
	for _, impServ := range info.ServicesToRemove {
		calcData.ServicesToRemove[impServ] = true
	}

	calcData.ServicesToAdd = make(map[string]bool, len(info.ServicesToAdd))
	for _, impServ := range info.ServicesToAdd {
		calcData.ServicesToAdd[impServ] = true
	}

	if strategy.CanSkip(calcData) {
		return false, nil, err
	}

	serviceIDs := append(info.Services, info.ServicesToRemove...)
	countersUpdated := false
	updatedServiceInfos := make(map[string]entitycounters.UpdatedServicesInfo)
	updateCountersModels := make(
		[]mongodriver.WriteModel,
		0,
		int(math.Min(
			float64(len(serviceIDs)),
			canopsis.DefaultBulkSize,
		)),
	)

	cursor, err := s.serviceCountersCollection.Find(ctx, bson.M{"_id": bson.M{"$in": serviceIDs}})
	if err != nil {
		return false, nil, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var counters entitycounters.EntityCounters
		err = cursor.Decode(&counters)
		if err != nil {
			return false, nil, err
		}

		if counters.Rule != nil && counters.Rule.InheritedEntityPattern != nil {
			calcData.Inherited, err = match.MatchEntityPattern(*counters.Rule.InheritedEntityPattern, entity)
			if err != nil {
				return false, nil, err
			}
		}

		calcData.Counters = counters.Copy()
		newCounters := strategy.Calculate(calcData)
		newOutput, err := s.templateExecutor.Execute(newCounters.OutputTemplate, newCounters)
		if err != nil {
			return false, nil, err
		}

		countersUpdated = true
		updateCountersModels = append(
			updateCountersModels,
			mongodriver.
				NewUpdateOneModel().
				SetFilter(bson.M{"_id": newCounters.ID}).
				SetUpdate(bson.M{"$inc": newCounters.Sub(counters), "$set": bson.M{"output": newOutput}}).
				SetUpsert(true),
		)
		newWorstState := newCounters.GetWorstState()
		if counters.Output == newOutput && counters.GetWorstState() == newWorstState {
			continue
		}

		updatedServiceInfos[counters.ID] = entitycounters.UpdatedServicesInfo{
			State:  newWorstState,
			Output: newOutput,
		}

		if len(updateCountersModels) == canopsis.DefaultBulkSize {
			_, err = s.serviceCountersCollection.BulkWrite(ctx, updateCountersModels)
			if err != nil {
				return false, nil, err
			}

			updateCountersModels = updateCountersModels[:0]
		}
	}

	if len(updateCountersModels) > 0 {
		_, err = s.serviceCountersCollection.BulkWrite(ctx, updateCountersModels)
		if err != nil {
			return false, nil, err
		}
	}

	_, err = s.entityCollection.UpdateOne(
		ctx,
		bson.M{"_id": entity.ID},
		bson.M{"$unset": bson.M{"services_to_add": 1, "services_to_remove": 1}})
	if err != nil {
		return false, nil, err
	}

	return countersUpdated, updatedServiceInfos, nil
}
