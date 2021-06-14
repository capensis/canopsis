package pbehaviorexception

import (
	"context"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type Validator interface {
	ValidateExceptionCreateRequest(sl validator.StructLevel)
	ValidateExceptionUpdateRequest(sl validator.StructLevel)
	ValidateExdateRequest(sl validator.StructLevel)
}

type baseValidator struct {
	dbClient     mongo.DbClient
	dbCollection mongo.DbCollection
}

func NewValidator(dbClient mongo.DbClient) Validator {
	return &baseValidator{
		dbClient:     dbClient,
		dbCollection: dbClient.Collection(pbehavior.ExceptionCollectionName),
	}
}

func (v *baseValidator) ValidateExceptionCreateRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(CreateRequest)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//check custom id if exists
	if r.ID != "" {
		foundException := pbehavior.Exception{}
		err := v.dbCollection.FindOne(ctx, bson.M{"_id": r.ID}).Decode(&foundException)
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

func (v *baseValidator) ValidateExceptionUpdateRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(UpdateRequest)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	v.validateName(ctx, sl, r.Name, r.ID)
}

func (v *baseValidator) validateName(ctx context.Context, sl validator.StructLevel, name, id string) {
	foundException := pbehavior.Exception{}
	err := v.dbCollection.FindOne(ctx, bson.M{"name": name}).Decode(&foundException)
	if err == nil {
		if foundException.ID != id {
			sl.ReportError(name, "Name", "Name", "unique", "")
		}
	} else if err != mongodriver.ErrNoDocuments {
		panic(err)
	}
}

func (baseValidator) ValidateExdateRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(ExdateRequest)

	if r.End.Before(r.Begin.Time) {
		sl.ReportError(r.End, "End", "End", "gtfield", "Begin")
	}
}
