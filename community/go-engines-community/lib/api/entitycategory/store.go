package entitycategory

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Store interface {
	Find(r ListRequest) (*AggregationResult, error)
	GetOneBy(id string) (*Category, error)
	Insert(r EditRequest) (*Category, error)
	Update(r EditRequest) (*Category, error)
	Delete(id string) (bool, error)
}

func NewStore(dbClient mongo.DbClient) Store {
	return &store{
		dbClient:      dbClient,
		dbCollection:  dbClient.Collection(mongo.EntityCategoryMongoCollection),
		defaultSortBy: "name",
	}
}

type store struct {
	dbClient      mongo.DbClient
	dbCollection  mongo.DbCollection
	defaultSortBy string
}

func (s *store) Find(r ListRequest) (*AggregationResult, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	filter := bson.M{}

	if r.Search != "" {
		searchRegexp := primitive.Regex{
			Pattern: fmt.Sprintf(".*%s.*", r.Search),
			Options: "i",
		}

		filter["$or"] = []bson.M{
			{"name": searchRegexp},
		}
	}

	sortBy := r.SortBy
	if sortBy == "" {
		sortBy = s.defaultSortBy
	}

	pipeline := []bson.M{{"$match": filter}}
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

func (s *store) GetOneBy(id string) (*Category, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

func (s *store) Insert(r EditRequest) (*Category, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

func (s *store) Update(r EditRequest) (*Category, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	now := types.CpsTime{Time: time.Now()}
	res, err := s.dbCollection.UpdateOne(ctx,
		bson.M{"_id": r.ID},
		bson.M{"$set": bson.M{
			"name":    r.Name,
			"author":  r.Author,
			"updated": now,
		}},
	)
	if err != nil {
		return nil, err
	}

	if res.MatchedCount == 0 {
		return nil, nil
	}

	return s.GetOneBy(r.ID)
}

func (s *store) Delete(id string) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
