package contextgraph

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"time"
)

const (
	actionDelete  = "delete"
	actionCreate  = "create"
	actionSet     = "set"
	actionUpdate  = "update"
	actionDisable = "disable"
	actionEnable  = "enable"

	statusPending = "pending"
	statusOngoing = "ongoing"
	statusFailed  = "failed"
	statusDone    = "done"
)

type JobStats struct{
	ExecTime time.Duration `bson:"-" json:"-"`
	Deleted  int64         `bson:"deleted" json:"deleted"`
	Updated  int64         `bson:"updated" json:"updated"`
}

type ImportJob struct {
	ID       string    `bson:"_id"`
	Creation time.Time `bson:"creation"`
	Status   string    `bson:"status"`
	Info     string    `bson:"info,omitempty"`
	ExecTime string    `bson:"exec_time,omitempty"`
	Stats    JobStats  `bson:"stats"`
}

type ImportResponse struct {
	ID string `json:"_id"`
}

type ConfigurationItem struct {
	ID               string                     `json:"_id" bson:"-"`
	Name             *string                    `json:"name" bson:"name,omitempty"`
	Depends          []string				    `json:"-" bson:"depends"`
	Impact			 []string				    `json:"-" bson:"impact"`
	EnableHistory    []string                   `json:"-" bson:"enable_history"`
	Measurements     []interface{}              `json:"measurements" bson:"measurements"`
	EntityPatterns   *pattern.EntityPatternList `bson:"entity_patterns,omitempty" json:"entity_patterns"`
	OutputTemplate   *string                    `bson:"output_template,omitempty" json:"output_template"`
	Infos            map[string]interface{}     `json:"infos" bson:"infos"`
	Type             *string                    `json:"type" bson:"type,omitempty"`
	Category         *string                    `json:"category" bson:"category,omitempty"`
	ImpactLevel      *int64                     `json:"impact_level" bson:"impact_level,omitempty"`
	Enabled          *bool                      `json:"enabled" bson:"enabled,omitempty"`
	Action           string                     `json:"action" bson:"-"`
	ActionProperties interface{}                `json:"action_properties" bson:"-"`
}

type Link struct {
	ID               string                 `json:"_id"`
	To               string                 `json:"to"`
	From             []string               `json:"from"`
	Infos            map[string]interface{} `json:"infos"`
	Action           string                 `json:"action"`
	ActionProperties interface{}            `json:"action_properties" bson:"-"`
}

// for swagger
type Request struct {
	Json struct{
		Cis   []ConfigurationItem `json:"cis"`
		Links []Link              `json:"links"`
	} `json:"json"`
}