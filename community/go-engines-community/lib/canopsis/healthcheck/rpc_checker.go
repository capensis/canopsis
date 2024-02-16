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

func NewRpcChecker(
	engine string,
	processor engine.MessageProcessor,
	encoder encoding.Encoder,
	createEvent func(types.Entity, types.Alarm) any,
) Checker {
	return &rpcChecker{
		processor:   processor,
		encoder:     encoder,
		engine:      engine,
		createEvent: createEvent,
	}
}

type rpcChecker struct {
	processor   engine.MessageProcessor
	encoder     encoding.Encoder
	engine      string
	createEvent func(types.Entity, types.Alarm) any
}

func (c *rpcChecker) Check(ctx context.Context) error {
	uuid := utils.NewID()
	now := datetime.NewCpsTime()
	event := c.createEvent(c.createEntity(uuid, now), c.createAlarm(uuid, now))
	b, err := c.encoder.Encode(event)
	if err != nil {
		return err
	}
	_, err = c.processor.Process(ctx, amqp.Delivery{
		Body: b,
	})
	return err
}

func (c *rpcChecker) createEntity(resource string, now datetime.CpsTime) types.Entity {
	return types.Entity{
		ID:            resource + "/" + eventComponent,
		Name:          resource,
		Enabled:       true,
		Type:          types.EntityTypeResource,
		Created:       now,
		LastEventDate: &now,
		Connector:     c.engine + "/" + c.engine,
		Component:     eventComponent,
		Healthcheck:   true,
	}
}

func (c *rpcChecker) createAlarm(resource string, now datetime.CpsTime) types.Alarm {
	output := "healthcheck" + resource
	stateStep := types.AlarmStep{
		Type:      types.AlarmStepStateIncrease,
		Timestamp: now,
		Author:    eventAuthor,
		Message:   output,
		Value:     eventState,
		Initiator: types.InitiatorSystem,
	}
	statusStep := types.AlarmStep{
		Type:      types.AlarmStepStatusIncrease,
		Timestamp: now,
		Author:    eventAuthor,
		Message:   output,
		Value:     types.AlarmStatusOngoing,
		Initiator: types.InitiatorSystem,
	}
	return types.Alarm{
		Time:     now,
		EntityID: resource + "/" + eventComponent,
		Value: types.AlarmValue{
			State:  &stateStep,
			Status: &statusStep,
			Steps: []types.AlarmStep{
				stateStep,
				statusStep,
			},
			Component:         eventComponent,
			Connector:         c.engine,
			ConnectorName:     c.engine,
			CreationDate:      now,
			DisplayName:       output,
			InitialOutput:     output,
			Output:            output,
			InitialLongOutput: output,
			LongOutput:        output,
			LastUpdateDate:    now,
			LastEventDate:     now,
			Resource:          resource,
		},
		Healthcheck: true,
	}
}
