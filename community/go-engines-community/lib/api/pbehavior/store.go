package pbehavior

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/api/common"
	"git.canopsis.net/canopsis/go-engines/lib/api/pagination"
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
	Insert(model *PBehavior) error
	Find(r ListRequest) (*AggregationResult, error)
	FindByEntityID(entityID string) ([]PBehavior, error)
	GetOneBy(filter bson.M) (*PBehavior, error)
	GetEIDs(pbhID string, request EIDsListRequest) (AggregationEIDsResult, error)
	Update(model *PBehavior) (bool, error)
	Delete(id string) (bool, error)
}

type store struct {
	dbClient      mongo.DbClient
	dbCollection  mongo.DbCollection
	entityMatcher pbehavior.EntityMatcher
	redisStore    redis.Store
	service       pbehavior.Service
	location      *time.Location
}

func NewStore(
	dbClient mongo.DbClient,
	entityMatcher pbehavior.EntityMatcher,
	redisStore redis.Store,
	service pbehavior.Service,
	location *time.Location,
) Store {
	return &store{
		dbClient:      dbClient,
		dbCollection:  dbClient.Collection(mongo.PbehaviorMongoCollection),
		entityMatcher: entityMatcher,
		redisStore:    redisStore,
		service:       service,
		location:      location,
	}
}

func (s *store) Insert(model *PBehavior) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

func (s *store) Find(r ListRequest) (*AggregationResult, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	err = s.fillActiveStatuses(result.Data)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *store) FindByEntityID(entityID string) ([]PBehavior, error) {
	pbhIDs, err := s.getMatchedPbhIDs(entityID)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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

	err = s.fillActiveStatuses(res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *store) getMatchedPbhIDs(entityID string) ([]string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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

func (s *store) GetOneBy(filter bson.M) (*PBehavior, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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

func (s *store) GetEIDs(pbhID string, request EIDsListRequest) (AggregationEIDsResult, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

func (s *store) Update(model *PBehavior) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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

	updatedModel, err := s.GetOneBy(bson.M{"_id": model.ID})
	if err != nil {
		return false, err
	}

	*model = *updatedModel

	return result.MatchedCount > 0, nil
}

func (s *store) Delete(id string) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

func (s *store) fillActiveStatuses(result []PBehavior) error {
	ok, err := s.redisStore.Restore(s.service)
	if err != nil {
		return err
	}

	if !ok {
		return nil
	}

	now := time.Now().In(s.location)
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
