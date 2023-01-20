package resolverule

import (
	"context"
	"errors"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/resolverule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type Store interface {
	Insert(ctx context.Context, r CreateRequest) (*Response, error)
	GetById(ctx context.Context, id string) (*Response, error)
	Find(ctx context.Context, query FilteredQuery) (*AggregationResult, error)
	Update(ctx context.Context, r UpdateRequest) (*Response, error)
	Delete(ctx context.Context, id string) (bool, error)
}

type store struct {
	dbClient     mongo.DbClient
	dbCollection mongo.DbCollection

	defaultSearchByFields []string
}

func NewStore(
	dbClient mongo.DbClient,
) Store {
	return &store{
		dbClient:     dbClient,
		dbCollection: dbClient.Collection(mongo.ResolveRuleMongoCollection),

		defaultSearchByFields: []string{"_id", "author.name", "name", "description"},
	}
}

func (s *store) Insert(ctx context.Context, request CreateRequest) (*Response, error) {
	now := types.NewCpsTime(time.Now().Unix())
	model := s.transformRequestToDocument(request.EditRequest)

	if request.ID == "" {
		request.ID = utils.NewID()
	}

	model.ID = request.ID
	model.Created = now
	model.Updated = now

	id := request.ID
	if id == "" {
		id = utils.NewID()
	}

	var res *Response
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		res = nil

		_, err := s.dbCollection.InsertOne(ctx, model)
		if err != nil {
			return err
		}

		err = s.updateFollowingPriorities(ctx, id, request.Priority)
		if err != nil {
			return err
		}

		res, err = s.GetById(ctx, id)
		return err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *store) GetById(ctx context.Context, id string) (*Response, error) {
	pipeline := []bson.M{{"$match": bson.M{"_id": id}}}
	pipeline = append(pipeline, author.Pipeline()...)

	cursor, err := s.dbCollection.Aggregate(ctx, pipeline)
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

func (s *store) Find(ctx context.Context, query FilteredQuery) (*AggregationResult, error) {
	pipeline := author.Pipeline()
	filter := common.GetSearchQuery(query.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	sortBy := "created"
	if query.SortBy != "" {
		sortBy = query.SortBy
	}

	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		query.Query,
		pipeline,
		common.GetSortQuery(sortBy, query.Sort),
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

func (s *store) Update(ctx context.Context, request UpdateRequest) (*Response, error) {
	model := s.transformRequestToDocument(request.EditRequest)
	model.Updated = types.CpsTime{Time: time.Now()}

	update := bson.M{"$set": model}

	unset := bson.M{}
	if request.CorporateAlarmPattern != "" || len(request.AlarmPattern) > 0 {
		unset["old_alarm_patterns"] = 1
	}

	if request.CorporateEntityPattern != "" || len(request.EntityPattern) > 0 {
		unset["old_entity_patterns"] = 1
	}

	if len(unset) > 0 {
		update["$unset"] = unset
	}

	var res *Response
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		res = nil
		originalModel := resolverule.Rule{}
		err := s.dbCollection.FindOne(ctx, bson.M{"_id": request.ID}).Decode(&originalModel)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}
			return err
		}

		_, err = s.dbCollection.UpdateOne(
			ctx,
			bson.M{"_id": request.ID},
			update,
		)
		if err != nil {
			return err
		}

		if originalModel.Priority != request.Priority {
			err := s.updateFollowingPriorities(ctx, request.ID, request.Priority)
			if err != nil {
				return err
			}
		}

		res, err = s.GetById(ctx, originalModel.ID)
		return err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	if id == resolverule.DefaultRule {
		return false, ErrDefaultRule
	}

	deleted, err := s.dbCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}

	return deleted > 0, nil
}

func (s *store) updateFollowingPriorities(ctx context.Context, id string, priority int) error {
	err := s.dbCollection.FindOne(ctx, bson.M{
		"_id":      bson.M{"$ne": id},
		"priority": priority,
	}).Err()
	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			return nil
		}
		return err
	}

	_, err = s.dbCollection.UpdateMany(
		ctx,
		bson.M{
			"_id":      bson.M{"$ne": id},
			"priority": bson.M{"$gte": priority},
		},
		bson.M{"$inc": bson.M{"priority": 1}},
	)

	return err
}

func (s *store) transformRequestToDocument(r EditRequest) resolverule.Rule {
	return resolverule.Rule{
		Name:        r.Name,
		Description: r.Description,
		Duration:    r.Duration,
		AlarmPatternFields: r.AlarmPatternFieldsRequest.ToModelWithoutFields(
			common.GetForbiddenFieldsInAlarmPattern(mongo.ResolveRuleMongoCollection),
			common.GetOnlyAbsoluteTimeCondFieldsInAlarmPattern(mongo.ResolveRuleMongoCollection),
		),
		EntityPatternFields: r.EntityPatternFieldsRequest.ToModelWithoutFields(
			common.GetForbiddenFieldsInEntityPattern(mongo.ResolveRuleMongoCollection),
		),
		Priority: r.Priority,
		Author:   r.Author,
	}
}
