package contextgraph

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/importcontextgraph"
	"time"
)

const (
	statusPending = "pending"
	statusOngoing = "ongoing"
	statusFailed  = "failed"
	statusDone    = "done"
)

type ImportJob struct {
	ID       string                   `bson:"_id" json:"_id"`
	Creation time.Time                `bson:"creation" json:"creation"`
	Status   string                   `bson:"status" json:"status"`
	Info     string                   `bson:"info,omitempty" json:"info"`
	ExecTime string                   `bson:"exec_time,omitempty" json:"exec_time"`
	Stats    importcontextgraph.Stats `bson:"stats" json:"stats"`
}

type ImportResponse struct {
	ID string `json:"_id"`
}

// Request is used only for swagger docs.
type Request struct {
	Json struct {
		Cis   []importcontextgraph.ConfigurationItem `json:"cis"`
		Links []importcontextgraph.Link              `json:"links"`
	} `json:"json"`
}
