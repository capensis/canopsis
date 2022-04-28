package pattern

import (
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"strings"
)

type Event [][]FieldCondition

type EventRegexMatches struct {
	Connector     RegexMatches
	ConnectorName RegexMatches
	Component     RegexMatches
	Resource      RegexMatches
	Output        RegexMatches
	EventType     RegexMatches
	SourceType    RegexMatches
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
	case "event_type":
		m.EventType = matches
	case "source_type":
		m.SourceType = matches
	}
}

func (m *EventRegexMatches) SetInfoRegexMatches(fieldName string, matches RegexMatches) {
	m.ExtraInfos[fieldName] = matches
}

func (p Event) Match(event types.Event) (bool, EventRegexMatches, error) {
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
				infoVal := getEventExtraInfoVal(event, infoName)

				switch v.FieldType {
				case FieldTypeString:
					if s, err := getStringValue(infoVal); err == nil {
						matched, regexMatches, err = cond.MatchString(s)
						if matched {
							eventRegexMatches.SetInfoRegexMatches(infoName, regexMatches)
						}
					}
				case FieldTypeInt:
					if i, err := getIntValue(infoVal); err == nil {
						matched, err = cond.MatchInt(i)
					}
				case FieldTypeBool:
					if b, err := getBoolValue(infoVal); err == nil {
						matched, err = cond.MatchBool(b)
					}
				case FieldTypeStringArray:
					if a, err := getStringArrayValue(infoVal); err == nil {
						matched, err = cond.MatchStringArray(a)
					}
				default:
					matched, err = cond.MatchRef(infoVal)
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
				matched, regexMatches, err = cond.MatchString(str)
				if matched {
					eventRegexMatches.SetRegexMatches(f, regexMatches)
				}
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

			if str, ok := getEventStringField(emptyEvent, f); ok {
				_, _, err = cond.MatchString(str)
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
	case "event_type":
		return event.EventType, true
	case "source_type":
		return event.SourceType, true
	default:
		return "", false
	}
}

func getEventExtraInfoName(f string) string {
	if n := strings.TrimPrefix(f, "extra."); n != f {
		return n
	}

	return ""
}

func getEventExtraInfoVal(event types.Event, f string) interface{} {
	if v, ok := event.ExtraInfos[f]; ok {
		return v
	}

	return nil
}
