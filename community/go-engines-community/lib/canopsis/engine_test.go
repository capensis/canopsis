package canopsis_test

import (
	"errors"
	"os"
	"syscall"
	"testing"
	"time"

	cps "git.canopsis.net/canopsis/go-engines/lib/canopsis"
	"git.canopsis.net/canopsis/go-engines/lib/log"
	"git.canopsis.net/canopsis/go-engines/lib/testutils"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/streadway/amqp"
)

type testEngine struct {
	cps.DefaultEngine
	Messages    []string
	TestChannel *chan amqp.Delivery
}

type testEngineBadConsumerChan struct {
	cps.DefaultEngine
}

func (eng *testEngineBadConsumerChan) WorkerProcess(msg amqp.Delivery) {
}

func (eng *testEngineBadConsumerChan) ConsumerChan() (<-chan amqp.Delivery, error) {
	return nil, errors.New("lel, i failed")
}

func (eng *testEngineBadConsumerChan) PeriodicalProcess() {

}

type testEngineInitializeFailed struct {
	testEngine
}

func (eng *testEngineInitializeFailed) Initialize() error {
	return errors.New("lel, i'm dumb")
}

func (eng *testEngine) ConsumerChan() (<-chan amqp.Delivery, error) {
	msgs := make(chan amqp.Delivery, 10)
	msgs <- amqp.Delivery{Body: []byte("coucou0")}
	msgs <- amqp.Delivery{Body: []byte("coucou1")}
	msgs <- amqp.Delivery{Body: []byte("coucou2")}
	msgs <- amqp.Delivery{Body: []byte("coucou3")}
	msgs <- amqp.Delivery{Body: []byte("coucou4")}
	msgs <- amqp.Delivery{Body: []byte("coucou5")}
	msgs <- amqp.Delivery{Body: []byte("coucou6")}
	msgs <- amqp.Delivery{Body: []byte("coucou7")}
	msgs <- amqp.Delivery{Body: []byte("coucou8")}
	msgs <- amqp.Delivery{Body: []byte("coucou9")}
	return msgs, nil
}

func (eng *testEngine) PeriodicalProcess() {
	eng.Logger().Info().Msg("PERIODICAL: working 1sec...")
	time.Sleep(time.Second)
	eng.Logger().Info().Msg("PERIODICAL: work finished")
}

func (eng *testEngine) WorkerProcess(msg amqp.Delivery) {
	eng.Messages = append(eng.Messages, string(msg.Body))
}

func TestDefaultEngine(t *testing.T) {
	testutils.SkipLongIfSet(t)

	Convey("Test different engines that MUST NOT PANIC", t, func() {
		waitChan := make(chan os.Signal, 1)

		ts := time.Now()

		logger := log.NewTestLogger()

		stopAfter2Sec := func() {
			time.Sleep(time.Second * 2)
			logger.Info().Msg("TEST: sending SIGINT")
			waitChan <- syscall.SIGINT
		}

		Convey("Test an engine with worker and periodical", func() {
			engine := testEngine{
				DefaultEngine: cps.NewDefaultEngine(time.Second*20, true, true, nil, logger),
				Messages:      make([]string, 0),
			}
			go stopAfter2Sec()

			logger.Info().Msg("TEST: start engine")
			exitStatus, err := cps.StartEngine(&engine, &waitChan)
			So(exitStatus, ShouldEqual, cps.ExitOK)
			So(err, ShouldBeNil)

			stopTime := time.Now().Sub(ts)
			So(stopTime, ShouldBeGreaterThan, time.Second*2)
			So(stopTime, ShouldBeLessThan, time.Second*20)
			So(len(engine.Messages), ShouldBeGreaterThan, 0)
			So(engine.Messages[0], ShouldEqual, "coucou0")
		})

		Convey("Test a bad engine", func() {
			engine := testEngine{
				DefaultEngine: cps.NewDefaultEngine(0, false, false, nil, log.NewTestLogger()),
			}

			exitStatus, err := cps.StartEngine(&engine, &waitChan)
			So(exitStatus, ShouldEqual, cps.ExitEngine)
			So(err, ShouldNotBeNil)
		})

		Convey("Test worker engine", func() {
			engine := testEngine{
				DefaultEngine: cps.NewDefaultEngine(0, true, false, nil, log.NewTestLogger()),
			}
			stopAfter2Sec()
			exitStatus, err := cps.StartEngine(&engine, &waitChan)
			So(exitStatus, ShouldEqual, cps.ExitOK)
			So(err, ShouldBeNil)
		})

		Convey("Test periodical engine", func() {
			engine := testEngine{
				DefaultEngine: cps.NewDefaultEngine(0, false, true, nil, log.NewTestLogger()),
			}
			stopAfter2Sec()
			exitStatus, err := cps.StartEngine(&engine, &waitChan)
			So(exitStatus, ShouldEqual, cps.ExitOK)
			So(err, ShouldBeNil)
		})

		Convey("Consumer failed", func() {
			engine := testEngineBadConsumerChan{
				DefaultEngine: cps.NewDefaultEngine(0, true, false, nil, log.NewTestLogger()),
			}
			stopAfter2Sec()
			exitStatus, err := cps.StartEngine(&engine, nil)
			So(exitStatus, ShouldEqual, cps.ExitEngine)
			So(err.Error(), ShouldEqual, "lel, i failed")
		})

		Convey("Initialized failed", func() {
			engine := testEngineInitializeFailed{
				testEngine: testEngine{
					DefaultEngine: cps.NewDefaultEngine(0, true, false, nil, log.NewTestLogger()),
				},
			}

			stopAfter2Sec()
			exitStatus, err := cps.StartEngine(&engine, nil)
			So(exitStatus, ShouldEqual, cps.ExitEngine)
			So(err.Error(), ShouldEqual, "lel, i'm dumb")
		})
	})
}

type WorkerPanicEngine struct {
	cps.DefaultEngine
}

func (eng *WorkerPanicEngine) ConsumerChan() (<-chan amqp.Delivery, error) {
	msgs := make(chan amqp.Delivery, 10)
	msgs <- amqp.Delivery{
		Body: []byte("pwet0"),
	}
	msgs <- amqp.Delivery{
		Body: []byte("pwet1"),
	}
	msgs <- amqp.Delivery{
		Body: []byte("pwet2"),
	}
	msgs <- amqp.Delivery{
		Body: []byte("pwet3"),
	}
	msgs <- amqp.Delivery{
		Body: []byte("panic"),
	}
	msgs <- amqp.Delivery{
		Body: []byte("pwet4"),
	}
	return msgs, nil
}

func (eng *WorkerPanicEngine) PeriodicalProcess() {
	eng.Logger().Info().Msg("periodical coucou")
}

func (eng *WorkerPanicEngine) WorkerProcess(msg amqp.Delivery) {
	if string(msg.Body) == "panic" {
		panic("jepanikay")
	} else {
		eng.Logger().Info().Msg(string(msg.Body))
	}
	time.Sleep(time.Second * 1)
}

func TestWorkerPanicEngine(t *testing.T) {
	testutils.SkipLongIfSet(t)
	Convey("Setup panic engine on worker process", t, func() {
		pe := WorkerPanicEngine{
			DefaultEngine: cps.NewDefaultEngine(
				time.Second*2,
				true,
				true,
				nil,
				log.NewTestLogger(),
			),
		}

		startEngine := func() {
			cps.StartEngine(&pe, nil)
		}

		So(startEngine, ShouldNotPanic)
	})

}

type PeriodicalPanicEngine struct {
	cps.DefaultEngine
}

func (eng *PeriodicalPanicEngine) ConsumerChan() (<-chan amqp.Delivery, error) {
	msgs := make(chan amqp.Delivery, 6)
	msgs <- amqp.Delivery{
		Body: []byte("pwet0"),
	}
	msgs <- amqp.Delivery{
		Body: []byte("pwet1"),
	}
	msgs <- amqp.Delivery{
		Body: []byte("pwet2"),
	}
	msgs <- amqp.Delivery{
		Body: []byte("pwet3"),
	}
	msgs <- amqp.Delivery{
		Body: []byte("pasdpanicabord"),
	}
	msgs <- amqp.Delivery{
		Body: []byte("pwet4"),
	}
	return msgs, nil
}

func (eng *PeriodicalPanicEngine) PeriodicalProcess() {
	panic("lavitesseetlefundabord")
}

func (eng *PeriodicalPanicEngine) WorkerProcess(msg amqp.Delivery) {
	if string(msg.Body) == "panic" {
		panic("jepanikay")
	} else {
		eng.Logger().Info().Msg(string(msg.Body))
	}
	time.Sleep(time.Second * 1)
}

func TestPeriodicalPanicEngine(t *testing.T) {
	testutils.SkipLongIfSet(t)
	Convey("Setup panic engine on periodical process", t, func() {
		pe := PeriodicalPanicEngine{
			DefaultEngine: cps.NewDefaultEngine(
				time.Second*2,
				true,
				true,
				nil,
				log.NewTestLogger(),
			),
		}

		startEngine := func() {
			cps.StartEngine(&pe, nil)
		}

		So(startEngine, ShouldNotPanic)
	})

}
