package api

import (
	"context"
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
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/user"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/userpreferences"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/view"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/viewgroup"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/viewtab"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/widget"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/widgetfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/action"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	libentityservice "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	libpbehavior "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
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

const (
	authObjAction = apisecurity.ObjAction

	authObjEntity         = apisecurity.ObjEntity
	authObjEntityService  = apisecurity.ObjEntityService
	authObjEntityCategory = apisecurity.ObjEntityCategory
	authObjContextGraph   = apisecurity.ObjContextGraph

	authObjPlaylist = apisecurity.ObjPlaylist

	authPermAlarmRead = apisecurity.PermAlarmRead

	authObjStateSettings = apisecurity.PermStateSettings

	authDataStorageRead   = apisecurity.PermDataStorageRead
	authDataStorageUpdate = apisecurity.PermDataStorageUpdate

	authEventFilter = apisecurity.ObjEventFilter

	authBroadcastMessage = apisecurity.ObjBroadcastMessage

	authAssociativeTable = apisecurity.ObjAssociativeTable

	authUserInterfaceUpdate = apisecurity.PermUserInterfaceUpdate
	authUserInterfaceDelete = apisecurity.PermUserInterfaceDelete

	authEvent = apisecurity.PermEvent

	authObjIdleRule = apisecurity.ObjIdleRule

	authObjNotification = apisecurity.PermNotification

	authMessageRateStatsRead = apisecurity.PermMessageRateStatsRead

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
	legacyUrl string,
	dbClient mongo.DbClient,
	pgPoolProvider postgres.PoolProvider,
	timezoneConfigProvider config.TimezoneConfigProvider,
	pbhEntityTypeResolver libpbehavior.EntityTypeResolver,
	pbhComputeChan chan<- libpbehavior.ComputeTask,
	entityPublChan chan<- libentityservice.ChangeEntityMessage,
	entityCleanerTaskChan chan<- entity.CleanTask,
	runInfoManager engine.RunInfoManager,
	exportExecutor export.TaskExecutor,
	techMetricsTaskExecutor techmetrics.TaskExecutor,
	actionLogger logger.ActionLogger,
	publisher amqp.Publisher,
	jobQueue contextgraph.JobQueue,
	userInterfaceConfig config.UserInterfaceConfigProvider,
	scenarioPriorityIntervals action.PriorityIntervals,
	filesRoot string,
	websocketHub websocket.Hub,
	broadcastMessageChan chan<- bool,
	metricsEntityMetaUpdater metrics.MetaUpdater,
	metricsUserMetaUpdater metrics.MetaUpdater,
	logger zerolog.Logger,
) {
	linksFetcher := common.NewLinksFetcher(legacyUrl)
	sessionStore := security.GetSessionStore()
	authMiddleware := security.GetAuthMiddleware()
	security.RegisterCallbackRoutes(router, dbClient)
	authApi := auth.NewApi(
		security.GetTokenService(),
		security.GetTokenProviders(),
		security.GetAuthProviders(),
		websocketHub,
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

		userApi := user.NewApi(user.NewStore(dbClient, security.GetPasswordEncoder(), websocketHub), actionLogger, logger,
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

		alarmStore := alarm.NewStore(dbClient, linksFetcher, logger)
		alarmAPI := alarm.NewApi(alarmStore, exportExecutor, timezoneConfigProvider, logger)
		alarmRouter := protected.Group("/alarms")
		{
			alarmRouter.GET(
				"",
				middleware.Authorize(authPermAlarmRead, permCan, enforcer),
				alarmAPI.List,
			)
			alarmRouter.GET(
				"/:id",
				middleware.Authorize(authPermAlarmRead, permCan, enforcer),
				alarmAPI.Get,
			)
		}
		protected.POST(
			"/alarm-details",
			middleware.Authorize(authPermAlarmRead, permCan, enforcer),
			alarmAPI.GetDetails,
		)
		protected.GET(
			"/manual-meta-alarms",
			middleware.Authorize(authPermAlarmRead, permCan, enforcer),
			alarmAPI.ListManual,
		)
		protected.GET(
			"/entityservice-alarms/:id",
			middleware.Authorize(authPermAlarmRead, permCan, enforcer),
			alarmAPI.ListByService,
		)
		protected.GET(
			"/component-alarms",
			middleware.Authorize(authPermAlarmRead, permCan, enforcer),
			alarmAPI.ListByComponent,
		)
		protected.GET(
			"/resolved-alarms",
			middleware.Authorize(authPermAlarmRead, permCan, enforcer),
			alarmAPI.ResolvedList,
		)
		protected.GET(
			"/open-alarms",
			middleware.Authorize(authPermAlarmRead, permCan, enforcer),
			alarmAPI.GetOpen,
		)
		protected.GET(
			"/alarm-counters",
			middleware.Authorize(authPermAlarmRead, permCan, enforcer),
			alarmAPI.Count,
		)
		alarmExportRouter := protected.Group("/alarm-export")
		{
			alarmExportRouter.POST(
				"",
				middleware.Authorize(authPermAlarmRead, permCan, enforcer),
				alarmAPI.StartExport,
			)
			alarmExportRouter.GET(
				"/:id/download",
				security.GetFileAuthMiddleware(),
				middleware.Authorize(authPermAlarmRead, permCan, enforcer),
				alarmAPI.DownloadExport,
			)
			alarmExportRouter.GET(
				"/:id",
				middleware.Authorize(authPermAlarmRead, permCan, enforcer),
				alarmAPI.GetExport,
			)
		}

		exportConfigurationAPI := exportconfiguration.NewApi(dbClient, logger)
		protected.POST(
			"/export-configuration",
			middleware.Authorize(apisecurity.PermExportConfigurations, permCan, enforcer),
			exportConfigurationAPI.Export,
		)

		entityAPI := entity.NewApi(
			entity.NewStore(dbClient, timezoneConfigProvider),
			exportExecutor,
			entityCleanerTaskChan,
			entityPublChan,
			metricsEntityMetaUpdater,
			actionLogger,
			logger,
		)

		entityExportRouter := protected.Group("/entity-export")
		{
			entityExportRouter.POST(
				"",
				middleware.Authorize(authObjEntity, permRead, enforcer),
				entityAPI.StartExport,
			)
			entityExportRouter.GET(
				"/:id/download",
				security.GetFileAuthMiddleware(),
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
				middleware.Authorize(authObjEntity, permRead, enforcer),
				entityAPI.List,
			)

			entityRouter.POST(
				"/clean",
				middleware.Authorize(authObjEntity, permDelete, enforcer),
				entityAPI.Clean,
			)

			entityRouter.GET(
				"/context-graph",
				middleware.Authorize(authObjEntity, permRead, enforcer),
				entityAPI.GetContextGraph,
			)

			entityRouter.GET(
				"/pbehaviors",
				middleware.Authorize(authObjEntity, permRead, enforcer),
				middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionRead, enforcer),
				pbehaviorApi.ListByEntityID,
			)

			entityRouter.GET(
				"/pbehavior-calendar",
				middleware.Authorize(authObjEntity, permRead, enforcer),
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

		entityserviceAPI := entityservice.NewApi(entityservice.NewStore(dbClient, linksFetcher, logger), entityPublChan,
			metricsEntityMetaUpdater, common.NewPatternFieldsTransformer(dbClient), actionLogger, logger)
		entityserviceRouter := protected.Group("/entityservices")
		{
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
			pbhTypeAuthorizeRead := middleware.Authorize(apisecurity.ObjPbehaviorType, model.PermissionRead, enforcer)
			pbhTypeAuthorizeCreate := middleware.Authorize(apisecurity.ObjPbehaviorType, model.PermissionCreate, enforcer)

			typeRouter.GET("", pbhTypeAuthorizeRead, pbehaviorTypeApi.List)
			typeRouter.POST("", pbhTypeAuthorizeCreate, pbehaviorTypeApi.Create)

			pbhTypeIDGroup := typeRouter.Group("")
			{
				pbhTypeAuthorizeUpdate := middleware.Authorize(apisecurity.ObjPbehaviorType, model.PermissionUpdate, enforcer)
				pbhTypeAuthorizeDelete := middleware.Authorize(apisecurity.ObjPbehaviorType, model.PermissionDelete, enforcer)

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
				linksFetcher,
				alarmStore,
				timezoneConfigProvider,
				logger,
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

		securityConfig := security.GetConfig().Security
		appInfoApi := appinfo.NewApi(enforcer, appinfo.NewStore(dbClient, securityConfig.AuthProviders,
			securityConfig.Cas.Title, securityConfig.Saml.Title))
		protected.GET("app-info", appInfoApi.GetAppInfo)
		appInfoRouter := protected.Group("/internal")
		{
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
			middleware.Authorize(apisecurity.PermHealthcheck, permCan, enforcer),
			engineinfo.GetRunInfo(ctx, runInfoManager),
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

		widgetAPI := widget.NewApi(widget.NewStore(dbClient), enforcer, common.NewPatternFieldsTransformer(dbClient), actionLogger)
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
			middleware.Authorize(apisecurity.ObjView, model.PermissionUpdate, enforcer),
			widgetFilterAPI.UpdatePositions,
		)

		viewGroupAPI := viewgroup.NewApi(viewgroup.NewStore(dbClient), actionLogger)
		viewGroupRouter := protected.Group("/view-groups")
		{
			viewGroupRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjViewGroup, permCreate, enforcer),
				middleware.SetAuthor(),
				viewGroupAPI.Create,
			)
			viewGroupRouter.GET(
				"",
				middleware.ProvideAuthorizedIds(permRead, enforcer),
				middleware.Authorize(apisecurity.ObjViewGroup, permRead, enforcer),
				viewGroupAPI.List,
			)
			viewGroupRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjViewGroup, permRead, enforcer),
				viewGroupAPI.Get,
			)
			viewGroupRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjViewGroup, permUpdate, enforcer),
				middleware.SetAuthor(),
				viewGroupAPI.Update,
			)
			viewGroupRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjViewGroup, permDelete, enforcer),
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

		scenarioAPI := scenario.NewApi(scenario.NewStore(dbClient), actionLogger, common.NewPatternFieldsTransformer(dbClient), logger, scenarioPriorityIntervals)
		scenarioRouter := protected.Group("/scenarios")
		{
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
			scenarioRouter.GET(
				"/minimal-priority",
				middleware.Authorize(authObjAction, permRead, enforcer),
				scenarioAPI.GetMinimalPriority,
			)
			scenarioRouter.POST(
				"/check-priority",
				middleware.Authorize(authObjAction, permRead, enforcer),
				scenarioAPI.CheckPriority,
			)
		}

		contextGraphAPIV1 := libcontextgraphV1.NewApi(conf, jobQueue, contextgraph.NewMongoStatusReporter(dbClient), logger)
		contextGraphRouter := protected.Group("/contextgraph")
		{
			contextGraphRouter.PUT(
				"import",
				middleware.Authorize(authObjContextGraph, permCreate, enforcer),
				contextGraphAPIV1.ImportAll,
			)
			contextGraphRouter.PUT(
				"import-partial",
				middleware.Authorize(authObjContextGraph, permCreate, enforcer),
				contextGraphAPIV1.ImportPartial,
			)
			contextGraphRouter.GET(
				"import/status/:id",
				middleware.Authorize(authObjContextGraph, permRead, enforcer),
				contextGraphAPIV1.Status,
			)
		}

		contextGraphAPIV2 := libcontextgraphV2.NewApi(conf, jobQueue, contextgraph.NewMongoStatusReporter(dbClient), logger)
		protected.PUT(
			"contextgraph-import",
			middleware.Authorize(authObjContextGraph, permCreate, enforcer),
			contextGraphAPIV2.ImportAll,
		)
		protected.PUT(
			"contextgraph-import-partial",
			middleware.Authorize(authObjContextGraph, permCreate, enforcer),
			contextGraphAPIV2.ImportPartial,
		)

		stateSettingsRouter := protected.Group("/state-settings")
		{
			stateSettingsApi := statesettings.NewApi(statesettings.NewStore(dbClient), actionLogger)
			stateSettingsRouter.PUT(
				"/:id",
				middleware.Authorize(authObjStateSettings, permCan, enforcer),
				stateSettingsApi.Update,
			)
			stateSettingsRouter.GET(
				"",
				middleware.Authorize(authObjStateSettings, permCan, enforcer),
				stateSettingsApi.List,
			)
		}

		notificationRouter := protected.Group("/notification")
		{
			notificationApi := notification.NewApi(notification.NewStore(dbClient), actionLogger)
			notificationRouter.PUT(
				"",
				middleware.Authorize(authObjNotification, permCan, enforcer),
				notificationApi.Update,
			)
			notificationRouter.GET(
				"",
				middleware.Authorize(authObjNotification, permCan, enforcer),
				notificationApi.Get,
			)
		}

		playlistRouter := protected.Group("/playlists")
		{
			playlistApi := playlist.NewApi(playlist.NewStore(dbClient), viewtab.NewStore(dbClient, widget.NewStore(dbClient)), enforcer, actionLogger)
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

		idleRuleAPI := idlerule.NewApi(idlerule.NewStore(dbClient), common.NewPatternFieldsTransformer(dbClient), actionLogger, logger)
		idleRuleRouter := protected.Group("/idle-rules")
		{
			idleRuleRouter.POST(
				"",
				middleware.Authorize(authObjIdleRule, permCreate, enforcer),
				middleware.SetAuthor(),
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
				middleware.SetAuthor(),
				idleRuleAPI.Update,
			)
			idleRuleRouter.DELETE(
				"/:id",
				middleware.Authorize(authObjIdleRule, permDelete, enforcer),
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
					middleware.Authorize(authObjAction, permCreate, enforcer),
					middleware.PreProcessBulk(conf, true),
					scenarioAPI.BulkCreate,
				)
				scenarioRouter.PUT(
					"",
					middleware.Authorize(authObjAction, permUpdate, enforcer),
					middleware.PreProcessBulk(conf, true),
					scenarioAPI.BulkUpdate,
				)
				scenarioRouter.DELETE(
					"",
					middleware.Authorize(authObjAction, permDelete, enforcer),
					middleware.PreProcessBulk(conf, false),
					scenarioAPI.BulkDelete,
				)
			}

			idleruleRouter := bulkRouter.Group("/idle-rules")
			{
				idleruleRouter.POST(
					"",
					middleware.Authorize(authObjIdleRule, permCreate, enforcer),
					middleware.PreProcessBulk(conf, true),
					idleRuleAPI.BulkCreate,
				)
				idleruleRouter.PUT(
					"",
					middleware.Authorize(authObjIdleRule, permUpdate, enforcer),
					middleware.PreProcessBulk(conf, true),
					idleRuleAPI.BulkUpdate,
				)
				idleruleRouter.DELETE(
					"",
					middleware.Authorize(authObjIdleRule, permDelete, enforcer),
					middleware.PreProcessBulk(conf, false),
					idleRuleAPI.BulkDelete,
				)
			}

			eventFilterRouter := bulkRouter.Group("/eventfilters")
			{
				eventFilterRouter.POST(
					"",
					middleware.Authorize(authEventFilter, permCreate, enforcer),
					middleware.PreProcessBulk(conf, true),
					eventFilterApi.BulkCreate,
				)
				eventFilterRouter.PUT(
					"",
					middleware.Authorize(authEventFilter, permUpdate, enforcer),
					middleware.PreProcessBulk(conf, true),
					eventFilterApi.BulkUpdate,
				)
				eventFilterRouter.DELETE(
					"",
					middleware.Authorize(authEventFilter, permDelete, enforcer),
					middleware.PreProcessBulk(conf, false),
					eventFilterApi.BulkDelete,
				)
			}

			entityserviceRouter := bulkRouter.Group("/entityservices")
			{
				entityserviceRouter.POST(
					"",
					middleware.Authorize(authObjEntityService, permCreate, enforcer),
					middleware.PreProcessBulk(conf, false),
					entityserviceAPI.BulkCreate,
				)
				entityserviceRouter.PUT(
					"",
					middleware.Authorize(authObjEntityService, permUpdate, enforcer),
					middleware.PreProcessBulk(conf, false),
					entityserviceAPI.BulkUpdate,
				)
				entityserviceRouter.DELETE(
					"",
					middleware.Authorize(authObjEntityService, permDelete, enforcer),
					middleware.PreProcessBulk(conf, false),
					entityserviceAPI.BulkDelete,
				)
			}

			userRouter := bulkRouter.Group("/users")
			{
				userRouter.POST(
					"",
					middleware.Authorize(apisecurity.PermAcl, permCreate, enforcer),
					middleware.PreProcessBulk(conf, false),
					userApi.BulkCreate,
				)
				userRouter.PUT(
					"",
					middleware.Authorize(apisecurity.PermAcl, permUpdate, enforcer),
					middleware.PreProcessBulk(conf, false),
					userApi.BulkUpdate,
				)
				userRouter.DELETE(
					"",
					middleware.Authorize(apisecurity.PermAcl, permDelete, enforcer),
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
		}

		dateStorageRouter := protected.Group("data-storage")
		{
			dateStorageAPI := datastorage.NewApi(datastorage.NewStore(dbClient, pgPoolProvider, logger))
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

		fileRouter := protected.Group("/file")
		{
			fileAPI := file.NewApi(enforcer, file.NewStore(dbClient, libfile.NewStorage(
				filesRoot,
				libfile.NewEtagEncoder(),
			), conf.File.UploadMaxSize))
			fileRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjFile, permCreate, enforcer),
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
				middleware.Authorize(apisecurity.ObjFile, permDelete, enforcer),
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
			middleware.Authorize(authObjEntity, permRead, enforcer),
			entityInfoDictionaryApi.ListKeys,
		)
		protected.GET("/entity-infos-dictionary/values",
			middleware.Authorize(authObjEntity, permRead, enforcer),
			entityInfoDictionaryApi.ListValues,
		)

		alarmTagRouter := protected.Group("/alarm-tags")
		{
			alarmTagAPI := alarmtag.NewApi(alarmtag.NewStore(dbClient))
			alarmTagRouter.GET(
				"",
				middleware.Authorize(authPermAlarmRead, permCan, enforcer),
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
