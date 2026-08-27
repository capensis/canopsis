package api

import (
	"path/filepath"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/patternfields"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/alarmexport"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/asyncjob"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/dataexport"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/dataimport"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/entityexport"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/entityimport"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/importcontextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
)

// AsyncJobs returns the handler and the maintainer of every job type:
// the handler runs a job, the maintainer recovers the ones a dead worker left behind and deletes the ones which expired.
func AsyncJobs(s Services) (map[string]asyncjob.Handler, []asyncjob.Maintainer) {
	entityImportStore := asyncjob.NewStore[entityimport.Job](s.PrimaryDbClient, mongo.EntityImportPublicCollection)
	entityImportDir := filepath.Join(s.Config.File.Dir, canopsis.SubDirEntityImportPublic)
	entityExportStore := asyncjob.NewStore[entityexport.Job](s.PrimaryDbClient, mongo.EntityExportPublicCollection)
	entityExportDir := filepath.Join(s.Config.File.Dir, canopsis.SubDirEntityExportPublic)
	alarmExportStore := asyncjob.NewStore[alarmexport.Job](s.PrimaryDbClient, mongo.AlarmExportPublicCollection)
	alarmExportDir := filepath.Join(s.Config.File.Dir, canopsis.SubDirAlarmExportPublic)

	handlers := map[string]asyncjob.Handler{
		asyncjob.TypeEntityImport: entityimport.NewHandler(
			entityImportStore,
			entityImportDir,
			importcontextgraph.NewWorker(
				s.PrimaryDbClient,
				importcontextgraph.NewEventPublisher(canopsis.DefaultExchangeName, canopsis.FIFOQueueName, json.NewEncoder(), canopsis.JsonContentType, libamqp.NewPooledPublisher(s.AMQPPubChPool)),
				metrics.NewNullMetaUpdater(),
				canopsis.ApiConnector,
				patternfields.NewTransformer(s.PrimaryDbClient),
				s.Logger,
			),
			s.Logger,
		),
		asyncjob.TypeEntityExport: entityexport.NewHandler(
			entityExportStore,
			s.SecondaryDbClient,
			entityExportDir,
			s.TimezoneConfigProvider,
			s.Logger,
		),
		asyncjob.TypeAlarmExport: alarmexport.NewHandler(
			alarmExportStore,
			s.SecondaryDbClient,
			alarmExportDir,
			s.TimezoneConfigProvider,
			s.Logger,
		),
	}

	maintainers := []asyncjob.Maintainer{
		asyncjob.NewMaintainer(entityImportStore, asyncjob.TypeEntityImport, s.JobPublisher, func(job entityimport.Job) string {
			return dataimport.Filepath(entityImportDir, job.ID)
		}),
		asyncjob.NewMaintainer(entityExportStore, asyncjob.TypeEntityExport, s.JobPublisher, func(job entityexport.Job) string {
			return dataexport.Filepath(entityExportDir, job.ID, job.Params.Format)
		}),
		asyncjob.NewMaintainer(alarmExportStore, asyncjob.TypeAlarmExport, s.JobPublisher, func(job alarmexport.Job) string {
			return dataexport.Filepath(alarmExportDir, job.ID, job.Params.Format)
		}),
	}

	return handlers, maintainers
}
