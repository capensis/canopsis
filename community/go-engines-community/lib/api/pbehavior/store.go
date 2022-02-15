package pbehavior

import (
	"context"
	"encoding/json"
	"errors"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	libtypes "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const limitMatch = 100

type Store interface {
	Insert(ctx context.Context, model *Response) error
	Find(ctx context.Context, r ListRequest) (*AggregationResult, error)
	FindByEntityID(ctx context.Context, entityID string) ([]Response, error)
	GetOneBy(ctx context.Context, filter bson.M) (*Response, error)
	FindEntities(ctx context.Context, pbhID string, request EntitiesListRequest) (*AggregationEntitiesResult, error)
	Update(ctx context.Context, model *Response) (bool, error)
	UpdateByFilter(ctx context.Context, model *Response, filters bson.M) (bool, error)
	Delete(ctx context.Context, id string) (bool, error)
	Count(context.Context, Filter, int) (*CountFilterResult, error)
}

type store struct {
	dbClient mongo.DbClient

	dbCollection, entitiesCollection mongo.DbCollection

	entityMatcher          pbehavior.EntityMatcher
	entityTypeResolver     pbehavior.EntityTypeResolver
	timezoneConfigProvider config.TimezoneConfigProvider
	defaultSortBy          string

	entitiesDefaultSearchByFields []string
	entitiesDefaultSortBy         string
}

func NewStore(
	dbClient mongo.DbClient,
	entityMatcher pbehavior.EntityMatcher,
	entityTypeResolver pbehavior.EntityTypeResolver,
	timezoneConfigProvider config.TimezoneConfigProvider,
) Store {
	return &store{
		dbClient:                      dbClient,
		dbCollection:                  dbClient.Collection(mongo.PbehaviorMongoCollection),
		entitiesCollection:            dbClient.Collection(mongo.EntityMongoCollection),
		entityMatcher:                 entityMatcher,
		entityTypeResolver:            entityTypeResolver,
		timezoneConfigProvider:        timezoneConfigProvider,
		defaultSortBy:                 "created",
		entitiesDefaultSearchByFields: []string{"_id", "name", "type"},
		entitiesDefaultSortBy:         "_id",
	}
}

func (s *store) Insert(ctx context.Context, model *Response) error {
	now := libtypes.NewCpsTime(time.Now().Unix())

	if model.ID == "" {
		model.ID = utils.NewID()
	}

	model.Created = &now
	model.Updated = &now
	model.Comments = make([]*pbehavior.Comment, 0)

	doc, err := s.transformModelToDocument(*model)
	if err != nil {
		return err
	}

	doc.ID = model.ID
	doc.Created = now
	doc.Updated = now

	// If model.Stop is nil, insert to mongo using map so that
	// tstop field can be cleared
	if model.Stop == nil {
		m := make(map[string]interface{})
		var p []byte

		p, err = bson.Marshal(doc)
		if err != nil {
			return err
		}

		err = bson.Unmarshal(p, &m)
		if err != nil {
			return err
		}

		delete(m, "tstop")
		_, err = s.dbCollection.InsertOne(ctx, m)
	} else {
		_, err = s.dbCollection.InsertOne(ctx, doc)
	}

	return err
}

func (s *store) Find(ctx context.Context, r ListRequest) (*AggregationResult, error) {
	mongoQuery := CreateMongoQuery(s.dbClient)
	pipeline, err := mongoQuery.CreateAggregationPipeline(ctx, r)
	if err != nil {
		return nil, err
	}

	cursor, err := s.dbCollection.Aggregate(ctx, pipeline,
		options.Aggregate().SetCollation(&options.Collation{Locale: "en"}))

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

	err = s.fillActiveStatuses(ctx, result.Data)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *store) FindByEntityID(ctx context.Context, entityID string) ([]Response, error) {
	pbhIDs, err := s.getMatchedPbhIDs(ctx, entityID)
	if err != nil {
		return nil, err
	}

	pipeline := []bson.M{{"$match": bson.M{"_id": bson.M{"$in": pbhIDs}}}}
	pipeline = append(pipeline, GetNestedObjectsPipeline()...)
	pipeline = append(pipeline, common.GetSortQuery("created", common.SortAsc))
	cursor, err := s.dbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	res := make([]Response, 0)
	err = cursor.All(ctx, &res)
	if err != nil {
		return nil, err
	}

	err = s.fillActiveStatuses(ctx, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *store) getMatchedPbhIDs(ctx context.Context, entityID string) ([]string, error) {
	cursor, err := s.dbCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	pbhIDs := make([]string, 0)
	filters := make(map[string]string)

	for cursor.Next(ctx) {
		var pbh struct {
			ID     string `bson:"_id"`
			Filter string `bson:"filter"`
		}
		err := cursor.Decode(&pbh)
		if err != nil {
			return nil, err
		}

		filters[pbh.ID] = pbh.Filter
		if len(filters) == limitMatch {
			ids, err := s.entityMatcher.MatchAll(ctx, entityID, filters)
			if err != nil {
				return nil, err
			}

			pbhIDs = append(pbhIDs, ids...)
			filters = make(map[string]string)
		}
	}

	if len(filters) > 0 {
		ids, err := s.entityMatcher.MatchAll(ctx, entityID, filters)
		if err != nil {
			return nil, err
		}

		pbhIDs = append(pbhIDs, ids...)
	}

	return pbhIDs, nil
}

func (s *store) GetOneBy(ctx context.Context, filter bson.M) (*Response, error) {
	pipeline := []bson.M{
		{"$match": filter},
	}
	pipeline = append(pipeline, GetNestedObjectsPipeline()...)
	cursor, err := s.dbCollection.Aggregate(ctx, pipeline)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	if cursor.Next(ctx) {
		var pbh Response
		err = cursor.Decode(&pbh)
		if err != nil {
			return nil, err
		}

		return &pbh, nil
	}

	return nil, nil
}

func (s *store) FindEntities(ctx context.Context, pbhID string, request EntitiesListRequest) (*AggregationEntitiesResult, error) {
	pbh, err := s.GetOneBy(ctx, bson.M{"_id": pbhID})
	if err != nil || pbh == nil {
		return nil, err
	}

	pipeline := []bson.M{
		{"$match": pbh.Filter},
	}
	filter := common.GetSearchQuery(request.Search, s.entitiesDefaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	sortBy := request.SortBy
	if sortBy == "" {
		sortBy = s.entitiesDefaultSortBy
	}

	project := []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.EntityCategoryMongoCollection,
			"localField":   "category",
			"foreignField": "_id",
			"as":           "category",
		}},
		{"$unwind": bson.M{"path": "$category", "preserveNullAndEmptyArrays": true}},
	}
	cursor, err := s.entitiesCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		request.Query,
		pipeline,
		common.GetSortQuery(sortBy, request.Sort),
		project,
	))

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	cursor.Next(ctx)

	result := AggregationEntitiesResult{}
	err = cursor.Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *store) Update(ctx context.Context, model *Response) (bool, error) {
	if model == nil {
		return true, nil
	}

	prevDoc := pbehavior.PBehavior{}
	err := s.dbCollection.FindOne(ctx, bson.M{"_id": model.ID}).Decode(&prevDoc)
	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			return false, nil
		}

		return false, err
	}

	now := libtypes.NewCpsTime(time.Now().Unix())
	doc, err := s.transformModelToDocument(*model)
	if err != nil {
		return false, err
	}

	created := prevDoc.Created
	model.Comments = prevDoc.Comments
	model.LastAlarmDate = prevDoc.LastAlarmDate
	model.Created = &created
	model.Updated = &now
	doc.Updated = now

	// If model.Stop is nil, insert to mongo using map so that
	// tstop field can be cleared
	var update interface{}

	if model.Stop == nil {
		m := make(map[string]interface{})
		p, err := bson.Marshal(doc)
		if err != nil {
			return false, err
		}

		err = bson.Unmarshal(p, &m)
		if err != nil {
			return false, err
		}

		delete(m, "tstop")
		update = bson.M{
			"$set":   m,
			"$unset": bson.M{"tstop": 1},
		}
	} else {
		update = bson.M{"$set": doc}
	}

	result, err := s.dbCollection.UpdateOne(ctx, bson.M{"_id": model.ID}, update)
	if err != nil {
		return false, err
	}

	return result.MatchedCount > 0, nil
}

func (s *store) UpdateByFilter(ctx context.Context, model *Response, filters bson.M) (bool, error) {
	doc, err := s.transformModelToDocument(*model)
	if err != nil {
		return false, err
	}

	doc.Updated = libtypes.NewCpsTime(time.Now().Unix())

	var update bson.M
	if model.Stop == nil {
		m := make(map[string]interface{})
		p, err := bson.Marshal(doc)
		if err != nil {
			return false, err
		}

		err = bson.Unmarshal(p, &m)
		if err != nil {
			return false, err
		}

		delete(m, "tstop")
		update = bson.M{
			"$set":   m,
			"$unset": bson.M{"tstop": 1},
		}
	} else {
		update = bson.M{"$set": doc}
	}

	result, err := s.dbCollection.UpdateOne(
		ctx,
		filters,
		update,
	)
	if err != nil {
		return false, err
	}

	updatedModel, err := s.GetOneBy(ctx, bson.M{"_id": model.ID})
	if err != nil {
		return false, err
	}

	*model = *updatedModel

	return result.MatchedCount > 0, nil
}

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	deleted, err := s.dbCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}

	return deleted > 0, nil
}

func (s *store) transformModelToDocument(model Response) (pbehavior.PBehavior, error) {
	exdates := make([]pbehavior.Exdate, len(model.Exdates))
	for i := range model.Exdates {
		exdates[i].Type = model.Exdates[i].Type.ID
		exdates[i].Begin = model.Exdates[i].Begin
		exdates[i].End = model.Exdates[i].End
	}

	exceptions := make([]string, len(model.Exceptions))
	for i := range model.Exceptions {
		exceptions[i] = model.Exceptions[i].ID
	}

	filter, err := json.Marshal(model.Filter)
	if err != nil {
		return pbehavior.PBehavior{}, err
	}

	return pbehavior.PBehavior{
		Author:     model.Author,
		Enabled:    model.Enabled,
		Filter:     string(filter),
		Name:       model.Name,
		Reason:     model.Reason.ID,
		RRule:      model.RRule,
		Start:      model.Start,
		Stop:       model.Stop,
		Type:       model.Type.ID,
		Exdates:    exdates,
		Exceptions: exceptions,
	}, nil
}

func (s *store) fillActiveStatuses(ctx context.Context, result []Response) error {
	location := s.timezoneConfigProvider.Get().Location
	now := time.Now().In(location)
	ids := make([]string, len(result))
	for i, pbh := range result {
		ids[i] = pbh.ID
	}

	typesByID, err := s.entityTypeResolver.GetPbehaviors(ctx, ids, now)
	if err != nil {
		if errors.Is(err, pbehavior.ErrNoComputed) || errors.Is(err, pbehavior.ErrRecomputeNeed) {
			return nil
		}

		return err
	}

	for i := range result {
		_, ok := typesByID[result[i].ID]
		result[i].IsActiveStatus = &ok
	}

	return nil
}

func (s store) countEntitiesCollection(ctx context.Context, filter Filter) (int64, error) {
	cursor, err := s.entitiesCollection.Aggregate(ctx, []bson.M{
		{"$match": filter.v},
		{"$count": "total_count"},
	})
	if err != nil {
		return 0, err
	}

	defer cursor.Close(ctx)

	ar := AggregationResult{}
	if cursor.Next(ctx) {
		err = cursor.Decode(&ar)
	}
	return ar.GetTotal(), err
}

func (s store) Count(ctx context.Context, filter Filter, timeout int) (*CountFilterResult, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	res, err := s.countEntitiesCollection(ctx, filter)
	return &CountFilterResult{
		TotalCount: res,
	}, err
}
