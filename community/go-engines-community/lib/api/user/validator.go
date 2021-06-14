package user

import (
	"context"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"git.canopsis.net/canopsis/go-engines/lib/security/model"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"strconv"
)

const minPasswordLength = 8
const maxPasswordLength = 255

type Validator interface {
	ValidateEditRequest(ctx context.Context, sl validator.StructLevel)
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

func (v *baseValidator) ValidateEditRequest(ctx context.Context, sl validator.StructLevel) {
	r := sl.Current().Interface().(EditRequest)
	// Validate name
	if r.Name != "" {
		// Check unique by id
		res := struct {
			ID string `bson:"_id"`
		}{}
		err := v.dbCollection.FindOne(ctx, bson.M{"_id": r.Name}).Decode(&res)
		if err == nil {
			if res.ID != r.ID {
				sl.ReportError(r.Name, "Name", "Name", "unique", "")
			}
		} else if err == mongodriver.ErrNoDocuments {
			// Check unique by name
			err := v.dbCollection.FindOne(ctx, bson.M{"crecord_name": r.Name}).Decode(&res)
			if err == nil {
				if res.ID != r.ID {
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
		if r.ID == "" {
			sl.ReportError(r.Password, "Password", "Password", "required", "")
		}
	} else {
		if len(r.Password) < minPasswordLength {
			sl.ReportError(r.Password, "Password", "Password", "min", strconv.Itoa(minPasswordLength))
		}
		if len(r.Password) > maxPasswordLength {
			sl.ReportError(r.Password, "Password", "Password", "max", strconv.Itoa(maxPasswordLength))
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
