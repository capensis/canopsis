package event

import (
	"context"
	"errors"
	"fmt"

	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

func NewMetaAlarmDetachProcessor(
	dbClient mongo.DbClient,
	ruleAdapter correlation.RulesAdapter,
	adapter libalarm.Adapter,
	alarmStatusService alarmstatus.Service,
	templateExecutor template.Executor,
) Processor {
	return &metaAlarmDetachProcessor{
		dbClient:           dbClient,
		alarmCollection:    dbClient.Collection(mongo.AlarmMongoCollection),
		adapter:            adapter,
		ruleAdapter:        ruleAdapter,
		alarmStatusService: alarmStatusService,
		templateExecutor:   templateExecutor,
	}
}

type metaAlarmDetachProcessor struct {
	dbClient           mongo.DbClient
	alarmCollection    mongo.DbCollection
	ruleAdapter        correlation.RulesAdapter
	adapter            libalarm.Adapter
	alarmStatusService alarmstatus.Service
	templateExecutor   template.Executor
}

func (p *metaAlarmDetachProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil {
		return result, nil
	}

	err := p.detachChildrenFromMetaAlarm(ctx, event)
	if err != nil {
		return result, err
	}

	result.Forward = false

	return result, nil
}

func (p *metaAlarmDetachProcessor) detachChildrenFromMetaAlarm(ctx context.Context, event rpc.AxeEvent) error {
	if len(event.Parameters.MetaAlarmChildren) == 0 {
		return nil
	}

	var updatedChildrenAlarms []types.Alarm
	var metaAlarm types.Alarm

	return p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		updatedChildrenAlarms = updatedChildrenAlarms[:0]
		metaAlarm = types.Alarm{}

		err := p.alarmCollection.FindOne(ctx, bson.M{"d": event.Entity.ID, "v.resolved": nil}).Decode(&metaAlarm)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		rule, err := p.ruleAdapter.GetRule(ctx, event.Parameters.MetaAlarmRuleID)
		if err != nil {
			return fmt.Errorf("cannot fetch meta alarm rule id=%q: %w", event.Parameters.MetaAlarmRuleID, err)
		} else if rule.ID == "" {
			return fmt.Errorf("meta alarm rule id=%q not found", event.Parameters.MetaAlarmRuleID)
		}

		alarms, err := getAlarmsWithEntityByMatch(ctx, p.alarmCollection, bson.M{
			"v.parents":  metaAlarm.EntityID,
			"d":          bson.M{"$in": event.Parameters.MetaAlarmChildren},
			"v.resolved": nil,
		})
		if err != nil {
			return err
		}

		eventsCount := types.CpsNumber(0)
		writeModels := make([]mongodriver.WriteModel, 0, len(alarms))
		childrenIds := make([]string, 0, len(alarms))

		for _, childAlarm := range alarms {
			if childAlarm.Alarm.RemoveParent(event.Entity.ID) {
				metaAlarm.RemoveChild(childAlarm.Entity.ID)
				childrenIds = append(childrenIds, childAlarm.Entity.ID)
				eventsCount -= childAlarm.Alarm.Value.EventsCount
				newStep := NewAlarmStep(types.AlarmStepMetaAlarmDetach, event.Parameters, !childAlarm.Alarm.Value.PbehaviorInfo.IsDefaultActive())
				newStep.Message = getMetaAlarmChildStepMsg(rule, metaAlarm, event)
				err = childAlarm.Alarm.Value.Steps.Add(newStep)
				if err != nil {
					return err
				}

				updatedChildrenAlarms = append(updatedChildrenAlarms, childAlarm.Alarm)
				writeModels = append(writeModels, mongodriver.NewUpdateOneModel().
					SetFilter(bson.M{
						"_id":       childAlarm.Alarm.ID,
						"v.parents": metaAlarm.EntityID,
					}).
					SetUpdate(bson.M{
						"$pull": bson.M{"v.parents": event.Entity.ID},
						"$push": bson.M{
							"v.unlinked_parents": event.Entity.ID,
							"v.steps":            newStep,
						},
					}))
			}
		}

		if len(updatedChildrenAlarms) == 0 {
			return nil
		}

		metaAlarmChildren, err := getAlarmsWithEntityByMatch(ctx, p.alarmCollection, bson.M{
			"v.parents":  metaAlarm.EntityID,
			"d":          bson.M{"$in": metaAlarm.Value.Children},
			"v.resolved": nil,
		})
		if err != nil {
			return err
		}

		var lastEventDate datetime.CpsTime // should be empty
		worstState := types.CpsNumber(types.AlarmStateOK)

		for _, childAlarm := range metaAlarmChildren {
			if childAlarm.Alarm.Value.State.Value > worstState {
				worstState = childAlarm.Alarm.Value.State.Value
			}

			if lastEventDate.Before(childAlarm.Alarm.Value.LastEventDate) {
				lastEventDate = childAlarm.Alarm.Value.LastEventDate
			}
		}

		infos := correlation.EventExtraInfosMeta{
			Rule:  rule,
			Count: int64(len(metaAlarmChildren)),
		}
		if len(metaAlarmChildren) > 0 {
			infos.Children = metaAlarmChildren[len(metaAlarmChildren)-1]
		}

		output := ""
		if rule.IsManual() {
			output = event.Parameters.Output
		} else {
			output, err = executeMetaAlarmOutputTpl(p.templateExecutor, infos)
			if err != nil {
				return err
			}
		}

		if output == "" {
			output = metaAlarm.Value.Output
		}

		setUpdate, pushUpdate, err := updateMetaAlarmState(&metaAlarm, *event.Entity, event.Parameters.Timestamp,
			worstState, output, p.alarmStatusService)
		if err != nil {
			return err
		}

		if setUpdate == nil {
			setUpdate = bson.M{}
		}

		metaAlarm.Value.Output = output
		setUpdate["v.output"] = output
		metaAlarm.Value.LastEventDate = lastEventDate
		setUpdate["v.last_event_date"] = lastEventDate
		update := bson.M{
			"$set":  setUpdate,
			"$inc":  bson.M{"v.events_count": eventsCount},
			"$pull": bson.M{"v.children": bson.M{"$in": childrenIds}},
		}
		if len(pushUpdate) > 0 {
			update["$push"] = pushUpdate
		}

		writeModels = append(writeModels, mongodriver.NewUpdateOneModel().
			SetFilter(bson.M{"_id": metaAlarm.ID}).
			SetUpdate(update))
		_, err = p.alarmCollection.BulkWrite(ctx, writeModels)

		return err
	})
}
