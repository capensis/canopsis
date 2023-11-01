package statesettings

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/priority"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statesetting"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

const StickySortField = "on_top"

type Store interface {
	GetById(ctx context.Context, id string) (*StateSetting, error)
	Find(ctx context.Context, query FilteredQuery) (*AggregationResult, error)
	Insert(ctx context.Context, r StateSetting) (*StateSetting, error)
	Update(ctx context.Context, r StateSetting) (*StateSetting, error)
	Delete(ctx context.Context, id string) (bool, error)
}

type store struct {
	dbClient     mongo.DbClient
	dbCollection mongo.DbCollection

	defaultSearchByFields []string
}

func NewStore(
	dbClient mongo.DbClient,
) Store {
	return &store{
		dbClient:              dbClient,
		dbCollection:          dbClient.Collection(mongo.StateSettingsMongoCollection),
		defaultSearchByFields: []string{"_id", "name"},
	}
}

func (s *store) GetById(ctx context.Context, id string) (*StateSetting, error) {
	var res StateSetting

	err := s.dbCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&res)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	return &res, nil
}

func (s *store) Find(ctx context.Context, query FilteredQuery) (*AggregationResult, error) {
	var pipeline []bson.M

	filter := common.GetSearchQuery(query.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	sortBy := "title"
	if query.SortBy != "" {
		sortBy = query.SortBy
	}

	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		query.Query,
		pipeline,
		s.getSortQuery(sortBy, query.Sort),
	))
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

	return &result, nil
}

func (s *store) Insert(ctx context.Context, r StateSetting) (*StateSetting, error) {
	r.ID = utils.NewID()

	var response *StateSetting

	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil

		_, err := s.dbCollection.InsertOne(ctx, r)
		if err != nil {
			if mongodriver.IsDuplicateKeyError(err) {
				return common.NewValidationError("title", "Title already exists.")
			}

			return err
		}

		err = priority.UpdateFollowing(ctx, s.dbCollection, r.ID, r.Priority)
		if err != nil {
			return err
		}

		response, err = s.GetById(ctx, r.ID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *store) Update(ctx context.Context, r StateSetting) (*StateSetting, error) {
	var response *StateSetting

	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil

		unset := make(bson.M)
		if r.Method == statesetting.MethodDependencies {
			unset["inherited_entity_pattern"] = 1
		}

		if r.Method == statesetting.MethodInherited {
			unset["state_thresholds"] = 1
		}

		res, err := s.dbCollection.UpdateOne(ctx, bson.M{"_id": r.ID}, bson.M{"$set": r, "$unset": unset})
		if err != nil {
			if mongodriver.IsDuplicateKeyError(err) {
				return common.NewValidationError("title", "Title already exists.")
			}

			return err
		}

		if err != nil || res.MatchedCount == 0 {
			return err
		}

		err = priority.UpdateFollowing(ctx, s.dbCollection, r.ID, r.Priority)
		if err != nil {
			return err
		}

		response, err = s.GetById(ctx, r.ID)
		return err
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	if id == statesetting.JUnitID || id == statesetting.ServiceID {
		return false, ErrDefaultRule
	}

	deleted, err := s.dbCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}

	return deleted > 0, nil
}

func (s *store) getSortQuery(sortBy, sort string) bson.M {
	sortDir := 1
	if sort == common.SortDesc {
		sortDir = -1
	}

	q := bson.D{{Key: StickySortField, Value: -1}, {Key: sortBy, Value: sortDir}}
	if sortBy != "_id" {
		q = append(q, bson.E{Key: "_id", Value: 1})
	}

	return bson.M{"$sort": q}
}
