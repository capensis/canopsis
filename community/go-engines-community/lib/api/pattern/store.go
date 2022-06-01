package pattern

import (
	"context"
	"errors"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"
)

type Store interface {
	Insert(ctx context.Context, r EditRequest) (*Response, error)
	GetById(ctx context.Context, id, userId string) (*Response, error)
	Find(ctx context.Context, r ListRequest, userId string) (*AggregationResult, error)
	Update(ctx context.Context, r EditRequest) (*Response, error)
	Delete(ctx context.Context, pattern Response) (bool, error)
	Count(ctx context.Context, r CountRequest, maxCount int64) (CountResponse, error)
}

type store struct {
	client     mongo.DbClient
	collection mongo.DbCollection

	linkedCollections []string

	defaultSearchByFields []string
	defaultSortBy         string

	pbhComputeChan chan<- pbehavior.ComputeTask

	logger zerolog.Logger
}

func NewStore(
	dbClient mongo.DbClient,
	pbhComputeChan chan<- pbehavior.ComputeTask,
	logger zerolog.Logger,
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
			mongo.InstructionMongoCollection,
			mongo.PbehaviorMongoCollection,
		},

		pbhComputeChan: pbhComputeChan,

		logger: logger,
	}
}

func (s *store) Insert(ctx context.Context, request EditRequest) (*Response, error) {
	now := types.NewCpsTime()
	model := transformRequestToModel(request)
	model.ID = utils.NewID()
	model.Created = now
	model.Updated = now

	var response *Response
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil
		_, err := s.collection.InsertOne(ctx, model)
		if err != nil {
			return err
		}

		response, err = s.GetById(ctx, model.ID, model.Author)
		return err
	})

	return response, err
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

	var response *Response
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil
		res, err := s.collection.UpdateOne(
			ctx,
			bson.M{"_id": request.ID},
			bson.M{"$set": model},
		)
		if err != nil || res.MatchedCount == 0 {
			return err
		}

		response, err = s.GetById(ctx, model.ID, model.Author)
		if err != nil || response == nil {
			return err
		}

		err = s.updateLinkedModels(ctx, *response)
		return err
	})

	return response, err
}

func (s *store) Delete(ctx context.Context, pattern Response) (bool, error) {
	res := false
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		res = false
		deleted, err := s.collection.DeleteOne(ctx, bson.M{"_id": pattern.ID})
		if err != nil || deleted == 0 {
			return err
		}

		err = s.cleanLinkedModels(ctx, pattern)
		if err != nil {
			return err
		}

		res = true
		return nil
	})

	return res, err
}

func (s *store) updateLinkedModels(ctx context.Context, pattern Response) error {
	if !pattern.IsCorporate {
		return nil
	}

	var filter, set bson.M
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
			"pbehavior_pattern":                 pattern.PbehaviorPattern,
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

	s.processPbehaviors(ctx, pattern)

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

func (s *store) Count(ctx context.Context, r CountRequest, maxCount int64) (CountResponse, error) {
	res := CountResponse{}
	g, ctx := errgroup.WithContext(ctx)
	var err error
	var alarmPatternQuery, entityPatternQuery, pbhPatternQuery bson.M
	var alarmPatternCount, entityPatternCount, pbhPatternCount Count
	if len(r.AlarmPattern) > 0 {
		alarmPatternQuery, err = r.AlarmPattern.ToMongoQuery("")
		if err != nil {
			return res, err
		}
	}
	if len(r.PbehaviorPattern) > 0 {
		pbhPatternQuery, err = r.PbehaviorPattern.ToMongoQuery("v")
		if err != nil {
			return res, err
		}
	}
	if len(r.EntityPattern) > 0 {
		entityPatternQuery, err = r.EntityPattern.ToMongoQuery("")
		if err != nil {
			return res, err
		}
	}

	if len(alarmPatternQuery) > 0 {
		g.Go(func() error {
			alarmPatternCount.Count, err = s.fetchCount(ctx, s.client.Collection(mongo.AlarmMongoCollection), alarmPatternQuery)
			alarmPatternCount.OverLimit = alarmPatternCount.Count > maxCount

			return err
		})
	}
	if len(pbhPatternQuery) > 0 {
		g.Go(func() error {
			var err error
			pbhPatternCount.Count, err = s.fetchCount(ctx, s.client.Collection(mongo.AlarmMongoCollection), pbhPatternQuery)
			pbhPatternCount.OverLimit = pbhPatternCount.Count > maxCount

			return err
		})
	}
	if len(entityPatternQuery) > 0 {
		g.Go(func() error {
			var err error
			entityPatternCount.Count, err = s.fetchCount(ctx, s.client.Collection(mongo.EntityMongoCollection), entityPatternQuery)
			entityPatternCount.OverLimit = entityPatternCount.Count > maxCount

			return err
		})
	}
	if err := g.Wait(); err != nil {
		return res, err
	}

	res.AlarmPattern = alarmPatternCount
	res.PbehaviorPattern = pbhPatternCount
	res.EntityPattern = entityPatternCount

	return res, nil
}

func (s *store) fetchCount(ctx context.Context, collection mongo.DbCollection, match bson.M) (int64, error) {
	cursor, err := collection.Aggregate(ctx, []bson.M{
		{"$match": match},
		{"$count": "total_count"},
	})
	if err != nil {
		return 0, err
	}

	defer cursor.Close(ctx)

	res := AggregationResult{}
	if cursor.Next(ctx) {
		err = cursor.Decode(&res)
	}

	return res.GetTotal(), err
}

func (s *store) processPbehaviors(ctx context.Context, pattern Response) {
	if pattern.Type != savedpattern.TypeEntity {
		return
	}

	cursor, err := s.client.Collection(mongo.PbehaviorMongoCollection).Find(ctx, bson.M{
		"corporate_entity_pattern": pattern.ID,
	}, options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		s.logger.Err(err).Msg("cannot fetch pbehaviors")
		return
	}
	defer cursor.Close(ctx)

	pbhIds := make([]string, 0)
	for cursor.Next(ctx) {
		pbh := pbehavior.PBehavior{}
		err := cursor.Decode(&pbh)
		if err != nil {
			s.logger.Err(err).Msg("cannot decode pbehavior")
			continue
		}
		pbhIds = append(pbhIds, pbh.ID)
	}

	if len(pbhIds) == 0 {
		return
	}

	task := pbehavior.ComputeTask{
		PbehaviorIds: pbhIds,
	}
	select {
	case s.pbhComputeChan <- task:
	default:
		s.logger.Err(errors.New("channel is full")).
			Strs("pbehavior", task.PbehaviorIds).
			Msg("fail to start pbehavior recompute")
	}
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
