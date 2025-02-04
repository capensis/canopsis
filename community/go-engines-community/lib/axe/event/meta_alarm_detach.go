package event

import (
	"context"
	"errors"
	"fmt"

	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
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

	_, err := p.detachChildrenFromMetaAlarm(ctx, event)
	if err != nil {
		return result, err
	}

	result.Forward = false

	return result, nil
}

func (p *metaAlarmDetachProcessor) detachChildrenFromMetaAlarm(ctx context.Context, event rpc.AxeEvent) (*types.Alarm, error) {
	if len(event.Parameters.MetaAlarmChildren) == 0 {
		return nil, nil
	}

	var updatedChildrenAlarms []types.Alarm
	var metaAlarm types.Alarm
	var err error

	err = p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		updatedChildrenAlarms = make([]types.Alarm, 0)

		err = p.alarmCollection.FindOne(ctx, bson.M{"d": event.Entity.ID, "v.resolved": nil}).Decode(&metaAlarm)
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

		for _, childAlarm := range alarms {
			if childAlarm.Alarm.RemoveParent(event.Entity.ID) {
				metaAlarm.RemoveChild(childAlarm.Entity.ID)
				eventsCount -= childAlarm.Alarm.Value.EventsCount
				updatedChildrenAlarms = append(updatedChildrenAlarms, childAlarm.Alarm)
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

		var lastEventDate types.CpsTime // should be empty
		worstState := types.CpsNumber(types.AlarmStateOK)

		for _, childAlarm := range metaAlarmChildren {
			if childAlarm.Alarm.Value.State.Value > worstState {
				worstState = childAlarm.Alarm.Value.State.Value
			}

			if lastEventDate.Before(childAlarm.Alarm.Value.LastEventDate) {
				lastEventDate = childAlarm.Alarm.Value.LastEventDate
			}
		}

		metaAlarm.PartialUpdateLastEventDate(lastEventDate)
		metaAlarm.IncrementEventsCount(eventsCount)
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

		metaAlarm.UpdateOutput(output)
		err = updateMetaAlarmState(&metaAlarm, *event.Entity, event.Parameters.Timestamp, worstState, output, p.alarmStatusService)
		if err != nil {
			return err
		}

		return p.adapter.PartialMassUpdateOpen(ctx, append([]types.Alarm{metaAlarm}, updatedChildrenAlarms...))
	})
	if err != nil {
		return nil, err
	}

	return &metaAlarm, nil
}
