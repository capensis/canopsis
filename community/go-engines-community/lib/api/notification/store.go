package notification

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const NotificationID = "notification"

type Store interface {
	Get(ctx context.Context) (Notification, error)
	Update(ctx context.Context, request Notification) (Notification, error)
}

type store struct {
	dbCollection mongo.DbCollection
}

func NewStore(
	dbClient mongo.DbClient,
) Store {
	return &store{
		dbCollection: dbClient.Collection(mongo.NotificationMongoCollection),
	}
}

func (s *store) Get(ctx context.Context) (Notification, error) {
	notification := Notification{}
	err := s.dbCollection.FindOne(ctx, bson.M{"_id": NotificationID}).Decode(&notification)

	return notification, err
}

func (s *store) Update(ctx context.Context, request Notification) (Notification, error) {
	notification := Notification{}
	err := s.dbCollection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": NotificationID},
		bson.M{"$set": request},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&notification)

	return notification, err
}
