package db

import (
	"fmt"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"go.mongodb.org/mongo-driver/bson"
)

func EntityPatternToMongoQuery(p pattern.Entity, prefix string) (bson.M, error) {
	groupQueries, err := getEntityPatternGroupMongoQueries(p, prefix)
	if err != nil || len(groupQueries) == 0 {
		return nil, err
	}
	return bson.M{"$or": groupQueries}, nil
}

func EntityPatternToNegativeMongoQuery(p pattern.Entity, prefix string) (bson.M, error) {
	groupQueries, err := getEntityPatternGroupMongoQueries(p, prefix)
	if err != nil || len(groupQueries) == 0 {
		return nil, err
	}
	return bson.M{"$nor": groupQueries}, nil
}

func EntityPatternToSql(p pattern.Entity, prefix string) (string, error) {
	if len(p) == 0 {
		return "", nil
	}

	if prefix != "" {
		prefix += "."
	}

	emptyEntity := &types.Entity{}
	groupQueries := make([]string, len(p))
	var err error

	fieldMap := map[string]string{"_id": "custom_id"}

	for i, group := range p {
		condQueries := make([]string, len(group))
		for j, cond := range group {
			if infoName := pattern.GetEntityInfoName(cond.Field); infoName != "" {
				switch cond.FieldType {
				case pattern.FieldTypeString:
					condQueries[j], err = cond.Condition.StringToSqlJson("infos", infoName)
				case pattern.FieldTypeInt:
					condQueries[j], err = cond.Condition.IntToSqlJson("infos", infoName)
				case pattern.FieldTypeBool:
					condQueries[j], err = cond.Condition.BoolToSqlJson("infos", infoName)
				case pattern.FieldTypeStringArray:
					condQueries[j], err = cond.Condition.StringArrayToSqlJson("infos", infoName)
				case "":
					condQueries[j], err = cond.Condition.RefToSqlJson("infos", infoName)
				default:
					err = fmt.Errorf("invalid field type for %q field: %s", cond.Field, cond.FieldType)
				}

				if err != nil {
					return "", err
				}

				continue
			}

			if infoName := pattern.GetEntityComponentInfoName(cond.Field); infoName != "" {
				switch cond.FieldType {
				case pattern.FieldTypeString:
					condQueries[j], err = cond.Condition.StringToSqlJson("component_infos", infoName)
				case pattern.FieldTypeInt:
					condQueries[j], err = cond.Condition.IntToSqlJson("component_infos", infoName)
				case pattern.FieldTypeBool:
					condQueries[j], err = cond.Condition.BoolToSqlJson("component_infos", infoName)
				case pattern.FieldTypeStringArray:
					condQueries[j], err = cond.Condition.StringArrayToSqlJson("component_infos", infoName)
				case "":
					condQueries[j], err = cond.Condition.RefToSqlJson("component_infos", infoName)
				default:
					err = fmt.Errorf("invalid field type for %q field: %s", cond.Field, cond.FieldType)
				}
				if err != nil {
					return "", err
				}

				continue
			}

			sqlField := cond.Field
			if v, ok := fieldMap[sqlField]; ok {
				sqlField = v
			}
			sqlField = prefix + sqlField
			if _, ok := emptyEntity.GetStringField(cond.Field); ok {
				condQueries[j], err = cond.Condition.StringToSql(sqlField)
			} else if _, ok := emptyEntity.GetIntField(cond.Field); ok {
				condQueries[j], err = cond.Condition.IntToSql(sqlField)
			} else {
				err = pattern.ErrUnsupportedField
			}
			if err != nil {
				return "", err
			}
		}

		groupQueries[i] = fmt.Sprintf("(%s)", strings.Join(condQueries, " AND "))
	}

	return strings.Join(groupQueries, " OR "), nil
}

func getEntityPatternGroupMongoQueries(p pattern.Entity, prefix string) ([]bson.M, error) {
	if len(p) == 0 {
		return nil, nil
	}

	if prefix != "" {
		prefix += "."
	}

	emptyEntity := &types.Entity{}
	groupQueries := make([]bson.M, len(p))
	var err error

	for i, group := range p {
		condQueries := make([]bson.M, len(group))
		for j, cond := range group {
			if infoName := pattern.GetEntityInfoName(cond.Field); infoName != "" {
				mongoField := prefix + "infos." + infoName + ".value"

				switch cond.FieldType {
				case pattern.FieldTypeString:
					condQueries[j], err = cond.Condition.StringToMongoQuery(mongoField, true)
				case pattern.FieldTypeInt:
					condQueries[j], err = cond.Condition.IntToMongoQuery(mongoField, true)
				case pattern.FieldTypeBool:
					condQueries[j], err = cond.Condition.BoolToMongoQuery(mongoField)
				case pattern.FieldTypeStringArray:
					condQueries[j], err = cond.Condition.StringArrayToMongoQuery(mongoField, true)
				case "":
					condQueries[j], err = cond.Condition.RefToMongoQuery(mongoField)
				default:
					err = fmt.Errorf("invalid field type for %q field: %s", cond.Field, cond.FieldType)
				}

				if err != nil {
					return nil, fmt.Errorf("invalid condition for %q field: %w", cond.Field, err)
				}

				continue
			}

			if infoName := pattern.GetEntityComponentInfoName(cond.Field); infoName != "" {
				mongoField := prefix + "component_infos." + infoName + ".value"

				switch cond.FieldType {
				case pattern.FieldTypeString:
					condQueries[j], err = cond.Condition.StringToMongoQuery(mongoField, true)
				case pattern.FieldTypeInt:
					condQueries[j], err = cond.Condition.IntToMongoQuery(mongoField, true)
				case pattern.FieldTypeBool:
					condQueries[j], err = cond.Condition.BoolToMongoQuery(mongoField)
				case pattern.FieldTypeStringArray:
					condQueries[j], err = cond.Condition.StringArrayToMongoQuery(mongoField, true)
				case "":
					condQueries[j], err = cond.Condition.RefToMongoQuery(mongoField)
				default:
					err = fmt.Errorf("invalid field type for %q field: %s", cond.Field, cond.FieldType)
				}

				if err != nil {
					return nil, fmt.Errorf("invalid condition for %q field: %w", cond.Field, err)
				}

				continue
			}

			mongoField := prefix + cond.Field
			if _, ok := emptyEntity.GetStringField(cond.Field); ok {
				condQueries[j], err = cond.Condition.StringToMongoQuery(mongoField, false)
			} else if _, ok := emptyEntity.GetIntField(cond.Field); ok {
				condQueries[j], err = cond.Condition.IntToMongoQuery(mongoField, false)
			} else if _, ok := emptyEntity.GetTimeField(cond.Field); ok {
				condQueries[j], err = cond.Condition.TimeToMongoQuery(mongoField)
			} else {
				err = pattern.ErrUnsupportedField
			}
			if err != nil {
				return nil, fmt.Errorf("invalid condition for %q field: %w", cond.Field, err)
			}
		}

		groupQueries[i] = bson.M{"$and": condQueries}
	}

	return groupQueries, nil
}
