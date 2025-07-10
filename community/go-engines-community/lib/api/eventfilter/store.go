package eventfilter

import (
	"cmp"
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	apiexternaldata "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/externaldatatable"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/priority"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/usernotification"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Store interface {
	Insert(ctx context.Context, request CreateRequest) (*Response, error)
	GetByID(ctx context.Context, id string) (*Response, error)
	Find(ctx context.Context, query FilteredQuery) (*AggregationResult, error)
	Update(ctx context.Context, request UpdateRequest) (*Response, error)
	Delete(ctx context.Context, id, userID string) (bool, error)
	FindFailures(ctx context.Context, id string, r FailureRequest) (*AggregationFailureResult, error)
	ReadFailures(ctx context.Context, id string) (bool, error)
}

type store struct {
	dbClient                mongo.DbClient
	dbCollection            mongo.DbCollection
	dbFailureCollection     mongo.DbCollection
	dbExdataTableCollection mongo.DbCollection
	authorProvider          author.Provider
	notificationStore       usernotification.Store
	defaultSearchByFields   []string
	defaultSortBy           string
}

func NewStore(
	dbClient mongo.DbClient,
	authorProvider author.Provider,
	notificationStore usernotification.Store,
) Store {
	return &store{
		dbClient:                dbClient,
		dbCollection:            dbClient.Collection(mongo.EventFilterRuleCollection),
		dbFailureCollection:     dbClient.Collection(mongo.EventFilterFailureCollection),
		dbExdataTableCollection: dbClient.Collection(mongo.ExternalDataTableCollection),
		authorProvider:          authorProvider,
		notificationStore:       notificationStore,
		defaultSearchByFields:   []string{"_id", "author.name", "description", "type"},
		defaultSortBy:           "created",
	}
}

func (s *store) Insert(ctx context.Context, request CreateRequest) (*Response, error) {
	model, err := s.transformRequestToDocument(ctx, request.EditRequest)
	if err != nil {
		return nil, err
	}

	model.ID = request.ID
	if model.ID == "" {
		model.ID = utils.NewID()
	}

	now := datetime.NewCpsTime()
	model.Created = &now
	model.Updated = &now

	var response *Response
	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil
		_, err := s.dbCollection.InsertOne(ctx, model)
		if err != nil {
			return err
		}

		err = priority.UpdateFollowing(ctx, s.dbCollection, model.ID, model.Priority)
		if err != nil {
			return err
		}

		response, err = s.GetByID(ctx, model.ID)
		return err
	})

	return response, err
}

func (s *store) GetByID(ctx context.Context, id string) (*Response, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"_id": id}},
	}
	pipeline = append(pipeline, s.getResponseLookups()...)
	pipeline = append(pipeline, s.authorProvider.Pipeline()...)

	cursor, err := s.dbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	if cursor.Next(ctx) {
		var rule Response
		err = cursor.Decode(&rule)
		if err != nil {
			return nil, err
		}

		return &rule, nil
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *store) Find(ctx context.Context, query FilteredQuery) (*AggregationResult, error) {
	pipeline := s.authorProvider.Pipeline()
	andCond := make([]bson.M, 0)
	filter := common.GetSearchQuery(query.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		andCond = append(andCond, filter)
	}

	if query.OnlyUnreadFailure {
		andCond = append(andCond, bson.M{"unread_failures_count": bson.M{"$gt": 0}})
	}

	if len(andCond) > 0 {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"$and": andCond}})
	}

	sort := common.GetSortQuery(cmp.Or(query.SortBy, s.defaultSortBy), query.Sort)
	project := s.getResponseLookups()
	project = append(project, sort)
	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		query.Query,
		pipeline,
		sort,
		project,
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

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *store) Update(ctx context.Context, request UpdateRequest) (*Response, error) {
	updated := datetime.NewCpsTime()
	model, err := s.transformRequestToDocument(ctx, request.EditRequest)
	if err != nil {
		return nil, err
	}

	model.ID = request.ID
	model.Created = nil
	model.Updated = &updated

	update := bson.M{"$set": model}
	unset := bson.M{
		"events_count":          "",
		"unread_failures_count": "",
	}

	if model.Start == nil || model.Start.IsZero() || model.Stop == nil || model.Stop.IsZero() {
		unset["start"] = ""
		unset["stop"] = ""
		unset["resolved_start"] = ""
		unset["resolved_stop"] = ""
		unset["next_resolved_start"] = ""
		unset["next_resolved_stop"] = ""
	}

	if len(unset) != 0 {
		update["$unset"] = unset
	}

	var response *Response
	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil
		_, err := s.dbCollection.UpdateOne(
			ctx,
			bson.M{"_id": model.ID},
			update,
		)
		if err != nil {
			return err
		}

		err = priority.UpdateFollowing(ctx, s.dbCollection, model.ID, model.Priority)
		if err != nil {
			return err
		}

		response, err = s.GetByID(ctx, model.ID)
		return err
	})
	if err != nil || response == nil {
		return nil, err
	}

	_, err = s.dbFailureCollection.UpdateMany(ctx, bson.M{"rule": request.ID, "unread": true}, bson.M{
		"$unset": bson.M{
			"unread": "",
		},
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *store) Delete(ctx context.Context, id, userID string) (bool, error) {
	_, err := s.dbFailureCollection.DeleteMany(ctx, bson.M{"rule": id})
	if err != nil {
		return false, err
	}

	var deleted int64

	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		deleted = 0

		// required to get the author in action log listener.
		res, err := s.dbCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"author": userID}})
		if err != nil || res.MatchedCount == 0 {
			return err
		}

		deleted, err = s.dbCollection.DeleteOne(ctx, bson.M{"_id": id})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil || deleted == 0 {
		return false, err
	}

	err = s.notificationStore.DeleteForEventFilterFailure(ctx, id)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *store) FindFailures(ctx context.Context, id string, r FailureRequest) (*AggregationFailureResult, error) {
	err := s.dbCollection.FindOne(ctx, bson.M{"_id": id}, options.FindOne().SetProjection(bson.M{"_id": 1})).Err()
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	match := bson.M{"rule": id}
	if r.Type != nil {
		match["type"] = r.Type
	}

	cursor, err := s.dbFailureCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		[]bson.M{{"$match": match}},
		common.GetSortQuery("t", common.SortDesc),
	))

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	var result AggregationFailureResult
	if cursor.Next(ctx) {
		err = cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *store) ReadFailures(ctx context.Context, id string) (bool, error) {
	ruleExists := false
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		ruleExists = false
		err := s.dbCollection.
			FindOneAndUpdate(
				ctx,
				bson.M{"_id": id},
				bson.M{"$unset": bson.M{"unread_failures_count": ""}},
				options.FindOneAndUpdate().SetProjection(bson.M{"_id": 1}),
			).
			Err()
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		ruleExists = true
		_, err = s.dbFailureCollection.UpdateMany(ctx, bson.M{"rule": id, "unread": true}, bson.M{"$unset": bson.M{
			"unread": "",
		}})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil || !ruleExists {
		return false, err
	}

	err = s.notificationStore.DeleteForEventFilterFailure(ctx, id)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *store) transformRequestToDocument(ctx context.Context, r EditRequest) (eventfilter.Rule, error) {
	exdates := make([]types.Exdate, len(r.Exdates))
	for i := range r.Exdates {
		exdates[i].Begin = r.Exdates[i].Begin
		exdates[i].End = r.Exdates[i].End
	}

	externalData, err := apiexternaldata.TransformRefParameters(ctx, r.ExternalData, s.dbExdataTableCollection)
	if err != nil {
		return eventfilter.Rule{}, err
	}

	return eventfilter.Rule{
		Author:              r.Author,
		Description:         r.Description,
		Type:                r.Type,
		Priority:            r.Priority,
		Enabled:             r.Enabled,
		Config:              r.Config,
		ExternalData:        externalData,
		EventPattern:        r.EventPattern,
		EntityPatternFields: r.EntityPatternFieldsRequest.ToModel(),
		RRule:               r.RRule,
		Start:               r.Start,
		Stop:                r.Stop,
		ResolvedStart:       r.Start,
		ResolvedStop:        r.Stop,
		Exdates:             exdates,
		Exceptions:          r.Exceptions,
	}, nil
}

func (s *store) getResponseLookups() []bson.M {
	pipeline := apiexternaldata.GetRefParametersLookups()
	pipeline = append(pipeline, bson.M{
		"$lookup": bson.M{
			"from":         mongo.PbehaviorExceptionMongoCollection,
			"localField":   "exceptions",
			"foreignField": "_id",
			"as":           "exceptions",
		},
	})

	return pipeline
}
