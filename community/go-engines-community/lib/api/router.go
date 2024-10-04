package api

import (
	"net/url"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/account"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/alarmtag"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/appinfo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/associativetable"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/broadcastmessage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/contextgraph"
	libcontextgraphV1 "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/contextgraph/v1"
	libcontextgraphV2 "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/contextgraph/v2"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/engineinfo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entitybasic"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entitycategory"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entityinfodictionary"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/export"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/exportconfiguration"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/file"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/flappingrule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/idlerule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/linkrule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/messageratestats"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/middleware"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/notification"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehaviorcomment"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehaviorexception"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehaviorics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehaviorreason"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehaviortimespan"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehaviortype"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/permission"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/playlist"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/resolverule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/role"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/scenario"
	apisecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/serviceweather"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/sessionauth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/sharetoken"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/statesettings"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/user"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/userpreferences"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/view"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/viewgroup"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/viewtab"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/widget"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/widgetfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/widgettemplate"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	libentityservice "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/link"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	libpbehavior "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template/validator"
	libfile "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/file"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	libsecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/proxy"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

const baseUrl = "/api/v4"

func RegisterRoutes(
	conf config.CanopsisConf,
	router gin.IRouter,
	security Security,
	enforcer libsecurity.Enforcer,
	linkGenerator link.Generator,
	dbClient mongo.DbClient,
	dbExportClient mongo.DbClient,
	pgPoolProvider postgres.PoolProvider,
	timezoneConfigProvider config.TimezoneConfigProvider,
	templateConfigProvider config.TemplateConfigProvider,
	pbhEntityTypeResolver libpbehavior.EntityTypeResolver,
	pbhComputeChan chan<- []string,
	entityPublChan chan<- libentityservice.ChangeEntityMessage,
	entityCleanerTaskChan chan<- entity.CleanTask,
	runInfoManager engine.RunInfoManager,
	exportExecutor export.TaskExecutor,
	techMetricsTaskExecutor techmetrics.TaskExecutor,
	actionLogger logger.ActionLogger,
	publisher amqp.Publisher,
	userInterfaceConfig config.UserInterfaceConfigProvider,
	filesRoot string,
	websocketHub websocket.Hub,
	websocketStore websocket.Store,
	broadcastMessageChan chan<- bool,
	metricsEntityMetaUpdater metrics.MetaUpdater,
	metricsUserMetaUpdater metrics.MetaUpdater,
	logger zerolog.Logger,
) {
	sessionStore := security.GetSessionStore()
	authMiddleware := security.GetAuthMiddleware()
	security.RegisterCallbackRoutes(router, dbClient)
	authApi := auth.NewApi(
		security.GetTokenService(),
		security.GetTokenProviders(),
		security.GetAuthProviders(),
		websocketStore,
		security.GetCookieOptions().FileAccessName,
		security.GetCookieOptions().MaxAge,
		logger,
	)
	sessionauthApi := sessionauth.NewApi(
		sessionStore,
		security.GetAuthProviders(),
		logger,
	)
	router.POST("/auth", sessionauthApi.LoginHandler())

	sessionProtected := router.Group("")
	{
		sessionProtected.Use(middleware.SessionAuth(dbClient, sessionStore), middleware.OnlyAuth())
		sessionProtected.GET("/logout", sessionauthApi.LogoutHandler())
	}

	unprotected := router.Group(baseUrl)
	{
		unprotected.POST("/login", authApi.Login)
		unprotected.POST("/logout", authApi.Logout)
	}

	protected := router.Group(baseUrl)
	{
		protected.Use(authMiddleware...)

		protected.Group("/ws").GET("", websocket.NewApi(websocketHub).Handler)

		accountRouter := protected.Group("/account/me")
		{
			accountRouter.Use(middleware.OnlyAuth())
			accountAPI := account.NewApi(account.NewStore(dbClient, security.GetPasswordEncoder()), actionLogger)
			accountRouter.GET("", accountAPI.Me)
			accountRouter.PUT("", accountAPI.Update)
		}
		protected.GET("/logged-user-count", authApi.GetLoggedUserCount)
		protected.GET("/file-access", authApi.GetFileAccess)

		userPreferencesRouter := protected.Group("/user-preferences")
		{
			userPreferencesRouter.Use(middleware.OnlyAuth())
			userPreferencesApi := userpreferences.NewApi(userpreferences.NewStore(dbClient), widget.NewStore(dbClient), enforcer, actionLogger)
			userPreferencesRouter.GET("/:id", userPreferencesApi.Get)
			userPreferencesRouter.PUT("", userPreferencesApi.Update)
		}

		userApi := user.NewApi(user.NewStore(dbClient, security.GetPasswordEncoder(), websocketStore), actionLogger, logger,
			metricsUserMetaUpdater)
		userRouter := protected.Group("/users")
		{
			userRouter.POST("",
				middleware.Authorize(apisecurity.PermAcl, model.PermissionCreate, enforcer),
				userApi.Create,
				middleware.ReloadEnforcerPolicyOnChange(enforcer),
			)
			userRouter.GET("",
				middleware.Authorize(apisecurity.PermAcl, model.PermissionRead, enforcer),
				userApi.List,
			)
			userRouter.GET("/:id",
				middleware.Authorize(apisecurity.PermAcl, model.PermissionRead, enforcer),
				userApi.Get,
			)
			userRouter.PUT("/:id",
				middleware.Authorize(apisecurity.PermAcl, model.PermissionUpdate, enforcer),
				userApi.Update,
				middleware.ReloadEnforcerPolicyOnChange(enforcer),
			)
			userRouter.DELETE("/:id",
				middleware.Authorize(apisecurity.PermAcl, model.PermissionDelete, enforcer),
				userApi.Delete,
			)
		}
		roleRouter := protected.Group("/roles")
		{
			roleApi := role.NewApi(role.NewStore(dbClient), actionLogger)
			roleRouter.POST("",
				middleware.Authorize(apisecurity.PermAcl, model.PermissionCreate, enforcer),
				roleApi.Create,
			)
			roleRouter.GET("",
				middleware.Authorize(apisecurity.PermAcl, model.PermissionRead, enforcer),
				roleApi.List,
			)
			roleRouter.GET("/:id",
				middleware.Authorize(apisecurity.PermAcl, model.PermissionRead, enforcer),
				roleApi.Get,
			)
			roleRouter.PUT("/:id",
				middleware.Authorize(apisecurity.PermAcl, model.PermissionUpdate, enforcer),
				roleApi.Update,
				middleware.ReloadEnforcerPolicyOnChange(enforcer),
			)
			roleRouter.DELETE("/:id",
				middleware.Authorize(apisecurity.PermAcl, model.PermissionDelete, enforcer),
				roleApi.Delete,
			)
		}
		permissionRouter := protected.Group("/permissions")
		{
			permissionApi := permission.NewApi(permission.NewStore(dbClient))
			permissionRouter.GET("",
				middleware.Authorize(apisecurity.PermAcl, model.PermissionRead, enforcer),
				permissionApi.List,
			)
		}

		sharetokenApi := sharetoken.NewApi(sharetoken.NewStore(dbClient, security.GetTokenGenerator()), actionLogger)
		sharetokenRouter := protected.Group("/share-tokens")
		{
			sharetokenRouter.POST("",
				middleware.Authorize(apisecurity.PermShareToken, model.PermissionCreate, enforcer),
				sharetokenApi.Create,
			)
			sharetokenRouter.GET("",
				middleware.Authorize(apisecurity.PermShareToken, model.PermissionRead, enforcer),
				sharetokenApi.List,
			)
			sharetokenRouter.DELETE("/:id",
				middleware.Authorize(apisecurity.PermShareToken, model.PermissionDelete, enforcer),
				sharetokenApi.Delete,
			)
		}

		alarmStore := alarm.NewStore(dbClient, dbExportClient, linkGenerator, timezoneConfigProvider, logger)
		alarmAPI := alarm.NewApi(alarmStore, exportExecutor, logger)
		alarmRouter := protected.Group("/alarms")
		{
			alarmRouter.GET(
				"",
				middleware.Authorize(apisecurity.PermAlarmRead, model.PermissionCan, enforcer),
				alarmAPI.List,
			)
			alarmRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.PermAlarmRead, model.PermissionCan, enforcer),
				alarmAPI.Get,
			)
		}
		protected.POST(
			"/alarm-details",
			middleware.Authorize(apisecurity.PermAlarmRead, model.PermissionCan, enforcer),
			alarmAPI.GetDetails,
		)
		protected.GET(
			"/alarm-links/:id",
			middleware.Authorize(apisecurity.PermAlarmRead, model.PermissionCan, enforcer),
			alarmAPI.GetLinks,
		)
		protected.GET(
			"/entityservice-alarms/:id",
			middleware.Authorize(apisecurity.PermAlarmRead, model.PermissionCan, enforcer),
			alarmAPI.ListByService,
		)
		protected.GET(
			"/component-alarms",
			middleware.Authorize(apisecurity.PermAlarmRead, model.PermissionCan, enforcer),
			alarmAPI.ListByComponent,
		)
		protected.GET(
			"/resolved-alarms",
			middleware.Authorize(apisecurity.PermAlarmRead, model.PermissionCan, enforcer),
			alarmAPI.ResolvedList,
		)
		protected.GET(
			"/open-alarms",
			middleware.Authorize(apisecurity.PermAlarmRead, model.PermissionCan, enforcer),
			alarmAPI.GetOpen,
		)
		protected.GET(
			"/alarm-counters",
			middleware.Authorize(apisecurity.PermAlarmRead, model.PermissionCan, enforcer),
			alarmAPI.Count,
		)
		exportExecutor.RegisterType("alarm", alarmStore.Export)
		alarmExportRouter := protected.Group("/alarm-export")
		{
			alarmExportRouter.POST(
				"",
				middleware.Authorize(apisecurity.PermAlarmRead, model.PermissionCan, enforcer),
				alarmAPI.StartExport,
			)
			alarmExportRouter.GET(
				"/:id/download",
				security.GetFileAuthMiddleware(),
				middleware.Authorize(apisecurity.PermAlarmRead, model.PermissionCan, enforcer),
				alarmAPI.DownloadExport,
			)
			alarmExportRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.PermAlarmRead, model.PermissionCan, enforcer),
				alarmAPI.GetExport,
			)
		}

		exportConfigurationAPI := exportconfiguration.NewApi(dbClient, logger)
		protected.POST(
			"/export-configuration",
			middleware.Authorize(apisecurity.PermExportConfigurations, model.PermissionCan, enforcer),
			exportConfigurationAPI.Export,
		)

		entityStore := entity.NewStore(dbClient, dbExportClient, timezoneConfigProvider)
		entityAPI := entity.NewApi(
			entityStore,
			exportExecutor,
			entityCleanerTaskChan,
			entityPublChan,
			metricsEntityMetaUpdater,
			actionLogger,
			logger,
		)

		exportExecutor.RegisterType("entity", entityStore.Export)
		entityExportRouter := protected.Group("/entity-export")
		{
			entityExportRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjEntity, model.PermissionRead, enforcer),
				entityAPI.StartExport,
			)
			entityExportRouter.GET(
				"/:id/download",
				security.GetFileAuthMiddleware(),
				middleware.Authorize(apisecurity.ObjEntity, model.PermissionRead, enforcer),
				entityAPI.DownloadExport,
			)
			entityExportRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjEntity, model.PermissionRead, enforcer),
				entityAPI.GetExport,
			)
		}

		protected.POST(
			"/pbehavior-timespans",
			middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionRead, enforcer),
			pbehaviortimespan.GetTimeSpans(pbehaviortimespan.NewService(dbClient, timezoneConfigProvider)),
		)
		protected.GET(
			"/pbehavior-ics/:id",
			middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionRead, enforcer),
			pbehaviorics.GetICS(pbehaviorics.NewStore(dbClient), pbehaviorics.NewService(timezoneConfigProvider)),
		)

		// event-filter API
		eventFilterApi := eventfilter.NewApi(
			eventfilter.NewStore(dbClient),
			actionLogger,
			logger,
			common.NewPatternFieldsTransformer(dbClient),
		)
		eventFilterRouter := protected.Group("/eventfilter/rules")
		{
			eventFilterRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjEventFilter, model.PermissionCreate, enforcer),
				middleware.SetAuthor(),
				eventFilterApi.Create)
			eventFilterRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjEventFilter, model.PermissionRead, enforcer),
				eventFilterApi.Get)
			eventFilterRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjEventFilter, model.PermissionDelete, enforcer),
				eventFilterApi.Delete)
			eventFilterRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjEventFilter, model.PermissionRead, enforcer),
				eventFilterApi.List)
			eventFilterRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjEventFilter, model.PermissionUpdate, enforcer),
				middleware.SetAuthor(),
				eventFilterApi.Update)
		}

		pbehaviorApi := pbehavior.NewApi(
			pbehavior.NewStore(
				dbClient,
				libpbehavior.NewEntityMatcher(dbClient),
				pbhEntityTypeResolver,
				libpbehavior.NewTypeComputer(libpbehavior.NewModelProvider(dbClient), json.NewDecoder()),
				timezoneConfigProvider,
			),
			pbhComputeChan,
			common.NewPatternFieldsTransformer(dbClient),
			actionLogger,
			logger,
		)
		pbehaviorRouter := protected.Group("/pbehaviors")
		{
			pbehaviorRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionCreate, enforcer),
				middleware.SetAuthor(),
				pbehaviorApi.Create)
			pbehaviorRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionRead, enforcer),
				pbehaviorApi.List)
			pbehaviorRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionRead, enforcer),
				pbehaviorApi.Get)
			pbehaviorRouter.GET(
				"/:id/entities",
				middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionRead, enforcer),
				pbehaviorApi.ListEntities)
			pbehaviorRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionUpdate, enforcer),
				middleware.SetAuthor(),
				pbehaviorApi.Update)
			pbehaviorRouter.PATCH(
				"/:id",
				middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionUpdate, enforcer),
				middleware.SetAuthor(),
				pbehaviorApi.Patch)
			pbehaviorRouter.DELETE(
				"",
				middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionDelete, enforcer),
				pbehaviorApi.DeleteByName)
			pbehaviorRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionDelete, enforcer),
				pbehaviorApi.Delete)
		}
		pbehaviorCommentRouter := protected.Group("/pbehavior-comments")
		{
			pbehaviorCommentAPI := pbehaviorcomment.NewApi(pbehaviorcomment.NewStore(dbClient))
			pbehaviorCommentRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionUpdate, enforcer),
				middleware.SetAuthor(),
				pbehaviorCommentAPI.Create,
			)
			pbehaviorCommentRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionUpdate, enforcer),
				pbehaviorCommentAPI.Delete,
			)
		}
		entityRouter := protected.Group("/entities")
		{
			entityRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjEntity, model.PermissionRead, enforcer),
				entityAPI.List,
			)

			entityRouter.POST(
				"/clean",
				middleware.Authorize(apisecurity.ObjEntity, model.PermissionDelete, enforcer),
				entityAPI.Clean,
			)

			entityRouter.GET(
				"/context-graph",
				middleware.Authorize(apisecurity.ObjEntity, model.PermissionRead, enforcer),
				entityAPI.GetContextGraph,
			)

			entityRouter.GET(
				"/pbehaviors",
				middleware.Authorize(apisecurity.ObjEntity, model.PermissionRead, enforcer),
				middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionRead, enforcer),
				pbehaviorApi.ListByEntityID,
			)

			entityRouter.GET(
				"/pbehavior-calendar",
				middleware.Authorize(apisecurity.ObjEntity, model.PermissionRead, enforcer),
				middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionRead, enforcer),
				pbehaviorApi.CalendarByEntityID,
			)
		}

		entitybasicsAPI := entitybasic.NewApi(entitybasic.NewStore(dbClient), entityPublChan, metricsEntityMetaUpdater,
			actionLogger, logger)
		entitybasicsRouter := protected.Group("/entitybasics")
		{
			entitybasicsRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjEntity, model.PermissionRead, enforcer),
				entitybasicsAPI.Get,
			)
			entitybasicsRouter.PUT(
				"",
				middleware.Authorize(apisecurity.ObjEntity, model.PermissionUpdate, enforcer),
				entitybasicsAPI.Update,
			)
			entitybasicsRouter.DELETE(
				"",
				middleware.Authorize(apisecurity.ObjEntity, model.PermissionDelete, enforcer),
				entitybasicsAPI.Delete,
			)
		}

		entityserviceAPI := entityservice.NewApi(entityservice.NewStore(dbClient, linkGenerator, logger), entityPublChan,
			metricsEntityMetaUpdater, common.NewPatternFieldsTransformer(dbClient), actionLogger, logger)
		entityserviceRouter := protected.Group("/entityservices")
		{
			entityserviceRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjEntityService, model.PermissionCreate, enforcer),
				entityserviceAPI.Create,
			)
			entityserviceRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjEntityService, model.PermissionRead, enforcer),
				entityserviceAPI.Get,
			)
			entityserviceRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjEntityService, model.PermissionUpdate, enforcer),
				entityserviceAPI.Update,
			)
			entityserviceRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjEntityService, model.PermissionDelete, enforcer),
				entityserviceAPI.Delete,
			)
			protected.GET(
				"/entityservice-dependencies",
				middleware.Authorize(apisecurity.ObjEntityService, model.PermissionRead, enforcer),
				entityserviceAPI.GetDependencies,
			)
			protected.GET(
				"/entityservice-impacts",
				middleware.Authorize(apisecurity.ObjEntityService, model.PermissionRead, enforcer),
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
			pbhTypeAuthorizeRead := middleware.Authorize(apisecurity.ObjPbehaviorType, model.PermissionRead, enforcer)
			pbhTypeAuthorizeCreate := middleware.Authorize(apisecurity.ObjPbehaviorType, model.PermissionCreate, enforcer)
			pbhTypeAuthorizeUpdate := middleware.Authorize(apisecurity.ObjPbehaviorType, model.PermissionUpdate, enforcer)
			pbhTypeAuthorizeDelete := middleware.Authorize(apisecurity.ObjPbehaviorType, model.PermissionDelete, enforcer)

			typeRouter.GET("", pbhTypeAuthorizeRead, pbehaviorTypeApi.List)
			typeRouter.POST("", pbhTypeAuthorizeCreate, pbehaviorTypeApi.Create)
			typeRouter.GET("/next-priority", pbhTypeAuthorizeRead, pbehaviorTypeApi.GetNextPriority)
			typeRouter.GET("/:id", pbhTypeAuthorizeRead, pbehaviorTypeApi.Get)
			typeRouter.PUT("/:id", pbhTypeAuthorizeUpdate, pbehaviorTypeApi.Update)
			typeRouter.DELETE("/:id", pbhTypeAuthorizeDelete, pbehaviorTypeApi.Delete)
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
				middleware.Authorize(apisecurity.ObjPbehaviorReason, model.PermissionCreate, enforcer),
				reasonAPI.Create)
			reasonRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjPbehaviorReason, model.PermissionRead, enforcer),
				reasonAPI.List)
			reasonRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjPbehaviorReason, model.PermissionRead, enforcer),
				reasonAPI.Get)
			reasonRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjPbehaviorReason, model.PermissionUpdate, enforcer),
				reasonAPI.Update)
			reasonRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjPbehaviorReason, model.PermissionDelete, enforcer),
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
				middleware.Authorize(apisecurity.ObjPbehaviorException, model.PermissionCreate, enforcer),
				exceptionAPI.Create)
			exceptionRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjPbehaviorException, model.PermissionRead, enforcer),
				exceptionAPI.List)
			exceptionRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjPbehaviorException, model.PermissionRead, enforcer),
				exceptionAPI.Get)
			exceptionRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjPbehaviorException, model.PermissionUpdate, enforcer),
				exceptionAPI.Update)
			exceptionRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjPbehaviorException, model.PermissionDelete, enforcer),
				exceptionAPI.Delete)
		}

		weatherRouter := protected.Group("/weather-services")
		{
			weatherAPI := serviceweather.NewApi(serviceweather.NewStore(
				dbClient,
				linkGenerator,
				alarmStore,
				timezoneConfigProvider,
				logger,
			))
			weatherRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjEntityService, model.PermissionRead, enforcer),
				weatherAPI.List,
			)
			weatherRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjEntityService, model.PermissionRead, enforcer),
				weatherAPI.EntityList,
			)
		}

		entityCategoryRouter := protected.Group("/entity-categories")
		{
			entityCategoryAPI := entitycategory.NewApi(entitycategory.NewStore(dbClient), actionLogger)
			entityCategoryRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjEntityCategory, model.PermissionCreate, enforcer),
				middleware.SetAuthor(),
				entityCategoryAPI.Create,
			)
			entityCategoryRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjEntityCategory, model.PermissionRead, enforcer),
				entityCategoryAPI.List,
			)
			entityCategoryRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjEntityCategory, model.PermissionRead, enforcer),
				entityCategoryAPI.Get,
			)
			entityCategoryRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjEntityCategory, model.PermissionUpdate, enforcer),
				middleware.SetAuthor(),
				entityCategoryAPI.Update,
			)
			entityCategoryRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjEntityCategory, model.PermissionDelete, enforcer),
				entityCategoryAPI.Delete,
			)
		}

		eventApi := event.NewApi(publisher, dbClient, userInterfaceConfig.Get().IsAllowChangeSeverityToInfo, logger)
		eventRouter := protected.Group("/event")
		{
			eventRouter.POST(
				"",
				middleware.Authorize(apisecurity.PermEvent, model.PermissionCan, enforcer),
				eventApi.Send)
		}

		securityConfig := security.GetConfig().Security
		appInfoApi := appinfo.NewApi(appinfo.NewStore(dbClient, pgPoolProvider, securityConfig.AuthProviders,
			securityConfig.Cas.Title, securityConfig.Saml.Title))
		protected.GET("app-info", appInfoApi.GetAppInfo)
		appInfoRouter := protected.Group("/internal")
		{
			appInfoRouter.PUT(
				"user_interface",
				middleware.Authorize(apisecurity.PermUserInterfaceUpdate, model.PermissionCan, enforcer),
				appInfoApi.UpdateUserInterface,
			)
			appInfoRouter.POST(
				"user_interface",
				middleware.Authorize(apisecurity.PermUserInterfaceUpdate, model.PermissionCan, enforcer),
				appInfoApi.UpdateUserInterface,
			)
			appInfoRouter.DELETE(
				"user_interface",
				middleware.Authorize(apisecurity.PermUserInterfaceDelete, model.PermissionCan, enforcer),
				appInfoApi.DeleteUserInterface,
			)
		}
		protected.GET(
			"/engine-runinfo",
			middleware.Authorize(apisecurity.PermHealthcheck, model.PermissionCan, enforcer),
			engineinfo.GetRunInfo(runInfoManager),
		)

		viewAPI := view.NewApi(view.NewStore(dbClient, viewtab.NewStore(dbClient, widget.NewStore(dbClient))), enforcer, actionLogger)
		viewRouter := protected.Group("/views")
		{
			viewRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjView, model.PermissionCreate, enforcer),
				middleware.SetAuthor(),
				viewAPI.Create,
				middleware.ReloadEnforcerPolicyOnChange(enforcer),
			)
			viewRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjView, model.PermissionRead, enforcer),
				middleware.ProvideAuthorizedIds(model.PermissionRead, enforcer),
				viewAPI.List,
			)
			viewRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjView, model.PermissionRead, enforcer),
				middleware.AuthorizeByID(model.PermissionRead, enforcer),
				viewAPI.Get,
			)
			viewRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjView, model.PermissionUpdate, enforcer),
				middleware.AuthorizeByID(model.PermissionUpdate, enforcer),
				middleware.SetAuthor(),
				viewAPI.Update,
			)
			viewRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjView, model.PermissionDelete, enforcer),
				middleware.AuthorizeByID(model.PermissionDelete, enforcer),
				viewAPI.Delete,
				middleware.ReloadEnforcerPolicyOnChange(enforcer),
			)
		}

		viewTabAPI := viewtab.NewApi(viewtab.NewStore(dbClient, widget.NewStore(dbClient)), enforcer, actionLogger)
		viewTabRouter := protected.Group("/view-tabs")
		{
			viewTabRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjView, model.PermissionUpdate, enforcer),
				middleware.SetAuthor(),
				viewTabAPI.Create,
			)
			viewTabRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjView, model.PermissionRead, enforcer),
				viewTabAPI.Get,
			)
			viewTabRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjView, model.PermissionUpdate, enforcer),
				middleware.SetAuthor(),
				viewTabAPI.Update,
			)
			viewTabRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjView, model.PermissionUpdate, enforcer),
				viewTabAPI.Delete,
			)
		}

		widgetAPI := widget.NewApi(
			widget.NewStore(dbClient),
			enforcer,
			widget.NewRequestTransformer(common.NewPatternFieldsTransformer(dbClient), dbClient),
			actionLogger,
		)
		widgetRouter := protected.Group("/widgets")
		{
			widgetRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjView, model.PermissionUpdate, enforcer),
				middleware.SetAuthor(),
				widgetAPI.Create,
			)
			widgetRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjView, model.PermissionRead, enforcer),
				widgetAPI.Get,
			)
			widgetRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjView, model.PermissionUpdate, enforcer),
				middleware.SetAuthor(),
				widgetAPI.Update,
			)
			widgetRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjView, model.PermissionUpdate, enforcer),
				widgetAPI.Delete,
			)
		}

		widgetFilterAPI := widgetfilter.NewApi(widgetfilter.NewStore(dbClient), enforcer, widgetfilter.NewPatternFieldsTransformer(dbClient), actionLogger)
		widgetFilterRouter := protected.Group("/widget-filters")
		{
			widgetFilterRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjView, model.PermissionRead, enforcer),
				widgetFilterAPI.List,
			)
			widgetFilterRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjView, model.PermissionRead, enforcer), // keep PermissionRead for private filters
				middleware.SetAuthor(),
				widgetFilterAPI.Create,
			)
			widgetFilterRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjView, model.PermissionRead, enforcer),
				widgetFilterAPI.Get,
			)
			widgetFilterRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjView, model.PermissionRead, enforcer), // keep PermissionRead for private filters
				middleware.SetAuthor(),
				widgetFilterAPI.Update,
			)
			widgetFilterRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjView, model.PermissionRead, enforcer),
				widgetFilterAPI.Delete,
			)
		}

		protected.PUT(
			"/widget-filter-positions",
			middleware.Authorize(apisecurity.ObjView, model.PermissionRead, enforcer),
			widgetFilterAPI.UpdatePositions,
		)

		widgetTemplateAPI := widgettemplate.NewApi(widgettemplate.NewStore(dbClient), actionLogger)
		widgetTemplateRouter := protected.Group("/widget-templates")
		{
			widgetTemplateRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjWidgetTemplate, model.PermissionRead, enforcer),
				widgetTemplateAPI.List,
			)
			widgetTemplateRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjWidgetTemplate, model.PermissionCreate, enforcer),
				middleware.SetAuthor(),
				widgetTemplateAPI.Create,
			)
			widgetTemplateRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjWidgetTemplate, model.PermissionRead, enforcer),
				widgetTemplateAPI.Get,
			)
			widgetTemplateRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjWidgetTemplate, model.PermissionUpdate, enforcer),
				middleware.SetAuthor(),
				widgetTemplateAPI.Update,
			)
			widgetTemplateRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjWidgetTemplate, model.PermissionDelete, enforcer),
				widgetTemplateAPI.Delete,
			)
		}

		viewGroupAPI := viewgroup.NewApi(viewgroup.NewStore(dbClient), actionLogger)
		viewGroupRouter := protected.Group("/view-groups")
		{
			viewGroupRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjViewGroup, model.PermissionCreate, enforcer),
				middleware.SetAuthor(),
				viewGroupAPI.Create,
			)
			viewGroupRouter.GET(
				"",
				middleware.ProvideAuthorizedIds(model.PermissionRead, enforcer),
				middleware.Authorize(apisecurity.ObjViewGroup, model.PermissionRead, enforcer),
				viewGroupAPI.List,
			)
			viewGroupRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjViewGroup, model.PermissionRead, enforcer),
				viewGroupAPI.Get,
			)
			viewGroupRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjViewGroup, model.PermissionUpdate, enforcer),
				middleware.SetAuthor(),
				viewGroupAPI.Update,
			)
			viewGroupRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjViewGroup, model.PermissionDelete, enforcer),
				viewGroupAPI.Delete,
			)
		}

		protected.POST(
			"/view-copy/:id",
			middleware.Authorize(apisecurity.ObjView, model.PermissionCreate, enforcer),
			middleware.Authorize(apisecurity.ObjView, model.PermissionRead, enforcer),
			middleware.AuthorizeByID(model.PermissionRead, enforcer),
			middleware.SetAuthor(),
			viewAPI.Copy,
			middleware.ReloadEnforcerPolicyOnChange(enforcer),
		)

		protected.PUT(
			"/view-positions",
			middleware.Authorize(apisecurity.ObjView, model.PermissionUpdate, enforcer),
			middleware.Authorize(apisecurity.ObjViewGroup, model.PermissionUpdate, enforcer),
			viewAPI.UpdatePositions,
		)

		protected.POST(
			"/view-export",
			middleware.Authorize(apisecurity.ObjView, model.PermissionRead, enforcer),
			middleware.Authorize(apisecurity.ObjViewGroup, model.PermissionRead, enforcer),
			viewAPI.Export,
		)

		protected.POST(
			"/view-import",
			middleware.Authorize(apisecurity.ObjView, model.PermissionUpdate, enforcer),
			middleware.Authorize(apisecurity.ObjViewGroup, model.PermissionUpdate, enforcer),
			viewAPI.Import,
			middleware.ReloadEnforcerPolicyOnChange(enforcer),
		)

		protected.POST(
			"/view-tab-copy/:id",
			middleware.Authorize(apisecurity.ObjView, model.PermissionUpdate, enforcer),
			middleware.Authorize(apisecurity.ObjView, model.PermissionRead, enforcer),
			middleware.SetAuthor(),
			viewTabAPI.Copy,
		)

		protected.PUT(
			"/view-tab-positions",
			middleware.Authorize(apisecurity.ObjView, model.PermissionUpdate, enforcer),
			viewTabAPI.UpdatePositions,
		)

		protected.POST(
			"/widget-copy/:id",
			middleware.Authorize(apisecurity.ObjView, model.PermissionUpdate, enforcer),
			middleware.Authorize(apisecurity.ObjView, model.PermissionRead, enforcer),
			middleware.SetAuthor(),
			widgetAPI.Copy,
		)

		protected.PUT(
			"/widget-grid-positions",
			middleware.Authorize(apisecurity.ObjView, model.PermissionUpdate, enforcer),
			widgetAPI.UpdateGridPositions,
		)

		// broadcast message API
		broadcastMessageApi := broadcastmessage.NewApi(
			broadcastmessage.NewStore(dbClient),
			broadcastMessageChan,
			actionLogger,
		)
		broadcastMessageRouter := protected.Group("/broadcast-message")
		{

			broadcastMessageRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjBroadcastMessage, model.PermissionCreate, enforcer),
				broadcastMessageApi.Create)
			broadcastMessageRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjBroadcastMessage, model.PermissionRead, enforcer),
				broadcastMessageApi.Get)
			broadcastMessageRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjBroadcastMessage, model.PermissionDelete, enforcer),
				broadcastMessageApi.Delete)
			broadcastMessageRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjBroadcastMessage, model.PermissionRead, enforcer),
				broadcastMessageApi.List)
			broadcastMessageRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjBroadcastMessage, model.PermissionUpdate, enforcer),
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
				middleware.Authorize(apisecurity.ObjAssociativeTable, model.PermissionUpdate, enforcer),
				associativeTableApi.Update,
			)
			associativeRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjAssociativeTable, model.PermissionRead, enforcer),
				associativeTableApi.Get,
			)
			associativeRouter.DELETE(
				"",
				middleware.Authorize(apisecurity.ObjAssociativeTable, model.PermissionDelete, enforcer),
				associativeTableApi.Delete,
			)
		}

		scenarioAPI := scenario.NewApi(scenario.NewStore(dbClient), actionLogger, common.NewPatternFieldsTransformer(dbClient), logger)
		scenarioRouter := protected.Group("/scenarios")
		{
			scenarioRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjAction, model.PermissionCreate, enforcer),
				middleware.SetAuthor(),
				scenarioAPI.Create,
			)
			scenarioRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjAction, model.PermissionRead, enforcer),
				scenarioAPI.List,
			)
			scenarioRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjAction, model.PermissionRead, enforcer),
				scenarioAPI.Get,
			)
			scenarioRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjAction, model.PermissionUpdate, enforcer),
				middleware.SetAuthor(),
				scenarioAPI.Update,
			)
			scenarioRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjAction, model.PermissionDelete, enforcer),
				scenarioAPI.Delete,
			)
		}

		contextGraphAPIV1 := libcontextgraphV1.NewApi(conf, contextgraph.NewMongoStatusReporter(dbClient), logger)
		contextGraphRouter := protected.Group("/contextgraph")
		{
			contextGraphRouter.PUT(
				"import",
				middleware.Authorize(apisecurity.ObjContextGraph, model.PermissionCreate, enforcer),
				contextGraphAPIV1.ImportAll,
			)
			contextGraphRouter.PUT(
				"import-partial",
				middleware.Authorize(apisecurity.ObjContextGraph, model.PermissionCreate, enforcer),
				contextGraphAPIV1.ImportPartial,
			)
			contextGraphRouter.GET(
				"import/status/:id",
				middleware.Authorize(apisecurity.ObjContextGraph, model.PermissionRead, enforcer),
				contextGraphAPIV1.Status,
			)
		}

		contextGraphAPIV2 := libcontextgraphV2.NewApi(conf, contextgraph.NewMongoStatusReporter(dbClient), logger)
		protected.PUT(
			"contextgraph-import",
			middleware.Authorize(apisecurity.ObjContextGraph, model.PermissionCreate, enforcer),
			contextGraphAPIV2.ImportAll,
		)
		protected.PUT(
			"contextgraph-import-partial",
			middleware.Authorize(apisecurity.ObjContextGraph, model.PermissionCreate, enforcer),
			contextGraphAPIV2.ImportPartial,
		)

		stateSettingsRouter := protected.Group("/state-settings")
		{
			stateSettingsApi := statesettings.NewApi(statesettings.NewStore(dbClient), actionLogger)
			stateSettingsRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.PermStateSettings, model.PermissionCan, enforcer),
				stateSettingsApi.Update,
			)
			stateSettingsRouter.GET(
				"",
				middleware.Authorize(apisecurity.PermStateSettings, model.PermissionCan, enforcer),
				stateSettingsApi.List,
			)
		}

		notificationRouter := protected.Group("/notification")
		{
			notificationApi := notification.NewApi(notification.NewStore(dbClient), actionLogger)
			notificationRouter.PUT(
				"",
				middleware.Authorize(apisecurity.PermNotification, model.PermissionCan, enforcer),
				notificationApi.Update,
			)
			notificationRouter.GET(
				"",
				middleware.Authorize(apisecurity.PermNotification, model.PermissionCan, enforcer),
				notificationApi.Get,
			)
		}

		playlistRouter := protected.Group("/playlists")
		{
			playlistApi := playlist.NewApi(playlist.NewStore(dbClient), viewtab.NewStore(dbClient, widget.NewStore(dbClient)), enforcer, actionLogger)
			playlistRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjPlaylist, model.PermissionCreate, enforcer),
				middleware.SetAuthor(),
				playlistApi.Create,
				middleware.ReloadEnforcerPolicyOnChange(enforcer),
			)
			playlistRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjPlaylist, model.PermissionRead, enforcer),
				middleware.ProvideAuthorizedIds(model.PermissionRead, enforcer),
				playlistApi.List,
			)
			playlistRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjPlaylist, model.PermissionRead, enforcer),
				middleware.AuthorizeByID(model.PermissionRead, enforcer),
				playlistApi.Get,
			)
			playlistRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjPlaylist, model.PermissionUpdate, enforcer),
				middleware.AuthorizeByID(model.PermissionUpdate, enforcer),
				middleware.SetAuthor(),
				playlistApi.Update,
			)
			playlistRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjPlaylist, model.PermissionDelete, enforcer),
				middleware.AuthorizeByID(model.PermissionDelete, enforcer),
				playlistApi.Delete,
				middleware.ReloadEnforcerPolicyOnChange(enforcer),
			)
		}

		idleRuleAPI := idlerule.NewApi(idlerule.NewStore(dbClient), common.NewPatternFieldsTransformer(dbClient), actionLogger, logger)
		idleRuleRouter := protected.Group("/idle-rules")
		{
			idleRuleRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjIdleRule, model.PermissionCreate, enforcer),
				middleware.SetAuthor(),
				idleRuleAPI.Create,
			)
			idleRuleRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjIdleRule, model.PermissionRead, enforcer),
				idleRuleAPI.List,
			)
			idleRuleRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjIdleRule, model.PermissionRead, enforcer),
				idleRuleAPI.Get,
			)
			idleRuleRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjIdleRule, model.PermissionUpdate, enforcer),
				middleware.SetAuthor(),
				idleRuleAPI.Update,
			)
			idleRuleRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjIdleRule, model.PermissionDelete, enforcer),
				idleRuleAPI.Delete,
			)
		}

		patternAPI := pattern.NewApi(pattern.NewStore(dbClient, pbhComputeChan, entityPublChan, logger), userInterfaceConfig,
			enforcer, actionLogger, logger)
		patternRouter := protected.Group("/patterns")
		{
			patternRouter.Use(middleware.OnlyAuth())
			patternRouter.POST(
				"",
				middleware.SetAuthor(),
				patternAPI.Create,
			)
			patternRouter.GET(
				"",
				patternAPI.List,
			)
			patternRouter.GET(
				"/:id",
				patternAPI.Get,
			)
			patternRouter.PUT(
				"/:id",
				middleware.SetAuthor(),
				patternAPI.Update,
			)
			patternRouter.DELETE(
				"/:id",
				patternAPI.Delete,
			)
		}
		protected.POST(
			"/patterns-count",
			middleware.OnlyAuth(),
			patternAPI.Count,
		)
		protected.POST(
			"/patterns-alarms",
			middleware.OnlyAuth(),
			patternAPI.GetAlarms,
		)

		linkRuleAPI := linkrule.NewApi(
			linkrule.NewStore(dbClient),
			common.NewPatternFieldsTransformer(dbClient),
			actionLogger,
			logger,
		)
		linkRuleRouter := protected.Group("/link-rules")
		{
			linkRuleRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjLinkRule, model.PermissionCreate, enforcer),
				middleware.SetAuthor(),
				linkRuleAPI.Create,
			)
			linkRuleRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjLinkRule, model.PermissionRead, enforcer),
				linkRuleAPI.List,
			)
			linkRuleRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjLinkRule, model.PermissionRead, enforcer),
				linkRuleAPI.Get,
			)
			linkRuleRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjLinkRule, model.PermissionUpdate, enforcer),
				middleware.SetAuthor(),
				linkRuleAPI.Update,
			)
			linkRuleRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjLinkRule, model.PermissionDelete, enforcer),
				linkRuleAPI.Delete,
			)
		}
		linkCategoryRouter := protected.Group("/link-categories")
		{
			linkCategoryRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjLinkRule, model.PermissionRead, enforcer),
				linkRuleAPI.GetCategories,
			)
		}

		bulkRouter := protected.Group("/bulk")
		{
			patternRouter := bulkRouter.Group("/patterns")
			{
				patternRouter.DELETE(
					"",
					middleware.PreProcessBulk(conf, false),
					patternAPI.BulkDelete,
				)
			}

			scenarioRouter := bulkRouter.Group("/scenarios")
			{
				scenarioRouter.POST(
					"",
					middleware.Authorize(apisecurity.ObjAction, model.PermissionCreate, enforcer),
					middleware.PreProcessBulk(conf, true),
					scenarioAPI.BulkCreate,
				)
				scenarioRouter.PUT(
					"",
					middleware.Authorize(apisecurity.ObjAction, model.PermissionUpdate, enforcer),
					middleware.PreProcessBulk(conf, true),
					scenarioAPI.BulkUpdate,
				)
				scenarioRouter.DELETE(
					"",
					middleware.Authorize(apisecurity.ObjAction, model.PermissionDelete, enforcer),
					middleware.PreProcessBulk(conf, false),
					scenarioAPI.BulkDelete,
				)
			}

			idleruleRouter := bulkRouter.Group("/idle-rules")
			{
				idleruleRouter.POST(
					"",
					middleware.Authorize(apisecurity.ObjIdleRule, model.PermissionCreate, enforcer),
					middleware.PreProcessBulk(conf, true),
					idleRuleAPI.BulkCreate,
				)
				idleruleRouter.PUT(
					"",
					middleware.Authorize(apisecurity.ObjIdleRule, model.PermissionUpdate, enforcer),
					middleware.PreProcessBulk(conf, true),
					idleRuleAPI.BulkUpdate,
				)
				idleruleRouter.DELETE(
					"",
					middleware.Authorize(apisecurity.ObjIdleRule, model.PermissionDelete, enforcer),
					middleware.PreProcessBulk(conf, false),
					idleRuleAPI.BulkDelete,
				)
			}

			eventFilterRouter := bulkRouter.Group("/eventfilters")
			{
				eventFilterRouter.POST(
					"",
					middleware.Authorize(apisecurity.ObjEventFilter, model.PermissionCreate, enforcer),
					middleware.PreProcessBulk(conf, true),
					eventFilterApi.BulkCreate,
				)
				eventFilterRouter.PUT(
					"",
					middleware.Authorize(apisecurity.ObjEventFilter, model.PermissionUpdate, enforcer),
					middleware.PreProcessBulk(conf, true),
					eventFilterApi.BulkUpdate,
				)
				eventFilterRouter.DELETE(
					"",
					middleware.Authorize(apisecurity.ObjEventFilter, model.PermissionDelete, enforcer),
					middleware.PreProcessBulk(conf, false),
					eventFilterApi.BulkDelete,
				)
			}

			entityserviceRouter := bulkRouter.Group("/entityservices")
			{
				entityserviceRouter.POST(
					"",
					middleware.Authorize(apisecurity.ObjEntityService, model.PermissionCreate, enforcer),
					middleware.PreProcessBulk(conf, false),
					entityserviceAPI.BulkCreate,
				)
				entityserviceRouter.PUT(
					"",
					middleware.Authorize(apisecurity.ObjEntityService, model.PermissionUpdate, enforcer),
					middleware.PreProcessBulk(conf, false),
					entityserviceAPI.BulkUpdate,
				)
				entityserviceRouter.DELETE(
					"",
					middleware.Authorize(apisecurity.ObjEntityService, model.PermissionDelete, enforcer),
					middleware.PreProcessBulk(conf, false),
					entityserviceAPI.BulkDelete,
				)
			}

			userRouter := bulkRouter.Group("/users")
			{
				userRouter.POST(
					"",
					middleware.Authorize(apisecurity.PermAcl, model.PermissionCreate, enforcer),
					middleware.PreProcessBulk(conf, false),
					userApi.BulkCreate,
				)
				userRouter.PUT(
					"",
					middleware.Authorize(apisecurity.PermAcl, model.PermissionUpdate, enforcer),
					middleware.PreProcessBulk(conf, false),
					userApi.BulkUpdate,
				)
				userRouter.DELETE(
					"",
					middleware.Authorize(apisecurity.PermAcl, model.PermissionDelete, enforcer),
					middleware.PreProcessBulk(conf, false),
					userApi.BulkDelete,
				)
			}

			pbehaviorRouter := bulkRouter.Group("/pbehaviors")
			{
				pbehaviorRouter.POST(
					"",
					middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionCreate, enforcer),
					middleware.PreProcessBulk(conf, true),
					pbehaviorApi.BulkCreate,
				)
				pbehaviorRouter.PUT(
					"",
					middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionUpdate, enforcer),
					middleware.PreProcessBulk(conf, true),
					pbehaviorApi.BulkUpdate,
				)
				pbehaviorRouter.DELETE(
					"",
					middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionDelete, enforcer),
					middleware.PreProcessBulk(conf, false),
					pbehaviorApi.BulkDelete,
				)
			}

			entityPbehaviorRouter := bulkRouter.Group("/entity-pbehaviors")
			{
				entityPbehaviorRouter.POST(
					"",
					middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionCreate, enforcer),
					middleware.PreProcessBulk(conf, true),
					pbehaviorApi.BulkEntityCreate,
				)
				entityPbehaviorRouter.DELETE(
					"",
					middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionDelete, enforcer),
					middleware.PreProcessBulk(conf, false),
					pbehaviorApi.BulkEntityDelete,
				)
			}

			entityRouter := bulkRouter.Group("/entities")
			{
				entityRouter.PUT(
					"/enable",
					middleware.Authorize(apisecurity.ObjEntity, model.PermissionUpdate, enforcer),
					middleware.PreProcessBulk(conf, false),
					entityAPI.BulkEnable,
				)
				entityRouter.PUT(
					"/disable",
					middleware.Authorize(apisecurity.ObjEntity, model.PermissionUpdate, enforcer),
					middleware.PreProcessBulk(conf, false),
					entityAPI.BulkDisable,
				)
			}

			linkRuleRouter := bulkRouter.Group("/link-rules")
			{
				linkRuleRouter.DELETE(
					"",
					middleware.Authorize(apisecurity.ObjLinkRule, model.PermissionDelete, enforcer),
					middleware.PreProcessBulk(conf, false),
					linkRuleAPI.BulkDelete,
				)
			}
		}

		dateStorageRouter := protected.Group("data-storage")
		{
			dateStorageAPI := datastorage.NewApi(datastorage.NewStore(dbClient, pgPoolProvider, logger))
			dateStorageRouter.GET(
				"",
				middleware.Authorize(apisecurity.PermDataStorageRead, model.PermissionCan, enforcer),
				dateStorageAPI.Get,
			)
			dateStorageRouter.PUT(
				"",
				middleware.Authorize(apisecurity.PermDataStorageUpdate, model.PermissionCan, enforcer),
				dateStorageAPI.Update,
			)
		}

		messageRateStatsRouter := protected.Group("/message-rate-stats")
		{
			messageRateStatsAPI := messageratestats.NewApi(messageratestats.NewStore(dbClient))
			messageRateStatsRouter.GET(
				"",
				middleware.Authorize(apisecurity.PermMessageRateStatsRead, model.PermissionCan, enforcer),
				messageRateStatsAPI.List,
			)
		}

		fileRouter := protected.Group("/file")
		{
			fileAPI := file.NewApi(enforcer, file.NewStore(dbClient, libfile.NewStorage(
				filesRoot,
				libfile.NewEtagEncoder(),
			), conf.File.UploadMaxSize))
			fileRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjFile, model.PermissionCreate, enforcer),
				fileAPI.Create,
			)
			fileRouter.GET(
				"",
				fileAPI.List,
			)
			fileRouter.GET(
				"/:id",
				security.GetFileAuthMiddleware(),
				fileAPI.Get,
			)
			fileRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjFile, model.PermissionDelete, enforcer),
				fileAPI.Delete,
			)
		}

		resolveRuleRouter := protected.Group("/resolve-rules")
		{
			resolveRuleAPI := resolverule.NewApi(
				resolverule.NewStore(dbClient),
				common.NewPatternFieldsTransformer(dbClient),
				actionLogger,
			)
			resolveRuleRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjResolveRule, model.PermissionCreate, enforcer),
				middleware.SetAuthor(),
				resolveRuleAPI.Create,
			)
			resolveRuleRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjResolveRule, model.PermissionRead, enforcer),
				resolveRuleAPI.List,
			)
			resolveRuleRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjResolveRule, model.PermissionRead, enforcer),
				resolveRuleAPI.Get,
			)
			resolveRuleRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjResolveRule, model.PermissionUpdate, enforcer),
				middleware.SetAuthor(),
				resolveRuleAPI.Update,
			)
			resolveRuleRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjResolveRule, model.PermissionDelete, enforcer),
				resolveRuleAPI.Delete,
			)
		}

		flappingRuleRouter := protected.Group("/flapping-rules")
		{
			flappingRuleAPI := flappingrule.NewApi(flappingrule.NewStore(dbClient), common.NewPatternFieldsTransformer(dbClient), actionLogger)
			flappingRuleRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjFlappingRule, model.PermissionCreate, enforcer),
				middleware.SetAuthor(),
				flappingRuleAPI.Create,
			)
			flappingRuleRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjFlappingRule, model.PermissionRead, enforcer),
				flappingRuleAPI.List,
			)
			flappingRuleRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjFlappingRule, model.PermissionRead, enforcer),
				flappingRuleAPI.Get,
			)
			flappingRuleRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjFlappingRule, model.PermissionUpdate, enforcer),
				middleware.SetAuthor(),
				flappingRuleAPI.Update,
			)
			flappingRuleRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjFlappingRule, model.PermissionDelete, enforcer),
				flappingRuleAPI.Delete,
			)
		}

		entityInfoDictionaryApi := entityinfodictionary.NewApi(entityinfodictionary.NewStore(dbClient), logger)
		protected.GET("/entity-infos-dictionary/keys",
			middleware.Authorize(apisecurity.ObjEntity, model.PermissionRead, enforcer),
			entityInfoDictionaryApi.ListKeys,
		)
		protected.GET("/entity-infos-dictionary/values",
			middleware.Authorize(apisecurity.ObjEntity, model.PermissionRead, enforcer),
			entityInfoDictionaryApi.ListValues,
		)

		alarmTagRouter := protected.Group("/alarm-tags")
		{
			alarmTagAPI := alarmtag.NewApi(alarmtag.NewStore(dbClient))
			alarmTagRouter.GET(
				"",
				middleware.Authorize(apisecurity.PermAlarmRead, model.PermissionCan, enforcer),
				alarmTagAPI.List,
			)
		}

		techMetricsRouter := protected.Group("/tech-metrics-export")
		{
			techMetricsAPI := techmetrics.NewApi(techMetricsTaskExecutor, timezoneConfigProvider)
			techMetricsRouter.POST(
				"",
				middleware.Authorize(apisecurity.PermTechMetrics, model.PermissionCan, enforcer),
				techMetricsAPI.StartExport,
			)
			techMetricsRouter.GET(
				"",
				middleware.Authorize(apisecurity.PermTechMetrics, model.PermissionCan, enforcer),
				techMetricsAPI.GetExport,
			)
			techMetricsRouter.GET(
				"/download",
				security.GetFileAuthMiddleware(),
				middleware.Authorize(apisecurity.PermTechMetrics, model.PermissionCan, enforcer),
				techMetricsAPI.DownloadExport,
			)
		}

		templateValidatorApi := template.NewApi(validator.NewValidator(timezoneConfigProvider), templateConfigProvider)
		templateValidatorRouter := protected.Group("/template-validator")
		{
			templateValidatorRouter.POST(
				"/declare-ticket-rules",
				middleware.OnlyAuth(),
				templateValidatorApi.ValidateDeclareTicketRules,
			)
			templateValidatorRouter.POST(
				"/scenarios",
				middleware.OnlyAuth(),
				templateValidatorApi.ValidateScenarios,
			)
		}
		protected.GET(
			"/template-vars",
			templateValidatorApi.GetEnvVars,
		)
	}
}

func GetProxy(
	legacyUrl *url.URL,
	security Security,
	enforcer libsecurity.Enforcer,
	accessConfig proxy.AccessConfig,
) []gin.HandlerFunc {
	authMiddleware := security.GetAuthMiddleware()

	return append(
		authMiddleware,
		middleware.ProxyAuthorize(enforcer, accessConfig),
		ReverseProxyHandler(legacyUrl),
	)
}
