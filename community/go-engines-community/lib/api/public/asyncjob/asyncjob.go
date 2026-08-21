// Package asyncjob provides a generic AMQP-driven runner for jobs which are too long to be handled in an http request:
// the API stores the job, publishes a signal, and a Worker picks it up and runs the Handler registered for the job type.
package asyncjob

import (
	"context"
)

const (
	TypeEntityImport = "entity_import"
	TypeEntityExport = "entity_export"
)

// Message is the signal published on job creation.
// It only carries the job type and its id: everything else is expected to be stored by the API, so a handler loads the job details itself.
type Message struct {
	Type  string `json:"type"`
	JobID string `json:"job_id"`
}

// Handler runs an already created job of a single type.
type Handler interface {
	Handle(ctx context.Context, jobID string) error
}
