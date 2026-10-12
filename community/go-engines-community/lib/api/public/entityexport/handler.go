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
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/link"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	liblink "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/link"
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

// linksField is the exported field holding the links of the entity: it's generated, not fetched.
const linksField = "links"

type handler struct {
	store                  *asyncjob.Store[Job]
	dbClient               mongo.DbClient
	dbCollection           mongo.DbCollection
	dbInfoPropCollection   mongo.DbCollection
	dir                    string
	timezoneConfigProvider config.TimezoneConfigProvider
	generator              link.Generator
	userProvider           liblink.UserProvider
	logger                 zerolog.Logger
}

// NewHandler returns the asyncjob.Handler for asyncjob.TypeEntityExport jobs.
func NewHandler(
	store *asyncjob.Store[Job],
	dbClient mongo.DbClient,
	dir string,
	timezoneConfigProvider config.TimezoneConfigProvider,
	generator link.Generator,
	userProvider liblink.UserProvider,
	logger zerolog.Logger,
) asyncjob.Handler {
	return &handler{
		store:                  store,
		dbClient:               dbClient,
		dbCollection:           dbClient.Collection(mongo.EntityMongoCollection),
		dbInfoPropCollection:   dbClient.Collection(mongo.EntityInfosPropertyCollection),
		dir:                    dir,
		timezoneConfigProvider: timezoneConfigProvider,
		generator:              generator,
		userProvider:           userProvider,
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

	job, ok, err := h.store.Claim(ctx, job.ID)
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
		if errors.Is(err, asyncjob.ErrLeaseLost) {
			h.logger.Warn().Str("job_id", job.ID).Msg("export job lease lost")

			return nil
		}

		if ctx.Err() != nil {
			return err
		}

		if updErr := h.store.SetFailed(ctx, job.ID, job.LeaseGeneration, err.Error(), time.Since(start)); updErr != nil {
			if errors.Is(updErr, asyncjob.ErrLeaseLost) {
				h.logger.Warn().Str("job_id", job.ID).Msg("export job lease lost")

				return nil
			}

			h.logger.Err(updErr).Str("job_id", job.ID).Msg("cannot mark export job as failed")
		}

		return err
	}

	err = h.store.SetDone(ctx, job.ID, job.LeaseGeneration, time.Since(start), nil)
	if err != nil {
		if errors.Is(err, asyncjob.ErrLeaseLost) {
			h.logger.Warn().Str("job_id", job.ID).Msg("export job lease lost")

			return nil
		}

		return err
	}

	return nil
}

func (h *handler) do(ctx context.Context, job Job) error {
	location := h.timezoneConfigProvider.Get().Location
	timeFormat := cmp.Or(validation.GetRealFormatTime(job.Params.TimeFormat), validation.DefaultTimeFormat)

	enrich, err := h.linksEnrich(ctx, job)
	if err != nil {
		return err
	}

	err = asyncjob.RunWithPing(
		ctx,
		func(ctx context.Context) error {
			return h.store.Ping(ctx, job.ID, job.LeaseGeneration)
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
				Enrich:    enrich,
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

// linksEnrich returns the dataexport.Enrich adding to every exported document the links
// generated for its entity, as the list endpoints do.
func (h *handler) linksEnrich(ctx context.Context, job Job) (dataexport.Enrich, error) {
	exported := slices.ContainsFunc(job.Params.Fields, func(f dataexport.Field) bool {
		return f.Name == linksField
	})
	if !exported {
		return nil, nil
	}

	user, err := h.userProvider.User(ctx, job.Author)
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context, doc bson.M) error {
		var e entity.Response

		b, err := bson.Marshal(doc)
		if err != nil {
			return err
		}

		err = bson.Unmarshal(b, &e)
		if err != nil {
			return err
		}

		doc[linksField] = h.generator.GenerateForEntity(ctx, e.Raw, user)

		return nil
	}, nil
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
