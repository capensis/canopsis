package api

import (
	"net/http"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/response"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
)

func infoHandler(responder response.Responder) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		bi := canopsis.GetBuildInfo()
		responder.Respond(w, r, http.StatusOK, map[string]any{
			"name":     "Canopsis Public API",
			"version":  Version,
			"build":    bi.Edition + " " + bi.Version + " " + bi.VcsRevision,
			"build_at": bi.Date.Format(time.RFC3339),
			"links": map[string]string{
				"docs": Prefix + "/docs",
			},
		})
	}
}
