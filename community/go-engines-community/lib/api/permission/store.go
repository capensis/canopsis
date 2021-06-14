package permission

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/go-engines/lib/api/common"
	"git.canopsis.net/canopsis/go-engines/lib/api/pagination"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	securitymodel "git.canopsis.net/canopsis/go-engines/lib/security/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Store interface {
	Find(request ListRequest) (*AggregationResult, error)
}

func NewStore(dbClient mongo.DbClient) Store {
	return &store{
		dbClient:     dbClient,
		dbCollection: dbClient.Collection(mongo.RightsMongoCollection),
	}
}

type store struct {
	dbClient     mongo.DbClient
	dbCollection mongo.DbCollection
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
			{"crecord_name": searchRegexp},
			{"desc": searchRegexp},
		}
	}

	sortBy := "name"
	if r.SortBy != "" {
		sortBy = r.SortBy
	}

	pipeline := []bson.M{
		{"$match": bson.M{"crecord_type": securitymodel.LineTypeObject}},
		{"$match": filter},
		{"$addFields": bson.M{
			"name":        "$crecord_name",
			"description": "$desc",
		}},
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
