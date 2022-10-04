package types

import (
	"strings"
)

type AlarmTag struct {
	ID      string  `bson:"_id" json:"_id"`
	Value   string  `bson:"value" json:"value"`
	Label   string  `bson:"label" json:"label"`
	Color   string  `bson:"color" json:"color"`
	Created CpsTime `bson:"created" json:"created"`
}

func TransformEventTags(eventTags map[string]string) []string {
	alarmTags := make([]string, 0, len(eventTags))
	for k, v := range eventTags {
		t := TransformEventTag(k, v)
		if t != "" {
			alarmTags = append(alarmTags, t)
		}
	}

	return alarmTags
}

func TransformEventTag(k, v string) string {
	k = strings.TrimSpace(k)
	v = strings.TrimSpace(v)

	if k == "" || v == "" {
		return k
	}

	return k + ": " + v
}
