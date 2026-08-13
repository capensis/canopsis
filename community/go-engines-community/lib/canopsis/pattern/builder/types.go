package builder

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
)

const (
	relativeTimeKindNewerThan = "newer_than"
	relativeTimeKindOlderThan = "older_than"
	relativeTimeKindBetween   = "between"
)

type Input struct {
	Error    string `json:"error,omitempty"`
	OffTopic bool   `json:"offtopic,omitempty"`
	Patterns []Spec `json:"patterns,omitempty"`
}

type Spec struct {
	Type   string            `json:"type"`
	Groups [][]ConditionSpec `json:"groups"`
}

type ConditionSpec struct {
	Field     string    `json:"field"`
	Key       string    `json:"key,omitempty"`
	FieldType string    `json:"field_type,omitempty"`
	Operator  string    `json:"operator"`
	Value     ValueSpec `json:"value"`
}

type ValueSpec struct {
	String  *string  `json:"string,omitempty"`
	Strings []string `json:"strings,omitempty"`
	Integer *int64   `json:"integer,omitempty"`
	Boolean *bool    `json:"boolean,omitempty"`

	Duration     *DurationSpec     `json:"duration,omitempty"`
	RelativeTime *RelativeTimeSpec `json:"relative_time,omitempty"`
	AbsoluteTime *AbsoluteTimeSpec `json:"absolute_time,omitempty"`
}

type DurationSpec struct {
	Value int64  `json:"value"`
	Unit  string `json:"unit"`
}

func (d DurationSpec) toDurationWithUnit() datetime.DurationWithUnit {
	return datetime.DurationWithUnit{Value: d.Value, Unit: d.Unit}
}

type RelativeTimeSpec struct {
	Kind string `json:"kind"`

	Duration *DurationSpec `json:"duration,omitempty"`
	From     *DurationSpec `json:"from,omitempty"`
	To       *DurationSpec `json:"to,omitempty"`
}

type AbsoluteTimeSpec struct {
	From int64 `json:"from"`
	To   int64 `json:"to"`
}

type Result struct {
	EventPattern          pattern.Event                 `json:"event_pattern,omitempty" bson:"event_pattern,omitempty"`
	AlarmPattern          pattern.Alarm                 `json:"alarm_pattern,omitempty" bson:"alarm_pattern,omitempty"`
	EntityPattern         pattern.Entity                `json:"entity_pattern,omitempty" bson:"entity_pattern,omitempty"`
	PbehaviorPattern      pattern.PbehaviorInfo         `json:"pbehavior_pattern,omitempty" bson:"pbehavior_pattern,omitempty"`
	WeatherServicePattern pattern.WeatherServicePattern `json:"weather_service_pattern,omitempty" bson:"weather_service_pattern,omitempty"`
}

func (r Result) IsZero() bool {
	return len(r.EventPattern) == 0 &&
		len(r.AlarmPattern) == 0 &&
		len(r.EntityPattern) == 0 &&
		len(r.PbehaviorPattern) == 0 &&
		len(r.WeatherServicePattern) == 0
}
