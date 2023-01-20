package pbehaviorexception

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface {
	Insert(ctx context.Context, model *Exception) error
	Find(ctx context.Context, r ListRequest) (*AggregationResult, error)
	GetOneBy(ctx context.Context, filter bson.M) (*Exception, error)
	Update(ctx context.Context, model *Exception) (bool, error)
	Delete(ctx context.Context, id string) (bool, error)
	IsLinked(ctx context.Context, id string) (bool, error)
}

func NewStore(dbClient mongo.DbClient) Store {
	return &store{
		dbCollection:          dbClient.Collection(mongo.PbehaviorExceptionMongoCollection),
		pbehaviorDbCollection: dbClient.Collection(mongo.PbehaviorMongoCollection),
		defaultSearchByFields: []string{"_id", "name", "description"},
		defaultSortBy:         "created",
	}
}

type store struct {
	dbCollection          mongo.DbCollection
	pbehaviorDbCollection mongo.DbCollection
	defaultSearchByFields []string
	defaultSortBy         string
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

func (s *store) GetOneBy(ctx context.Context, filter bson.M) (*Exception, error) {
	pipeline := []bson.M{
		{"$match": filter},
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
