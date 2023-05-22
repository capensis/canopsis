package account

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/role"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/password"
	"go.mongodb.org/mongo-driver/bson"
)

type Store interface {
	GetOneBy(ctx context.Context, id string) (*User, error)
	Update(ctx context.Context, r EditRequest) (*User, error)
}

type store struct {
	client          mongo.DbClient
	collection      mongo.DbCollection
	passwordEncoder password.Encoder
	authorProvider  author.Provider
}

func NewStore(db mongo.DbClient, passwordEncoder password.Encoder, authorProvider author.Provider) Store {
	return &store{
		client:          db,
		collection:      db.Collection(mongo.UserCollection),
		passwordEncoder: passwordEncoder,
		authorProvider:  authorProvider,
	}
}

func (s *store) GetOneBy(ctx context.Context, id string) (*User, error) {
	cursor, err := s.collection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{"_id": id}},
		// Find role
		{"$lookup": bson.M{
			"from":         mongo.RoleCollection,
			"localField":   "role",
			"foreignField": "_id",
			"as":           "role",
		}},
		{"$unwind": bson.M{"path": "$role", "preserveNullAndEmptyArrays": true}},
		// Find permissions
		{"$addFields": bson.M{
			"permissions": bson.M{"$objectToArray": "$role.permissions"},
		}},
		{"$unwind": bson.M{"path": "$permissions", "preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from":         mongo.PermissionCollection,
			"localField":   "permissions.k",
			"foreignField": "_id",
			"as":           "permissions.model",
		}},
		{"$unwind": bson.M{"path": "$permissions.model", "preserveNullAndEmptyArrays": true}},
		{"$sort": bson.M{"permissions.model.name": 1}},
		{"$group": bson.M{
			"_id":                       "$_id",
			"name":                      bson.M{"$first": "$name"},
			"lastname":                  bson.M{"$first": "$lastname"},
			"firstname":                 bson.M{"$first": "$firstname"},
			"email":                     bson.M{"$first": "$email"},
			"role":                      bson.M{"$first": "$role"},
			"ui_language":               bson.M{"$first": "$ui_language"},
			"ui_theme":                  bson.M{"$first": "$ui_theme"},
			"ui_tours":                  bson.M{"$first": "$ui_tours"},
			"ui_groups_navigation_type": bson.M{"$first": "$ui_groups_navigation_type"},
			"enable":                    bson.M{"$first": "$enable"},
			"defaultview":               bson.M{"$first": "$defaultview"},
			"external_id":               bson.M{"$first": "$external_id"},
			"source":                    bson.M{"$first": "$source"},
			"authkey":                   bson.M{"$first": "$authkey"},
			"paused_executions":         bson.M{"$first": "$paused_executions"},
			"permissions": bson.M{"$push": bson.M{"$cond": bson.M{
				"if": "$permissions.model",
				"then": bson.M{"$mergeObjects": bson.A{
					"$permissions.model",
					bson.M{"bitmask": "$permissions.v"},
				}},
				"else": "$$REMOVE",
			}}},
		}},
		// Find defaultview
		{"$lookup": bson.M{
			"from":         mongo.ViewMongoCollection,
			"localField":   "defaultview",
			"foreignField": "_id",
			"as":           "defaultview",
		}},
		{"$unwind": bson.M{"path": "$defaultview", "preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from":         mongo.ViewMongoCollection,
			"localField":   "role.defaultview",
			"foreignField": "_id",
			"as":           "role.defaultview",
		}},
		{"$unwind": bson.M{"path": "$role.defaultview", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			"username": "$name",
		}},
		{"$addFields": bson.M{
			"display_name": s.authorProvider.GetDisplayNameQuery(""),
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

		for i := range user.Permissions {
			user.Permissions[i].Actions = role.TransformBitmaskToActions(user.Permissions[i].Bitmask, user.Permissions[i].Type)
		}

		return user, nil
	}

	return nil, nil
}

func (s *store) Update(ctx context.Context, r EditRequest) (*User, error) {
	var user *User
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		user = nil
		res, err := s.collection.UpdateOne(ctx,
			bson.M{"_id": r.ID},
			bson.M{"$set": r.getUpdateBson(s.passwordEncoder)},
		)
		if err != nil || res.MatchedCount == 0 {
			return err
		}

		user, err = s.GetOneBy(ctx, r.ID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}
