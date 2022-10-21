package engine_test

import (
	"context"
	"errors"
	"testing"
	"time"

	libengine "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	mock_engine "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/engine"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
)

const waitTimeout = time.Second
const interval = time.Millisecond * 100
const inaccuracy = interval / 100

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
	engine.AddPeriodicalWorker(mockPeriodicalWorker)

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
		if sub < interval-inaccuracy || sub >= 2*(interval-inaccuracy) {
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

func TestEngine_Run_GivenErrorOnPeriodicalProcess_ShouldStopEngine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan bool)
	defer close(done)

	mockPeriodicalWorker := mock_engine.NewMockPeriodicalWorker(ctrl)
	expectedErr := &testErr{msg: "test error"}
	mockPeriodicalWorker.EXPECT().GetInterval().Return(interval).AnyTimes()
	mockPeriodicalWorker.EXPECT().Work(gomock.Any()).Return(expectedErr)

	engine := libengine.New(nil, nil, zerolog.Logger{})
	engine.AddPeriodicalWorker(mockPeriodicalWorker)

	var err error
	go func() {
		err = engine.Run(ctx)
		done <- true
	}()

	waitDone(t, done)

	testErr := &testErr{}
	if !errors.As(err, &testErr) || testErr.Error() != expectedErr.Error() {
		t.Errorf("expected error %v but got %v", expectedErr, err)
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
	expectedErr := &testErr{msg: "test error"}
	mockConsumer.EXPECT().Consume(gomock.Any()).Return(expectedErr)

	engine := libengine.New(nil, nil, zerolog.Logger{})
	engine.AddConsumer(mockConsumer)

	var err error
	go func() {
		err = engine.Run(ctx)
		done <- true
	}()

	waitDone(t, done)

	testErr := &testErr{}
	if !errors.As(err, &testErr) || testErr.Error() != expectedErr.Error() {
		t.Errorf("expected error %v but got %v", expectedErr, err)
	}
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
