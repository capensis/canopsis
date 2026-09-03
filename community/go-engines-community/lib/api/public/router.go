package api

import (
	"embed"
	"fmt"
	"net/http"
	"path/filepath"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/docs"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/patternfields"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/alarmexport"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/asyncjob"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/bulk"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/crud"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/dataexport"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/dataimport"
	publicdocs "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/docs"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/entityexport"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/entityimport"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/middleware"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/pbehaviorreason"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/pbehaviortype"
	apisecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"github.com/go-chi/chi/v5"
)

const (
	Version    = "v5"
	Prefix     = "/api/" + Version
	bulkPrefix = "/bulk"
)

//go:embed docs/*.yaml
var docsFile embed.FS

func (s *Services) authorize(obj, act string) func(next http.Handler) http.Handler {
	return middleware.Authorize(obj, act, s.Enforcer, s.Responder)
}

// RegisterRoutes registers the API routes. The extra functions are called with the
// authenticated router of the API prefix, they allow to register additional routes.
func RegisterRoutes(s Services, opts Options, extra ...func(r chi.Router)) (*chi.Mux, error) {
	r := chi.NewRouter()

	if opts.LogOptions.Debug {
		defer func() {
			_ = chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
				s.Logger.Debug().Str("method", method).Str("route", route).Send()

				return nil
			})
		}()
	}

	r.Use(middleware.Logger(s.Responder, s.Logger, opts.LogBody, opts.LogBodyOnError))
	r.Use(middleware.Recovery(s.Logger))
	r.Use(middleware.CacheControl())
	r.Use(middleware.Maintenance(config.NewMaintenanceAdapter(s.PrimaryDbClient), s.Enforcer, s.Responder))
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		s.Responder.RespondError(w, r, httperror.ErrNotFound)
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		s.Responder.RespondError(w, r, httperror.ErrMethodNotAllowed)
	})
	r.Route(Prefix, func(r chi.Router) {
		r.Get("/", infoHandler(s.Responder))
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(s.AuthProviders, s.Responder))
			authorize := s.authorize
			{
				eventAPI := event.NewAPI(s.AMQPPubChPool, s.Responder, s.Logger)
				r.With(authorize(apisecurity.PermEvent, model.PermissionCan)).
					Post("/event", eventAPI.Send)
			}
			{
				pbhTypeAPI := pbehaviortype.NewAPI(
					pbehaviortype.NewStore(s.PrimaryDbClient),
					s.Binder,
					s.Responder,
				)
				RegisterCRUDRoutes(r, "/pbehavior-types", apisecurity.ObjPbehaviorType, pbhTypeAPI, s)
			}
			{
				pbhReasonAPI := pbehaviorreason.NewAPI(
					pbehaviorreason.NewStore(s.PrimaryDbClient),
					s.Binder,
					s.Responder,
				)
				RegisterCRUDRoutes(r, "/pbehavior-reasons", apisecurity.ObjPbehaviorReason, pbhReasonAPI, s)
			}
			{
				pbhAPI := pbehavior.NewAPI(
					pbehavior.NewStore(s.PrimaryDbClient, s.TimezoneConfigProvider, patternfields.NewTransformer(s.PrimaryDbClient)),
					s.Binder,
					s.Responder,
					amqp.NewPooledPublisher(s.AMQPPubChPool),
					json.NewEncoder(),
					json.NewDecoder(),
					s.Logger,
				)
				permObj := apisecurity.ObjPbehavior

				RegisterBulkRoutes(r, "/pbehaviors", permObj, pbhAPI, s)

				r.With(authorize(permObj, model.PermissionUpdate)).
					With(middleware.SetAuthor(s.Responder)).
					Patch("/pbehaviors/{id}", pbhAPI.Patch)

				r.Route(bulkPrefix+"/connector-pbehaviors", func(r chi.Router) {
					r.With(authorize(permObj, model.PermissionCreate)).
						With(middleware.PreProcessBulk(s.APIConfigProvider, s.Responder, true)).
						Post("/", pbhAPI.BulkConnectorCreate)

					r.With(authorize(permObj, model.PermissionDelete)).
						With(middleware.PreProcessBulk(s.APIConfigProvider, s.Responder, true)).
						Delete("/", pbhAPI.BulkConnectorDelete)

					r.With(authorize(permObj, model.PermissionCreate)).
						With(authorize(permObj, model.PermissionDelete)).
						With(middleware.PreProcessBulk(s.APIConfigProvider, s.Responder, true)).
						Put("/", pbhAPI.BulkConnectorEdit)
				})
			}
			{
				entityAPI := entity.NewAPI(entity.NewStore(s.PrimaryDbClient, patternfields.NewTransformer(s.PrimaryDbClient)), s.Binder, s.Responder)
				r.Route("/entities", func(r chi.Router) {
					r.With(authorize(apisecurity.ObjEntity, model.PermissionRead)).
						Get("/", entityAPI.List)

					r.With(authorize(apisecurity.ObjEntity, model.PermissionRead)).
						Post("/search", entityAPI.Search)
				})
			}
			{
				entityImportAPI := entityimport.NewAPI(
					asyncjob.NewStore[entityimport.Job](s.PrimaryDbClient, mongo.EntityImportPublicCollection),
					s.JobPublisher,
					filepath.Join(s.Config.File.Dir, canopsis.SubDirEntityImportPublic),
					s.Config.File.ImportMaxSize,
					s.Binder,
					s.Responder,
					s.Logger,
				)
				RegisterDataImportRoutes(r, "/entity", apisecurity.ObjEntity, model.PermissionCreate, entityImportAPI, s)
			}
			{
				entityExportAPI := entityexport.NewAPI(
					asyncjob.NewStore[entityexport.Job](s.PrimaryDbClient, mongo.EntityExportPublicCollection),
					s.JobPublisher,
					filepath.Join(s.Config.File.Dir, canopsis.SubDirEntityExportPublic),
					patternfields.NewTransformer(s.PrimaryDbClient),
					s.Binder,
					s.Responder,
					s.Logger,
				)
				RegisterDataExportRoutes(r, "/entity", apisecurity.ObjEntity, model.PermissionRead, entityExportAPI, s)
			}
			{
				alarmAPI := alarm.NewAPI(alarm.NewStore(s.PrimaryDbClient, patternfields.NewTransformer(s.PrimaryDbClient)), s.Binder, s.Responder)
				r.Route("/alarms", func(r chi.Router) {
					r.With(authorize(apisecurity.PermAlarmRead, model.PermissionCan)).
						Get("/", alarmAPI.List)

					r.With(authorize(apisecurity.PermAlarmRead, model.PermissionCan)).
						Post("/search", alarmAPI.Search)

					r.With(authorize(apisecurity.PermAlarmRead, model.PermissionCan)).
						Get("/{id}/steps", alarmAPI.Steps)
				})
			}
			{
				alarmExportAPI := alarmexport.NewAPI(
					asyncjob.NewStore[alarmexport.Job](s.PrimaryDbClient, mongo.AlarmExportPublicCollection),
					s.JobPublisher,
					filepath.Join(s.Config.File.Dir, canopsis.SubDirAlarmExportPublic),
					patternfields.NewTransformer(s.PrimaryDbClient),
					s.Binder,
					s.Responder,
					s.Logger,
				)
				RegisterDataExportRoutes(r, "/alarm", apisecurity.PermAlarmRead, model.PermissionCan, alarmExportAPI, s)
			}

			for _, f := range extra {
				f(r)
			}
		})
	})

	if opts.EnableDocs {
		err := RegisterDocs(r, s)
		if err != nil {
			return nil, err
		}
	}

	return r, nil
}

// RegisterCRUDRoutes registers CRUD routes for the pattern.
func RegisterCRUDRoutes(
	r chi.Router,
	pattern, permObj string,
	a crud.API,
	s Services,
) {
	authorize := s.authorize
	r.Route(pattern, func(r chi.Router) {
		r.With(authorize(permObj, model.PermissionRead)).
			Get("/", a.List)

		r.With(authorize(permObj, model.PermissionRead)).
			Get("/{id}", a.Get)

		r.With(authorize(permObj, model.PermissionCreate)).
			With(middleware.SetAuthor(s.Responder)).
			Post("/", a.Create)

		r.With(authorize(permObj, model.PermissionUpdate)).
			With(middleware.SetAuthor(s.Responder)).
			Put("/{id}", a.Update)

		r.With(authorize(permObj, model.PermissionDelete)).
			Delete("/{id}", a.Delete)
	})
}

// RegisterBulkRoutes registers bulk routes for the pattern.
func RegisterBulkRoutes(
	r chi.Router,
	pattern, permObj string,
	a bulk.API,
	s Services,
) {
	RegisterCRUDRoutes(r, pattern, permObj, a, s)

	authorize := s.authorize
	pattern = bulkPrefix + pattern
	r.Route(pattern, func(r chi.Router) {
		r.With(authorize(permObj, model.PermissionCreate)).
			With(middleware.PreProcessBulk(s.APIConfigProvider, s.Responder, true)).
			Post("/", a.BulkCreate)

		r.With(authorize(permObj, model.PermissionUpdate)).
			With(middleware.PreProcessBulk(s.APIConfigProvider, s.Responder, true)).
			Put("/", a.BulkUpdate)

		r.With(authorize(permObj, model.PermissionDelete)).
			With(middleware.PreProcessBulk(s.APIConfigProvider, s.Responder, false)).
			Delete("/", a.BulkDelete)
	})
}

func RegisterDataExportRoutes(
	r chi.Router,
	pattern, permObj, permAct string,
	a dataexport.API,
	s Services,
) {
	authorize := s.authorize
	r.Route(pattern+"-exports", func(r chi.Router) {
		r.With(authorize(permObj, permAct)).
			Post("/", a.Start)
		r.With(authorize(permObj, permAct)).
			Get("/{id}", a.Status)
		r.With(authorize(permObj, permAct)).
			Get("/{id}/download", a.Download)
	})
}

func RegisterDataImportRoutes(
	r chi.Router,
	pattern, permObj, permAct string,
	a dataimport.API,
	s Services,
) {
	authorize := s.authorize
	r.Route(pattern+"-imports", func(r chi.Router) {
		r.With(authorize(permObj, permAct)).
			Post("/", a.Start)
		r.With(authorize(permObj, permAct)).
			Get("/{id}", a.Status)
	})
}

func RegisterDocs(r chi.Router, s Services) error {
	content, err := docsFile.ReadFile("docs/swagger.yaml")
	if err != nil {
		return fmt.Errorf("cannot read swagger: %w", err)
	}

	schemasContent, err := docsFile.ReadFile("docs/schemas_swagger.yaml")
	if err != nil {
		return fmt.Errorf("cannot read swagger: %w", err)
	}

	r.Get(Prefix+"/docs/*", func(w http.ResponseWriter, r *http.Request) {
		path := chi.URLParam(r, "*")
		if path == "swagger.yaml" {
			publicdocs.GetHandler(
				s.Responder,
				s.Logger,
				[][]byte{schemasContent},
				[][]byte{content},
			)(w, r)

			return
		}

		backup := r.URL.Path
		defer func() {
			r.URL.Path = backup
		}()

		r.URL.Path = "static/" + path

		http.FileServer(http.FS(docs.UIFiles)).ServeHTTP(w, r)
	})

	return nil
}
