package appinfo

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type Validator struct {
	dbColorThemeCollection mongo.DbCollection
}

func NewValidator(dbClient mongo.DbClient) *Validator {
	return &Validator{
		dbColorThemeCollection: dbClient.Collection(mongo.ColorThemeCollection),
	}
}

func (v *Validator) ValidateRequest(ctx context.Context, sl validator.StructLevel) {
	r := sl.Current().Interface().(UserInterfaceConf)

	if r.DefaultColorTheme != "" {
		err := v.dbColorThemeCollection.FindOne(ctx, bson.M{"_id": r.DefaultColorTheme}).Err()
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				sl.ReportError(r.DefaultColorTheme, "DefaultColorTheme", "DefaultColorTheme", "not_exist", "")
			} else {
				panic(err)
			}
		}
	}
}
