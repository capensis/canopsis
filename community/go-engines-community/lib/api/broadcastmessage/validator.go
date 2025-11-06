package broadcastmessage

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
)

type baseValidator struct {
	dbClient mongo.DbClient
}

func (v *baseValidator) Validate(ctx context.Context, sl validator.StructLevel) {
	r := sl.Current().Interface().(CreateRequest)

	if r.ID != "" {
		err := v.dbClient.Collection(mongo.BroadcastMessageCollection).FindOne(ctx, bson.M{"_id": r.ID}).Err()
		if err == nil {
			sl.ReportError("_id", "ID", "ID", "exist", "")
		} else if !errors.Is(err, mongodriver.ErrNoDocuments) {
			panic(err)
		}
	}

}

func NewValidator(dbClient mongo.DbClient) *baseValidator {
	return &baseValidator{dbClient: dbClient}
}
