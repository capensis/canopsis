package statecounters

import (
	"context"
	"fmt"
	"math"
	"strings"
	"text/template"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
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

		err = s.pubChannel.Publish(
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

func (s *service) UpdateServiceState(serviceID string, serviceInfo UpdatedServicesInfo) error {
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

	err = s.pubChannel.Publish(
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
			float64(len(entity.ImpactedServicesToRemove)+len(entity.ImpactedServices)),
			canopsis.DefaultBulkSize,
		)),
	)

	serviceToRemove := make(map[string]bool, len(entity.ImpactedServicesToRemove))
	for _, impServ := range entity.ImpactedServicesToRemove {
		serviceToRemove[impServ] = true
	}

	serviceToAdd := make(map[string]bool, len(entity.ImpactedServicesToAdd))
	for _, impServ := range entity.ImpactedServicesToAdd {
		serviceToAdd[impServ] = true
	}

	var cursor mongo.Cursor
	var err error

	var inSlice []string
	if changeType == types.AlarmChangeTypeEnabled {
		inSlice = entity.ImpactedServicesToAdd
	} else {
		inSlice = append(entity.ImpactedServices, entity.ImpactedServicesToRemove...)
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
			counters.PbehaviorCounters = make(map[string]int64)
		}

		switch changeType {
		case types.AlarmChangeTypeEnabled:
			if !serviceToRemove[counters.ID] && !entity.PbehaviorInfo.IsActive() {
				counters.PbehaviorCounters[entity.PbehaviorInfo.TypeID]++
			}
		case types.AlarmChangeTypeNone:
			if serviceToRemove[counters.ID] {
				if alarm != nil {
					counters.All--
					if isActive {
						counters.DecrementAlarmCounters(curState, acked)
					}
				}

				if !isActive {
					counters.PbehaviorCounters[entity.PbehaviorInfo.TypeID]--
				}

				break
			}

			if serviceToAdd[counters.ID] {
				if alarm != nil {
					counters.All++
					if isActive {
						counters.IncrementAlarmCounters(curState, acked)
					}
				}

				if !isActive {
					counters.PbehaviorCounters[entity.PbehaviorInfo.TypeID]++
				}
			}
		case types.AlarmChangeTypeCreate, types.AlarmChangeTypeCreateAndPbhEnter:
			if serviceToRemove[counters.ID] {
				continue
			}

			counters.All++
			if isActive {
				counters.IncrementAlarmCounters(curState, acked)
			} else if alarmChange.PreviousPbehaviorTypeID != alarm.Value.PbehaviorInfo.TypeID {
				counters.PbehaviorCounters[entity.PbehaviorInfo.TypeID]++
			}
		case types.AlarmChangeTypePbhEnter:
			if serviceToRemove[counters.ID] && alarm != nil {
				counters.All--
				counters.DecrementAlarmCounters(curState, acked)

				break
			}

			if serviceToAdd[counters.ID] {
				if alarm != nil {
					counters.All++
					if isActive {
						counters.IncrementAlarmCounters(curState, acked)
					}
				}

				if !isActive {
					counters.PbehaviorCounters[entity.PbehaviorInfo.TypeID]++
				}

				break
			}

			if alarm != nil && !isActive {
				counters.DecrementAlarmCounters(curState, acked)
				counters.PbehaviorCounters[entity.PbehaviorInfo.TypeID]++
			} else if !isActive {
				counters.PbehaviorCounters[entity.PbehaviorInfo.TypeID]++
			}
		case types.AlarmChangeTypePbhLeaveAndEnter:
			if serviceToRemove[counters.ID] {
				if alarm != nil {
					counters.All--
					if alarmChange.PreviousPbehaviorCannonicalType != pbehavior.TypeActive {
						counters.DecrementAlarmCounters(curState, acked)
					}
				}

				if alarmChange.PreviousPbehaviorCannonicalType != pbehavior.TypeActive {
					counters.PbehaviorCounters[alarmChange.PreviousPbehaviorTypeID]--
				}

				break
			}

			if serviceToAdd[counters.ID] {
				if alarm != nil {
					counters.All++
					if isActive {
						counters.IncrementAlarmCounters(curState, acked)
					}
				}

				if !isActive {
					counters.PbehaviorCounters[entity.PbehaviorInfo.TypeID]++
				}

				break
			}

			if alarmChange.PreviousPbehaviorCannonicalType != pbehavior.TypeActive {
				counters.PbehaviorCounters[alarmChange.PreviousPbehaviorTypeID]--
			}

			if alarm != nil {
				if isActive {
					if alarmChange.PreviousPbehaviorCannonicalType != pbehavior.TypeActive {
						counters.IncrementAlarmCounters(curState, acked)
					}
				} else {
					if alarmChange.PreviousPbehaviorCannonicalType == pbehavior.TypeActive {
						counters.DecrementAlarmCounters(curState, acked)
					}
				}
			}

			if !isActive {
				counters.PbehaviorCounters[entity.PbehaviorInfo.TypeID]++
			}
		case types.AlarmChangeTypePbhLeave:
			if serviceToRemove[counters.ID] {
				if alarm != nil {
					counters.All--
				}

				if alarmChange.PreviousPbehaviorCannonicalType == pbehavior.TypeActive {
					if alarm != nil {
						counters.DecrementAlarmCounters(curState, acked)
					}
				} else {
					counters.PbehaviorCounters[alarmChange.PreviousPbehaviorTypeID]--
				}

				break
			}

			if serviceToAdd[counters.ID] {
				if alarm != nil {
					counters.All++
					counters.IncrementAlarmCounters(curState, acked)
				}

				break
			}

			if alarmChange.PreviousPbehaviorCannonicalType != pbehavior.TypeActive {
				if alarm != nil {
					counters.IncrementAlarmCounters(curState, acked)
				}

				counters.PbehaviorCounters[alarmChange.PreviousPbehaviorTypeID]--
			}
		case types.AlarmChangeTypeStateIncrease,
			types.AlarmChangeTypeStateDecrease,
			types.AlarmChangeTypeChangeState:
			if serviceToRemove[counters.ID] {
				counters.All--
				if isActive {
					counters.DecrementAlarmCounters(prevState, acked)
				} else {
					counters.PbehaviorCounters[entity.PbehaviorInfo.TypeID]--
				}

				break
			}

			if serviceToAdd[counters.ID] {
				counters.All++
				if isActive {
					counters.IncrementAlarmCounters(curState, acked)
				} else {
					counters.PbehaviorCounters[entity.PbehaviorInfo.TypeID]++
				}

				break
			}

			if isActive {
				counters.DecrementState(prevState)
				counters.IncrementState(curState)
			}
		case types.AlarmChangeTypeResolve:
			if serviceToAdd[counters.ID] {
				continue
			}

			counters.All--
			if isActive {
				counters.DecrementAlarmCounters(curState, acked)
			} else if serviceToRemove[counters.ID] || !entity.Enabled {
				counters.PbehaviorCounters[entity.PbehaviorInfo.TypeID]--
			}
		case types.AlarmChangeTypeAck:
			if serviceToRemove[counters.ID] {
				if isActive {
					counters.DecrementState(curState)
					counters.NotAcknowledged--
				} else {
					counters.PbehaviorCounters[entity.PbehaviorInfo.TypeID]--
				}

				break
			}

			if serviceToAdd[counters.ID] {
				if isActive {
					counters.IncrementState(curState)
					counters.Acknowledged++
				} else {
					counters.PbehaviorCounters[entity.PbehaviorInfo.TypeID]++
				}

				break
			}

			counters.Acknowledged++
			counters.NotAcknowledged--
		case types.AlarmChangeTypeAckremove:
			if serviceToRemove[counters.ID] {
				if isActive {
					counters.DecrementState(curState)
					counters.Acknowledged--
				} else {
					counters.PbehaviorCounters[entity.PbehaviorInfo.TypeID]--
				}

				break
			}

			if serviceToAdd[counters.ID] {
				if isActive {
					counters.IncrementState(curState)
					counters.NotAcknowledged++
				} else {
					counters.PbehaviorCounters[entity.PbehaviorInfo.TypeID]++
				}

				break
			}

			counters.NotAcknowledged++
			counters.Acknowledged--
		}

		output, err := s.getServiceOutput(counters)
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

	if len(entity.ImpactedServicesToAdd) > 0 || len(entity.ImpactedServicesToRemove) > 0 {
		_, err = s.entityCollection.UpdateOne(
			ctx, bson.M{"_id": entity.ID},
			bson.M{"$unset": bson.M{"impacted_services_to_add": 1, "impacted_services_to_remove": 1}})
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
		PbehaviorCounters: make(map[string]int64),
	}

	if len(event.Entity.Depends) == 0 {
		output, err := s.getServiceOutput(counters)
		if err != nil {
			return nil, err
		}

		updatedServiceStates[event.Entity.ID] = UpdatedServicesInfo{
			State:  counters.GetWorstState(),
			Output: output,
		}

		_, err = s.serviceCountersCollection.UpdateOne(ctx, bson.M{"_id": event.GetEID()}, bson.M{"$set": counters}, options.Update().SetUpsert(true))
		return updatedServiceStates, err
	}

	cursor, err := s.entityCollection.Aggregate(ctx, []bson.M{
		{
			"$match": bson.M{"_id": bson.M{"$in": event.Entity.Depends}},
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

		if depEnt.Alarm.ID != "" {
			if depEnt.Alarm.IsResolved() {
				continue
			}

			counters.All++
			if depEnt.Alarm.IsInActivePeriod() {
				counters.IncrementAlarmCounters(int(depEnt.Alarm.CurrentState()), depEnt.Alarm.IsAck())
			} else {
				counters.PbehaviorCounters[depEnt.Alarm.Value.PbehaviorInfo.TypeID]++
			}
		} else {
			if !depEnt.Entity.PbehaviorInfo.IsActive() {
				counters.PbehaviorCounters[depEnt.Entity.PbehaviorInfo.TypeID]++
			}
		}
	}

	_, err = s.serviceCountersCollection.UpdateOne(ctx, bson.M{"_id": event.GetEID()}, bson.M{"$set": counters}, options.Update().SetUpsert(true))
	if err != nil {
		return nil, err
	}

	output, err := s.getServiceOutput(counters)
	if err != nil {
		return nil, err
	}

	updatedServiceStates[event.Entity.ID] = UpdatedServicesInfo{
		State:  counters.GetWorstState(),
		Output: output,
	}

	return updatedServiceStates, nil
}

func (s *service) getServiceOutput(counters EntityServiceCounters) (string, error) {
	tpl, err := template.New("template").Parse(counters.OutputTemplate)
	if err != nil {
		return "", fmt.Errorf(
			"unable to parse output template for service %s: %w", counters.OutputTemplate, err)
	}

	b := strings.Builder{}
	err = tpl.Execute(&b, counters)
	if err != nil {
		return "", fmt.Errorf(
			"unable to execute output template for service %s: %w",
			counters.OutputTemplate, err)
	}

	return b.String(), nil
}
