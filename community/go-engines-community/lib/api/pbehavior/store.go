package pbehavior

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	libentity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/mongoquery"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	apipattern "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/patternfields"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/db"
	libmatch "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	libtypes "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	libredis "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/timespan"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/go-playground/validator/v10"
	"github.com/kylelemons/godebug/pretty"
	"github.com/redis/go-redis/v9"
	librrule "github.com/teambition/rrule-go"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/sync/errgroup"
)

const (
	nextEventMaxMonths = 1

	lockValue          = 1
	lockTickInterval   = 30 * time.Second
	lockExpirationTime = lockTickInterval + 10*time.Second
)

type Store interface {
	Insert(ctx context.Context, r CreateRequest) (*Response, error)
	Find(ctx context.Context, r ListRequest) (*AggregationResult, error)
	FindByEntityID(ctx context.Context, entity libtypes.Entity, r FindByEntityIDRequest) ([]Response, error)
	CalendarByEntityID(ctx context.Context, entity libtypes.Entity, r CalendarByEntityIDRequest) ([]CalendarResponse, error)
	GetOneBy(ctx context.Context, id string) (*Response, error)
	FindEntities(ctx context.Context, pbhID string, request EntitiesListRequest) (*AggregationEntitiesResult, error)
	Update(ctx context.Context, r UpdateRequest) (*Response, bool, error)
	UpdateByPatch(ctx context.Context, r PatchRequest) (*Response, bool, error)
	Delete(ctx context.Context, id, userID string) (bool, bool, error)
	DeleteByName(ctx context.Context, name, userID string) (string, bool, error)
	FindEntity(ctx context.Context, entityId string) (*libtypes.Entity, error)
	EntityInsert(ctx context.Context, r BulkEntityCreateRequestItem) (*Response, error)
	EntityDelete(ctx context.Context, r BulkEntityDeleteRequestItem) (string, error)
	ConnectorCreate(ctx context.Context, r BulkConnectorCreateRequestItem) (*Response, error)
	ConnectorDelete(ctx context.Context, r BulkConnectorDeleteRequestItem) (string, error)
	ExecPatternAndUpdate(ctx context.Context, id string, pattern pattern.Entity) (*apipattern.CountResponse, error)
	ExecPatternsAndUpdate(ctx context.Context) error
}

type store struct {
	dbClient           mongo.DbClient
	readDbClient       mongo.DbClient
	redisClient        redis.Cmdable
	dbCollection       mongo.DbCollection
	entityDbCollection mongo.DbCollection

	authorProvider         author.Provider
	entityTypeResolver     pbehavior.EntityTypeResolver
	pbhTypeComputer        pbehavior.TypeComputer
	timezoneConfigProvider config.TimezoneConfigProvider

	websocketHub                  websocket.Hub
	userInterfaceConfigProvider   config.UserInterfaceConfigProvider
	defaultSortBy                 string
	entitiesDefaultSearchByFields []string
	entitiesDefaultSortBy         string

	transformer    patternfields.Transformer
	dupErrorParser validation.DuplicateErrorParser

	workers int
}

func NewStore(
	dbClient mongo.DbClient,
	readDbClient mongo.DbClient,
	redisClient redis.Cmdable,
	entityTypeResolver pbehavior.EntityTypeResolver,
	pbhTypeComputer pbehavior.TypeComputer,
	timezoneConfigProvider config.TimezoneConfigProvider,
	authorProvider author.Provider,
	transformer patternfields.Transformer,
	websocketHub websocket.Hub,
	userInterfaceConfigProvider config.UserInterfaceConfigProvider,
) Store {
	return &store{
		dbClient:                      dbClient,
		dbCollection:                  dbClient.Collection(mongo.PbehaviorMongoCollection),
		entityDbCollection:            dbClient.Collection(mongo.EntityMongoCollection),
		readDbClient:                  readDbClient,
		redisClient:                   redisClient,
		entityTypeResolver:            entityTypeResolver,
		pbhTypeComputer:               pbhTypeComputer,
		timezoneConfigProvider:        timezoneConfigProvider,
		authorProvider:                authorProvider,
		transformer:                   transformer,
		websocketHub:                  websocketHub,
		userInterfaceConfigProvider:   userInterfaceConfigProvider,
		defaultSortBy:                 "created",
		entitiesDefaultSearchByFields: []string{"_id", "name", "type"},
		entitiesDefaultSortBy:         "_id",
		dupErrorParser:                validation.NewDuplicateErrorParser(),
		workers:                       10,
	}
}

func (s *store) Insert(ctx context.Context, r CreateRequest) (*Response, error) {
	now := datetime.NewCpsTime()
	doc := s.transformRequestToDocument(r.EditRequest)
	doc.ID = cmp.Or(r.ID, utils.NewID())
	loc := s.timezoneConfigProvider.Get().Location
	var err error
	if doc.Timezone != "" {
		loc, err = time.LoadLocation(doc.Timezone)
		if err != nil {
			return nil, fmt.Errorf("invalid timezone %q: %w", doc.Timezone, err)
		}
	}

	rruleEnd, err := pbehavior.GetRruleEnd(*r.Start, r.RRule, loc)
	if err != nil {
		return nil, err
	}

	doc.Created = &now
	doc.Updated = &now
	doc.Comments = make([]pbehavior.Comment, 0)
	doc.RRuleEnd = rruleEnd
	if r.ExecPattern {
		doc.EntityPatternFields, doc.Aliases, err = s.transformPatternRequestToModel(ctx, r.EntityRequest, r)
		if err != nil {
			return nil, err
		}

		_, doc.PatternMs, err = s.execPattern(ctx, doc.EntityPattern)
		if err != nil {
			return nil, err
		}

		doc.PatternExecAt = &now
	}

	var pbh *Response
	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		pbh = nil

		doc.EntityPatternFields, doc.Aliases, err = s.transformPatternRequestToModel(ctx, r.EntityRequest, r)
		if err != nil {
			return err
		}

		_, err := s.dbCollection.InsertOne(ctx, doc)
		if err != nil {
			if mongodriver.IsDuplicateKeyError(err) {
				return s.dupErrorParser.Parse(err, Response{})
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
	pipeline = append(pipeline, mongoquery.GetSortQuery("created", pagination.SortAsc))
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

	span := timespan.New(r.From.Time, r.To.Time)
	computed, err := s.pbhTypeComputer.ComputeByIds(ctx, span, pbhIDs, s.timezoneConfigProvider.Get().Location)
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

	sort.Slice(res, s.sortCalendarResponse(res))

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
	filter := mongoquery.GetSearchQuery(request.Search, s.entitiesDefaultSearchByFields)
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
		mongoquery.GetSortQuery(cmp.Or(request.SortBy, s.entitiesDefaultSortBy), request.Sort),
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

func (s *store) Update(ctx context.Context, r UpdateRequest) (*Response, bool, error) {
	now := datetime.NewCpsTime()
	doc := s.transformRequestToDocument(r.EditRequest)
	doc.Updated = &now

	unset := bson.M{
		"rrule_cstart": "",
	}

	if r.Stop == nil {
		unset["tstop"] = ""
	}

	loc := s.timezoneConfigProvider.Get().Location
	var err error
	if doc.Timezone != "" {
		loc, err = time.LoadLocation(doc.Timezone)
		if err != nil {
			return nil, false, fmt.Errorf("invalid timezone %q: %w", doc.Timezone, err)
		}
	}

	rruleEnd, err := pbehavior.GetRruleEnd(*r.Start, r.RRule, loc)
	if err != nil {
		return nil, false, err
	}
	if rruleEnd == nil {
		unset["rrule_end"] = ""
	} else {
		doc.RRuleEnd = rruleEnd
	}

	if r.ExecPattern {
		doc.EntityPatternFields, doc.Aliases, err = s.transformPatternRequestToModel(ctx, r.EntityRequest, r)
		if err != nil {
			return nil, false, err
		}

		_, doc.PatternMs, err = s.execPattern(ctx, doc.EntityPattern)
		if err != nil {
			return nil, false, err
		}

		doc.PatternExecAt = &now
	}

	update := make(bson.M)
	if len(unset) > 0 {
		update["$unset"] = unset
	}

	var pbh *Response
	var recomputeInherited bool

	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		pbh = nil
		recomputeInherited = false

		prevPbh := pbehavior.PBehavior{}
		err := s.dbCollection.FindOne(ctx, bson.M{"_id": r.ID}).Decode(&prevPbh)
		if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
			return err
		}

		recomputeInherited = prevPbh.Inherited

		if prevPbh.Origin != "" {
			if len(prevPbh.Entities) > 0 {
				return httperror.NewForbiddenError("The external pbehavior cannot be modified.")
			}

			var fieldErrs validator.ValidationErrors
			if !*r.Enabled {
				fieldErrs = append(fieldErrs, validation.NewFieldError("unchangeable", "Enabled", "Enabled"))
			}

			if r.RRule != "" {
				fieldErrs = append(fieldErrs, validation.NewFieldError("unchangeable", "RRule", "RRule"))
			}

			if r.Timezone != "" {
				fieldErrs = append(fieldErrs, validation.NewFieldError("unchangeable", "Timezone", "Timezone"))
			}

			if len(r.Exdates) > 0 {
				fieldErrs = append(fieldErrs, validation.NewFieldError("unchangeable", "Exdates", "Exdates"))
			}

			if len(r.Exceptions) > 0 {
				fieldErrs = append(fieldErrs, validation.NewFieldError("unchangeable", "Exceptions", "Exceptions"))
			}

			if r.CorporateEntityPattern != "" {
				fieldErrs = append(fieldErrs, validation.NewFieldError("unchangeable", "CorporateEntityPattern", "CorporateEntityPattern"))
			}

			if diff := pretty.Compare(prevPbh.EntityPattern, r.EntityPattern); diff != "" {
				fieldErrs = append(fieldErrs, validation.NewFieldError("unchangeable", "EntityPattern", "EntityPattern"))
			}

			if len(fieldErrs) > 0 {
				return validation.NewError(fieldErrs, r)
			}
		}

		doc.EntityPatternFields, doc.Aliases, err = s.transformPatternRequestToModel(ctx, r.EntityRequest, r)
		if err != nil {
			return err
		}

		update["$set"] = doc

		_, err = s.dbCollection.UpdateOne(ctx, bson.M{"_id": r.ID}, update)
		if err != nil {
			if mongodriver.IsDuplicateKeyError(err) {
				return s.dupErrorParser.Parse(err, Response{})
			}

			return err
		}

		pbh, err = s.GetOneBy(ctx, r.ID)

		return err
	})

	if pbh != nil {
		recomputeInherited = recomputeInherited || pbh.Inherited
	}

	return pbh, recomputeInherited, err
}

func (s *store) UpdateByPatch(ctx context.Context, r PatchRequest) (*Response, bool, error) {
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
	if r.Timezone != nil {
		set["timezone"] = *r.Timezone
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

	if r.Inherited != nil {
		set["inherited"] = *r.Inherited
	}

	if rruleUpdated {
		unset["rrule_cstart"] = ""
	}

	var pbh *Response
	var recomputeInherited bool

	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		pbh = nil
		recomputeInherited = false

		prevPbh := pbehavior.PBehavior{}
		err := s.dbCollection.FindOne(ctx, bson.M{"_id": r.ID}).Decode(&prevPbh)
		if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
			return err
		}

		recomputeInherited = prevPbh.Inherited

		if prevPbh.Origin != "" {
			if len(prevPbh.Entities) > 0 {
				return httperror.NewForbiddenError("The external pbehavior cannot be modified.")
			}

			var fieldErrs validator.ValidationErrors
			if r.Enabled != nil && !*r.Enabled {
				fieldErrs = append(fieldErrs, validation.NewFieldError("unchangeable", "Enabled", "Enabled"))
			}

			if r.RRule != nil && *r.RRule != "" {
				fieldErrs = append(fieldErrs, validation.NewFieldError("unchangeable", "RRule", "RRule"))
			}

			if r.Timezone != nil && *r.Timezone != "" {
				fieldErrs = append(fieldErrs, validation.NewFieldError("unchangeable", "Timezone", "Timezone"))
			}

			if len(r.Exdates) > 0 {
				fieldErrs = append(fieldErrs, validation.NewFieldError("unchangeable", "Exdates", "Exdates"))
			}

			if len(r.Exceptions) > 0 {
				fieldErrs = append(fieldErrs, validation.NewFieldError("unchangeable", "Exceptions", "Exceptions"))
			}

			if r.CorporateEntityPattern != nil && *r.CorporateEntityPattern != "" {
				fieldErrs = append(fieldErrs, validation.NewFieldError("unchangeable", "CorporateEntityPattern", "CorporateEntityPattern"))
			}

			if r.EntityPattern != nil {
				if diff := pretty.Compare(prevPbh.EntityPattern, r.EntityPattern); diff != "" {
					fieldErrs = append(fieldErrs, validation.NewFieldError("unchangeable", "EntityPattern", "EntityPattern"))
				}
			}

			if len(fieldErrs) > 0 {
				return validation.NewError(fieldErrs, r)
			}
		}

		if r.CorporateEntityPattern != nil || r.EntityPattern != nil {
			corpPattern := ""
			if r.CorporateEntityPattern != nil {
				corpPattern = *r.CorporateEntityPattern
			}

			epf, aliases, err := s.transformPatternRequestToModel(ctx, patternfields.EntityRequest{
				CorporateEntityPattern: corpPattern,
				EntityPattern:          r.EntityPattern,
			}, r)
			if err != nil {
				return err
			}

			set["aliases"] = aliases
			set["entity_pattern"] = epf.EntityPattern
			set["corporate_entity_pattern"] = epf.CorporateEntityPattern
			set["corporate_entity_pattern_title"] = epf.CorporateEntityPatternTitle
		}

		update := bson.M{"$set": set}
		if len(unset) > 0 {
			update["$unset"] = unset
		}

		_, err = s.dbCollection.UpdateOne(ctx, bson.M{"_id": r.ID}, update)
		if err != nil {
			if mongodriver.IsDuplicateKeyError(err) {
				return s.dupErrorParser.Parse(err, Response{})
			}

			return err
		}

		pbh, err = s.GetOneBy(ctx, r.ID)
		if err != nil || pbh == nil {
			return err
		}

		if rruleUpdated {
			loc := s.timezoneConfigProvider.Get().Location
			if pbh.Timezone != "" {
				loc, err = time.LoadLocation(pbh.Timezone)
				if err != nil {
					return fmt.Errorf("invalid timezone %q: %w", pbh.Timezone, err)
				}
			}

			pbh.RRuleEnd, err = pbehavior.GetRruleEnd(*pbh.Start, pbh.RRule, loc)
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

	if pbh != nil {
		recomputeInherited = recomputeInherited || pbh.Inherited
	}

	return pbh, recomputeInherited, err
}

func (s *store) Delete(ctx context.Context, id, userID string) (bool, bool, error) {
	var deleted int64
	var recomputeInherited bool

	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		deleted = 0
		recomputeInherited = false

		var prevPbh pbehavior.PBehavior

		// required to get the author in action log listener.
		err := s.dbCollection.FindOneAndUpdate(
			ctx,
			bson.M{"_id": id},
			bson.M{"$set": bson.M{"author": userID}},
			options.FindOneAndUpdate().SetReturnDocument(options.Before).SetProjection(bson.M{"inherited": 1}),
		).Decode(&prevPbh)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		recomputeInherited = prevPbh.Inherited

		deleted, err = s.dbCollection.DeleteOne(ctx, bson.M{"_id": id})
		return err
	})

	return deleted > 0, recomputeInherited, err
}

func (s *store) DeleteByName(ctx context.Context, name, userID string) (string, bool, error) {
	pbh := pbehavior.PBehavior{}
	var recomputeInherited bool

	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		pbh = pbehavior.PBehavior{}
		recomputeInherited = false

		// required to get the author in action log listener.
		err := s.dbCollection.FindOneAndUpdate(
			ctx,
			bson.M{"name": name},
			bson.M{"$set": bson.M{"author": userID}},
			options.FindOneAndUpdate().SetReturnDocument(options.Before).SetProjection(bson.M{"_id": 1, "inherited": 1}),
		).Decode(&pbh)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		recomputeInherited = pbh.Inherited

		_, err = s.dbCollection.DeleteOne(ctx, bson.M{"_id": pbh.ID})

		return err
	})
	if err != nil {
		return "", false, err
	}

	return pbh.ID, recomputeInherited, nil
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
		Comments: make([]pbehavior.Comment, 0),
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
		doc.Comments = append(doc.Comments, pbehavior.Comment{
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
			options.UpdateOne().SetUpsert(true),
		)
		if err != nil {
			if mongodriver.IsDuplicateKeyError(err) {
				return s.dupErrorParser.Parse(err, Response{})
			}

			return err
		}

		if updateRes.UpsertedCount == 0 {
			return validation.NewError(
				validator.ValidationErrors{validation.NewFieldError("exist", "Origin", "Origin")},
				r,
			)
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
				"$or": bson.A{
					bson.M{"tstop": nil},
					bson.M{"tstop": bson.M{"$gte": now}},
				},
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
		Comments: []pbehavior.Comment{
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
							"else": bson.M{"$literal": doc},
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
				return s.dupErrorParser.Parse(err, Response{})
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

func (s *store) ExecPatternAndUpdate(ctx context.Context, id string, pattern pattern.Entity) (*apipattern.CountResponse, error) {
	conf := s.userInterfaceConfigProvider.Get()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(conf.CheckCountRequestTimeout)*time.Second)
	defer cancel()
	updateStats := false
	if id != "" {
		pbh, err := s.GetOneBy(ctx, id)
		if err != nil || pbh == nil {
			return nil, err
		}

		updateStats = reflect.DeepEqual(pbh.EntityPattern, pattern)
	}

	count, ms, err := s.execPattern(ctx, pattern)
	if err != nil {
		return nil, err
	}

	if updateStats {
		now := datetime.NewCpsTime()
		_, err = s.dbCollection.UpdateOne(ctx,
			bson.M{
				"_id": id,
				"$or": []bson.M{
					{"pattern_exec_at": nil},
					{"pattern_exec_at": bson.M{"$lt": now}},
				},
			},
			bson.M{"$set": bson.M{
				"pattern_ms":      ms,
				"pattern_exec_at": now,
			}},
		)
		if err != nil {
			return nil, err
		}
	}

	res := apipattern.CountResponse{
		Count:     count,
		OverLimit: count > int64(conf.MaxMatchedItems),
		Millisecs: ms,
	}

	return &res, nil
}

func (s *store) ExecPatternsAndUpdate(ctx context.Context) (resErr error) {
	res := s.redisClient.SetNX(ctx, libredis.ApiPbhPatternCountLockKey, lockValue, lockExpirationTime)
	if err := res.Err(); err != nil {
		return err
	}

	if !res.Val() {
		return nil
	}

	defer func() {
		err := s.redisClient.Del(ctx, libredis.ApiPbhPatternCountLockKey).Err()
		if err != nil && resErr == nil {
			resErr = err
		}
	}()

	g, gCtx := errgroup.WithContext(ctx)
	ch := make(chan pbehavior.PBehavior)
	g.Go(func() error {
		defer close(ch)
		cursor, err := s.readDbClient.Collection(mongo.PbehaviorMongoCollection).
			Find(gCtx, bson.M{}, options.Find().SetProjection(bson.M{"entity_pattern": 1}))
		if err != nil {
			return err
		}

		defer cursor.Close(gCtx)
		for cursor.Next(gCtx) {
			pbh := pbehavior.PBehavior{}
			err = cursor.Decode(&pbh)
			if err != nil {
				return err
			}

			select {
			case <-gCtx.Done():
			case ch <- pbh:
			}
		}

		if err = cursor.Err(); err != nil {
			return err
		}

		return nil
	})
	now := datetime.NewCpsTime()
	done := make(chan struct{})
	defer close(done)
	for i := 0; i < s.workers; i++ {
		g.Go(func() error {
			for {
				select {
				case <-gCtx.Done():
					return nil
				case pbh, ok := <-ch:
					if !ok {
						select {
						case <-gCtx.Done():
						case done <- struct{}{}:
						}

						return nil
					}

					conf := s.userInterfaceConfigProvider.Get()
					execCtx, cancel := context.WithTimeout(gCtx, time.Duration(conf.CheckCountRequestTimeout)*time.Second)
					_, ms, err := s.execPattern(execCtx, pbh.EntityPattern)
					if err != nil {
						cancel()

						return err
					}

					cancel()
					_, err = s.dbCollection.UpdateOne(gCtx,
						bson.M{
							"_id": pbh.ID,
							"$or": []bson.M{
								{"pattern_exec_at": nil},
								{"pattern_exec_at": bson.M{"$lt": now}},
							},
						},
						bson.M{"$set": bson.M{
							"pattern_ms":      ms,
							"pattern_exec_at": now,
						}},
					)
					if err != nil {
						return err
					}
				}
			}
		})
	}

	g.Go(func() error {
		ticker := time.NewTicker(lockTickInterval)
		defer ticker.Stop()
		doneCount := 0
		for {
			select {
			case <-done:
				doneCount++
				if doneCount == s.workers {
					return nil
				}
			case <-gCtx.Done():
				return nil
			case <-ticker.C:
				err := s.redisClient.SetEx(gCtx, libredis.ApiCleanEntitiesLockKey, lockValue, lockExpirationTime).Err()
				if err != nil {
					return err
				}
			}
		}
	})

	err := g.Wait()
	if err != nil {
		s.websocketHub.Send(websocket.RoomPbhPatterns, map[string]bool{"ok": false})

		return err
	}

	s.websocketHub.Send(websocket.RoomPbhPatterns, map[string]bool{"ok": true})

	return nil
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
		Timezone:   r.Timezone,
		Start:      r.Start,
		Stop:       r.Stop,
		Type:       r.Type,
		Exdates:    exdates,
		Exceptions: r.Exceptions,
		Color:      r.Color,
		Inherited:  r.Inherited,
	}
}

func (s *store) transformResponse(ctx context.Context, result []Response) error {
	err := s.fillActiveStatuses(ctx, result)
	if err != nil {
		return err
	}

	defLoc := s.timezoneConfigProvider.Get().Location
	after := time.Now().In(defLoc)
	before := after.AddDate(0, nextEventMaxMonths, 0)
	locs := make(map[string]*time.Location)
	for i, v := range result {
		if v.RRule == "" || v.RRuleEnd != nil && v.RRuleEnd.Time.Before(after) {
			continue
		}

		loc := defLoc
		if v.Timezone != "" {
			var ok bool
			if loc, ok = locs[v.Timezone]; !ok {
				loc, err = time.LoadLocation(v.Timezone)
				if err != nil {
					return fmt.Errorf("invalid timezone %q for pbehavior id=%q: %w", v.Timezone, v.ID, err)
				}

				locs[v.Timezone] = loc
			}
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
	now := time.Now()
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

func (s *store) execPattern(ctx context.Context, entityPattern pattern.Entity) (int64, int64, error) {
	q, err := db.EntityPatternToMongoQuery(entityPattern, "")
	if err != nil {
		return 0, 0, err
	}

	pipeline := []bson.M{
		{"$match": q},
		{"$count": "total_count"},
	}
	collection := s.readDbClient.Collection(mongo.EntityMongoCollection)
	start := time.Now()
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, 0, err
	}

	defer cursor.Close(ctx)
	res := struct {
		Count int64 `bson:"total_count"`
	}{}
	if cursor.Next(ctx) {
		err = cursor.Decode(&res)
		if err != nil {
			return 0, 0, err
		}
	}

	if err = cursor.Err(); err != nil {
		return 0, 0, err
	}

	return res.Count, max(time.Since(start).Milliseconds(), 1), nil
}

func (s *store) sortCalendarResponse(response []CalendarResponse) func(i, j int) bool {
	return func(i, j int) bool {
		if response[i].From.Before(response[j].From) {
			return true
		}

		if response[i].From.After(response[j].From) {
			return false
		}

		if response[i].Type.Priority > response[j].Type.Priority {
			return true
		}

		if response[i].Type.Priority < response[j].Type.Priority {
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

func (s *store) transformPatternRequestToModel(ctx context.Context, er patternfields.EntityRequest, r any) (epf savedpattern.EntityPatternFields, aliasPropIDs []string, err error) {
	return s.transformer.TransformEntityRequest(ctx, er, r, s.dbCollection.Name())
}
