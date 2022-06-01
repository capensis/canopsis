package entitycategory

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type Store interface {
	Find(ctx context.Context, r ListRequest) (*AggregationResult, error)
	GetOneBy(ctx context.Context, id string) (*Category, error)
	Insert(ctx context.Context, r EditRequest) (*Category, error)
	Update(ctx context.Context, r EditRequest) (*Category, error)
	Delete(ctx context.Context, id string) (bool, error)
}

func NewStore(dbClient mongo.DbClient) Store {
	return &store{
		dbClient:              dbClient,
		dbCollection:          dbClient.Collection(mongo.EntityCategoryMongoCollection),
		defaultSearchByFields: []string{"_id", "name", "author"},
		defaultSortBy:         "name",
	}
}

type store struct {
	dbClient              mongo.DbClient
	dbCollection          mongo.DbCollection
	defaultSearchByFields []string
	defaultSortBy         string
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

	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		pipeline,
		common.GetSortQuery(sortBy, r.Sort),
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

func (s *store) GetOneBy(ctx context.Context, id string) (*Category, error) {
	res := s.dbCollection.FindOne(ctx, bson.M{"_id": id})
	if err := res.Err(); err != nil {
		if err == mongodriver.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	category := &Category{}
	err := res.Decode(category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *store) Insert(ctx context.Context, r EditRequest) (*Category, error) {
	now := types.CpsTime{Time: time.Now()}
	category := Category{
		ID:      utils.NewID(),
		Name:    r.Name,
		Author:  r.Author,
		Created: &now,
		Updated: &now,
	}
	_, err := s.dbCollection.InsertOne(ctx, category)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (s *store) Update(ctx context.Context, r EditRequest) (*Category, error) {
	now := types.CpsTime{Time: time.Now()}
	var result *Category

	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		result = nil
		res, err := s.dbCollection.UpdateOne(ctx,
			bson.M{"_id": r.ID},
			bson.M{"$set": bson.M{
				"name":    r.Name,
				"author":  r.Author,
				"updated": now,
			}},
		)
		if err != nil || res.MatchedCount == 0 {
			return err
		}

		result, err = s.GetOneBy(ctx, r.ID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	entityCollection := s.dbClient.Collection(mongo.EntityMongoCollection)
	res := entityCollection.FindOne(ctx, bson.M{"category": id})
	if err := res.Err(); err != nil {
		if err != mongodriver.ErrNoDocuments {
			return false, err
		}
	} else {
		return false, ErrLinkedCategoryToEntity
	}

	delCount, err := s.dbCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}

	return delCount > 0, nil
}
