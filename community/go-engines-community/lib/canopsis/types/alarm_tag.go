package types

import (
	"fmt"
	"strings"
)

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

	return fmt.Sprintf("%s: %s", k, v)
}
