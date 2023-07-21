package che

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

const minInfoLength = 2

func NewInfosDictionaryPeriodicalWorker(
	client mongo.DbClient,
	periodicalInterval time.Duration,
	logger zerolog.Logger,
) engine.PeriodicalWorker {
	return &infosDictionaryPeriodicalWorker{
		entityCollection:          client.Collection(mongo.EntityMongoCollection),
		entityInfosDictCollection: client.Collection(mongo.EntityInfosDictionaryCollection),
		PeriodicalInterval:        periodicalInterval,
		Logger:                    logger,
	}
}

type infosDictionaryPeriodicalWorker struct {
	entityCollection          mongo.DbCollection
	entityInfosDictCollection mongo.DbCollection
	PeriodicalInterval        time.Duration
	Logger                    zerolog.Logger
}

func (w *infosDictionaryPeriodicalWorker) GetInterval() time.Duration {
	return w.PeriodicalInterval
}

func (w *infosDictionaryPeriodicalWorker) Work(ctx context.Context) {
	now := types.NewCpsTime()

	dictCursor, err := w.entityCollection.Aggregate(
		ctx,
		[]bson.M{
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
			{
				"$project": bson.M{
					"_id": 0,
					"k":   "$_id.k",
					"v":   "$_id.v",
				},
			},
		},
	)
	if err != nil {
		w.Logger.Error().Err(err).Msg("unable to load entity infos data")
		return
	}

	defer dictCursor.Close(ctx)

	writeModels := make([]mongodriver.WriteModel, 0, canopsis.DefaultBulkSize)
	bulkBytesSize := 0

	for dictCursor.Next(ctx) {
		var info struct {
			Key   string `bson:"k"`
			Value string `bson:"v"`
		}

		err = dictCursor.Decode(&info)
		if err != nil {
			w.Logger.Error().Err(err).Msg("unable to decode entity infos data")
			return
		}

		newModel := mongodriver.
			NewUpdateOneModel().
			SetFilter(bson.M{"k": info.Key, "v": info.Value}).
			SetUpdate(bson.M{
				"$set":         bson.M{"last_update": now},
				"$setOnInsert": bson.M{"_id": utils.NewID()},
			}).
			SetUpsert(true)

		b, err := bson.Marshal(newModel)
		if err != nil {
			w.Logger.Error().Err(err).Msg("unable to marshal entity infos data")
			return
		}

		newModelLen := len(b)
		if bulkBytesSize+newModelLen > canopsis.DefaultBulkBytesSize {
			_, err := w.entityInfosDictCollection.BulkWrite(ctx, writeModels)
			if err != nil {
				w.Logger.Error().Err(err).Msg("unable to bulk write entity infos dictionary")
				return
			}

			writeModels = writeModels[:0]
			bulkBytesSize = 0
		}

		bulkBytesSize += newModelLen
		writeModels = append(writeModels, newModel)

		if len(writeModels) == canopsis.DefaultBulkSize {
			_, err := w.entityInfosDictCollection.BulkWrite(ctx, writeModels)
			if err != nil {
				w.Logger.Error().Err(err).Msg("unable to bulk write entity infos dictionary")
				return
			}

			writeModels = writeModels[:0]
			bulkBytesSize = 0
		}
	}

	if len(writeModels) > 0 {
		_, err := w.entityInfosDictCollection.BulkWrite(ctx, writeModels)
		if err != nil {
			w.Logger.Error().Err(err).Msg("unable to bulk write entity infos dictionary")
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
		w.Logger.Error().Err(err).Msg("unable to find outdated entity infos dictionary documents")
		return
	}

	defer delCursor.Close(ctx)

	ids := make([]string, 0, canopsis.DefaultBulkSize)

	for delCursor.Next(ctx) {
		var info struct {
			ID string `bson:"_id"`
		}

		err = delCursor.Decode(&info)
		if err != nil {
			w.Logger.Error().Err(err).Msg("unable to decode entity infos data")
			return
		}

		ids = append(ids, info.ID)

		if len(ids) == canopsis.DefaultBulkSize {
			_, err = w.entityInfosDictCollection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})
			if err != nil {
				w.Logger.Error().Err(err).Msg("unable to delete outdated entity infos dictionary documents")
				return
			}

			ids = ids[:0]
		}
	}

	if len(ids) > 0 {
		_, err = w.entityInfosDictCollection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})
		if err != nil {
			w.Logger.Error().Err(err).Msg("unable to delete outdated entity infos dictionary documents")
			return
		}
	}
}
