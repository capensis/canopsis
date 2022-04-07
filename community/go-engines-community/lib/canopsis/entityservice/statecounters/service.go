package statecounters

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	libamqp "github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"strings"
	"text/template"
	"time"
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
		return fmt.Errorf("unable to serialize service event: %v", err)
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
		return fmt.Errorf("unable to send service event: %v", err)
	}

	return nil
}

func (s *service) UpdateServiceCounters(ctx context.Context, entity types.Entity, alarm *types.Alarm, alarmChange types.AlarmChange) (map[string]UpdatedServicesInfo, error) {
	switch alarmChange.Type {
	case types.AlarmChangeTypeCreate, types.AlarmChangeTypeCreateAndPbhEnter,
		types.AlarmChangeTypePbhEnter, types.AlarmChangeTypePbhLeave, types.AlarmChangeTypePbhLeaveAndEnter,
		types.AlarmChangeTypeChangeState, types.AlarmChangeTypeStateDecrease, types.AlarmChangeTypeStateIncrease,
		types.AlarmChangeTypeResolve, types.AlarmChangeTypeAck, types.AlarmChangeTypeAckremove, types.AlarmChangeTypeNone:
	default:
		return nil, nil
	}

	prevState := int(alarmChange.PreviousState)
	curState := types.AlarmStateOK
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

	inSlice := append(entity.ImpactedServices, entity.ImpactedServicesToRemove...)
	if len(inSlice) == 0 {
		return updatedServiceInfos, nil
	}

	cursor, err = s.serviceCountersCollection.Find(context.Background(), bson.M{"_id": bson.M{"$in": inSlice}})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

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
		case types.AlarmChangeTypeNone:
			if serviceToRemove[counters.ID] {
				if alarm != nil {
					counters.All--
					if alarm.IsInActivePeriod() {
						counters.DecrementAlarmCounters(curState, acked)
					} else {
						counters.PbehaviorCounters[alarm.Value.PbehaviorInfo.TypeID]--
					}
				} else if !entity.PbehaviorInfo.IsActive() {
					counters.PbehaviorCounters[entity.PbehaviorInfo.TypeID]--
				}
			} else if serviceToAdd[counters.ID] {
				if alarm != nil {
					counters.All++
					if alarm.IsInActivePeriod() {
						counters.IncrementAlarmCounters(curState, acked)
					} else {
						counters.PbehaviorCounters[alarm.Value.PbehaviorInfo.TypeID]++
					}
				} else if !entity.PbehaviorInfo.IsActive() {
					counters.PbehaviorCounters[entity.PbehaviorInfo.TypeID]++
				}
			}
		case types.AlarmChangeTypeCreate, types.AlarmChangeTypeCreateAndPbhEnter:
			if serviceToRemove[counters.ID] {
				continue
			} else {
				counters.All++
				if alarm.IsInActivePeriod() {
					counters.IncrementAlarmCounters(curState, acked)
				} else {
					if alarmChange.PreviousPbehaviorTypeID != alarm.Value.PbehaviorInfo.TypeID {
						counters.PbehaviorCounters[alarm.Value.PbehaviorInfo.TypeID]++
					}
				}
			}
		case types.AlarmChangeTypePbhEnter:
			if serviceToRemove[counters.ID] && alarm != nil {
				counters.All--
				counters.DecrementAlarmCounters(curState, acked)
			} else if serviceToAdd[counters.ID] {
				if alarm != nil {
					counters.All++
					if alarm.IsInActivePeriod() {
						counters.IncrementAlarmCounters(curState, acked)
					} else {
						counters.PbehaviorCounters[alarm.Value.PbehaviorInfo.TypeID]++
					}
				} else if !entity.PbehaviorInfo.IsActive() {
					counters.PbehaviorCounters[entity.PbehaviorInfo.TypeID]++
				}
			} else {
				if alarm != nil && !alarm.IsInActivePeriod() {
					counters.DecrementAlarmCounters(curState, acked)
					counters.PbehaviorCounters[alarm.Value.PbehaviorInfo.TypeID]++
				} else if !entity.PbehaviorInfo.IsActive() {
					counters.PbehaviorCounters[entity.PbehaviorInfo.TypeID]++
				}
			}
		case types.AlarmChangeTypePbhLeaveAndEnter:
			if serviceToRemove[counters.ID] {
				if alarm != nil {
					counters.All--
					if alarmChange.PreviousPbehaviorCannonicalType != pbehavior.TypeActive {
						counters.DecrementAlarmCounters(curState, acked)
					} else {
						counters.PbehaviorCounters[alarmChange.PreviousPbehaviorTypeID]--
					}
				} else if alarmChange.PreviousPbehaviorCannonicalType != pbehavior.TypeActive {
					counters.PbehaviorCounters[alarmChange.PreviousPbehaviorTypeID]--
				}
			} else if serviceToAdd[counters.ID] {
				if alarm != nil {
					counters.All++
					if alarm.IsInActivePeriod() {
						counters.IncrementAlarmCounters(curState, acked)
					} else {
						counters.PbehaviorCounters[alarm.Value.PbehaviorInfo.TypeID]++
					}
				} else if !entity.PbehaviorInfo.IsActive() {
					counters.PbehaviorCounters[entity.PbehaviorInfo.TypeID]++
				}
			} else {
				if alarmChange.PreviousPbehaviorCannonicalType != pbehavior.TypeActive {
					counters.PbehaviorCounters[alarmChange.PreviousPbehaviorTypeID]--
				}

				if alarm != nil {
					if alarm.IsInActivePeriod() {
						if alarmChange.PreviousPbehaviorCannonicalType != pbehavior.TypeActive {
							counters.IncrementAlarmCounters(curState, acked)
						}
					} else {
						counters.PbehaviorCounters[alarm.Value.PbehaviorInfo.TypeID]++
						if alarmChange.PreviousPbehaviorCannonicalType == pbehavior.TypeActive {
							counters.DecrementAlarmCounters(curState, acked)
						}
					}
				} else {
					if !entity.PbehaviorInfo.IsActive() {
						counters.PbehaviorCounters[entity.PbehaviorInfo.TypeID]++
					}
				}
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
			} else if serviceToAdd[counters.ID] {
				if alarm != nil {
					counters.All++
					counters.IncrementAlarmCounters(curState, acked)
				}
			} else {
				if alarmChange.PreviousPbehaviorCannonicalType != pbehavior.TypeActive {
					if alarm != nil {
						counters.IncrementAlarmCounters(curState, acked)
					}

					counters.PbehaviorCounters[alarmChange.PreviousPbehaviorTypeID]--
				}
			}
		case types.AlarmChangeTypeStateIncrease,
			types.AlarmChangeTypeStateDecrease,
			types.AlarmChangeTypeChangeState:
			if serviceToRemove[counters.ID] {
				counters.All--
				if alarm.IsInActivePeriod() {
					counters.DecrementAlarmCounters(prevState, acked)
				}
			} else if serviceToAdd[counters.ID] {
				counters.All++
				if alarm.IsInActivePeriod() {
					counters.IncrementAlarmCounters(curState, acked)
				} else {
					counters.PbehaviorCounters[alarm.Value.PbehaviorInfo.TypeID]++
				}
			} else {
				if alarm.IsInActivePeriod() {
					counters.DecrementState(prevState)
					counters.IncrementState(curState)
				}
			}
		case types.AlarmChangeTypeResolve:
			if serviceToAdd[counters.ID] {
				continue
			} else {
				counters.All--
				if alarm.IsInActivePeriod() {
					counters.DecrementAlarmCounters(curState, acked)
				} else if serviceToRemove[counters.ID] {
					counters.PbehaviorCounters[entity.PbehaviorInfo.TypeID]--
				}
			}
		case types.AlarmChangeTypeAck:
			if serviceToRemove[counters.ID] {
				if alarm.IsInActivePeriod() {
					counters.DecrementState(curState)
					counters.NotAcknowledged--
				} else {
					counters.PbehaviorCounters[alarm.Value.PbehaviorInfo.TypeID]--
				}
			} else if serviceToAdd[counters.ID] {
				if alarm.IsInActivePeriod() {
					counters.IncrementState(curState)
					counters.Acknowledged++
				} else {
					counters.PbehaviorCounters[alarm.Value.PbehaviorInfo.TypeID]++
				}
			} else {
				counters.Acknowledged++
				counters.NotAcknowledged--
			}
		case types.AlarmChangeTypeAckremove:
			if serviceToRemove[counters.ID] {
				if alarm.IsInActivePeriod() {
					counters.DecrementState(curState)
					counters.Acknowledged--
				} else {
					counters.PbehaviorCounters[alarm.Value.PbehaviorInfo.TypeID]--
				}
			} else if serviceToAdd[counters.ID] {
				if alarm.IsInActivePeriod() {
					counters.IncrementState(curState)
					counters.NotAcknowledged++
				} else {
					counters.PbehaviorCounters[alarm.Value.PbehaviorInfo.TypeID]++
				}
			} else {
				counters.NotAcknowledged++
				counters.Acknowledged--
			}
		}

		output, err := s.getServiceOutput(counters)
		if err != nil {
			return nil, err
		}

		updatedServiceInfos[counters.ID] = UpdatedServicesInfo{
			State:  counters.GetWorstState(),
			Output: output,
		}

		updateCountersModels = append(updateCountersModels,
			mongodriver.
				NewUpdateOneModel().
				SetFilter(bson.M{"_id": counters.ID}).
				SetUpdate(bson.M{"$set": counters}).
				SetUpsert(true),
		)

		if len(updateCountersModels) >= canopsis.DefaultBulkSize {
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
		ID:             event.Entity.ID,
		OutputTemplate: counters.OutputTemplate,
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
				"from":         mongo.AlarmMongoCollection,
				"localField":   "entity._id",
				"foreignField": "d",
				"as":           "alarm",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$alarm",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$match": bson.M{"alarm.v.resolved": nil},
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
