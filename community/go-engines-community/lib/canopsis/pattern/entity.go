package pattern

import (
	"fmt"
	"strings"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"go.mongodb.org/mongo-driver/bson"
)

type Entity [][]FieldCondition

type EntityRegexMatches struct {
	ID             RegexMatches
	Name           RegexMatches
	Category       RegexMatches
	Type           RegexMatches
	Infos          map[string]RegexMatches
	ComponentInfos map[string]RegexMatches
}

func NewEntityRegexMatches() EntityRegexMatches {
	return EntityRegexMatches{
		Infos:          make(map[string]RegexMatches),
		ComponentInfos: make(map[string]RegexMatches),
	}
}

func (m *EntityRegexMatches) SetRegexMatches(fieldName string, matches RegexMatches) {
	switch fieldName {
	case "_id":
		m.ID = matches
	case "name":
		m.Name = matches
	case "category":
		m.Category = matches
	case "type":
		m.Type = matches
	}
}

func (m *EntityRegexMatches) SetInfoRegexMatches(fieldName string, matches RegexMatches) {
	m.Infos[fieldName] = matches
}

func (m *EntityRegexMatches) SetComponentInfoRegexMatches(fieldName string, matches RegexMatches) {
	m.ComponentInfos[fieldName] = matches
}

func (p Entity) Match(entity types.Entity) (bool, EntityRegexMatches, error) {
	entityRegexMatches := NewEntityRegexMatches()

	if len(p) == 0 {
		return true, entityRegexMatches, nil
	}

	for _, group := range p {
		matched := false

		for _, v := range group {
			f := v.Field
			cond := v.Condition
			var err error
			matched = false

			var regexMatches map[string]string

			if infoName := getEntityInfoName(f); infoName != "" {
				infoVal, ok := getEntityInfoVal(entity, infoName)
				if v.FieldType == "" {
					matched, err = cond.MatchRef(infoVal)
				} else if ok {
					switch v.FieldType {
					case FieldTypeString:
						var s string
						if s, err = getStringValue(infoVal); err == nil {
							matched, regexMatches, err = cond.MatchString(s)
							if matched {
								entityRegexMatches.SetInfoRegexMatches(infoName, regexMatches)
							}
						}
					case FieldTypeInt:
						var i int64
						if i, err = getIntValue(infoVal); err == nil {
							matched, err = cond.MatchInt(i)
						}
					case FieldTypeBool:
						var b bool
						if b, err = getBoolValue(infoVal); err == nil {
							matched, err = cond.MatchBool(b)
						}
					case FieldTypeStringArray:
						var a []string
						if a, err = getStringArrayValue(infoVal); err == nil {
							matched, err = cond.MatchStringArray(a)
						}
					default:
						return false, entityRegexMatches, fmt.Errorf("invalid field type for %q field: %s", f, v.FieldType)
					}
				}

				if err != nil {
					return false, entityRegexMatches, fmt.Errorf("invalid condition for %q field: %w", f, err)
				}

				if !matched {
					break
				}

				continue
			}

			if infoName := getEntityComponentInfoName(f); infoName != "" {
				infoVal, ok := getEntityComponentInfoVal(entity, infoName)
				if v.FieldType == "" {
					matched, err = cond.MatchRef(infoVal)
				} else if ok {
					switch v.FieldType {
					case FieldTypeString:
						var s string
						if s, err = getStringValue(infoVal); err == nil {
							matched, regexMatches, err = cond.MatchString(s)
							if matched {
								entityRegexMatches.SetComponentInfoRegexMatches(infoName, regexMatches)
							}
						}
					case FieldTypeInt:
						var i int64
						if i, err = getIntValue(infoVal); err == nil {
							matched, err = cond.MatchInt(i)
						}
					case FieldTypeBool:
						var b bool
						if b, err = getBoolValue(infoVal); err == nil {
							matched, err = cond.MatchBool(b)
						}
					case FieldTypeStringArray:
						var a []string
						if a, err = getStringArrayValue(infoVal); err == nil {
							matched, err = cond.MatchStringArray(a)
						}
					default:
						return false, entityRegexMatches, fmt.Errorf("invalid field type for %q field: %s", f, v.FieldType)
					}
				}

				if err != nil {
					return false, entityRegexMatches, fmt.Errorf("invalid condition for %q field: %w", f, err)
				}

				if !matched {
					break
				}

				continue
			}

			if str, ok := getEntityStringField(entity, f); ok {
				matched, regexMatches, err = cond.MatchString(str)
				if matched {
					entityRegexMatches.SetRegexMatches(f, regexMatches)
				}
			} else if i, ok := getEntityIntField(entity, f); ok {
				matched, err = cond.MatchInt(i)
			} else if t, ok := getEntityTimeField(entity, f); ok {
				matched, err = cond.MatchTime(t)
			} else {
				err = ErrUnsupportedField
			}

			if err != nil {
				return false, entityRegexMatches, fmt.Errorf("invalid condition for %q field: %w", f, err)
			}

			if !matched {
				break
			}
		}

		if matched {
			return true, entityRegexMatches, nil
		}
	}

	return false, entityRegexMatches, nil
}

func (p Entity) Validate(forbiddenFields []string) bool {
	emptyEntity := types.Entity{}
	forbiddenFieldsMap := make(map[string]bool, len(forbiddenFields))
	for _, field := range forbiddenFields {
		forbiddenFieldsMap[field] = true
	}

	for _, group := range p {
		if len(group) == 0 {
			return false
		}

		for _, v := range group {
			f := v.Field
			cond := v.Condition
			var err error

			if isForbiddenEntityField(v, forbiddenFieldsMap) {
				return false
			}

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

			if infoName := getEntityComponentInfoName(f); infoName != "" {
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

func (p Entity) ToMongoQuery(prefix string) (bson.M, error) {
	groupQueries, err := p.getGroupMongoQueries(prefix)
	if err != nil || len(groupQueries) == 0 {
		return nil, err
	}
	return bson.M{"$or": groupQueries}, nil
}

func (p Entity) ToNegativeMongoQuery(prefix string) (bson.M, error) {
	groupQueries, err := p.getGroupMongoQueries(prefix)
	if err != nil || len(groupQueries) == 0 {
		return nil, err
	}
	return bson.M{"$nor": groupQueries}, nil
}

func (p Entity) HasField(field string) bool {
	for _, group := range p {
		for _, condition := range group {
			if condition.Field == field {
				return true
			}
		}
	}

	return false
}

func (p Entity) HasInfosField() bool {
	for _, group := range p {
		for _, condition := range group {
			if infoName := getEntityInfoName(condition.Field); infoName != "" {
				return true
			}
		}
	}

	return false
}

func (p Entity) HasComponentInfosField() bool {
	for _, group := range p {
		for _, condition := range group {
			if infoName := getEntityComponentInfoName(condition.Field); infoName != "" {
				return true
			}
		}
	}

	return false
}

func (p Entity) RemoveFields(fields []string) Entity {
	forbiddenFieldsMap := make(map[string]bool, len(fields))
	for _, field := range fields {
		forbiddenFieldsMap[field] = true
	}

	newGroups := make(Entity, 0, len(p))
	for _, group := range p {
		newGroup := make([]FieldCondition, 0, len(group))
		for _, condition := range group {
			if isForbiddenEntityField(condition, forbiddenFieldsMap) {
				continue
			}

			newGroup = append(newGroup, condition)
		}
		if len(newGroup) > 0 {
			newGroups = append(newGroups, newGroup)
		}
	}

	if len(newGroups) > 0 {
		return newGroups
	}

	return nil
}

func (p Entity) getGroupMongoQueries(prefix string) ([]bson.M, error) {
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
				f = prefix + "infos." + infoName + ".value"

				condQueries[j], err = cond.Condition.ToMongoQuery(f)
				if err != nil {
					return nil, fmt.Errorf("invalid condition for %q field: %w", f, err)
				}

				conds := getTypeMongoQuery(f, cond.FieldType)

				if len(conds) > 0 {
					conds = append(conds, condQueries[j])
					condQueries[j] = bson.M{"$and": conds}
				}

				continue
			}

			if infoName := getEntityComponentInfoName(f); infoName != "" {
				f = prefix + "component_infos." + infoName + ".value"

				condQueries[j], err = cond.Condition.ToMongoQuery(f)
				if err != nil {
					return nil, fmt.Errorf("invalid condition for %q field: %w", f, err)
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
				return nil, fmt.Errorf("invalid condition for %q field: %w", f, err)
			}
		}

		groupQueries[i] = bson.M{"$and": condQueries}
	}

	return groupQueries, nil
}

func (p Entity) ToSql(prefix string) (string, error) {
	if len(p) == 0 {
		return "", nil
	}

	if prefix != "" {
		prefix += "."
	}

	groupQueries := make([]string, len(p))
	var err error

	fieldMap := map[string]string{"_id": "custom_id"}

	for i, group := range p {
		condQueries := make([]string, len(group))
		for j, cond := range group {
			f := cond.Field
			if v, ok := fieldMap[f]; ok {
				f = v
			}

			if infoName := getEntityInfoName(f); infoName != "" {
				condQueries[j], err = cond.Condition.ToSqlJson("infos", infoName, cond.FieldType)
				if err != nil {
					return "", err
				}

				continue
			}

			if infoName := getEntityComponentInfoName(f); infoName != "" {
				condQueries[j], err = cond.Condition.ToSqlJson("component_infos", infoName, cond.FieldType)
				if err != nil {
					return "", err
				}

				continue
			}

			f = prefix + f
			condQueries[j], err = cond.Condition.ToSql(f)
			if err != nil {
				return "", err
			}
		}

		groupQueries[i] = fmt.Sprintf("(%s)", strings.Join(condQueries, " AND "))
	}

	return strings.Join(groupQueries, " OR "), nil
}

func getEntityStringField(entity types.Entity, f string) (string, bool) {
	switch f {
	case "_id":
		return entity.ID, true
	case "name":
		return entity.Name, true
	case "category":
		return entity.Category, true
	case "type":
		return entity.Type, true
	case "connector":
		return entity.Connector, true
	case "component":
		return entity.Component, true
	default:
		return "", false
	}
}

func getEntityIntField(entity types.Entity, f string) (int64, bool) {
	switch f {
	case "impact_level":
		return entity.ImpactLevel, true
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

func getEntityInfoVal(entity types.Entity, f string) (interface{}, bool) {
	if v, ok := entity.Infos[f]; ok {
		return v.Value, true
	}

	return nil, false
}

func getEntityInfoName(f string) string {
	if n := strings.TrimPrefix(f, "infos."); n != f {
		return n
	}

	return ""
}

func getEntityComponentInfoVal(entity types.Entity, f string) (interface{}, bool) {
	if v, ok := entity.ComponentInfos[f]; ok {
		return v.Value, true
	}

	return nil, false
}

func getEntityComponentInfoName(f string) string {
	if n := strings.TrimPrefix(f, "component_infos."); n != f {
		return n
	}

	return ""
}

func isForbiddenEntityField(condition FieldCondition, forbiddenFieldsMap map[string]bool) bool {
	return forbiddenFieldsMap[condition.Field] ||
		forbiddenFieldsMap["infos"] && strings.HasPrefix(condition.Field, "infos") ||
		forbiddenFieldsMap["component_infos"] && strings.HasPrefix(condition.Field, "component_infos")
}
