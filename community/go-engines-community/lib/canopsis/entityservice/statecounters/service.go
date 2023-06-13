package statecounters

import (
	"context"
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
		serviceCountersCollection: client.Collection(mongo.EntityServiceCountersMongoCollection),
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
			Component:     serv.Name,
			Connector:     types.ConnectorEngineService,
			ConnectorName: types.ConnectorEngineService,
			Timestamp:     types.CpsTime{Time: time.Now()},
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

	prevState := int(alarmChange.PreviousState)
	curState := types.AlarmStateOK
	isActive := entity.PbehaviorInfo.IsActive()
	acked := false
	if alarm != nil {
		curState = int(alarm.CurrentState())
		acked = alarm.IsAck()
	}
	changeType := alarmChange.Type

	updatedServiceInfos := make(map[string]UpdatedServicesInfo)
	updateCountersModels := make(
		[]mongodriver.WriteModel,
		0,
		int(math.Min(
			float64(len(entity.ServicesToRemove)+len(entity.Services)),
			canopsis.DefaultBulkSize,
		)),
	)

	servicesToRemove := make(map[string]bool, len(entity.ServicesToRemove))
	for _, impServ := range entity.ServicesToRemove {
		servicesToRemove[impServ] = true
	}

	servicesToAdd := make(map[string]bool, len(entity.ServicesToAdd))
	for _, impServ := range entity.ServicesToAdd {
		servicesToAdd[impServ] = true
	}

	var cursor mongo.Cursor
	var err error

	var inSlice []string
	if changeType == types.AlarmChangeTypeEnabled {
		inSlice = entity.ServicesToAdd
	} else {
		inSlice = append(entity.Services, entity.ServicesToRemove...)
	}

	if len(inSlice) == 0 {
		return updatedServiceInfos, nil
	}

	cursor, err = s.serviceCountersCollection.Find(context.Background(), bson.M{"_id": bson.M{"$in": inSlice}})
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
				continue
			}

			counters.All++

			counters.IncrementAlarmCounters(curState, acked, isActive)
			if !isActive {
				if alarmChange.PreviousPbehaviorTypeID != alarm.Value.PbehaviorInfo.TypeID {
					counters.IncrementPbhCounters(pbhTypeID)
				}
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
				continue
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
				counters.IncrementAlarmCounters(curState, acked, isActive)
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
				counters.IncrementAlarmCounters(curState, acked, isActive)
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
			SetUpdate(bson.M{"$set": counters}).
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

	if len(entity.ServicesToAdd) > 0 || len(entity.ServicesToRemove) > 0 {
		_, err = s.entityCollection.UpdateOne(
			ctx, bson.M{"_id": entity.ID},
			bson.M{"$unset": bson.M{"services_to_add": 1, "services_to_remove": 1}})
		if err != nil {
			return nil, err
		}
	}

	return updatedServiceInfos, nil
}

func (s *service) RecomputeEntityServiceCounters(ctx context.Context, event types.Event) (map[string]UpdatedServicesInfo, error) {
	if event.Entity == nil {
		return nil, nil
	}

	updatedServiceStates := make(map[string]UpdatedServicesInfo)
	var counters EntityServiceCounters
	err := s.serviceCountersCollection.FindOne(ctx, bson.M{"_id": event.Entity.ID}).Decode(&counters)
	if err != nil && err != mongodriver.ErrNoDocuments {
		return nil, err
	}

	counters = EntityServiceCounters{
		ID:                event.Entity.ID,
		OutputTemplate:    counters.OutputTemplate,
		PbehaviorCounters: make(map[string]int),
	}

	cursor, err := s.entityCollection.Aggregate(ctx, []bson.M{
		{
			"$match": bson.M{"services": event.Entity.ID},
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

	_, err = s.serviceCountersCollection.UpdateOne(ctx, bson.M{"_id": event.GetEID()}, bson.M{"$set": counters}, options.Update().SetUpsert(true))
	if err != nil {
		return nil, err
	}

	output, err := s.templateExecutor.Execute(counters.OutputTemplate, counters)
	if err != nil {
		return nil, err
	}

	updatedServiceStates[event.Entity.ID] = UpdatedServicesInfo{
		State:  counters.GetWorstState(),
		Output: output,
	}

	return updatedServiceStates, nil
}
