package contextgraph

import (
	"fmt"
	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/importcontextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/streadway/amqp"
	"time"
)

type rmqPublisher struct {
	exchange, queue string
	encoder         encoding.Encoder
	contentType     string
	amqpPublisher   libamqp.Publisher
}

func NewEventPublisher(
	exchange, queue string,
	encoder encoding.Encoder,
	contentType string,
	amqpPublisher libamqp.Publisher,
) EventPublisher {
	return &rmqPublisher{
		exchange:      exchange,
		queue:         queue,
		encoder:       encoder,
		contentType:   contentType,
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
		ExecutionTime: execTime,
	})
}

func (p *rmqPublisher) SendPerfDataEvent(uuid string, stats importcontextgraph.Stats, state types.CpsNumber) error {
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

func (p *rmqPublisher) sendEvent(event types.Event) error {
	bevent, err := p.encoder.Encode(event)
	if err != nil {
		return fmt.Errorf("error while encoding event %+v", err)
	}

	return p.amqpPublisher.Publish(
		p.exchange,
		p.queue,
		false,
		false,
		amqp.Publishing{
			ContentType:  p.contentType,
			Body:         bevent,
			DeliveryMode: amqp.Persistent,
		},
	)
}
