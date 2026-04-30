package contextgraph

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/workers"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/importcontextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

const (
	defaultThdWarnMinPerImport = 30 * time.Minute
	defaultThdCritMinPerImport = 60 * time.Minute

	reportCleanTickInterval = time.Hour
	reportCleanInterval     = 24 * time.Hour
	abandonedTickInterval   = 4 * time.Minute
	abandonedInterval       = 5 * time.Minute
)

type worker struct {
	reporter            StatusReporter
	publisher           EventPublisher
	logger              zerolog.Logger
	filePattern         string
	thdWarnMinPerImport time.Duration
	thdCritMinPerImport time.Duration
	worker              importcontextgraph.Worker
	jobPublisher        workers.JobPublisher
}

func NewImportWorker(
	conf config.CanopsisConf,
	publisher EventPublisher,
	reporter StatusReporter,
	importWorker importcontextgraph.Worker,
	jobPublisher workers.JobPublisher,
	logger zerolog.Logger,
) ImportWorker {
	w := &worker{
		publisher:    publisher,
		reporter:     reporter,
		filePattern:  filepath.Join(conf.File.Dir, canopsis.SubDirImport, filePattern),
		worker:       importWorker,
		jobPublisher: jobPublisher,
		logger:       logger,
	}

	thdWarnMinPerImport, err := time.ParseDuration(conf.ImportCtx.ThdWarnMinPerImport)
	if err != nil {
		logger.Warn().Err(err).Msg("Can't parse thdWarnMinPerImport value, use default")
		thdWarnMinPerImport = defaultThdWarnMinPerImport
	}

	thdCritMinPerImport, err := time.ParseDuration(conf.ImportCtx.ThdCritMinPerImport)
	if err != nil {
		logger.Warn().Err(err).Msg("Can't parse thdCritMinPerImport value, use default")
		thdCritMinPerImport = defaultThdCritMinPerImport
	}

	w.thdWarnMinPerImport = thdWarnMinPerImport
	w.thdCritMinPerImport = thdCritMinPerImport

	return w
}

func (w *worker) ProcessAbandonedJob(ctx context.Context) {
	ticker := time.NewTicker(abandonedTickInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ok, err := w.reporter.HasAbandoned(ctx, abandonedInterval)
			if err != nil {
				w.logger.Err(err).Msg("failed to get import job")
				continue
			}

			if !ok {
				continue
			}

			err = w.jobPublisher.Publish(ctx, "")
			if err != nil {
				w.logger.Err(err).Msg("failed to publish import job")
				continue
			}
		}
	}
}

func (w *worker) DeleteOldJobs(ctx context.Context) {
	ticker := time.NewTicker(reportCleanTickInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := w.reporter.Clean(ctx, reportCleanInterval)
			if err != nil {
				w.logger.Err(err).Msg("failed to clean import reports")
			}
		}
	}
}

func (w *worker) ProcessFirstJob(ctx context.Context) error {
	job, err := w.reporter.GetFirst(ctx, abandonedInterval)
	if err != nil {
		return err
	}

	if job.ID == "" {
		return nil
	}

	done := make(chan struct{})
	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		err := w.processJob(gctx, job)
		close(done)

		return err
	})

	g.Go(func() error {
		ticket := time.NewTicker(abandonedTickInterval)
		defer ticket.Stop()
		for {
			select {
			case <-done:
				return nil
			case <-ticket.C:
				err := w.reporter.ReportOngoing(gctx, job)
				if err != nil {
					return err
				}
			}
		}
	})

	err = g.Wait()
	if err != nil {
		return err
	}

	// to start next job
	err = w.jobPublisher.Publish(ctx, "")
	if err != nil {
		return fmt.Errorf("failed to publish import job: %w", err)
	}

	return nil
}

func (w *worker) processJob(ctx context.Context, job ImportJob) error {
	startTime := time.Now()
	stats, jobErr := w.doJob(ctx, job)
	stats.ExecTime = time.Since(startTime)

	resultState := types.AlarmStateOK
	if jobErr != nil {
		w.logger.Err(jobErr).Str("job_id", job.ID).Msg("error during the import")
		resultState = types.AlarmStateCritical
		ok, err := w.reporter.ReportError(ctx, job, stats.ExecTime, jobErr)
		if err != nil {
			return fmt.Errorf("failed to update import info: %w", err)
		}

		if !ok {
			return nil
		}
	} else {
		w.logger.Info().Str("job_id", job.ID).Msg("import done")
		ok, err := w.reporter.ReportDone(ctx, job, stats)
		if err != nil {
			return fmt.Errorf("failed to update import info: %w", err)
		}

		if !ok {
			return nil
		}
	}

	perfDataState := types.AlarmStateOK
	if stats.ExecTime > w.thdCritMinPerImport {
		perfDataState = types.AlarmStateMajor
	} else if stats.ExecTime > w.thdWarnMinPerImport {
		perfDataState = types.AlarmStateMinor
	}

	if perfDataState != types.AlarmStateOK {
		err := w.publisher.SendPerfDataEvent(ctx, job.ID, stats, types.CpsNumber(perfDataState))
		if err != nil {
			return fmt.Errorf("failed to send perf data: %w", err)
		}
	}

	if resultState != types.AlarmStateOK {
		err := w.publisher.SendImportResultEvent(ctx, job.ID, stats.ExecTime, types.CpsNumber(resultState))
		if err != nil {
			return fmt.Errorf("failed to send import result event: %w", err)
		}
	}

	return nil
}

func (w *worker) doJob(ctx context.Context, job ImportJob) (importcontextgraph.Stats, error) {
	w.logger.Info().Str("job_id", job.ID).Msg("processing import")
	filename := fmt.Sprintf(w.filePattern, job.ID)

	if job.IsPartial {
		return w.worker.WorkPartial(ctx, filename, job.Source)
	}

	return w.worker.Work(ctx, filename, job.Source)
}
