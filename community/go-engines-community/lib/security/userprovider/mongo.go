// Package userprovider contains user storages.
package userprovider

import (
	"context"
	"errors"

	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongoProvider decorates request to mongo db.
type mongoProvider struct {
	collection libmongo.DbCollection
}

// NewMongoProvider creates new provider.
func NewMongoProvider(db libmongo.DbClient) security.UserProvider {
	return &mongoProvider{
		collection: db.Collection(libmongo.UserCollection),
	}
}

func (p *mongoProvider) FindByUsername(ctx context.Context, username string) (*security.User, error) {
	return p.findByFilter(ctx, bson.M{
		"name":   username,
		"source": bson.M{"$in": bson.A{"", nil}},
	})
}

func (p *mongoProvider) FindByAuthApiKey(ctx context.Context, apiKey string) (*security.User, error) {
	return p.findByFilter(ctx, bson.M{
		"authkey": apiKey,
	})
}

func (p *mongoProvider) FindByID(ctx context.Context, id string) (*security.User, error) {
	var objID interface{}
	var err error
	objID, err = primitive.ObjectIDFromHex(id)
	if err != nil {
		objID = id
	}

	return p.findByFilter(ctx, bson.M{
		"_id": objID,
	})
}

func (p *mongoProvider) FindByExternalSource(
	ctx context.Context,
	externalID string,
	source security.Source,
) (*security.User, error) {
	return p.findByFilter(ctx, bson.M{
		"external_id": externalID,
		"source":      source,
	})
}

func (p *mongoProvider) Save(ctx context.Context, u *security.User) error {
	if u.ID == "" {
		u.ID = utils.NewID()
		u.AuthApiKey = utils.NewID()
	}

	_, err := p.collection.UpdateOne(
		ctx,
		bson.M{"_id": u.ID},
		bson.M{"$set": u},
		options.Update().SetUpsert(true),
	)

	if err != nil {
		return err
	}

	return nil
}

// findByFilter returns User or nil if no user matches filter.
func (p *mongoProvider) findByFilter(ctx context.Context, f interface{}) (*security.User, error) {
	var u security.User
	err := p.collection.FindOne(ctx, f).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}
