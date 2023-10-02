package link

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

var ErrNoRule = errors.New("rule not found")
var ErrNotMatchedAlarm = errors.New("alarms aren't matched to rule")

type Generator interface {
	GenerateForAlarm(ctx context.Context, alarm types.Alarm, entity types.Entity, user User) (LinksByCategory, error)
	GenerateForAlarms(ctx context.Context, alarmIds []string, user User) (map[string]LinksByCategory, error)
	GenerateForEntities(ctx context.Context, entityIds []string, user User) (map[string]LinksByCategory, error)
	GenerateCombinedForAlarmsByRule(ctx context.Context, ruleId string, alarmIds []string, user User) ([]Link, error)
	Load(ctx context.Context) error
}

type User struct {
	Email      string   `bson:"email"`
	Username   string   `bson:"username"`
	Firstname  string   `bson:"firstname"`
	Lastname   string   `bson:"lastname"`
	ExternalID string   `bson:"external_id"`
	Source     string   `bson:"source"`
	Roles      []string `bson:"roles"`
}
