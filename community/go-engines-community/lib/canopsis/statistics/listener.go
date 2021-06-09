package statistics

import (
	"context"
	"strconv"
	"sync"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const secondsInMinute = 60
const redisTimeout = 5 * time.Second

type statsListener struct {
	mongoClient   mongo.DbClient
	redisMx       sync.Mutex
	redisClient   redis.Cmdable
	flushInterval time.Duration
	// storeIntervals is collection => minutes map.
	storeIntervals map[string]int64
	statsMx        sync.Mutex
	stats          map[int]counts
	logger         zerolog.Logger
}

type counts struct {
	Received int64
	Dropped  int64
}

func NewStatsListener(
	mongoClient mongo.DbClient,
	redisClient redis.Cmdable,
	flushInterval time.Duration,
	storeIntervals map[string]int64,
	logger zerolog.Logger,
) StatsListener {
	return &statsListener{
		mongoClient:    mongoClient,
		redisClient:    redisClient,
		flushInterval:  flushInterval,
		storeIntervals: storeIntervals,
		stats:          make(map[int]counts),
		logger:         logger,
	}
}

func (l *statsListener) Listen(ctx context.Context, channel <-chan Message) {
	tickerMongo := time.NewTicker(l.flushInterval)
	defer tickerMongo.Stop()
	tickerRedis := time.NewTicker(redisTimeout)
	defer tickerRedis.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case m, ok := <-channel:
			if !ok {
				return
			}
			l.save(m)
		case <-tickerRedis.C:
			l.saveToRedis(ctx)
		case <-tickerMongo.C:
			l.saveToDB(ctx)
		}
	}
}

func (l *statsListener) save(m Message) {
	l.statsMx.Lock()
	defer l.statsMx.Unlock()

	var received, dropped int64
	ts := m.Timestamp
	minute := int(ts) / secondsInMinute * secondsInMinute

	if v, ok := l.stats[minute]; ok {
		received = v.Received
		dropped = v.Dropped
	}

	received += m.Received
	dropped += m.Dropped
	l.stats[minute] = counts{
		Received: received,
		Dropped:  dropped,
	}
}

func (l *statsListener) saveToRedis(ctx context.Context) {
	l.statsMx.Lock()
	l.redisMx.Lock()
	defer l.statsMx.Unlock()
	defer l.redisMx.Unlock()
	for minute, counts := range l.stats {
		key := strconv.Itoa(minute)
		result := l.redisClient.HIncrBy(ctx, key, "received", counts.Received)
		if result.Err() != nil {
			l.logger.Error().Err(result.Err()).Str("redis_key", key).Int64("value", counts.Received).Msg("Failed to save statistics in redis")
			return
		}

		result = l.redisClient.HIncrBy(ctx, key, "dropped", counts.Dropped)
		if result.Err() != nil {
			l.logger.Error().Err(result.Err()).Str("redis_key", key).Int64("value", counts.Dropped).Msg("Failed to save statistics in redis")
			return
		}
	}

	l.stats = make(map[int]counts)
}

func (l *statsListener) saveToDB(ctx context.Context) {
	l.redisMx.Lock()
	defer l.redisMx.Unlock()

	var err error
	l.logger.Debug().Msg("flush")
	keysResult := l.redisClient.Keys(ctx, "*")
	if keysResult.Err() != nil {
		l.logger.Error().Err(keysResult.Err()).Msg("Failed to flush statistics: failed to get data from redis")
		return
	}

	operations := make(map[string][]mongodriver.WriteModel)

	for _, key := range keysResult.Val() {
		var (
			result                    *redis.StringCmd
			minute, received, dropped int
		)

		result = l.redisClient.HGet(ctx, key, "received")
		if err = result.Err(); err != nil {
			l.logger.Error().Err(result.Err()).Str("redis_key", key).Msg("Failed to flush statistics: failed to get received value from redis")
			break
		}

		received, err = strconv.Atoi(result.Val())
		if err != nil {
			l.logger.Error().Err(err).Str("redis_key", key).Msg("Failed to flush statistics: failed to convert redis value to int")
			break
		}

		result = l.redisClient.HGet(ctx, key, "dropped")
		if err = result.Err(); err != nil {
			l.logger.Error().Err(result.Err()).Str("redis_key", key).Msg("Failed to flush statistics: failed to get dropped value from redis")
			break
		}

		dropped, err = strconv.Atoi(result.Val())
		if err != nil {
			l.logger.Error().Err(err).Str("redis_key", key).Msg("Failed to flush statistics: failed to convert redis value to int")
			break
		}

		minute, err = strconv.Atoi(key)
		if err != nil {
			l.logger.Error().Err(err).Str("redis_key", key).Msg("Failed to flush statistics: failed to convert redis value to int")
			break
		}

		for collection, interval := range l.storeIntervals {
			id := int64(minute) / (interval * secondsInMinute) * (interval * secondsInMinute)
			operation := mongodriver.NewUpdateOneModel()
			operation.SetFilter(bson.M{"_id": id})
			operation.SetUpdate(bson.M{"$inc": bson.M{"received": received, "dropped": dropped}})
			operation.SetUpsert(true)

			operations[collection] = append(operations[collection], operation)
		}
	}

	if err != nil {
		l.logger.Error().Err(err).Msg("Failed to flush statistics: Skip flush")
		return
	}

	for collection, writeModels := range operations {
		if len(writeModels) == 0 {
			continue
		}

		collection := l.mongoClient.Collection(collection)
		_, err = collection.BulkWrite(ctx, writeModels, &options.BulkWriteOptions{})
		if err != nil {
			l.logger.Error().Err(err).Msg("Failed to flush statistics: failed to save statistics to mongo")
			return
		}
	}

	flushAllResult := l.redisClient.FlushDB(ctx)
	if flushAllResult.Err() != nil {
		l.logger.Error().Err(flushAllResult.Err()).Msg("Failed to flush statistics: failed to remove old data from redis")
		return
	}
}
