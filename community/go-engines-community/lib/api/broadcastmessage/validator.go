package broadcastmessage

import (
	"context"

	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type broadcastMessageValidator struct {
	dbClient mongo.DbClient
}

func (v *broadcastMessageValidator) Validate(ctx context.Context, sl validator.StructLevel) {
	r := sl.Current().Interface().(BroadcastMessage)

	if r.ID != "" {
		err := v.dbClient.Collection(mongo.BroadcastMessageMongoCollection).FindOne(ctx, bson.M{"_id": r.ID}).Err()
		if err == nil {
			sl.ReportError("_id", "ID", "ID", "unique", "")
		} else if err != mongodriver.ErrNoDocuments {
			panic(err)
		}
	}

}

func NewValidator(dbClient mongo.DbClient) *broadcastMessageValidator {
	return &broadcastMessageValidator{dbClient: dbClient}
}
