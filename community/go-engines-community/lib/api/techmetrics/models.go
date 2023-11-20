package techmetrics

import libtime "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/time"

type ExportResponse struct {
	// Possible values.
	//   * `0` - None
	//   * `1` - Running
	//   * `2` - Succeeded
	//   * `3` - Failed
	//   * `4` - Disabled
	Status   int              `json:"status"`
	Created  *libtime.CpsTime `json:"created,omitempty" swaggertype:"integer"`
	Duration *int             `json:"duration,omitempty"`
}
