package entity

import (
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type searchPattern [][]pattern.FieldCondition

func (p searchPattern) ToMongoQuery() (bson.M, map[string]struct{}, map[string]struct{}, error) {
	if len(p) == 0 {
		return nil, nil, nil, nil
	}

	additionalMatchLookups := make(map[string]struct{})
	additionalMatchComputedFields := make(map[string]struct{})

	groupQueries := make([]bson.M, len(p))
	var err error

	now := time.Now()

	for i, group := range p {
		condQueries := make([]bson.M, len(group))
		for j, cond := range group {
			if infoName := pattern.GetEntityInfoName(cond.Field); infoName != "" {
				mongoField := "infos." + infoName + ".value"

				switch cond.FieldType {
				case pattern.FieldTypeString:
					condQueries[j], err = cond.Condition.StringToMongoQuery(mongoField, true)
				case pattern.FieldTypeInt:
					condQueries[j], err = cond.Condition.IntToMongoQuery(mongoField, true)
				case pattern.FieldTypeBool:
					condQueries[j], err = cond.Condition.BoolToMongoQuery(mongoField)
				case pattern.FieldTypeStringArray:
					condQueries[j], err = cond.Condition.StringArrayToMongoQuery(mongoField, true)
				case pattern.FieldTypeTimestamp:
					condQueries[j], err = cond.Condition.TimeToMongoQuery(mongoField, now)
				case "":
					condQueries[j], err = cond.Condition.RefToMongoQuery(mongoField)
				default:
					err = fmt.Errorf("invalid field type for %q field: %s", cond.Field, cond.FieldType)
				}

				if err != nil {
					return nil, nil, nil, fmt.Errorf("invalid condition for %q field: %w", cond.Field, err)
				}

				continue
			}

			if infoName := pattern.GetEntityComponentInfoName(cond.Field); infoName != "" {
				mongoField := "component_infos." + infoName + ".value"

				switch cond.FieldType {
				case pattern.FieldTypeString:
					condQueries[j], err = cond.Condition.StringToMongoQuery(mongoField, true)
				case pattern.FieldTypeInt:
					condQueries[j], err = cond.Condition.IntToMongoQuery(mongoField, true)
				case pattern.FieldTypeBool:
					condQueries[j], err = cond.Condition.BoolToMongoQuery(mongoField)
				case pattern.FieldTypeStringArray:
					condQueries[j], err = cond.Condition.StringArrayToMongoQuery(mongoField, true)
				case pattern.FieldTypeTimestamp:
					condQueries[j], err = cond.Condition.TimeToMongoQuery(mongoField, now)
				case "":
					condQueries[j], err = cond.Condition.RefToMongoQuery(mongoField)
				default:
					err = fmt.Errorf("invalid field type for %q field: %s", cond.Field, cond.FieldType)
				}

				if err != nil {
					return nil, nil, nil, fmt.Errorf("invalid condition for %q field: %w", cond.Field, err)
				}

				continue
			}

			switch cond.Field {
			case "_id", "name", "type", "category", "component", "connector",
				"import_source", "pbehavior_info.name", "pbehavior_info.reason",
				"pbehavior_info.type", "pbehavior_info.canonical_type":
				condQueries[j], err = cond.Condition.StringToMongoQuery(cond.Field, false)
			case "impact_level":
				condQueries[j], err = cond.Condition.IntToMongoQuery(cond.Field, true)
			case "impact_state", "state", "status":
				condQueries[j], err = cond.Condition.IntToMongoQuery(cond.Field, true)

				additionalMatchLookups["alarm"] = struct{}{}
				additionalMatchComputedFields[cond.Field] = struct{}{}
			case "ko_events", "ok_events":
				condQueries[j], err = cond.Condition.IntToMongoQuery(cond.Field, false)

				additionalMatchLookups["event_stats"] = struct{}{}
				additionalMatchComputedFields[cond.Field] = struct{}{}
			case "idle_since", "imported", "last_pbehavior_date", "last_event_date":
				condQueries[j], err = cond.Condition.TimeToMongoQuery(cond.Field, now)
			case "alarm_last_update_date":
				condQueries[j], err = cond.Condition.TimeToMongoQuery(cond.Field, now)

				additionalMatchLookups["alarm"] = struct{}{}
				additionalMatchComputedFields[cond.Field] = struct{}{}
			default:
				return nil, nil, nil, pattern.ErrUnsupportedField
			}

			if err != nil {
				return nil, nil, nil, err
			}
		}

		groupQueries[i] = bson.M{"$and": condQueries}
	}

	return bson.M{"$or": groupQueries}, additionalMatchLookups, additionalMatchComputedFields, nil
}
