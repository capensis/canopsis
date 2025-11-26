package template

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"slices"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/mongoquery"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/db"
	libtypes "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Store interface {
	FindData(ctx context.Context, r ListDataRequest) (AggregationDataResult, error)
	GetData(ctx context.Context, id string) (DataResponse, error)
	CreateData(ctx context.Context, r EditDataRequest) (DataResponse, error)
	UpdateData(ctx context.Context, r EditDataRequest) (DataResponse, error)
	DeleteData(ctx context.Context, id, author string) (bool, error)
	FindTest(ctx context.Context, r ListTestRequest, userID string) (AggregationTestResult, error)
	GetTest(ctx context.Context, id, userID string) (TestResponse, error)
	CreateTest(ctx context.Context, r EditTestRequest) (TestResponse, error)
	UpdateTest(ctx context.Context, r EditTestRequest) (TestResponse, error)
	DeleteTest(ctx context.Context, id, author string) (bool, error)
}

func NewStore(
	client mongo.DbClient,
	authorProvider author.Provider,
	enforcer security.Enforcer,
	testTypePermMapping map[int][]any,
	decoder encoding.Decoder,
) Store {
	return &store{
		client:              client,
		testDataCollection:  client.Collection(mongo.TemplateTestDataCollection),
		testCollection:      client.Collection(mongo.TemplateTestCollection),
		alarmCollection:     client.Collection(mongo.AlarmMongoCollection),
		entityCollection:    client.Collection(mongo.EntityMongoCollection),
		userCollection:      client.Collection(mongo.UserCollection),
		authorProvider:      authorProvider,
		enforcer:            enforcer,
		testTypePermMapping: testTypePermMapping,
		decoder:             decoder,
		defaultSortBy:       "name",
		defaultSearchByFields: []string{
			"name",
			"description",
		},
		collectionNamesByType: map[int]string{
			TypeTestEventFilterRule:   mongo.EventFilterRuleCollection,
			TypeTestLinkRule:          mongo.LinkRuleMongoCollection,
			TypeTestActionScenario:    mongo.ScenarioCollection,
			TypeTestWidget:            mongo.WidgetMongoCollection,
			TypeTestDeclareTicketRule: mongo.DeclareTicketRuleCollection,
			TypeTestDynamicInfosRule:  mongo.DynamicInfosRulesMongoCollection,
			TypeTestInstruction:       mongo.InstructionMongoCollection,
			TypeTestJob:               mongo.JobMongoCollection,
			TypeTestMetaAlarmRule:     mongo.MetaAlarmRulesMongoCollection,
			TypeTestWebhookTokenRule:  mongo.WebhookTokenRuleCollection,
		},
		dupErrorParser: validation.NewDuplicateErrorParser(),
	}
}

type store struct {
	client                mongo.DbClient
	testDataCollection    mongo.DbCollection
	testCollection        mongo.DbCollection
	alarmCollection       mongo.DbCollection
	entityCollection      mongo.DbCollection
	userCollection        mongo.DbCollection
	authorProvider        author.Provider
	enforcer              security.Enforcer
	decoder               encoding.Decoder
	testTypePermMapping   map[int][]any
	defaultSortBy         string
	defaultSearchByFields []string
	collectionNamesByType map[int]string
	dupErrorParser        validation.DuplicateErrorParser
}

func (s *store) FindData(ctx context.Context, r ListDataRequest) (AggregationDataResult, error) {
	var res AggregationDataResult
	beforeLimit := make([]bson.M, 0)
	if r.Type != nil {
		beforeLimit = append(beforeLimit, bson.M{"$match": bson.M{"type": r.Type}})
	}

	if r.EventPattern != "" {
		var eventPattern pattern.Event
		err := s.decoder.Decode([]byte(r.EventPattern), &eventPattern)
		if err != nil {
			return res, validation.NewError(
				validator.ValidationErrors{validation.NewFieldError("event_pattern", "EventPattern", "EventPattern")},
				r,
			)
		}

		q, err := db.EventPatternToMongoQuery(eventPattern, "body")
		if err != nil {
			return res, validation.NewError(
				validator.ValidationErrors{validation.NewFieldError("event_pattern", "EventPattern", "EventPattern")},
				r,
			)
		}

		beforeLimit = append(beforeLimit, bson.M{"$match": bson.M{"$and": []bson.M{
			{"type": TypeTestDataEvent},
			q,
		}}})
	}

	filter := mongoquery.GetSearchQuery(r.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		beforeLimit = append(beforeLimit, bson.M{"$match": filter})
	}

	cursor, err := s.testDataCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		beforeLimit,
		mongoquery.GetSortQuery(cmp.Or(r.SortBy, s.defaultSortBy), r.Sort),
	))
	if err != nil {
		return res, err
	}

	defer cursor.Close(ctx)
	if cursor.Next(ctx) {
		err = cursor.Decode(&res)
		if err != nil {
			return res, err
		}
	}

	if err = cursor.Err(); err != nil {
		return res, err
	}

	return res, nil
}

func (s *store) GetData(ctx context.Context, id string) (DataResponse, error) {
	res := DataResponse{}
	err := s.testDataCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&res)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return res, nil
		}

		return res, err
	}

	return res, nil
}

func (s *store) CreateData(ctx context.Context, r EditDataRequest) (DataResponse, error) {
	now := datetime.NewCpsTime()
	model := DataModel{
		ID:          utils.NewID(),
		Type:        *r.Type,
		Name:        r.Name,
		Description: r.Description,
		Body:        r.Body,
		Headers:     r.Headers,
		Author:      r.Author,
		Created:     &now,
		Updated:     &now,
	}
	var res DataResponse
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		res = DataResponse{}

		_, err := s.testDataCollection.InsertOne(ctx, model)
		if err != nil {
			if mongodriver.IsDuplicateKeyError(err) {
				return s.dupErrorParser.Parse(err, res)
			}

			return err
		}

		res, err = s.GetData(ctx, model.ID)

		return err
	})

	return res, err
}

func (s *store) UpdateData(ctx context.Context, r EditDataRequest) (DataResponse, error) {
	now := datetime.NewCpsTime()
	model := DataModel{
		Type:        *r.Type,
		Name:        r.Name,
		Description: r.Description,
		Body:        r.Body,
		Headers:     r.Headers,
		Author:      r.Author,
		Updated:     &now,
	}
	var res DataResponse
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		res = DataResponse{}

		prev, err := s.GetData(ctx, r.ID)
		if err != nil || prev.ID == "" {
			return err
		}

		if prev.Type != model.Type {
			return validation.NewError(
				validator.ValidationErrors{validation.NewFieldError("unchangeable", "Type", "Type")},
				r,
			)
		}

		_, err = s.testDataCollection.UpdateOne(ctx, bson.M{"_id": r.ID}, bson.M{"$set": model})
		if err != nil {
			if mongodriver.IsDuplicateKeyError(err) {
				return s.dupErrorParser.Parse(err, res)
			}

			return err
		}

		res, err = s.GetData(ctx, r.ID)

		return err
	})

	return res, err
}

func (s *store) DeleteData(ctx context.Context, id, author string) (bool, error) {
	var res bool
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		res = false

		isUsed, err := s.testDataIsUsed(ctx, id)
		if err != nil {
			return err
		}

		if isUsed {
			return httperror.NewConflictError("The test data cannot be deleted because it is referenced by a test.")
		}

		// required to get the author in action log listener
		ur, err := s.testDataCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"author": author}})
		if err != nil || ur.MatchedCount == 0 {
			return err
		}

		d, err := s.testDataCollection.DeleteOne(ctx, bson.M{"_id": id})
		if err != nil || d == 0 {
			return err
		}

		res = true

		return nil
	})

	return res, err
}

func (s *store) FindTest(ctx context.Context, r ListTestRequest, userID string) (AggregationTestResult, error) {
	var res AggregationTestResult
	beforeLimit := make([]bson.M, 0)
	types, err := s.getAuthorizedTestTypes(userID)
	if err != nil {
		return res, err
	}

	if r.Type == nil {
		beforeLimit = append(beforeLimit, bson.M{"$match": bson.M{"type": bson.M{"$in": types}}})
	} else {
		if !slices.Contains(types, *r.Type) {
			return res, validation.NewError(
				validator.ValidationErrors{validation.NewFieldError("not_accessible", "Type", "Type")},
				r,
			)
		}

		beforeLimit = append(beforeLimit, bson.M{"$match": bson.M{"type": r.Type}})
	}

	if len(r.IDs) > 0 {
		beforeLimit = append(beforeLimit, bson.M{"$match": bson.M{"_id": bson.M{"$in": r.IDs}}})
	}

	if r.Rule != "" {
		beforeLimit = append(beforeLimit, bson.M{"$match": bson.M{"rule._id": r.Rule}})
	}

	filter := mongoquery.GetSearchQuery(r.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		beforeLimit = append(beforeLimit, bson.M{"$match": filter})
	}

	sort := mongoquery.GetSortQuery(cmp.Or(r.SortBy, s.defaultSortBy), r.Sort)
	afterLimit := s.getTestNestedObjectsPipeline()
	afterLimit = append(afterLimit, sort)
	cursor, err := s.testCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		beforeLimit,
		sort,
		afterLimit,
	))
	if err != nil {
		return res, err
	}

	defer cursor.Close(ctx)
	if cursor.Next(ctx) {
		err = cursor.Decode(&res)
		if err != nil {
			return res, err
		}
	}

	if err = cursor.Err(); err != nil {
		return res, err
	}

	return res, nil
}

func (s *store) GetTest(ctx context.Context, id, userID string) (TestResponse, error) {
	res := TestResponse{}
	types, err := s.getAuthorizedTestTypes(userID)
	if err != nil {
		return res, err
	}

	pipeline := []bson.M{
		{"$match": bson.M{
			"_id":  id,
			"type": bson.M{"$in": types},
		}},
	}
	pipeline = append(pipeline, s.getTestNestedObjectsPipeline()...)
	cursor, err := s.testCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return res, err
	}

	defer cursor.Close(ctx)
	if !cursor.Next(ctx) {
		return res, nil
	}

	err = cursor.Decode(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *store) CreateTest(ctx context.Context, r EditTestRequest) (TestResponse, error) {
	types, err := s.getAuthorizedTestTypes(r.Author)
	if err != nil {
		return TestResponse{}, err
	}

	if !slices.Contains(types, *r.Type) {
		return TestResponse{}, validation.NewError(
			validator.ValidationErrors{validation.NewFieldError("not_accessible", "Type", "Type")},
			r,
		)
	}

	now := datetime.NewCpsTime()
	model := TestModel{
		ID:          utils.NewID(),
		Name:        r.Name,
		Description: r.Description,
		Type:        *r.Type,
		Author:      r.Author,
		Created:     &now,
		Updated:     &now,
	}
	model.Data.Event = r.Data.Event
	model.Data.Response = r.Data.Response
	model.Data.Responses = r.Data.Responses
	model.Data.User = r.Data.User
	ruleCollectionName, ok := s.collectionNamesByType[*r.Type]
	if !ok {
		return TestResponse{}, fmt.Errorf("uknown test type %d", *r.Type)
	}

	ruleCollection := s.client.Collection(ruleCollectionName)
	var res TestResponse
	err = s.client.WithTransaction(ctx, func(ctx context.Context) error {
		res = TestResponse{}

		rule := struct {
			Name string `bson:"name"`
		}{}
		err = ruleCollection.
			FindOne(ctx, bson.M{"_id": r.Rule}, options.FindOne().SetProjection(bson.M{
				"name": bson.M{"$ifNull": bson.A{
					"$name",
					"$description",
					"$title",
				}},
			})).
			Decode(&rule)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return validation.NewError(
					validator.ValidationErrors{validation.NewFieldError("not_exist", "Rule", "Rule")},
					r,
				)
			}

			return err
		}

		model.Rule.ID = r.Rule
		model.Rule.Name = rule.Name
		model.Data.Alarm, model.Data.Entity, err = s.validateTestData(ctx, r, nil)
		if err != nil {
			return err
		}

		_, err = s.testCollection.InsertOne(ctx, model)
		if err != nil {
			if mongodriver.IsDuplicateKeyError(err) {
				return s.dupErrorParser.Parse(err, res)
			}

			return err
		}

		res, err = s.GetTest(ctx, model.ID, r.Author)

		return err
	})

	return res, err
}

func (s *store) UpdateTest(ctx context.Context, r EditTestRequest) (TestResponse, error) {
	types, err := s.getAuthorizedTestTypes(r.Author)
	if err != nil {
		return TestResponse{}, err
	}

	if !slices.Contains(types, *r.Type) {
		return TestResponse{}, validation.NewError(
			validator.ValidationErrors{validation.NewFieldError("not_accessible", "Type", "Type")},
			r,
		)
	}

	now := datetime.NewCpsTime()
	model := TestModel{
		Name:        r.Name,
		Description: r.Description,
		Type:        *r.Type,
		Author:      r.Author,
		Updated:     &now,
	}
	model.Data.Event = r.Data.Event
	model.Data.Response = r.Data.Response
	model.Data.Responses = r.Data.Responses
	model.Data.User = r.Data.User
	var res TestResponse
	err = s.client.WithTransaction(ctx, func(ctx context.Context) error {
		res = TestResponse{}

		var prev TestModel
		err = s.testCollection.FindOne(ctx, bson.M{"_id": r.ID, "type": bson.M{"$in": types}}).Decode(&prev)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		if prev.Type != model.Type {
			return validation.NewError(
				validator.ValidationErrors{validation.NewFieldError("unchangeable", "Type", "Type")},
				r,
			)
		}

		if prev.Rule.ID != r.Rule {
			return validation.NewError(
				validator.ValidationErrors{validation.NewFieldError("unchangeable", "Rule", "Rule")},
				r,
			)
		}

		model.Rule = prev.Rule
		model.Data.Alarm, model.Data.Entity, err = s.validateTestData(ctx, r, &prev)
		if err != nil {
			return err
		}

		updateRes, err := s.testCollection.UpdateOne(ctx,
			bson.M{"_id": r.ID, "type": bson.M{"$in": types}},
			bson.M{"$set": model},
		)
		if err != nil || updateRes.MatchedCount == 0 {
			if mongodriver.IsDuplicateKeyError(err) {
				return s.dupErrorParser.Parse(err, res)
			}

			return err
		}

		res, err = s.GetTest(ctx, r.ID, r.Author)

		return err
	})

	return res, err
}

func (s *store) DeleteTest(ctx context.Context, id, author string) (bool, error) {
	types, err := s.getAuthorizedTestTypes(author)
	if err != nil {
		return false, err
	}

	filter := bson.M{"_id": id, "type": bson.M{"$in": types}}
	var res bool
	err = s.client.WithTransaction(ctx, func(ctx context.Context) error {
		res = false
		// required to get the author in action log listener
		ur, err := s.testCollection.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"author": author}})
		if err != nil || ur.MatchedCount == 0 {
			return err
		}

		d, err := s.testCollection.DeleteOne(ctx, filter)
		if err != nil || d == 0 {
			return err
		}

		res = true

		return nil
	})

	return res, err
}

func (s *store) getAuthorizedTestTypes(userID string) ([]int, error) {
	types := make([]int, 0, len(s.testTypePermMapping))
	for t, perm := range s.testTypePermMapping {
		rvals := make([]any, 0, len(perm)+1)
		rvals = append(rvals, userID)
		rvals = append(rvals, perm...)
		ok, err := s.enforcer.Enforce(rvals...)
		if err != nil {
			return nil, err
		}

		if ok {
			types = append(types, t)
		}
	}

	return types, nil
}

func (s *store) validateTestData(ctx context.Context, r EditTestRequest, prevTest *TestModel) (*libtypes.AlarmWithEntityField, *libtypes.Entity, error) {
	var alarm *libtypes.AlarmWithEntityField
	var entity *libtypes.Entity
	if r.Data.Event != "" {
		err := s.testDataCollection.FindOne(ctx, bson.M{"_id": r.Data.Event, "type": TypeTestDataEvent},
			options.FindOne().SetProjection(bson.M{"_id": 1})).Err()
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil, nil, validation.NewError(
					validator.ValidationErrors{validation.NewFieldError("not_exist", "Event", "Data.Event")},
					r,
				)
			}

			return nil, nil, err
		}
	}

	if r.Data.Response != "" {
		err := s.testDataCollection.FindOne(ctx, bson.M{"_id": r.Data.Response, "type": TypeTestDataResponse},
			options.FindOne().SetProjection(bson.M{"_id": 1})).Err()
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil, nil, validation.NewError(
					validator.ValidationErrors{validation.NewFieldError("not_exist", "Response", "Data.Response")},
					r,
				)
			}

			return nil, nil, err
		}
	}

	if len(r.Data.Responses) > 0 {
		ids := make([]string, 0, len(r.Data.Responses))
		for _, id := range r.Data.Responses {
			ids = append(ids, id)
		}

		ids = utils.Unique(ids)
		cursor, err := s.testDataCollection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}, "type": TypeTestDataResponse},
			options.Find().SetProjection(bson.M{"_id": 1}))
		if err != nil {
			return nil, nil, err
		}

		var responses []DataModel
		err = cursor.All(ctx, &responses)
		if err != nil {
			return nil, nil, err
		}

		if len(ids) != len(responses) {
			return nil, nil, validation.NewError(
				validator.ValidationErrors{validation.NewFieldError("not_exist", "Responses", "Data.Responses")},
				r,
			)
		}
	}

	if r.Data.User != "" {
		err := s.userCollection.FindOne(ctx, bson.M{"_id": r.Data.User},
			options.FindOne().SetProjection(bson.M{"_id": 1})).Err()
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil, nil, validation.NewError(
					validator.ValidationErrors{validation.NewFieldError("not_exist", "User", "Data.User")},
					r,
				)
			}

			return nil, nil, err
		}
	}

	if r.Data.Alarm != "" {
		if prevTest != nil && prevTest.Data.Alarm != nil && r.Data.Alarm == prevTest.Data.Alarm.ID {
			alarm = prevTest.Data.Alarm
		} else {
			alarm = &libtypes.AlarmWithEntityField{}
			err := s.alarmCollection.FindOne(ctx, bson.M{"_id": r.Data.Alarm},
				options.FindOne().SetProjection(bson.M{"v.steps": 0})).Decode(alarm)
			if err != nil {
				if errors.Is(err, mongodriver.ErrNoDocuments) {
					return nil, nil, validation.NewError(
						validator.ValidationErrors{validation.NewFieldError("not_exist", "Alarm", "Data.Alarm")},
						r,
					)
				}

				return nil, nil, err
			}
		}
	}

	if r.Data.Entity != "" {
		if prevTest != nil && prevTest.Data.Entity != nil && r.Data.Entity == prevTest.Data.Entity.ID {
			entity = prevTest.Data.Entity
		} else {
			entity = &libtypes.Entity{}
			err := s.entityCollection.FindOne(ctx, bson.M{"_id": r.Data.Entity}).Decode(entity)
			if err != nil {
				if errors.Is(err, mongodriver.ErrNoDocuments) {
					return nil, nil, validation.NewError(
						validator.ValidationErrors{validation.NewFieldError("not_exist", "Entity", "Data.Entity")},
						r,
					)
				}

				return nil, nil, err
			}
		}
	}

	return alarm, entity, nil
}

func (s *store) testDataIsUsed(ctx context.Context, id string) (bool, error) {
	cursor, err := s.testCollection.Aggregate(ctx, []bson.M{
		{"$addFields": bson.M{
			"responses": bson.M{"$objectToArray": "$data.responses"},
		}},
		{"$unwind": bson.M{"path": "$responses", "preserveNullAndEmptyArrays": true}},
		{"$match": bson.M{"$or": []bson.M{
			{"data.event": id},
			{"data.response": id},
			{"responses.v": id},
		}}},
		{"$limit": 1},
		{"$project": bson.M{"_id": 1}},
	})
	if err != nil {
		return false, err
	}

	defer cursor.Close(ctx)
	isUsed := cursor.Next(ctx)
	if err = cursor.Err(); err != nil {
		return false, err
	}

	return isUsed, nil
}

func (s *store) getTestNestedObjectsPipeline() []bson.M {
	pipeline := s.authorProvider.PipelineForField("data.user")

	return append(pipeline, []bson.M{
		{"$addFields": bson.M{
			"data.alarm": bson.M{"$cond": bson.M{
				"if": "$data.alarm",
				"then": bson.M{
					"_id":          "$data.alarm._id",
					"display_name": "$data.alarm.v.display_name",
				},
				"else": "$data.alarm",
			}},
		}},
		{"$lookup": bson.M{
			"from":         mongo.TemplateTestDataCollection,
			"localField":   "data.event",
			"foreignField": "_id",
			"as":           "data.event",
			"pipeline": []bson.M{
				{"$match": bson.M{"type": TypeTestDataEvent}},
			},
		}},
		{"$unwind": bson.M{"path": "$data.event", "preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from":         mongo.TemplateTestDataCollection,
			"localField":   "data.response",
			"foreignField": "_id",
			"as":           "data.response",
			"pipeline": []bson.M{
				{"$match": bson.M{"type": TypeTestDataResponse}},
			},
		}},
		{"$unwind": bson.M{"path": "$data.response", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			"doc":       "$$ROOT",
			"responses": bson.M{"$objectToArray": "$data.responses"},
		}},
		{"$unwind": bson.M{"path": "$responses", "preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from":         mongo.TemplateTestDataCollection,
			"localField":   "responses.v",
			"foreignField": "_id",
			"as":           "responses.v",
		}},
		{"$unwind": bson.M{"path": "$responses.v", "preserveNullAndEmptyArrays": true}},
		{"$group": bson.M{
			"_id": "$_id",
			"doc": bson.M{"$first": "$doc"},
			"responses": bson.M{"$push": bson.M{"$cond": bson.M{
				"if":   "$responses.k",
				"then": "$responses",
				"else": "$$REMOVE",
			}}},
		}},
		{"$addFields": bson.M{
			"doc.data.responses": bson.M{"$arrayToObject": "$responses"},
		}},
		{"$replaceRoot": bson.M{"newRoot": "$doc"}},
	}...)
}
