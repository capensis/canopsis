package notification

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const NotificationID = "notification"

type Store interface {
	Get() (Notification, error)
	Update(request Notification) (Notification, error)
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

func (s *store) Get() (Notification, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	notification := Notification{}
	err := s.dbCollection.FindOne(ctx, bson.M{"_id": NotificationID}).Decode(&notification)

	return notification, err
}

func (s *store) Update(request Notification) (Notification, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	notification := Notification{}
	err := s.dbCollection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": NotificationID},
		bson.M{"$set": request},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&notification)

	return notification, err
}
