package logger

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// DeleteOne updates author before delete document.
func DeleteOne(ctx context.Context, id, userID string, collection mongo.DbCollection) (int64, error) {
	// required to get the author in action log listener.
	res, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"author": userID}})
	if err != nil || res.MatchedCount == 0 {
		return 0, err
	}

	return collection.DeleteOne(ctx, bson.M{"_id": id})
}

// DeleteByFilter updates author before delete document.
func DeleteByFilter(ctx context.Context, filter bson.M, userID string, collection mongo.DbCollection) (int64, error) {
	// required to get the author in action log listener.
	res, err := collection.UpdateMany(ctx, filter, bson.M{"$set": bson.M{"author": userID}})
	if err != nil || res.MatchedCount == 0 {
		return 0, err
	}

	return collection.DeleteMany(ctx, filter)
}
