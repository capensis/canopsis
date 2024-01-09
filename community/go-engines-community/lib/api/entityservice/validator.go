package entityservice

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type Validator interface {
	ValidateEditRequest(ctx context.Context, sl validator.StructLevel)
	ValidateCreateRequest(sl validator.StructLevel)
	ValidateUpdateRequest(sl validator.StructLevel)
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

	if r.CorporateEntityPattern == "" && len(r.EntityPattern) > 0 &&
		!match.ValidateEntityPattern(r.EntityPattern, common.GetForbiddenFieldsInEntityPattern(mongo.EntityMongoCollection)) {
		sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "entity_pattern", "")
	}
}

func (v *basicValidator) ValidateCreateRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(CreateRequest)
	if len(r.EntityPattern) == 0 && r.CorporateEntityPattern == "" {
		sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "required", "")
	}
}

func (v *basicValidator) ValidateUpdateRequest(sl validator.StructLevel) {
	corporateEntityPattern := ""
	var entityPattern pattern.Entity
	switch r := sl.Current().Interface().(type) {
	case UpdateRequest:
		entityPattern = r.EntityPattern
		corporateEntityPattern = r.CorporateEntityPattern
	case BulkUpdateRequestItem:
		entityPattern = r.EntityPattern
		corporateEntityPattern = r.CorporateEntityPattern
	}

	if len(entityPattern) == 0 && corporateEntityPattern == "" {
		sl.ReportError(entityPattern, "EntityPattern", "EntityPattern", "required", "")
	}
}

func (v *basicValidator) validateCategory(ctx context.Context, sl validator.StructLevel, category string) {
	if category != "" {
		err := v.dbClient.Collection(mongo.EntityCategoryMongoCollection).
			FindOne(ctx, bson.M{"_id": category}).Err()
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				sl.ReportError(category, "Category", "Category", "not_exist", "")
			} else {
				panic(err)
			}
		}
	}
}
