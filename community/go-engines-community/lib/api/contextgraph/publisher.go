package contextgraph

import (
	"fmt"
	"time"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/streadway/amqp"
)

type rmqPublisher struct {
	encoder       encoding.Encoder
	amqpPublisher libamqp.Publisher
}

func NewRMQPublisher(
	encoder encoding.Encoder,
	amqpPublisher libamqp.Publisher,
) EventPublisher {
	return &rmqPublisher{
		encoder:       encoder,
		amqpPublisher: amqpPublisher,
	}
}

func (p *rmqPublisher) SendImportResultEvent(uuid string, execTime time.Duration, state types.CpsNumber) error {
	output := fmt.Sprintf("Import %s succeed.", uuid)

	if state != types.AlarmStateOK {
		output = fmt.Sprintf("Import %s failed.", uuid)
	}

	stateType := new(types.CpsNumber)
	*stateType = 1

	return p.sendEvent(types.Event{
		Connector:     "taskhandler",
		ConnectorName: "task_importctx",
		Component:     "job",
		EventType:     types.EventTypeCheck,
		SourceType:    types.SourceTypeResource,
		Resource:      uuid,
		Timestamp:     types.NewCpsTime(time.Now().Unix()),
		State:         state,
		StateType:     stateType,
		Output:        output,
		ExtraInfos: map[string]interface{}{
			"execution_time": execTime,
		},
	})
}

func (p *rmqPublisher) SendPerfDataEvent(uuid string, stats JobStats, state types.CpsNumber) error {
	stateType := new(types.CpsNumber)
	*stateType = 1

	return p.sendEvent(types.Event{
		Connector:     "Taskhandler",
		ConnectorName: "task_importctx",
		Component:     uuid,
		EventType:     types.EventTypeCheck,
		SourceType:    types.SourceTypeResource,
		Resource:      "task_importctx/report",
		State:         state,
		StateType:     stateType,
		Output:        fmt.Sprintf("execution : %f sec, updated ent : %d, deleted ent : %d", stats.ExecTime.Seconds(), stats.Updated, stats.Deleted),
		PerfDataArray: []types.PerfData{
			{
				Metric: "execution_time",
				Unit:   "GAUGE",
				Value:  stats.ExecTime.Seconds(),
			},
			{
				Metric: "ent_updated",
				Unit:   "GAUGE",
				Value:  float64(stats.Updated),
			},
			{
				Metric: "ent_deleted",
				Unit:   "GAUGE",
				Value:  float64(stats.Deleted),
			},
		},
	})
}

func (p *rmqPublisher) SendUpdateEntityServiceEvent(serviceId string) error {
	return p.sendEvent(types.Event{
		EventType:     types.EventTypeRecomputeEntityService,
		Connector:     types.ConnectorEngineService,
		ConnectorName: types.ConnectorEngineService,
		Component:     serviceId,
		Timestamp:     types.CpsTime{Time: time.Now()},
		Author:        canopsis.DefaultEventAuthor,
		SourceType:    types.SourceTypeService,
	})
}

func (p *rmqPublisher) sendEvent(event types.Event) error {
	bevent, err := p.encoder.Encode(event)
	if err != nil {
		return fmt.Errorf("error while encoding event %+v", err)
	}

	return p.amqpPublisher.Publish(
		"canopsis.events",
		"",
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json", // this type is mandatory to avoid bad conversions into Python.
			Body:         bevent,
			DeliveryMode: amqp.Persistent,
		},
	)
}
