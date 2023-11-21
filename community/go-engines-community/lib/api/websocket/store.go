package websocket

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type Store interface {
	UpdateConnections(ctx context.Context, conns []UserConnection) error
	GetConnections(ctx context.Context, ids []string) (map[string]int64, error)
	GetActiveConnections(ctx context.Context) (int64, error)
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

func (s *store) UpdateConnections(ctx context.Context, conns []UserConnection) error {
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

func (s *store) GetConnections(ctx context.Context, ids []string) (map[string]int64, error) {
	cursor, err := s.collection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"user":    bson.M{"$in": ids},
			"updated": bson.M{"$gt": datetime.CpsTime{Time: time.Now().Add(-s.readInterval)}},
		}},
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

func (s *store) GetActiveConnections(ctx context.Context) (int64, error) {
	cursor, err := s.collection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"updated": bson.M{"$gt": datetime.CpsTime{Time: time.Now().Add(-s.readInterval)}},
		}},
		{"$group": bson.M{
			"_id":    "$user",
			"tokens": bson.M{"$addToSet": "$token"},
		}},
		{"$group": bson.M{
			"_id":   nil,
			"count": bson.M{"$sum": bson.M{"$size": "$tokens"}},
		}},
	})
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

	return 0, nil
}

func (s *store) GetUsers(ctx context.Context) ([]string, error) {
	cursor, err := s.collection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"updated": bson.M{"$gt": datetime.CpsTime{Time: time.Now().Add(-s.readInterval)}},
		}},
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
