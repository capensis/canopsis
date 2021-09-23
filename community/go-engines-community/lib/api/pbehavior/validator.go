package pbehavior

import (
	"context"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		rType   *string
		rStart  *types.CpsTime
		rStop   *types.CpsTime
		rFilter interface{}
	)
	slr := sl.Current().Interface()
	switch r := slr.(type) {
	case EditRequest:
		rType = &r.Type
		rStart = &r.Start
		rStop = r.Stop
		rFilter = r.Filter
	case PatchRequest:
		rType = r.Type
		rStart = r.Start
		rFilter = r.Filter
		if r.Stop.isSet {
			rStop = r.Stop.CpsTime
		}
	}

	// Only pause pbehavior have optional Stop
	if rStop == nil && rType != nil && *rType != "" {
		foundType := pbehavior.Type{}
		err := v.dbClient.Collection(pbehavior.TypeCollectionName).
			FindOne(ctx, bson.M{"_id": rType}).Decode(&foundType)
		if err == nil {
			if foundType.Type != pbehavior.TypePause {
				sl.ReportError(rStop, "Stop", "Stop", "required", "")
			}
		} else if err == mongodriver.ErrNoDocuments {
			sl.ReportError(rStop, "Type", "Type", "not_exist", "")
		} else {
			panic(err)
		}
	}
	// Stop must be > Start
	if rStop != nil && rStart != nil && rStop.Before(rStart.Time) {
		sl.ReportError(rStop, "Stop", "Stop", "gtfield", "Start")
	}
	// Filter must be valid mongo filter
	if rFilter != nil {
		_, err := v.dbClient.Collection(validationCollection).
			Find(ctx, rFilter)
		if err != nil {
			sl.ReportError(rStop, "Filter", "Filter", "entityfilter", "Filter")
		}
	}
}
