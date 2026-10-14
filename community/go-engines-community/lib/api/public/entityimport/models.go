package entityimport

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/asyncjob"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/importcontextgraph"
)

type Job struct {
	asyncjob.Base `bson:",inline"`

	Source string `bson:"source" json:"source"`
	Stats  Stats  `bson:"stats" json:"stats"`
}

type Stats struct {
	Updated int64 `bson:"updated" json:"updated"`
	Deleted int64 `bson:"deleted" json:"deleted"`
}

// RequestItem is one item of the body of an import request, declared for the swagger generator only.
type RequestItem importcontextgraph.EntityConfiguration
