package statesettings

import (
	"context"

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
	GetById(ctx context.Context, id string) (*Response, error)
	Find(ctx context.Context, query FilteredQuery) (*AggregationResult, error)
	Insert(ctx context.Context, r EditRequest) (*Response, error)
	Update(ctx context.Context, r EditRequest) (*Response, error)
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
		defaultSearchByFields: []string{"_id", "title"},
	}
}

func (s *store) GetById(ctx context.Context, id string) (*Response, error) {
	pipeline := []bson.M{{"$match": bson.M{"_id": id}}, addEditableAndDeletableFields()}

	cursor, err := s.dbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var res Response

		err = cursor.Decode(&res)
		if err != nil {
			return nil, err
		}

		return &res, nil
	}

	return nil, nil
}

func (s *store) Find(ctx context.Context, query FilteredQuery) (*AggregationResult, error) {
	var pipeline []bson.M

	filter := common.GetSearchQuery(query.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	pipeline = append(pipeline, addEditableAndDeletableFields())

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

func (s *store) Insert(ctx context.Context, r EditRequest) (*Response, error) {
	r.ID = utils.NewID()

	var response *Response

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

func (s *store) Update(ctx context.Context, r EditRequest) (*Response, error) {
	var response *Response

	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil

		unset := make(bson.M)
		if r.Method == statesetting.MethodDependencies {
			unset["inherited_entity_pattern"] = 1
		} else if r.Method == statesetting.MethodInherited {
			unset["state_thresholds"] = 1
		}

		res, err := s.dbCollection.UpdateOne(ctx, bson.M{"_id": r.ID}, bson.M{"$set": r, "$unset": unset})
		if err != nil {
			if mongodriver.IsDuplicateKeyError(err) {
				return common.NewValidationError("title", "Title already exists.")
			}

			return err
		}

		if res.MatchedCount == 0 {
			return nil
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

func addEditableAndDeletableFields() bson.M {
	return bson.M{
		"$addFields": bson.M{
			"editable": bson.M{"$cond": bson.M{
				"if":   bson.M{"$eq": bson.A{"$_id", statesetting.ServiceID}},
				"then": false,
				"else": true,
			}},
			"deletable": bson.M{"$cond": bson.M{
				"if":   bson.M{"$in": bson.A{"$_id", bson.A{statesetting.ServiceID, statesetting.JUnitID}}},
				"then": false,
				"else": true,
			}},
		},
	}
}
