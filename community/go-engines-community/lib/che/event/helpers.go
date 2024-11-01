package event

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

func updateMetaAlarmInfos(
	ctx context.Context,
	entityID string,
	newInfos map[string]types.Info,
	alarmCollection, entityCollection mongo.DbCollection,
) error {
	if len(newInfos) == 0 {
		return nil
	}

	cursor, err := alarmCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"d":          entityID,
			"v.resolved": nil,
			"v.parents":  bson.M{"$nin": bson.A{nil, bson.A{}}},
		}},
		{"$lookup": bson.M{
			"from":         mongo.AlarmMongoCollection,
			"localField":   "v.parents",
			"foreignField": "d",
			"as":           "parents",
			"pipeline": []bson.M{
				{"$match": bson.M{
					"v.resolved": nil,
					"cinfos":     bson.M{"$nin": bson.A{nil, bson.A{}}},
				}},
			},
		}},
		{"$unwind": "$parents"},
		{"$project": bson.M{
			"_id":    "$parents.d",
			"cinfos": "$parents.cinfos",
		}},
	})
	if err != nil {
		return fmt.Errorf("cannot find parents: %w", err)
	}

	defer cursor.Close(ctx)
	writeModels := make([]mongodriver.WriteModel, 0)
	for cursor.Next(ctx) {
		parent := struct {
			ID                      string   `bson:"_id"`
			EntityInfosFromChildren []string `bson:"cinfos"`
		}{}
		err = cursor.Decode(&parent)
		if err != nil {
			return fmt.Errorf("cannot decode parent: %w", err)
		}

		update := bson.M{}
		for _, infoName := range parent.EntityInfosFromChildren {
			if info, ok := newInfos[infoName]; ok {
				update["infos."+infoName+".value"] = info.Value
			}
		}

		if len(update) > 0 {
			writeModels = append(writeModels, mongodriver.NewUpdateOneModel().
				SetFilter(bson.M{"_id": parent.ID}).
				SetUpdate(bson.M{"$set": update}))
		}
	}

	if len(writeModels) > 0 {
		_, err = entityCollection.BulkWrite(ctx, writeModels)
		if err != nil {
			return fmt.Errorf("cannot update parent infos: %w", err)
		}
	}

	return nil
}
