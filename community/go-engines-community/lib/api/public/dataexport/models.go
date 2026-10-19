package dataexport

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/asyncjob"
)

const DownloadNameTimeLayout = "2006-01-02T15-04-05-MST"

type Field struct {
	Label string `json:"label" binding:"required,max=10000"`
	Name  string `json:"name" binding:"required,max=10000"`
}

type Options struct {
	// Dir is the directory the export file is written to, created if it doesn't exist.
	Dir string
	// Name is the export file name, without extension: the format appends its own.
	Name string
	// Fields are the exported fields, in order.
	Fields []Field
	// Transform converts a value before it's written, eg. to format a timestamp.
	// It's called with the field and the raw value, and may be nil.
	Transform Transform
	// Enrich adds to a document the fields which aren't fetched from the database, eg. the links
	// of the resource. It's called with the decoded document, before its fields are selected,
	// and may be nil.
	Enrich Enrich
}

type Job interface {
	JobID() string
	Response() asyncjob.Response
	JobStatus() string
	DownloadName() string
	MimeType() (mimeType string, format string)
}

type Defaulter[T any] interface {
	*T
	SetDefaults()
}
