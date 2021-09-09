package pbehavior

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type Validator struct {
	dbClient mongo.DbClient
}

const validationCollection = "unexisted_collection_for_validation"

func NewValidator(client mongo.DbClient) *Validator {
	return &Validator{dbClient: client}
}

func (v *Validator) ValidateEditRequest(ctx context.Context, sl validator.StructLevel) {
	request := sl.Current().Interface().(EditRequest)

	// Only pause pbehavior have optional Stop
	if request.Stop == nil && request.Type != "" {
		foundType := pbehavior.Type{}
		err := v.dbClient.Collection(pbehavior.TypeCollectionName).
			FindOne(ctx, bson.M{"_id": request.Type}).Decode(&foundType)
		if err == nil {
			if foundType.Type != pbehavior.TypePause {
				sl.ReportError(request.Stop, "Stop", "Stop", "required", "")
			}
		} else if err != mongodriver.ErrNoDocuments {
			panic(err)
		}
	}
	// Stop must be > Start
	if request.Stop != nil && request.Stop.Before(request.Start.Time) {
		sl.ReportError(request.Stop, "Stop", "Stop", "gtfield", "Start")
	}
	// Filter must be valid mongo filter
	if request.Filter != nil {
		_, err := v.dbClient.Collection(validationCollection).
			Find(ctx, request.Filter)
		if err != nil {
			sl.ReportError(request.Stop, "Filter", "Filter", "entityfilter", "Filter")
		}
	}
}
