package pbehaviorexception

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	ics "github.com/apognu/gocal"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const maxImportDateInYears = 20

type Store interface {
	Insert(ctx context.Context, r CreateRequest) (*Response, error)
	Find(ctx context.Context, r ListRequest) (*AggregationResult, error)
	GetById(ctx context.Context, id string) (*Response, error)
	Update(ctx context.Context, r UpdateRequest) (*Response, error)
	Delete(ctx context.Context, id, userId string) (bool, error)
	IsLinked(ctx context.Context, id string) (bool, error)
	Import(ctx context.Context, name, pbhType, userId string, f multipart.File, fh *multipart.FileHeader) (*Response, error)
}

func NewStore(dbClient mongo.DbClient, timezoneConfigProvider config.TimezoneConfigProvider, authorProvider author.Provider) Store {
	return &store{
		dbClient:              dbClient,
		dbCollection:          dbClient.Collection(mongo.PbehaviorExceptionMongoCollection),
		pbehaviorDbCollection: dbClient.Collection(mongo.PbehaviorMongoCollection),
		typeDbCollection:      dbClient.Collection(mongo.PbehaviorTypeMongoCollection),
		authorProvider:        authorProvider,
		defaultSearchByFields: []string{"_id", "name", "description"},
		defaultSortBy:         "created",

		timezoneConfigProvider: timezoneConfigProvider,
	}
}

type store struct {
	dbClient              mongo.DbClient
	dbCollection          mongo.DbCollection
	pbehaviorDbCollection mongo.DbCollection
	typeDbCollection      mongo.DbCollection
	defaultSearchByFields []string
	defaultSortBy         string

	timezoneConfigProvider config.TimezoneConfigProvider
	authorProvider         author.Provider
}

func (s *store) Insert(ctx context.Context, r CreateRequest) (*Response, error) {
	doc, err := s.transformRequestToDoc(ctx, r.EditRequest)
	if err != nil {
		return nil, err
	}

	doc.ID = cmp.Or(r.ID, utils.NewID())
	doc.Created = datetime.NewCpsTime()
	doc.Updated = datetime.NewCpsTime()

	var res *Response

	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		res = nil

		_, err := s.dbCollection.InsertOne(ctx, doc)
		if err != nil {
			return err
		}

		res, err = s.GetById(ctx, doc.ID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *store) Find(ctx context.Context, r ListRequest) (*AggregationResult, error) {
	pipeline := make([]bson.M, 0)
	filter := common.GetSearchQuery(r.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	sort := common.GetSortQuery(cmp.Or(r.SortBy, s.defaultSortBy), r.Sort)
	project := getNestedObjectsPipeline()
	project = append(project, s.authorProvider.Pipeline()...)
	project = append(project, sort)
	if r.WithFlags {
		project = append(project, getDeletablePipeline()...)
	}

	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		pipeline,
		sort,
		project,
	), options.Aggregate().SetCollation(&options.Collation{Locale: "en"}))

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	cursor.Next(ctx)

	var result AggregationResult
	err = cursor.Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *store) GetById(ctx context.Context, id string) (*Response, error) {
	pipeline := []bson.M{{"$match": bson.M{"_id": id}}}

	pipeline = append(pipeline, s.authorProvider.Pipeline()...)
	pipeline = append(pipeline, getNestedObjectsPipeline()...)

	cursor, err := s.dbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	if cursor.Next(ctx) {
		var res Response

		err := cursor.Decode(&res)
		if err != nil {
			return nil, err
		}

		return &res, nil
	}

	return nil, nil
}

func (s *store) Update(ctx context.Context, r UpdateRequest) (*Response, error) {
	doc, err := s.transformRequestToDoc(ctx, r.EditRequest)
	if err != nil {
		return nil, err
	}

	doc.ID = r.ID
	doc.Updated = datetime.NewCpsTime()

	var res *Response

	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		res = nil

		result, err := s.dbCollection.UpdateOne(ctx, bson.M{"_id": doc.ID}, bson.M{"$set": doc})
		if err != nil || result.MatchedCount == 0 {
			return err
		}

		res, err = s.GetById(ctx, doc.ID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *store) Delete(ctx context.Context, id, userId string) (bool, error) {
	isLinked, err := s.IsLinked(ctx, id)
	if err != nil {
		return false, err
	}

	if isLinked {
		return false, ErrLinkedException
	}

	var deleted int64

	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		deleted = 0

		// required to get the author in action log listener.
		res, err := s.dbCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"author": userId}})
		if err != nil || res.MatchedCount == 0 {
			return err
		}

		deleted, err = s.dbCollection.DeleteOne(ctx, bson.M{"_id": id})
		return err
	})

	return deleted > 0, err
}

// IsLinked checks if there is pbehavior with linked exception.
func (s *store) IsLinked(ctx context.Context, id string) (bool, error) {
	res := s.pbehaviorDbCollection.FindOne(ctx, bson.M{"exceptions": id})
	if err := res.Err(); err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (s *store) Import(ctx context.Context, name, pbhType, userId string, f multipart.File, fh *multipart.FileHeader) (*Response, error) {
	err := s.typeDbCollection.FindOne(ctx, bson.M{"_id": pbhType}).Err()
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, common.NewValidationError("type", "Type doesn't exist.")
		}
		return nil, err
	}
	err = s.dbCollection.FindOne(ctx, bson.M{"name": name}).Err()
	if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
		return nil, err
	}
	if err == nil {
		return nil, common.NewValidationError("name", "Name already exists.")
	}

	contentType := fh.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, err
	}

	switch mediaType {
	case "application/json":
		return s.importJson(ctx, name, pbhType, userId, f)
	case "text/calendar":
		return s.importICS(ctx, name, pbhType, userId, f)
	default:
		return nil, common.NewValidationError("file", "File is not supported.")
	}
}

func (s *store) importJson(
	ctx context.Context,
	name, pbhType, userId string,
	r io.Reader,
) (*Response, error) {
	dates := make(map[string]string)
	d := json.NewDecoder(r)
	err := d.Decode(&dates)
	if err != nil {
		return nil, common.NewValidationError("file", "File is not supported.")
	}

	location := s.timezoneConfigProvider.Get().Location
	exdates := make([]pbehavior.Exdate, 0, len(dates))
	now := datetime.NewCpsTime()
	for dateStr := range dates {
		start, err := time.ParseInLocation(time.DateOnly, dateStr, location)
		if err != nil {
			return nil, common.NewValidationError("file", "File is not supported.")
		}
		end := start.AddDate(0, 0, 1)
		if end.Before(now.Time) {
			continue
		}

		exdates = append(exdates, pbehavior.Exdate{
			Exdate: types.Exdate{
				Begin: datetime.CpsTime{Time: start},
				End:   datetime.CpsTime{Time: end},
			},
			Type: pbhType,
		})
	}

	if len(exdates) == 0 {
		return nil, common.NewValidationError("file", "File is empty.")
	}

	doc := pbehavior.Exception{
		ID:      utils.NewID(),
		Name:    name,
		Author:  userId,
		Exdates: exdates,
		Created: datetime.NewCpsTime(),
	}

	var response *Response
	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil

		_, err := s.dbCollection.InsertOne(ctx, doc)
		if err != nil {
			return err
		}

		response, err = s.GetById(ctx, doc.ID)
		return err
	})

	return response, err
}

func (s *store) importICS(
	ctx context.Context,
	name, pbhType, userId string,
	r io.Reader,
) (*Response, error) {
	cal := ics.NewParser(r)
	now := datetime.NewCpsTime()
	intervalStart := now.Time
	intervalEnd := intervalStart.AddDate(maxImportDateInYears, 0, 0)
	cal.Start = &intervalStart
	cal.End = &intervalEnd
	err := cal.Parse()
	if err != nil {
		return nil, common.NewValidationError("file", "File is not supported.")
	}

	exdates := make([]pbehavior.Exdate, 0, len(cal.Events))
	location := s.timezoneConfigProvider.Get().Location
	for _, event := range cal.Events {
		if event.Start == nil || event.End == nil {
			return nil, common.NewValidationError("file", "File is not valid.")
		}

		start := adjustCalendarTime(*event.Start, location)
		end := adjustCalendarTime(*event.End, location)
		if end.Before(now.Time) {
			continue
		}

		exdates = append(exdates, pbehavior.Exdate{
			Exdate: types.Exdate{
				Begin: datetime.CpsTime{Time: start},
				End:   datetime.CpsTime{Time: end},
			},
			Type: pbhType,
		})
	}

	if len(exdates) == 0 {
		return nil, common.NewValidationError("file", "File is empty.")
	}

	doc := pbehavior.Exception{
		ID:      utils.NewID(),
		Name:    name,
		Author:  userId,
		Exdates: exdates,
		Created: now,
	}

	var response *Response
	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil

		_, err := s.dbCollection.InsertOne(ctx, doc)
		if err != nil {
			return err
		}

		response, err = s.GetById(ctx, doc.ID)
		return err
	})

	return response, err
}

func getNestedObjectsPipeline() []bson.M {
	return []bson.M{
		// Lookup exdate type
		{"$unwind": "$exdates"},
		{"$lookup": bson.M{
			"from":         mongo.PbehaviorTypeMongoCollection,
			"localField":   "exdates.type",
			"foreignField": "_id",
			"as":           "exdates.type",
		}},
		{"$unwind": "$exdates.type"},
		{"$sort": bson.M{"exdates.begin": 1}},
		{"$group": bson.M{
			"_id":         "$_id",
			"name":        bson.M{"$first": "$name"},
			"description": bson.M{"$first": "$description"},
			"author":      bson.M{"$first": "$author"},
			"created":     bson.M{"$first": "$created"},
			"updated":     bson.M{"$first": "$updated"},
			"deletable":   bson.M{"$first": "$deletable"},
			"exdates":     bson.M{"$push": "$exdates"},
		}},
	}
}

func getDeletablePipeline() []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from": mongo.PbehaviorMongoCollection,
			"let":  bson.M{"exception": "$_id"},
			"pipeline": []bson.M{
				{"$match": bson.M{"$expr": bson.M{"$and": []bson.M{
					{"$isArray": "$exceptions"},
					{"$in": bson.A{"$$exception", "$exceptions"}},
				}}}},
				{"$limit": 1},
				{"$project": bson.M{"_id": 1}},
			},
			"as": "pbhs",
		}},
		{"$lookup": bson.M{
			"from": mongo.EventFilterRuleCollection,
			"let":  bson.M{"exception": "$_id"},
			"pipeline": []bson.M{
				{"$match": bson.M{"$expr": bson.M{"$and": []bson.M{
					{"$isArray": "$exceptions"},
					{"$in": bson.A{"$$exception", "$exceptions"}},
				}}}},
				{"$limit": 1},
				{"$project": bson.M{"_id": 1}},
			},
			"as": "efs",
		}},
		{"$addFields": bson.M{
			"deletable": bson.M{"$and": []bson.M{
				{"$eq": bson.A{bson.M{"$size": "$pbhs"}, 0}},
				{"$eq": bson.A{bson.M{"$size": "$efs"}, 0}},
			}},
		}},
		{"$project": bson.M{
			"pbhs": 0,
			"efs":  0,
		}},
	}
}

func adjustCalendarTime(t time.Time, location *time.Location) time.Time {
	if t.Location() != time.UTC {
		return t
	}

	y, m, d := t.Date()
	return time.Date(y, m, d, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), location)
}

func (s *store) transformRequestToDoc(ctx context.Context, r EditRequest) (*pbehavior.Exception, error) {
	var err error
	var exception pbehavior.Exception

	exception.Name = r.Name
	exception.Description = r.Description
	exception.Author = r.Author
	exception.Exdates, err = s.transformExdatesRequestToModel(ctx, r.Exdates)
	if err != nil {
		return nil, err
	}

	return &exception, nil
}

func (s *store) transformExdatesRequestToModel(ctx context.Context, r []ExdateRequest) ([]pbehavior.Exdate, error) {
	if len(r) == 0 {
		return []pbehavior.Exdate{}, nil
	}

	pbhTypes := make([]string, len(r))
	for i := range r {
		pbhTypes[i] = r[i].Type
	}

	res, err := s.typeDbCollection.Find(ctx, bson.M{"_id": bson.M{"$in": pbhTypes}})
	if err != nil {
		return nil, err
	}

	defer res.Close(ctx)
	typesByID := make(map[string]*pbehavior.Type)
	for res.Next(ctx) {
		var t pbehavior.Type
		err = res.Decode(&t)
		if err != nil {
			return nil, err
		}
		if _, ok := typesByID[t.ID]; !ok {
			typesByID[t.ID] = &t
		}
	}

	exdates := make([]pbehavior.Exdate, len(r))
	for i := range r {
		if t, ok := typesByID[r[i].Type]; ok {
			exdates[i] = pbehavior.Exdate{
				Exdate: types.Exdate{
					Begin: r[i].Begin,
					End:   r[i].End,
				},
				Type: t.ID,
			}
		} else {
			return nil, ErrTypeNotExists
		}
	}

	return exdates, nil
}
