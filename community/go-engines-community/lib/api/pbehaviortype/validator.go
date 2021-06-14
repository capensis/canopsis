package pbehaviortype

import (
	"context"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"strings"
)

type Validator interface {
	ValidateTypeCreateRequest(sl validator.StructLevel)
	ValidateTypeUpdateRequest(sl validator.StructLevel)
}

type baseValidator struct {
	dbClient     mongo.DbClient
	dbCollection mongo.DbCollection
}

func NewValidator(dbClient mongo.DbClient) Validator {
	return &baseValidator{
		dbClient:     dbClient,
		dbCollection: dbClient.Collection(pbehavior.TypeCollectionName),
	}
}

func (v *baseValidator) ValidateTypeCreateRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(CreateRequest)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//check custom id if exists
	if r.ID != "" {
		foundType := pbehavior.Type{}
		err := v.dbCollection.FindOne(ctx, bson.M{"_id": r.ID}).Decode(&foundType)
		if err == nil {
			sl.ReportError("_id", "ID", "ID", "unique", "")
		} else {
			if err != mongodriver.ErrNoDocuments {
				panic(err)
			}
		}
	}

	v.validateType(sl, r.Type)
	v.validateName(ctx, sl, r.Name, r.ID)
	v.validatePriority(ctx, sl, *r.Priority, r.ID)
}

func (v *baseValidator) ValidateTypeUpdateRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(UpdateRequest)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	v.validateType(sl, r.Type)
	v.validateName(ctx, sl, r.Name, r.ID)
	v.validatePriority(ctx, sl, *r.Priority, r.ID)
}

func (v *baseValidator) validateType(sl validator.StructLevel, rType string) {
	switch rType {
	case pbehavior.TypeActive,
		pbehavior.TypeInactive,
		pbehavior.TypeMaintenance,
		pbehavior.TypePause:
	/*do nothing*/
	default:
		types := []string{
			pbehavior.TypeActive,
			pbehavior.TypeInactive,
			pbehavior.TypeMaintenance,
			pbehavior.TypePause,
		}
		param := strings.Join(types, " ")
		sl.ReportError(rType, "Type", "type", "oneof", param)
	}
}

func (v *baseValidator) validateName(ctx context.Context, sl validator.StructLevel, name, id string) {
	foundType := pbehavior.Type{}
	err := v.dbCollection.FindOne(ctx, bson.M{"name": name}).Decode(&foundType)
	if err == nil {
		if foundType.ID != id {
			sl.ReportError(name, "Name", "Name", "unique", "")
		}
	} else if err != mongodriver.ErrNoDocuments {
		panic(err)
	}
}

func (v *baseValidator) validatePriority(ctx context.Context, sl validator.StructLevel, priority int, id string) {
	foundType := pbehavior.Type{}
	err := v.dbCollection.FindOne(ctx, bson.M{"priority": priority}).Decode(&foundType)
	if err == nil {
		if foundType.ID != id {
			sl.ReportError(priority, "Priority", "Priority", "unique", "")
		}
	} else if err != mongodriver.ErrNoDocuments {
		panic(err)
	}
}