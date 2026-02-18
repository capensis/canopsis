package api

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/middleware"
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
	Logger            zerolog.Logger
	Responder         response.Responder
	PrimaryDbClient   mongo.DbClient
	AMQPPublisher     amqp.Channel
	Config            config.CanopsisConf
	SecurityConfig    security.Config
	Enforcer          security.Enforcer
	AuthProviders     []security.HttpProvider
	APIConfigProvider config.ApiConfigProvider
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
		})
	})

	return r, nil
}
