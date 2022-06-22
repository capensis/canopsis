package contextgraph

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/importcontextgraph"
	"time"
)

type StatusReporter interface {
	GetStatus(ctx context.Context, id string) (ImportJob, error)
	ReportCreate(ctx context.Context, job *ImportJob) error
	ReportOngoing(ctx context.Context, job ImportJob) error
	ReportDone(ctx context.Context, job ImportJob, stats importcontextgraph.Stats) error
	ReportError(ctx context.Context, job ImportJob, execDuration time.Duration, err error) error
	Clean(ctx context.Context, interval time.Duration) error
}

type ImportWorker interface {
	Run(ctx context.Context)
}

type JobQueue interface {
	Push(job ImportJob)
	Pop() ImportJob
}
