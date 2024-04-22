package savedpattern

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
)

const (
	TypeAlarm          = "alarm"
	TypeEntity         = "entity"
	TypePbehavior      = "pbehavior"
	TypeWeatherService = "weather_service"
)

type SavedPattern struct {
	ID                    string                        `bson:"_id"`
	Title                 string                        `bson:"title"`
	Type                  string                        `bson:"type"`
	IsCorporate           bool                          `bson:"is_corporate"`
	AlarmPattern          pattern.Alarm                 `bson:"alarm_pattern,omitempty"`
	EntityPattern         pattern.Entity                `bson:"entity_pattern,omitempty"`
	PbehaviorPattern      pattern.PbehaviorInfo         `bson:"pbehavior_pattern,omitempty"`
	WeatherServicePattern pattern.WeatherServicePattern `bson:"weather_service_pattern,omitempty"`
	Author                string                        `bson:"author"`
	Created               datetime.CpsTime              `bson:"created,omitempty"`
	Updated               datetime.CpsTime              `bson:"updated,omitempty"`
}

type AlarmPatternFields struct {
	AlarmPattern pattern.Alarm `bson:"alarm_pattern" json:"alarm_pattern,omitempty"`

	CorporateAlarmPattern      string `bson:"corporate_alarm_pattern" json:"corporate_alarm_pattern,omitempty"`
	CorporateAlarmPatternTitle string `bson:"corporate_alarm_pattern_title" json:"corporate_alarm_pattern_title,omitempty"`
}

type EntityPatternFields struct {
	EntityPattern pattern.Entity `bson:"entity_pattern" json:"entity_pattern,omitempty"`

	CorporateEntityPattern      string `bson:"corporate_entity_pattern" json:"corporate_entity_pattern,omitempty"`
	CorporateEntityPatternTitle string `bson:"corporate_entity_pattern_title" json:"corporate_entity_pattern_title,omitempty"`
}

type PbehaviorPatternFields struct {
	PbehaviorPattern pattern.PbehaviorInfo `bson:"pbehavior_pattern" json:"pbehavior_pattern,omitempty"`

	CorporatePbehaviorPattern      string `bson:"corporate_pbehavior_pattern" json:"corporate_pbehavior_pattern,omitempty"`
	CorporatePbehaviorPatternTitle string `bson:"corporate_pbehavior_pattern_title" json:"corporate_pbehavior_pattern_title,omitempty"`
}

type WeatherServicePatternFields struct {
	WeatherServicePattern pattern.WeatherServicePattern `bson:"weather_service_pattern" json:"weather_service_pattern,omitempty"`

	CorporateWeatherServicePattern      string `bson:"corporate_weather_service_pattern" json:"corporate_weather_service_pattern,omitempty"`
	CorporateWeatherServicePatternTitle string `bson:"corporate_weather_service_pattern_title" json:"corporate_weather_service_pattern_title,omitempty"`
}
