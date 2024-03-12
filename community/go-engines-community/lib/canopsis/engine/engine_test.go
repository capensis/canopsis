package engine_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	libengine "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	mock_engine "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/engine"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
)

const waitTimeout = time.Second
const interval = 100 * time.Millisecond
const inaccuracy = 5 * time.Millisecond

func TestEngine_Run_GivenPeriodicalProcess_ShouldRunIt(t *testing.T) {
	const timesToRun = 2
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan bool)
	defer close(done)

	mockPeriodicalWorker := mock_engine.NewMockPeriodicalWorker(ctrl)
	mockPeriodicalWorker.EXPECT().GetInterval().Return(interval).AnyTimes()
	workTimes := make([]time.Time, 0)
	mockPeriodicalWorker.EXPECT().Work(gomock.Any()).
		Do(func(_ context.Context) {
			workTimes = append(workTimes, time.Now())
			if len(workTimes) == timesToRun {
				cancel()
			}
		}).
		Times(timesToRun)

	engine := libengine.New(nil, nil, zerolog.Nop())
	engine.AddPeriodicalWorker("test", mockPeriodicalWorker)

	var start time.Time
	var err error

	go func() {
		start = time.Now()
		err = engine.Run(ctx)
		done <- true
	}()

	waitDone(t, done)

	if err != nil {
		t.Errorf("expected no error but got %v", err)
		return
	}

	if len(workTimes) != timesToRun {
		t.Errorf("expected %v periodical executions but got %v", timesToRun, len(workTimes))
		return
	}

	for _, date := range workTimes {
		sub := date.Sub(start)
		if sub < interval-inaccuracy || sub > interval+inaccuracy {
			t.Errorf("expected %v between periodical executions but got %v", interval, sub)
			return
		}

		start = date
	}
}

func TestEngine_Run_GivenConsumer_ShouldRunIt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan bool)
	defer close(done)

	mockConsumer := mock_engine.NewMockConsumer(ctrl)
	mockConsumer.EXPECT().Consume(gomock.Any()).Do(func(_ context.Context) {
		cancel()
	})

	engine := libengine.New(nil, nil, zerolog.Nop())
	engine.AddConsumer(mockConsumer)

	var err error
	go func() {
		err = engine.Run(ctx)
		done <- true
	}()

	waitDone(t, done)

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
}

func TestEngine_Run_GivenErrorOnConsumer_ShouldStopEngine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan bool)
	defer close(done)

	mockConsumer := mock_engine.NewMockConsumer(ctrl)
	expectedErr := &testError{msg: "test error"}
	mockConsumer.EXPECT().Consume(gomock.Any()).Return(expectedErr)

	engine := libengine.New(nil, nil, zerolog.Logger{})
	engine.AddConsumer(mockConsumer)

	var err error
	go func() {
		err = engine.Run(ctx)
		done <- true
	}()

	waitDone(t, done)

	testErr := &testError{}
	if !errors.As(err, &testErr) || testErr.Error() != expectedErr.Error() {
		t.Errorf("expected error %v but got %v", expectedErr, err)
	}
}

func TestEngine_Run_GivenRoutine_ShouldRunIt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	expectedDuration := 100 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), expectedDuration)
	defer cancel()
	done := make(chan bool)
	defer close(done)

	var duration1Mx, duration2Mx sync.Mutex
	var duration1, duration2 time.Duration
	engine := libengine.New(nil, nil, zerolog.Nop())
	engine.AddRoutine(func(ctx context.Context) error {
		started := time.Now()
		defer func() {
			duration1Mx.Lock()
			duration1 = time.Since(started)
			duration1Mx.Unlock()
		}()

		<-ctx.Done()
		return nil
	})
	engine.AddRoutine(func(ctx context.Context) error {
		started := time.Now()
		defer func() {
			duration2Mx.Lock()
			duration2 = time.Since(started)
			duration2Mx.Unlock()
		}()

		<-ctx.Done()
		return nil
	})

	var err error

	go func() {
		err = engine.Run(ctx)
		done <- true
	}()

	waitDone(t, done)

	if err != nil {
		t.Errorf("expected no error but got %v", err)
		return
	}

	duration1Mx.Lock()
	if duration1 < expectedDuration-inaccuracy || duration1 > expectedDuration+inaccuracy {
		t.Errorf("expected %s but got %s", expectedDuration, duration1)
	}
	duration1Mx.Unlock()
	duration2Mx.Lock()
	if duration2 < expectedDuration-inaccuracy || duration2 > expectedDuration+inaccuracy {
		t.Errorf("expected %s but got %s", expectedDuration, duration2)
	}
	duration2Mx.Unlock()
}

func waitDone(t *testing.T, done <-chan bool) {
	select {
	case <-time.After(waitTimeout):
		t.Error("timeout expired")
	case _, ok := <-done:
		if !ok {
			t.Error("channel closed")
		}
	}
}
