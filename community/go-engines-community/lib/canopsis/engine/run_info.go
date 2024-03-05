package engine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	libredis "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/redis/go-redis/v9"
)

// NewRunInfoManager creates new run info manager.
func NewRunInfoManager(client redis.Cmdable, key ...string) RunInfoManager {
	k := libredis.RunInfoKey
	if len(key) == 1 {
		k = key[0]
	} else if len(key) > 1 {
		panic("too much arguments")
	}

	return &runInfoManager{
		client: client,
		key:    k,
	}
}

// runInfoManager implements redis storage.
type runInfoManager struct {
	client redis.Cmdable
	key    string
}

func (m *runInfoManager) SaveInstance(ctx context.Context, info InstanceRunInfo, expiration time.Duration) error {
	if info.ID == "" {
		return errors.New("id is required")
	}
	if info.Name == "" {
		return errors.New("name is required")
	}

	b, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("cannot marshal info: %w", err)
	}

	res := m.client.Set(ctx, m.getCacheKey(info), b, expiration)
	if err := res.Err(); err != nil {
		return fmt.Errorf("cannot save info to cache: %w", err)
	}

	return nil
}

func (m *runInfoManager) GetEngines(ctx context.Context) ([]RunInfo, error) {
	infosByName := make(map[string]*RunInfo)
	var cursor uint64
	processedKeys := make(map[string]bool)

	for {
		res := m.client.Scan(ctx, cursor, fmt.Sprintf("%s*", m.key), 50)
		if err := res.Err(); err != nil {
			return nil, fmt.Errorf("cannot scan cache keys: %w", err)
		}

		var keys []string
		keys, cursor = res.Val()
		unprocessedKeys := make([]string, 0)
		for _, key := range keys {
			if !processedKeys[key] {
				unprocessedKeys = append(unprocessedKeys, key)
				processedKeys[key] = true
			}
		}

		if len(unprocessedKeys) > 0 {
			resGet := m.client.MGet(ctx, unprocessedKeys...)
			if err := resGet.Err(); err != nil {
				return nil, fmt.Errorf("cannot fetch infos from cache: %w", err)
			}

			for i, v := range resGet.Val() {
				if v == nil {
					continue
				}

				s, ok := v.(string)
				if !ok {
					return nil, fmt.Errorf("expect string for key=%q but got type=%T", unprocessedKeys[i], v)
				}

				var info InstanceRunInfo
				err := json.Unmarshal([]byte(s), &info)
				if err != nil {
					return nil, fmt.Errorf("cannot unmarshal info key=%q: %w", unprocessedKeys[i], err)
				}

				sort.Strings(info.RpcConsumeQueues)
				sort.Strings(info.RpcPublishQueues)

				if v, ok := infosByName[info.Name]; ok {
					v.Instances++

					if v.ConsumeQueue != info.ConsumeQueue || v.PublishQueue != info.PublishQueue ||
						!reflect.DeepEqual(v.RpcConsumeQueues, info.RpcConsumeQueues) ||
						!reflect.DeepEqual(v.RpcPublishQueues, info.RpcPublishQueues) {
						v.HasDiffConfig = true
					}

					if v.Time.Before(info.Time) {
						v.QueueLength = info.QueueLength
						v.Time = info.Time

						if v.HasDiffConfig {
							v.ConsumeQueue = info.ConsumeQueue
							v.PublishQueue = info.PublishQueue
							v.RpcConsumeQueues = info.RpcConsumeQueues
							v.RpcPublishQueues = info.RpcPublishQueues
						}
					}
				} else {
					infosByName[info.Name] = &RunInfo{
						Name:             info.Name,
						ConsumeQueue:     info.ConsumeQueue,
						PublishQueue:     info.PublishQueue,
						RpcConsumeQueues: info.RpcConsumeQueues,
						RpcPublishQueues: info.RpcPublishQueues,
						Instances:        1,
						QueueLength:      info.QueueLength,
						Time:             info.Time,
						HasDiffConfig:    false,
					}
				}
			}
		}

		if cursor == 0 {
			break
		}
	}

	infos := make([]RunInfo, len(infosByName))
	i := 0
	for _, info := range infosByName {
		infos[i] = *info
		i++
	}

	sort.Slice(infos, func(i, j int) bool {
		return infos[i].Name < infos[j].Name
	})

	return infos, nil
}

func (m *runInfoManager) getCacheKey(info InstanceRunInfo) string {
	return strings.Join([]string{m.key, info.Name, info.ID}, libredis.KeyDelimiter)
}
