package contextgraph

import (
	"context"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"time"
)

type EventPublisher interface {
	SendImportResultEvent(uuid string, execTime time.Duration, state types.CpsNumber) error
	SendPerfDataEvent(uuid string, stats JobStats, state types.CpsNumber) error
	SendUpdateEntityServiceEvent(serviceId string) error
}

type StatusReporter interface {
	GetStatus(id string) (ImportJob, error)
	ReportCreate(job *ImportJob) error
	ReportOngoing(job ImportJob) error
	ReportDone(job ImportJob, stats JobStats) error
	ReportError(job ImportJob, execDuration time.Duration, err error) error
}

type ImportWorker interface {
	Run(ctx context.Context)
}

type JobQueue interface {
	Push(job ImportJob)
	Pop() ImportJob
}