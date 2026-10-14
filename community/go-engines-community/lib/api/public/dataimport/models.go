package dataimport

import (
	"path/filepath"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/asyncjob"
)

// Filepath returns the path of the file uploaded for a job, in dir.
// The API and the handler of a job type are the only ones to know it.
func Filepath(dir, jobID string) string {
	return filepath.Join(dir, jobID+".json")
}

type Job interface {
	JobID() string
	Response() asyncjob.Response
}

type Query struct {
	Source string `form:"source" binding:"required"`
	Author string `json:"author" swaggerignore:"true"`
}
