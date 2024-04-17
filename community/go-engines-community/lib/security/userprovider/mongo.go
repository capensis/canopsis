// Package userprovider contains user storages.
package userprovider

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongoProvider decorates request to mongo db.
type mongoProvider struct {
	client         libmongo.DbClient
	collection     libmongo.DbCollection
	configProvider config.ApiConfigProvider
}

// NewMongoProvider creates new provider.
func NewMongoProvider(db libmongo.DbClient, configProvider config.ApiConfigProvider) security.UserProvider {
	return &mongoProvider{
		client:         db,
		collection:     db.Collection(libmongo.UserCollection),
		configProvider: configProvider,
	}
}

func (p *mongoProvider) FindWithoutPermission(ctx context.Context, perm string) ([]security.User, error) {
	var users []security.User

	cursor, err := p.collection.Aggregate(ctx, []bson.M{
		{
			"$lookup": bson.M{
				"from":         "role",
				"localField":   "roles",
				"foreignField": "_id",
				"as":           "roles_objects",
			},
		},
		{
			"$match": bson.M{
				"roles_objects": bson.M{
					"$elemMatch": bson.M{
						"permissions." + perm: bson.M{"$exists": false},
					},
				},
			},
		},
	})
	if err != nil {
		return users, err
	}

	return users, cursor.All(ctx, &users)
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
	externalID, source string,
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

	u.DisplayName = ""
	err := p.client.WithTransaction(ctx, func(ctx context.Context) error {
		_, err := p.collection.UpdateOne(
			ctx,
			bson.M{"_id": u.ID},
			bson.M{"$set": u},
			options.Update().SetUpsert(true),
		)

		if err != nil {
			return err
		}

		newUser, err := p.findByFilter(ctx, bson.M{"_id": u.ID})
		if err != nil {
			return err
		}

		*u = *newUser

		return nil
	})

	return err
}

func (p *mongoProvider) UpdateHashedPassword(ctx context.Context, id, hash string) error {
	res, err := p.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{
			"password": hash,
		}},
	)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}

// findByFilter returns User or nil if no user matches filter.
func (p *mongoProvider) findByFilter(ctx context.Context, match bson.M) (*security.User, error) {
	cursor, err := p.collection.Aggregate(ctx, []bson.M{
		{"$match": match},
		{"$addFields": bson.M{
			"username": "$name",
		}},
		{"$addFields": bson.M{
			"display_name": p.getDisplayNameQuery(),
		}},
	})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var u security.User
		err := cursor.Decode(&u)
		if err != nil {
			return nil, err
		}

		return &u, nil
	}

	return nil, nil
}

func (p *mongoProvider) getDisplayNameQuery() bson.M {
	authorScheme := p.configProvider.Get().AuthorScheme
	concat := make([]any, len(authorScheme))
	for i, v := range authorScheme {
		if len(v) > 0 && v[0] == '$' {
			concat[i] = bson.M{"$ifNull": bson.A{v, ""}}
		} else {
			concat[i] = v
		}
	}

	return bson.M{"$concat": concat}
}
