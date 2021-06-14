package account

import "git.canopsis.net/canopsis/go-engines/lib/canopsis/types"

type User struct {
	ID         string `bson:"_id" json:"_id"`
	Name       string `bson:"crecord_name" json:"crecord_name"`
	Lastname   string `bson:"lastname" json:"lastname"`
	Firstname  string `bson:"firstname" json:"firstname"`
	Email      string `bson:"mail" json:"mail"`
	IsEnabled  bool   `bson:"enable" json:"enable"`
	ExternalID string `bson:"external_id" json:"external_id"`
	Source     string `bson:"source" json:"source"`
	AuthApiKey string `bson:"authkey" json:"authkey"`
	Role       string `bson:"role" json:"role"`
	Rights     map[string]struct {
		Bitmask int `bson:"checksum" json:"checksum"`
	} `bson:"rights" json:"rights"`
	UILanguage           string          `bson:"ui_language" json:"ui_language"`
	GroupsNavigationType string          `bson:"groupsNavigationType" json:"groupsNavigationType"`
	Tours                map[string]bool `bson:"tours" json:"tours"`
	PausedExecutions     []struct {
		ID              string        `bson:"_id" json:"_id"`
		AlarmName       string        `bson:"alarm_name" json:"alarm_name"`
		InstructionName string        `bson:"instruction_name" json:"instruction_name"`
		Paused          types.CpsTime `bson:"paused" json:"paused" swaggertype:"integer"`
	} `bson:"paused_executions,omitempty" json:"paused_executions,omitempty"`
	DefaultView string `bson:"defaultview" json:"defaultview"`
}
