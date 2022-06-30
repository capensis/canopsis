package che

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"time"
)

const minInfoLength = 2

type infosDictionaryPeriodicalWorker struct {
	Client             mongo.DbClient
	PeriodicalInterval time.Duration
	Logger             zerolog.Logger
}

func (w *infosDictionaryPeriodicalWorker) GetInterval() time.Duration {
	return w.PeriodicalInterval
}

func (w *infosDictionaryPeriodicalWorker) Work(ctx context.Context) {
	cursor, err := w.Client.Collection(mongo.EntityMongoCollection).Aggregate(
		ctx,
		[]bson.M{
			{
				"$project": bson.M{
					"infos": bson.M{
						"$map": bson.M{
							"input": bson.M{
								"$objectToArray": "$infos",
							},
							"in": "$$this",
						},
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
				"$match": bson.M{
					"infos.v.value": bson.M{
						"$type": bson.A{"string"},
					},
				},
			},
			{
				"$addFields": bson.M{
					"valueLen": bson.M{
						"$strLenCP": "$infos.v.value",
					},
				},
			},
			{
				"$match": bson.M{
					"valueLen": bson.M{
						"$gt": minInfoLength,
					},
				},
			},
			{
				"$group": bson.M{
					"_id": "$infos.k",
					"values": bson.M{
						"$addToSet": "$infos.v.value",
					},
				},
			},
		},
	)
	if err != nil {
		w.Logger.Error().Err(err).Msg("unable to load entity infos data")
		return
	}

	defer cursor.Close(ctx)

	writeModels := make([]mongodriver.WriteModel, 0, canopsis.DefaultBulkSize)
	bulkBytesSize := 0
	var ids []string

	for cursor.Next(ctx) {
		var info struct {
			ID     string        `bson:"_id"`
			Values []interface{} `bson:"values"`
		}

		err = cursor.Decode(&info)
		if err != nil {
			w.Logger.Error().Err(err).Msg("unable to decode entity infos data")
			return
		}

		ids = append(ids, info.ID)

		newModel := mongodriver.
			NewUpdateOneModel().
			SetFilter(bson.M{"_id": info.ID}).
			SetUpdate(bson.M{"$set": info}).
			SetUpsert(true)

		b, err := bson.Marshal(newModel)
		if err != nil {
			w.Logger.Error().Err(err).Msg("unable to marshal entity infos data model")
			return
		}

		newModelLen := len(b)
		if bulkBytesSize+newModelLen > canopsis.DefaultBulkBytesSize {
			_, err := w.Client.Collection(mongo.InfosDictionaryCollection).BulkWrite(ctx, writeModels)
			if err != nil {
				w.Logger.Error().Err(err).Msg("unable to bulk write entity infos")
				return
			}

			writeModels = writeModels[:0]
			bulkBytesSize = 0
		}

		bulkBytesSize += newModelLen
		writeModels = append(writeModels, newModel)

		if len(writeModels) == canopsis.DefaultBulkSize {
			_, err := w.Client.Collection(mongo.InfosDictionaryCollection).BulkWrite(ctx, writeModels)
			if err != nil {
				w.Logger.Error().Err(err).Msg("unable to bulk write entity infos")
				return
			}

			writeModels = writeModels[:0]
			bulkBytesSize = 0
		}
	}

	if len(writeModels) > 0 {
		_, err := w.Client.Collection(mongo.InfosDictionaryCollection).BulkWrite(ctx, writeModels)
		if err != nil {
			w.Logger.Error().Err(err).Msg("unable to bulk write entity infos")
			return
		}
	}

	_, err = w.Client.Collection(mongo.InfosDictionaryCollection).DeleteMany(ctx, bson.M{"_id": bson.M{"$nin": ids}})
	if err != nil {
		w.Logger.Error().Err(err).Msg("unable to delete entity infos")
		return
	}
}
