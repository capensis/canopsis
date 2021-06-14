package playlist

import (
	"context"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
)

type Validator interface {
	ValidateEditRequest(ctx context.Context, sl validator.StructLevel)
}

type baseValidator struct {
	client         mongo.DbClient
	viewCollection mongo.DbCollection
}

func NewPlaylistValidator(dbClient mongo.DbClient) Validator {
	return &baseValidator{
		client:         dbClient,
		viewCollection: dbClient.Collection(mongo.ViewMongoCollection),
	}
}

func (v *baseValidator) ValidateEditRequest(ctx context.Context, sl validator.StructLevel) {
	r := sl.Current().Interface().(EditRequest)

	if len(r.TabsList) > 0 {
		cursor, err := v.viewCollection.Aggregate(ctx, []bson.M{
			{"$unwind": "$tabs"},
			{"$match": bson.M{"tabs._id": bson.M{"$in": r.TabsList}}},
			{"$group": bson.M{"_id": nil, "count": bson.M{"$sum": 1}}},
		})
		if err != nil {
			panic(err)
		}

		res := struct {
			Count int64 `bson:"count"`
		}{}
		if cursor.Next(ctx) {
			err := cursor.Decode(&res)
			if err != nil {
				panic(err)
			}
		}

		if res.Count < int64(len(r.TabsList)) {
			sl.ReportError(r.TabsList, "TabsList", "TabsList", "not_exist", "")
		}
	}
}
