package entityexport

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/asyncjob"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/dataexport"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
)

// dateFields are the entity fields holding a timestamp.
var dateFields = []string{
	"last_event_date",
	"created",
	"updated",
	"imported",
	"pbehavior_info.timestamp",
}

type handler struct {
	store                  *asyncjob.Store[Job]
	dbClient               mongo.DbClient
	dbCollection           mongo.DbCollection
	dbInfoPropCollection   mongo.DbCollection
	dir                    string
	timezoneConfigProvider config.TimezoneConfigProvider
	logger                 zerolog.Logger
}

// NewHandler returns the asyncjob.Handler for asyncjob.TypeEntityExport jobs.
func NewHandler(
	store *asyncjob.Store[Job],
	dbClient mongo.DbClient,
	dir string,
	timezoneConfigProvider config.TimezoneConfigProvider,
	logger zerolog.Logger,
) asyncjob.Handler {
	return &handler{
		store:                  store,
		dbClient:               dbClient,
		dbCollection:           dbClient.Collection(mongo.EntityMongoCollection),
		dbInfoPropCollection:   dbClient.Collection(mongo.EntityInfosPropertyCollection),
		dir:                    dir,
		timezoneConfigProvider: timezoneConfigProvider,
		logger:                 logger,
	}
}

func (h *handler) Handle(ctx context.Context, jobID string) error {
	job, err := h.store.Get(ctx, jobID)
	if err != nil {
		if errors.Is(err, httperror.ErrNotFound) {
			return nil
		}

		return err
	}

	if job.Status == asyncjob.StatusOngoing {
		return nil
	}

	ok, err := h.store.SetOngoing(ctx, job.ID)
	if err != nil {
		return err
	}

	if !ok {
		// The job isn't pending anymore: it's already been taken by another worker.
		h.logger.Warn().Str("job_id", job.ID).Msg("export job is not pending, skipped")

		return nil
	}

	start := time.Now()
	err = h.do(ctx, job)
	if err != nil {
		if ctx.Err() != nil {
			return err
		}

		if updErr := h.store.SetFailed(ctx, job.ID, err.Error(), time.Since(start)); updErr != nil {
			h.logger.Err(updErr).Str("job_id", job.ID).Msg("cannot mark export job as failed")
		}

		return err
	}

	err = h.store.SetDone(ctx, job.ID, time.Since(start), nil)
	if err != nil {
		return err
	}

	return nil
}

func (h *handler) do(ctx context.Context, job Job) error {
	location := h.timezoneConfigProvider.Get().Location
	timeFormat := cmp.Or(validation.GetRealFormatTime(job.Params.TimeFormat), validation.DefaultTimeFormat)

	err := asyncjob.RunWithPing(
		ctx,
		func(ctx context.Context) error {
			return h.store.Ping(ctx, job.ID)
		},
		func(ctx context.Context) error {
			beforeLimit, afterLimit, sort, err := entity.SearchPipeline(job.Params.SearchRequest)
			if err != nil {
				return err
			}

			pipeline := beforeLimit
			pipeline = append(pipeline, sort)
			pipeline = append(pipeline, afterLimit...)
			cursor, err := h.dbCollection.Aggregate(ctx, pipeline)
			if err != nil {
				return err
			}

			opts := dataexport.Options{
				Dir:       h.dir,
				Name:      job.ID,
				Fields:    job.Params.Fields,
				Transform: dataexport.WalkObjects(dataexport.NewDateTransform(h.isDateField, location, timeFormat)),
			}

			switch job.Params.Format {
			case "csv":
				err = dataexport.ToCsv(ctx, cursor, dataexport.CsvOptions{
					Options:   opts,
					Delimiter: job.Params.Delimiter,
				})
			case "json":
				err = dataexport.ToJSON(ctx, cursor, opts)
			default:
				_ = cursor.Close(ctx)

				return fmt.Errorf("unsupported format: %q", job.Params.Format)
			}

			return err
		},
	)

	return err
}

// isDateField reports whether the entity field holds a timestamp:
// the dateFields, and the info value fields whose property is of timestamp type.
func (h *handler) isDateField(ctx context.Context, field string) (bool, error) {
	if slices.Contains(dateFields, field) {
		return true, nil
	}

	name, ok := infoValueName(field)
	if !ok {
		return false, nil
	}

	err := h.dbInfoPropCollection.FindOne(ctx, bson.M{"name": name, "type": types.EntityInfoTypeTimestamp}).Err()
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// infoValueName returns the info property name of an "infos.<name>.value" field.
func infoValueName(field string) (string, bool) {
	parts := strings.Split(field, ".")
	if len(parts) != 3 || parts[2] != "value" {
		return "", false
	}

	if parts[0] != "infos" && parts[0] != "component_infos" {
		return "", false
	}

	return parts[1], true
}
