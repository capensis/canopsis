package contextgraph

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"time"
)

type EventPublisher interface {
	SendImportResultEvent(uuid string, execTime time.Duration, state types.CpsNumber) error
	SendPerfDataEvent(uuid string, stats JobStats, state types.CpsNumber) error
	SendUpdateEntityServiceEvent(serviceId string) error
}

type StatusReporter interface {
	GetStatus(ctx context.Context, id string) (ImportJob, error)
	ReportCreate(ctx context.Context, job *ImportJob) error
	ReportOngoing(ctx context.Context, job ImportJob) error
	ReportDone(ctx context.Context, job ImportJob, stats JobStats) error
	ReportError(ctx context.Context, job ImportJob, execDuration time.Duration, err error) error
}

type ImportWorker interface {
	Run(ctx context.Context)
}

type JobQueue interface {
	Push(job ImportJob)
	Pop() ImportJob
}
