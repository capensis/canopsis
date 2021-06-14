package sessionstats

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/session/stats"
)

type pingRequest struct {
	Visible *bool    `json:"visible" binding:"required"`
	Path    []string `json:"path" binding:"required,notblank,max=2"`
}

type changePathRequest struct {
	Path []string `json:"path" binding:"required,notblank,max=2"`
}

type listRequest struct {
	IsActive      *bool          `form:"active"`
	Usernames     []string       `form:"usernames[]"`
	StartedAfter  *types.CpsTime `form:"started_after"`
	StoppedBefore *types.CpsTime `form:"stopped_before"`
}

type listResponse struct {
	Stats []stats.Stats `json:"sessions"`
}
