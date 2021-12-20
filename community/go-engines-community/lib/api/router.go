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
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/file"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/flappingrule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/idlerule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/messageratestats"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/middleware"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/notification"
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
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/statesettings"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/user"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/userpreferences"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/view"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/viewgroup"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/action"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	libentityservice "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	libpbehavior "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	libfile "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/file"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	libsecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

const baseUrl = "/api/v4"

const (
	authObjPbh          = apisecurity.ObjPbehavior
	authObjPbhType      = apisecurity.ObjPbehaviorType
	authObjPbhReason    = apisecurity.ObjPbehaviorReason
	authObjPbhException = apisecurity.ObjPbehaviorException

	authObjAction = apisecurity.ObjAction

	authObjEntity         = apisecurity.ObjEntity
	authObjEntityService  = apisecurity.ObjEntityService
	authObjEntityCategory = apisecurity.ObjEntityCategory
	authObjContextGraph   = apisecurity.ObjContextGraph

	authObjView      = apisecurity.ObjView
	authObjViewGroup = apisecurity.ObjViewGroup
	authObjPlaylist  = apisecurity.ObjPlaylist

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
	dbClient mongo.DbClient,
	timezoneConfigProvider config.TimezoneConfigProvider,
	pbhEntityTypeResolver libpbehavior.EntityTypeResolver,
	pbhComputeChan chan<- libpbehavior.ComputeTask,
	entityPublChan chan<- libentityservice.ChangeEntityMessage,
	entityCleanerTaskChan chan<- entity.CleanTask,
	runInfoManager engine.RunInfoManager,
	exportExecutor export.TaskExecutor,
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
	sessionStore := security.GetSessionStore()
	authMiddleware := security.GetAuthMiddleware()
	security.RegisterCallbackRoutes(router)
	authApi := auth.NewApi(
		security.GetTokenService(),
		security.GetTokenStore(),
		security.GetAuthProviders(),
		security.GetSessionStore(),
		websocketHub,
		security.GetCookieOptions().FileAccessName,
		security.GetCookieOptions().MaxAge,
		security.GetCookieOptions().Secure,
		logger,
	)
	sessionauthApi := sessionauth.NewApi(
		sessionStore,
		security.GetAuthProviders(),
		websocketHub,
		security.GetTokenStore(),
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

		protected.GET("/account/me", account.NewApi(account.NewStore(dbClient)).Me)
		protected.GET("/logged-user-count", authApi.GetLoggedUserCount)
		protected.GET("/file-access", authApi.GetFileAccess)

		userPreferencesRouter := protected.Group("/user-preferences")
		{
			userPreferencesRouter.Use(middleware.OnlyAuth())
			userPreferencesApi := userpreferences.NewApi(userpreferences.NewStore(dbClient), actionLogger)
			userPreferencesRouter.GET("/:id", userPreferencesApi.Get)
			userPreferencesRouter.PUT("", userPreferencesApi.Update)
		}

		userApi := user.NewApi(user.NewStore(dbClient, security.GetPasswordEncoder()), actionLogger, metricsUserMetaUpdater)
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

		entityAPI := entity.NewApi(entity.NewStore(dbClient), exportExecutor, entityCleanerTaskChan, logger)
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
			entityAPI := entity.NewApi(entity.NewStore(dbClient), exportExecutor, entityCleanerTaskChan, logger)
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
				"/pbehaviors",
				middleware.Authorize(authObjEntity, permRead, enforcer),
				middleware.Authorize(authObjPbh, permRead, enforcer),
				pbehaviorApi.ListByEntityID,
			)
		}
		entitybasicsRouter := protected.Group("/entitybasics")
		{
			entitybasicsAPI := entitybasic.NewApi(entitybasic.NewStore(dbClient), entityPublChan, metricsEntityMetaUpdater,
				actionLogger, logger)
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

		entityserviceAPI := entityservice.NewApi(entityservice.NewStore(dbClient), entityPublChan, metricsEntityMetaUpdater, actionLogger, logger)
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

		appInfoApi := appinfo.NewApi(enforcer, appinfo.NewStore(dbClient, security.GetConfig().Security.AuthProviders))
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

		scenarioAPI := scenario.NewApi(scenario.NewStore(dbClient), actionLogger, scenarioPriorityIntervals)
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

		idleRuleAPI := idlerule.NewApi(idlerule.NewStore(dbClient), actionLogger, userInterfaceConfig)
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
			idleRuleRouter.POST(
				"/count",
				middleware.Authorize(authObjPbh, permCreate, enforcer),
				idleRuleAPI.CountPatterns)
		}

		bulkRouter := protected.Group("/bulk")
		{
			scenarioRouter := bulkRouter.Group("/scenarios")
			{
				scenarioRouter.POST(
					"",
					middleware.Authorize(authObjView, permCreate, enforcer),
					middleware.PreProcessBulk(conf),
					scenarioAPI.BulkCreate,
				)
				scenarioRouter.PUT(
					"",
					middleware.Authorize(authObjView, permUpdate, enforcer),
					middleware.PreProcessBulk(conf),
					scenarioAPI.BulkUpdate,
				)
				scenarioRouter.DELETE(
					"",
					middleware.Authorize(authObjView, permDelete, enforcer),
					middleware.PreProcessBulk(conf),
					scenarioAPI.BulkDelete,
				)
			}

			idleruleRouter := bulkRouter.Group("/idle-rules")
			{
				idleruleRouter.POST(
					"",
					middleware.Authorize(authObjView, permCreate, enforcer),
					middleware.PreProcessBulk(conf),
					idleRuleAPI.BulkCreate,
				)
				idleruleRouter.PUT(
					"",
					middleware.Authorize(authObjView, permUpdate, enforcer),
					middleware.PreProcessBulk(conf),
					idleRuleAPI.BulkUpdate,
				)
				idleruleRouter.DELETE(
					"",
					middleware.Authorize(authObjView, permDelete, enforcer),
					middleware.PreProcessBulk(conf),
					idleRuleAPI.BulkDelete,
				)
			}

			eventFilterRouter := bulkRouter.Group("/eventfilters")
			{
				eventFilterRouter.POST(
					"",
					middleware.Authorize(authObjView, permCreate, enforcer),
					middleware.PreProcessBulk(conf),
					eventFilterApi.BulkCreate,
				)
				eventFilterRouter.PUT(
					"",
					middleware.Authorize(authObjView, permUpdate, enforcer),
					middleware.PreProcessBulk(conf),
					eventFilterApi.BulkUpdate,
				)
				eventFilterRouter.DELETE(
					"",
					middleware.Authorize(authObjView, permDelete, enforcer),
					middleware.PreProcessBulk(conf),
					eventFilterApi.BulkDelete,
				)
			}

			entityserviceRouter := bulkRouter.Group("/entityservices")
			{
				entityserviceRouter.POST(
					"",
					middleware.Authorize(authObjView, permCreate, enforcer),
					middleware.PreProcessBulk(conf),
					entityserviceAPI.BulkCreate,
				)
				entityserviceRouter.PUT(
					"",
					middleware.Authorize(authObjView, permUpdate, enforcer),
					middleware.PreProcessBulk(conf),
					entityserviceAPI.BulkUpdate,
				)
				entityserviceRouter.DELETE(
					"",
					middleware.Authorize(authObjView, permDelete, enforcer),
					middleware.PreProcessBulk(conf),
					entityserviceAPI.BulkDelete,
				)
			}

			userRouter := bulkRouter.Group("/users")
			{
				userRouter.POST(
					"",
					middleware.Authorize(authObjView, permCreate, enforcer),
					userApi.BulkCreate,
				)
				userRouter.PUT(
					"",
					middleware.Authorize(authObjView, permUpdate, enforcer),
					userApi.BulkUpdate,
				)
				userRouter.DELETE(
					"",
					middleware.Authorize(authObjView, permDelete, enforcer),
					userApi.BulkDelete,
				)
			}

			viewRouter := bulkRouter.Group("/views")
			{
				viewRouter.POST(
					"",
					middleware.Authorize(authObjView, permCreate, enforcer),
					middleware.PreProcessBulk(conf),
					viewAPI.BulkCreate,
					middleware.ReloadEnforcerPolicyOnChange(enforcer),
				)
				viewRouter.PUT(
					"",
					middleware.Authorize(authObjView, permUpdate, enforcer),
					middleware.ProvideAuthorizedIds(permUpdate, enforcer),
					middleware.PreProcessBulk(conf),
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
					middleware.PreProcessBulk(conf),
					viewGroupAPI.BulkCreate,
				)
				viewGroupRouter.PUT(
					"",
					middleware.Authorize(authObjViewGroup, permUpdate, enforcer),
					middleware.PreProcessBulk(conf),
					viewGroupAPI.BulkUpdate,
				)
				viewGroupRouter.DELETE(
					"",
					middleware.Authorize(authObjViewGroup, permDelete, enforcer),
					viewGroupAPI.BulkDelete,
				)
			}

			pbehaviorRouter := bulkRouter.Group("/pbehaviors")
			{
				pbehaviorRouter.POST(
					"",
					middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionCreate, enforcer),
					middleware.PreProcessBulk(conf),
					pbehaviorApi.BulkCreate,
				)
				pbehaviorRouter.PUT(
					"",
					middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionUpdate, enforcer),
					middleware.PreProcessBulk(conf),
					pbehaviorApi.BulkUpdate,
				)
				pbehaviorRouter.DELETE(
					"",
					middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionDelete, enforcer),
					pbehaviorApi.BulkDelete,
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
			getFileRouter := fileRouter.Group("", security.GetFileAuthMiddleware()...)
			getFileRouter.GET(
				"",
				fileAPI.List,
			)
			getFileRouter.GET(
				"/:id",
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
			resolveRuleAPI := resolverule.NewApi(resolverule.NewStore(dbClient), actionLogger)
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
			flappingRuleAPI := flappingrule.NewApi(flappingrule.NewStore(dbClient), actionLogger)
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
	}
}
