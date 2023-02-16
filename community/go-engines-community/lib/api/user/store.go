package user

import (
	"context"
	"fmt"
	"math"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	securitymodel "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/password"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type Store interface {
	Find(ctx context.Context, r ListRequest) (*AggregationResult, error)
	GetOneBy(ctx context.Context, id string) (*User, error)
	Insert(ctx context.Context, r CreateRequest) (*User, error)
	Update(ctx context.Context, r UpdateRequest) (*User, error)
	Delete(ctx context.Context, id string) (bool, error)

	BulkInsert(ctx context.Context, requests []CreateRequest) error
	BulkUpdate(ctx context.Context, requests []BulkUpdateRequestItem) error
	BulkDelete(ctx context.Context, ids []string) error
}

func NewStore(dbClient mongo.DbClient, passwordEncoder password.Encoder, websocketHub websocket.Hub) Store {
	return &store{
		client:                 dbClient,
		collection:             dbClient.Collection(mongo.RightsMongoCollection),
		userPrefCollection:     dbClient.Collection(mongo.UserPreferencesMongoCollection),
		patternCollection:      dbClient.Collection(mongo.PatternMongoCollection),
		widgetFilterCollection: dbClient.Collection(mongo.WidgetFiltersMongoCollection),
		shareTokenCollection:   dbClient.Collection(mongo.ShareTokenMongoCollection),

		passwordEncoder: passwordEncoder,
		websocketHub:    websocketHub,

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
	websocketHub    websocket.Hub

	defaultSearchByFields []string
	defaultSortBy         string
}

func (s *store) Find(ctx context.Context, r ListRequest) (*AggregationResult, error) {
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

	activeUsers := s.websocketHub.GetUsers()
	for i := range res.Data {
		activeConnects := len(activeUsers[res.Data[i].ID])
		res.Data[i].ActiveConnects = &activeConnects
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

func (s *store) Update(ctx context.Context, r UpdateRequest) (*User, error) {
	var user *User
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		user = nil
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

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
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

func (s *store) BulkInsert(ctx context.Context, requests []CreateRequest) error {
	var err error
	writeModels := make([]mongodriver.WriteModel, 0, int(math.Min(float64(canopsis.DefaultBulkSize), float64(len(requests)))))

	for _, r := range requests {
		writeModels = append(
			writeModels,
			mongodriver.NewInsertOneModel().SetDocument(r.getBson(s.passwordEncoder)),
		)

		if len(writeModels) == canopsis.DefaultBulkSize {
			_, err = s.collection.BulkWrite(ctx, writeModels)
			if err != nil {
				return err
			}

			writeModels = writeModels[:0]
		}
	}

	if len(writeModels) > 0 {
		_, err = s.collection.BulkWrite(ctx, writeModels)
	}

	return err
}

func (s *store) BulkUpdate(ctx context.Context, requests []BulkUpdateRequestItem) error {
	var err error
	writeModels := make([]mongodriver.WriteModel, 0, int(math.Min(float64(canopsis.DefaultBulkSize), float64(len(requests)))))

	for _, r := range requests {
		writeModels = append(
			writeModels,
			mongodriver.
				NewUpdateOneModel().
				SetFilter(bson.M{"_id": r.ID, "crecord_type": securitymodel.LineTypeSubject}).
				SetUpdate(bson.M{"$set": r.getBson(s.passwordEncoder)}),
		)

		if len(writeModels) == canopsis.DefaultBulkSize {
			_, err = s.collection.BulkWrite(ctx, writeModels)
			if err != nil {
				return err
			}

			writeModels = writeModels[:0]
		}
	}

	if len(writeModels) > 0 {
		_, err = s.collection.BulkWrite(ctx, writeModels)
	}

	return err
}

func (s *store) BulkDelete(ctx context.Context, ids []string) error {
	_, err := s.collection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})

	return err
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
