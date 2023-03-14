package link

import (
	"context"
	"errors"
)

var ErrNoRule = errors.New("rule not found")
var ErrNotMatchedAlarm = errors.New("alarms aren't matched to rule")

type Generator interface {
	GenerateForAlarms(ctx context.Context, alarmIds []string) (map[string]LinksByCategory, error)
	GenerateForEntities(ctx context.Context, entityIds []string) (map[string]LinksByCategory, error)
	GenerateCombinedForAlarmsByRule(ctx context.Context, ruleId string, alarmIds []string) ([]Link, error)
	Load(ctx context.Context) error
}
