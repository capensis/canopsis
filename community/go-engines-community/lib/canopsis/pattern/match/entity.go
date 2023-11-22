package match

import (
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type EntityRegexMatches struct {
	ID             pattern.RegexMatches
	Name           pattern.RegexMatches
	Category       pattern.RegexMatches
	Type           pattern.RegexMatches
	Infos          map[string]pattern.RegexMatches
	ComponentInfos map[string]pattern.RegexMatches
}

func NewEntityRegexMatches() EntityRegexMatches {
	return EntityRegexMatches{
		Infos:          make(map[string]pattern.RegexMatches),
		ComponentInfos: make(map[string]pattern.RegexMatches),
	}
}

func (m *EntityRegexMatches) SetRegexMatches(fieldName string, matches pattern.RegexMatches) {
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

func (m *EntityRegexMatches) SetInfoRegexMatches(fieldName string, matches pattern.RegexMatches) {
	m.Infos[fieldName] = matches
}

func (m *EntityRegexMatches) SetComponentInfoRegexMatches(fieldName string, matches pattern.RegexMatches) {
	m.ComponentInfos[fieldName] = matches
}

func ValidateEntityPattern(p pattern.Entity, forbiddenFields []string) bool {
	emptyEntity := &types.Entity{}

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

			if pattern.IsForbiddenEntityField(v, forbiddenFieldsMap) {
				return false
			}

			if infoName := pattern.GetEntityInfoName(f); infoName != "" {
				switch v.FieldType {
				case pattern.FieldTypeString:
					_, err = cond.MatchString("")
				case pattern.FieldTypeInt:
					_, err = cond.MatchInt(0)
				case pattern.FieldTypeBool:
					_, err = cond.MatchBool(false)
				case pattern.FieldTypeStringArray:
					_, err = cond.MatchStringArray([]string{})
				case "":
					_, err = cond.MatchRef(nil)
				default:
					return false
				}

				if err != nil {
					return false
				}

				continue
			}

			if infoName := pattern.GetEntityComponentInfoName(f); infoName != "" {
				switch v.FieldType {
				case pattern.FieldTypeString:
					_, err = cond.MatchString("")
				case pattern.FieldTypeInt:
					_, err = cond.MatchInt(0)
				case pattern.FieldTypeBool:
					_, err = cond.MatchBool(false)
				case pattern.FieldTypeStringArray:
					_, err = cond.MatchStringArray([]string{})
				case "":
					_, err = cond.MatchRef(nil)
				default:
					return false
				}

				if err != nil {
					return false
				}

				continue
			}

			if str, ok := emptyEntity.GetStringField(f); ok {
				_, err = cond.MatchString(str)
			} else if i, ok := emptyEntity.GetIntField(f); ok {
				_, err = cond.MatchInt(i)
			} else if t, ok := emptyEntity.GetTimeField(f); ok {
				_, err = cond.MatchTime(t)
			} else {
				err = pattern.ErrUnsupportedField
			}

			if err != nil {
				return false
			}
		}
	}

	return true
}

func MatchEntityPattern(p pattern.Entity, entity *types.Entity) (bool, error) {
	if len(p) == 0 {
		return true, nil
	}

	for idx := range p {
		matched := false

		for jdx := range p[idx] {
			f := p[idx][jdx].Field
			cond := p[idx][jdx].Condition
			var err error
			matched = false

			if infoName := pattern.GetEntityInfoName(f); infoName != "" {
				infoVal, ok := getEntityInfoVal(entity, infoName)
				if p[idx][jdx].FieldType == "" {
					matched, err = cond.MatchRef(infoVal)
				} else if ok {
					switch p[idx][jdx].FieldType {
					case pattern.FieldTypeString:
						var s string
						if s, err = pattern.GetStringValue(infoVal); err == nil {
							matched, err = cond.MatchString(s)
						}
					case pattern.FieldTypeInt:
						var i int64
						if i, err = pattern.GetIntValue(infoVal); err == nil {
							matched, err = cond.MatchInt(i)
						}
					case pattern.FieldTypeBool:
						var b bool
						if b, err = pattern.GetBoolValue(infoVal); err == nil {
							matched, err = cond.MatchBool(b)
						}
					case pattern.FieldTypeStringArray:
						var a []string
						if a, err = pattern.GetStringArrayValue(infoVal); err == nil {
							matched, err = cond.MatchStringArray(a)
						}
					default:
						return false, fmt.Errorf("invalid field type for %q field: %s", f, p[idx][jdx].FieldType)
					}
				}

				if err != nil {
					return false, fmt.Errorf("invalid condition for %q field: %w", f, err)
				}

				if !matched {
					break
				}

				continue
			}

			if infoName := pattern.GetEntityComponentInfoName(f); infoName != "" {
				infoVal, ok := getEntityComponentInfoVal(entity, infoName)
				if p[idx][jdx].FieldType == "" {
					matched, err = cond.MatchRef(infoVal)
				} else if ok {
					switch p[idx][jdx].FieldType {
					case pattern.FieldTypeString:
						var s string
						if s, err = pattern.GetStringValue(infoVal); err == nil {
							matched, err = cond.MatchString(s)
						}
					case pattern.FieldTypeInt:
						var i int64
						if i, err = pattern.GetIntValue(infoVal); err == nil {
							matched, err = cond.MatchInt(i)
						}
					case pattern.FieldTypeBool:
						var b bool
						if b, err = pattern.GetBoolValue(infoVal); err == nil {
							matched, err = cond.MatchBool(b)
						}
					case pattern.FieldTypeStringArray:
						var a []string
						if a, err = pattern.GetStringArrayValue(infoVal); err == nil {
							matched, err = cond.MatchStringArray(a)
						}
					default:
						return false, fmt.Errorf("invalid field type for %q field: %s", f, p[idx][jdx].FieldType)
					}
				}

				if err != nil {
					return false, fmt.Errorf("invalid condition for %q field: %w", f, err)
				}

				if !matched {
					break
				}

				continue
			}

			if str, ok := entity.GetStringField(f); ok {
				matched, err = cond.MatchString(str)
			} else if i, ok := entity.GetIntField(f); ok {
				matched, err = cond.MatchInt(i)
			} else if t, ok := entity.GetTimeField(f); ok {
				matched, err = cond.MatchTime(t)
			} else {
				err = pattern.ErrUnsupportedField
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

func MatchEntityPatternWithRegexMatches(p pattern.Entity, entity *types.Entity) (bool, EntityRegexMatches, error) {
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

			if infoName := pattern.GetEntityInfoName(f); infoName != "" {
				infoVal, ok := getEntityInfoVal(entity, infoName)
				if v.FieldType == "" {
					matched, err = cond.MatchRef(infoVal)
				} else if ok {
					switch v.FieldType {
					case pattern.FieldTypeString:
						var s string
						if s, err = pattern.GetStringValue(infoVal); err == nil {
							matched, regexMatches, err = cond.MatchStringWithRegexpMatches(s)
							if matched {
								entityRegexMatches.SetInfoRegexMatches(infoName, regexMatches)
							}
						}
					case pattern.FieldTypeInt:
						var i int64
						if i, err = pattern.GetIntValue(infoVal); err == nil {
							matched, err = cond.MatchInt(i)
						}
					case pattern.FieldTypeBool:
						var b bool
						if b, err = pattern.GetBoolValue(infoVal); err == nil {
							matched, err = cond.MatchBool(b)
						}
					case pattern.FieldTypeStringArray:
						var a []string
						if a, err = pattern.GetStringArrayValue(infoVal); err == nil {
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

			if infoName := pattern.GetEntityComponentInfoName(f); infoName != "" {
				infoVal, ok := getEntityComponentInfoVal(entity, infoName)
				if v.FieldType == "" {
					matched, err = cond.MatchRef(infoVal)
				} else if ok {
					switch v.FieldType {
					case pattern.FieldTypeString:
						var s string
						if s, err = pattern.GetStringValue(infoVal); err == nil {
							matched, regexMatches, err = cond.MatchStringWithRegexpMatches(s)
							if matched {
								entityRegexMatches.SetComponentInfoRegexMatches(infoName, regexMatches)
							}
						}
					case pattern.FieldTypeInt:
						var i int64
						if i, err = pattern.GetIntValue(infoVal); err == nil {
							matched, err = cond.MatchInt(i)
						}
					case pattern.FieldTypeBool:
						var b bool
						if b, err = pattern.GetBoolValue(infoVal); err == nil {
							matched, err = cond.MatchBool(b)
						}
					case pattern.FieldTypeStringArray:
						var a []string
						if a, err = pattern.GetStringArrayValue(infoVal); err == nil {
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

			if str, ok := entity.GetStringField(f); ok {
				matched, regexMatches, err = cond.MatchStringWithRegexpMatches(str)
				if matched {
					entityRegexMatches.SetRegexMatches(f, regexMatches)
				}
			} else if i, ok := entity.GetIntField(f); ok {
				matched, err = cond.MatchInt(i)
			} else if t, ok := entity.GetTimeField(f); ok {
				matched, err = cond.MatchTime(t)
			} else {
				err = pattern.ErrUnsupportedField
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

func getEntityInfoVal(entity *types.Entity, f string) (interface{}, bool) {
	if v, ok := entity.Infos[f]; ok {
		return v.Value, true
	}

	return nil, false
}

func getEntityComponentInfoVal(entity *types.Entity, f string) (interface{}, bool) {
	if v, ok := entity.ComponentInfos[f]; ok {
		return v.Value, true
	}

	return nil, false
}
