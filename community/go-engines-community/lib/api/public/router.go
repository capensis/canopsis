package api

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/binding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/crud"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/middleware"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/pbehaviorreason"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/pbehaviortype"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/response"
	apisecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

const Version = "v5"
const Prefix = "/api/" + Version

type Services struct {
	Logger                 zerolog.Logger
	Binder                 binding.Binder
	Responder              response.Responder
	PrimaryDbClient        mongo.DbClient
	AMQPPublisher          amqp.Channel
	Config                 config.CanopsisConf
	SecurityConfig         security.Config
	Enforcer               security.Enforcer
	AuthProviders          []security.HttpProvider
	APIConfigProvider      config.ApiConfigProvider
	TimezoneConfigProvider config.TimezoneConfigProvider
	AuthorProvider         author.Provider
}

func (s *Services) authorize(obj, act string) func(next http.Handler) http.Handler {
	return middleware.Authorize(obj, act, s.Enforcer, s.Responder)
}

func RegisterRoutes(s Services, opts Options) (*chi.Mux, error) {
	r := chi.NewRouter()
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
				eventAPI := event.NewAPI(s.AMQPPublisher, s.Responder, s.Logger)
				r.With(authorize(apisecurity.PermEvent, model.PermissionCan)).
					Post("/event", eventAPI.Send)
			}
			{
				pbhTypeAPI := pbehaviortype.NewAPI(
					pbehaviortype.NewStore(s.PrimaryDbClient, s.AuthorProvider),
					s.Binder,
					s.Responder,
				)
				registerCRUDRoutes(r, "/pbehavior-types", apisecurity.ObjPbehaviorType, pbhTypeAPI, s)
			}
			{
				pbhReasonAPI := pbehaviorreason.NewAPI(
					pbehaviorreason.NewStore(s.PrimaryDbClient, s.AuthorProvider),
					s.Binder,
					s.Responder,
				)
				registerCRUDRoutes(r, "/pbehavior-reasons", apisecurity.ObjPbehaviorReason, pbhReasonAPI, s)
			}
		})
	})

	return r, nil
}

func registerCRUDRoutes(r chi.Router, pattern, permObj string, a crud.API, s Services) {
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
