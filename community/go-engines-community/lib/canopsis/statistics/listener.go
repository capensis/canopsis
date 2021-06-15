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

	l.logger.Debug().Msg("flush")
	operations := make(map[string]map[int64]mongodriver.WriteModel)
	minIds := make(map[string]int64)
	maxIds := make(map[string]int64)
	for collection := range l.storeIntervals {
		operations[collection] = make(map[int64]mongodriver.WriteModel)
	}

	var cursor uint64
	processedKeys := make(map[string]bool)

	for {
		res := l.redisClient.Scan(ctx, cursor, "*", 60)
		if err := res.Err(); err != nil {
			l.logger.Error().Err(err).Msg("Failed to flush statistics: failed to get data from redis")
			return
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

		for _, key := range unprocessedKeys {
			var (
				result                    *redis.StringCmd
				minute, received, dropped int
				err                       error
			)

			result = l.redisClient.HGet(ctx, key, "received")
			if err := result.Err(); err != nil {
				l.logger.Error().Err(result.Err()).Str("redis_key", key).Msg("Failed to flush statistics: failed to get received value from redis")
				break
			}

			received, err = strconv.Atoi(result.Val())
			if err != nil {
				l.logger.Error().Err(err).Str("redis_key", key).Msg("Failed to flush statistics: failed to convert redis value to int")
				break
			}

			result = l.redisClient.HGet(ctx, key, "dropped")
			if err := result.Err(); err != nil {
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
				operations[collection][id] = mongodriver.NewUpdateOneModel().
					SetFilter(bson.M{"_id": id}).
					SetUpdate(bson.M{"$inc": bson.M{"received": received, "dropped": dropped}}).
					SetUpsert(true)
				if minIds[collection] == 0 || minIds[collection] > id {
					minIds[collection] = id
				}
				if maxIds[collection] == 0 || maxIds[collection] < id {
					maxIds[collection] = id
				}
			}
		}

		if cursor == 0 {
			break
		}
	}

	for collectionName, writeModelsByID := range operations {
		if len(writeModelsByID) == 0 {
			continue
		}

		collection := l.mongoClient.Collection(collectionName)
		interval := l.storeIntervals[collectionName] * secondsInMinute
		lastSavedID, err := getLastID(ctx, collection)

		if err != nil {
			l.logger.Error().Err(err).Msg("Failed to flush statistics: failed to fetch last id")
			return
		}

		writeModels := fillEmptyIds(writeModelsByID, lastSavedID, minIds[collectionName],
			maxIds[collectionName], interval)
		_, err = collection.BulkWrite(ctx, writeModels)
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

// fillEmptyIds adds write models with empty counters if there are ids between minID and maxID
// or between lastSavedID and maxID.
func fillEmptyIds(
	writeModelsByID map[int64]mongodriver.WriteModel,
	lastSavedID, minID, maxID, interval int64,
) []mongodriver.WriteModel {
	writeModels := make([]mongodriver.WriteModel, 0)
	var id int64
	if lastSavedID == 0 || lastSavedID > minID {
		id = minID
	} else {
		id = lastSavedID
	}

	for ; id <= maxID; id += interval {
		if writeModel, ok := writeModelsByID[id]; ok {
			writeModels = append(writeModels, writeModel)
		} else {
			writeModels = append(writeModels, mongodriver.NewUpdateOneModel().
				SetFilter(bson.M{"_id": id}).
				SetUpdate(bson.M{"$setOnInsert": bson.M{"_id": id, "received": 0, "dropped": 0}}).
				SetUpsert(true))
		}
	}

	return writeModels
}

// getLastID fetches max id from collection.
func getLastID(ctx context.Context, collection mongo.DbCollection) (int64, error) {
	cursor, err := collection.Aggregate(ctx, []bson.M{{"$group": bson.M{
		"_id": nil,
		"max": bson.M{"$max": "$_id"},
		"min": bson.M{"$min": "$_id"},
	}}})
	if err != nil {
		return 0, err
	}

	defer cursor.Close(ctx)

	res := struct {
		Max int64 `bson:"max"`
	}{}

	if cursor.Next(ctx) {
		err := cursor.Decode(&res)
		if err != nil {
			return 0, err
		}
	}

	return res.Max, nil
}
