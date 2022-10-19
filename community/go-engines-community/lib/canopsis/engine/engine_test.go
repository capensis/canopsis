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

func TestEngine_Run_GivenPeriodicalProcess_ShouldRunIt(t *testing.T) {
	const (
		timesToRun = 2
		interval   = 200 * time.Millisecond
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPeriodicalWorker := mock_engine.NewMockPeriodicalWorker(ctrl)
	finish := (timesToRun+1)*interval - time.Millisecond
	mockPeriodicalWorker.EXPECT().GetInterval().Return(interval).AnyTimes()
	mockPeriodicalWorker.EXPECT().Work(gomock.Any()).Times(timesToRun)

	engine := libengine.New(nil, nil, zerolog.Logger{})
	engine.AddPeriodicalWorker(mockPeriodicalWorker)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		_ = engine.Run(ctx)
	}()

	time.Sleep(finish)
	cancel()
}

func TestEngine_Run_GivenConsumer_ShouldRunIt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockConsumer := mock_engine.NewMockConsumer(ctrl)
	mockConsumer.EXPECT().Consume(gomock.Any())

	engine := libengine.New(nil, nil, zerolog.Logger{})
	engine.AddConsumer(mockConsumer)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = engine.Run(ctx)
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()
}

func TestEngine_Run_GivenErrorOnPeriodicalProcess_ShouldStopEngine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPeriodicalWorker := mock_engine.NewMockPeriodicalWorker(ctrl)
	interval := 500 * time.Millisecond
	expectedErr := errors.New("test err")
	mockPeriodicalWorker.EXPECT().GetInterval().Return(interval).AnyTimes()
	mockPeriodicalWorker.EXPECT().Work(gomock.Any()).Return(expectedErr)

	engine := libengine.New(nil, nil, zerolog.Logger{})
	engine.AddPeriodicalWorker(mockPeriodicalWorker)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := engine.Run(ctx)

	if err != expectedErr {
		t.Errorf("expected error %v but got %v", expectedErr, err)
	}
}

func TestEngine_Run_GivenErrorOnConsumer_ShouldStopEngine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockConsumer := mock_engine.NewMockConsumer(ctrl)
	expectedErr := errors.New("test err")
	mockConsumer.EXPECT().Consume(gomock.Any()).Return(expectedErr)

	engine := libengine.New(nil, nil, zerolog.Logger{})
	engine.AddConsumer(mockConsumer)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := engine.Run(ctx)

	if err != expectedErr {
		t.Errorf("expected error %v but got %v", expectedErr, err)
	}
}
