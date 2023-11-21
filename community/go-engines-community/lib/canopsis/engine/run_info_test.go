package engine_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/rs/zerolog"
)

func TestNewRunInfoManager(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := redis.NewSession(ctx, redis.EngineRunInfo, zerolog.Nop(), 0, 0)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	err = client.FlushDB(ctx).Err()
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	manager := engine.NewRunInfoManager(client)
	expiration := time.Second * 5
	now := time.Now()

	instances := []engine.InstanceRunInfo{
		{
			ID:               "test-axe-1",
			Name:             "test-axe",
			ConsumeQueue:     "test-consume-axe-1",
			PublishQueue:     "test-publish-axe-1",
			RpcConsumeQueues: []string{"test-rpc-consume-axe"},
			RpcPublishQueues: []string{"test-rpc-publish-axe"},
			QueueLength:      9,
			Time:             datetime.CpsTime{Time: now.Add(-2 * time.Second)},
		},
		{
			ID:               "test-axe-2",
			Name:             "test-axe",
			ConsumeQueue:     "test-consume-axe-2",
			PublishQueue:     "test-publish-axe-2",
			RpcConsumeQueues: []string{"test-rpc-consume-axe"},
			RpcPublishQueues: []string{"test-rpc-publish-axe"},
			QueueLength:      10,
			Time:             datetime.CpsTime{Time: now.Add(-time.Second)},
		},
		{
			ID:               "test-che-1",
			Name:             "test-che",
			ConsumeQueue:     "test-consume-che",
			PublishQueue:     "test-publish-che",
			RpcConsumeQueues: []string{"test-rpc-consume-che"},
			RpcPublishQueues: []string{"test-rpc-publish-che"},
			QueueLength:      11,
			Time:             datetime.CpsTime{Time: now.Add(-3 * time.Second)},
		},
	}
	expected := []engine.RunInfo{
		{
			Name:             "test-axe",
			ConsumeQueue:     "test-consume-axe-2",
			PublishQueue:     "test-publish-axe-2",
			RpcConsumeQueues: []string{"test-rpc-consume-axe"},
			RpcPublishQueues: []string{"test-rpc-publish-axe"},
			Instances:        2,
			QueueLength:      10,
			Time:             instances[1].Time,
			HasDiffConfig:    true,
		},
		{
			Name:             "test-che",
			ConsumeQueue:     "test-consume-che",
			PublishQueue:     "test-publish-che",
			RpcConsumeQueues: []string{"test-rpc-consume-che"},
			RpcPublishQueues: []string{"test-rpc-publish-che"},
			Instances:        1,
			QueueLength:      11,
			Time:             instances[2].Time,
			HasDiffConfig:    false,
		},
	}

	for _, v := range instances {
		err = manager.SaveInstance(ctx, v, expiration)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}
	}

	engines, err := manager.GetEngines(ctx)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	times := make([]datetime.CpsTime, len(engines))
	expectedTimes := make([]datetime.CpsTime, len(expected))
	for i, v := range engines {
		times[i] = v.Time
		engines[i].Time = datetime.CpsTime{}
	}
	for i, v := range expected {
		expectedTimes[i] = v.Time
		expected[i].Time = datetime.CpsTime{}
	}

	if reflect.DeepEqual(engines, expected) {
		for i := range times {
			if times[i].Unix() != expectedTimes[i].Unix() {
				t.Errorf("expected time %+v for %q but got %+v", expectedTimes[i].Unix(), expected[i].Name, times[i].Unix())
			}
		}
	} else {
		t.Errorf("expected %+v but got %+v", expected, engines)
	}

	err = client.FlushDB(ctx).Err()
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
}
