package axe

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"time"
)

type entityServiceStateWorker struct {
	pipeline                  []bson.M
	periodicalInterval        time.Duration
	dbClient                  mongo.DbClient
	entityCollection          mongo.DbCollection
	serviceCountersCollection mongo.DbCollection
	logger                    zerolog.Logger
}

func NewEntityServiceStateWorker(
	dbClient mongo.DbClient,
	logger zerolog.Logger,
	interval time.Duration,
) engine.PeriodicalWorker {
	return &entityServiceStateWorker{
		entityCollection:          dbClient.Collection(mongo.EntityMongoCollection),
		serviceCountersCollection: dbClient.Collection(mongo.EntityServiceCountersMongoCollection),
		periodicalInterval:        interval,
		dbClient:                  dbClient,
		logger:                    logger,
		pipeline: []bson.M{
			{
				"$match": bson.M{
					"$or": []bson.M{
						{
							"$and": []bson.M{
								{
									"impacted_services_to_add": bson.M{"$exists": true},
								},
								{
									"$expr": bson.M{"$gt": bson.A{bson.M{"$size": "$impacted_services_to_add"}, 0}},
								},
							},
						},
						{
							"$and": []bson.M{
								{
									"impacted_services_to_remove": bson.M{"$exists": true},
								},
								{
									"$expr": bson.M{"$gt": bson.A{bson.M{"$size": "$impacted_services_to_remove"}, 0}},
								},
							},
						},
					},
					"last_event_date": bson.M{"$lt": time.Now().Unix() - 60},
				},
			},
			{
				"$project": bson.M{
					"entity": "$$ROOT",
					"_id":    0,
				},
			},
			{
				"$lookup": bson.M{
					"from":         mongo.AlarmMongoCollection,
					"localField":   "entity._id",
					"foreignField": "d",
					"as":           "alarm",
				},
			},
			{
				"$unwind": bson.M{
					"path":                       "$alarm",
					"preserveNullAndEmptyArrays": true,
				},
			},
		},
	}
}

func (w *entityServiceStateWorker) GetInterval() time.Duration {
	return w.periodicalInterval
}

func (w *entityServiceStateWorker) Work(ctx context.Context) {
	fmt.Printf("tick\n")
	err := w.dbClient.WithTransaction(ctx, func(tCtx context.Context) error {
		cursor, err := w.entityCollection.Aggregate(tCtx, w.pipeline)
		if err != nil {
			return err
		}

		defer cursor.Close(tCtx)

		updateCountersWriteModels := make([]mongodriver.WriteModel, 0, canopsis.DefaultBulkSize)
		updateEntitiesWriteModels := make([]mongodriver.WriteModel, 0, canopsis.DefaultBulkSize)

		for cursor.Next(tCtx) {
			fmt.Printf("found\n")
			var depEnt types.AlarmWithEntity
			err := cursor.Decode(&depEnt)
			if err != nil {
				return err
			}

			state := int(depEnt.Alarm.CurrentState())

			for _, impServ := range depEnt.Entity.ImpactedServicesToRemove {
				updateCountersWriteModels = append(
					updateCountersWriteModels,
					mongodriver.
						NewUpdateOneModel().
						SetFilter(bson.M{"_id": impServ}).
						SetUpdate(bson.M{
							"$inc": bson.M{common.StateTitles[state]: -1},
						}).
						SetUpsert(true))
			}

			for _, impServ := range depEnt.Entity.ImpactedServicesToAdd {
				updateCountersWriteModels = append(
					updateCountersWriteModels,
					mongodriver.
						NewUpdateOneModel().
						SetFilter(bson.M{"_id": impServ}).
						SetUpdate(bson.M{
							"$inc": bson.M{common.StateTitles[state]: 1},
						}).
						SetUpsert(true))
			}

			updateEntitiesWriteModels = append(
				updateEntitiesWriteModels,
				mongodriver.
					NewUpdateOneModel().
					SetFilter(bson.M{"_id": depEnt.Entity.ID}).
					SetUpdate(bson.M{
						"$unset": bson.M{
							"impacted_services_to_add":    1,
							"impacted_services_to_remove": 1,
						},
					}),
			)

			if len(updateCountersWriteModels) >= canopsis.DefaultBulkSize {
				_, err := w.serviceCountersCollection.BulkWrite(ctx, updateCountersWriteModels)
				if err != nil {
					return err
				}

				updateCountersWriteModels = updateCountersWriteModels[:0]
			}

			if len(updateEntitiesWriteModels) >= canopsis.DefaultBulkSize {
				_, err := w.entityCollection.BulkWrite(ctx, updateEntitiesWriteModels)
				if err != nil {
					return err
				}

				updateEntitiesWriteModels = updateEntitiesWriteModels[:0]
			}
		}

		if len(updateCountersWriteModels) > 0 {
			_, err := w.serviceCountersCollection.BulkWrite(ctx, updateCountersWriteModels)
			if err != nil {
				return err
			}
		}

		if len(updateEntitiesWriteModels) > 0 {
			_, err := w.entityCollection.BulkWrite(ctx, updateEntitiesWriteModels)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		w.logger.Err(err).Msg("cannot load alarm status rules")
	}
}
