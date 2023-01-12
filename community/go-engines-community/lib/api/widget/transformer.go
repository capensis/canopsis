package widget

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type RequestTransformer struct {
	patternFieldsTransformer common.PatternFieldsTransformer
	widgetTemplateCollection mongo.DbCollection
}

func NewRequestTransformer(
	patternFieldsTransformer common.PatternFieldsTransformer,
	client mongo.DbClient,
) *RequestTransformer {
	return &RequestTransformer{
		patternFieldsTransformer: patternFieldsTransformer,
		widgetTemplateCollection: client.Collection(mongo.WidgetTemplateMongoCollection),
	}
}

func (t *RequestTransformer) Transform(ctx context.Context, r *EditRequest) error {
	err := t.transformPatternFields(ctx, r)
	if err != nil {
		return err
	}

	return t.transformColumnFields(ctx, r)
}

func (t *RequestTransformer) transformPatternFields(ctx context.Context, r *EditRequest) error {
	var err error
	for i := range r.Filters {
		r.Filters[i].AlarmPatternFieldsRequest, err = t.patternFieldsTransformer.TransformAlarmPatternFieldsRequest(ctx, r.Filters[i].AlarmPatternFieldsRequest)
		if err != nil {
			if err == common.ErrNotExistCorporateAlarmPattern {
				return common.NewValidationError(fmt.Sprintf("filters.%d.corporate_alarm_pattern", i), err)
			}
			return err
		}
		r.Filters[i].EntityPatternFieldsRequest, err = t.patternFieldsTransformer.TransformEntityPatternFieldsRequest(ctx, r.Filters[i].EntityPatternFieldsRequest)
		if err != nil {
			if err == common.ErrNotExistCorporateEntityPattern {
				return common.NewValidationError(fmt.Sprintf("filters.%d.corporate_entity_pattern", i), err)
			}
			return err
		}
		r.Filters[i].PbehaviorPatternFieldsRequest, err = t.patternFieldsTransformer.TransformPbehaviorPatternFieldsRequest(ctx, r.Filters[i].PbehaviorPatternFieldsRequest)
		if err != nil {
			if err == common.ErrNotExistCorporatePbehaviorPattern {
				return common.NewValidationError(fmt.Sprintf("filters.%d.corporate_pbehavior_pattern", i), err)
			}
			return err
		}
	}

	return nil
}

func (t *RequestTransformer) transformColumnFields(ctx context.Context, r *EditRequest) error {
	columnParametersByType := view.GetWidgetColumnParameters()[r.Type]
	for tplType, columnParameters := range columnParametersByType {
		for _, parameter := range columnParameters {
			parameters := r.Parameters.RemainParameters
			key := parameter
			parts := strings.Split(parameter, ".")
			if len(parts) > 1 {
				key = parts[len(parts)-1]
				var ok bool
				for i := 0; i < len(parts)-1; i++ {
					parameters, ok = parameters[parts[i]].(map[string]any)
					if !ok {
						break
					}
				}
				if !ok {
					continue
				}
			}

			tplId, ok := parameters[key+"Template"].(string)
			if !ok || tplId == "" {
				continue
			}
			tpl := view.WidgetTemplate{}
			err := t.widgetTemplateCollection.FindOne(ctx, bson.M{
				"_id":  tplId,
				"type": tplType,
			}).Decode(&tpl)
			if err != nil {
				if errors.Is(err, mongodriver.ErrNoDocuments) {
					return common.NewValidationError("parameters."+parameter+"Template", errors.New("Template doesn't exist."))
				}
				return err
			}

			parameters[key+"TemplateTitle"] = tpl.Title
			parameters[key] = tpl.Columns
		}
	}

	return nil
}
