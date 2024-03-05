package entitybasic

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

type basicValidator struct {
	dbClient mongo.DbClient
}

func NewValidator(client mongo.DbClient) Validator {
	return &basicValidator{dbClient: client}
}

func (v *basicValidator) ValidateEditRequest(ctx context.Context, sl validator.StructLevel) {
	r := sl.Current().Interface().(EditRequest)

	// Validate category
	if r.Category != "" {
		err := v.dbClient.Collection(mongo.EntityCategoryMongoCollection).
			FindOne(ctx, bson.M{"_id": r.Category}).Err()
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				sl.ReportError(r.Category, "Category", "Category", "not_exist", "")
			} else {
				panic(err)
			}
		}
	}
}
