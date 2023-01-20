package user

import (
	"context"
	"strconv"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/password"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type Validator interface {
	ValidateRequest(ctx context.Context, sl validator.StructLevel)
	ValidateBulkUpdateRequestItem(ctx context.Context, sl validator.StructLevel)
}

type baseValidator struct {
	dbClient     mongo.DbClient
	dbCollection mongo.DbCollection
}

func NewValidator(dbClient mongo.DbClient) Validator {
	return &baseValidator{
		dbClient:     dbClient,
		dbCollection: dbClient.Collection(mongo.RightsMongoCollection),
	}
}

func (v *baseValidator) ValidateBulkUpdateRequestItem(ctx context.Context, sl validator.StructLevel) {
	r := sl.Current().Interface().(BulkUpdateRequestItem)

	v.validateEditRequest(ctx, sl, r.ID, r.EditRequest)
}

func (v *baseValidator) ValidateRequest(ctx context.Context, sl validator.StructLevel) {
	r := sl.Current().Interface().(Request)

	v.validateEditRequest(ctx, sl, r.ID, r.EditRequest)
}

func (v *baseValidator) validateEditRequest(ctx context.Context, sl validator.StructLevel, id string, r EditRequest) {
	if r.Name != "" {
		// Check unique by id
		res := struct {
			ID string `bson:"_id"`
		}{}
		err := v.dbCollection.FindOne(ctx, bson.M{"_id": r.Name}).Decode(&res)
		if err == nil {
			if res.ID != id {
				sl.ReportError(r.Name, "Name", "Name", "unique", "")
			}
		} else if err == mongodriver.ErrNoDocuments {
			// Check unique by name
			err := v.dbCollection.FindOne(ctx, bson.M{"crecord_name": r.Name}).Decode(&res)
			if err == nil {
				if res.ID != id {
					sl.ReportError(r.Name, "Name", "Name", "unique", "")
				}
			} else if err != mongodriver.ErrNoDocuments {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	// Validate password
	if r.Password == "" {
		if id == "" {
			sl.ReportError(r.Password, "Password", "Password", "required", "")
		}
	} else {
		if len(r.Password) < password.MinLength {
			sl.ReportError(r.Password, "Password", "Password", "min", strconv.Itoa(password.MinLength))
		}
		if len(r.Password) > password.MaxLength {
			sl.ReportError(r.Password, "Password", "Password", "max", strconv.Itoa(password.MaxLength))
		}
	}
	// Validate default view
	if r.DefaultView != "" {
		err := v.dbClient.Collection(mongo.ViewMongoCollection).FindOne(ctx, bson.M{"_id": r.DefaultView}).Err()
		if err != nil {
			if err == mongodriver.ErrNoDocuments {
				sl.ReportError(r.DefaultView, "DefaultView", "DefaultView", "not_exist", "")
			} else {
				panic(err)
			}
		}
	}
	// Validate role
	if r.Role != "" {
		err := v.dbCollection.FindOne(ctx, bson.M{"_id": r.Role, "crecord_type": model.LineTypeRole}).Err()
		if err != nil {
			if err == mongodriver.ErrNoDocuments {
				sl.ReportError(r.Role, "Role", "Role", "not_exist", "")
			} else {
				panic(err)
			}
		}
	}
}
