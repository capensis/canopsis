package pbehaviorreason

import (
	"context"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type Validator struct {
	dbClient     mongo.DbClient
	dbCollection mongo.DbCollection
}

func (v Validator) ValidateReasonCreateRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(CreateRequest)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//check custom id if exists
	if r.ID != "" {
		foundReason := pbehavior.Reason{}
		err := v.dbCollection.FindOne(ctx, bson.M{"_id": r.ID}).Decode(&foundReason)
		if err == nil {
			sl.ReportError("_id", "ID", "ID", "unique", "")
		} else {
			if err != mongodriver.ErrNoDocuments {
				panic(err)
			}
		}
	}

	v.validateName(ctx, sl, r.Name, r.ID)
}

func (v Validator) ValidateReasonUpdateRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(UpdateRequest)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	v.validateName(ctx, sl, r.Name, r.ID)
}

func (v Validator) validateName(ctx context.Context, sl validator.StructLevel, name, id string) {
	foundReason := pbehavior.Reason{}
	err := v.dbCollection.FindOne(ctx, bson.M{"name": name}).Decode(&foundReason)
	if err == nil {
		if foundReason.ID != id {
			sl.ReportError(name, "Name", "Name", "unique", "")
		}
	} else {
		if err != mongodriver.ErrNoDocuments {
			panic(err)
		}
	}
}

func NewValidator(client mongo.DbClient) Validator {
	return Validator{
		dbClient:     client,
		dbCollection: client.Collection(pbehavior.ReasonCollectionName),
	}
}
