package pbehavior

import (
	"context"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
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

func (v *Validator) ValidateCreateRequest(sl validator.StructLevel) {
	request := sl.Current().Interface().(CreateRequest)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//check custom id if exists
	if request.ID != "" {
		foundPBehavior := pbehavior.PBehavior{}
		err := v.dbClient.Collection(pbehavior.PBehaviorCollectionName).FindOne(ctx, bson.M{"_id": request.ID}).Decode(&foundPBehavior)
		if err == nil {
			sl.ReportError("_id", "ID", "ID", "unique", "")
		} else {
			if err != mongodriver.ErrNoDocuments {
				panic(err)
			}
		}
	}
}

func (v *Validator) ValidateEditRequest(sl validator.StructLevel) {
	request := sl.Current().Interface().(EditRequest)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Only pause pbehavior have optional Stop
	if request.Stop == nil && request.Type != "" {
		foundType := pbehavior.Type{}
		err := v.dbClient.Collection(pbehavior.TypeCollectionName).
			FindOne(ctx, bson.M{"_id": request.Type}).Decode(&foundType)
		if err == nil {
			if foundType.Type != pbehavior.TypePause {
				sl.ReportError(request.Stop, "Stop", "Stop", "required", "")
			}
		} else if err == mongodriver.ErrNoDocuments {
			sl.ReportError(request.Stop, "Type", "Type", "not_exist", "")
		} else {
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
