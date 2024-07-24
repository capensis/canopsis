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
	pipeline := []bson.M{
		{"$match": bson.M{"_id": id}},
		// Find permissions
		{"$lookup": bson.M{
			"from":         mongo.RoleCollection,
			"localField":   "roles",
			"foreignField": "_id",
			"as":           "role",
		}},
		{"$unwind": bson.M{"path": "$role", "preserveNullAndEmptyArrays": true}},
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
			"_id":  "$_id",
			"data": bson.M{"$first": "$$ROOT"},
			"permissions": bson.M{"$push": bson.M{"$cond": bson.M{
				"if": "$permissions.model",
				"then": bson.M{"$mergeObjects": bson.A{
					"$permissions.model",
					bson.M{"bitmask": "$permissions.v"},
				}},
				"else": "$$REMOVE",
			}}},
		}},
		{"$replaceRoot": bson.M{"newRoot": bson.M{"$mergeObjects": bson.A{
			"$data",
			bson.M{"permissions": "$permissions"},
		}}}},
		// Find defaultview
		{"$lookup": bson.M{
			"from":         mongo.ViewMongoCollection,
			"localField":   "defaultview",
			"foreignField": "_id",
			"as":           "defaultview",
		}},
		{"$unwind": bson.M{"path": "$defaultview", "preserveNullAndEmptyArrays": true}},
		{"$unwind": bson.M{
			"path":                       "$roles",
			"preserveNullAndEmptyArrays": true,
			"includeArrayIndex":          "role_index",
		}},
		{"$lookup": bson.M{
			"from":         mongo.RoleCollection,
			"localField":   "roles",
			"foreignField": "_id",
			"as":           "roles",
		}},
		{"$unwind": bson.M{"path": "$roles", "preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from":         mongo.ViewMongoCollection,
			"localField":   "roles.defaultview",
			"foreignField": "_id",
			"as":           "roles.defaultview",
		}},
		{"$unwind": bson.M{"path": "$roles.defaultview", "preserveNullAndEmptyArrays": true}},
		{"$sort": bson.M{"role_index": 1}},
		{"$group": bson.M{
			"_id":  "$_id",
			"data": bson.M{"$first": "$$ROOT"},
			"roles": bson.M{"$push": bson.M{
				"$cond": bson.M{
					"if":   "$roles._id",
					"then": "$roles",
					"else": "$$REMOVE",
				},
			}},
		}},
		{"$replaceRoot": bson.M{"newRoot": bson.M{"$mergeObjects": bson.A{
			"$data",
			bson.M{"roles": "$roles"},
		}}}},
		{"$addFields": bson.M{
			"username": "$name",
		}},
		{"$addFields": bson.M{
			"display_name": s.authorProvider.GetDisplayNameQuery(""),
		}},
		{
			"$addFields": bson.M{
				"ui_theme": bson.M{
					"$cond": bson.M{
						"if": bson.M{
							"$or": bson.A{
								bson.M{"$eq": bson.A{"$ui_theme", ""}},
								bson.M{"$eq": bson.A{bson.M{"$ifNull": bson.A{"$ui_theme", ""}}, ""}},
							},
						},
						"then": "canopsis",
						"else": "$ui_theme",
					},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         mongo.ColorThemeCollection,
				"localField":   "ui_theme",
				"foreignField": "_id",
				"as":           "ui_theme",
			},
		},
		{
			"$unwind": bson.M{"path": "$ui_theme", "preserveNullAndEmptyArrays": true},
		},
	}
	pipeline = append(pipeline, s.authorProvider.Pipeline()...)
	pipeline = append(pipeline, s.authorProvider.PipelineForField("ui_theme.author")...)

	cursor, err := s.collection.Aggregate(ctx, pipeline)
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

		permIndexes := make(map[string]int, len(user.Permissions))
		k := 0
		for _, perm := range user.Permissions {
			if idx, ok := permIndexes[perm.ID]; ok {
				user.Permissions[idx].Bitmask |= perm.Bitmask
			} else {
				user.Permissions[k] = perm
				permIndexes[perm.ID] = k
				k++
			}
		}
		user.Permissions = user.Permissions[:k]

		for i, perm := range user.Permissions {
			user.Permissions[i].Actions = role.TransformBitmaskToActions(perm.Bitmask, perm.Type)
		}

		return user, nil
	}

	return nil, nil
}

func (s *store) Update(ctx context.Context, r EditRequest) (*User, error) {
	var user *User
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		user = nil
		updateDoc, err := r.getUpdateBson(s.passwordEncoder)
		if err != nil {
			return err
		}

		res, err := s.collection.UpdateOne(ctx,
			bson.M{"_id": r.ID},
			bson.M{"$set": updateDoc},
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
