package engine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

const defaultKey = "engine-run-info"

// NewRunInfoManager creates new run info manager.
func NewRunInfoManager(client redis.Cmdable, key ...string) RunInfoManager {
	k := defaultKey
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

func (m *runInfoManager) Save(ctx context.Context, info RunInfo, expiration time.Duration) error {
	if info.Name == "" {
		return errors.New("name is required")
	}

	b, err := json.Marshal(info)
	if err != nil {
		return err
	}

	res := m.client.Set(ctx, m.getKey(info.Name), b, expiration)
	if err := res.Err(); err != nil {
		return err
	}

	return nil
}

func (m *runInfoManager) Get(ctx context.Context, engineName string) (*RunInfo, error) {
	if engineName == "" {
		return nil, errors.New("name is required")
	}

	res := m.client.Get(ctx, m.getKey(engineName))
	if err := res.Err(); err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var info RunInfo
	err := json.Unmarshal([]byte(res.Val()), &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (m *runInfoManager) GetAll(ctx context.Context) ([]RunInfo, error) {
	infos := make([]RunInfo, 0)
	var cursor uint64
	processedKeys := make(map[string]bool)

	for {
		res := m.client.Scan(ctx, cursor, fmt.Sprintf("%s*", m.key), 50)
		if err := res.Err(); err != nil {
			return nil, err
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
				return nil, err
			}

			for _, v := range resGet.Val() {
				if s, ok := v.(string); ok {
					var info RunInfo
					err := json.Unmarshal([]byte(s), &info)
					if err != nil {
						return nil, err
					}

					infos = append(infos, info)
				} else {
					return nil, fmt.Errorf("unknown value type")
				}
			}
		}

		if cursor == 0 {
			break
		}
	}

	return infos, nil
}

func (m *runInfoManager) GetGraph(ctx context.Context) (*RunInfoGraph, error) {
	infos, err := m.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	graph := RunInfoGraph{
		Nodes: infos,
	}
	for i := range infos {
		for j := range infos {
			if infos[i].PublishQueue == infos[j].ConsumeQueue {
				graph.Edges = append(graph.Edges, Edge{
					From: infos[i].Name,
					To:   infos[j].Name,
				})
			}
		}
	}

	return &graph, nil
}

func (m *runInfoManager) ClearAll(ctx context.Context) error {
	var cursor uint64
	for {
		res := m.client.Scan(ctx, cursor, fmt.Sprintf("%s*", m.key), 50)
		if err := res.Err(); err != nil {
			return err
		}

		var keys []string
		keys, cursor = res.Val()
		if len(keys) > 0 {
			resDel := m.client.Del(ctx, keys...)
			if err := resDel.Err(); err != nil {
				return err
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}

func (m *runInfoManager) getKey(name string) string {
	return fmt.Sprintf("%s[%s]", m.key, name)
}
