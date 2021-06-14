package heartbeat

import (
	"context"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/heartbeat"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type Validator struct {
	dbClient mongo.DbClient
}

func NewValidator(client mongo.DbClient) *Validator {
	return &Validator{dbClient: client}
}

func (v *Validator) ValidateCreateRequest(sl validator.StructLevel) {
	request := sl.Current().Interface().(CreateRequest)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//check custom id if exists
	if request.ID != "" {
		foundHeartbeat := heartbeat.Heartbeat{}
		err := v.dbClient.Collection(heartbeat.HeartbeatCollectionName).FindOne(ctx, bson.M{"_id": request.ID}).Decode(&foundHeartbeat)
		if err == nil {
			sl.ReportError("_id", "ID", "ID", "unique", "")
		} else {
			if err != mongodriver.ErrNoDocuments {
				panic(err)
			}
		}
	}

	foundHeartbeat := heartbeat.Heartbeat{}
	err := v.dbClient.Collection(heartbeat.HeartbeatCollectionName).FindOne(ctx, bson.M{"name": request.Name}).Decode(&foundHeartbeat)
	if err == nil {
		sl.ReportError("name", "Name", "Name", "unique", "")
	} else {
		if err != mongodriver.ErrNoDocuments {
			panic(err)
		}
	}
}
