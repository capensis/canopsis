package user

import (
	"context"
	"errors"
	"strconv"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/password"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type Validator interface {
	ValidateCreateRequest(ctx context.Context, sl validator.StructLevel)
	ValidateUpdateRequest(ctx context.Context, sl validator.StructLevel)
	ValidateBulkUpdateRequestItem(ctx context.Context, sl validator.StructLevel)
}

type baseValidator struct {
	dbCollection           mongo.DbCollection
	dbRoleCollection       mongo.DbCollection
	dbViewCollection       mongo.DbCollection
	dbColorThemeCollection mongo.DbCollection
}

func NewValidator(dbClient mongo.DbClient) Validator {
	return &baseValidator{
		dbCollection:           dbClient.Collection(mongo.UserCollection),
		dbRoleCollection:       dbClient.Collection(mongo.RoleCollection),
		dbViewCollection:       dbClient.Collection(mongo.ViewMongoCollection),
		dbColorThemeCollection: dbClient.Collection(mongo.ColorThemeCollection),
	}
}

func (v *baseValidator) ValidateBulkUpdateRequestItem(ctx context.Context, sl validator.StructLevel) {
	r := sl.Current().Interface().(BulkUpdateRequestItem)

	v.validateEditRequest(ctx, sl, r.ID, r.EditRequest)
	v.validatePassword(sl, r.EditRequest, r.ID)
}

func (v *baseValidator) ValidateCreateRequest(ctx context.Context, sl validator.StructLevel) {
	r := sl.Current().Interface().(CreateRequest)

	v.validateEditRequest(ctx, sl, "", r.EditRequest)

	// Validate source and external_id
	if r.Source == "" && r.ExternalID != "" {
		sl.ReportError(r.Source, "Source", "Source", "required_with", "ExternalID")
	}

	if r.ExternalID == "" && r.Source != "" {
		sl.ReportError(r.ExternalID, "ExternalID", "ExternalID", "required_with", "Source")
	}

	// Validate password
	if r.Source != "" {
		if r.Password != "" {
			sl.ReportError(r.Source, "Source", "Source", "required_not_both", "Password")
		}
	} else {
		v.validatePassword(sl, r.EditRequest, "")
	}
}

func (v *baseValidator) ValidateUpdateRequest(ctx context.Context, sl validator.StructLevel) {
	r := sl.Current().Interface().(UpdateRequest)

	v.validateEditRequest(ctx, sl, r.ID, r.EditRequest)
	v.validatePassword(sl, r.EditRequest, r.ID)
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
		} else if errors.Is(err, mongodriver.ErrNoDocuments) {
			// Check unique by name
			err := v.dbCollection.FindOne(ctx, bson.M{"name": r.Name}).Decode(&res)
			if err == nil {
				if res.ID != id {
					sl.ReportError(r.Name, "Name", "Name", "unique", "")
				}
			} else if !errors.Is(err, mongodriver.ErrNoDocuments) {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	// Validate default view
	if r.DefaultView != "" {
		err := v.dbViewCollection.FindOne(ctx, bson.M{"_id": r.DefaultView}).Err()
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				sl.ReportError(r.DefaultView, "DefaultView", "DefaultView", "not_exist", "")
			} else {
				panic(err)
			}
		}
	}
	// Validate role
	if len(r.Roles) > 0 {
		c, err := v.dbRoleCollection.CountDocuments(ctx, bson.M{"_id": bson.M{"$in": r.Roles}})
		if err != nil {
			panic(err)
		}
		if int(c) < len(r.Roles) {
			sl.ReportError(r.Roles, "Roles", "Roles", "not_exist", "")
		}
	}
	// Validate UITheme
	if r.UITheme != "" {
		err := v.dbColorThemeCollection.FindOne(ctx, bson.M{"_id": r.UITheme}).Err()
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				sl.ReportError(r.UITheme, "UITheme", "UITheme", "not_exist", "")
			} else {
				panic(err)
			}
		}
	}
}

func (v *baseValidator) validatePassword(sl validator.StructLevel, r EditRequest, id string) {
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
}
