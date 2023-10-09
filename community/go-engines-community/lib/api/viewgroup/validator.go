package viewgroup

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type Validator interface {
	ValidateEditRequest(ctx context.Context, sl validator.StructLevel)
}

func NewValidator(dbClient mongo.DbClient) Validator {
	return &baseValidator{
		collection: dbClient.Collection(mongo.ViewGroupMongoCollection),
	}
}

type baseValidator struct {
	collection mongo.DbCollection
}

func (v *baseValidator) ValidateEditRequest(ctx context.Context, sl validator.StructLevel) {
	r := sl.Current().Interface().(EditRequest)

	err := v.collection.FindOne(ctx, bson.M{"title": r.Title, "is_private": false}).Err()
	if err == nil {
		sl.ReportError(r.Title, "Title", "Title", "unique", "")
	} else if !errors.Is(err, mongodriver.ErrNoDocuments) {
		panic(err)
	}
}
