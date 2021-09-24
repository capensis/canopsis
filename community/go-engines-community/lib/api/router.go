package api

import (
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/api/account"
	"git.canopsis.net/canopsis/go-engines/lib/api/alarm"
	"git.canopsis.net/canopsis/go-engines/lib/api/auth"
	"git.canopsis.net/canopsis/go-engines/lib/api/engineinfo"
	"git.canopsis.net/canopsis/go-engines/lib/api/entity"
	"git.canopsis.net/canopsis/go-engines/lib/api/export"
	"git.canopsis.net/canopsis/go-engines/lib/api/heartbeat"
	"git.canopsis.net/canopsis/go-engines/lib/api/logger"
	"git.canopsis.net/canopsis/go-engines/lib/api/middleware"
	"git.canopsis.net/canopsis/go-engines/lib/api/pbehavior"
	"git.canopsis.net/canopsis/go-engines/lib/api/pbehaviorcomment"
	"git.canopsis.net/canopsis/go-engines/lib/api/pbehaviorexception"
	"git.canopsis.net/canopsis/go-engines/lib/api/pbehaviorics"
	"git.canopsis.net/canopsis/go-engines/lib/api/pbehaviorreason"
	"git.canopsis.net/canopsis/go-engines/lib/api/pbehaviortimespan"
	"git.canopsis.net/canopsis/go-engines/lib/api/pbehaviortype"
	"git.canopsis.net/canopsis/go-engines/lib/api/scenario"
	"git.canopsis.net/canopsis/go-engines/lib/api/sessionstats"
	"git.canopsis.net/canopsis/go-engines/lib/api/watcherweather"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/engine"
	libpbehavior "git.canopsis.net/canopsis/go-engines/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"git.canopsis.net/canopsis/go-engines/lib/redis"
	libsecurity "git.canopsis.net/canopsis/go-engines/lib/security"
	"git.canopsis.net/canopsis/go-engines/lib/security/proxy"
	"git.canopsis.net/canopsis/go-engines/lib/security/session/stats"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

const baseUrl = "/api/v4"

const (
	authObjPbh          = "api_pbehavior"
	authObjPbhType      = "api_pbehaviortype"
	authObjPbhReason    = "api_pbehaviorreason"
	authObjPbhException = "api_pbehaviorexception"
	authObjWatcher      = "api_watcher"
	authObjHeartbeat    = "api_heartbeat"
	authObjAction       = "api_action"
	authPermAlarmRead   = "api_alarm_read"
	authPermEntityRead  = "api_entity_read"
	authEngine          = "api_engine"

	permRead   = "read"
	permCreate = "create"
	permUpdate = "update"
	permDelete = "delete"
	permCan    = "can"
)

func RegisterRoutes(
	router gin.IRouter,
	security Security,
	enforcer libsecurity.Enforcer,
	dbClient mongo.DbClient,
	location *time.Location,
	pbhStore redis.Store,
	pbhService libpbehavior.Service,
	pbhComputeChan chan<- libpbehavior.ComputeTask,
	runInfoManager engine.RunInfoManager,
	exportExecutor export.TaskExecutor,
	actionLogger logger.ActionLogger,
	logger zerolog.Logger,
) {
	sessionStore := security.GetSessionStore()
	authMiddleware := security.GetAuthMiddleware()
	security.RegisterCallbackRoutes(router)
	authApi := auth.NewApi(
		sessionStore,
		security.GetAuthProviders(),
	)
	router.POST("/auth", authApi.LoginHandler())
	sessionStatsApi := sessionstats.NewApi(sessionStore, stats.NewManager(dbClient, security.GetConfig().Session.StatsFrame))
	sessionProtected := router.Group("")
	{
		sessionProtected.Use(middleware.SessionAuth(dbClient, sessionStore), middleware.OnlyAuth())
		sessionProtected.GET("/logout", authApi.LogoutHandler())

		{
			sessionProtected.GET("/api/v2/sessionstart", sessionStatsApi.StartHandler())
			sessionProtected.POST("/api/v2/keepalive", sessionStatsApi.PingHandler())
			sessionProtected.POST("/api/v2/session_tracepath", sessionStatsApi.ChangePathHandler())
		}
	}

	getStatsHandlers := append(
		authMiddleware,
		middleware.SessionAuth(dbClient, sessionStore),
		middleware.OnlyAuth(),
		sessionStatsApi.ListHandler(),
	)
	router.GET("/api/v2/sessions", getStatsHandlers...)

	protected := router.Group(baseUrl)
	{
		protected.Use(authMiddleware...)

		protected.GET("/account/me", account.NewApi(account.NewStore(dbClient)).Me())
		protected.GET("/sessions-count", authApi.GetSessionsCount())

		alarmAPI := alarm.NewApi(alarm.NewStore(dbClient, GetLegacyURL()), exportExecutor)
		alarmRouter := protected.Group("/alarms")
		{
			alarmRouter.GET(
				"",
				middleware.Authorize(authPermAlarmRead, permCan, enforcer),
				alarmAPI.List,
			)
		}
		alarmCountersRouter := protected.Group("/alarm-counters")
		{
			alarmCountersRouter.GET(
				"",
				middleware.Authorize(authPermAlarmRead, permCan, enforcer),
				alarmAPI.Count,
			)
		}
		alarmExportRouter := protected.Group("/alarm-export")
		{
			alarmExportRouter.POST(
				"",
				middleware.Authorize(authPermAlarmRead, permCan, enforcer),
				alarmAPI.StartExport,
			)
			alarmExportRouter.GET(
				"/:id/download",
				middleware.Authorize(authPermAlarmRead, permCan, enforcer),
				alarmAPI.DownloadExport,
			)
			alarmExportRouter.GET(
				"/:id",
				middleware.Authorize(authPermAlarmRead, permCan, enforcer),
				alarmAPI.GetExport,
			)
		}

		entityAPI := entity.NewApi(entity.NewStore(dbClient), exportExecutor)
		entityExportRouter := protected.Group("/entity-export")
		{
			entityExportRouter.POST(
				"",
				middleware.Authorize(authPermEntityRead, permCan, enforcer),
				entityAPI.StartExport,
			)
			entityExportRouter.GET(
				"/:id/download",
				middleware.Authorize(authPermEntityRead, permCan, enforcer),
				entityAPI.DownloadExport,
			)
			entityExportRouter.GET(
				"/:id",
				middleware.Authorize(authPermEntityRead, permCan, enforcer),
				entityAPI.GetExport,
			)
		}

		protected.POST(
			"/pbehavior-timespans",
			middleware.Authorize(authObjPbh, permCreate, enforcer),
			pbehaviortimespan.GetTimeSpans(pbehaviortimespan.NewService(dbClient, location)),
		)
		protected.GET(
			"/pbehavior-ics/:id",
			middleware.Authorize(authObjPbh, permRead, enforcer),
			pbehaviorics.GetICS(pbehaviorics.NewStore(dbClient), pbehaviorics.NewService(location)),
		)

		pbehaviorApi := pbehavior.NewApi(
			pbehavior.NewModelTransformer(
				dbClient,
				pbehaviorreason.NewModelTransformer(),
				pbehaviorexception.NewModelTransformer(dbClient),
			),
			pbehavior.NewStore(
				dbClient,
				libpbehavior.NewEntityMatcher(dbClient),
				pbhStore,
				pbhService,
				location,
			),
			pbhComputeChan,
			actionLogger,
			logger,
		)
		pbehaviorRouter := protected.Group("/pbehaviors")
		{
			pbehaviorRouter.POST(
				"",
				middleware.Authorize(authObjPbh, permCreate, enforcer),
				pbehaviorApi.Create)
			pbehaviorRouter.GET(
				"",
				middleware.Authorize(authObjPbh, permRead, enforcer),
				pbehaviorApi.List)
			pbehaviorRouter.GET(
				"/:id",
				middleware.Authorize(authObjPbh, permRead, enforcer),
				pbehaviorApi.Get)
			pbehaviorRouter.GET(
				"/:id/eids",
				middleware.Authorize(authObjPbh, permRead, enforcer),
				pbehaviorApi.GetEIDs)
			pbehaviorRouter.PUT(
				"/:id",
				middleware.Authorize(authObjPbh, permUpdate, enforcer),
				pbehaviorApi.Update)
			pbehaviorRouter.PATCH(
				"/:id",
				middleware.Authorize(authObjPbh, permUpdate, enforcer),
				pbehaviorApi.Patch)
			pbehaviorRouter.DELETE(
				"/:id",
				middleware.Authorize(authObjPbh, permDelete, enforcer),
				pbehaviorApi.Delete)
		}
		pbehaviorCommentRouter := protected.Group("/pbehavior-comments")
		{
			pbehaviorCommentAPI := pbehaviorcomment.NewApi(pbehaviorcomment.NewModelTransformer(), pbehaviorcomment.NewStore(dbClient))
			pbehaviorCommentRouter.POST(
				"",
				middleware.Authorize(authObjPbh, permUpdate, enforcer),
				pbehaviorCommentAPI.Create,
			)
			pbehaviorCommentRouter.DELETE(
				"/:id",
				middleware.Authorize(authObjPbh, permUpdate, enforcer),
				pbehaviorCommentAPI.Delete,
			)
		}
		entityRouter := protected.Group("/entities")
		{
			entityRouter.GET(
				"/pbehaviors",
				middleware.Authorize(authPermEntityRead, permCan, enforcer),
				middleware.Authorize(authObjPbh, permRead, enforcer),
				pbehaviorApi.ListByEntityID,
			)
		}

		typeRouter := protected.Group("/pbehavior-types")
		{
			pbehaviorTypeApi := pbehaviortype.NewApi(
				pbehaviortype.NewModelTransformer(),
				pbehaviortype.NewStore(dbClient),
				pbhComputeChan,
				actionLogger,
				logger,
			)
			pbhTypeAuthorizeRead := middleware.Authorize(authObjPbhType, permRead, enforcer)
			pbhTypeAuthorizeCreate := middleware.Authorize(authObjPbhType, permCreate, enforcer)

			typeRouter.GET("", pbhTypeAuthorizeRead, pbehaviorTypeApi.List)
			typeRouter.POST("", pbhTypeAuthorizeCreate, pbehaviorTypeApi.Create)

			pbhTypeIDGroup := typeRouter.Group("")
			{
				pbhTypeAuthorizeUpdate := middleware.Authorize(authObjPbhType, permUpdate, enforcer)
				pbhTypeAuthorizeDelete := middleware.Authorize(authObjPbhType, permDelete, enforcer)

				pbhTypeIDGroup.GET("/:id", pbhTypeAuthorizeRead, pbehaviorTypeApi.Get)
				pbhTypeIDGroup.PUT("/:id", pbhTypeAuthorizeUpdate, pbehaviorTypeApi.Update)
				pbhTypeIDGroup.DELETE("/:id", pbhTypeAuthorizeDelete, pbehaviorTypeApi.Delete)
			}
		}
		reasonRouter := protected.Group("/pbehavior-reasons")
		{
			reasonAPI := pbehaviorreason.NewApi(
				pbehaviorreason.NewModelTransformer(),
				pbehaviorreason.NewStore(dbClient),
				pbhComputeChan,
				actionLogger,
				logger,
			)
			reasonRouter.POST(
				"",
				middleware.Authorize(authObjPbhReason, permCreate, enforcer),
				reasonAPI.Create)
			reasonRouter.GET(
				"",
				middleware.Authorize(authObjPbhReason, permRead, enforcer),
				reasonAPI.List)
			reasonRouter.GET(
				"/:id",
				middleware.Authorize(authObjPbhReason, permRead, enforcer),
				reasonAPI.Get)
			reasonRouter.PUT(
				"/:id",
				middleware.Authorize(authObjPbhReason, permUpdate, enforcer),
				reasonAPI.Update)
			reasonRouter.DELETE(
				"/:id",
				middleware.Authorize(authObjPbhReason, permDelete, enforcer),
				reasonAPI.Delete)
		}
		exceptionRouter := protected.Group("/pbehavior-exceptions")
		{
			exceptionAPI := pbehaviorexception.NewApi(
				pbehaviorexception.NewModelTransformer(dbClient),
				pbehaviorexception.NewStore(dbClient),
				pbhComputeChan,
				actionLogger,
				logger,
			)
			exceptionRouter.POST(
				"",
				middleware.Authorize(authObjPbhException, permCreate, enforcer),
				exceptionAPI.Create)
			exceptionRouter.GET(
				"",
				middleware.Authorize(authObjPbhException, permRead, enforcer),
				exceptionAPI.List)
			exceptionRouter.GET(
				"/:id",
				middleware.Authorize(authObjPbhException, permRead, enforcer),
				exceptionAPI.Get)
			exceptionRouter.PUT(
				"/:id",
				middleware.Authorize(authObjPbhException, permUpdate, enforcer),
				exceptionAPI.Update)
			exceptionRouter.DELETE(
				"/:id",
				middleware.Authorize(authObjPbhException, permDelete, enforcer),
				exceptionAPI.Delete)
		}

		weatherRouter := protected.Group("/weather-watchers")
		{
			statsStore := watcherweather.NewStatsStore(dbClient, location)

			weartherAPI := watcherweather.NewApi(watcherweather.NewStore(
				dbClient,
				GetLegacyURL(),
				statsStore,
				location,
			))
			weatherRouter.GET(
				"",
				middleware.Authorize(authObjWatcher, permRead, enforcer),
				weartherAPI.List,
			)
			weatherRouter.GET(
				"/:id",
				middleware.Authorize(authObjWatcher, permRead, enforcer),
				weartherAPI.EntityList,
			)
		}

		heartbeatAPI := heartbeat.NewApi(
			heartbeat.NewStore(dbClient),
			heartbeat.NewModelTransformer(),
			actionLogger,
		)
		heartbeatRouter := protected.Group("/heartbeats")
		{
			heartbeatRouter.POST(
				"",
				middleware.Authorize(authObjHeartbeat, permCreate, enforcer),
				heartbeatAPI.Create,
			)
			heartbeatRouter.GET(
				"",
				middleware.Authorize(authObjHeartbeat, permRead, enforcer),
				heartbeatAPI.List,
			)
			heartbeatRouter.GET(
				"/:id",
				middleware.Authorize(authObjHeartbeat, permRead, enforcer),
				heartbeatAPI.Get,
			)
			heartbeatRouter.PUT(
				"/:id",
				middleware.Authorize(authObjHeartbeat, permUpdate, enforcer),
				heartbeatAPI.Update,
			)
			heartbeatRouter.DELETE(
				"/:id",
				middleware.Authorize(authObjHeartbeat, permDelete, enforcer),
				heartbeatAPI.Delete,
			)
		}

		bulkRouter := protected.Group("/bulk")
		{
			heartbeatRouter := bulkRouter.Group("/heartbeats")
			{
				heartbeatRouter.POST(
					"",
					middleware.Authorize(authObjHeartbeat, permCreate, enforcer),
					heartbeatAPI.BulkCreate,
				)
				heartbeatRouter.PUT(
					"",
					middleware.Authorize(authObjHeartbeat, permUpdate, enforcer),
					heartbeatAPI.BulkUpdate,
				)
				heartbeatRouter.DELETE(
					"",
					middleware.Authorize(authObjHeartbeat, permDelete, enforcer),
					heartbeatAPI.BulkDelete,
				)
			}
		}

		protected.GET(
			"/engine-runinfo",
			middleware.Authorize(authEngine, permCan, enforcer),
			engineinfo.GetRunInfo(runInfoManager),
		)

		scenarioRouter := protected.Group("/scenarios")
		{
			scenarioAPI := scenario.NewApi(scenario.NewStore(dbClient), actionLogger)
			scenarioRouter.POST(
				"",
				middleware.Authorize(authObjAction, permCreate, enforcer),
				scenarioAPI.Create,
			)
			scenarioRouter.GET(
				"",
				middleware.Authorize(authObjAction, permRead, enforcer),
				scenarioAPI.List,
			)
			scenarioRouter.GET(
				"/:id",
				middleware.Authorize(authObjAction, permRead, enforcer),
				scenarioAPI.Get,
			)
			scenarioRouter.PUT(
				"/:id",
				middleware.Authorize(authObjAction, permUpdate, enforcer),
				scenarioAPI.Update,
			)
			scenarioRouter.DELETE(
				"/:id",
				middleware.Authorize(authObjAction, permDelete, enforcer),
				scenarioAPI.Delete,
			)
		}
	}
}

func GetProxy(
	security Security,
	enforcer libsecurity.Enforcer,
	accessConfig proxy.AccessConfig,
) []gin.HandlerFunc {
	authMiddleware := security.GetAuthMiddleware()

	return append(
		authMiddleware,
		middleware.ProxyAuthorize(enforcer, accessConfig),
		ReverseProxyHandler(),
	)
}
