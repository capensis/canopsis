package account

import (
	"context"

	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type Store interface {
	GetOneBy(id string) (*User, error)
}

type store struct {
	db         mongo.DbClient
	collection mongo.DbCollection
}

func NewStore(db mongo.DbClient) Store {
	return &store{
		db:         db,
		collection: db.Collection(mongo.RightsMongoCollection),
	}
}

func (s *store) GetOneBy(id string) (*User, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cursor, err := s.collection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{"_id": id}},
		{"$lookup": bson.M{
			"from":         mongo.RightsMongoCollection,
			"localField":   "role",
			"foreignField": "_id",
			"as":           "rights",
		}},
		{"$unwind": bson.M{"path": "$rights", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			"rights": "$rights.rights",
		}},
	})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		user := &User{}
		err := cursor.Decode(user)
		if err != nil {
			return nil, err
		}

		return user, nil
	}

	return nil, nil
}
