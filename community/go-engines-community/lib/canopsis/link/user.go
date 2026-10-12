package link

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var ErrUserNotFound = errors.New("user not found")

type UserProvider interface {
	User(ctx context.Context, id string) (User, error)
}

func NewUserProvider(client mongo.DbClient) UserProvider {
	return &userProvider{
		userCollection: client.Collection(mongo.UserCollection),
	}
}

type userProvider struct {
	userCollection mongo.DbCollection
}

func (p *userProvider) User(ctx context.Context, id string) (user User, err error) {
	cursor, err := p.userCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{"_id": id}},
		{"$addFields": bson.M{"username": "$name"}},
	})
	if err != nil {
		return user, err
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		err = cursor.Decode(&user)

		return user, err
	}

	return user, ErrUserNotFound
}
