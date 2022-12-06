package contextgraph

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/importcontextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/gin-gonic/gin"
)

type API interface {
	ImportAll(c *gin.Context)
	ImportPartial(c *gin.Context)
	Status(c *gin.Context)
}

type EventPublisher interface {
	SendImportResultEvent(ctx context.Context, uuid string, execTime time.Duration, state types.CpsNumber) error
	SendPerfDataEvent(ctx context.Context, uuid string, stats importcontextgraph.Stats, state types.CpsNumber) error
}

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
