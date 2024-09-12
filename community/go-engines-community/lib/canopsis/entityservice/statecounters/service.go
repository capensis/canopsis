package statecounters

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	libamqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type service struct {
	dbClient                  mongo.DbClient
	serviceCountersCollection mongo.DbCollection
	entityCollection          mongo.DbCollection
	encoder                   encoding.Encoder
	pubChannel                amqp.Publisher
	templateExecutor          template.Executor
	logger                    zerolog.Logger
	pubExchangeName           string
	pubQueueName              string
}

type UpdatedServicesInfo struct {
	State  int
	Output string
}

func NewStateCountersService(
	client mongo.DbClient,
	pubChannel amqp.Publisher,
	pubExchangeName, pubQueueName string,
	encoder encoding.Encoder,
	templateExecutor template.Executor,
	logger zerolog.Logger,
) StateCountersService {
	return &service{
		dbClient:                  client,
		serviceCountersCollection: client.Collection(mongo.EntityServiceCountersCollection),
		entityCollection:          client.Collection(mongo.EntityMongoCollection),
		encoder:                   encoder,
		pubChannel:                pubChannel,
		logger:                    logger,
		pubExchangeName:           pubExchangeName,
		pubQueueName:              pubQueueName,
		templateExecutor:          templateExecutor,
	}
}

type ServiceCountersConf struct {
	Entity      types.Entity
	Alarm       *types.Alarm
	AlarmChange types.AlarmChange
}

func (s *service) RecomputeAllServices(ctx context.Context) error {
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

		event := types.Event{
			EventType:     types.EventTypeRecomputeEntityService,
			SourceType:    types.SourceTypeService,
			Component:     serv.ID,
			Connector:     types.ConnectorEngineService,
			ConnectorName: types.ConnectorEngineService,
			Timestamp:     types.CpsTime{Time: time.Now()},
			Initiator:     types.InitiatorSystem,
		}

		body, err := s.encoder.Encode(event)
		if err != nil {
			return fmt.Errorf("unable to serialize service event: %w", err)
		}

		err = s.pubChannel.PublishWithContext(
			ctx,
			s.pubExchangeName,
			s.pubQueueName,
			false,
			false,
			libamqp.Publishing{
				Body:        body,
				ContentType: "application/json",
			},
		)
		if err != nil {
			return fmt.Errorf("unable to send service event: %w", err)
		}
	}

	return nil
}

func (s *service) UpdateServiceState(ctx context.Context, serviceID string, serviceInfo UpdatedServicesInfo) error {
	event := types.Event{
		EventType:     types.EventTypeCheck,
		SourceType:    types.SourceTypeService,
		Component:     serviceID,
		Connector:     types.ConnectorEngineService,
		ConnectorName: types.ConnectorEngineService,
		State:         types.CpsNumber(serviceInfo.State),
		Output:        serviceInfo.Output,
		Timestamp:     types.CpsTime{Time: time.Now()},
		Initiator:     types.InitiatorSystem,
	}

	body, err := s.encoder.Encode(event)
	if err != nil {
		return fmt.Errorf("unable to serialize service event: %w", err)
	}

	err = s.pubChannel.PublishWithContext(
		ctx,
		s.pubExchangeName,
		s.pubQueueName,
		false,
		false,
		libamqp.Publishing{
			Body:        body,
			ContentType: "application/json",
		},
	)
	if err != nil {
		return fmt.Errorf("unable to send service event: %w", err)
	}

	return nil
}

func (s *service) UpdateServiceCounters(ctx context.Context, entity types.Entity, alarm *types.Alarm, alarmChange types.AlarmChange) (map[string]UpdatedServicesInfo, error) {
	switch alarmChange.Type {
	case types.AlarmChangeTypeCreate, types.AlarmChangeTypeCreateAndPbhEnter, types.AlarmChangeTypeEnabled,
		types.AlarmChangeTypePbhEnter, types.AlarmChangeTypePbhLeave, types.AlarmChangeTypePbhLeaveAndEnter,
		types.AlarmChangeTypeChangeState, types.AlarmChangeTypeStateDecrease, types.AlarmChangeTypeStateIncrease,
		types.AlarmChangeTypeResolve, types.AlarmChangeTypeAck, types.AlarmChangeTypeAckremove, types.AlarmChangeTypeNone:
	default:
		return nil, nil
	}

	changeType := alarmChange.Type
	prevState := int(alarmChange.PreviousState)
	curState := types.AlarmStateOK
	isActive := entity.PbehaviorInfo.IsActive()
	acked := false
	if alarm != nil {
		curState = int(alarm.CurrentState())
		acked = alarm.IsAck()
	}

	// some services data may be changed if another event for the same entity was processed or because of transaction retry.
	// Always get the fresh data.
	entityServiceData := struct {
		Services         []string `bson:"services"`
		ServicesToAdd    []string `bson:"services_to_add"`
		ServicesToRemove []string `bson:"services_to_remove"`
	}{}
	err := s.entityCollection.FindOne(
		ctx,
		bson.M{"_id": entity.ID},
		options.FindOne().SetProjection(bson.M{
			"services":           1,
			"services_to_add":    1,
			"services_to_remove": 1,
		}),
	).Decode(&entityServiceData)
	if err != nil {
		return nil, err
	}

	updatedServiceInfos := make(map[string]UpdatedServicesInfo)
	updateCountersModels := make(
		[]mongodriver.WriteModel,
		0,
		int(math.Min(
			float64(len(entityServiceData.ServicesToRemove)+len(entityServiceData.Services)),
			canopsis.DefaultBulkSize,
		)),
	)

	servicesToRemove := make(map[string]bool, len(entityServiceData.ServicesToRemove))
	for _, impServ := range entityServiceData.ServicesToRemove {
		servicesToRemove[impServ] = true
	}

	servicesToAdd := make(map[string]bool, len(entityServiceData.ServicesToAdd))
	for _, impServ := range entityServiceData.ServicesToAdd {
		servicesToAdd[impServ] = true
	}

	var inSlice []string
	if changeType == types.AlarmChangeTypeEnabled {
		inSlice = entityServiceData.ServicesToAdd
	} else {
		inSlice = append(entityServiceData.Services, entityServiceData.ServicesToRemove...)
	}

	if len(inSlice) == 0 {
		return updatedServiceInfos, nil
	}

	cursor, err := s.serviceCountersCollection.Find(ctx, bson.M{"_id": bson.M{"$in": inSlice}})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var newModel mongodriver.WriteModel
	bulkBytesSize := 0

	for cursor.Next(ctx) {
		var counters EntityServiceCounters

		err = cursor.Decode(&counters)
		if err != nil {
			return nil, err
		}

		oldCounters := counters
		oldCounters.PbehaviorCounters = make(map[string]int, len(counters.PbehaviorCounters))
		for k, v := range counters.PbehaviorCounters {
			oldCounters.PbehaviorCounters[k] = v
		}

		if counters.PbehaviorCounters == nil {
			counters.PbehaviorCounters = make(map[string]int)
		}

		if servicesToAdd[counters.ID] {
			counters.Depends++
		}

		if servicesToRemove[counters.ID] {
			counters.Depends--
		}

		pbhTypeID := entity.PbehaviorInfo.TypeID

		switch changeType {
		case types.AlarmChangeTypeEnabled:
			if !servicesToRemove[counters.ID] && !entity.PbehaviorInfo.IsActive() {
				counters.IncrementPbhCounters(pbhTypeID)
			}
		case types.AlarmChangeTypeNone:
			if servicesToRemove[counters.ID] {
				if alarm != nil {
					counters.All--
					counters.DecrementAlarmCounters(curState, acked, isActive)
				}

				if !isActive {
					counters.DecrementPbhCounters(pbhTypeID)
				}

				break
			}

			if servicesToAdd[counters.ID] {
				if alarm != nil {
					counters.All++
					counters.IncrementAlarmCounters(curState, acked, isActive)
				}

				if !isActive {
					counters.IncrementPbhCounters(pbhTypeID)
				}
			}
		case types.AlarmChangeTypeCreate, types.AlarmChangeTypeCreateAndPbhEnter:
			if servicesToRemove[counters.ID] {
				break
			}

			counters.All++

			counters.IncrementAlarmCounters(curState, false, isActive)
			if !isActive && alarm != nil && alarmChange.PreviousPbehaviorTypeID != alarm.Value.PbehaviorInfo.TypeID {
				counters.IncrementPbhCounters(pbhTypeID)
			}
		case types.AlarmChangeTypePbhEnter:
			if servicesToRemove[counters.ID] && alarm != nil {
				counters.All--
				counters.DecrementAlarmCounters(curState, acked, true)

				break
			}

			if servicesToAdd[counters.ID] {
				if alarm != nil {
					counters.All++
					counters.IncrementAlarmCounters(curState, acked, isActive)
				}

				if !isActive {
					counters.IncrementPbhCounters(pbhTypeID)
				}

				break
			}

			if !isActive {
				counters.IncrementPbhCounters(pbhTypeID)
				if alarm != nil {
					counters.DecrementAlarmCounters(curState, acked, true)
					counters.IncrementAlarmCounters(types.AlarmStateOK, acked, false)
				}
			}
		case types.AlarmChangeTypePbhLeaveAndEnter:
			if servicesToRemove[counters.ID] {
				if alarm != nil {
					counters.All--
					counters.DecrementAlarmCounters(curState, acked, alarmChange.PreviousPbehaviorCannonicalType == pbehavior.TypeActive)
				}

				if alarmChange.PreviousPbehaviorCannonicalType != pbehavior.TypeActive {
					counters.DecrementPbhCounters(alarmChange.PreviousPbehaviorTypeID)
				}

				break
			}

			if servicesToAdd[counters.ID] {
				if alarm != nil {
					counters.All++
					counters.IncrementAlarmCounters(curState, acked, isActive)
				}

				if !isActive {
					counters.IncrementPbhCounters(pbhTypeID)
				}

				break
			}

			if alarmChange.PreviousPbehaviorCannonicalType != pbehavior.TypeActive {
				counters.DecrementPbhCounters(alarmChange.PreviousPbehaviorTypeID)
			}

			if alarm != nil {
				if isActive && alarmChange.PreviousPbehaviorCannonicalType != pbehavior.TypeActive {
					counters.IncrementAlarmCounters(curState, acked, true)
					counters.DecrementAlarmCounters(types.AlarmStateOK, acked, false)
				}

				if !isActive && alarmChange.PreviousPbehaviorCannonicalType == pbehavior.TypeActive {
					counters.DecrementAlarmCounters(curState, acked, false)
					counters.IncrementAlarmCounters(types.AlarmStateOK, acked, true)
				}
			}

			if !isActive {
				counters.IncrementPbhCounters(pbhTypeID)
			}
		case types.AlarmChangeTypePbhLeave:
			if servicesToRemove[counters.ID] {
				if alarm != nil {
					counters.All--
					counters.DecrementAlarmCounters(curState, acked, alarmChange.PreviousPbehaviorCannonicalType == pbehavior.TypeActive)
				}

				if alarmChange.PreviousPbehaviorCannonicalType != pbehavior.TypeActive {
					counters.DecrementPbhCounters(alarmChange.PreviousPbehaviorTypeID)
				}

				break
			}

			if servicesToAdd[counters.ID] {
				if alarm != nil {
					counters.All++
					counters.IncrementAlarmCounters(curState, acked, true)
				}

				break
			}

			if alarmChange.PreviousPbehaviorCannonicalType != pbehavior.TypeActive {
				if alarm != nil {
					counters.IncrementAlarmCounters(curState, acked, true)
					counters.DecrementAlarmCounters(types.AlarmStateOK, acked, false)
				}

				counters.DecrementPbhCounters(alarmChange.PreviousPbehaviorTypeID)
			}
		case types.AlarmChangeTypeStateIncrease,
			types.AlarmChangeTypeStateDecrease,
			types.AlarmChangeTypeChangeState:
			if servicesToRemove[counters.ID] {
				counters.All--
				counters.DecrementAlarmCounters(prevState, acked, isActive)
				if !isActive {
					counters.DecrementPbhCounters(pbhTypeID)
				}

				break
			}

			if servicesToAdd[counters.ID] {
				counters.All++
				counters.IncrementAlarmCounters(curState, acked, isActive)
				if !isActive {
					counters.IncrementPbhCounters(pbhTypeID)
				}

				break
			}

			if isActive {
				counters.DecrementState(prevState)
				counters.IncrementState(curState)
			}
		case types.AlarmChangeTypeResolve:
			if servicesToAdd[counters.ID] {
				break
			}

			counters.All--
			counters.DecrementAlarmCounters(curState, acked, isActive)

			if !isActive && (servicesToRemove[counters.ID] || !entity.Enabled) {
				counters.DecrementPbhCounters(pbhTypeID)
			}
		case types.AlarmChangeTypeAck:
			if servicesToRemove[counters.ID] {
				counters.All--
				counters.DecrementAlarmCounters(curState, false, isActive)

				if !isActive {
					counters.DecrementPbhCounters(pbhTypeID)
				}

				break
			}

			if servicesToAdd[counters.ID] {
				counters.All++
				counters.IncrementAlarmCounters(curState, true, isActive)

				if !isActive {
					counters.IncrementPbhCounters(pbhTypeID)
				}

				break
			}

			if isActive {
				counters.Acknowledged++
				counters.NotAcknowledged--
			} else {
				counters.AcknowledgedUnderPbh++
			}
		case types.AlarmChangeTypeAckremove:
			if servicesToRemove[counters.ID] {
				counters.All--
				counters.DecrementAlarmCounters(curState, true, isActive)

				if !isActive {
					counters.DecrementPbhCounters(pbhTypeID)
				}

				break
			}

			if servicesToAdd[counters.ID] {
				counters.All++
				counters.IncrementAlarmCounters(curState, false, isActive)
				if !isActive {
					counters.IncrementPbhCounters(pbhTypeID)
				}

				break
			}

			if isActive {
				counters.Acknowledged--
				counters.NotAcknowledged++
			} else {
				counters.AcknowledgedUnderPbh--
			}
		}

		countersChanged := false
		diffCounters := counters.Sub(oldCounters)
		for _, v := range diffCounters {
			if v != 0 {
				countersChanged = true
				break
			}
		}

		if !countersChanged {
			continue
		}

		output, err := s.templateExecutor.Execute(counters.OutputTemplate, counters)
		if err != nil {
			return nil, err
		}

		updatedServiceInfos[counters.ID] = UpdatedServicesInfo{
			State:  counters.GetWorstState(),
			Output: output,
		}

		newModel = mongodriver.
			NewUpdateOneModel().
			SetFilter(bson.M{"_id": counters.ID}).
			SetUpdate(bson.M{"$inc": diffCounters}).
			SetUpsert(true)

		b, err := bson.Marshal(newModel)
		if err != nil {
			return nil, err
		}

		newModelLen := len(b)
		if bulkBytesSize+newModelLen > canopsis.DefaultBulkBytesSize {
			_, err = s.serviceCountersCollection.BulkWrite(ctx, updateCountersModels)
			if err != nil {
				return nil, err
			}

			updateCountersModels = updateCountersModels[:0]
			bulkBytesSize = 0
		}

		bulkBytesSize += newModelLen
		updateCountersModels = append(
			updateCountersModels,
			newModel,
		)

		if len(updateCountersModels) == canopsis.DefaultBulkSize {
			_, err = s.serviceCountersCollection.BulkWrite(ctx, updateCountersModels)
			if err != nil {
				return nil, err
			}

			updateCountersModels = updateCountersModels[:0]
		}
	}

	if len(updateCountersModels) > 0 {
		_, err = s.serviceCountersCollection.BulkWrite(ctx, updateCountersModels)
		if err != nil {
			return nil, err
		}
	}

	if len(entityServiceData.ServicesToAdd) > 0 || len(entityServiceData.ServicesToRemove) > 0 {
		_, err = s.entityCollection.UpdateOne(
			ctx, bson.M{"_id": entity.ID},
			bson.M{"$unset": bson.M{"services_to_add": 1, "services_to_remove": 1}})
		if err != nil {
			return nil, err
		}
	}

	return updatedServiceInfos, nil
}

func (s *service) RecomputeEntityServiceCounters(ctx context.Context, entity types.Entity) (map[string]UpdatedServicesInfo, error) {
	updatedServiceStates := make(map[string]UpdatedServicesInfo)
	var counters EntityServiceCounters
	err := s.serviceCountersCollection.FindOne(ctx, bson.M{"_id": entity.ID}).Decode(&counters)
	if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
		return nil, err
	}

	counters = EntityServiceCounters{
		ID:                entity.ID,
		OutputTemplate:    counters.OutputTemplate,
		PbehaviorCounters: make(map[string]int),
	}

	cursor, err := s.entityCollection.Aggregate(ctx, []bson.M{
		{
			"$match": bson.M{"services": entity.ID},
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

		counters.Depends++
		isActive := depEnt.Entity.PbehaviorInfo.IsActive()

		if depEnt.Alarm.ID != "" {
			if depEnt.Alarm.IsResolved() {
				continue
			}

			counters.All++
			counters.IncrementAlarmCounters(int(depEnt.Alarm.CurrentState()), depEnt.Alarm.IsAck(), isActive)
		}

		if !isActive {
			counters.IncrementPbhCounters(depEnt.Entity.PbehaviorInfo.TypeID)
		}
	}

	_, err = s.serviceCountersCollection.UpdateOne(ctx, bson.M{"_id": entity.ID}, bson.M{"$set": counters}, options.Update().SetUpsert(true))
	if err != nil {
		return nil, err
	}

	output, err := s.templateExecutor.Execute(counters.OutputTemplate, counters)
	if err != nil {
		return nil, err
	}

	updatedServiceStates[entity.ID] = UpdatedServicesInfo{
		State:  counters.GetWorstState(),
		Output: output,
	}

	return updatedServiceStates, nil
}
