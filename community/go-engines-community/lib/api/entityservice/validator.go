package entityservice

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/oldpattern"
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

	v.validateCategory(ctx, sl, r.Category)
	v.validatePatterns(sl, r.EntityPatterns)
}

func (v *basicValidator) validateCategory(ctx context.Context, sl validator.StructLevel, category string) {
	if category != "" {
		err := v.dbClient.Collection(mongo.EntityCategoryMongoCollection).
			FindOne(ctx, bson.M{"_id": category}).Err()
		if err != nil {
			if err == mongodriver.ErrNoDocuments {
				sl.ReportError(category, "Category", "Category", "not_exist", "")
			} else {
				panic(err)
			}
		}
	}
}

func (v *basicValidator) validatePatterns(sl validator.StructLevel, patterns oldpattern.EntityPatternList) {
	if !patterns.IsValid() {
		sl.ReportError(patterns, "EntityPatterns", "EntityPatterns", "entitypattern_invalid", "")
	} else {
		if patterns.IsSet() {
			query := patterns.AsMongoDriverQuery()["$or"].([]bson.M)
			if len(query) == 0 {
				sl.ReportError(patterns, "EntityPatterns", "EntityPatterns", "entitypattern_empty", "")
			} else {
				for _, q := range query {
					if len(q) == 0 {
						sl.ReportError(patterns, "EntityPatterns", "EntityPatterns", "entitypattern_contains_empty", "")
					}
				}
			}
		}
	}
}
