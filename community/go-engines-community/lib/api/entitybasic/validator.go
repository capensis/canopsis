package entitybasic

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
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
	validTypes := []string{types.EntityTypeResource, types.EntityTypeComponent, types.EntityTypeConnector}

	// Validate category
	if r.Category != "" {
		err := v.dbClient.Collection(mongo.EntityCategoryMongoCollection).
			FindOne(ctx, bson.M{"_id": r.Category}).Err()
		if err != nil {
			if err == mongodriver.ErrNoDocuments {
				sl.ReportError(r.Category, "Category", "Category", "not_exist", "")
			} else {
				panic(err)
			}
		}
	}

	// Validate impacts
	if len(r.Impacts) > 0 {
		if v.hasDuplicates(r.Impacts) {
			sl.ReportError(r.Impacts, "Impacts", "Impacts", "has_duplicates", "")
		} else {
			ok, err := v.checkExist(ctx, r.Impacts, validTypes)
			if err != nil {
				panic(err)
			}

			if !ok {
				sl.ReportError(r.Impacts, "Impacts", "Impacts", "not_exist", "")
			}
		}
	}

	// Validate depends
	if len(r.Depends) > 0 {
		if v.hasDuplicates(r.Depends) {
			sl.ReportError(r.Depends, "Depends", "Depends", "has_duplicates", "")
		} else {
			ok, err := v.checkExist(ctx, r.Depends, validTypes)
			if err != nil {
				panic(err)
			}

			if !ok {
				sl.ReportError(r.Depends, "Depends", "Depends", "not_exist", "")
			} else if len(r.Impacts) > 0 {
				hash := make(map[string]bool)
				for _, e := range r.Impacts {
					hash[e] = true
				}
				for _, e := range r.Depends {
					if hash[e] {
						sl.ReportError(r.Depends, "Depends", "Depends", "has_duplicates_with", "Impacts")
						break
					}
				}
			}
		}
	}
}

func (v *basicValidator) hasDuplicates(values []string) bool {
	counts := make(map[string]int, len(values))
	for _, id := range values {
		counts[id]++
	}

	for _, c := range counts {
		if c > 1 {
			return true
		}
	}

	return false
}

func (v *basicValidator) checkExist(ctx context.Context, ids []string, validTypes []string) (bool, error) {
	count, err := v.dbClient.Collection(mongo.EntityMongoCollection).CountDocuments(ctx, bson.M{
		"_id":          bson.M{"$in": ids},
		"type":         bson.M{"$in": validTypes},
		"soft_deleted": bson.M{"$in": bson.A{false, nil}},
	})
	if err != nil {
		return false, err
	}

	return len(ids) == int(count), nil
}
