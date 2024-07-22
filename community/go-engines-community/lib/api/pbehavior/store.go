package pbehavior

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	libentity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/db"
	libmatch "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	libtypes "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/timespan"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/kylelemons/godebug/pretty"
	librrule "github.com/teambition/rrule-go"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	nextEventMaxMonths = 1
)

type Store interface {
	Insert(ctx context.Context, r CreateRequest) (*Response, error)
	Find(ctx context.Context, r ListRequest) (*AggregationResult, error)
	FindByEntityID(ctx context.Context, entity libtypes.Entity, r FindByEntityIDRequest) ([]Response, error)
	CalendarByEntityID(ctx context.Context, entity libtypes.Entity, r CalendarByEntityIDRequest) ([]CalendarResponse, error)
	GetOneBy(ctx context.Context, id string) (*Response, error)
	FindEntities(ctx context.Context, pbhID string, request EntitiesListRequest) (*AggregationEntitiesResult, error)
	Update(ctx context.Context, r UpdateRequest) (*Response, error)
	UpdateByPatch(ctx context.Context, r PatchRequest) (*Response, error)
	Delete(ctx context.Context, id, userID string) (bool, error)
	DeleteByName(ctx context.Context, name, userID string) (string, error)
	FindEntity(ctx context.Context, entityId string) (*libtypes.Entity, error)
	EntityInsert(ctx context.Context, r BulkEntityCreateRequestItem) (*Response, error)
	EntityDelete(ctx context.Context, r BulkEntityDeleteRequestItem) (string, error)
	ConnectorCreate(ctx context.Context, r BulkConnectorCreateRequestItem) (*Response, error)
	ConnectorDelete(ctx context.Context, r BulkConnectorDeleteRequestItem) (string, error)
}

type store struct {
	dbClient mongo.DbClient

	dbCollection       mongo.DbCollection
	entityDbCollection mongo.DbCollection

	authorProvider         author.Provider
	entityTypeResolver     pbehavior.EntityTypeResolver
	pbhTypeComputer        pbehavior.TypeComputer
	timezoneConfigProvider config.TimezoneConfigProvider
	defaultSortBy          string

	entitiesDefaultSearchByFields []string
	entitiesDefaultSortBy         string

	dupErrorRegexp *regexp.Regexp
}

func NewStore(
	dbClient mongo.DbClient,
	entityTypeResolver pbehavior.EntityTypeResolver,
	pbhTypeComputer pbehavior.TypeComputer,
	timezoneConfigProvider config.TimezoneConfigProvider,
	authorProvider author.Provider,
) Store {
	return &store{
		dbClient:                      dbClient,
		dbCollection:                  dbClient.Collection(mongo.PbehaviorMongoCollection),
		entityDbCollection:            dbClient.Collection(mongo.EntityMongoCollection),
		entityTypeResolver:            entityTypeResolver,
		pbhTypeComputer:               pbhTypeComputer,
		timezoneConfigProvider:        timezoneConfigProvider,
		authorProvider:                authorProvider,
		defaultSortBy:                 "created",
		entitiesDefaultSearchByFields: []string{"_id", "name", "type"},
		entitiesDefaultSortBy:         "_id",
		dupErrorRegexp:                regexp.MustCompile(`{ ([^:]+)`),
	}
}

func (s *store) Insert(ctx context.Context, r CreateRequest) (*Response, error) {
	now := datetime.NewCpsTime()
	doc := s.transformRequestToDocument(r.EditRequest)
	doc.ID = cmp.Or(r.ID, utils.NewID())

	rruleEnd, err := pbehavior.GetRruleEnd(*r.Start, r.RRule, s.timezoneConfigProvider.Get().Location)
	if err != nil {
		return nil, err
	}

	doc.Created = &now
	doc.Updated = &now
	doc.Comments = make([]*pbehavior.Comment, 0)
	doc.RRuleEnd = rruleEnd

	var pbh *Response
	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		pbh = nil

		_, err := s.dbCollection.InsertOne(ctx, doc)
		if err != nil {
			if mongodriver.IsDuplicateKeyError(err) {
				return s.parseDupError(err)
			}

			return err
		}

		pbh, err = s.GetOneBy(ctx, doc.ID)
		return err
	})

	return pbh, err
}

func (s *store) Find(ctx context.Context, r ListRequest) (*AggregationResult, error) {
	mongoQuery := CreateMongoQuery(s.dbClient, s.authorProvider)
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

	err = s.transformResponse(ctx, result.Data)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *store) FindByEntityID(ctx context.Context, entity libtypes.Entity, r FindByEntityIDRequest) ([]Response, error) {
	pbhIDs, err := s.getMatchedPbhIDs(ctx, entity)
	if err != nil {
		return nil, err
	}

	pipeline := []bson.M{{"$match": bson.M{"_id": bson.M{"$in": pbhIDs}}}}
	pipeline = append(pipeline, GetNestedObjectsPipeline(s.authorProvider)...)
	pipeline = append(pipeline, common.GetSortQuery("created", common.SortAsc))
	if r.WithFlags {
		pipeline = append(pipeline, bson.M{"$addFields": bson.M{
			"editable": bson.M{"$cond": bson.M{
				"if":   "$origin",
				"then": false,
				"else": true,
			}},
		}})
	}
	cursor, err := s.dbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	res := make([]Response, 0)
	err = cursor.All(ctx, &res)
	if err != nil {
		return nil, err
	}

	err = s.transformResponse(ctx, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *store) CalendarByEntityID(ctx context.Context, entity libtypes.Entity, r CalendarByEntityIDRequest) ([]CalendarResponse, error) {
	pbhIDs, err := s.getMatchedPbhIDs(ctx, entity)
	if err != nil {
		return nil, err
	}

	location := s.timezoneConfigProvider.Get().Location
	span := timespan.New(r.From.In(location).Time, r.To.In(location).Time)
	computed, err := s.pbhTypeComputer.ComputeByIds(ctx, span, pbhIDs)
	if err != nil {
		return nil, err
	}

	res := make([]CalendarResponse, 0, len(computed.ComputedPbehaviors))
	for pbhId, computedPbehavior := range computed.ComputedPbehaviors {
		for _, computedType := range computedPbehavior.Types {
			res = append(res, CalendarResponse{
				ID:    pbhId,
				Title: computedPbehavior.Name,
				Color: computedPbehavior.Color,
				From:  datetime.CpsTime{Time: computedType.Span.From()},
				To:    datetime.CpsTime{Time: computedType.Span.To()},
				Type:  computed.TypesByID[computedType.ID],
			})
		}
	}

	sort.Slice(res, sortCalendarResponse(res))

	return res, nil
}

func (s *store) GetOneBy(ctx context.Context, id string) (*Response, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"_id": id}},
	}
	pipeline = append(pipeline, GetNestedObjectsPipeline(s.authorProvider)...)
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
	pbh, err := s.GetOneBy(ctx, pbhID)
	if err != nil || pbh == nil {
		return nil, err
	}

	match, err := db.EntityPatternToMongoQuery(pbh.EntityPattern, "")
	if err != nil {
		return nil, err
	}

	if len(match) == 0 {
		return &AggregationEntitiesResult{
			Data: []libentity.Entity{},
		}, nil
	}

	pipeline := []bson.M{
		{"$match": match},
	}
	filter := common.GetSearchQuery(request.Search, s.entitiesDefaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
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
	cursor, err := s.entityDbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		request.Query,
		pipeline,
		common.GetSortQuery(cmp.Or(request.SortBy, s.entitiesDefaultSortBy), request.Sort),
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

func (s *store) Update(ctx context.Context, r UpdateRequest) (*Response, error) {
	now := datetime.NewCpsTime()
	doc := s.transformRequestToDocument(r.EditRequest)
	doc.Updated = &now

	unset := bson.M{
		"rrule_cstart": "",
	}

	if r.Stop == nil {
		unset["tstop"] = ""
	}

	rruleEnd, err := pbehavior.GetRruleEnd(*r.Start, r.RRule, s.timezoneConfigProvider.Get().Location)
	if err != nil {
		return nil, err
	}
	if rruleEnd == nil {
		unset["rrule_end"] = ""
	} else {
		doc.RRuleEnd = rruleEnd
	}

	update := bson.M{"$set": doc}
	if len(unset) > 0 {
		update["$unset"] = unset
	}

	var pbh *Response
	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		pbh = nil

		prevPbh := pbehavior.PBehavior{}
		err := s.dbCollection.FindOne(ctx, bson.M{"_id": r.ID}).Decode(&prevPbh)
		if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
			return err
		}

		if prevPbh.Origin != "" {
			if len(prevPbh.Entities) > 0 {
				return common.NewValidationError("_id", "Cannot update an external pbehavior.")
			}

			valErr := common.NewValidationError("_id", "Cannot update a pbehavior with origin.")
			if !*r.Enabled || r.RRule != "" || len(r.Exdates) > 0 || len(r.Exceptions) > 0 || r.CorporateEntityPattern != "" {
				return valErr
			}

			if diff := pretty.Compare(prevPbh.EntityPattern, r.EntityPattern); diff != "" {
				return valErr
			}
		}

		_, err = s.dbCollection.UpdateOne(ctx, bson.M{"_id": r.ID}, update)
		if err != nil {
			if mongodriver.IsDuplicateKeyError(err) {
				return s.parseDupError(err)
			}

			return err
		}

		pbh, err = s.GetOneBy(ctx, r.ID)

		return err
	})

	return pbh, err
}

func (s *store) UpdateByPatch(ctx context.Context, r PatchRequest) (*Response, error) {
	set := bson.M{
		"author":  r.Author,
		"updated": datetime.NewCpsTime(),
	}
	unset := bson.M{}
	rruleUpdated := false
	if r.Name != nil {
		set["name"] = *r.Name
	}
	if r.Enabled != nil {
		set["enabled"] = *r.Enabled
	}
	if r.Reason != nil {
		set["reason"] = *r.Reason
	}
	if r.Type != nil {
		set["type_"] = *r.Type
	}
	if r.RRule != nil {
		set["rrule"] = *r.RRule
		rruleUpdated = true
	}
	if r.Start != nil {
		set["tstart"] = *r.Start
		rruleUpdated = true
	}
	if r.Stop.isSet {
		if r.Stop.val == nil {
			unset["tstop"] = ""
		} else {
			set["tstop"] = *r.Stop.val
		}
	}
	if r.Exdates != nil {
		exdates := make([]pbehavior.Exdate, len(r.Exdates))
		for i := range r.Exdates {
			exdates[i].Type = r.Exdates[i].Type
			exdates[i].Begin = r.Exdates[i].Begin
			exdates[i].End = r.Exdates[i].End
		}

		set["exdates"] = exdates
	}
	if r.Exceptions != nil {
		set["exceptions"] = r.Exceptions
	}
	if r.Color != nil {
		set["color"] = *r.Color
	}
	if r.CorporatePattern != nil {
		set["entity_pattern"] = r.CorporatePattern.EntityPattern.RemoveFields(common.GetForbiddenFieldsInEntityPattern(mongo.PbehaviorMongoCollection))
		set["corporate_entity_pattern"] = r.CorporatePattern.ID
		set["corporate_entity_pattern_title"] = r.CorporatePattern.Title
	} else if r.EntityPattern != nil {
		set["entity_pattern"] = r.EntityPattern
		unset["corporate_entity_pattern"] = ""
		unset["corporate_entity_pattern_title"] = ""
	}

	if rruleUpdated {
		unset["rrule_cstart"] = ""
	}

	update := bson.M{"$set": set}
	if len(unset) > 0 {
		update["$unset"] = unset
	}

	var pbh *Response
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		pbh = nil

		prevPbh := pbehavior.PBehavior{}
		err := s.dbCollection.FindOne(ctx, bson.M{"_id": r.ID}).Decode(&prevPbh)
		if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
			return err
		}

		if prevPbh.Origin != "" {
			if len(prevPbh.Entities) > 0 {
				return common.NewValidationError("_id", "Cannot update an external pbehavior.")
			}

			valErr := common.NewValidationError("_id", "Cannot update a pbehavior with origin.")
			if r.Enabled != nil && !*r.Enabled ||
				r.RRule != nil && *r.RRule != "" ||
				len(r.Exdates) > 0 ||
				len(r.Exceptions) > 0 ||
				r.CorporateEntityPattern != nil && *r.CorporateEntityPattern != "" {

				return valErr
			}

			if r.EntityPattern != nil {
				if diff := pretty.Compare(prevPbh.EntityPattern, r.EntityPattern); diff != "" {
					return valErr
				}
			}
		}

		_, err = s.dbCollection.UpdateOne(ctx, bson.M{"_id": r.ID}, update)
		if err != nil {
			if mongodriver.IsDuplicateKeyError(err) {
				return s.parseDupError(err)
			}

			return err
		}

		pbh, err = s.GetOneBy(ctx, r.ID)
		if err != nil || pbh == nil {
			return err
		}

		if rruleUpdated {
			pbh.RRuleEnd, err = pbehavior.GetRruleEnd(*pbh.Start, pbh.RRule, s.timezoneConfigProvider.Get().Location)
			if err != nil {
				return err
			}

			if pbh.RRuleEnd == nil {
				_, err = s.dbCollection.UpdateOne(ctx, bson.M{"_id": r.ID}, bson.M{"$unset": bson.M{"rrule_end": ""}})
			} else {
				_, err = s.dbCollection.UpdateOne(ctx, bson.M{"_id": r.ID}, bson.M{"$set": bson.M{"rrule_end": pbh.RRuleEnd}})
			}
			if err != nil {
				return err
			}
		}

		return nil
	})

	return pbh, err
}

func (s *store) Delete(ctx context.Context, id, userID string) (bool, error) {
	var deleted int64

	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		deleted = 0

		// required to get the author in action log listener.
		res, err := s.dbCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"author": userID}})
		if err != nil || res.MatchedCount == 0 {
			return err
		}

		deleted, err = s.dbCollection.DeleteOne(ctx, bson.M{"_id": id})
		return err
	})

	return deleted > 0, err
}

func (s *store) DeleteByName(ctx context.Context, name, userID string) (string, error) {
	pbh := pbehavior.PBehavior{}

	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		pbh = pbehavior.PBehavior{}

		// required to get the author in action log listener.
		err := s.dbCollection.FindOneAndUpdate(
			ctx,
			bson.M{"name": name},
			bson.M{"$set": bson.M{"author": userID}},
		).Decode(&pbh)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		_, err = s.dbCollection.DeleteOne(ctx, bson.M{"_id": pbh.ID})

		return err
	})
	if err != nil {
		return "", err
	}

	return pbh.ID, nil
}

func (s *store) FindEntity(ctx context.Context, entityId string) (*libtypes.Entity, error) {
	entity := libtypes.Entity{}
	err := s.entityDbCollection.FindOne(ctx, bson.M{"_id": entityId}).Decode(&entity)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &entity, nil
}

func (s *store) EntityInsert(ctx context.Context, r BulkEntityCreateRequestItem) (*Response, error) {
	now := datetime.NewCpsTime()
	doc := pbehavior.PBehavior{
		ID:       utils.NewID(),
		Author:   r.Author,
		Comments: make([]*pbehavior.Comment, 0),
		Enabled:  true,
		Name:     r.Name,
		Reason:   r.Reason,
		Start:    r.Start,
		Stop:     r.Stop,
		Type:     r.Type,
		Created:  &now,
		Updated:  &now,
		Color:    r.Color,
		Origin:   r.Origin,
		Entity:   r.Entity,
		EntityPatternFields: savedpattern.EntityPatternFields{
			EntityPattern: pattern.Entity{
				{
					{
						Field: "_id",
						Condition: pattern.Condition{
							Type:  pattern.ConditionEqual,
							Value: r.Entity,
						},
					},
				},
			},
		},
	}

	if r.Comment != "" {
		doc.Comments = append(doc.Comments, &pbehavior.Comment{
			ID:        utils.NewID(),
			Author:    r.Author,
			Timestamp: &now,
			Message:   r.Comment,
		})
	}

	var pbh *Response
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		pbh = nil
		updateRes, err := s.dbCollection.UpdateOne(ctx,
			bson.M{
				"origin": r.Origin,
				"entity": r.Entity,
				"tstart": bson.M{"$lte": now},
				"tstop":  bson.M{"$gte": now},
			},
			bson.M{
				"$setOnInsert": doc,
			},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			if mongodriver.IsDuplicateKeyError(err) {
				return s.parseDupError(err)
			}

			return err
		}

		if updateRes.UpsertedCount == 0 {
			return common.NewValidationError("entity", "Pbehavior for origin already exists.")
		}

		pbh, err = s.GetOneBy(ctx, doc.ID)

		return err
	})

	return pbh, err
}

func (s *store) EntityDelete(ctx context.Context, r BulkEntityDeleteRequestItem) (string, error) {
	now := datetime.NewCpsTime()
	id := ""
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		id = ""
		var pbh pbehavior.PBehavior
		err := s.dbCollection.
			FindOneAndDelete(ctx, bson.M{
				"entity": r.Entity,
				"origin": r.Origin,
				"tstart": bson.M{"$lte": now},
				"tstop":  bson.M{"$gte": now},
			}, options.FindOneAndDelete().SetProjection(bson.M{"_id": 1})).
			Decode(&pbh)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		id = pbh.ID

		return nil
	})

	return id, err
}

func (s *store) ConnectorCreate(ctx context.Context, r BulkConnectorCreateRequestItem) (*Response, error) {
	entities := make([]string, 0, len(r.Entities))
	exists := make(map[string]struct{}, len(r.Entities))
	for _, entity := range r.Entities {
		if _, ok := exists[entity]; !ok {
			exists[entity] = struct{}{}
			entities = append(entities, entity)
		}
	}

	now := datetime.NewCpsTime()
	doc := pbehavior.PBehavior{
		ID:       utils.NewID(),
		Author:   r.Author,
		Enabled:  true,
		Name:     r.Name,
		Reason:   r.Reason,
		Start:    r.Start,
		Stop:     r.Stop,
		Type:     r.Type,
		Created:  &now,
		Updated:  &now,
		Color:    r.Color,
		Origin:   r.Origin,
		Entities: entities,
		EntityPatternFields: savedpattern.EntityPatternFields{
			EntityPattern: pattern.Entity{
				{
					{
						Field: "_id",
						Condition: pattern.Condition{
							Type:  pattern.ConditionIsOneOf,
							Value: entities,
						},
					},
				},
			},
		},
		Comments: []*pbehavior.Comment{
			{
				ID:        utils.NewID(),
				Author:    r.Author,
				Origin:    r.Origin,
				Timestamp: &now,
				Message:   r.Comment,
			},
		},
	}

	var pbh *Response
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		pbh = nil
		var findDoc pbehavior.PBehavior
		err := s.dbCollection.FindOneAndUpdate(ctx,
			bson.M{
				"origin":           r.Origin,
				"tstart":           r.Start,
				"tstop":            r.Stop,
				"comments.message": r.Comment,
				"entities":         bson.M{"$ne": nil},
			},
			[]bson.M{
				{"$set": bson.M{
					"entities": bson.M{"$cond": bson.M{
						"if": "$_id",
						"then": bson.M{"$setUnion": bson.A{
							"$entities",
							r.Entities,
						}},
						"else": nil,
					}},
				}},
				{"$set": bson.M{
					"entity_pattern": bson.M{"$cond": bson.M{
						"if": "$_id",
						"then": bson.A{bson.A{bson.M{
							"field": "_id",
							"cond": bson.M{
								"type":  pattern.ConditionIsOneOf,
								"value": "$entities",
							},
						}}},
						"else": nil,
					}},
				}},
				{
					"$replaceRoot": bson.M{"newRoot": bson.M{
						"$cond": bson.M{
							"if": "$_id",
							"then": bson.M{"$mergeObjects": bson.A{
								bson.M{
									"entity_pattern": "$entity_pattern",
									"entities":       "$entities",
									"author":         r.Author,
									"updated":        now,
								},
								"$$ROOT",
							}},
							"else": doc,
						},
					}},
				},
			},
			options.FindOneAndUpdate().
				SetUpsert(true).
				SetReturnDocument(options.After).
				SetProjection(bson.M{"_id": 1}),
		).Decode(&findDoc)
		if err != nil {
			if mongodriver.IsDuplicateKeyError(err) {
				return s.parseDupError(err)
			}

			return err
		}

		pbh, err = s.GetOneBy(ctx, findDoc.ID)

		return err
	})

	return pbh, err
}

func (s *store) ConnectorDelete(ctx context.Context, r BulkConnectorDeleteRequestItem) (string, error) {
	now := datetime.NewCpsTime()
	id := ""
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		id = ""
		var pbh pbehavior.PBehavior
		err := s.dbCollection.FindOneAndUpdate(ctx,
			bson.M{
				"origin":           r.Origin,
				"tstart":           r.Start,
				"tstop":            r.Stop,
				"comments.message": r.Comment,
				"entities":         bson.M{"$ne": nil},
			},
			bson.M{
				"$pullAll": bson.M{
					"entity_pattern.0.0.cond.value": r.Entities,
					"entities":                      r.Entities,
				},
				"$set": bson.M{
					"author":  r.Author,
					"updated": now,
				},
			},
			options.FindOneAndUpdate().SetProjection(bson.M{"_id": 1}),
		).Decode(&pbh)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		id = pbh.ID
		// Delete if entities are empty
		_, err = s.dbCollection.DeleteOne(ctx, bson.M{
			"_id":      id,
			"entities": bson.A{},
		})

		return err
	})

	return id, err
}

func (s *store) getMatchedPbhIDs(ctx context.Context, entity libtypes.Entity) ([]string, error) {
	cursor, err := s.dbCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	pbhIDs := make([]string, 0)

	for cursor.Next(ctx) {
		var pbh pbehavior.PBehavior
		err := cursor.Decode(&pbh)
		if err != nil {
			return nil, err
		}

		if len(pbh.EntityPattern) == 0 {
			continue
		}

		matched, err := libmatch.MatchEntityPattern(pbh.EntityPattern, &entity)
		if err != nil {
			return nil, err
		}

		if matched {
			pbhIDs = append(pbhIDs, pbh.ID)
		}
	}

	return pbhIDs, nil
}

func (s *store) transformRequestToDocument(r EditRequest) pbehavior.PBehavior {
	exdates := make([]pbehavior.Exdate, len(r.Exdates))
	for i := range r.Exdates {
		exdates[i].Type = r.Exdates[i].Type
		exdates[i].Begin = r.Exdates[i].Begin
		exdates[i].End = r.Exdates[i].End
	}

	return pbehavior.PBehavior{
		Author:     r.Author,
		Enabled:    *r.Enabled,
		Name:       r.Name,
		Reason:     r.Reason,
		RRule:      r.RRule,
		Start:      r.Start,
		Stop:       r.Stop,
		Type:       r.Type,
		Exdates:    exdates,
		Exceptions: r.Exceptions,
		Color:      r.Color,

		EntityPatternFields: r.EntityPatternFieldsRequest.ToModelWithoutFields(common.GetForbiddenFieldsInEntityPattern(mongo.PbehaviorMongoCollection)),
	}
}

func (s *store) transformResponse(ctx context.Context, result []Response) error {
	err := s.fillActiveStatuses(ctx, result)
	if err != nil {
		return err
	}

	loc := s.timezoneConfigProvider.Get().Location
	after := time.Now().In(loc)
	before := after.AddDate(0, nextEventMaxMonths, 0)
	for i, v := range result {
		if v.RRule == "" || v.RRuleEnd != nil && v.RRuleEnd.Time.Before(after) {
			continue
		}

		rOption, err := librrule.StrToROption(v.RRule)
		if err != nil {
			continue
		}

		if v.RRuleComputedStart != nil && v.RRuleComputedStart.Time.Before(after) {
			rOption.Dtstart = v.RRuleComputedStart.Time.In(loc)
		} else {
			rOption.Dtstart = v.Start.Time.In(loc)
		}
		r, err := librrule.NewRRule(*rOption)
		if err != nil {
			continue
		}

		iterator := r.Iterator()
		var next time.Time
		for {
			event, ok := iterator()
			if !ok || event.After(before) {
				break
			}
			if !event.Before(after) {
				next = event
				break
			}
		}

		if !next.IsZero() {
			if v.Stop != nil {
				d := v.Stop.Sub(v.Start.Time)
				result[i].Stop = &datetime.CpsTime{Time: next.Add(d)}
			}

			result[i].Start = &datetime.CpsTime{Time: next}
		}
	}

	return nil
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

func (s *store) parseDupError(err error) error {
	match := s.dupErrorRegexp.FindStringSubmatch(err.Error())
	if len(match) > 1 {
		matchedStr := match[1]

		switch matchedStr {
		case "name":
			return common.NewValidationError("name", "Name already exists.")
		case "_id":
			return common.NewValidationError("_id", "ID already exists.")
		default:
			return common.NewValidationError(matchedStr, matchedStr+" already exists.")
		}
	}

	return fmt.Errorf("can't parse duplication error: %w", err)
}

func sortCalendarResponse(response []CalendarResponse) func(i, j int) bool {
	return func(i, j int) bool {
		dateLeft := utils.DateOf(response[i].From.Time)
		dateRight := utils.DateOf(response[j].From.Time)

		if dateLeft.Before(dateRight) {
			return true
		}
		if dateLeft.After(dateRight) {
			return false
		}

		if response[i].Type.Priority > response[j].Type.Priority {
			return true
		}
		if response[i].Type.Priority < response[j].Type.Priority {
			return false
		}

		if response[i].From.Before(response[j].From) {
			return true
		}
		if response[i].From.After(response[j].From) {
			return false
		}

		if response[i].To.Before(response[j].To) {
			return true
		}
		if response[i].To.After(response[j].To) {
			return false
		}

		return response[i].Title < response[j].Title
	}
}
