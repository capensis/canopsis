package techmetrics

import "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"

type ExportResponse struct {
	Status   int            `json:"status"`
	Created  *types.CpsTime `json:"created,omitempty" swaggertype:"integer"`
	Duration *int           `json:"duration,omitempty"`
}
