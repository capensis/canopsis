package fifo

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/scheduler"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type ackMessageProcessor struct {
	FeaturePrintEventOnError bool
	Scheduler                scheduler.Scheduler
	Decoder                  encoding.Decoder
	Logger                   zerolog.Logger
}

func (p *ackMessageProcessor) Process(ctx context.Context, d amqp.Delivery) ([]byte, error) {
	msg := d.Body
	var event types.Event
	err := p.Decoder.Decode(msg, &event)
	if err != nil {
		p.logError(err, "cannot decode event", msg)
		return nil, nil
	}

	p.Logger.Debug().Msgf("valid input event: %v", string(msg))
	err = p.Scheduler.AckEvent(ctx, event)
	if err != nil {
		p.logError(err, "cannot process event", msg)
		return nil, nil
	}

	return nil, nil
}

func (p *ackMessageProcessor) logError(err error, errMsg string, msg []byte) {
	if p.FeaturePrintEventOnError {
		p.Logger.Err(err).Str("event", string(msg)).Msg(errMsg)
	} else {
		p.Logger.Err(err).Msg(errMsg)
	}
}
