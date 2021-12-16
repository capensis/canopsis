package widgetfilter

import (
	"context"
	"encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
)

const validationCollection = "unexisted_collection_for_validation"

type Validator interface {
	ValidateEditRequest(ctx context.Context, sl validator.StructLevel)
}

type baseValidator struct {
	client mongo.DbClient
}

func NewValidator(client mongo.DbClient) Validator {
	return &baseValidator{
		client: client,
	}
}

func (v *baseValidator) ValidateEditRequest(ctx context.Context, sl validator.StructLevel) {
	r := sl.Current().Interface().(EditRequest)

	if r.Query != "" {
		var query bson.M
		err := json.Unmarshal([]byte(r.Query), &query)
		if err == nil {
			_, err = v.client.Collection(validationCollection).Find(ctx, query)
		}

		if err != nil {
			sl.ReportError(r.Query, "Query", "Query", "filter", "")
		}
	}
}
