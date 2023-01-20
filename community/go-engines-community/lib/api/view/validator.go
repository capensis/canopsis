package view

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type Validator interface {
	ValidateEditRequest(ctx context.Context, sl validator.StructLevel)
}

type baseValidator struct {
	dbClient mongo.DbClient
}

func NewValidator(dbClient mongo.DbClient) Validator {
	return &baseValidator{
		dbClient: dbClient,
	}
}

func (v *baseValidator) ValidateEditRequest(ctx context.Context, sl validator.StructLevel) {
	r := sl.Current().Interface().(EditRequest)
	// Validate group
	if r.Group != "" {
		err := v.dbClient.Collection(mongo.ViewGroupMongoCollection).FindOne(ctx, bson.M{"_id": r.Group}).Err()
		if err != nil {
			if err == mongodriver.ErrNoDocuments {
				sl.ReportError(r.Group, "Group", "Group", "not_exist", "")
			} else {
				panic(err)
			}
		}
	}
}

func ValidateEditPositionRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(EditPositionRequest)

	if len(r.Items) > 0 {
		exists := make(map[string]bool, len(r.Items))
		existsView := make(map[string]bool, len(r.Items))
		for _, item := range r.Items {
			if exists[item.ID] {
				sl.ReportError(r.Items, "Items", "Item", "has_duplicates", "")
				return
			}

			exists[item.ID] = true

			for _, view := range item.Views {
				if existsView[view] {
					sl.ReportError(r.Items, "Items", "Item", "has_duplicates", "")
					return
				}

				existsView[view] = true
			}
		}
	}
}
