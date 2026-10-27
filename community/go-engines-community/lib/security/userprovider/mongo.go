// Package userprovider contains user storages.
package userprovider

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// mongoProvider decorates request to mongo db.
type mongoProvider struct {
	client         libmongo.DbClient
	collection     libmongo.DbCollection
	configProvider config.ApiConfigProvider
}

type roleNamesReader struct {
	collection libmongo.DbCollection
}

type roleNamesAggregationResult struct {
	Roles []struct {
		Name string `bson:"name"`
	} `bson:"role_names"`
}

// NewRoleNamesReader creates a reader for current role names assigned to users.
func NewRoleNamesReader(db libmongo.DbClient) security.UserRoleReader {
	return &roleNamesReader{collection: db.Collection(libmongo.UserCollection)}
}

func (r *roleNamesReader) GetRoleNames(ctx context.Context, userID string) ([]string, error) {
	if userID == "" {
		return []string{}, nil
	}

	cursor, err := r.collection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{"_id": userID}},
		{"$project": bson.M{"roles": 1}},
		{"$lookup": bson.M{
			"from":         libmongo.RoleCollection,
			"localField":   "roles",
			"foreignField": "_id",
			"pipeline":     []bson.M{{"$project": bson.M{"_id": 0, "name": 1}}},
			"as":           "role_names",
		}},
		{"$project": bson.M{"_id": 0, "role_names.name": 1}},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate role names for user %q: %w", userID, err)
	}
	defer cursor.Close(ctx)

	if !cursor.Next(ctx) {
		if err := cursor.Err(); err != nil {
			return nil, fmt.Errorf("failed to read role names for user %q: %w", userID, err)
		}
		if err := ctx.Err(); err != nil {
			return nil, fmt.Errorf("failed to read role names for user %q: %w", userID, err)
		}
		return []string{}, nil
	}

	var result roleNamesAggregationResult
	if err := cursor.Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode role names for user %q: %w", userID, err)
	}
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("failed to read role names for user %q: %w", userID, err)
	}
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("failed to read role names for user %q: %w", userID, err)
	}

	roleNames := make([]string, 0, len(result.Roles))
	for _, role := range result.Roles {
		if role.Name != "" {
			roleNames = append(roleNames, role.Name)
		}
	}
	slices.Sort(roleNames)

	return roleNames, nil
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
	return p.findByFilter(ctx, bson.M{
		"_id": id,
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
			options.UpdateOne().SetUpsert(true),
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
			"language": "$ui_language",
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
