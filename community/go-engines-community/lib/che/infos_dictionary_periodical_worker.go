package che

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

const minInfoLength = 2

func NewInfosDictionaryPeriodicalWorker(
	client mongo.DbClient,
	periodicalInterval time.Duration,
	logger zerolog.Logger,
) engine.PeriodicalWorker {
	return &infosDictionaryPeriodicalWorker{
		Client:             client,
		PeriodicalInterval: periodicalInterval,
		Logger:             logger,
	}
}

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
						"$type": "string",
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
			{
				"$out": mongo.InfosDictionaryCollection,
			},
		},
	)
	if err != nil {
		w.Logger.Error().Err(err).Msg("unable to load entity infos data")
		return
	}

	defer cursor.Close(ctx)
}
