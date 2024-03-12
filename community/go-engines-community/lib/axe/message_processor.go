package axe

import (
	"context"
	"runtime/trace"
	"time"

	libevent "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/axe/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
)

type MessageProcessor struct {
	FeaturePrintEventOnError bool

	EventProcessor    libevent.Processor
	Encoder           encoding.Encoder
	Decoder           encoding.Decoder
	TechMetricsSender techmetrics.Sender
	AlarmCollection   mongo.DbCollection
	Logger            zerolog.Logger
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

	res, err := p.EventProcessor.Process(ctx, p.transformEvent(event))
	if err != nil {
		if engine.IsConnectionError(err) {
			return nil, err
		}

		p.logError(err, "cannot process event", "", msg)
		return nil, nil
	}

	if event.Healthcheck {
		_, err := p.AlarmCollection.DeleteMany(ctx, bson.M{"healthcheck": true})
		if err != nil {
			p.logError(err, "cannot delete temporary alarm", "", d.Body)
		}
	}

	if !res.Forward {
		return nil, nil
	}

	if res.Alarm.ID != "" {
		event.Alarm = &res.Alarm
	}

	if res.Entity.ID != "" {
		event.Entity = &res.Entity
	}

	if !res.AlarmChange.IsZero() {
		event.AlarmChange = &res.AlarmChange
	}

	event.IsInstructionMatched = res.IsInstructionMatched
	// Encode and publish the event to the next engine
	var bevent []byte
	trace.WithRegion(ctx, "encode-event", func() {
		bevent, err = p.Encoder.Encode(event)
	})
	if err != nil {
		alarmID := ""
		if event.Alarm != nil {
			alarmID = event.Alarm.ID
		}
		p.logError(err, "cannot encode event", alarmID, msg)
		return nil, nil
	}

	return bevent, nil
}

func (p *MessageProcessor) transformEvent(event types.Event) rpc.AxeEvent {
	params := rpc.AxeParameters{
		Output:              event.Output,
		Author:              event.Author,
		User:                event.UserID,
		Role:                event.Role,
		Initiator:           event.Initiator,
		Timestamp:           event.Timestamp,
		State:               &event.State,
		TicketInfo:          event.TicketInfo,
		PbehaviorInfo:       event.PbehaviorInfo,
		Execution:           event.Execution,
		Instruction:         event.Instruction,
		LongOutput:          event.LongOutput,
		Connector:           event.Connector,
		ConnectorName:       event.ConnectorName,
		Tags:                event.Tags,
		IdleRuleApply:       event.IdleRuleApply,
		MetaAlarmRuleID:     event.MetaAlarmRuleID,
		MetaAlarmValuePath:  event.MetaAlarmValuePath,
		DisplayName:         event.DisplayName,
		MetaAlarmChildren:   event.MetaAlarmChildren,
		StateSettingUpdated: event.StateSettingUpdated,
	}

	if event.Duration > 0 {
		params.Duration = &datetime.DurationWithUnit{
			Value: int64(event.Duration),
			Unit:  datetime.DurationUnitSecond,
		}
	}

	if event.AlarmChange != nil {
		params.Trigger = string(event.AlarmChange.Type)
	}

	return rpc.AxeEvent{
		EventType:   event.EventType,
		Parameters:  params,
		Alarm:       event.Alarm,
		AlarmID:     event.AlarmID,
		Entity:      event.Entity,
		Healthcheck: event.Healthcheck,
	}
}

func (p *MessageProcessor) logError(err error, errMsg string, alarmID string, msg []byte) {
	if p.FeaturePrintEventOnError {
		p.Logger.Err(err).Str("event", string(msg)).Str("alarm_id", alarmID).Msg(errMsg)
	} else {
		p.Logger.Err(err).Str("alarm_id", alarmID).Msg(errMsg)
	}
}
