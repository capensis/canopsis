package pattern

import (
	"fmt"
	"strings"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"go.mongodb.org/mongo-driver/bson"
)

type Entity [][]FieldCondition

func (p Entity) Match(entity types.Entity) (bool, error) {
	for _, group := range p {
		matched := false

		for _, v := range group {
			f := v.Field
			cond := v.Condition
			var err error
			matched = false

			if infoName := getEntityInfoName(f); infoName != "" {
				infoVal := getEntityInfoVal(entity, infoName)

				switch v.FieldType {
				case FieldTypeString:
					if s, ok := infoVal.(string); ok {
						matched, _, err = cond.MatchString(s)
					}
				case FieldTypeInt:
					if i, ok := infoVal.(int); ok {
						matched, err = cond.MatchInt(i)
					}
				case FieldTypeBool:
					if b, ok := infoVal.(bool); ok {
						matched, err = cond.MatchBool(b)
					}
				case FieldTypeStringArray:
					if a, ok := infoVal.([]string); ok {
						matched, err = cond.MatchStringArray(a)
					}
				default:
					matched, err = cond.MatchRef(infoVal)
				}

				if err != nil {
					return false, fmt.Errorf("invalid condition for %q field: %w", f, err)
				}

				if !matched {
					break
				}

				continue
			}

			if str, ok := getEntityStringField(entity, f); ok {
				matched, _, err = cond.MatchString(str)
			} else if i, ok := getEntityIntField(entity, f); ok {
				matched, err = cond.MatchInt(i)
			} else if t, ok := getEntityTimeField(entity, f); ok {
				matched, err = cond.MatchTime(t)
			} else if a, ok := getEntityStringArrayField(entity, f); ok {
				matched, err = cond.MatchStringArray(a)
			} else {
				err = ErrUnsupportedField
			}

			if err != nil {
				return false, fmt.Errorf("invalid condition for %q field: %w", f, err)
			}

			if !matched {
				break
			}
		}

		if matched {
			return true, nil
		}
	}

	return false, nil
}

func (p Entity) Validate() bool {
	emptyEntity := types.Entity{}

	for _, group := range p {
		if len(group) == 0 {
			return false
		}

		for _, v := range group {
			f := v.Field
			cond := v.Condition
			var err error

			if infoName := getEntityInfoName(f); infoName != "" {
				switch v.FieldType {
				case FieldTypeString:
					_, _, err = cond.MatchString("")
				case FieldTypeInt:
					_, err = cond.MatchInt(0)
				case FieldTypeBool:
					_, err = cond.MatchBool(false)
				case FieldTypeStringArray:
					_, err = cond.MatchStringArray([]string{})
				default:
					_, err = cond.MatchRef(nil)
				}

				if err != nil {
					return false
				}

				continue
			}

			if str, ok := getEntityStringField(emptyEntity, f); ok {
				_, _, err = cond.MatchString(str)
			} else if i, ok := getEntityIntField(emptyEntity, f); ok {
				_, err = cond.MatchInt(i)
			} else if t, ok := getEntityTimeField(emptyEntity, f); ok {
				_, err = cond.MatchTime(t)
			} else if a, ok := getEntityStringArrayField(emptyEntity, f); ok {
				_, err = cond.MatchStringArray(a)
			} else {
				err = ErrUnsupportedField
			}

			if err != nil {
				return false
			}
		}
	}

	return true
}

func (p Entity) ToMongoQuery(prefix string) ([]bson.M, error) {
	pipeline := make([]bson.M, 0)
	if len(p) == 0 {
		return nil, nil
	}

	if prefix != "" {
		prefix += "."
	}

	groupQueries := make([]bson.M, len(p))
	var err error

	for i, group := range p {
		condQueries := make([]bson.M, len(group))
		for j, cond := range group {
			f := cond.Field

			if infoName := getEntityInfoName(f); infoName != "" {
				f = prefix + "infos." + infoName + ".val"

				condQueries[j], err = cond.Condition.ToMongoQuery(f)
				if err != nil {
					return nil, err
				}

				conds := getTypeMongoQuery(f, cond.FieldType)

				if len(conds) > 0 {
					conds = append(conds, condQueries[j])
					condQueries[j] = bson.M{"$and": conds}
				}

				continue
			}

			f = prefix + f
			condQueries[j], err = cond.Condition.ToMongoQuery(f)
			if err != nil {
				return nil, err
			}
		}

		groupQueries[i] = bson.M{"$and": condQueries}
	}

	pipeline = append(pipeline, bson.M{"$match": bson.M{"$or": groupQueries}})

	return pipeline, nil
}

func getEntityStringField(entity types.Entity, f string) (string, bool) {
	switch f {
	case "_id":
		return entity.ID, true
	case "name":
		return entity.Name, true
	case "category":
		return entity.Category, true
	default:
		return "", false
	}
}

func getEntityIntField(entity types.Entity, f string) (int, bool) {
	switch f {
	case "impact_level":
		return int(entity.ImpactLevel), true
	default:
		return 0, false
	}
}

func getEntityTimeField(entity types.Entity, field string) (time.Time, bool) {
	switch field {
	case "last_event_date":
		if entity.LastEventDate != nil {
			return entity.LastEventDate.Time, true
		}

		return time.Time{}, true
	default:
		return time.Time{}, false
	}
}

func getEntityStringArrayField(entity types.Entity, f string) ([]string, bool) {
	switch f {
	case "impact":
		return entity.Impacts, true
	case "depends":
		return entity.Depends, true
	default:
		return nil, false
	}
}

func getEntityInfoVal(entity types.Entity, f string) interface{} {
	if v, ok := entity.Infos[f]; ok {
		return v.Value
	}

	return nil
}

func getEntityInfoName(f string) string {
	if n := strings.TrimPrefix(f, "infos."); n != f {
		return n
	}

	return ""
}
