package docs

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/docs"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/response"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"github.com/rs/zerolog"
	"go.yaml.in/yaml/v3"
)

func GetHandler(responder response.Responder, logger zerolog.Logger, generatedSchemasContents [][]byte, contents [][]byte) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		mergedContent, err := docs.MergeContent(contents)
		if err != nil {
			responder.RespondError(w, r, err)

			return
		}

		mergedContent, err = docs.MergeSchemas(generatedSchemasContents, mergedContent)
		if err != nil {
			responder.RespondError(w, r, err)

			return
		}

		if info, ok := mergedContent["info"].(map[string]any); ok {
			buildInfo := canopsis.GetBuildInfo()

			info["version"] = buildInfo.Version
			mergedContent["info"] = info
		}

		b, err := yaml.Marshal(mergedContent)
		if err != nil {
			responder.RespondError(w, r, err)

			return
		}

		w.Header().Set("Content-Type", response.MIMEYAML)
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(b) //nolint:gosec

		if err != nil {
			logger.Err(err).Msg("failed to write response")
		}
	}
}
