package template

import (
	"cmp"
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
)

type Store interface {
	FindData(ctx context.Context, r ListDataRequest) (AggregationDataResult, error)
	GetData(ctx context.Context, id string) (DataResponse, error)
	CreateData(ctx context.Context, r EditDataRequest) (DataResponse, error)
	UpdateData(ctx context.Context, r EditDataRequest) (DataResponse, error)
	DeleteData(ctx context.Context, id, author string) (bool, error)
}

func NewStore(client mongo.DbClient) Store {
	return &store{
		client:        client,
		collection:    client.Collection(mongo.TemplateDataCollection),
		defaultSortBy: "name",
		defaultSearchByFields: []string{
			"name",
			"description",
		},
	}
}

type store struct {
	client                mongo.DbClient
	collection            mongo.DbCollection
	defaultSortBy         string
	defaultSearchByFields []string
}

func (s *store) FindData(ctx context.Context, r ListDataRequest) (AggregationDataResult, error) {
	var res AggregationDataResult
	pipeline := make([]bson.M, 0)
	if r.Type != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"type": r.Type}})
	}

	filter := common.GetSearchQuery(r.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	cursor, err := s.collection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		pipeline,
		common.GetSortQuery(cmp.Or(r.SortBy, s.defaultSortBy), r.Sort),
	))
	if err != nil {
		return res, err
	}

	defer cursor.Close(ctx)
	if cursor.Next(ctx) {
		err = cursor.Decode(&res)
		if err != nil {
			return res, err
		}
	}

	if err = cursor.Err(); err != nil {
		return res, err
	}

	return res, nil
}

func (s *store) GetData(ctx context.Context, id string) (DataResponse, error) {
	res := DataResponse{}
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&res)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return res, nil
		}

		return res, err
	}

	return res, nil
}

func (s *store) CreateData(ctx context.Context, r EditDataRequest) (DataResponse, error) {
	now := datetime.NewCpsTime()
	model := DataModel{
		ID:          utils.NewID(),
		Type:        *r.Type,
		Name:        r.Name,
		Description: r.Description,
		Body:        r.Body,
		Headers:     r.Headers,
		Author:      r.Author,
		Created:     &now,
		Updated:     &now,
	}
	var res DataResponse
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		res = DataResponse{}

		_, err := s.collection.InsertOne(ctx, model)
		if err != nil {
			if mongodriver.IsDuplicateKeyError(err) {
				return common.NewValidationError("name", "Name already exists.")
			}

			return err
		}

		res, err = s.GetData(ctx, model.ID)

		return err
	})

	return res, err
}

func (s *store) UpdateData(ctx context.Context, r EditDataRequest) (DataResponse, error) {
	now := datetime.NewCpsTime()
	model := DataModel{
		Type:        *r.Type,
		Name:        r.Name,
		Description: r.Description,
		Body:        r.Body,
		Headers:     r.Headers,
		Author:      r.Author,
		Updated:     &now,
	}
	var res DataResponse
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		res = DataResponse{}

		prev, err := s.GetData(ctx, r.ID)
		if err != nil || prev.ID == "" {
			return err
		}

		if prev.Type != model.Type {
			return common.NewValidationError("type", "Type cannot be changed.")
		}

		_, err = s.collection.UpdateOne(ctx, bson.M{"_id": r.ID}, bson.M{"$set": model})
		if err != nil {
			if mongodriver.IsDuplicateKeyError(err) {
				return common.NewValidationError("name", "Name already exists.")
			}

			return err
		}

		res, err = s.GetData(ctx, r.ID)

		return err
	})

	return res, err
}

func (s *store) DeleteData(ctx context.Context, id, author string) (bool, error) {
	var res bool
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		res = false
		// required to get the author in action log listener
		ur, err := s.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"author": author}})
		if err != nil || ur.MatchedCount == 0 {
			return err
		}

		d, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
		if err != nil || d == 0 {
			return err
		}

		res = true

		return nil
	})

	return res, err
}
