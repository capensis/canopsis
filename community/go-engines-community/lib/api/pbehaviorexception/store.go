package pbehaviorexception

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
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
	Insert(ctx context.Context, model *Exception) error
	Find(ctx context.Context, r ListRequest) (*AggregationResult, error)
	GetOneById(ctx context.Context, id string) (*Exception, error)
	Update(ctx context.Context, model *Exception) (bool, error)
	Delete(ctx context.Context, id string) (bool, error)
	IsLinked(ctx context.Context, id string) (bool, error)
	Import(ctx context.Context, name, pbhType string, file *multipart.FileHeader) (*Exception, error)
}

func NewStore(dbClient mongo.DbClient, timezoneConfigProvider config.TimezoneConfigProvider) Store {
	return &store{
		dbClient:              dbClient,
		dbCollection:          dbClient.Collection(mongo.PbehaviorExceptionMongoCollection),
		pbehaviorDbCollection: dbClient.Collection(mongo.PbehaviorMongoCollection),
		typeDbCollection:      dbClient.Collection(mongo.PbehaviorTypeMongoCollection),
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
}

func (s *store) Insert(ctx context.Context, model *Exception) error {
	if model.ID == "" {
		model.ID = utils.NewID()
	}

	created := types.CpsTime{Time: time.Now()}
	exdates := make([]pbehavior.Exdate, len(model.Exdates))
	for i := range model.Exdates {
		exdates[i].Type = model.Exdates[i].Type.ID
		exdates[i].Begin = model.Exdates[i].Begin
		exdates[i].End = model.Exdates[i].End
	}

	_, err := s.dbCollection.InsertOne(ctx, pbehavior.Exception{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description,
		Exdates:     exdates,
		Created:     &created,
	})

	if err != nil {
		return err
	}

	model.Created = created

	return nil
}

func (s *store) Find(ctx context.Context, r ListRequest) (*AggregationResult, error) {
	pipeline := make([]bson.M, 0)
	filter := common.GetSearchQuery(r.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	sortBy := r.SortBy
	if sortBy == "" {
		sortBy = s.defaultSortBy
	}

	sort := common.GetSortQuery(sortBy, r.Sort)
	project := getNestedObjectsPipeline()
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

func (s *store) GetOneById(ctx context.Context, id string) (*Exception, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"_id": id}},
	}
	pipeline = append(pipeline, getNestedObjectsPipeline()...)
	cursor, err := s.dbCollection.Aggregate(ctx, pipeline)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	if cursor.Next(ctx) {
		var exception Exception
		err = cursor.Decode(&exception)
		if err != nil {
			return nil, err
		}

		return &exception, err
	}

	return nil, nil
}

func (s *store) Update(ctx context.Context, model *Exception) (bool, error) {
	exdates := make([]pbehavior.Exdate, len(model.Exdates))
	for i := range model.Exdates {
		exdates[i].Type = model.Exdates[i].Type.ID
		exdates[i].Begin = model.Exdates[i].Begin
		exdates[i].End = model.Exdates[i].End
	}

	res := s.dbCollection.FindOneAndUpdate(ctx, bson.M{"_id": model.ID}, bson.M{"$set": pbehavior.Exception{
		Name:        model.Name,
		Description: model.Description,
		Exdates:     exdates,
	}})

	if err := res.Err(); err != nil {
		if err == mongodriver.ErrNoDocuments {
			return false, nil
		}

		return false, err
	}

	var v struct{ Created *types.CpsTime }
	err := res.Decode(&v)
	if err != nil {
		return false, err
	}

	model.Created = *v.Created

	return true, nil
}

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	isLinked, err := s.IsLinked(ctx, id)
	if err != nil {
		return false, err
	}

	if isLinked {
		return false, ErrLinkedException
	}

	deleted, err := s.dbCollection.DeleteOne(ctx, bson.M{"_id": id})

	return deleted > 0, err
}

// IsLinked checks if there is pbehavior with linked exception.
func (s *store) IsLinked(ctx context.Context, id string) (bool, error) {
	res := s.pbehaviorDbCollection.FindOne(ctx, bson.M{"exceptions": id})
	if err := res.Err(); err != nil {
		if err == mongodriver.ErrNoDocuments {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (s *store) Import(ctx context.Context, name, pbhType string, file *multipart.FileHeader) (*Exception, error) {
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

	h, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer h.Close()

	content, err := io.ReadAll(h)
	if err != nil {
		return nil, err
	}

	contentType := ""
	if v, ok := file.Header["Content-Type"]; ok && len(v) > 0 {
		contentType = v[0]
	}

	switch contentType {
	case "application/json":
		return s.importJson(ctx, name, pbhType, content)
	case "text/calendar":
		return s.importICS(ctx, name, pbhType, content)
	default:
		return nil, common.NewValidationError("file", "File is not supported.")
	}
}

func (s *store) importJson(
	ctx context.Context,
	name, pbhType string,
	input []byte,
) (*Exception, error) {
	dates := make(map[string]string)
	err := json.Unmarshal(input, &dates)
	if err != nil {
		return nil, common.NewValidationError("file", "File is not supported.")
	}

	location := s.timezoneConfigProvider.Get().Location
	exdates := make([]pbehavior.Exdate, len(dates))
	i := 0
	for dateStr := range dates {
		date, err := time.ParseInLocation("2006-01-02", dateStr, location)
		if err != nil {
			return nil, common.NewValidationError("file", "File is not supported.")
		}

		exdates[i] = pbehavior.Exdate{
			Exdate: types.Exdate{
				Begin: types.CpsTime{Time: date},
				End:   types.CpsTime{Time: date.AddDate(0, 0, 1)},
			},
			Type: pbhType,
		}

		i++
	}

	now := types.NewCpsTime()
	doc := pbehavior.Exception{
		ID:      utils.NewID(),
		Name:    name,
		Exdates: exdates,
		Created: &now,
	}

	var response *Exception
	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil

		_, err := s.dbCollection.InsertOne(ctx, doc)
		if err != nil {
			return err
		}

		response, err = s.GetOneById(ctx, doc.ID)
		return err
	})

	return response, err
}

func (s *store) importICS(
	ctx context.Context,
	name, pbhType string,
	input []byte,
) (*Exception, error) {
	cal := ics.NewParser(bytes.NewReader(input))
	intervalStart := time.Now()
	intervalEnd := intervalStart.AddDate(maxImportDateInYears, 0, 0)
	cal.Start = &intervalStart
	cal.End = &intervalEnd
	err := cal.Parse()
	if err != nil {
		return nil, common.NewValidationError("file", "File is not supported.")
	}

	exdates := make([]pbehavior.Exdate, len(cal.Events))
	location := s.timezoneConfigProvider.Get().Location
	for i, event := range cal.Events {
		if event.Start == nil || event.End == nil {
			return nil, common.NewValidationError("file", "File is not valid.")
		}

		start := adjustCalendarTime(*event.Start, location)
		end := adjustCalendarTime(*event.End, location)

		exdates[i] = pbehavior.Exdate{
			Exdate: types.Exdate{
				Begin: types.CpsTime{Time: start},
				End:   types.CpsTime{Time: end},
			},
			Type: pbhType,
		}
	}

	now := types.NewCpsTime()
	doc := pbehavior.Exception{
		ID:      utils.NewID(),
		Name:    name,
		Exdates: exdates,
		Created: &now,
	}

	var response *Exception
	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil

		_, err := s.dbCollection.InsertOne(ctx, doc)
		if err != nil {
			return err
		}

		response, err = s.GetOneById(ctx, doc.ID)
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
		{"$group": bson.M{
			"_id":         "$_id",
			"name":        bson.M{"$first": "$name"},
			"description": bson.M{"$first": "$description"},
			"created":     bson.M{"$first": "$created"},
			"deletable":   bson.M{"$first": "$deletable"},
			"exdates":     bson.M{"$push": "$exdates"},
		}},
	}
}

func getDeletablePipeline() []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.PbehaviorMongoCollection,
			"localField":   "_id",
			"foreignField": "exceptions",
			"as":           "pbh",
		}},
		{"$addFields": bson.M{
			"deletable": bson.M{"$eq": bson.A{bson.M{"$size": "$pbh"}, 0}},
		}},
		{"$project": bson.M{
			"pbh": 0,
		}},
		{"$lookup": bson.M{
			"from":         mongo.EventFilterRulesMongoCollection,
			"localField":   "_id",
			"foreignField": "exceptions",
			"as":           "ef",
		}},
		{"$addFields": bson.M{
			"deletable": bson.M{"$and": bson.A{bson.M{"$eq": bson.A{bson.M{"$size": "$ef"}, 0}}, "$deletable"}},
		}},
		{"$project": bson.M{
			"ef": 0,
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
