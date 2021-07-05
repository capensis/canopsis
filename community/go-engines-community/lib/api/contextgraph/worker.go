package contextgraph

import (
	"context"
	"errors"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/importcontextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
	"time"
)

const (
	DefaultThdWarnMinPerImport = 30 * time.Minute
	DefaultThdCritMinPerImport = 60 * time.Minute
)

type worker struct {
	importQueue         JobQueue
	reporter            StatusReporter
	publisher           EventPublisher
	logger              zerolog.Logger
	filePattern         string
	thdWarnMinPerImport time.Duration
	thdCritMinPerImport time.Duration
	worker              importcontextgraph.Worker
}

func NewImportWorker(
	conf config.CanopsisConf,
	publisher EventPublisher,
	reporter StatusReporter,
	queue JobQueue,
	importWorker importcontextgraph.Worker,
	logger zerolog.Logger,
) ImportWorker {
	w := &worker{
		importQueue: queue,
		publisher:   publisher,
		reporter:    reporter,
		filePattern: conf.ImportCtx.FilePattern,
		worker:      importWorker,
	}

	thdWarnMinPerImport, err := time.ParseDuration(conf.ImportCtx.ThdWarnMinPerImport)
	if err != nil {
		logger.Warn().Err(err).Msg("Can't parse thdWarnMinPerImport value, use default")
		thdWarnMinPerImport = DefaultThdWarnMinPerImport
	}

	thdCritMinPerImport, err := time.ParseDuration(conf.ImportCtx.ThdCritMinPerImport)
	if err != nil {
		logger.Warn().Err(err).Msg("Can't parse thdCritMinPerImport value, use default")
		thdCritMinPerImport = DefaultThdCritMinPerImport
	}

	w.thdWarnMinPerImport = thdWarnMinPerImport
	w.thdCritMinPerImport = thdCritMinPerImport

	return w
}

func (w *worker) Run(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			job := w.importQueue.Pop()
			if job.ID == "" {
				continue
			}

			err := w.reporter.ReportOngoing(ctx, job)
			if err != nil {
				w.logger.Err(err).Str("job_id", job.ID).Msg("Import-ctx: Failed to update import info")

				err = w.publisher.SendImportResultEvent(job.ID, 0, types.AlarmStateCritical)
				if err != nil {
					w.logger.Err(err).Str("job_id", job.ID).Msg("Import-ctx: Failed to send import result event")
				}

				continue
			}

			startTime := time.Now()
			stats, err := w.doJob(ctx, job)
			stats.ExecTime = time.Since(startTime)

			resultState := types.AlarmStateOK
			if err != nil {
				w.logger.Err(err).Str("job_id", job.ID).Msg("Import-ctx: Error during the import.")

				resultState = types.AlarmStateCritical
				if errors.Is(err, importcontextgraph.ErrNotImplemented) {
					resultState = types.AlarmStateMinor
				}

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

			perfDataState := types.AlarmStateOK
			if stats.ExecTime > w.thdCritMinPerImport {
				perfDataState = types.AlarmStateMajor
			} else if stats.ExecTime > w.thdWarnMinPerImport {
				perfDataState = types.AlarmStateMinor
			}

			err = w.publisher.SendPerfDataEvent(job.ID, stats, types.CpsNumber(perfDataState))
			if err != nil {
				w.logger.Err(err).Str("job_id", job.ID).Msg("Import-ctx: Failed to send perf data")
			}

			err = w.publisher.SendImportResultEvent(job.ID, stats.ExecTime, types.CpsNumber(resultState))
			if err != nil {
				w.logger.Err(err).Str("job_id", job.ID).Msg("Import-ctx: Failed to send import result event")
			}
		}
	}
}

func (w *worker) doJob(ctx context.Context, job ImportJob) (importcontextgraph.Stats, error) {
	w.logger.Info().Str("job_id", job.ID).Msg("Import-ctx: Processing import")
	filename := fmt.Sprintf(w.filePattern, job.ID)
	return w.worker.Work(ctx, filename)
}
