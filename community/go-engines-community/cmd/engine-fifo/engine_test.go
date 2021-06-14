package main_test

import (
	"context"
	"testing"

	fifo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/cmd/engine-fifo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	"github.com/streadway/amqp"

	. "github.com/smartystreets/goconvey/convey"
)

type EngineTestFIFO struct {
	*fifo.EngineFIFO
}

func (e *EngineTestFIFO) ConsumerChan() (<-chan amqp.Delivery, error) {
	ch := make(chan amqp.Delivery)
	return ch, nil
}

func (e *EngineTestFIFO) AckConsumerChan() (<-chan amqp.Delivery, error) {
	ch := make(chan amqp.Delivery)
	return ch, nil
}

func testNewEngineFIFO() (*EngineTestFIFO, error) {
	options := fifo.Options{
		PublishToQueue:   canopsis.CheQueueName,
		ConsumeFromQueue: canopsis.FIFOQueueName,
	}
	depMaker := fifo.DependencyMaker{}
	references := depMaker.GetDefaultReferences(context.Background(), options, log.NewTestLogger())

	engine := EngineTestFIFO{
		EngineFIFO: fifo.NewEngineFIFO(options, references),
	}
	_, err := engine.ConsumerChan()
	if err != nil {
		return nil, err
	}
	_, err = engine.AckConsumerChan()

	return &engine, err
}

func TestInitializeFIFO(t *testing.T) {
	Convey("Setup", t, func() {
		engine, err := testNewEngineFIFO()
		So(err, ShouldBeNil)

		Convey("Initialize will set correct active state", func() {
			err := engine.Initialize(context.Background())
			So(err, ShouldBeNil)
		})
	})
}
