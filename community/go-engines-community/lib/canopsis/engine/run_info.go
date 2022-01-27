package engine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	libredis "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/go-redis/redis/v8"
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

	res := m.client.Set(ctx, m.GetCacheKey(info), b, expiration)
	if err := res.Err(); err != nil {
		return fmt.Errorf("cannot save info to cache: %w", err)
	}

	return nil
}

func (m *runInfoManager) GetEngineQueues(ctx context.Context) ([]RunInfo, error) {
	recentInfos := make(map[string]InstanceRunInfo)
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
				if s, ok := v.(string); ok {
					var info InstanceRunInfo
					err := json.Unmarshal([]byte(s), &info)
					if err != nil {
						return nil, fmt.Errorf("cannot unmarshal info key=%q: %w", unprocessedKeys[i], err)
					}

					if recentInfo, ok := recentInfos[info.Name]; !ok || recentInfo.Time.Before(info.Time) {
						recentInfos[info.Name] = info
					}
				} else {
					return nil, fmt.Errorf("expect string for key=%q but got type=%T", unprocessedKeys[i], v)
				}
			}
		}

		if cursor == 0 {
			break
		}
	}

	infos := make([]RunInfo, len(recentInfos))
	i := 0
	for _, info := range recentInfos {
		infos[i] = RunInfo{
			Name:         info.Name,
			ConsumeQueue: info.ConsumeQueue,
			PublishQueue: info.PublishQueue,
		}
		i++
	}

	return infos, nil
}

func (m *runInfoManager) GetCacheKey(info InstanceRunInfo) string {
	return strings.Join([]string{m.key, info.Name, info.ID}, libredis.KeyDelimiter)
}
