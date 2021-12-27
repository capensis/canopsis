package account

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/role"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"go.mongodb.org/mongo-driver/bson"
)

type Store interface {
	GetOneBy(ctx context.Context, id string) (*User, error)
}

type store struct {
	collection mongo.DbCollection
}

func NewStore(db mongo.DbClient) Store {
	return &store{
		collection: db.Collection(mongo.RightsMongoCollection),
	}
}

func (s *store) GetOneBy(ctx context.Context, id string) (*User, error) {
	cursor, err := s.collection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{"_id": id, "crecord_type": model.LineTypeSubject}},
		// Find role
		{"$graphLookup": bson.M{
			"from":             mongo.RightsMongoCollection,
			"startWith":        "$role",
			"connectFromField": "role",
			"connectToField":   "_id",
			"as":               "role",
		}},
		{"$unwind": bson.M{"path": "$role", "preserveNullAndEmptyArrays": true}},
		// Find permissions
		{"$addFields": bson.M{
			"permissions": bson.M{"$map": bson.M{
				"input": bson.M{"$objectToArray": "$role.rights"},
				"as":    "each",
				"in":    "$$each.k",
			}},
			"bitmasks": bson.M{"$map": bson.M{
				"input": bson.M{"$objectToArray": "$role.rights"},
				"as":    "each",
				"in": bson.M{
					"k": "$$each.k",
					"v": "$$each.v.checksum",
				},
			}},
			"role": bson.M{
				"_id":  "$role._id",
				"name": "$role.crecord_name",
			},
		}},
		{"$graphLookup": bson.M{
			"from":             mongo.RightsMongoCollection,
			"startWith":        "$permissions",
			"connectFromField": "permissions",
			"connectToField":   "_id",
			"as":               "permissions",
		}},
		{"$addFields": bson.M{
			"permissions": bson.M{"$map": bson.M{
				"input": "$permissions",
				"as":    "each",
				"in": bson.M{
					"_id":         "$$each._id",
					"name":        "$$each.crecord_name",
					"description": "$$each.desc",
					"type":        "$$each.type",
					"bitmask": bson.M{"$arrayElemAt": bson.A{
						bson.M{"$map": bson.M{
							"input": bson.M{"$filter": bson.M{
								"input": "$bitmasks",
								"as":    "bitmask",
								"cond": bson.M{
									"$eq": bson.A{"$$each._id", "$$bitmask.k"},
								},
							}},
							"as": "bitmask",
							"in": "$$bitmask.v",
						}},
						0,
					}},
				},
			}},
		}},
		{"$unwind": bson.M{"path": "$permissions", "preserveNullAndEmptyArrays": true}},
		{"$sort": bson.M{"permissions.name": 1}},
		{"$group": bson.M{
			"_id":                       "$_id",
			"name":                      bson.M{"$first": "$crecord_name"},
			"lastname":                  bson.M{"$first": "$lastname"},
			"firstname":                 bson.M{"$first": "$firstname"},
			"email":                     bson.M{"$first": "$mail"},
			"role":                      bson.M{"$first": "$role"},
			"ui_language":               bson.M{"$first": "$ui_language"},
			"ui_tours":                  bson.M{"$first": "$tours"},
			"ui_groups_navigation_type": bson.M{"$first": "$groupsNavigationType"},
			"enable":                    bson.M{"$first": "$enable"},
			"defaultview":               bson.M{"$first": "$defaultview"},
			"external_id":               bson.M{"$first": "$external_id"},
			"source":                    bson.M{"$first": "$source"},
			"authkey":                   bson.M{"$first": "$authkey"},
			"paused_executions":         bson.M{"$first": "$paused_executions"},
			"permissions":               bson.M{"$push": "$permissions"},
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
