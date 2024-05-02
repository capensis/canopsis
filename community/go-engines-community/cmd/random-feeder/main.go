// nolint
// no lint until #4082 has not fixed
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"time"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

const (
	numberOfConnectors = 10
	numberOfComponents = 1000
)

type resourceData func() (int64, int64, int64)

type References struct {
	channelPub *amqp.Channel
}

type Feeder struct {
	flags      Flags
	references References
	logger     zerolog.Logger

	oldResourcesMap map[int]int
	newResourcesMap map[int]int
}

func (f *Feeder) getEvent(state, Ci, ci, ri int64) types.Event {
	return types.Event{
		Connector:     "feeder" + strconv.Itoa(int(Ci)),
		ConnectorName: "feeder" + strconv.Itoa(int(Ci)) + "_inst0",
		Component:     "feeder" + strconv.Itoa(int(Ci)) + "_" + strconv.Itoa(int(ci)),
		Resource:      "feeder" + strconv.Itoa(int(Ci)) + "_" + strconv.Itoa(int(ri)),
		State:         types.CpsNumber(state),
		SourceType:    types.SourceTypeResource,
		EventType:     types.EventTypeCheck,
	}
}

func (f *Feeder) sendBytes(ctx context.Context, content []byte) error {
	return f.references.channelPub.PublishWithContext(
		ctx,
		f.flags.ExchangeName,
		"Engine_che",
		false,
		false,
		amqp.Publishing{
			Body:        content,
			ContentType: "application/json",
		},
	)
}

func (f *Feeder) send(ctx context.Context, state, Ci, ci, ri int64) error {
	evt := f.getEvent(state, Ci, ci, ri)

	bevt, _ := json.Marshal(evt)

	return f.sendBytes(ctx, bevt)
}

func (f *Feeder) adjust(target float64, sent, tsent int64) int64 {
	eps := float64(sent * time.Second.Nanoseconds() / tsent)
	variance := math.Abs(target-eps) * 100.0 / float64(target)
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
		return fmt.Errorf("amqp session: %v", err)
	}

	channelPub, err := amqpSession.Channel()
	if err != nil {
		return fmt.Errorf("amqp pub channel: %v", err)
	}

	if err := channelPub.Confirm(false); err != nil {
		return fmt.Errorf("confirm: %v", err)
	}

	f.references = References{
		channelPub: channelPub,
	}

	return nil
}

func (f *Feeder) sendMessages(ctx context.Context, eventsPerSecond float64, callback resourceData) error {
	if err := f.setupAmqp(); err != nil {
		return err
	}

	sleepEvery := int64(25)
	pubcount := int64(0)

	nanosecSleep := time.Duration(float64(time.Second.Nanoseconds()) / eventsPerSecond * float64(sleepEvery))

	tstart := time.Now().UnixNano()

	checkEvery := int64(100)
	changeStateEvery := int64(100 / f.flags.Alarms)

	stateMap := make(map[string]int)

	for {
		connectorId, componentId, resourceId := callback()

		eid := fmt.Sprintf("%d%d", componentId, resourceId)
		state := stateMap[eid]

		if (componentId*connectorId*resourceId)%changeStateEvery == 0 {
			if state == types.AlarmStateOK {
				state = types.AlarmStateCritical
			} else {
				state = types.AlarmStateOK
			}
		}

		stateMap[eid] = state

		err := f.send(ctx, int64(state), connectorId, componentId, resourceId)
		if err != nil {
			return err
		}

		pubcount++

		if pubcount%checkEvery == 0 {
			tsent := time.Now().UnixNano() - tstart
			adj := f.adjust(eventsPerSecond, checkEvery, tsent)
			nanosecSleep = time.Duration(nanosecSleep.Nanoseconds() + adj)
			tstart = time.Now().UnixNano()
		}

		if pubcount%sleepEvery == 0 {
			if nanosecSleep.Nanoseconds() > 0 {
				time.Sleep(nanosecSleep)
			}
		}
	}
}

func (f *Feeder) feed(ctx context.Context, eventsPerSecond float64, newResourcesPerSec float64) {
	go func() {
		err := f.sendMessages(ctx, eventsPerSecond-newResourcesPerSec, func() (int64, int64, int64) {
			connectorId := 1 + rand.Intn(numberOfConnectors)
			componentId := rand.Intn(numberOfComponents)

			resourceId := rand.Intn(f.oldResourcesMap[connectorId])

			return int64(connectorId), int64(componentId), int64(resourceId)
		})

		if err != nil {
			f.logger.Fatal().Err(err).Msg("failed to send events")
		}
	}()

	ticker := time.NewTicker(time.Millisecond * time.Duration(1000/newResourcesPerSec))
	go func() {
		for {
			<-ticker.C

			connectorId := 1 + rand.Intn(numberOfConnectors)
			componentId := rand.Intn(numberOfComponents)

			f.newResourcesMap[connectorId]++
			resourceId := f.newResourcesMap[connectorId]

			err := f.send(ctx, 1, int64(connectorId), int64(componentId), int64(resourceId))
			if err != nil {
				f.logger.Fatal().Err(err).Msg("failed to send event")
				return
			}
		}
	}()
}

func (f *Feeder) Run(ctx context.Context) error {
	if err := f.setupAmqp(); err != nil {
		return err
	}

	f.feed(ctx, float64(f.flags.EventsPerSec), float64(f.flags.NewResourcesPerSec))

	time.Sleep(time.Second * time.Duration(f.flags.FeederTime))

	return nil
}

func NewFeeder(logger zerolog.Logger) (*Feeder, error) {
	oldResourcesMap := make(map[int]int)
	newResourcesMap := make(map[int]int)

	for i := 1; i < 11; i++ {
		oldResourcesMap[i] = i * 10
		newResourcesMap[i] = i * 10
	}

	f := Feeder{
		logger:          logger,
		oldResourcesMap: oldResourcesMap,
		newResourcesMap: newResourcesMap,
	}

	f.flags.ParseArgs()
	if f.flags.Version {
		canopsis.PrintVersionInfo()
		os.Exit(0)
	}

	return &f, nil
}

func main() {
	logger := log.NewLogger(false)

	feeder, err := NewFeeder(logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("feeder init error")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	err = feeder.Run(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("feeder run error")
	}
}
