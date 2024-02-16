package pattern

import (
	"fmt"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type Event [][]FieldCondition

type EventRegexMatches struct {
	Connector     RegexMatches
	ConnectorName RegexMatches
	Component     RegexMatches
	Resource      RegexMatches
	Output        RegexMatches
	LongOutput    RegexMatches
	EventType     RegexMatches
	SourceType    RegexMatches
	Author        RegexMatches
	Initiator     RegexMatches
	ExtraInfos    map[string]RegexMatches
}

func NewEventRegexMatches() EventRegexMatches {
	return EventRegexMatches{ExtraInfos: make(map[string]RegexMatches)}
}

func (m *EventRegexMatches) SetRegexMatches(fieldName string, matches RegexMatches) {
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
	case "author":
		m.Author = matches
	case "initiator":
		m.Initiator = matches
	}
}

func (m *EventRegexMatches) SetInfoRegexMatches(fieldName string, matches RegexMatches) {
	m.ExtraInfos[fieldName] = matches
}

func (p Event) Match(event types.Event) (bool, error) {
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

			if infoName := getEventExtraInfoName(f); infoName != "" {
				infoVal, ok := getEventExtraInfoVal(event, infoName)
				if v.FieldType == "" {
					matched, err = cond.MatchRef(infoVal)
				} else if ok {
					switch v.FieldType {
					case FieldTypeString:
						var s string
						if s, err = getStringValue(infoVal); err == nil {
							matched, err = cond.MatchString(s)
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

			if str, ok := getEventStringField(event, f); ok {
				matched, err = cond.MatchString(str)
			} else if i, ok := getEventIntField(event, f); ok {
				matched, err = cond.MatchInt(i)
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

func (p Event) MatchWithRegexMatches(event types.Event) (bool, EventRegexMatches, error) {
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

			if infoName := getEventExtraInfoName(f); infoName != "" {
				infoVal, ok := getEventExtraInfoVal(event, infoName)
				if v.FieldType == "" {
					matched, err = cond.MatchRef(infoVal)
				} else if ok {
					switch v.FieldType {
					case FieldTypeString:
						var s string
						if s, err = getStringValue(infoVal); err == nil {
							matched, regexMatches, err = cond.MatchStringWithRegexpMatches(s)
							if matched {
								eventRegexMatches.SetInfoRegexMatches(infoName, regexMatches)
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

			if str, ok := getEventStringField(event, f); ok {
				matched, regexMatches, err = cond.MatchStringWithRegexpMatches(str)
				if matched {
					eventRegexMatches.SetRegexMatches(f, regexMatches)
				}
			} else if i, ok := getEventIntField(event, f); ok {
				matched, err = cond.MatchInt(i)
			} else {
				err = ErrUnsupportedField
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

func (p Event) Validate() bool {
	emptyEvent := types.Event{}

	for _, group := range p {
		if len(group) == 0 {
			return false
		}

		for _, v := range group {
			f := v.Field
			cond := v.Condition
			var err error

			if infoName := getEventExtraInfoName(f); infoName != "" {
				switch v.FieldType {
				case FieldTypeString:
					_, err = cond.MatchString("")
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

			if str, ok := getEventStringField(emptyEvent, f); ok {
				_, err = cond.MatchString(str)
			} else if i, ok := getEventIntField(emptyEvent, f); ok {
				_, err = cond.MatchInt(i)
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

func (p Event) HasField(field string) bool {
	for _, group := range p {
		for _, condition := range group {
			if condition.Field == field {
				return true
			}
		}
	}

	return false
}

func (p Event) HasExtraInfosField() bool {
	for _, group := range p {
		for _, condition := range group {
			if infoName := getEventExtraInfoName(condition.Field); infoName != "" {
				return true
			}
		}
	}

	return false
}

func getEventStringField(event types.Event, f string) (string, bool) {
	switch f {
	case "connector":
		return event.Connector, true
	case "connector_name":
		return event.ConnectorName, true
	case "component":
		return event.Component, true
	case "resource":
		return event.Resource, true
	case "output":
		return event.Output, true
	case "long_output":
		return event.LongOutput, true
	case "event_type":
		return event.EventType, true
	case "source_type":
		return event.SourceType, true
	case "author":
		return event.Author, true
	case "initiator":
		return event.Initiator, true
	default:
		return "", false
	}
}

func getEventIntField(event types.Event, f string) (int64, bool) {
	switch f {
	case "state":
		return int64(event.State), true
	default:
		return 0, false
	}
}

func getEventExtraInfoName(f string) string {
	if n := strings.TrimPrefix(f, "extra."); n != f {
		return n
	}

	return ""
}

func getEventExtraInfoVal(event types.Event, f string) (interface{}, bool) {
	if v, ok := event.ExtraInfos[f]; ok {
		return v, true
	}

	return nil, false
}
