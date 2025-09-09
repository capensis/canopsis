package scenario

import (
	"cmp"
	"context"
	"errors"
	"strconv"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/priority"
	libaction "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/action"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Store interface {
	Insert(ctx context.Context, r CreateRequest) (*Scenario, error)
	Find(ctx context.Context, q FilteredQuery) (*AggregationResult, error)
	GetOneBy(ctx context.Context, id string) (*Scenario, error)
	Update(ctx context.Context, r UpdateRequest) (*Scenario, error)
	Delete(ctx context.Context, id, userID string) (bool, error)
}

type store struct {
	dbClient              mongo.DbClient
	collection            mongo.DbCollection
	transformer           common.PatternFieldsTransformer
	authorProvider        author.Provider
	defaultSearchByFields []string
	defaultSortBy         string
}

func NewStore(db mongo.DbClient, authorProvider author.Provider, transformer common.PatternFieldsTransformer) Store {
	return &store{
		dbClient:              db,
		collection:            db.Collection(mongo.ScenarioCollection),
		transformer:           transformer,
		authorProvider:        authorProvider,
		defaultSearchByFields: []string{"_id", "name", "author.name"},
		defaultSortBy:         "created",
	}
}

func (s *store) Find(ctx context.Context, r FilteredQuery) (*AggregationResult, error) {
	pipeline := s.authorProvider.Pipeline()
	filter := common.GetSearchQuery(r.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	pipeline = append(pipeline, getNestedObjectsPipeline()...)
	cursor, err := s.collection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		pipeline,
		s.getSort(r),
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

	return &res, nil
}

func (s *store) GetOneBy(ctx context.Context, id string) (*Scenario, error) {
	pipeline := []bson.M{{"$match": bson.M{"_id": id}}}
	pipeline = append(pipeline, getNestedObjectsPipeline()...)
	pipeline = append(pipeline, s.authorProvider.Pipeline()...)

	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	if cursor.Next(ctx) {
		res := &Scenario{}
		err := cursor.Decode(res)
		if err != nil {
			return nil, err
		}

		return res, nil
	}

	return nil, nil
}

func (s *store) Insert(ctx context.Context, r CreateRequest) (*Scenario, error) {
	now := datetime.NewCpsTime()

	if r.ID == "" {
		r.ID = utils.NewID()
	}

	model := s.transformEditRequestToModel(r.EditRequest)
	model.ID = r.ID
	model.Created = now
	model.Updated = now

	var result *Scenario
	var err error

	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		result = nil

		model.Actions, model.Aliases, err = s.transformActionRequestToModel(ctx, r.Actions)
		if err != nil {
			return err
		}

		_, err := s.collection.InsertOne(ctx, model)
		if err != nil {
			return err
		}

		err = priority.UpdateFollowing(ctx, s.collection, model.ID, model.Priority)
		if err != nil {
			return err
		}

		result, err = s.GetOneBy(ctx, model.ID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *store) Update(ctx context.Context, r UpdateRequest) (*Scenario, error) {
	now := datetime.NewCpsTime()

	model := s.transformEditRequestToModel(r.EditRequest)
	model.Updated = now

	var result *Scenario
	var err error

	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		result = nil

		model.Actions, model.Aliases, err = s.transformActionRequestToModel(ctx, r.Actions)
		if err != nil {
			return err
		}

		res, err := s.collection.UpdateOne(ctx, bson.M{"_id": r.ID}, bson.M{"$set": model})
		if err != nil || res.MatchedCount == 0 {
			return err
		}

		err = priority.UpdateFollowing(ctx, s.collection, r.ID, r.Priority)
		if err != nil {
			return err
		}

		result, err = s.GetOneBy(ctx, r.ID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *store) Delete(ctx context.Context, id, userID string) (bool, error) {
	var deleted int64

	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		deleted = 0

		// required to get the author in action log listener.
		res, err := s.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"author": userID}})
		if err != nil || res.MatchedCount == 0 {
			return err
		}

		deleted, err = s.collection.DeleteOne(ctx, bson.M{"_id": id})
		return err
	})

	return deleted > 0, err
}

func (s *store) getSort(r FilteredQuery) bson.M {
	sortBy := cmp.Or(r.SortBy, s.defaultSortBy)
	if sortBy == "delay" {
		sortBy = "delay.value"
	}

	return common.GetSortQuery(sortBy, r.Sort)
}

func (s *store) transformEditRequestToModel(r EditRequest) libaction.Scenario {
	triggers := make([]string, 0, len(r.Triggers))
	for _, triggerRequest := range r.Triggers {
		if triggerRequest.Type == string(types.AlarmChangeEventsCount) {
			triggers = append(triggers, triggerRequest.Type+strconv.Itoa(triggerRequest.Threshold))
		} else {
			triggers = append(triggers, triggerRequest.Type)
		}
	}

	return libaction.Scenario{
		Name:                 r.Name,
		Author:               r.Author,
		Enabled:              *r.Enabled,
		DisableDuringPeriods: r.DisableDuringPeriods,
		Triggers:             triggers,
		Priority:             r.Priority,
		Delay:                r.Delay,
	}
}

func (s *store) transformActionRequestToModel(ctx context.Context, r []ActionRequest) ([]libaction.Action, []string, error) {
	var err error
	var valErr common.ValidationError

	actions := make([]libaction.Action, len(r))

	uniqueAliasesMap := make(map[string]bool)
	uniqueAliases := make([]string, 0)

	for idx, actionRequest := range r {
		transformedAlarmPatternFieldsRequest, err := s.transformer.TransformAlarmPatternFieldsRequest(ctx, actionRequest.AlarmPatternFieldsRequest)
		if err != nil {
			if errors.As(err, &valErr) {
				return nil, nil, valErr.AddFieldPrefix("actions." + strconv.Itoa(idx))
			}

			return nil, nil, err
		}

		transformEntityPatternFieldsRequest, err := s.transformer.TransformEntityPatternFieldsRequest(ctx, actionRequest.EntityPatternFieldsRequest)
		if err != nil {
			if errors.As(err, &valErr) {
				return nil, nil, valErr.AddFieldPrefix("actions." + strconv.Itoa(idx))
			}

			return nil, nil, err
		}

		for _, alias := range transformEntityPatternFieldsRequest.Aliases {
			if !uniqueAliasesMap[alias] {
				uniqueAliasesMap[alias] = true
				uniqueAliases = append(uniqueAliases, alias)
			}
		}

		actions[idx] = libaction.Action{
			Type:       r[idx].Type,
			Comment:    r[idx].Comment,
			Parameters: r[idx].Parameters,
			EntityPatternFields: transformEntityPatternFieldsRequest.ToModelWithoutFields(
				common.GetForbiddenFieldsInEntityPattern(mongo.ScenarioCollection),
			),
			AlarmPatternFields: transformedAlarmPatternFieldsRequest.ToModelWithoutFields(
				common.GetForbiddenFieldsInAlarmPattern(mongo.ScenarioCollection),
				common.GetOnlyAbsoluteTimeCondFieldsInAlarmPattern(mongo.ScenarioCollection),
			),
			DropScenarioIfNotMatched: *r[idx].DropScenarioIfNotMatched,
			EmitTrigger:              *r[idx].EmitTrigger,
		}
	}

	return actions, uniqueAliases, err
}

func getNestedObjectsPipeline() []bson.M {
	return []bson.M{
		{"$unwind": bson.M{
			"path":                       "$actions",
			"preserveNullAndEmptyArrays": true,
			"includeArrayIndex":          "action_index",
		}},
		{"$lookup": bson.M{
			"from":         mongo.PbehaviorTypeMongoCollection,
			"localField":   "actions.parameters.type",
			"foreignField": "_id",
			"as":           "actions.parameters.type",
		}},
		{"$unwind": bson.M{"path": "$actions.parameters.type", "preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from":         mongo.PbehaviorReasonMongoCollection,
			"localField":   "actions.parameters.reason",
			"foreignField": "_id",
			"as":           "actions.parameters.reason",
		}},
		{"$unwind": bson.M{"path": "$actions.parameters.reason", "preserveNullAndEmptyArrays": true}},
		{"$sort": bson.M{"action_index": 1}},
		{"$group": bson.M{
			"_id":     "$_id",
			"data":    bson.M{"$first": "$$ROOT"},
			"actions": bson.M{"$push": "$actions"},
		}},
		{"$replaceRoot": bson.M{
			"newRoot": bson.M{"$mergeObjects": bson.A{
				"$data",
				bson.M{"actions": "$actions"},
			}},
		}},
	}
}
