package techmetrics

import "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"

type Settings struct {
	Enabled          bool   `bson:"enabled" json:"enabled"`
	DumpKeepInterval string `bson:"dumpkeepinterval" json:"-"`
}

type ExportResponse struct {
	// Possible values.
	//   * `0` - None
	//   * `1` - Running
	//   * `2` - Succeeded
	//   * `3` - Failed
	//   * `4` - Disabled
	Status   int               `json:"status"`
	Created  *datetime.CpsTime `json:"created,omitempty" swaggertype:"integer"`
	Duration *int              `json:"duration,omitempty"`
}
