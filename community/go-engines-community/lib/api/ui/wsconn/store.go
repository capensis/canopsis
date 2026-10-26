package wsconn

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/websocket"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
)

type Store interface {
	SyncConnections(ctx context.Context, conns []websocket.ConnectionAuthInfo) error
	GetConnections(ctx context.Context, userIDs []string) (map[string]int64, error)
	GetTopConnectedUser(ctx context.Context) (*UserConnections, error)
	CountActiveConnections(ctx context.Context) (int64, error)
	GetUsers(ctx context.Context) ([]string, error)
}

func NewStore(
	client mongo.DbClient,
	ttl time.Duration,
) Store {
	return &store{
		collection:     client.Collection(mongo.WebsocketConnectionMongoCollection),
		deleteInterval: ttl,
		readInterval:   2 * ttl,
	}
}

type store struct {
	collection     mongo.DbCollection
	deleteInterval time.Duration
	readInterval   time.Duration
}

// UserConnections contains the distinct active token count for one user.
type UserConnections struct {
	UserID string `bson:"_id"`
	Count  int64  `bson:"count"`
}

func (s *store) freshnessFilter() bson.M {
	return bson.M{
		"updated": bson.M{"$gt": datetime.CpsTime{Time: time.Now().Add(-s.readInterval)}},
	}
}

// freshConnectionsPipeline returns new stages on each call so callers can
// safely extend the aggregation without sharing mutable BSON values.
func (s *store) freshConnectionsPipeline() []bson.M {
	return []bson.M{
		{"$match": s.freshnessFilter()},
		{"$group": bson.M{
			"_id":    "$user",
			"tokens": bson.M{"$addToSet": "$token"},
		}},
	}
}

func (s *store) GetTopConnectedUser(ctx context.Context) (*UserConnections, error) {
	pipeline := s.freshConnectionsPipeline()
	pipeline = append(pipeline,
		bson.M{"$project": bson.M{
			"_id":   1,
			"count": bson.M{"$size": "$tokens"},
		}},
		bson.M{"$match": bson.M{"count": bson.M{"$gt": 0}}},
		// Keep candidates ranked while checking existence so MongoDB can stop
		// the lookup after it finds the first existing user. The limit must
		// remain after the existence filter so orphaned connections fall back.
		bson.M{"$sort": bson.D{{Key: "count", Value: -1}, {Key: "_id", Value: 1}}},
		bson.M{"$lookup": bson.M{
			"from":         mongo.UserCollection,
			"localField":   "_id",
			"foreignField": "_id",
			"pipeline": bson.A{
				bson.M{"$match": bson.M{"enabled": true}},
				bson.M{"$project": bson.M{"_id": 1}},
			},
			"as": "existingUser",
		}},
		bson.M{"$match": bson.M{"existingUser.0": bson.M{"$exists": true}}},
		bson.M{"$limit": 1},
		bson.M{"$project": bson.M{"_id": 1, "count": 1}},
	)

	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		res := &UserConnections{}
		if err := cursor.Decode(res); err != nil {
			return nil, err
		}
		return res, nil
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *store) SyncConnections(ctx context.Context, conns []websocket.ConnectionAuthInfo) error {
	writeModels := make([]mongodriver.WriteModel, len(conns))
	now := datetime.NewCpsTime()
	for i, conn := range conns {
		writeModels[i] = mongodriver.NewUpdateManyModel().
			SetFilter(bson.M{"_id": conn.ID}).
			SetUpsert(true).
			SetUpdate(bson.M{
				"$set": bson.M{
					"updated": now,
				},
				"$setOnInsert": bson.M{
					"_id":     conn.ID,
					"user":    conn.UserID,
					"token":   conn.Token,
					"created": now,
				},
			})
	}

	if len(writeModels) > 0 {
		_, err := s.collection.BulkWrite(ctx, writeModels)
		if err != nil {
			return err
		}
	}

	_, err := s.collection.DeleteMany(ctx, bson.M{
		"updated": bson.M{"$lt": datetime.CpsTime{Time: now.Add(-s.deleteInterval)}},
	})
	return err
}

func (s *store) GetConnections(ctx context.Context, userIDs []string) (map[string]int64, error) {
	match := s.freshnessFilter()
	match["user"] = bson.M{"$in": userIDs}
	cursor, err := s.collection.Aggregate(ctx, []bson.M{
		{"$match": match},
		{"$group": bson.M{
			"_id":    "$user",
			"tokens": bson.M{"$addToSet": "$token"},
		}},
		{"$project": bson.M{
			"user":  "$_id",
			"count": bson.M{"$size": "$tokens"},
		}},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	res := make(map[string]int64)
	for cursor.Next(ctx) {
		item := struct {
			User  string `bson:"user"`
			Count int64  `bson:"count"`
		}{}
		err = cursor.Decode(&item)
		if err != nil {
			return nil, err
		}

		res[item.User] = item.Count
	}

	return res, nil
}

func (s *store) CountActiveConnections(ctx context.Context) (int64, error) {
	pipeline := s.freshConnectionsPipeline()
	pipeline = append(pipeline, bson.M{"$group": bson.M{
		"_id":   nil,
		"count": bson.M{"$sum": bson.M{"$size": "$tokens"}},
	}})
	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		res := struct {
			Count int64 `bson:"count"`
		}{}
		err = cursor.Decode(&res)
		return res.Count, err
	}
	if err := cursor.Err(); err != nil {
		return 0, err
	}

	return 0, nil
}

func (s *store) GetUsers(ctx context.Context) ([]string, error) {
	cursor, err := s.collection.Aggregate(ctx, []bson.M{
		{"$match": s.freshnessFilter()},
		{"$group": bson.M{
			"_id": "$user",
		}},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	res := make([]string, 0)
	for cursor.Next(ctx) {
		item := struct {
			ID string `bson:"_id"`
		}{}
		err = cursor.Decode(&item)
		if err != nil {
			return nil, err
		}

		res = append(res, item.ID)
	}

	return res, nil
}
