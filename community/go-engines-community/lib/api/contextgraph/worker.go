package contextgraph

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/importcontextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
)

const (
	defaultThdWarnMinPerImport = 30 * time.Minute
	defaultThdCritMinPerImport = 60 * time.Minute

	queueCheckTickInterval  = time.Second
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
	workerV1            importcontextgraph.Worker
	workerV2            importcontextgraph.Worker
}

func NewImportWorker(
	conf config.CanopsisConf,
	publisher EventPublisher,
	reporter StatusReporter,
	workerV1 importcontextgraph.Worker,
	workerV2 importcontextgraph.Worker,
	logger zerolog.Logger,
) ImportWorker {
	w := &worker{
		publisher:   publisher,
		reporter:    reporter,
		filePattern: conf.ImportCtx.FilePattern,
		workerV1:    workerV1,
		workerV2:    workerV2,
		logger:      logger,
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

func (w *worker) Run(ctx context.Context) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(queueCheckTickInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				w.processFirstJob(ctx)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(reportCleanTickInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := w.reporter.Clean(ctx, reportCleanInterval)
				if err != nil {
					w.logger.Err(err).Msg("Import-ctx: Failed to clean import reports")
				}
			}
		}
	}()

	wg.Wait()
}

func (w *worker) processFirstJob(ctx context.Context) {
	job, err := w.reporter.GetFirst(ctx, abandonedInterval)
	if err != nil {
		w.logger.Err(err).Msg("Import-ctx: Failed to get import info")
		return
	}

	if job.ID == "" {
		return
	}

	done := make(chan struct{})
	go func() {
		w.processJob(ctx, job)
		close(done)
	}()

	ticket := time.NewTicker(abandonedTickInterval)
	defer ticket.Stop()
	for {
		select {
		case <-done:
			return
		case <-ticket.C:
			err := w.reporter.ReportOngoing(ctx, job)
			if err != nil {
				w.logger.Err(err).Str("job_id", job.ID).Msg("Import-ctx: Failed to update import info")
			}
		}
	}
}

func (w *worker) processJob(ctx context.Context, job ImportJob) {
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

		ok, err := w.reporter.ReportError(ctx, job, stats.ExecTime, err)
		if err != nil {
			w.logger.Err(err).Str("job_id", job.ID).Msg("Import-ctx: Failed to update import info")
		}
		if !ok {
			return
		}
	} else {
		w.logger.Info().Str("job_id", job.ID).Msg("Import-ctx: import done")

		ok, err := w.reporter.ReportDone(ctx, job, stats)
		if err != nil {
			w.logger.Err(err).Str("job_id", job.ID).Msg("Import-ctx: Failed to update import info")
		}
		if !ok {
			return
		}
	}

	perfDataState := types.AlarmStateOK
	if stats.ExecTime > w.thdCritMinPerImport {
		perfDataState = types.AlarmStateMajor
	} else if stats.ExecTime > w.thdWarnMinPerImport {
		perfDataState = types.AlarmStateMinor
	}

	if perfDataState != types.AlarmStateOK {
		err = w.publisher.SendPerfDataEvent(ctx, job.ID, stats, types.CpsNumber(perfDataState))
		if err != nil {
			w.logger.Err(err).Str("job_id", job.ID).Msg("Import-ctx: Failed to send perf data")
		}
	}

	if resultState != types.AlarmStateOK {
		err = w.publisher.SendImportResultEvent(ctx, job.ID, stats.ExecTime, types.CpsNumber(resultState))
		if err != nil {
			w.logger.Err(err).Str("job_id", job.ID).Msg("Import-ctx: Failed to send import result event")
		}
	}
}

func (w *worker) doJob(ctx context.Context, job ImportJob) (importcontextgraph.Stats, error) {
	w.logger.Info().Str("job_id", job.ID).Msg("Import-ctx: Processing import")
	filename := fmt.Sprintf(w.filePattern, job.ID)

	if job.IsOld {
		if job.IsPartial {
			return w.workerV1.WorkPartial(ctx, filename, job.Source)
		}

		return w.workerV1.Work(ctx, filename, job.Source)
	}

	if job.IsPartial {
		return w.workerV2.WorkPartial(ctx, filename, job.Source)
	}

	return w.workerV2.Work(ctx, filename, job.Source)
}
