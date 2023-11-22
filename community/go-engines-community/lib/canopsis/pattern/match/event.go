package match

import (
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type EventRegexMatches struct {
	Connector     pattern.RegexMatches
	ConnectorName pattern.RegexMatches
	Component     pattern.RegexMatches
	Resource      pattern.RegexMatches
	Output        pattern.RegexMatches
	LongOutput    pattern.RegexMatches
	EventType     pattern.RegexMatches
	SourceType    pattern.RegexMatches
	ExtraInfos    map[string]pattern.RegexMatches
}

func NewEventRegexMatches() EventRegexMatches {
	return EventRegexMatches{ExtraInfos: make(map[string]pattern.RegexMatches)}
}

func (m *EventRegexMatches) SetRegexMatches(fieldName string, matches pattern.RegexMatches) {
	switch fieldName {
	case "connector":
		m.Connector = matches
	case "connector_name":
		m.ConnectorName = matches
	case "component":
		m.Component = matches
	case "resource":
		m.Resource = matches
	case "output":
		m.Output = matches
	case "long_output":
		m.LongOutput = matches
	case "event_type":
		m.EventType = matches
	case "source_type":
		m.SourceType = matches
	}
}

func (m *EventRegexMatches) SetInfoRegexMatches(fieldName string, matches pattern.RegexMatches) {
	m.ExtraInfos[fieldName] = matches
}

func ValidateEventPattern(p pattern.Event) bool {
	emptyEvent := types.Event{}

	for _, group := range p {
		if len(group) == 0 {
			return false
		}

		for _, v := range group {
			f := v.Field
			cond := v.Condition
			var err error

			if infoName := pattern.GetEventExtraInfoName(f); infoName != "" {
				switch v.FieldType {
				case pattern.FieldTypeString:
					_, err = cond.MatchString("")
				case pattern.FieldTypeInt:
					_, err = cond.MatchInt(0)
				case pattern.FieldTypeBool:
					_, err = cond.MatchBool(false)
				case pattern.FieldTypeStringArray:
					_, err = cond.MatchStringArray([]string{})
				default:
					_, err = cond.MatchRef(nil)
				}

				if err != nil {
					return false
				}

				continue
			}

			if str, ok := emptyEvent.GetStringField(f); ok {
				_, err = cond.MatchString(str)
			} else if i, ok := emptyEvent.GetIntField(f); ok {
				_, err = cond.MatchInt(i)
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

func MatchEventPattern(p pattern.Event, event *types.Event) (bool, error) {
	if len(p) == 0 {
		return true, nil
	}

	for _, group := range p {
		matched := false

		for _, v := range group {
			f := v.Field
			cond := v.Condition
			var err error
			matched = false

			if infoName := pattern.GetEventExtraInfoName(f); infoName != "" {
				infoVal, ok := event.GetExtraInfoVal(infoName)
				if v.FieldType == "" {
					matched, err = cond.MatchRef(infoVal)
				} else if ok {
					switch v.FieldType {
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
						return false, fmt.Errorf("invalid field type for %q field: %s", f, v.FieldType)
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

			if str, ok := event.GetStringField(f); ok {
				matched, err = cond.MatchString(str)
			} else if i, ok := event.GetIntField(f); ok {
				matched, err = cond.MatchInt(i)
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

func MatchEventPatternWithRegexMatches(p pattern.Event, event *types.Event) (bool, EventRegexMatches, error) {
	eventRegexMatches := NewEventRegexMatches()

	if len(p) == 0 {
		return true, eventRegexMatches, nil
	}

	for _, group := range p {
		matched := false

		for _, v := range group {
			f := v.Field
			cond := v.Condition
			var err error
			matched = false

			var regexMatches map[string]string

			if infoName := pattern.GetEventExtraInfoName(f); infoName != "" {
				infoVal, ok := event.GetExtraInfoVal(infoName)
				if v.FieldType == "" {
					matched, err = cond.MatchRef(infoVal)
				} else if ok {
					switch v.FieldType {
					case pattern.FieldTypeString:
						var s string
						if s, err = pattern.GetStringValue(infoVal); err == nil {
							matched, regexMatches, err = cond.MatchStringWithRegexpMatches(s)
							if matched {
								eventRegexMatches.SetInfoRegexMatches(infoName, regexMatches)
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
						return false, eventRegexMatches, fmt.Errorf("invalid field type for %q field: %s", f, v.FieldType)
					}
				}

				if err != nil {
					return false, eventRegexMatches, fmt.Errorf("invalid condition for %q field: %w", f, err)
				}

				if !matched {
					break
				}

				continue
			}

			if str, ok := event.GetStringField(f); ok {
				matched, regexMatches, err = cond.MatchStringWithRegexpMatches(str)
				if matched {
					eventRegexMatches.SetRegexMatches(f, regexMatches)
				}
			} else if i, ok := event.GetIntField(f); ok {
				matched, err = cond.MatchInt(i)
			} else {
				err = pattern.ErrUnsupportedField
			}

			if err != nil {
				return false, eventRegexMatches, fmt.Errorf("invalid condition for %q field: %w", f, err)
			}

			if !matched {
				break
			}
		}

		if matched {
			return true, eventRegexMatches, nil
		}
	}

	return false, eventRegexMatches, nil
}
