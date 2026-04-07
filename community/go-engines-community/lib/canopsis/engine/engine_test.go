package engine_test

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"
	"time"

	libengine "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	mock_engine "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/engine"
	"github.com/rs/zerolog"
	"go.uber.org/mock/gomock"
)

const waitTimeout = time.Second
const interval = 100 * time.Millisecond

func TestEngine_Run_GivenPeriodicalProcess_ShouldRunIt_SyntheticClock(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		const timesToRun = 2

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mockPeriodicalWorker := mock_engine.NewMockPeriodicalWorker(ctrl)
		mockPeriodicalWorker.EXPECT().GetInterval().Return(interval).AnyTimes()

		workTimes := make([]time.Time, 0, timesToRun)
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

		go func() {
			if err := engine.Run(ctx); err != nil {
				t.Errorf("expected no error but got %v", err)
			}
		}()

		// park until all goroutines are blocked on the ticker.
		synctest.Wait()

		for i := 0; i < timesToRun; i++ {
			time.Sleep(interval)
			synctest.Wait() // let Work() execute and record the time
		}

		// wait for engine's shut down by now due cancelled context
		synctest.Wait()

		if len(workTimes) != timesToRun {
			t.Errorf("expected %v periodical executions but got %v", timesToRun, len(workTimes))
			return
		}

		// Check that each subsequent execution happens exactly at the expected interval.
		for i := 1; i < len(workTimes); i++ {
			sub := workTimes[i].Sub(workTimes[i-1])
			if sub != interval { // exact equality, no jitter in fake clock
				t.Errorf("expected interval %v but got %v between run %d and %d", interval, sub, i-1, i)
			}
		}
	})
}

func TestEngine_Run_GivenConsumer_ShouldRunIt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(t.Context())
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

	done := make(chan bool)
	defer close(done)

	mockConsumer := mock_engine.NewMockConsumer(ctrl)
	expectedErr := &testError{msg: "test error"}
	mockConsumer.EXPECT().Consume(gomock.Any()).Return(expectedErr)

	engine := libengine.New(nil, nil, zerolog.Nop())
	engine.AddConsumer(mockConsumer)

	var err error
	go func() {
		err = engine.Run(t.Context())
		done <- true
	}()

	waitDone(t, done)

	testErr := &testError{}
	if !errors.As(err, &testErr) || testErr.Error() != expectedErr.Error() {
		t.Errorf("expected error %v but got %v", expectedErr, err)
	}
}

func TestEngine_Run_GivenRoutine_ShouldRunIt_SyntheticClock(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expectedDuration := 100 * time.Millisecond
		ctx, cancel := context.WithTimeout(t.Context(), expectedDuration)
		defer cancel()

		var duration1, duration2 time.Duration
		engine := libengine.New(nil, nil, zerolog.Nop())
		engine.AddRoutine(func(ctx context.Context) error {
			started := time.Now()
			defer func() { duration1 = time.Since(started) }()

			<-ctx.Done()
			return nil
		})
		engine.AddRoutine(func(ctx context.Context) error {
			started := time.Now()
			defer func() { duration2 = time.Since(started) }()

			<-ctx.Done()
			return nil
		})

		var err error
		go func() {
			err = engine.Run(ctx)
		}()

		synctest.Wait() // park until all goroutines are blocked on the ticker
		time.Sleep(expectedDuration)
		synctest.Wait() // routines finish, engine.Run returns

		if err != nil {
			t.Errorf("expected no error but got %v", err)
			return
		}

		if duration1 != expectedDuration {
			t.Errorf("routine 1: expected %v but got %v", expectedDuration, duration1)
		}
		if duration2 != expectedDuration {
			t.Errorf("routine 2: expected %v but got %v", expectedDuration, duration2)
		}
	})
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
