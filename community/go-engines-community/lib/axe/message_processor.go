package axe

import (
	"context"
	"runtime/trace"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmtag"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
)

type MessageProcessor struct {
	FeaturePrintEventOnError bool

	EventProcessor         alarm.EventProcessor
	TechMetricsSender      techmetrics.Sender
	RemediationRpcClient   engine.RPCClient
	TimezoneConfigProvider config.TimezoneConfigProvider
	Encoder                encoding.Encoder
	Decoder                encoding.Decoder
	Logger                 zerolog.Logger
	PbehaviorAdapter       pbehavior.Adapter
	TagUpdater             alarmtag.Updater
	AutoInstructionMatcher AutoInstructionMatcher
	AlarmCollection        mongo.DbCollection
}

func (p *MessageProcessor) Process(parentCtx context.Context, d amqp.Delivery) ([]byte, error) {
	eventMetric := techmetrics.AxeEventMetric{}
	eventMetric.Timestamp = time.Now()

	ctx, task := trace.NewTask(parentCtx, "axe.WorkerProcess")
	defer task.End()

	msg := d.Body
	trace.Logf(ctx, "event_size", "%d", len(msg))
	var err error
	var event types.Event
	trace.WithRegion(ctx, "decode-event", func() {
		err = p.Decoder.Decode(msg, &event)
	})
	if err != nil {
		p.logError(err, "cannot decode event", "", msg)
		return nil, nil
	}

	trace.Log(ctx, "event.event_type", event.EventType)
	trace.Log(ctx, "event.timestamp", event.Timestamp.String())
	trace.Log(ctx, "event.source_type", event.SourceType)
	trace.Log(ctx, "event.connector", event.Connector)
	trace.Log(ctx, "event.connector_name", event.ConnectorName)
	trace.Log(ctx, "event.component", event.Component)
	trace.Log(ctx, "event.resource", event.Resource)

	defer func() {
		eventMetric.EventType = event.EventType
		if event.AlarmChange != nil {
			eventMetric.AlarmChangeType = string(event.AlarmChange.Type)
		}
		if event.Entity != nil {
			eventMetric.EntityType = event.Entity.Type
		}

		eventMetric.Interval = time.Since(eventMetric.Timestamp)
		p.TechMetricsSender.SendAxeEvent(eventMetric)
	}()

	alarmChange, err := p.EventProcessor.Process(ctx, &event)

	alarmID := ""
	if event.Alarm != nil {
		alarmID = event.Alarm.ID
	}

	if err != nil {
		if engine.IsConnectionError(err) {
			return nil, err
		}

		p.logError(err, "cannot process event", alarmID, msg)
		return nil, nil
	}
	event.AlarmChange = &alarmChange

	if event.Healthcheck {
		_, err := p.AlarmCollection.DeleteMany(ctx, bson.M{"healthcheck": true})
		if err != nil {
			p.logError(err, "cannot delete temporary alarm", alarmID, d.Body)
		}
	} else {
		err = p.handleRemediation(ctx, event, msg)
		if err != nil {
			return nil, err
		}

		p.updatePbhLastAlarmDate(ctx, event)
		p.updateTags(event)
		event.IsInstructionMatched = p.isInstructionMatched(event, msg)
	}

	// Encode and publish the event to the next engine
	var bevent []byte
	trace.WithRegion(ctx, "encode-event", func() {
		bevent, err = p.Encoder.Encode(event)
	})
	if err != nil {
		p.logError(err, "cannot encode event", alarmID, msg)
		return nil, nil
	}

	return bevent, nil
}

// updatePbhLastAlarmDate updates last time in pbehavior when it was applied on alarm.
func (p *MessageProcessor) updatePbhLastAlarmDate(ctx context.Context, event types.Event) {
	if event.AlarmChange.Type != types.AlarmChangeTypeCreateAndPbhEnter &&
		event.AlarmChange.Type != types.AlarmChangeTypePbhEnter &&
		event.AlarmChange.Type != types.AlarmChangeTypePbhLeaveAndEnter {
		return
	}

	go func() {
		err := p.PbehaviorAdapter.UpdateLastAlarmDate(ctx, event.PbehaviorInfo.ID, types.CpsTime{Time: time.Now()})
		if err != nil {
			p.Logger.Err(err).Msg("cannot update pbehavior")
		}
	}()
}

func (p *MessageProcessor) handleRemediation(ctx context.Context, event types.Event, msg []byte) error {
	if p.RemediationRpcClient == nil || event.Alarm == nil || event.Entity == nil || event.AlarmChange == nil {
		return nil
	}

	switch event.AlarmChange.Type {
	case types.AlarmChangeTypeCreate,
		types.AlarmChangeTypeCreateAndPbhEnter,
		types.AlarmChangeTypeStateIncrease,
		types.AlarmChangeTypeStateDecrease,
		types.AlarmChangeTypeChangeState,
		types.AlarmChangeTypeUnsnooze,
		types.AlarmChangeTypeActivate,
		types.AlarmChangeTypePbhEnter,
		types.AlarmChangeTypePbhLeave,
		types.AlarmChangeTypePbhLeaveAndEnter,
		types.AlarmChangeTypeResolve:
	default:
		return nil
	}

	alarmID := ""
	if event.Alarm != nil {
		alarmID = event.Alarm.ID
	}

	body, err := p.Encoder.Encode(types.RPCRemediationEvent{
		Alarm:       event.Alarm,
		Entity:      event.Entity,
		AlarmChange: *event.AlarmChange,
	})
	if err != nil {
		p.logError(err, "cannot encode remediation event", alarmID, msg)
		return nil
	}

	err = p.RemediationRpcClient.Call(ctx, engine.RPCMessage{
		CorrelationID: event.Alarm.ID,
		Body:          body,
	})
	if err != nil {
		if engine.IsConnectionError(err) {
			return err
		}

		p.logError(err, "cannot send rpc call to remediation", alarmID, msg)
	}

	return nil
}

func (p *MessageProcessor) updateTags(event types.Event) {
	if event.EventType == types.EventTypeCheck {
		p.TagUpdater.Add(event.Tags)
	}
}

func (p *MessageProcessor) isInstructionMatched(event types.Event, msg []byte) (res bool) {
	if event.Alarm == nil || event.Entity == nil || event.AlarmChange == nil {
		return false
	}
	triggers := types.GetTriggers(event.AlarmChange.Type)
	if len(triggers) == 0 {
		return false
	}

	matched, err := p.AutoInstructionMatcher.Match(triggers, types.AlarmWithEntity{Alarm: *event.Alarm, Entity: *event.Entity})
	if err != nil {
		p.logError(err, "cannot match auto instructions", event.Alarm.ID, msg)
		return false
	}

	return matched
}

func (p *MessageProcessor) logError(err error, errMsg string, alarmID string, msg []byte) {
	if p.FeaturePrintEventOnError {
		p.Logger.Err(err).Str("event", string(msg)).Str("alarm_id", alarmID).Msg(errMsg)
	} else {
		p.Logger.Err(err).Str("alarm_id", alarmID).Msg(errMsg)
	}
}
