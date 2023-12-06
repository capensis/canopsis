package role

import (
	"context"
	"errors"
	"sort"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	securitymodel "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

const tplLimit = 100

type Store interface {
	Find(ctx context.Context, r ListRequest) (*AggregationResult, error)
	GetOneBy(ctx context.Context, id string) (*Response, error)
	Insert(ctx context.Context, r CreateRequest) (*Response, error)
	Update(ctx context.Context, id string, r EditRequest) (*Response, error)
	Delete(ctx context.Context, id string) (bool, error)
	GetTemplates(ctx context.Context) ([]Template, error)
}

func NewStore(dbClient mongo.DbClient) Store {
	return &store{
		dbClient:               dbClient,
		dbCollection:           dbClient.Collection(mongo.RoleCollection),
		dbPermissionCollection: dbClient.Collection(mongo.PermissionCollection),
		dbUserCollection:       dbClient.Collection(mongo.UserCollection),
		dbTemplateCollection:   dbClient.Collection(mongo.RoleTemplateCollection),
		defaultSearchByFields:  []string{"_id", "name", "description"},
		defaultSortBy:          "name",
	}
}

type store struct {
	dbClient               mongo.DbClient
	dbCollection           mongo.DbCollection
	dbPermissionCollection mongo.DbCollection
	dbUserCollection       mongo.DbCollection
	dbTemplateCollection   mongo.DbCollection
	defaultSearchByFields  []string
	defaultSortBy          string
}

func (s *store) Find(ctx context.Context, r ListRequest) (*AggregationResult, error) {
	pipeline := make([]bson.M, 0)
	filter := common.GetSearchQuery(r.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	sortBy := "name"
	if r.SortBy != "" {
		sortBy = r.SortBy
	}

	pipeline = append(pipeline, getNestedObjectsPipeline()...)
	if r.Permission != "" {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"permissions._id": r.Permission}})
	}
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
			if r.WithFlags {
				isNotAdmin := res.Data[i].ID != security.RoleAdmin
				res.Data[i].Editable = &isNotAdmin
				res.Data[i].Deletable = &isNotAdmin
			}
		}
	}

	return &res, nil
}

func (s *store) GetOneBy(ctx context.Context, id string) (*Response, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"_id": id}},
	}
	pipeline = append(pipeline, getNestedObjectsPipeline()...)
	cursor, err := s.dbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		role := &Response{}
		err := cursor.Decode(role)
		if err != nil {
			return nil, err
		}

		fillRolePermissions(role)

		return role, nil
	}

	return nil, nil
}

func (s *store) Insert(ctx context.Context, r CreateRequest) (*Response, error) {
	types, err := getTypes(ctx, s.dbPermissionCollection, r.Permissions)
	if err != nil {
		return nil, err
	}

	var role *Response
	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		role = nil
		_, err = s.dbCollection.InsertOne(ctx, bson.M{
			"_id":         r.Name,
			"name":        r.Name,
			"description": r.Description,
			"defaultview": r.DefaultView,
			"permissions": transformPermissionsToDoc(r.Permissions, types),
			"auth_config": r.AuthConfig,
		})
		if err != nil {
			return err
		}
		role, err = s.GetOneBy(ctx, r.Name)
		return err
	})
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (s *store) Update(ctx context.Context, id string, r EditRequest) (*Response, error) {
	types, err := getTypes(ctx, s.dbPermissionCollection, r.Permissions)
	if err != nil {
		return nil, err
	}

	var role *Response
	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		role = nil
		res, err := s.dbCollection.UpdateOne(ctx,
			bson.M{"_id": id},
			bson.M{"$set": bson.M{
				"description": r.Description,
				"defaultview": r.DefaultView,
				"permissions": transformPermissionsToDoc(r.Permissions, types),
				"auth_config": r.AuthConfig,
			}},
		)
		if err != nil || res.MatchedCount == 0 {
			return nil
		}
		role, err = s.GetOneBy(ctx, id)
		return err
	})
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	res := s.dbUserCollection.FindOne(ctx, bson.M{"roles": id})
	if err := res.Err(); err != nil {
		if !errors.Is(err, mongodriver.ErrNoDocuments) {
			return false, err
		}
	} else {
		return false, ErrLinkedToUser
	}

	delCount, err := s.dbCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}

	return delCount > 0, nil
}

func (s *store) GetTemplates(ctx context.Context) ([]Template, error) {
	cursor, err := s.dbTemplateCollection.Aggregate(ctx, []bson.M{
		{"$addFields": bson.M{
			"permissions": bson.M{"$objectToArray": "$permissions"},
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
			"_id":         "$_id",
			"name":        bson.M{"$first": "$name"},
			"description": bson.M{"$first": "$description"},
			"permissions": bson.M{"$push": bson.M{"$cond": bson.M{
				"if": "$permissions.model",
				"then": bson.M{"$mergeObjects": bson.A{
					"$permissions.model",
					bson.M{"bitmask": "$permissions.v"},
				}},
				"else": "$$REMOVE",
			}}},
		}},
		{"$limit": tplLimit},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	tpls := make([]Template, 0)
	for cursor.Next(ctx) {
		var tpl Template
		err := cursor.Decode(&tpl)
		if err != nil {
			return nil, err
		}

		for i := range tpl.Permissions {
			tpl.Permissions[i].Actions = TransformBitmaskToActions(tpl.Permissions[i].Bitmask, tpl.Permissions[i].Type)
		}
		tpls = append(tpls, tpl)
	}

	return tpls, nil
}

func getNestedObjectsPipeline() []bson.M {
	return []bson.M{
		{"$addFields": bson.M{
			"permissions": bson.M{"$objectToArray": "$permissions"},
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
			"_id":         "$_id",
			"name":        bson.M{"$first": "$name"},
			"description": bson.M{"$first": "$description"},
			"defaultview": bson.M{"$first": "$defaultview"},
			"auth_config": bson.M{"$first": "$auth_config"},
			"permissions": bson.M{"$push": bson.M{"$cond": bson.M{
				"if": "$permissions.model",
				"then": bson.M{"$mergeObjects": bson.A{
					"$permissions.model",
					bson.M{"bitmask": "$permissions.v"},
				}},
				"else": "$$REMOVE",
			}}},
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

func fillRolePermissions(role *Response) {
	for i := range role.Permissions {
		role.Permissions[i].Actions = TransformBitmaskToActions(role.Permissions[i].Bitmask, role.Permissions[i].Type)
	}
}

func TransformBitmaskToActions(bitmask int64, permType string) []string {
	actions := make([]string, 0)
	switch permType {
	case securitymodel.ObjectTypeCRUD:
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
	case securitymodel.ObjectTypeRW:
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

func transformPermissionsToDoc(permissions map[string][]string, types map[string]string) map[string]int64 {
	doc := make(map[string]int64, len(permissions))
	actionsBitmasks := map[string]int64{
		securitymodel.PermissionCreate: securitymodel.PermissionBitmaskCreate,
		securitymodel.PermissionRead:   securitymodel.PermissionBitmaskRead,
		securitymodel.PermissionUpdate: securitymodel.PermissionBitmaskUpdate,
		securitymodel.PermissionDelete: securitymodel.PermissionBitmaskDelete,
	}

	for id, actions := range permissions {
		bitmask := int64(0)
		if len(actions) == 0 {
			permType, knownType := types[id]
			if !knownType {
				continue
			}
			switch permType {
			case securitymodel.ObjectTypeCRUD, securitymodel.ObjectTypeRW:
				continue
			}
			bitmask = 1
		} else {
			for _, action := range actions {
				bitmask = bitmask | actionsBitmasks[action]
			}
		}

		doc[id] = bitmask
	}

	return doc
}

func getTypes(ctx context.Context, permissionCollection mongo.DbCollection, permissions map[string][]string) (map[string]string, error) {
	ids := make([]string, 0)
	for id := range permissions {
		ids = append(ids, id)
	}
	cursor, err := permissionCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"_id": bson.M{"$in": ids},
			"type": bson.M{"$in": bson.A{
				nil,
				"",
				securitymodel.ObjectTypeCRUD,
				securitymodel.ObjectTypeRW,
			}},
		},
		}, {"$group": bson.M{
			"_id": nil, "d": bson.M{"$push": bson.M{
				"k": "$_id", "v": bson.M{"$ifNull": []string{"$type", ""}},
			}},
		}}, {"$project": bson.M{
			"_id": 0, "d": bson.M{"$arrayToObject": "$d"},
		}},
	})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var types struct {
		Data map[string]string `bson:"d"`
	}

	if !cursor.Next(ctx) {
		return types.Data, nil
	}

	if err := cursor.Decode(&types); err != nil {
		return nil, err
	}
	return types.Data, nil
}
