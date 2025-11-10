package alarmtag

import (
	"cmp"
	"context"
	"errors"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/mongoquery"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/patternfields"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmtag"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Store interface {
	Find(ctx context.Context, r ListRequest) (*AggregationResult, error)
	FindLabels(ctx context.Context, r ListLabelsRequest) (*AggregationLabelResult, error)
	GetByID(ctx context.Context, id string) (*Response, error)
	Create(ctx context.Context, r CreateRequest) (*Response, error)
	Update(ctx context.Context, r UpdateRequest) (*Response, error)
	Delete(ctx context.Context, id, userID string) (bool, error)
}

func NewStore(dbClient mongo.DbClient, authorProvider author.Provider, transformer patternfields.Transformer) Store {
	return &store{
		client:          dbClient,
		collection:      dbClient.Collection(mongo.AlarmTagCollection),
		labelCollection: dbClient.Collection(mongo.AlarmTagColorCollection),
		authorProvider:  authorProvider,
		transformer:     transformer,

		defaultSearchByFields: []string{"value"},
		defaultSortBy:         "value",

		dupErrorParser: validation.NewDuplicateErrorParser(map[string]string{
			"value": "Value already exists.",
		}),
	}
}

type store struct {
	client          mongo.DbClient
	collection      mongo.DbCollection
	labelCollection mongo.DbCollection
	authorProvider  author.Provider
	transformer     patternfields.Transformer

	defaultSearchByFields []string
	defaultSortBy         string

	dupErrorParser validation.DuplicateErrorParser
}

func (s *store) Find(ctx context.Context, r ListRequest) (*AggregationResult, error) {
	var pipeline []bson.M
	var match []bson.M
	if len(r.Values) > 0 {
		match = append(match, bson.M{"value": bson.M{"$in": r.Values}})
	}

	filter := mongoquery.GetSearchQuery(r.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		match = append(match, filter)
	}

	if len(match) > 0 {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"$and": match}})
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
		mongoquery.GetSortQuery(cmp.Or(r.SortBy, s.defaultSortBy), r.Sort),
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

func (s *store) FindLabels(ctx context.Context, r ListLabelsRequest) (*AggregationLabelResult, error) {
	var pipeline []bson.M
	var match []bson.M
	if len(r.IDs) > 0 {
		match = append(match, bson.M{"_id": bson.M{"$in": r.IDs}})
	}

	filter := mongoquery.GetSearchQuery(r.Search, []string{"_id"})
	if len(filter) > 0 {
		match = append(match, filter)
	}

	if len(match) > 0 {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"$and": match}})
	}

	cursor, err := s.labelCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		pipeline,
		mongoquery.GetSortQuery("_id", common.SortAsc),
	))
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	res := AggregationLabelResult{}
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
	var response *Response

	now := datetime.NewCpsTime()
	tag := alarmtag.AlarmTag{
		ID:      utils.NewID(),
		Type:    alarmtag.TypeInternal,
		Value:   r.Value,
		Color:   r.Color,
		Author:  r.Author,
		Created: now,
		Updated: now,
	}

	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil

		err := s.transformPatternRequestsToModel(ctx, r.EntityRequest, r.AlarmRequest, &tag)
		if err != nil {
			return err
		}

		_, err = s.collection.InsertOne(ctx, tag)
		if err != nil {
			if mongodriver.IsDuplicateKeyError(err) {
				return s.dupErrorParser.Parse(err)
			}

			return err
		}

		response, err = s.GetByID(ctx, tag.ID)
		if err != nil {
			return err
		}

		label := tag.Value
		splitLabel := strings.Split(label, ":")
		if len(splitLabel) > 0 {
			label = splitLabel[0]
		}

		_, err = s.labelCollection.UpdateOne(
			ctx,
			bson.M{"_id": label},
			bson.M{"$setOnInsert": bson.M{"color": tag.Color}},
			options.UpdateOne().SetUpsert(true),
		)

		return err
	})

	return response, err
}

func (s *store) Update(ctx context.Context, r UpdateRequest) (*Response, error) {
	var response *Response

	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil

		now := datetime.NewCpsTime()

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

			err = s.transformPatternRequestsToModel(ctx, r.EntityRequest, r.AlarmRequest, &tag)
			if err != nil {
				return err
			}

			tag.Color = r.Color
			tag.Author = r.Author
			tag.Updated = now

			updateResult, err = s.collection.UpdateOne(ctx, bson.M{"_id": r.ID}, bson.M{"$set": tag})
		}

		if err != nil || updateResult.MatchedCount == 0 {
			if mongodriver.IsDuplicateKeyError(err) {
				return s.dupErrorParser.Parse(err)
			}

			return err
		}

		response, err = s.GetByID(ctx, r.ID)
		return err
	})

	return response, err
}

func (s *store) Delete(ctx context.Context, id, userID string) (bool, error) {
	var deleted int64

	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		deleted = 0

		// required to get the author in action log listener.
		result, err := s.collection.UpdateOne(ctx, bson.M{"_id": id, "type": alarmtag.TypeInternal}, bson.M{"$set": bson.M{"author": userID}})
		if err != nil || result.MatchedCount == 0 {
			return err
		}

		deleted, err = s.collection.DeleteOne(ctx, bson.M{"_id": id, "type": alarmtag.TypeInternal})
		return err
	})

	return deleted > 0, err
}

func (s *store) transformPatternRequestsToModel(
	ctx context.Context,
	entityPatternReq patternfields.EntityRequest,
	alarmPatternReq patternfields.AlarmRequest,
	model *alarmtag.AlarmTag,
) error {
	transformedEntityPatternRequest, err := s.transformer.TransformEntityRequest(ctx, entityPatternReq)
	if err != nil {
		return err
	}

	transformedAlarmPatternRequest, err := s.transformer.TransformAlarmRequest(ctx, alarmPatternReq)
	if err != nil {
		return err
	}

	model.Aliases = transformedEntityPatternRequest.Aliases
	model.EntityPatternFields = transformedEntityPatternRequest.ToModelWithoutFields(
		patternfields.GetForbiddenFieldsInEntityPattern(mongo.AlarmTagCollection),
	)
	model.AlarmPatternFields = transformedAlarmPatternRequest.ToModelWithoutFields(
		patternfields.GetForbiddenFieldsInAlarmPattern(mongo.AlarmTagCollection),
		patternfields.GetOnlyAbsoluteTimeCondFieldsInAlarmPattern(mongo.AlarmTagCollection),
	)

	return nil
}
