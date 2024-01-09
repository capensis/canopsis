package alarmtag

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmtag"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type Store interface {
	Find(ctx context.Context, r ListRequest) (*AggregationResult, error)
	GetByID(ctx context.Context, id string) (*Response, error)
	Create(ctx context.Context, r CreateRequest) (*Response, error)
	Update(ctx context.Context, r UpdateRequest) (*Response, error)
	Delete(ctx context.Context, id string) (bool, error)
}

func NewStore(dbClient mongo.DbClient, authorProvider author.Provider) Store {
	return &store{
		client:         dbClient,
		collection:     dbClient.Collection(mongo.AlarmTagCollection),
		authorProvider: authorProvider,

		defaultSearchByFields: []string{"value"},
		defaultSortBy:         "value",
	}
}

type store struct {
	client         mongo.DbClient
	collection     mongo.DbCollection
	authorProvider author.Provider

	defaultSearchByFields []string
	defaultSortBy         string
}

func (s *store) Find(ctx context.Context, r ListRequest) (*AggregationResult, error) {
	var pipeline []bson.M
	var match []bson.M
	if len(r.Values) > 0 {
		match = append(match, bson.M{"value": bson.M{"$in": r.Values}})
	}

	filter := common.GetSearchQuery(r.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		match = append(match, filter)
	}

	if len(match) > 0 {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"$and": match}})
	}

	sortBy := s.defaultSortBy
	if r.SortBy != "" {
		sortBy = r.SortBy
	}

	project := s.authorProvider.Pipeline()
	if r.WithFlags {
		project = append(project, bson.M{
			"$addFields": bson.M{
				"deletable": bson.M{"$eq": bson.A{"$type", alarmtag.TypeInternal}},
			},
		})
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

	return &res, nil
}

func (s *store) GetByID(ctx context.Context, id string) (*Response, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"_id": id}},
	}
	pipeline = append(pipeline, s.authorProvider.Pipeline()...)
	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		response := Response{}
		err = cursor.Decode(&response)
		if err != nil {
			return nil, err
		}

		return &response, nil
	}

	return nil, nil
}

func (s *store) Create(ctx context.Context, r CreateRequest) (*Response, error) {
	model := transformCreateRequestToModel(r)
	model.ID = utils.NewID()
	model.Type = alarmtag.TypeInternal
	now := datetime.NewCpsTime()
	model.Created = now
	model.Updated = now
	var response *Response
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil

		_, err := s.collection.InsertOne(ctx, model)
		if err != nil {
			if mongodriver.IsDuplicateKeyError(err) {
				return common.NewValidationError("value", "Value already exists.")
			}

			return err
		}

		response, err = s.GetByID(ctx, model.ID)
		return err
	})

	return response, err
}

func (s *store) Update(ctx context.Context, r UpdateRequest) (*Response, error) {
	now := datetime.NewCpsTime()
	var response *Response
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil
		var tag alarmtag.AlarmTag
		err := s.collection.FindOne(ctx, bson.M{"_id": r.ID}).Decode(&tag)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}
			return err
		}

		var updateResult *mongodriver.UpdateResult
		switch tag.Type {
		case alarmtag.TypeExternal:
			updateResult, err = s.collection.UpdateOne(ctx,
				bson.M{
					"_id": r.ID,
				},
				bson.M{"$set": bson.M{
					"color":   r.Color,
					"author":  r.Author,
					"updated": now,
				}},
			)
		case alarmtag.TypeInternal:
			if len(r.EntityPattern) == 0 && r.CorporateEntityPattern == "" &&
				len(r.AlarmPattern) == 0 && r.CorporateAlarmPattern == "" {
				return common.NewValidationError("alarm_pattern", "AlarmPattern or EntityPattern is required.")
			}

			tag.Color = r.Color
			tag.Author = r.Author
			tag.Updated = now
			tag.EntityPatternFields = r.EntityPatternFieldsRequest.ToModelWithoutFields(
				common.GetForbiddenFieldsInEntityPattern(mongo.AlarmTagCollection),
			)
			tag.AlarmPatternFields = r.AlarmPatternFieldsRequest.ToModelWithoutFields(
				common.GetForbiddenFieldsInAlarmPattern(mongo.AlarmTagCollection),
				common.GetOnlyAbsoluteTimeCondFieldsInAlarmPattern(mongo.AlarmTagCollection),
			)
			updateResult, err = s.collection.UpdateOne(ctx,
				bson.M{
					"_id": r.ID,
				},
				bson.M{"$set": tag},
			)
		}

		if err != nil || updateResult.MatchedCount == 0 {
			if mongodriver.IsDuplicateKeyError(err) {
				return common.NewValidationError("value", "Value already exists.")
			}

			return err
		}

		response, err = s.GetByID(ctx, r.ID)
		return err
	})

	return response, err
}

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	res, err := s.collection.DeleteOne(ctx, bson.M{"_id": id, "type": alarmtag.TypeInternal})
	return res > 0, err
}

func transformCreateRequestToModel(r CreateRequest) alarmtag.AlarmTag {
	return alarmtag.AlarmTag{
		Value:  r.Value,
		Color:  r.Color,
		Author: r.Author,
		EntityPatternFields: r.EntityPatternFieldsRequest.ToModelWithoutFields(
			common.GetForbiddenFieldsInEntityPattern(mongo.AlarmTagCollection),
		),
		AlarmPatternFields: r.AlarmPatternFieldsRequest.ToModelWithoutFields(
			common.GetForbiddenFieldsInAlarmPattern(mongo.AlarmTagCollection),
			common.GetOnlyAbsoluteTimeCondFieldsInAlarmPattern(mongo.AlarmTagCollection),
		),
	}
}
