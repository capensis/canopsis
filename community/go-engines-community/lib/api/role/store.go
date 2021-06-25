package role

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	securitymodel "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"sort"
)

type Store interface {
	Find(ctx context.Context, r ListRequest) (*AggregationResult, error)
	GetOneBy(ctx context.Context, id string) (*Role, error)
	Insert(ctx context.Context, r CreateRequest) (*Role, error)
	Update(ctx context.Context, id string, r EditRequest) (*Role, error)
	Delete(ctx context.Context, id string) (bool, error)
}

func NewStore(dbClient mongo.DbClient) Store {
	return &store{
		dbClient:              dbClient,
		dbCollection:          dbClient.Collection(mongo.RightsMongoCollection),
		defaultSearchByFields: []string{"_id", "crecord_name", "description"},
		defaultSortBy:         "name",
	}
}

type store struct {
	dbClient              mongo.DbClient
	dbCollection          mongo.DbCollection
	defaultSearchByFields []string
	defaultSortBy         string
}

func (s *store) Find(ctx context.Context, r ListRequest) (*AggregationResult, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"crecord_type": securitymodel.LineTypeRole}},
		{"$addFields": bson.M{
			"name": "$crecord_name",
		}},
	}
	filter := common.GetSearchQuery(r.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	sortBy := "name"
	if r.SortBy != "" {
		sortBy = r.SortBy
	}

	pipeline = append(pipeline, getNestedObjectsPipeline()...)
	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		pipeline,
		common.GetSortQuery(sortBy, r.Sort),
	))

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	res := AggregationResult{}

	if cursor.Next(ctx) {
		err := cursor.Decode(&res)
		if err != nil {
			return nil, err
		}

		for i := range res.Data {
			fillRolePermissions(&res.Data[i])
		}
	}

	return &res, nil
}

func (s *store) GetOneBy(ctx context.Context, id string) (*Role, error) {
	pipeline := []bson.M{
		{"$match": bson.M{
			"_id":          id,
			"crecord_type": securitymodel.LineTypeRole,
		}},
	}
	pipeline = append(pipeline, getNestedObjectsPipeline()...)
	cursor, err := s.dbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		role := &Role{}
		err := cursor.Decode(role)
		if err != nil {
			return nil, err
		}

		fillRolePermissions(role)

		return role, nil
	}

	return nil, nil
}

func (s *store) Insert(ctx context.Context, r CreateRequest) (*Role, error) {
	_, err := s.dbCollection.InsertOne(ctx, bson.M{
		"_id":          r.Name,
		"crecord_name": r.Name,
		"crecord_type": securitymodel.LineTypeRole,
		"description":  r.Description,
		"defaultview":  r.DefaultView,
		"rights":       transformPermissionsToDoc(r.Permissions),
	})
	if err != nil {
		return nil, err
	}

	return s.GetOneBy(ctx, r.Name)
}

func (s *store) Update(ctx context.Context, id string, r EditRequest) (*Role, error) {
	res, err := s.dbCollection.UpdateOne(ctx,
		bson.M{"_id": id, "crecord_type": securitymodel.LineTypeRole},
		bson.M{"$set": bson.M{
			"description": r.Description,
			"defaultview": r.DefaultView,
			"rights":      transformPermissionsToDoc(r.Permissions),
		}},
	)
	if err != nil {
		return nil, err
	}

	if res.MatchedCount == 0 {
		return nil, nil
	}

	return s.GetOneBy(ctx, id)
}

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	res := s.dbCollection.FindOne(ctx, bson.M{
		"crecord_type": securitymodel.LineTypeSubject,
		"role":         id,
	})
	if err := res.Err(); err != nil {
		if err != mongodriver.ErrNoDocuments {
			return false, err
		}
	} else {
		return false, ErrLinkedToUser
	}

	delCount, err := s.dbCollection.DeleteOne(ctx, bson.M{
		"_id":          id,
		"crecord_type": securitymodel.LineTypeRole,
	})
	if err != nil {
		return false, err
	}

	return delCount > 0, nil
}

func getNestedObjectsPipeline() []bson.M {
	return []bson.M{
		{"$addFields": bson.M{
			"permissions": bson.M{"$map": bson.M{
				"input": bson.M{"$objectToArray": "$rights"},
				"as":    "each",
				"in":    "$$each.k",
			}},
			"bitmasks": bson.M{"$map": bson.M{
				"input": bson.M{"$objectToArray": "$rights"},
				"as":    "each",
				"in": bson.M{
					"k": "$$each.k",
					"v": "$$each.v.checksum",
				},
			}},
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
			"_id":         "$_id",
			"name":        bson.M{"$first": "$crecord_name"},
			"description": bson.M{"$first": "$description"},
			"defaultview": bson.M{"$first": "$defaultview"},
			"permissions": bson.M{"$push": "$permissions"},
		}},
		{"$lookup": bson.M{
			"from":         mongo.ViewMongoCollection,
			"localField":   "defaultview",
			"foreignField": "_id",
			"as":           "defaultview",
		}},
		{"$unwind": bson.M{"path": "$defaultview", "preserveNullAndEmptyArrays": true}},
	}
}

func fillRolePermissions(role *Role) {
	for i := range role.Permissions {
		role.Permissions[i].Actions = TransformBitmaskToActions(role.Permissions[i].Bitmask, role.Permissions[i].Type)
	}
}

func TransformBitmaskToActions(bitmask int64, roleType string) []string {
	actions := make([]string, 0)
	switch roleType {
	case securitymodel.LineObjectTypeCRUD:
		actionsBitmasks := map[string]int64{
			securitymodel.PermissionCreate: securitymodel.PermissionBitmaskCreate,
			securitymodel.PermissionRead:   securitymodel.PermissionBitmaskRead,
			securitymodel.PermissionUpdate: securitymodel.PermissionBitmaskUpdate,
			securitymodel.PermissionDelete: securitymodel.PermissionBitmaskDelete,
		}
		for action, actionBitmask := range actionsBitmasks {
			if bitmask&actionBitmask == actionBitmask {
				actions = append(actions, action)
			}
		}
	case securitymodel.LineObjectTypeRW:
		actionsBitmasks := map[string]int64{
			securitymodel.PermissionRead:   securitymodel.PermissionBitmaskRead,
			securitymodel.PermissionUpdate: securitymodel.PermissionBitmaskUpdate,
			securitymodel.PermissionDelete: securitymodel.PermissionBitmaskDelete,
		}
		for action, actionBitmask := range actionsBitmasks {
			if bitmask&actionBitmask == actionBitmask {
				actions = append(actions, action)
			}
		}
	}

	sort.Strings(actions)

	return actions
}

func transformPermissionsToDoc(permissions map[string][]string) map[string]interface{} {
	rights := make(map[string]interface{}, len(permissions))
	actionsBitmasks := map[string]int64{
		securitymodel.PermissionCreate: securitymodel.PermissionBitmaskCreate,
		securitymodel.PermissionRead:   securitymodel.PermissionBitmaskRead,
		securitymodel.PermissionUpdate: securitymodel.PermissionBitmaskUpdate,
		securitymodel.PermissionDelete: securitymodel.PermissionBitmaskDelete,
	}

	for id, actions := range permissions {
		bitmask := int64(0)
		if len(actions) == 0 {
			bitmask = 1
		} else {
			for _, action := range actions {
				bitmask = bitmask | actionsBitmasks[action]
			}
		}

		rights[id] = map[string]int64{
			"checksum": bitmask,
		}
	}

	return rights
}
