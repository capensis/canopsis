package api

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/account"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/appinfo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/associativetable"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/broadcastmessage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/contextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/engineinfo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entitybasic"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entitycategory"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/export"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/heartbeat"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/idlerule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/messageratestats"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/middleware"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehaviorcomment"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehaviorexception"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehaviorics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehaviorreason"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehaviortimespan"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehaviortype"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/permission"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/playlist"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/role"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/scenario"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/serviceweather"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/sessionstats"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/statesettings"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/user"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/userpreferences"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/view"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/viewgroup"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	libentityservice "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	libpbehavior "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	libsecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/proxy"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/session/stats"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

const baseUrl = "/api/v4"

const (
	authObjPbh              = "api_pbehavior"
	authObjPbhType          = "api_pbehaviortype"
	authObjPbhReason        = "api_pbehaviorreason"
	authObjPbhException     = "api_pbehaviorexception"
	authObjHeartbeat        = "api_heartbeat"
	authObjAction           = "api_action"
	authObjEntity           = "api_entity"
	authObjEntityService    = "api_entityservice"
	authObjEntityCategory   = "api_entitycategory"
	authObjView             = "api_view"
	authObjViewGroup        = "api_viewgroup"
	authObjPlaylist         = "api_playlist"
	authPermAlarmRead       = "api_alarm_read"
	authEngine              = "api_engine"
	authObjContextGraph     = "api_contextgraph"
	authAcl                 = "api_acl"
	authObjStateSettings    = "api_state_settings"
	authDataStorageRead     = "api_datastorage_read"
	authDataStorageUpdate   = "api_datastorage_update"
	authEventFilter         = "api_eventfilter"
	authBroadcastMessage    = "api_broadcast_message"
	authAssociativeTable    = "api_associative_table"
	authAppInfoRead         = "api_app_info_read"
	authUserInterfaceUpdate = "api_user_interface_update"
	authUserInterfaceDelete = "api_user_interface_delete"
	authEvent               = "api_event"
	authObjIdleRule         = "api_idlerule"

	authMessageRateStatsRead = "api_message_rate_stats_read"

	permRead   = model.PermissionRead
	permCreate = model.PermissionCreate
	permUpdate = model.PermissionUpdate
	permDelete = model.PermissionDelete
	permCan    = model.PermissionCan
)

func RegisterRoutes(
	ctx context.Context,
	conf config.CanopsisConf,
	router gin.IRouter,
	security Security,
	enforcer libsecurity.Enforcer,
	dbClient mongo.DbClient,
	timezoneConfigProvider config.TimezoneConfigProvider,
	pbhEntityTypeResolver libpbehavior.EntityTypeResolver,
	pbhComputeChan chan<- libpbehavior.ComputeTask,
	entityPublChan chan<- libentityservice.ChangeEntityMessage,
	runInfoManager engine.RunInfoManager,
	exportExecutor export.TaskExecutor,
	actionLogger logger.ActionLogger,
	publisher amqp.Publisher,
	jobQueue contextgraph.JobQueue,
	userInterfaceConfig config.UserInterfaceConfigProvider,
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

		protected.GET("/account/me", account.NewApi(account.NewStore(dbClient)).Me)
		protected.GET("/sessions-count", authApi.GetSessionsCount())

		userPreferencesRouter := protected.Group("/user-preferences")
		{
			userPreferencesRouter.Use(middleware.OnlyAuth())
			userPreferencesApi := userpreferences.NewApi(userpreferences.NewStore(dbClient), actionLogger)
			userPreferencesRouter.GET("/:id", userPreferencesApi.Get)
			userPreferencesRouter.PUT("", userPreferencesApi.Update)
		}

		userRouter := protected.Group("/users")
		{
			userRouter.Use(middleware.Authorize(authAcl, permCan, enforcer))
			userApi := user.NewApi(user.NewStore(dbClient, security.GetPasswordEncoder()), actionLogger)
			userRouter.POST("",
				userApi.Create,
				middleware.ReloadEnforcerPolicyOnChange(enforcer),
			)
			userRouter.GET("", userApi.List)
			userRouter.GET("/:id", userApi.Get)
			userRouter.PUT("/:id",
				userApi.Update,
				middleware.ReloadEnforcerPolicyOnChange(enforcer),
			)
			userRouter.DELETE("/:id", userApi.Delete)
		}
		roleRouter := protected.Group("/roles")
		{
			roleRouter.Use(middleware.Authorize(authAcl, permCan, enforcer))
			roleApi := role.NewApi(role.NewStore(dbClient), actionLogger)
			roleRouter.POST("", roleApi.Create)
			roleRouter.GET("", roleApi.List)
			roleRouter.GET("/:id", roleApi.Get)
			roleRouter.PUT("/:id",
				roleApi.Update,
				middleware.ReloadEnforcerPolicyOnChange(enforcer),
			)
			roleRouter.DELETE("/:id", roleApi.Delete)
		}
		permissionRouter := protected.Group("/permissions")
		{
			permissionRouter.Use(middleware.Authorize(authAcl, permCan, enforcer))
			permissionApi := permission.NewApi(permission.NewStore(dbClient))
			permissionRouter.GET("", permissionApi.List)
		}

		alarmAPI := alarm.NewApi(alarm.NewStore(dbClient, GetLegacyURL()), exportExecutor, timezoneConfigProvider)
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
				middleware.Authorize(authObjEntity, permRead, enforcer),
				entityAPI.StartExport,
			)
			entityExportRouter.GET(
				"/:id/download",
				middleware.Authorize(authObjEntity, permRead, enforcer),
				entityAPI.DownloadExport,
			)
			entityExportRouter.GET(
				"/:id",
				middleware.Authorize(authObjEntity, permRead, enforcer),
				entityAPI.GetExport,
			)
		}

		protected.POST(
			"/pbehavior-timespans",
			middleware.Authorize(authObjPbh, permCreate, enforcer),
			pbehaviortimespan.GetTimeSpans(pbehaviortimespan.NewService(dbClient, timezoneConfigProvider)),
		)
		protected.GET(
			"/pbehavior-ics/:id",
			middleware.Authorize(authObjPbh, permRead, enforcer),
			pbehaviorics.GetICS(pbehaviorics.NewStore(dbClient), pbehaviorics.NewService(timezoneConfigProvider)),
		)

		// event-filter API
		eventFilterApi := eventfilter.NewApi(
			eventfilter.NewStore(dbClient),
			actionLogger,
		)
		eventFilterRouter := protected.Group("/eventfilter/rules")
		{
			eventFilterRouter.POST(
				"",
				middleware.Authorize(authEventFilter, permCreate, enforcer),
				middleware.SetAuthor(),
				eventFilterApi.Create)
			eventFilterRouter.GET(
				"/:id",
				middleware.Authorize(authEventFilter, permRead, enforcer),
				eventFilterApi.Get)
			eventFilterRouter.DELETE(
				"/:id",
				middleware.Authorize(authEventFilter, permDelete, enforcer),
				eventFilterApi.Delete)
			eventFilterRouter.GET(
				"",
				middleware.Authorize(authEventFilter, permRead, enforcer),
				eventFilterApi.List)
			eventFilterRouter.PUT(
				"/:id",
				middleware.Authorize(authEventFilter, permUpdate, enforcer),
				middleware.SetAuthor(),
				eventFilterApi.Update)
		}

		pbehaviorApi := pbehavior.NewApi(
			pbehavior.NewModelTransformer(
				dbClient,
				pbehaviorreason.NewModelTransformer(),
				pbehaviorexception.NewModelTransformer(dbClient),
			),
			pbehavior.NewStore(
				dbClient,
				libpbehavior.NewEntityMatcher(dbClient),
				pbhEntityTypeResolver,
				timezoneConfigProvider,
			),
			pbhComputeChan,
			userInterfaceConfig,
			actionLogger,
			logger,
		)
		pbehaviorRouter := protected.Group("/pbehaviors")
		{
			pbehaviorRouter.POST(
				"",
				middleware.Authorize(authObjPbh, permCreate, enforcer),
				middleware.SetAuthor(),
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
				"/:id/entities",
				middleware.Authorize(authObjPbh, permRead, enforcer),
				pbehaviorApi.ListEntities)
			pbehaviorRouter.PUT(
				"/:id",
				middleware.Authorize(authObjPbh, permUpdate, enforcer),
				middleware.SetAuthor(),
				pbehaviorApi.Update)
			pbehaviorRouter.PATCH(
				"/:id",
				middleware.Authorize(authObjPbh, permUpdate, enforcer),
				middleware.SetAuthor(),
				pbehaviorApi.Patch)
			pbehaviorRouter.DELETE(
				"",
				middleware.Authorize(authObjPbh, permDelete, enforcer),
				pbehaviorApi.DeleteByName)
			pbehaviorRouter.DELETE(
				"/:id",
				middleware.Authorize(authObjPbh, permDelete, enforcer),
				pbehaviorApi.Delete)
			pbehaviorRouter.POST(
				"/count",
				middleware.Authorize(authObjPbh, permCreate, enforcer),
				pbehaviorApi.CountFilter)
		}
		pbehaviorCommentRouter := protected.Group("/pbehavior-comments")
		{
			pbehaviorCommentAPI := pbehaviorcomment.NewApi(pbehaviorcomment.NewModelTransformer(), pbehaviorcomment.NewStore(dbClient))
			pbehaviorCommentRouter.POST(
				"",
				middleware.Authorize(authObjPbh, permUpdate, enforcer),
				middleware.SetAuthor(),
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
			entityAPI := entity.NewApi(entity.NewStore(dbClient), exportExecutor)
			entityRouter.GET(
				"",
				middleware.Authorize(authObjEntity, permRead, enforcer),
				entityAPI.List,
			)

			entityRouter.GET(
				"/pbehaviors",
				middleware.Authorize(authObjEntity, permRead, enforcer),
				middleware.Authorize(authObjPbh, permRead, enforcer),
				pbehaviorApi.ListByEntityID,
			)
		}
		entitybasicsRouter := protected.Group("/entitybasics")
		{
			entitybasicsAPI := entitybasic.NewApi(entitybasic.NewStore(dbClient), entityPublChan, actionLogger, logger)
			entitybasicsRouter.GET(
				"",
				middleware.Authorize(authObjEntity, permRead, enforcer),
				entitybasicsAPI.Get,
			)
			entitybasicsRouter.PUT(
				"",
				middleware.Authorize(authObjEntity, permUpdate, enforcer),
				entitybasicsAPI.Update,
			)
			entitybasicsRouter.DELETE(
				"",
				middleware.Authorize(authObjEntity, permDelete, enforcer),
				entitybasicsAPI.Delete,
			)
		}
		entityserviceRouter := protected.Group("/entityservices")
		{
			entityserviceAPI := entityservice.NewApi(entityservice.NewStore(dbClient), entityPublChan, actionLogger, logger)
			entityserviceRouter.POST(
				"",
				middleware.Authorize(authObjEntityService, permCreate, enforcer),
				entityserviceAPI.Create,
			)
			entityserviceRouter.GET(
				"/:id",
				middleware.Authorize(authObjEntityService, permRead, enforcer),
				entityserviceAPI.Get,
			)
			entityserviceRouter.PUT(
				"/:id",
				middleware.Authorize(authObjEntityService, permUpdate, enforcer),
				entityserviceAPI.Update,
			)
			entityserviceRouter.DELETE(
				"/:id",
				middleware.Authorize(authObjEntityService, permDelete, enforcer),
				entityserviceAPI.Delete,
			)
			protected.GET(
				"/entityservice-dependencies",
				middleware.Authorize(authObjEntityService, permRead, enforcer),
				entityserviceAPI.GetDependencies,
			)
			protected.GET(
				"/entityservice-impacts",
				middleware.Authorize(authObjEntityService, permRead, enforcer),
				entityserviceAPI.GetImpacts,
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

		weatherRouter := protected.Group("/weather-services")
		{
			statsStore := serviceweather.NewStatsStore(dbClient)
			weatherAPI := serviceweather.NewApi(serviceweather.NewStore(
				dbClient,
				GetLegacyURL(),
				statsStore,
				timezoneConfigProvider,
				pbhEntityTypeResolver,
			))
			weatherRouter.GET(
				"",
				middleware.Authorize(authObjEntityService, permRead, enforcer),
				weatherAPI.List,
			)
			weatherRouter.GET(
				"/:id",
				middleware.Authorize(authObjEntityService, permRead, enforcer),
				weatherAPI.EntityList,
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
				middleware.SetAuthor(),
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
				middleware.SetAuthor(),
				heartbeatAPI.Update,
			)
			heartbeatRouter.DELETE(
				"/:id",
				middleware.Authorize(authObjHeartbeat, permDelete, enforcer),
				heartbeatAPI.Delete,
			)
		}

		entityCategoryRouter := protected.Group("/entity-categories")
		{
			entityCategoryAPI := entitycategory.NewApi(entitycategory.NewStore(dbClient), actionLogger)
			entityCategoryRouter.POST(
				"",
				middleware.Authorize(authObjEntityCategory, permCreate, enforcer),
				middleware.SetAuthor(),
				entityCategoryAPI.Create,
			)
			entityCategoryRouter.GET(
				"",
				middleware.Authorize(authObjEntityCategory, permRead, enforcer),
				entityCategoryAPI.List,
			)
			entityCategoryRouter.GET(
				"/:id",
				middleware.Authorize(authObjEntityCategory, permRead, enforcer),
				entityCategoryAPI.Get,
			)
			entityCategoryRouter.PUT(
				"/:id",
				middleware.Authorize(authObjEntityCategory, permUpdate, enforcer),
				middleware.SetAuthor(),
				entityCategoryAPI.Update,
			)
			entityCategoryRouter.DELETE(
				"/:id",
				middleware.Authorize(authObjEntityCategory, permDelete, enforcer),
				entityCategoryAPI.Delete,
			)
		}

		eventApi := event.NewApi(publisher, dbClient, userInterfaceConfig.Get().IsAllowChangeSeverityToInfo, logger)
		eventRouter := protected.Group("/event")
		{
			eventRouter.POST(
				"",
				middleware.Authorize(authEvent, permCan, enforcer),
				eventApi.Send)
		}

		appInfoApi := appinfo.NewApi(appinfo.NewStore(dbClient, security.GetConfig().Security.AuthProviders))
		appInfoRouter := protected.Group("/internal")
		{
			appInfoRouter.GET("login_info", appInfoApi.LoginInfo)
			appInfoRouter.GET(
				"app_info",
				middleware.Authorize(authAppInfoRead, permCan, enforcer),
				appInfoApi.GetAppInfo,
			)

			appInfoRouter.PUT(
				"user_interface",
				middleware.Authorize(authUserInterfaceUpdate, permCan, enforcer),
				appInfoApi.UpdateUserInterface,
			)
			appInfoRouter.POST(
				"user_interface",
				middleware.Authorize(authUserInterfaceUpdate, permCan, enforcer),
				appInfoApi.UpdateUserInterface,
			)
			appInfoRouter.DELETE(
				"user_interface",
				middleware.Authorize(authUserInterfaceDelete, permCan, enforcer),
				appInfoApi.DeleteUserInterface,
			)
		}
		protected.GET(
			"/engine-runinfo",
			middleware.Authorize(authEngine, permCan, enforcer),
			engineinfo.GetRunInfo(ctx, runInfoManager),
		)

		viewAPI := view.NewApi(view.NewStore(dbClient), actionLogger)
		viewRouter := protected.Group("/views")
		{
			viewRouter.POST(
				"",
				middleware.Authorize(authObjView, permCreate, enforcer),
				middleware.SetAuthor(),
				viewAPI.Create,
				middleware.ReloadEnforcerPolicyOnChange(enforcer),
			)
			viewRouter.GET(
				"",
				middleware.Authorize(authObjView, permRead, enforcer),
				middleware.ProvideAuthorizedIds(permRead, enforcer),
				viewAPI.List,
			)
			viewRouter.GET(
				"/:id",
				middleware.Authorize(authObjView, permRead, enforcer),
				middleware.AuthorizeByID(permRead, enforcer),
				viewAPI.Get,
			)
			viewRouter.PUT(
				"/:id",
				middleware.Authorize(authObjView, permUpdate, enforcer),
				middleware.AuthorizeByID(permUpdate, enforcer),
				middleware.SetAuthor(),
				viewAPI.Update,
			)
			viewRouter.DELETE(
				"/:id",
				middleware.Authorize(authObjView, permDelete, enforcer),
				middleware.AuthorizeByID(permDelete, enforcer),
				viewAPI.Delete,
				middleware.ReloadEnforcerPolicyOnChange(enforcer),
			)
		}

		viewGroupAPI := viewgroup.NewApi(viewgroup.NewStore(dbClient), actionLogger)
		viewGroupRouter := protected.Group("/view-groups")
		{
			viewGroupRouter.POST(
				"",
				middleware.Authorize(authObjViewGroup, permCreate, enforcer),
				middleware.SetAuthor(),
				viewGroupAPI.Create,
			)
			viewGroupRouter.GET(
				"",
				middleware.ProvideAuthorizedIds(permRead, enforcer),
				middleware.Authorize(authObjViewGroup, permRead, enforcer),
				viewGroupAPI.List,
			)
			viewGroupRouter.GET(
				"/:id",
				middleware.Authorize(authObjViewGroup, permRead, enforcer),
				viewGroupAPI.Get,
			)
			viewGroupRouter.PUT(
				"/:id",
				middleware.Authorize(authObjViewGroup, permUpdate, enforcer),
				middleware.SetAuthor(),
				viewGroupAPI.Update,
			)
			viewGroupRouter.DELETE(
				"/:id",
				middleware.Authorize(authObjViewGroup, permDelete, enforcer),
				viewGroupAPI.Delete,
			)
		}

		protected.PUT(
			"/view-positions",
			middleware.Authorize(authObjView, permUpdate, enforcer),
			middleware.Authorize(authObjViewGroup, permUpdate, enforcer),
			middleware.ProvideAuthorizedIds(permUpdate, enforcer),
			viewAPI.UpdatePositions,
		)

		// broadcast message API
		broadcastMessageApi := broadcastmessage.NewApi(
			broadcastmessage.NewStore(dbClient),
			actionLogger,
		)
		broadcastMessageRouter := protected.Group("/broadcast-message")
		{

			broadcastMessageRouter.POST(
				"",
				middleware.Authorize(authBroadcastMessage, permCreate, enforcer),
				broadcastMessageApi.Create)
			broadcastMessageRouter.GET(
				"/:id",
				middleware.Authorize(authBroadcastMessage, permRead, enforcer),
				broadcastMessageApi.Get)
			broadcastMessageRouter.DELETE(
				"/:id",
				middleware.Authorize(authBroadcastMessage, permDelete, enforcer),
				broadcastMessageApi.Delete)
			broadcastMessageRouter.GET(
				"",
				middleware.Authorize(authBroadcastMessage, permRead, enforcer),
				broadcastMessageApi.List)
			broadcastMessageRouter.PUT(
				"/:id",
				middleware.Authorize(authBroadcastMessage, permUpdate, enforcer),
				broadcastMessageApi.Update)
			// can not make typical format like /api/v4/broadcast-message/active
			// because it would be failed with conflict error apart of get /:id route
			router.GET(baseUrl+"/active-broadcast-message", broadcastMessageApi.GetActive)
		}

		associativeTableApi := associativetable.NewApi(
			associativetable.NewStore(dbClient),
			actionLogger,
		)
		associativeRouter := protected.Group("/associativetable")
		{
			associativeRouter.POST(
				"",
				middleware.Authorize(authAssociativeTable, permUpdate, enforcer),
				associativeTableApi.Update,
			)
			associativeRouter.GET(
				"",
				middleware.Authorize(authAssociativeTable, permRead, enforcer),
				associativeTableApi.Get,
			)
			associativeRouter.DELETE(
				"",
				middleware.Authorize(authAssociativeTable, permDelete, enforcer),
				associativeTableApi.Delete,
			)
		}

		scenarioRouter := protected.Group("/scenarios")
		{
			scenarioAPI := scenario.NewApi(scenario.NewStore(dbClient), actionLogger)
			scenarioRouter.POST(
				"",
				middleware.Authorize(authObjAction, permCreate, enforcer),
				middleware.SetAuthor(),
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
				middleware.SetAuthor(),
				scenarioAPI.Update,
			)
			scenarioRouter.DELETE(
				"/:id",
				middleware.Authorize(authObjAction, permDelete, enforcer),
				scenarioAPI.Delete,
			)
		}

		contextGraphRouter := protected.Group("/contextgraph")
		{
			contextGraphAPI := contextgraph.NewApi(conf, jobQueue, contextgraph.NewMongoStatusReporter(dbClient), logger)
			contextGraphRouter.PUT(
				"import",
				middleware.Authorize(authObjContextGraph, permCreate, enforcer),
				contextGraphAPI.Import,
			)
			contextGraphRouter.GET(
				"import/status/:id",
				middleware.Authorize(authObjContextGraph, permRead, enforcer),
				contextGraphAPI.Status,
			)
		}

		stateSettingsRouter := protected.Group("/state-settings")
		{
			stateSettingsApi := statesettings.NewApi(statesettings.NewStore(dbClient), actionLogger)
			stateSettingsRouter.PUT(
				"/:id",
				middleware.Authorize(authObjStateSettings, permCan, enforcer),
				stateSettingsApi.Update,
			)
			stateSettingsRouter.GET(
				"/",
				middleware.Authorize(authObjStateSettings, permCan, enforcer),
				stateSettingsApi.List,
			)
		}

		playlistRouter := protected.Group("/playlists")
		{
			playlistApi := playlist.NewApi(playlist.NewStore(dbClient), actionLogger)
			playlistRouter.POST(
				"",
				middleware.Authorize(authObjPlaylist, permCreate, enforcer),
				middleware.SetAuthor(),
				playlistApi.Create,
				middleware.ReloadEnforcerPolicyOnChange(enforcer),
			)
			playlistRouter.GET(
				"",
				middleware.Authorize(authObjPlaylist, permRead, enforcer),
				middleware.ProvideAuthorizedIds(permRead, enforcer),
				playlistApi.List,
			)
			playlistRouter.GET(
				"/:id",
				middleware.Authorize(authObjPlaylist, permRead, enforcer),
				middleware.AuthorizeByID(permRead, enforcer),
				playlistApi.Get,
			)
			playlistRouter.PUT(
				"/:id",
				middleware.Authorize(authObjPlaylist, permUpdate, enforcer),
				middleware.AuthorizeByID(permUpdate, enforcer),
				middleware.SetAuthor(),
				playlistApi.Update,
			)
			playlistRouter.DELETE(
				"/:id",
				middleware.Authorize(authObjPlaylist, permDelete, enforcer),
				middleware.AuthorizeByID(permDelete, enforcer),
				playlistApi.Delete,
				middleware.ReloadEnforcerPolicyOnChange(enforcer),
			)
		}

		bulkRouter := protected.Group("/bulk")
		{
			heartbeatRouter := bulkRouter.Group("/heartbeats")
			{
				heartbeatRouter.POST(
					"",
					middleware.Authorize(authObjHeartbeat, permCreate, enforcer),
					middleware.SetAuthorToBulk(),
					heartbeatAPI.BulkCreate,
				)
				heartbeatRouter.PUT(
					"",
					middleware.Authorize(authObjHeartbeat, permUpdate, enforcer),
					middleware.SetAuthorToBulk(),
					heartbeatAPI.BulkUpdate,
				)
				heartbeatRouter.DELETE(
					"",
					middleware.Authorize(authObjHeartbeat, permDelete, enforcer),
					heartbeatAPI.BulkDelete,
				)
			}

			viewRouter := bulkRouter.Group("/views")
			{
				viewRouter.POST(
					"",
					middleware.Authorize(authObjView, permCreate, enforcer),
					middleware.SetAuthorToBulk(),
					viewAPI.BulkCreate,
					middleware.ReloadEnforcerPolicyOnChange(enforcer),
				)
				viewRouter.PUT(
					"",
					middleware.Authorize(authObjView, permUpdate, enforcer),
					middleware.ProvideAuthorizedIds(permUpdate, enforcer),
					middleware.SetAuthorToBulk(),
					viewAPI.BulkUpdate,
				)
				viewRouter.DELETE(
					"",
					middleware.Authorize(authObjView, permDelete, enforcer),
					middleware.ProvideAuthorizedIds(permDelete, enforcer),
					viewAPI.BulkDelete,
					middleware.ReloadEnforcerPolicyOnChange(enforcer),
				)
			}

			viewGroupRouter := bulkRouter.Group("/view-groups")
			{
				viewGroupRouter.POST(
					"",
					middleware.Authorize(authObjViewGroup, permCreate, enforcer),
					middleware.SetAuthorToBulk(),
					viewGroupAPI.BulkCreate,
				)
				viewGroupRouter.PUT(
					"",
					middleware.Authorize(authObjViewGroup, permUpdate, enforcer),
					middleware.SetAuthorToBulk(),
					viewGroupAPI.BulkUpdate,
				)
				viewGroupRouter.DELETE(
					"",
					middleware.Authorize(authObjViewGroup, permDelete, enforcer),
					viewGroupAPI.BulkDelete,
				)
			}
		}

		dateStorageRouter := protected.Group("data-storage")
		{
			dateStorageAPI := datastorage.NewApi(datastorage.NewStore(dbClient))
			dateStorageRouter.GET(
				"",
				middleware.Authorize(authDataStorageRead, permCan, enforcer),
				dateStorageAPI.Get,
			)
			dateStorageRouter.PUT(
				"",
				middleware.Authorize(authDataStorageUpdate, permCan, enforcer),
				dateStorageAPI.Update,
			)
		}

		messageRateStatsRouter := protected.Group("/message-rate-stats")
		{
			messageRateStatsAPI := messageratestats.NewApi(messageratestats.NewStore(dbClient))
			messageRateStatsRouter.GET(
				"",
				middleware.Authorize(authMessageRateStatsRead, permCan, enforcer),
				messageRateStatsAPI.List,
			)
		}

		idleRuleRouter := protected.Group("/idle-rules")
		{
			idleRuleAPI := idlerule.NewApi(idlerule.NewStore(dbClient), actionLogger, userInterfaceConfig)
			idleRuleRouter.POST(
				"",
				middleware.Authorize(authObjIdleRule, permCreate, enforcer),
				idleRuleAPI.Create,
			)
			idleRuleRouter.GET(
				"",
				middleware.Authorize(authObjIdleRule, permRead, enforcer),
				idleRuleAPI.List,
			)
			idleRuleRouter.GET(
				"/:id",
				middleware.Authorize(authObjIdleRule, permRead, enforcer),
				idleRuleAPI.Get,
			)
			idleRuleRouter.PUT(
				"/:id",
				middleware.Authorize(authObjIdleRule, permUpdate, enforcer),
				idleRuleAPI.Update,
			)
			idleRuleRouter.DELETE(
				"/:id",
				middleware.Authorize(authObjIdleRule, permDelete, enforcer),
				idleRuleAPI.Delete,
			)
			idleRuleRouter.POST(
				"/count",
				middleware.Authorize(authObjPbh, permCreate, enforcer),
				idleRuleAPI.CountPatterns)
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
