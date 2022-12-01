package contextgraph

import (
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/importcontextgraph"
)

const (
	statusPending = "pending"
	statusOngoing = "ongoing"
	statusFailed  = "failed"
	statusDone    = "done"
)

type ImportJob struct {
	ID       string    `bson:"_id" json:"_id"`
	Creation time.Time `bson:"creation" json:"creation"`
	Status   string    `bson:"status" json:"status"`
	Info     string    `bson:"info,omitempty" json:"info"`
	ExecTime string    `bson:"exec_time,omitempty" json:"exec_time"`
	Source   string    `bson:"source" json:"source"`

	Stats importcontextgraph.Stats `bson:"stats" json:"stats"`

	IsPartial bool `bson:"is_partial" json:"-"`
	IsOld     bool `bson:"is_old" json:"-"`
}

type ImportResponse struct {
	ID string `json:"_id"`
}

type ImportQuery struct {
	Source string `form:"source" binding:"required"`
}

// OldImportRequest is used only for swagger docs.
type OldImportRequest struct {
	Json struct {
		Cis   []importcontextgraph.OldConfigurationItem `json:"cis"`
		Links []importcontextgraph.Link                 `json:"links"`
	} `json:"json"`
}

// ImportRequest is used only for swagger docs.
type ImportRequest struct {
	Json struct {
		Cis   []importcontextgraph.ConfigurationItem `json:"cis"`
		Links []importcontextgraph.Link              `json:"links"`
	} `json:"json"`
}
