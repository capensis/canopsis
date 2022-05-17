package pattern

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type Store interface {
	Insert(ctx context.Context, r EditRequest) (*Response, error)
	GetById(ctx context.Context, id, userId string) (*Response, error)
	Find(ctx context.Context, r ListRequest, userId string) (*AggregationResult, error)
	Update(ctx context.Context, r EditRequest) (*Response, error)
	Delete(ctx context.Context, pattern Response) (bool, error)
}

type store struct {
	client     mongo.DbClient
	collection mongo.DbCollection

	linkedCollections []string

	defaultSearchByFields []string
	defaultSortBy         string
}

func NewStore(
	dbClient mongo.DbClient,
) Store {
	return &store{
		client:     dbClient,
		collection: dbClient.Collection(mongo.PatternMongoCollection),

		defaultSearchByFields: []string{"_id", "author", "title"},
		defaultSortBy:         "created",

		linkedCollections: []string{
			mongo.WidgetFiltersMongoCollection,
			mongo.EventFilterRulesMongoCollection,
			mongo.MetaAlarmRulesMongoCollection,
		},
	}
}

func (s *store) Insert(ctx context.Context, request EditRequest) (*Response, error) {
	now := types.NewCpsTime()
	model := transformRequestToModel(request)
	model.ID = utils.NewID()
	model.Created = now
	model.Updated = now

	_, err := s.collection.InsertOne(ctx, model)
	if err != nil {
		return nil, err
	}

	return s.GetById(ctx, model.ID, model.Author)
}

func (s *store) GetById(ctx context.Context, id, userId string) (*Response, error) {
	pipeline := []bson.M{{"$match": bson.M{
		"_id": id,
		"$or": []bson.M{
			{"author": userId},
			{"is_corporate": true},
		},
	}}}
	pipeline = append(pipeline, getAuthorPipeline()...)
	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	if cursor.Next(ctx) {
		response := Response{}
		err := cursor.Decode(&response)
		if err != nil {
			return nil, err
		}

		return &response, nil
	}

	return nil, nil
}

func (s *store) Find(ctx context.Context, request ListRequest, userId string) (*AggregationResult, error) {
	pipeline := make([]bson.M, 0)
	match := make([]bson.M, 0)

	if request.Corporate == nil {
		match = append(match, bson.M{"$or": []bson.M{
			{"author": userId},
			{"is_corporate": true},
		}})
	} else if *request.Corporate {
		match = append(match, bson.M{"is_corporate": true})
	} else {
		match = append(match, bson.M{"author": userId, "is_corporate": false})
	}

	if request.Type != "" {
		match = append(match, bson.M{"type": request.Type})
	}

	if len(match) > 0 {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"$and": match}})
	}

	filter := common.GetSearchQuery(request.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	sortBy := s.defaultSortBy
	if request.SortBy != "" {
		sortBy = request.SortBy
	}

	cursor, err := s.collection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		request.Query,
		pipeline,
		common.GetSortQuery(sortBy, request.Sort),
		getAuthorPipeline(),
	))

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	var result AggregationResult
	if cursor.Next(ctx) {
		err = cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
	}

	return &result, nil
}

func (s *store) Update(ctx context.Context, request EditRequest) (*Response, error) {
	now := types.NewCpsTime()
	model := transformRequestToModel(request)
	model.ID = request.ID
	model.Updated = now

	res, err := s.collection.UpdateOne(
		ctx,
		bson.M{"_id": request.ID},
		bson.M{"$set": model},
	)
	if err != nil || res.MatchedCount == 0 {
		return nil, err
	}

	pattern, err := s.GetById(ctx, model.ID, model.Author)
	if err != nil || pattern == nil {
		return nil, err
	}

	err = s.updateLinkedModels(ctx, *pattern)
	if err != nil {
		return nil, err
	}

	return pattern, nil
}

func (s *store) Delete(ctx context.Context, pattern Response) (bool, error) {
	deleted, err := s.collection.DeleteOne(ctx, bson.M{"_id": pattern.ID})
	if err != nil || deleted == 0 {
		return false, err
	}

	err = s.cleanLinkedModels(ctx, pattern)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *store) updateLinkedModels(ctx context.Context, pattern Response) error {
	if !pattern.IsCorporate {
		return nil
	}

	filter := bson.M{}
	set := bson.M{}
	switch pattern.Type {
	case savedpattern.TypeAlarm:
		filter = bson.M{"corporate_alarm_pattern": pattern.ID}
		set = bson.M{
			"alarm_pattern":                 pattern.AlarmPattern,
			"corporate_alarm_pattern_title": pattern.Title,
		}
	case savedpattern.TypeEntity:
		filter = bson.M{"corporate_entity_pattern": pattern.ID}
		set = bson.M{
			"entity_pattern":                 pattern.EntityPattern,
			"corporate_entity_pattern_title": pattern.Title,
		}
	case savedpattern.TypePbehavior:
		filter = bson.M{"corporate_pbehavior_pattern": pattern.ID}
		set = bson.M{
			"pbehavior_pattern":               pattern.PbehaviorPattern,
			"corporate_pbehavior_pattern_title": pattern.Title,
		}
	default:
		return fmt.Errorf("unknown pattern type id=%s: %q", pattern.ID, pattern.Type)
	}

	if pattern.Type == savedpattern.TypeEntity {
		_, err := s.client.Collection(mongo.MetaAlarmRulesMongoCollection).UpdateMany(ctx, bson.M{"corporate_total_entity_pattern": pattern.ID}, bson.M{
			"$set": bson.M{
				"total_entity_pattern":                 pattern.EntityPattern,
				"corporate_total_entity_pattern_title": pattern.Title,
			},
		})
		if err != nil {
			return err
		}
	}

	for _, collection := range s.linkedCollections {
		_, err := s.client.Collection(collection).UpdateMany(ctx, filter, bson.M{
			"$set": set,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *store) cleanLinkedModels(ctx context.Context, pattern Response) error {
	if !pattern.IsCorporate {
		return nil
	}

	f := ""
	switch pattern.Type {
	case savedpattern.TypeAlarm:
		f = "corporate_alarm_pattern"
	case savedpattern.TypeEntity:
		f = "corporate_entity_pattern"
	case savedpattern.TypePbehavior:
		f = "corporate_pbehavior_pattern"
	default:
		return fmt.Errorf("unknown pattern type for deleted pattern id=%s: %q", pattern.ID, pattern.Type)
	}

	if pattern.Type == savedpattern.TypeEntity {
		_, err := s.client.Collection(mongo.MetaAlarmRulesMongoCollection).UpdateMany(ctx, bson.M{"corporate_total_entity_pattern": pattern.ID}, bson.M{
			"$unset": bson.M{
				"corporate_total_entity_pattern":       "",
				"corporate_total_entity_pattern_title": "",
			},
		})
		if err != nil {
			return err
		}
	}

	for _, collection := range s.linkedCollections {
		_, err := s.client.Collection(collection).UpdateMany(ctx, bson.M{f: pattern.ID}, bson.M{
			"$unset": bson.M{
				f:            "",
				f + "_title": "",
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func getAuthorPipeline() []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.RightsMongoCollection,
			"localField":   "author",
			"foreignField": "_id",
			"as":           "author",
		}},
		{"$unwind": bson.M{"path": "$author", "preserveNullAndEmptyArrays": true}},
	}
}

func transformRequestToModel(request EditRequest) savedpattern.SavedPattern {
	model := savedpattern.SavedPattern{
		Title:       request.Title,
		Type:        request.Type,
		IsCorporate: *request.IsCorporate,
		Author:      request.Author,
	}

	switch request.Type {
	case savedpattern.TypeAlarm:
		model.AlarmPattern = request.AlarmPattern
	case savedpattern.TypeEntity:
		model.EntityPattern = request.EntityPattern
	case savedpattern.TypePbehavior:
		model.PbehaviorPattern = request.PbehaviorPattern
	}

	return model
}
