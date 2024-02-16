package account

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
	ValidateEditRequest(ctx context.Context, sl validator.StructLevel)
}

type baseValidator struct {
	dbViewCollection       mongo.DbCollection
	dbColorThemeCollection mongo.DbCollection
}

func NewValidator(dbClient mongo.DbClient) Validator {
	return &baseValidator{
		dbViewCollection:       dbClient.Collection(mongo.ViewMongoCollection),
		dbColorThemeCollection: dbClient.Collection(mongo.ColorThemeCollection),
	}
}

func (v *baseValidator) ValidateEditRequest(ctx context.Context, sl validator.StructLevel) {
	r := sl.Current().Interface().(EditRequest)
	// Validate password
	if r.Password != "" {
		if len(r.Password) < password.MinLength {
			sl.ReportError(r.Password, "Password", "Password", "min", strconv.Itoa(password.MinLength))
		}
		if len(r.Password) > password.MaxLength {
			sl.ReportError(r.Password, "Password", "Password", "max", strconv.Itoa(password.MaxLength))
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
