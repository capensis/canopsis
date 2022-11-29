package che

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"

	"github.com/rs/zerolog"
)

type softDeletePeriodicalWorker struct {
	collection         mongo.DbCollection
	PeriodicalInterval time.Duration
	Logger             zerolog.Logger
}

func (w *softDeletePeriodicalWorker) GetInterval() time.Duration {
	return w.PeriodicalInterval
}

func (w *softDeletePeriodicalWorker) Work(ctx context.Context) {
	cursor, err := w.collection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"soft_deleted": true,
				},
			},
			{
				"$lookup": bson.M{
					"from": mongo.AlarmMongoCollection,
					"let":  bson.M{"id": "$_id"},
					"pipeline": []bson.M{
						{
							"$match": bson.M{"$and": []bson.M{
								{"$expr": bson.M{"$eq": bson.A{"$d", "$$id"}}},
								{"v.resolved": nil},
							}},
						},
						{"$limit": 1},
					},
					"as": "alarm",
				},
			},
			{
				"$match": bson.M{
					"alarm": bson.M{
						"$size": 0,
					},
				},
			},
		},
	)
	if err != nil {
		w.Logger.Error().Err(err).Msg("unable to load soft deleted entities")
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var ent struct {
			ID string `bson:"_id"`
		}

		err = cursor.Decode(&ent)
		if err != nil {
			w.Logger.Error().Err(err).Msg("unable to decode an entity")
		}

		_, err = w.collection.DeleteOne(ctx, bson.M{"_id": ent.ID, "soft_deleted": true})
		if err != nil {
			w.Logger.Error().Err(err).Msg("unable to decode an entity")
		}
	}
}
