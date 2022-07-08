package pbehavior

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	libtypes "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/timespan"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface {
	Insert(ctx context.Context, r CreateRequest) (*Response, error)
	Find(ctx context.Context, r ListRequest) (*AggregationResult, error)
	FindByEntityID(ctx context.Context, entity libtypes.Entity) ([]Response, error)
	CalendarByEntityID(ctx context.Context, entity libtypes.Entity, r CalendarByEntityIDRequest) ([]CalendarResponse, error)
	GetOneBy(ctx context.Context, id string) (*Response, error)
	FindEntities(ctx context.Context, pbhID string, request EntitiesListRequest) (*AggregationEntitiesResult, error)
	Update(ctx context.Context, r UpdateRequest) (*Response, error)
	UpdateByPatch(ctx context.Context, r PatchRequest) (*Response, error)
	Delete(ctx context.Context, id string) (bool, error)
	DeleteByName(ctx context.Context, name string) (string, error)
	FindEntity(ctx context.Context, entityId string) (*libtypes.Entity, error)
}

type store struct {
	dbClient mongo.DbClient

	dbCollection, entityCollection mongo.DbCollection

	entityMatcher          pbehavior.EntityMatcher
	entityTypeResolver     pbehavior.EntityTypeResolver
	pbhTypeComputer        pbehavior.TypeComputer
	timezoneConfigProvider config.TimezoneConfigProvider
	defaultSortBy          string

	entitiesDefaultSearchByFields []string
	entitiesDefaultSortBy         string
}

func NewStore(
	dbClient mongo.DbClient,
	entityMatcher pbehavior.EntityMatcher,
	entityTypeResolver pbehavior.EntityTypeResolver,
	pbhTypeComputer pbehavior.TypeComputer,
	timezoneConfigProvider config.TimezoneConfigProvider,
) Store {
	return &store{
		dbClient:                      dbClient,
		dbCollection:                  dbClient.Collection(mongo.PbehaviorMongoCollection),
		entityCollection:              dbClient.Collection(mongo.EntityMongoCollection),
		entityMatcher:                 entityMatcher,
		entityTypeResolver:            entityTypeResolver,
		pbhTypeComputer:               pbhTypeComputer,
		timezoneConfigProvider:        timezoneConfigProvider,
		defaultSortBy:                 "created",
		entitiesDefaultSearchByFields: []string{"_id", "name", "type"},
		entitiesDefaultSortBy:         "_id",
	}
}

func (s *store) Insert(ctx context.Context, r CreateRequest) (*Response, error) {
	now := libtypes.NewCpsTime(time.Now().Unix())
	doc := s.transformRequestToDocument(r.EditRequest)
	doc.ID = r.ID
	if doc.ID == "" {
		doc.ID = utils.NewID()
	}

	doc.Created = &now
	doc.Updated = &now
	doc.Comments = make([]*pbehavior.Comment, 0)

	var pbh *Response
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		pbh = nil
		_, err := s.dbCollection.InsertOne(ctx, doc)
		if err != nil {
			return err
		}

		pbh, err = s.GetOneBy(ctx, doc.ID)
		return err
	})

	return pbh, err
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

func (s *store) FindByEntityID(ctx context.Context, entity libtypes.Entity) ([]Response, error) {
	pbhIDs, err := s.getMatchedPbhIDs(ctx, entity)
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

	res := make([]Response, 0)
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
				From:  libtypes.CpsTime{Time: computedType.Span.From()},
				To:    libtypes.CpsTime{Time: computedType.Span.To()},
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
	pipeline = append(pipeline, GetNestedObjectsPipeline()...)
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

	var match interface{}
	if len(pbh.OldMongoQuery) > 0 {
		match = pbh.OldMongoQuery
	} else {
		match, err = pbh.EntityPattern.ToMongoQuery("")
		if err != nil {
			return nil, err
		}
	}
	pipeline := []bson.M{
		{"$match": match},
	}
	filter := common.GetSearchQuery(request.Search, s.entitiesDefaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	sortBy := request.SortBy
	if sortBy == "" {
		sortBy = s.entitiesDefaultSortBy
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
	cursor, err := s.entityCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		request.Query,
		pipeline,
		common.GetSortQuery(sortBy, request.Sort),
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
	now := libtypes.NewCpsTime(time.Now().Unix())
	doc := s.transformRequestToDocument(r.EditRequest)
	doc.Updated = &now

	unset := bson.M{}

	if r.Stop == nil {
		unset["tstop"] = ""
	}

	if r.CorporateEntityPattern != "" || len(r.EntityPattern) > 0 {
		unset["old_mongo_query"] = ""
	}

	update := bson.M{"$set": doc}
	if len(unset) > 0 {
		update["$unset"] = unset
	}

	var pbh *Response
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		pbh = nil
		_, err := s.dbCollection.UpdateOne(ctx, bson.M{"_id": r.ID}, update)
		if err != nil {
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
		"updated": libtypes.NewCpsTime(),
	}
	unset := bson.M{}
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
	}
	if r.Start != nil {
		set["tstart"] = *r.Start
	}
	if r.Stop.isSet {
		if r.Stop.val == nil {
			unset["tstop"] = ""
		} else {
			set["tstop"] = *r.Stop.val
		}
	}
	if r.Exdates != nil {
		set["exdates"] = r.Exdates
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

	update := bson.M{"$set": set}
	if len(unset) > 0 {
		update["$unset"] = unset
	}

	var pbh *Response
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		pbh = nil
		_, err := s.dbCollection.UpdateOne(ctx, bson.M{"_id": r.ID}, update)
		if err != nil {
			return err
		}

		pbh, err = s.GetOneBy(ctx, r.ID)
		return err
	})

	return pbh, err
}

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	deleted, err := s.dbCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}

	return deleted > 0, nil
}

func (s *store) DeleteByName(ctx context.Context, name string) (string, error) {
	pbh := pbehavior.PBehavior{}
	err := s.dbCollection.FindOne(ctx, bson.M{"name": name}).Decode(&pbh)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return "", nil
		}
		return "", err
	}

	deleted, err := s.dbCollection.DeleteOne(ctx, bson.M{"_id": pbh.ID})
	if err != nil || deleted == 0 {
		return "", err
	}

	return pbh.ID, nil
}

func (s *store) FindEntity(ctx context.Context, entityId string) (*libtypes.Entity, error) {
	entity := libtypes.Entity{}
	err := s.entityCollection.FindOne(ctx, bson.M{"_id": entityId}).Decode(&entity)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &entity, nil
}

func (s *store) getMatchedPbhIDs(ctx context.Context, entity libtypes.Entity) ([]string, error) {
	cursor, err := s.dbCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	pbhIDs := make([]string, 0)
	filters := make(map[string]interface{})

	for cursor.Next(ctx) {
		var pbh pbehavior.PBehavior
		err := cursor.Decode(&pbh)
		if err != nil {
			return nil, err
		}

		if len(pbh.EntityPattern) > 0 {
			matched, _, err := pbh.EntityPattern.Match(entity)
			if err != nil {
				return nil, err
			}

			if matched {
				pbhIDs = append(pbhIDs, pbh.ID)
			}

			continue
		}

		var oldMongoQuery map[string]interface{}
		err = json.Unmarshal([]byte(pbh.OldMongoQuery), &oldMongoQuery)
		if err != nil {
			return nil, err
		}
		filters[pbh.ID] = oldMongoQuery
	}

	if len(filters) > 0 {
		ids, err := s.entityMatcher.MatchAll(ctx, entity.ID, filters)
		if err != nil {
			return nil, err
		}

		pbhIDs = append(pbhIDs, ids...)
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

	exceptions := make([]string, len(r.Exceptions))
	copy(exceptions, r.Exceptions)

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
		Exceptions: exceptions,
		Color:      r.Color,

		EntityPatternFields: r.EntityPatternFieldsRequest.ToModelWithoutFields(common.GetForbiddenFieldsInEntityPattern(mongo.PbehaviorMongoCollection)),
	}
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
