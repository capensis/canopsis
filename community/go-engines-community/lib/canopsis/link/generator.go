package link

import (
	"context"
	"errors"
)

var ErrNoRule = errors.New("rule not found")
var ErrNotMatchedAlarm = errors.New("alarms aren't matched to rule")

type Generator interface {
	GenerateForAlarms(ctx context.Context, alarmIds []string, user User) (map[string]LinksByCategory, error)
	GenerateForEntities(ctx context.Context, entityIds []string, user User) (map[string]LinksByCategory, error)
	GenerateCombinedForAlarmsByRule(ctx context.Context, ruleId string, alarmIds []string, user User) ([]Link, error)
	Load(ctx context.Context) error
}

type User struct {
	Email      string `bson:"email"`
	Username   string `bson:"username"`
	Firstname  string `bson:"firstname"`
	Lastname   string `bson:"lastname"`
	ExternalID string `bson:"external_id"`
	Source     string `bson:"source"`
	Role       string `bson:"role"`
}
