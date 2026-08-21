package api

import (
	"path/filepath"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/patternfields"
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
	importStore := asyncjob.NewStore[entityimport.Job](s.PrimaryDbClient, mongo.EntityImportPublicCollection)
	importDir := filepath.Join(s.Config.File.Dir, canopsis.SubDirEntityImportPublic)
	exportStore := asyncjob.NewStore[entityexport.Job](s.PrimaryDbClient, mongo.EntityExportPublicCollection)
	exportDir := filepath.Join(s.Config.File.Dir, canopsis.SubDirEntityExportPublic)

	handlers := map[string]asyncjob.Handler{
		asyncjob.TypeEntityImport: entityimport.NewHandler(
			importStore,
			importDir,
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
			exportStore,
			s.SecondaryDbClient,
			exportDir,
			s.TimezoneConfigProvider,
			s.Logger,
		),
	}

	maintainers := []asyncjob.Maintainer{
		asyncjob.NewMaintainer(importStore, asyncjob.TypeEntityImport, s.JobPublisher, func(job entityimport.Job) string {
			return dataimport.Filepath(importDir, job.ID)
		}),
		asyncjob.NewMaintainer(exportStore, asyncjob.TypeEntityExport, s.JobPublisher, func(job entityexport.Job) string {
			return dataexport.Filepath(exportDir, job.ID, job.Params.Format)
		}),
	}

	return handlers, maintainers
}
