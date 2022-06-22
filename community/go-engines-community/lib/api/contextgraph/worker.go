package contextgraph

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/importcontextgraph"
	"github.com/rs/zerolog"
	"time"
)

const (
	queueCheckTickInterval  = time.Second
	reportCleanTickInterval = time.Hour
	reportCleanInterval     = 24 * time.Hour
)

type worker struct {
	importQueue JobQueue
	reporter    StatusReporter
	logger      zerolog.Logger
	filePattern string
	worker      importcontextgraph.Worker
}

func NewImportWorker(
	conf config.CanopsisConf,
	reporter StatusReporter,
	queue JobQueue,
	importWorker importcontextgraph.Worker,
	logger zerolog.Logger,
) ImportWorker {
	return &worker{
		importQueue: queue,
		reporter:    reporter,
		filePattern: conf.ImportCtx.FilePattern,
		worker:      importWorker,
		logger:      logger,
	}
}

func (w *worker) Run(ctx context.Context) {
	ticker := time.NewTicker(queueCheckTickInterval)
	defer ticker.Stop()
	cleanTicker := time.NewTicker(reportCleanTickInterval)
	defer cleanTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-cleanTicker.C:
			err := w.reporter.Clean(ctx, reportCleanInterval)
			if err != nil {
				w.logger.Err(err).Msg("Import-ctx: Failed to clean import reports")
			}
		case <-ticker.C:
			job := w.importQueue.Pop()
			if job.ID == "" {
				continue
			}

			err := w.reporter.ReportOngoing(ctx, job)
			if err != nil {
				w.logger.Err(err).Str("job_id", job.ID).Msg("Import-ctx: Failed to update import info")

				continue
			}

			startTime := time.Now()
			stats, err := w.doJob(ctx, job)
			stats.ExecTime = time.Since(startTime)

			if err != nil {
				w.logger.Err(err).Str("job_id", job.ID).Msg("Import-ctx: Error during the import.")

				err = w.reporter.ReportError(ctx, job, stats.ExecTime, err)
				if err != nil {
					w.logger.Err(err).Str("job_id", job.ID).Msg("Import-ctx: Failed to update import info")
				}
			} else {
				w.logger.Info().Str("job_id", job.ID).Msg("Import-ctx: import done")

				err = w.reporter.ReportDone(ctx, job, stats)
				if err != nil {
					w.logger.Err(err).Str("job_id", job.ID).Msg("Import-ctx: Failed to update import info")
				}
			}
		}
	}
}

func (w *worker) doJob(ctx context.Context, job ImportJob) (importcontextgraph.Stats, error) {
	w.logger.Info().Str("job_id", job.ID).Msg("Import-ctx: Processing import")
	filename := fmt.Sprintf(w.filePattern, job.ID)
	if job.IsPartial {
		return w.worker.WorkPartial(ctx, filename, job.Source)
	}
	return w.worker.Work(ctx, filename, job.Source)
}
