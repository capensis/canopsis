package canopsis

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"
	"time"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
)

// Returned values on exit
const (
	ExitOK = iota
	ExitPanic
	ExitEngine
)

// Engine makes engines more predictible in the way they are built.
// Used with StartEngine(), goroutine spawn is handled for you
// and you only have to implement your working logic.
//
// You can use the DefaultEngine struct to avoid implementing
// Initialize(), Stop(), Continue(), Confirm*Stop() and Has*Process()
type Engine interface {
	// PeriodicalWaitTime returns the duration to wait between
	// two run of PeriodicalProcess()
	PeriodicalWaitTime() time.Duration

	// Continue returns the state of the engine.
	// true to continue running
	// false to stop as soon as possible, meaning when functions finish
	// their job.
	Continue() bool

	// ConsumerChan must return a new channel from amqp.Channel.Consume()
	// The worker process will stop looping when the channel is closed,
	// so you must close yourself all channels.
	ConsumerChan() (<-chan amqp.Delivery, error)

	// Initialize anything you want. If err != nil then the Engine will not start.
	Initialize(ctx context.Context) error
	Started(ctx context.Context, runInfo engine.RunInfo)

	// WorkerProcess must implement the actual consumer processing.
	WorkerProcess(context.Context, amqp.Delivery)

	// PeriodicalProcess must implement the actual "beat processing".
	// Prefer using channels that you will close for long
	// PeriodicalWaitTime() values.
	PeriodicalProcess(ctx context.Context)

	SetWaitStopChan(chan os.Signal)
	AskStop(int)
	// Stop handles stopping the engine, and also waiting
	// for goroutines to finish.
	// It returns the engine's exit status, that should be used in os.Exit.
	Stop() int

	// ConfirmPeriodicalStop is useful to mutate a variable
	// or an internal chan to confirm that the Periodical process
	// is now stopped.
	ConfirmPeriodicalStop()
	ConfirmWorkerStop()

	RunPeriodicalProcess() bool
	RunWorkerProcess() bool

	// Do NOT override when using DefaultEngine.
	// Use this method to send a ConfirmStop without doing anything else.
	AcknowledgeConfirmStop()

	// WorkerEnd is called at the end of the worker routine, when one of the
	// following things happened:
	//  - the AMQP consumer channel closed
	//  - the engine has been stopped
	//  - WorkerProcess panicked
	// It should recover from panics, and notify the engine that the worker
	// routine has stopped (with the method ConfirmWorkerStop).
	WorkerEnd()

	// PeriodicalEnd is called at the end of the periodical routine, when one
	// of the following things happened:
	//  - the engine has been stopped
	//  - PeriodicalProcess panicked
	// It should recover from panics, and notify the engine that the periodical
	// routine has stopped (with the method ConfirmWorkerStop).
	PeriodicalEnd()

	Logger() *zerolog.Logger

	GetRunInfo() engine.RunInfo
}

// DefaultEngine provides basic functions and will behave as an Engine
// with PeriodicalProcess and WorkerProcess.
//
// Stop() will automatically close the Sub channel and waiting
// for confirmations to come on ConfirmChan.
type DefaultEngine struct {
	mc             sync.RWMutex
	cont           bool
	sigchan        chan os.Signal
	msa            sync.RWMutex
	stopAsked      bool
	returnValue    int
	ConfirmChan    chan bool
	Sub            libamqp.Channel
	Sleep          time.Duration
	RunWorker      bool
	RunPeriodical  bool
	Debug          bool
	logger         zerolog.Logger
	RunInfoManager engine.RunInfoManager
}

// NewDefaultEngine returns an default engine implementation.
// Check engine_test.go in this package to see how you can use it.
func NewDefaultEngine(
	sleep time.Duration,
	runworker, runperiodical bool,
	sub libamqp.Channel,
	logger zerolog.Logger,
	runInfoManager ...engine.RunInfoManager,
) DefaultEngine {
	var rim engine.RunInfoManager
	if len(runInfoManager) == 1 {
		rim = runInfoManager[0]
	} else if len(runInfoManager) > 1 {
		panic("too much arguments")
	}

	return DefaultEngine{
		cont:           true,
		returnValue:    ExitOK,
		Sleep:          sleep,
		ConfirmChan:    make(chan bool, 2),
		Sub:            sub,
		RunWorker:      runworker,
		RunPeriodical:  runperiodical,
		logger:         logger,
		RunInfoManager: rim,
	}
}

func (de *DefaultEngine) GetRunInfo() engine.RunInfo {
	return engine.RunInfo{}
}

// Initialize does nothing here besides returning nil.
func (de *DefaultEngine) Initialize(ctx context.Context) error {
	return nil
}

// Started is called before StartEngine waits for stop signal.
func (de *DefaultEngine) Started(ctx context.Context, runInfo engine.RunInfo) {
	de.Logger().Info().Msg("Engine started")
	de.saveRunInfo(ctx, runInfo)
}

func (de *DefaultEngine) saveRunInfo(ctx context.Context, runInfo engine.RunInfo) {
	if de.RunInfoManager != nil && runInfo.Name != "" {
		go func(info engine.RunInfo) {
			err := de.RunInfoManager.Save(ctx, info, de.Sleep)
			if err != nil {
				de.Logger().Error().Err(err).Msg("cannot save run info")
				return
			}

			ticker := time.NewTicker(de.Sleep)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					if !de.Continue() {
						return
					}

					err := de.RunInfoManager.Save(ctx, info, de.Sleep)
					if err != nil {
						de.Logger().Error().Err(err).Msg("cannot save run info")
						return
					}
				}
			}
		}(runInfo)
	}
}

// Continue returns the DefaultEngine.cont flag.
func (de *DefaultEngine) Continue() bool {
	de.mc.RLock()
	defer de.mc.RUnlock()
	return de.cont
}

// PeriodicalWaitTime returns the Sleep attribute.
func (de *DefaultEngine) PeriodicalWaitTime() time.Duration {
	return de.Sleep
}

// RunPeriodicalProcess returns DefaultEngine.RunPeriodical
func (de DefaultEngine) RunPeriodicalProcess() bool {
	return de.RunPeriodical
}

// RunWorkerProcess returns DefaultEngine.RunWorker
func (de DefaultEngine) RunWorkerProcess() bool {
	return de.RunWorker
}

// SetWaitStopChan keeps in memory our sigint/term chan.
func (de *DefaultEngine) SetWaitStopChan(sigchan chan os.Signal) {
	de.sigchan = sigchan
}

func (de *DefaultEngine) getStopAsked() bool {
	de.msa.RLock()
	defer de.msa.RUnlock()
	return de.stopAsked
}

// AskStop sends SIGINT into the channel where StartEngine waits for sigint/term,
// leading to a proper engine stop.
func (de *DefaultEngine) AskStop(state int) {
	if !de.getStopAsked() {
		de.returnValue = state
		de.sigchan <- syscall.SIGINT
	}
	de.msa.Lock()
	de.stopAsked = true
	de.msa.Unlock()
}

// Stop set cont flag to false, close the Sub channel then wait
// for two values from ConfirmChan.
// It returns the engine's exit status, that should be used in os.Exit.
func (de *DefaultEngine) Stop() int {
	de.Logger().Info().Msg("engine stop called, waiting for processes to finish.")
	de.mc.Lock()
	de.cont = false
	de.mc.Unlock()
	if de.Sub != nil {
		de.Logger().Debug().Msg("Closing amqp.Channel")
		err := de.Sub.Close()
		if err != nil {
			de.Logger().Error().Err(err).Msg("Error while closing amqp.Channel")
		}
	}
	<-de.ConfirmChan
	<-de.ConfirmChan

	return de.returnValue
}

// AcknowledgeConfirmStop sends true to ConfirmChan once for each call.
// Do NOT override this method when using DefaultEngine in your struct.
func (de *DefaultEngine) AcknowledgeConfirmStop() {
	de.ConfirmChan <- true
}

// ConfirmWorkerStop sends one value to ConfirmChan.
func (de *DefaultEngine) ConfirmWorkerStop() {
	de.AcknowledgeConfirmStop()
}

// ConfirmPeriodicalStop sends one value to ConfirmChan.
func (de *DefaultEngine) ConfirmPeriodicalStop() {
	de.AcknowledgeConfirmStop()
}

// StartEngine handles starting the WorkerProcess and PeriodicalProcess
// of your Engine.
// It first calls Engine.Initialize(), returns an error if any, then
// proceed to signals binding.
//
// Engine.PeriodicalProcess and Engine.WorkerProcess are launched into
// separate goroutines.
//
// Only SIGTERM or SIGINT will trigger the Engine.Stop() method.
//
// waitChan is optional: if nil a chan will be instanciated and
// managed for you. This is mainly for tests.
func StartEngine(ctx context.Context, engine Engine, waitChan *chan os.Signal) (int, error) {
	if !engine.RunPeriodicalProcess() && !engine.RunWorkerProcess() {
		return ExitEngine, errors.New("impossible engine configuration: must have at least one Process()")
	}

	err := engine.Initialize(ctx)
	if err != nil {
		return ExitEngine, err
	}

	if waitChan == nil {
		w := make(chan os.Signal, 1)
		waitChan = &w
	}

	engine.SetWaitStopChan(*waitChan)
	signal.Notify(*waitChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	if engine.RunWorkerProcess() {
		channel, err := engine.ConsumerChan()
		if err != nil {
			return ExitEngine, err
		}
		go routineWorker(ctx, engine, channel)
	} else {
		engine.AcknowledgeConfirmStop()
	}

	if engine.RunPeriodicalProcess() {
		go routinePeriodical(ctx, engine)
	} else {
		engine.AcknowledgeConfirmStop()
	}

	engine.Started(ctx, engine.GetRunInfo())

	<-*waitChan
	return engine.Stop(), nil
}

func (de *DefaultEngine) WorkerEnd() {
	exitCode := ExitOK
	if r := recover(); r != nil {
		de.Logger().Error().Msgf("worker recovered from panic: %v", r)
		debug.PrintStack()
		exitCode = ExitPanic
	}

	de.AskStop(exitCode)
	de.ConfirmWorkerStop()
}

func (de *DefaultEngine) PeriodicalEnd() {
	exitCode := ExitOK
	if r := recover(); r != nil {
		de.Logger().Error().Msgf("periodical recovered from panic: %v", r)
		debug.PrintStack()
		exitCode = ExitPanic
	}

	de.AskStop(exitCode)
	de.ConfirmPeriodicalStop()
}

func (de *DefaultEngine) Logger() *zerolog.Logger {
	return &de.logger
}

func (de *DefaultEngine) AckMessage(msg amqp.Delivery) {
	err := de.Sub.Ack(msg.DeliveryTag, false)
	if err != nil {
		de.logger.Debug().Err(err).Msg("fail to ack message")
	} else {
		de.logger.Debug().
			Uint64("tag", msg.DeliveryTag).
			Str("id", msg.MessageId).
			Str("exchange", msg.Exchange).
			Str("routing_key", msg.RoutingKey).
			Str("consumer_tag", msg.ConsumerTag).Msg("ack message")
	}
}

func (de *DefaultEngine) NackMessage(msg amqp.Delivery) {
	err := de.Sub.Nack(msg.DeliveryTag, false, true)
	if err != nil {
		de.logger.Debug().Err(err).Msg("fail to nack message")
	} else {
		de.logger.Debug().
			Uint64("tag", msg.DeliveryTag).
			Str("id", msg.MessageId).
			Str("exchange", msg.Exchange).
			Str("routing_key", msg.RoutingKey).
			Str("consumer_tag", msg.ConsumerTag).Msg("nack message")
	}
}

// ProcessWorkerError nacks message if external services are not reachable and stops engine.
// It acks messages on other errors and continues engine working.
func (de *DefaultEngine) ProcessWorkerError(err error, msg amqp.Delivery) {
	if mongo.IsConnectionError(err) || redis.IsConnectionError(err) {
		de.NackMessage(msg)
		de.AskStop(ExitEngine)
		return
	}

	de.AckMessage(msg)
}

func routineWorker(ctx context.Context, engine Engine, channel <-chan amqp.Delivery) {
	defer engine.WorkerEnd()

	timeout := make(chan bool, 1)
	go func() {
		for engine.Continue() {
			time.Sleep(time.Second * 1)
			timeout <- true
		}
	}()

	for engine.Continue() {
		select {
		case msg, ok := <-channel:
			if !ok {
				engine.Logger().Error().
					Msg("the rabbitmq channel has been closed")
				return
			}

			engine.WorkerProcess(ctx, msg)
		case <-timeout:
		}
	}
}

// watch engine.Continue every 2 second, during a maximum amount of time
// given by engine.PeriodicalWaitTime(), and return engine.Continue().
func waitAndCheck(engine Engine) bool {
	sleepCountDown := engine.PeriodicalWaitTime()
	for sleepCountDown > 0 && engine.Continue() {
		wait := time.Second * 2
		time.Sleep(wait)
		sleepCountDown = sleepCountDown - wait
	}
	return engine.Continue()
}

func routinePeriodical(ctx context.Context, engine Engine) {
	defer engine.PeriodicalEnd()

	for waitAndCheck(engine) {
		engine.PeriodicalProcess(ctx)
	}
}
