package che

import (
	"context"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	minInfoLength   = 2
	batchSize       = 100000
	buildingTimeout = 5 * time.Minute
)

// A composite id is used, because it works faster with a lot of bulk upserts instead of filter and uuid
type infosDictID struct {
	Key   string `bson:"k"`
	Value string `bson:"v"`
}

type infosDictDoc struct {
	ID infosDictID `bson:"_id"`
}

func NewInfosDictionaryPeriodicalWorker(
	client mongo.DbClient,
	periodicalInterval time.Duration,
	logger zerolog.Logger,
) engine.PeriodicalWorker {
	return &infosDictionaryPeriodicalWorker{
		entityCollection:          client.Collection(mongo.EntityMongoCollection),
		entityInfosDictCollection: client.Collection(mongo.EntityInfosDictionaryCollection),
		configCollection:          client.Collection(mongo.ConfigurationMongoCollection),
		periodicalInterval:        periodicalInterval,
		logger:                    logger,
	}
}

type infosDictionaryPeriodicalWorker struct {
	entityCollection          mongo.DbCollection
	entityInfosDictCollection mongo.DbCollection
	configCollection          mongo.DbCollection
	periodicalInterval        time.Duration
	logger                    zerolog.Logger
}

func (w *infosDictionaryPeriodicalWorker) GetInterval() time.Duration {
	return w.periodicalInterval
}

func (w *infosDictionaryPeriodicalWorker) Work(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, buildingTimeout)
	defer cancel()

	conf := config.CanopsisConf{}

	err := w.configCollection.FindOne(ctx, bson.M{"_id": config.ConfigKeyName}).Decode(&conf)
	if err != nil {
		w.logger.Error().Err(err).Msg("unable to get canopsis config")
		return
	}

	if !conf.Global.BuildEntityInfosDictionary {
		return
	}

	now := types.NewCpsTime()

	entCursor, err := w.entityCollection.Find(
		ctx,
		bson.M{},
		options.Find().SetSort(bson.M{"_id": 1}).SetProjection(bson.M{"_id": 1}),
	)
	if err != nil {
		w.logger.Error().Err(err).Msg("unable to load entity infos data")
		return
	}

	defer entCursor.Close(ctx)

	entIds := make([]string, 0, batchSize)

	for entCursor.Next(ctx) {
		var doc struct {
			ID string `bson:"_id"`
		}

		err = entCursor.Decode(&doc)
		if err != nil {
			w.logger.Error().Err(err).Msg("unable to load entity infos data")
			return
		}

		entIds = append(entIds, doc.ID)
		if len(entIds) == batchSize {
			err = w.buildDictionary(ctx, entIds, now)
			if err != nil {
				w.logger.Error().Err(err).Msg("failed to build entity infos dictionary")
				return
			}

			entIds = entIds[:0]
		}
	}

	if len(entIds) > 0 {
		err = w.buildDictionary(ctx, entIds, now)
		if err != nil {
			w.logger.Error().Err(err).Msg("failed to build entity infos dictionary")
			return
		}
	}

	delCursor, err := w.entityInfosDictCollection.Find(ctx, bson.M{
		"$or": bson.A{
			bson.M{"last_update": bson.M{"$lt": now}},
			bson.M{"last_update": bson.M{"$exists": false}},
		}},
	)
	if err != nil {
		w.logger.Error().Err(err).Msg("unable to find outdated entity infos dictionary documents")
		return
	}

	defer delCursor.Close(ctx)

	// use any, because old dictionary ids contained strings instead of composite ids
	ids := make([]any, 0, canopsis.DefaultBulkSize)

	for delCursor.Next(ctx) {
		var info struct {
			ID any `bson:"_id"`
		}

		err = delCursor.Decode(&info)
		if err != nil {
			w.logger.Error().Err(err).Msg("unable to decode entity infos data")
			return
		}

		ids = append(ids, info.ID)

		if len(ids) == canopsis.DefaultBulkSize {
			_, err = w.entityInfosDictCollection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})
			if err != nil {
				w.logger.Error().Err(err).Msg("unable to delete outdated entity infos dictionary documents")
				return
			}

			ids = ids[:0]
		}
	}

	if len(ids) > 0 {
		_, err = w.entityInfosDictCollection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})
		if err != nil {
			w.logger.Error().Err(err).Msg("unable to delete outdated entity infos dictionary documents")
			return
		}
	}
}

func (w *infosDictionaryPeriodicalWorker) buildDictionary(ctx context.Context, entIds []string, t types.CpsTime) error {
	cursor, err := w.entityCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{"enabled": true, "_id": bson.M{"$in": entIds}},
			},
			{
				"$project": bson.M{
					"infos": bson.M{
						"$objectToArray": "$infos",
					},
				},
			},
			{
				"$unwind": "$infos",
			},
			{
				"$unwind": "$infos.v.value",
			},
			{
				"$addFields": bson.M{
					"valueLen": bson.M{
						"$cond": bson.M{
							"if":   bson.M{"$eq": bson.A{bson.M{"$type": "$infos.v.value"}, "string"}},
							"then": bson.M{"$strLenCP": "$infos.v.value"},
							"else": 0,
						},
					},
				},
			},
			{
				"$project": bson.M{
					"k": "$infos.k",
					"v": bson.M{
						"$cond": bson.M{
							"if":   bson.M{"$gt": bson.A{"$valueLen", minInfoLength}},
							"then": "$infos.v.value",
							"else": "$$REMOVE",
						},
					},
				},
			},
			{
				"$group": bson.M{
					"_id": bson.M{
						"k": "$k",
						"v": "$v",
					},
				},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("unable to load entity infos data: %w", err)
	}

	defer cursor.Close(ctx)

	writeModels := make([]mongodriver.WriteModel, 0, canopsis.DefaultBulkSize)
	bulkBytesSize := 0

	for cursor.Next(ctx) {
		var info infosDictDoc

		err = cursor.Decode(&info)
		if err != nil {
			return fmt.Errorf("unable to decode entity infos data: %w", err)
		}

		newModel := mongodriver.
			NewUpdateOneModel().
			SetFilter(bson.M{"_id": info.ID}).
			SetUpdate(bson.M{"$set": bson.M{"last_update": t}}).
			SetUpsert(true)

		b, err := bson.Marshal(newModel)
		if err != nil {
			return fmt.Errorf("unable to marshal entity infos data: %w", err)
		}

		newModelLen := len(b)
		if bulkBytesSize+newModelLen > canopsis.DefaultBulkBytesSize {
			_, err = w.entityInfosDictCollection.BulkWrite(ctx, writeModels)
			if err != nil {
				return fmt.Errorf("unable to bulk write entity infos dictionary: %w", err)
			}

			writeModels = writeModels[:0]
			bulkBytesSize = 0
		}

		bulkBytesSize += newModelLen
		writeModels = append(writeModels, newModel)

		if len(writeModels) == canopsis.DefaultBulkSize {
			_, err = w.entityInfosDictCollection.BulkWrite(ctx, writeModels)
			if err != nil {
				return fmt.Errorf("unable to bulk write entity infos dictionary: %w", err)
			}

			writeModels = writeModels[:0]
			bulkBytesSize = 0
		}
	}

	if len(writeModels) > 0 {
		_, err = w.entityInfosDictCollection.BulkWrite(ctx, writeModels)
		if err != nil {
			return fmt.Errorf("unable to bulk write entity infos dictionary: %w", err)
		}
	}

	return nil
}
