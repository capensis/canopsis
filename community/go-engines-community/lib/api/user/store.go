package user

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	securitymodel "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/password"
	"go.mongodb.org/mongo-driver/bson"
)

type Store interface {
	Find(ctx context.Context, r ListRequest, userID string) (*AggregationResult, error)
	GetOneBy(ctx context.Context, id string) (*User, error)
	Insert(ctx context.Context, r CreateRequest) (*User, error)
	Update(ctx context.Context, r UpdateRequest, userID string) (*User, error)
	Delete(ctx context.Context, id, userID string) (bool, error)
}

func NewStore(dbClient mongo.DbClient, passwordEncoder password.Encoder, websocketStore websocket.Store) Store {
	return &store{
		client:                 dbClient,
		collection:             dbClient.Collection(mongo.RightsMongoCollection),
		userPrefCollection:     dbClient.Collection(mongo.UserPreferencesMongoCollection),
		patternCollection:      dbClient.Collection(mongo.PatternMongoCollection),
		widgetFilterCollection: dbClient.Collection(mongo.WidgetFiltersMongoCollection),
		shareTokenCollection:   dbClient.Collection(mongo.ShareTokenMongoCollection),

		passwordEncoder: passwordEncoder,
		websocketStore:  websocketStore,

		defaultSearchByFields: []string{"_id", "crecord_name", "firstname", "lastname", "role.name"},
		defaultSortBy:         "name",
	}
}

type store struct {
	client                 mongo.DbClient
	collection             mongo.DbCollection
	userPrefCollection     mongo.DbCollection
	patternCollection      mongo.DbCollection
	widgetFilterCollection mongo.DbCollection
	shareTokenCollection   mongo.DbCollection

	passwordEncoder password.Encoder
	websocketStore  websocket.Store

	defaultSearchByFields []string
	defaultSortBy         string
}

func (s *store) Find(ctx context.Context, r ListRequest, curUserID string) (*AggregationResult, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"crecord_type": securitymodel.LineTypeSubject}},
	}
	pipeline = append(pipeline, getRenameFieldsPipeline()...)
	project := make([]bson.M, 0)

	filter := common.GetSearchQuery(r.Search, s.defaultSearchByFields)
	if len(filter) > 0 || r.Permission != "" || r.SortBy == "role.name" {
		pipeline = append(pipeline, getRolePipeline()...)
	} else {
		project = append(project, getRolePipeline()...)
	}

	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	if r.Permission != "" {
		pipeline = append(pipeline, bson.M{"$match": bson.M{fmt.Sprintf("role.rights.%s", r.Permission): bson.M{"$exists": true}}})
	}

	project = append(project, getViewPipeline()...)

	sortBy := s.defaultSortBy
	if r.SortBy != "" {
		sortBy = r.SortBy
	}

	cursor, err := s.collection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		pipeline,
		common.GetSortQuery(sortBy, r.Sort),
		project,
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
	}

	ids := make([]string, len(res.Data))
	for i, v := range res.Data {
		ids[i] = v.ID
	}

	conns, err := s.websocketStore.GetConnections(ctx, ids)
	if err != nil {
		return nil, err
	}

	var onlyOneAdmin bool
	var lastAdminID string
	if r.WithFlags {
		onlyOneAdmin, lastAdminID, err = s.checkLastAdmin(ctx)
		if err != nil {
			return nil, err
		}
	}

	for i := range res.Data {
		activeConnects := conns[res.Data[i].ID]
		res.Data[i].ActiveConnects = &activeConnects
		if r.WithFlags {
			deletable := res.Data[i].ID != curUserID && (!onlyOneAdmin || res.Data[i].ID != lastAdminID)
			res.Data[i].Deletable = &deletable
		}
	}

	return &res, nil
}

func (s *store) GetOneBy(ctx context.Context, id string) (*User, error) {
	pipeline := []bson.M{
		{"$match": bson.M{
			"_id":          id,
			"crecord_type": securitymodel.LineTypeSubject,
		}},
	}
	pipeline = append(pipeline, getNestedObjectsPipeline()...)
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

		return user, nil
	}

	return nil, nil
}

func (s *store) Insert(ctx context.Context, r CreateRequest) (*User, error) {
	var user *User
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		user = nil
		_, err := s.collection.InsertOne(ctx, r.getBson(s.passwordEncoder))
		if err != nil {
			return err
		}

		user, err = s.GetOneBy(ctx, r.Name)
		return err
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *store) Update(ctx context.Context, r UpdateRequest, curUserID string) (*User, error) {
	if r.ID == curUserID && r.IsEnabled != nil && !*r.IsEnabled {
		return nil, common.NewValidationError("enable", "user cannot disable itself")
	}

	var user *User
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		user = nil
		onlyOneAdmin, lastAdminID, err := s.checkLastAdmin(ctx)
		if err != nil {
			return err
		}

		if onlyOneAdmin && lastAdminID == r.ID {
			if r.Role != security.RoleAdmin {
				return common.NewValidationError("role", "last admin cannot be edited")
			}

			if r.IsEnabled != nil && !*r.IsEnabled {
				return common.NewValidationError("enable", "last admin cannot be disabled")
			}
		}

		res, err := s.collection.UpdateOne(ctx,
			bson.M{"_id": r.ID, "crecord_type": securitymodel.LineTypeSubject},
			bson.M{"$set": r.getBson(s.passwordEncoder)},
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

func (s *store) Delete(ctx context.Context, id, userID string) (bool, error) {
	if id == userID {
		return false, common.NewValidationError("_id", "user cannot delete itself")
	}

	onlyOneAdmin, lastAdminID, err := s.checkLastAdmin(ctx)
	if err != nil {
		return false, err
	}

	if onlyOneAdmin && id == lastAdminID {
		return false, common.NewValidationError("_id", "last admin cannot be deleted")
	}

	delCount, err := s.collection.DeleteOne(ctx, bson.M{
		"_id":          id,
		"crecord_type": securitymodel.LineTypeSubject,
	})
	if err != nil {
		return false, err
	}

	if delCount == 0 {
		return false, nil
	}

	err = s.deleteUserPreferences(ctx, id)
	if err != nil {
		return false, err
	}

	err = s.deletePatterns(ctx, id)
	if err != nil {
		return false, err
	}

	err = s.deleteWidgetFilters(ctx, id)
	if err != nil {
		return false, err
	}

	err = s.deleteShareTokens(ctx, id)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *store) deleteUserPreferences(ctx context.Context, id string) error {
	_, err := s.userPrefCollection.DeleteMany(ctx, bson.M{
		"user": id,
	})

	return err
}

func (s *store) deletePatterns(ctx context.Context, id string) error {
	_, err := s.patternCollection.DeleteMany(ctx, bson.M{
		"author":       id,
		"is_corporate": false,
	})

	return err
}

func (s *store) deleteWidgetFilters(ctx context.Context, id string) error {
	_, err := s.widgetFilterCollection.DeleteMany(ctx, bson.M{
		"author":     id,
		"is_private": true,
	})

	return err
}

func (s *store) deleteShareTokens(ctx context.Context, id string) error {
	_, err := s.shareTokenCollection.DeleteMany(ctx, bson.M{
		"user": id,
	})

	return err
}

func (s *store) checkLastAdmin(ctx context.Context) (bool, string, error) {
	cursor, err := s.collection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"crecord_type": securitymodel.LineTypeSubject,
			"enable":       true,
			"role":         security.RoleAdmin,
		}},
		{"$group": bson.M{
			"_id":     nil,
			"count":   bson.M{"$sum": 1},
			"last_id": bson.M{"$first": "$_id"},
		}},
	})
	if err != nil {
		return false, "", err
	}

	defer cursor.Close(ctx)
	res := struct {
		Count  int64  `bson:"count"`
		LastID string `bson:"last_id"`
	}{}
	if cursor.Next(ctx) {
		err = cursor.Decode(&res)
		if err != nil {
			return false, "", err
		}
	}

	return res.Count <= 1, res.LastID, nil
}

func getNestedObjectsPipeline() []bson.M {
	pipeline := getRenameFieldsPipeline()
	pipeline = append(pipeline, getRolePipeline()...)
	pipeline = append(pipeline, getViewPipeline()...)

	return pipeline
}

func getRolePipeline() []bson.M {
	return []bson.M{
		{"$graphLookup": bson.M{
			"from":             mongo.RightsMongoCollection,
			"startWith":        "$role",
			"connectFromField": "role",
			"connectToField":   "_id",
			"as":               "role",
			"maxDepth":         0,
		}},
		{"$unwind": bson.M{"path": "$role", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			"role": bson.M{
				"_id":         "$role._id",
				"name":        "$role.crecord_name",
				"rights":      "$role.rights",
				"defaultview": "$role.defaultview",
			},
		}},
	}
}

func getViewPipeline() []bson.M {
	return []bson.M{
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
	}
}

func getRenameFieldsPipeline() []bson.M {
	return []bson.M{
		{"$addFields": bson.M{
			"name": "$crecord_name",
		}},
	}
}
