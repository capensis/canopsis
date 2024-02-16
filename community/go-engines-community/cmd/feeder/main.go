package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/signal"
	"strconv"
	"time"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	cps "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type References struct {
	channelPub *amqp.Channel
}

type Feeder struct {
	flags      Flags
	references References
	logger     zerolog.Logger
}

func (f *Feeder) getDirtyEvent(state, Ci, ci, ri int64) types.Event {
	return types.Event{
		Connector:     "feeder2",
		ConnectorName: "feeder2_inst" + strconv.Itoa(int(Ci)),
		Component:     "feeder2_" + strconv.Itoa(int(ci)),
		Resource:      "feeder2_" + strconv.Itoa(int(ri)),
		State:         types.CpsNumber(state),
		SourceType:    types.SourceTypeResource,
		EventType:     types.EventTypeCheck,
	}
}

func (f *Feeder) getCompatEvent(state, Ci, ci, ri int64) types.Event {
	return types.Event{
		Connector:     "feeder2",
		ConnectorName: "feeder2_inst" + strconv.Itoa(int(Ci)),
		Component:     "feeder2_" + strconv.Itoa(int(ci)),
		Resource:      "feeder2_" + strconv.Itoa(int(ri)),
		State:         types.CpsNumber(state),
		SourceType:    types.SourceTypeResource,
		EventType:     types.EventTypeCheck,
		Timestamp:     datetime.NewCpsTime(),
	}
}

func (f *Feeder) sendBytes(ctx context.Context, content []byte, rk string) error {
	return f.references.channelPub.PublishWithContext(
		ctx,
		f.flags.ExchangeName,
		rk,
		false,
		false,
		amqp.Publishing{
			Body:        content,
			ContentType: "application/json",
		},
	)
}

func (f *Feeder) send(ctx context.Context, state, Ci, ci, ri int64) error {
	var evt types.Event
	if f.flags.DirtyEvent {
		evt = f.getDirtyEvent(state, Ci, ci, ri)
	} else {
		evt = f.getCompatEvent(state, Ci, ci, ri)
	}
	bevt, _ := json.Marshal(evt)
	return f.sendBytes(ctx, bevt, evt.GetCompatRK())
}

func (f *Feeder) adjust(target float64, sent, tsent int64) int64 {
	eps := float64(sent * time.Second.Nanoseconds() / tsent)
	variance := math.Abs(target-eps) * 100.0 / target
	adj := int64(0)

	if variance > 1 {
		if eps < target {
			adj = -tsent / 100
		} else {
			adj = tsent / 100
		}
	}

	return adj
}

func (f *Feeder) setupAmqp() error {
	amqpSession, err := libamqp.NewSession()
	if err != nil {
		return fmt.Errorf("amqp session: %w", err)
	}

	channelPub, err := amqpSession.Channel()
	if err != nil {
		return fmt.Errorf("amqp pub channel: %w", err)
	}

	if err := channelPub.Confirm(false); err != nil {
		return fmt.Errorf("confirm: %w", err)
	}

	f.references = References{
		channelPub: channelPub,
	}

	return nil
}

func (f *Feeder) Run(ctx context.Context) error {
	var err error
	switch f.flags.Mode {
	case "file":
		if f.flags.PubHTTP {
			f.modeSendEventHTTP()
		} else if f.flags.PubAMQP {
			err = f.modeSendEvent(ctx)
		}
	case "feeder":
		err = f.modeFeeder(ctx)
	default:
		err = fmt.Errorf("unknown mode \"%s\": please check help (-h)", f.flags.Mode)
	}

	return err
}

func NewFeeder(logger zerolog.Logger) (*Feeder, error) {
	f := Feeder{
		logger: logger,
	}
	if err := f.flags.ParseArgs(); err != nil {
		return nil, err
	}

	if f.flags.Version {
		cps.PrintVersionInfo()
		os.Exit(0)
	}

	return &f, nil
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	logger := log.NewLogger(false)

	feeder, err := NewFeeder(logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("feeder init error")
	}

	if err = feeder.Run(ctx); err != nil {
		logger.Fatal().Err(err).Msg("feeder run error")
	}
}
