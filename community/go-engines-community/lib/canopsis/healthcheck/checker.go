package healthcheck

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	eventComponent = "healthcheck"
	eventAuthor    = "healthcheck"
	eventState     = types.AlarmStateCritical
)

type Checker interface {
	Check(ctx context.Context) error
}

func NewChecker(
	engine string,
	processor engine.MessageProcessor,
	encoder encoding.Encoder,
	eventWithEntity bool,
	eventWithAlarm bool,
) Checker {
	return &checker{
		processor:       processor,
		encoder:         encoder,
		engine:          engine,
		eventWithEntity: eventWithEntity,
		eventWithAlarm:  eventWithAlarm,
	}
}

type checker struct {
	processor       engine.MessageProcessor
	encoder         encoding.Encoder
	engine          string
	eventWithEntity bool
	eventWithAlarm  bool
}

func (c *checker) Check(ctx context.Context) error {
	event := c.createEvent()
	b, err := c.encoder.Encode(event)
	if err != nil {
		return err
	}
	_, err = c.processor.Process(ctx, amqp.Delivery{
		Body: b,
	})
	return err
}

func (c *checker) createEvent() types.Event {
	now := datetime.NewCpsTime()
	event := types.Event{
		EventType:     types.EventTypeCheck,
		State:         eventState,
		Connector:     c.engine,
		ConnectorName: c.engine,
		Component:     eventComponent,
		Resource:      utils.NewID(),
		SourceType:    types.SourceTypeResource,
		Output:        "check engine-" + c.engine,
		Author:        eventAuthor,
		Timestamp:     now,
		Initiator:     types.InitiatorSystem,
		Healthcheck:   true,
	}
	if c.eventWithEntity {
		entity := types.Entity{
			ID:            event.GetEID(),
			Name:          event.Resource,
			Enabled:       true,
			Type:          types.EntityTypeResource,
			Created:       now,
			LastEventDate: &now,
			Connector:     event.Connector + "/" + event.ConnectorName,
			Component:     event.Component,
			Healthcheck:   true,
		}
		event.Entity = &entity
	}
	if c.eventWithAlarm {
		output := "healthcheck" + event.Resource
		stateStep := types.AlarmStep{
			Type:      types.AlarmStepStateIncrease,
			Timestamp: now,
			Author:    event.Author,
			Message:   output,
			Value:     event.State,
			Initiator: event.Initiator,
		}
		statusStep := types.AlarmStep{
			Type:      types.AlarmStepStatusIncrease,
			Timestamp: now,
			Author:    event.Author,
			Message:   output,
			Value:     types.AlarmStatusOngoing,
			Initiator: event.Initiator,
		}
		alarm := types.Alarm{
			Time:     now,
			EntityID: event.GetEID(),
			Value: types.AlarmValue{
				State:  &stateStep,
				Status: &statusStep,
				Steps: []types.AlarmStep{
					stateStep,
					statusStep,
				},
				Component:         event.Component,
				Connector:         event.Connector,
				ConnectorName:     event.ConnectorName,
				CreationDate:      now,
				DisplayName:       output,
				InitialOutput:     output,
				Output:            output,
				InitialLongOutput: output,
				LongOutput:        output,
				LastUpdateDate:    now,
				LastEventDate:     now,
				Resource:          event.Resource,
			},
			Healthcheck: true,
		}
		event.Alarm = &alarm
		event.AlarmChange = &types.AlarmChange{
			Type: types.AlarmChangeTypeNone,
		}
	}

	return event
}
