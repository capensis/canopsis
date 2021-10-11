package engine_test

import (
	"context"
	"reflect"
	"sort"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
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
			Time:             types.CpsTime{Time: now.Add(-2 * time.Second)},
		},
		{
			ID:               "test-axe-2",
			Name:             "test-axe",
			ConsumeQueue:     "test-consume-axe-2",
			PublishQueue:     "test-publish-axe-2",
			RpcConsumeQueues: []string{"test-rpc-consume-axe"},
			RpcPublishQueues: []string{"test-rpc-publish-axe"},
			QueueLength:      10,
			Time:             types.CpsTime{Time: now.Add(-time.Second)},
		},
		{
			ID:               "test-che-1",
			Name:             "test-che",
			ConsumeQueue:     "test-consume-che",
			PublishQueue:     "test-publish-che",
			RpcConsumeQueues: []string{"test-rpc-consume-che"},
			RpcPublishQueues: []string{"test-rpc-publish-che"},
			QueueLength:      11,
			Time:             types.CpsTime{Time: now.Add(-3 * time.Second)},
		},
	}
	expected := []engine.RunInfo{
		{
			Name:         "test-axe",
			ConsumeQueue: "test-consume-axe-2",
			PublishQueue: "test-publish-axe-2",
		},
		{
			Name:         "test-che",
			ConsumeQueue: "test-consume-che",
			PublishQueue: "test-publish-che",
		},
	}

	for _, v := range instances {
		err = manager.SaveInstance(ctx, v, expiration)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}
	}

	engines, err := manager.GetEngineQueues(ctx)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	sort.Slice(engines, func(i, j int) bool {
		return engines[i].Name < engines[j].Name
	})

	if !reflect.DeepEqual(engines, expected) {
		t.Errorf("expected %+v but got %+v", expected, engines)
	}

	err = client.FlushDB(ctx).Err()
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
}
