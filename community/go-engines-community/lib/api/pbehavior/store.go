package pbehavior

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/api/common"
	"git.canopsis.net/canopsis/go-engines/lib/api/pagination"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/config"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/pbehavior"
	libtypes "git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"git.canopsis.net/canopsis/go-engines/lib/redis"
	"git.canopsis.net/canopsis/go-engines/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const limitMatch = 100

type Store interface {
	Insert(ctx context.Context, model *PBehavior) error
	Find(ctx context.Context, r ListRequest) (*AggregationResult, error)
	FindByEntityID(ctx context.Context, entityID string) ([]PBehavior, error)
	GetOneBy(ctx context.Context, filter bson.M) (*PBehavior, error)
	GetEIDs(ctx context.Context, pbhID string, request EIDsListRequest) (AggregationEIDsResult, error)
	Update(ctx context.Context, model *PBehavior) (bool, error)
	Delete(ctx context.Context, id string) (bool, error)
	Count(context.Context, Filter, int) (*CountFilterResult, error)
}

type store struct {
	dbClient mongo.DbClient

	dbCollection, entitiesCollection mongo.DbCollection

	entityMatcher          pbehavior.EntityMatcher
	redisStore             redis.Store
	service                pbehavior.Service
	timezoneConfigProvider config.TimezoneConfigProvider
	defaultSortBy          string
}

func NewStore(
	dbClient mongo.DbClient,
	entityMatcher pbehavior.EntityMatcher,
	redisStore redis.Store,
	service pbehavior.Service,
	timezoneConfigProvider config.TimezoneConfigProvider,
) Store {
	return &store{
		dbClient:               dbClient,
		dbCollection:           dbClient.Collection(mongo.PbehaviorMongoCollection),
		entitiesCollection:     dbClient.Collection(mongo.EntityMongoCollection),
		entityMatcher:          entityMatcher,
		redisStore:             redisStore,
		service:                service,
		timezoneConfigProvider: timezoneConfigProvider,
		defaultSortBy:          "created",
	}
}

func (s *store) Insert(ctx context.Context, model *PBehavior) error {
	if model.ID == "" {
		model.ID = utils.NewID()
	}

	now := libtypes.NewCpsTime(time.Now().Unix())
	doc, err := s.transformModelToDocument(model)
	if err != nil {
		return err
	}

	doc.ID = model.ID
	doc.Created = now
	doc.Updated = now
	_, err = s.dbCollection.InsertOne(ctx, doc)
	if err != nil {
		return err
	}

	model.Created = &now
	model.Updated = &now
	model.Comments = make(pbehavior.Comments, 0)

	return nil
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

func (s *store) FindByEntityID(ctx context.Context, entityID string) ([]PBehavior, error) {
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

	res := make([]PBehavior, 0)
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
			ids, err := s.getMatchedPbhIDsByFilters(ctx, entityID, filters)
			if err != nil {
				return nil, err
			}

			pbhIDs = append(pbhIDs, ids...)
			filters = make(map[string]string)
		}
	}

	if len(filters) > 0 {
		ids, err := s.getMatchedPbhIDsByFilters(ctx, entityID, filters)
		if err != nil {
			return nil, err
		}

		pbhIDs = append(pbhIDs, ids...)
	}

	return pbhIDs, nil
}

func (s *store) getMatchedPbhIDsByFilters(
	ctx context.Context,
	entityID string,
	filters map[string]string,
) ([]string, error) {
	pbhIDs := make([]string, 0)
	match, err := s.entityMatcher.MatchAll(ctx, entityID, filters)
	if err != nil {
		return nil, err
	}

	for pbhID, matched := range match {
		if matched {
			pbhIDs = append(pbhIDs, pbhID)
		}
	}

	return pbhIDs, nil
}

func (s *store) GetOneBy(ctx context.Context, filter bson.M) (*PBehavior, error) {
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
		var pbh PBehavior
		err = cursor.Decode(&pbh)
		if err != nil {
			return nil, err
		}

		return &pbh, nil
	}

	return nil, nil
}

func (s *store) GetEIDs(ctx context.Context, pbhID string, request EIDsListRequest) (AggregationEIDsResult, error) {
	var filter bson.M

	result := AggregationEIDsResult{
		Data:       make([]EID, 0),
		TotalCount: 0,
	}

	if request.Search != "" {
		searchRegexp := primitive.Regex{
			Pattern: fmt.Sprintf(".*%s.*", request.Search),
			Options: "i",
		}

		filter = bson.M{"d": searchRegexp}
	} else {
		filter = bson.M{}
	}

	sortBy := request.SortBy
	if sortBy == "" {
		sortBy = "t"
	}

	collection := s.dbClient.Collection(mongo.AlarmMongoCollection)
	pipeline := pagination.CreateAggregationPipeline(
		request.Query,
		[]bson.M{
			{
				"$match": bson.M{
					"$and": []bson.M{
						{"v.pbehavior_info.id": pbhID},
						filter,
					},
				},
			},
			{
				"$project": bson.M{
					"id": "$d",
					"t":  1,
				},
			},
		},
		common.GetSortQuery(sortBy, request.Sort),
	)
	cursor, err := collection.Aggregate(ctx, pipeline)

	if err != nil {
		return result, err
	}

	defer cursor.Close(ctx)
	cursor.Next(ctx)

	err = cursor.Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (s *store) Update(ctx context.Context, model *PBehavior) (bool, error) {
	doc, err := s.transformModelToDocument(model)
	if err != nil {
		return false, err
	}

	doc.Updated = libtypes.NewCpsTime(time.Now().Unix())
	result, err := s.dbCollection.UpdateOne(
		ctx,
		bson.M{"_id": model.ID},
		bson.M{"$set": doc},
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

func (s *store) transformModelToDocument(model *PBehavior) (*pbehavior.PBehavior, error) {
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
		return nil, err
	}

	return &pbehavior.PBehavior{
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

func (s *store) fillActiveStatuses(ctx context.Context, result []PBehavior) error {
	ok, err := s.redisStore.Restore(ctx, s.service)
	if err != nil {
		return err
	}

	if !ok {
		return nil
	}

	location := s.timezoneConfigProvider.Get().Location
	now := time.Now().In(location)
	ids := make([]string, len(result))
	for i, pbh := range result {
		ids[i] = pbh.ID
	}

	statusesByID, err := s.service.GetPbehaviorStatus(context.Background(), ids, now)
	if err != nil {
		return err
	}

	for i := range result {
		v := statusesByID[result[i].ID]
		result[i].IsActiveStatus = &v
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
