package contextgraph

import (
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/importcontextgraph"
)

const (
	StatusPending = "pending"
	StatusOngoing = "ongoing"
	StatusFailed  = "failed"
	StatusDone    = "done"
)

type ImportJob struct {
	ID       string     `bson:"_id" json:"_id"`
	Creation time.Time  `bson:"creation" json:"creation"`
	LastPing *time.Time `bson:"last_ping,omitempty" json:"last_ping"`
	Status   string     `bson:"status" json:"status"`
	Info     string     `bson:"info,omitempty" json:"info"`
	ExecTime string     `bson:"exec_time,omitempty" json:"exec_time"`
	Source   string     `bson:"source" json:"source"`

	Stats importcontextgraph.Stats `bson:"stats" json:"stats"`

	IsPartial bool `bson:"is_partial" json:"-"`
}

type ImportResponse struct {
	ID string `json:"_id"`
}

type ImportQuery struct {
	Source string `form:"source" binding:"required"`
}
