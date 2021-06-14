package pbehaviorexception

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/go-engines/lib/api/common"
	"git.canopsis.net/canopsis/go-engines/lib/api/pagination"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"git.canopsis.net/canopsis/go-engines/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Store interface {
	Insert(model *Exception) error
	Find(r ListRequest) (*AggregationResult, error)
	GetOneBy(filter bson.M) (*Exception, error)
	Update(model *Exception) (bool, error)
	Delete(id string) (bool, error)
	IsLinked(id string) (bool, error)
}

func NewStore(dbClient mongo.DbClient) Store {
	return &store{
		dbClient:            dbClient,
		dbCollection:        dbClient.Collection(pbehavior.ExceptionCollectionName),
		pbehaviorCollection: dbClient.Collection(pbehavior.PBehaviorCollectionName),
		defaultSortBy:       "created",
	}
}

type store struct {
	dbClient            mongo.DbClient
	dbCollection        mongo.DbCollection
	pbehaviorCollection mongo.DbCollection
	defaultSortBy       string
}

func (s *store) Insert(model *Exception) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

func (s *store) Find(r ListRequest) (*AggregationResult, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var filter bson.M
	if r.Search != "" {
		searchRegexp := primitive.Regex{
			Pattern: fmt.Sprintf(".*%s.*", r.Search),
			Options: "i",
		}

		filter = bson.M{"$or": []bson.M{
			{"name": searchRegexp},
			{"description": searchRegexp},
		}}
	} else {
		filter = bson.M{}
	}

	sortBy := r.SortBy
	if sortBy == "" {
		sortBy = s.defaultSortBy
	}

	sort := common.GetSortQuery(sortBy, r.Sort)
	pipeline := []bson.M{{"$match": filter}}
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

func (s *store) GetOneBy(filter bson.M) (*Exception, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pipeline := []bson.M{
		{"$match": filter},
	}
	pipeline = append(pipeline, getNestedObjectsPipeline()...)
	cursor, err := s.dbCollection.Aggregate(ctx, pipeline)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var exception Exception
		err = cursor.Decode(&exception)
		if err != nil {
			return nil, err
		}

		return &exception, err
	}

	return nil, nil
}

func (s *store) Update(model *Exception) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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

func (s *store) Delete(id string) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	isLinked, err := s.IsLinked(id)
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
func (s *store) IsLinked(id string) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	res := s.pbehaviorCollection.FindOne(ctx, bson.M{"exceptions": id})
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
			"from":         pbehavior.TypeCollectionName,
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
			"from":         pbehavior.PBehaviorCollectionName,
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
	}
}
