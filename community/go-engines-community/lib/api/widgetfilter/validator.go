package widgetfilter

import (
	"context"
	"errors"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type Validator struct {
	collection mongo.DbCollection
}

func NewValidator(client mongo.DbClient) *Validator {
	return &Validator{
		collection: client.Collection(mongo.WidgetFiltersMongoCollection),
	}
}

func (v *Validator) ValidateEditRequest(ctx context.Context, sl validator.StructLevel) {
	r := sl.Current().Interface().(EditRequest)

	if len(r.AlarmPattern) == 0 && r.CorporateAlarmPattern == "" &&
		len(r.EntityPattern) == 0 && r.CorporateEntityPattern == "" &&
		len(r.PbehaviorPattern) == 0 && r.CorporatePbehaviorPattern == "" {

		if r.ID != "" {
			err := v.collection.FindOne(ctx, bson.M{"_id": r.ID, "old_mongo_query": bson.M{"$ne": nil}}).Err()
			if err == nil {
				return
			} else if !errors.Is(err, mongodriver.ErrNoDocuments) {
				panic(err)
			}
		}

		sl.ReportError(r.AlarmPattern, "AlarmPattern", "AlarmPattern", "required", "")
		sl.ReportError(r.CorporateAlarmPattern, "CorporateAlarmPattern", "CorporateAlarmPattern", "required", "")
		sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "required", "")
		sl.ReportError(r.CorporateEntityPattern, "CorporateEntityPattern", "CorporateEntityPattern", "required", "")
		sl.ReportError(r.PbehaviorPattern, "PbehaviorPattern", "PbehaviorPattern", "required", "")
		sl.ReportError(r.CorporatePbehaviorPattern, "CorporatePbehaviorPattern", "CorporatePbehaviorPattern", "required", "")
	}
}
