package priority

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UpdateFollowing(
	ctx context.Context,
	collection mongo.DbCollection,
	id string,
	priority int64,
) error {
	if priority <= 0 {
		return nil
	}

	cursor, err := collection.Find(ctx, bson.M{
		"_id":      bson.M{"$ne": id},
		"priority": bson.M{"$gte": priority},
	}, options.Find().SetSort(bson.M{"priority": 1}))
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)

	writeModels := make([]mongodriver.WriteModel, 0)
	seq := priority
	for cursor.Next(ctx) {
		item := struct {
			ID       string `bson:"_id"`
			Priority int64  `bson:"priority"`
		}{}
		err = cursor.Decode(&item)
		if err != nil {
			return err
		}

		// Do nothing if there is no conflict
		if seq == priority && item.Priority != priority {
			return nil
		}
		// Do not update priorities after gap
		if seq != item.Priority {
			break
		}

		seq++
		writeModels = append(writeModels, mongodriver.NewUpdateOneModel().
			SetFilter(bson.M{"_id": item.ID}).
			SetUpdate(bson.M{"$set": bson.M{"priority": seq}}))
	}

	if len(writeModels) > 0 {
		_, err = collection.BulkWrite(ctx, writeModels)
		if err != nil {
			return err
		}
	}

	return nil
}
