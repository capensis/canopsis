package api

import (
	"context"
	"path/filepath"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/account"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/alarmaction"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/alarmtag"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/appinfo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/associativetable"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/broadcastmessage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/colortheme"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/contextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/dbexport"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entitybasic"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entitycategory"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entitycomment"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entityinfodictionary"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entityinfosproperty"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/export"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/exportconfiguration"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/externaldatatable"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/file"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/flappingrule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/healthcheck"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/icon"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/idlerule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/linkrule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/maintenance"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/messageratestats"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/middleware"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/notification"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/patternfields"
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
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/workers"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	libentity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	libentityservice "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	libevent "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/link"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	libpbehavior "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statesetting"
	libtemplate "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template/validator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/usernotification"
	libfile "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/file"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	libsecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/userprovider"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

const BaseUrl = "/api/v4"

const mimeTypeSvg = "image/svg+xml"

const cacheExpiration = time.Second

// RegisterRoutes
// nolint: contextcheck
func RegisterRoutes(
	ctx context.Context,
	conf config.CanopsisConf,
	router gin.IRouter,
	security Security,
	enforcer libsecurity.Enforcer,
	errorResponder httperror.Responder,
	linkGenerator link.Generator,
	primaryDbClient mongo.DbClient,
	secondaryDbClient mongo.DbClient,
	dbExportClient mongo.DbClient,
	pgPoolProvider postgres.PoolProvider,
	amqpPublisher amqp.Channel,
	lockRedisSession redis.Cmdable,
	apiConfigProvider config.ApiConfigProvider,
	timezoneConfigProvider config.TimezoneConfigProvider,
	templateConfigProvider config.TemplateConfigProvider,
	pbhEntityTypeResolver libpbehavior.EntityTypeResolver,
	pbhComputeChan chan<- rpc.PbehaviorRecomputeEvent,
	entityPublChan chan<- libentityservice.ChangeEntityMessage,
	entityCleanerTaskChan chan<- libentity.CleanTask,
	exportTaskExecutor export.TaskExecutor,
	techMetricsTaskExecutor techmetrics.TaskExecutor,
	userInterfaceConfig config.UserInterfaceConfigProvider,
	websocketHub websocket.Hub,
	websocketStore websocket.Store,
	broadcastMessageChan chan<- bool,
	metricsEntityMetaUpdater metrics.MetaUpdater,
	metricsUserMetaUpdater metrics.MetaUpdater,
	authorProvider author.Provider,
	healthcheckStore healthcheck.Store,
	tplExecutor libtemplate.Executor,
	stateSettingsUpdatesChan chan statesetting.RuleUpdatedMessage,
	enableSameServiceNames bool,
	eventGenerator libevent.Generator,
	securityConfig libsecurity.Config,
	exdataImportWorker externaldatatable.ImportWorker,
	notifStore usernotification.Store,
	externalDataContainer *externaldata.GetterContainer,
	tplTestTypePermMapping map[int][]any,
	logger zerolog.Logger,
) error {
	sessionStore := security.GetSessionStore()
	authMiddleware := security.GetAuthMiddleware()
	err := security.RegisterCallbackRoutes(ctx, router, primaryDbClient, sessionStore)
	if err != nil {
		return err
	}

	maintenanceAdapter := config.NewMaintenanceAdapter(primaryDbClient)
	userInterfaceAdapter := config.NewUserInterfaceAdapter(primaryDbClient)

	authApi := auth.NewApi(
		security.GetTokenService(),
		security.GetTokenProviders(),
		security.GetAuthProviders(),
		websocketStore,
		maintenanceAdapter,
		enforcer,
		security.GetCookieOptions().FileAccessName,
		security.GetCookieOptions().MaxAge,
		errorResponder,
		logger,
	)
	sessionauthApi := sessionauth.NewApi(
		sessionStore,
		security.GetAuthProviders(),
		maintenanceAdapter,
		enforcer,
		errorResponder,
		logger,
	)
	router.POST("/auth", sessionauthApi.LoginHandler())

	sessionProtected := router.Group("")
	{
		sessionProtected.Use(middleware.SessionAuth(primaryDbClient, apiConfigProvider, sessionStore, errorResponder), middleware.OnlyAuth(errorResponder))
		sessionProtected.GET("/logout", sessionauthApi.LogoutHandler())
	}

	unprotected := router.Group(BaseUrl)
	{
		unprotected.POST("/login", authApi.Login)
		unprotected.POST("/logout", authApi.Logout)
	}

	protected := router.Group(BaseUrl)
	{
		protected.Use(authMiddleware...)

		protected.Group("/ws").GET("", websocket.NewApi(websocketHub, errorResponder).Handler)

		accountRouter := protected.Group("/account/me")
		{
			accountRouter.Use(middleware.OnlyAuth(errorResponder))
			accountAPI := account.NewApi(
				account.NewStore(primaryDbClient, security.GetPasswordEncoder(), authorProvider, userInterfaceAdapter, securityConfig),
				errorResponder)
			accountRouter.GET("", accountAPI.Me)
			accountRouter.PUT("", accountAPI.Update)
		}
		protected.GET("/logged-user-count", authApi.GetLoggedUserCount)
		protected.GET("/file-access", authApi.GetFileAccess)

		userPreferencesRouter := protected.Group("/user-preferences")
		{
			userPreferencesRouter.Use(middleware.OnlyAuth(errorResponder))
			userPreferencesApi := userpreferences.NewApi(userpreferences.NewStore(primaryDbClient, authorProvider),
				widget.NewStore(primaryDbClient, authorProvider, enforcer,
					patternfields.NewTransformer(primaryDbClient), validator.NewValidator(tplExecutor),
					templateConfigProvider), enforcer, errorResponder)
			userPreferencesRouter.GET("/:id", userPreferencesApi.Get)
			userPreferencesRouter.PUT("", userPreferencesApi.Update)
		}

		userApi := user.NewApi(user.NewStore(primaryDbClient, security.GetPasswordEncoder(), websocketStore, authorProvider, securityConfig),
			metricsUserMetaUpdater, errorResponder)
		userRouter := protected.Group("/users")
		{
			userRouter.POST("",
				middleware.Authorize(apisecurity.PermAcl, model.PermissionCreate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				userApi.Create,
				middleware.ReloadEnforcerPolicyOnChange(enforcer, errorResponder),
			)
			userRouter.GET("",
				middleware.Authorize(apisecurity.PermAcl, model.PermissionRead, enforcer, errorResponder),
				userApi.List,
			)
			userRouter.GET("/:id",
				middleware.Authorize(apisecurity.PermAcl, model.PermissionRead, enforcer, errorResponder),
				userApi.Get,
			)
			userRouter.PUT("/:id",
				middleware.Authorize(apisecurity.PermAcl, model.PermissionUpdate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				userApi.Update,
				middleware.ReloadEnforcerPolicyOnChange(enforcer, errorResponder),
			)
			userRouter.PATCH("/:id",
				middleware.Authorize(apisecurity.PermAcl, model.PermissionUpdate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				userApi.Patch,
				middleware.ReloadEnforcerPolicyOnChange(enforcer, errorResponder),
			)
			userRouter.DELETE("/:id",
				middleware.Authorize(apisecurity.PermAcl, model.PermissionDelete, enforcer, errorResponder),
				userApi.Delete,
			)
		}
		roleApi := role.NewApi(role.NewStore(primaryDbClient, authorProvider), errorResponder)
		roleRouter := protected.Group("/roles")
		{
			roleRouter.POST("",
				middleware.Authorize(apisecurity.PermAcl, model.PermissionCreate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				roleApi.Create,
			)
			roleRouter.GET("",
				middleware.Authorize(apisecurity.PermAcl, model.PermissionRead, enforcer, errorResponder),
				roleApi.List,
			)
			roleRouter.GET("/:id",
				middleware.Authorize(apisecurity.PermAcl, model.PermissionRead, enforcer, errorResponder),
				roleApi.Get,
			)
			roleRouter.PUT("/:id",
				middleware.Authorize(apisecurity.PermAcl, model.PermissionUpdate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				roleApi.Update,
				middleware.ReloadEnforcerPolicyOnChange(enforcer, errorResponder),
			)
			roleRouter.DELETE("/:id",
				middleware.Authorize(apisecurity.PermAcl, model.PermissionDelete, enforcer, errorResponder),
				roleApi.Delete,
			)
		}
		protected.GET("/role-templates",
			middleware.Authorize(apisecurity.PermAcl, model.PermissionRead, enforcer, errorResponder),
			roleApi.ListTemplates,
		)
		permissionRouter := protected.Group("/permissions")
		{
			permissionApi := permission.NewApi(permission.NewStore(primaryDbClient), errorResponder)
			permissionRouter.GET("",
				middleware.Authorize(apisecurity.PermAcl, model.PermissionRead, enforcer, errorResponder),
				permissionApi.List,
			)
		}

		sharetokenApi := sharetoken.NewApi(
			sharetoken.NewStore(primaryDbClient, security.GetTokenGenerator(), authorProvider), errorResponder)
		sharetokenRouter := protected.Group("/share-tokens")
		{
			sharetokenRouter.POST("",
				middleware.Authorize(apisecurity.PermShareToken, model.PermissionCreate, enforcer, errorResponder),
				sharetokenApi.Create,
			)
			sharetokenRouter.GET("",
				middleware.Authorize(apisecurity.PermShareToken, model.PermissionRead, enforcer, errorResponder),
				sharetokenApi.List,
			)
			sharetokenRouter.DELETE("/:id",
				middleware.Authorize(apisecurity.PermShareToken, model.PermissionDelete, enforcer, errorResponder),
				sharetokenApi.Delete,
			)
		}

		alarmStore := alarm.NewStore(secondaryDbClient, dbExportClient, linkGenerator, patternfields.NewTransformer(primaryDbClient),
			timezoneConfigProvider, authorProvider, tplExecutor, json.NewDecoder(), logger)
		alarmAPI := alarm.NewApi(alarmStore, exportTaskExecutor, json.NewEncoder(), errorResponder, logger)
		alarmActionAPI := alarmaction.NewApi(alarmaction.NewStore(primaryDbClient, amqpPublisher, canopsis.DefaultExchangeName,
			canopsis.FIFOQueueName, json.NewEncoder(), canopsis.JsonContentType, eventGenerator, logger),
			errorResponder)
		alarmRouter := protected.Group("/alarms")
		{
			alarmRouter.GET(
				"",
				middleware.Authorize(apisecurity.PermAlarmRead, model.PermissionCan, enforcer, errorResponder),
				alarmAPI.List,
			)
			alarmRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.PermAlarmRead, model.PermissionCan, enforcer, errorResponder),
				alarmAPI.Get,
			)
			alarmRouter.PUT(
				"/:id/ack",
				middleware.Authorize(apisecurity.PermAlarmUpdate, model.PermissionCan, enforcer, errorResponder),
				alarmActionAPI.Ack,
			)
			alarmRouter.PUT(
				"/:id/ackremove",
				middleware.Authorize(apisecurity.PermAlarmUpdate, model.PermissionCan, enforcer, errorResponder),
				alarmActionAPI.AckRemove,
			)
			alarmRouter.PUT(
				"/:id/snooze",
				middleware.Authorize(apisecurity.PermAlarmUpdate, model.PermissionCan, enforcer, errorResponder),
				alarmActionAPI.Snooze,
			)
			alarmRouter.PUT(
				"/:id/cancel",
				middleware.Authorize(apisecurity.PermAlarmUpdate, model.PermissionCan, enforcer, errorResponder),
				alarmActionAPI.Cancel,
			)
			alarmRouter.PUT(
				"/:id/uncancel",
				middleware.Authorize(apisecurity.PermAlarmUpdate, model.PermissionCan, enforcer, errorResponder),
				alarmActionAPI.Uncancel,
			)
			alarmRouter.PUT(
				"/:id/assocticket",
				middleware.Authorize(apisecurity.PermAlarmUpdate, model.PermissionCan, enforcer, errorResponder),
				alarmActionAPI.AssocTicket,
			)
			alarmRouter.PUT(
				"/:id/comment",
				middleware.Authorize(apisecurity.PermAlarmUpdate, model.PermissionCan, enforcer, errorResponder),
				alarmActionAPI.Comment,
			)
			alarmRouter.PUT(
				"/:id/changestate",
				middleware.Authorize(apisecurity.PermAlarmUpdate, model.PermissionCan, enforcer, errorResponder),
				alarmActionAPI.ChangeState,
			)
			alarmRouter.PUT(
				"/:id/bookmark",
				middleware.Authorize(apisecurity.PermAlarmUpdate, model.PermissionCan, enforcer, errorResponder),
				alarmActionAPI.AddBookmark,
			)
			alarmRouter.DELETE(
				"/:id/bookmark",
				middleware.Authorize(apisecurity.PermAlarmUpdate, model.PermissionCan, enforcer, errorResponder),
				alarmActionAPI.RemoveBookmark,
			)
		}
		protected.POST(
			"/alarm-details",
			middleware.Authorize(apisecurity.PermAlarmRead, model.PermissionCan, enforcer, errorResponder),
			alarmAPI.GetDetails,
		)
		protected.GET(
			"/alarm-links/:id",
			middleware.Authorize(apisecurity.PermAlarmRead, model.PermissionCan, enforcer, errorResponder),
			alarmAPI.GetLinks,
		)
		protected.GET(
			"/entityservice-alarms/:id",
			middleware.Authorize(apisecurity.PermAlarmRead, model.PermissionCan, enforcer, errorResponder),
			alarmAPI.ListByService,
		)
		protected.GET(
			"/component-alarms",
			middleware.Authorize(apisecurity.PermAlarmRead, model.PermissionCan, enforcer, errorResponder),
			alarmAPI.ListByComponent,
		)
		protected.GET(
			"/resolved-alarms",
			middleware.Authorize(apisecurity.PermAlarmRead, model.PermissionCan, enforcer, errorResponder),
			alarmAPI.ResolvedList,
		)
		protected.GET(
			"/open-alarms",
			middleware.Authorize(apisecurity.PermAlarmRead, model.PermissionCan, enforcer, errorResponder),
			alarmAPI.GetOpen,
		)
		protected.GET(
			"/alarm-counters",
			middleware.Authorize(apisecurity.PermAlarmRead, model.PermissionCan, enforcer, errorResponder),
			alarmAPI.Count,
		)
		protected.GET(
			"/alarm-display-names",
			middleware.Authorize(apisecurity.PermAlarmRead, model.PermissionCan, enforcer, errorResponder),
			alarmAPI.GetDisplayNames,
		)
		exportTaskExecutor.RegisterType("alarm", alarmStore.Export)
		alarmExportRouter := protected.Group("/alarm-export")
		{
			alarmExportRouter.POST(
				"",
				middleware.Authorize(apisecurity.PermAlarmRead, model.PermissionCan, enforcer, errorResponder),
				alarmAPI.StartExport,
			)
			alarmExportRouter.GET(
				"/:id/download",
				security.GetFileAuthMiddleware(),
				middleware.Authorize(apisecurity.PermAlarmRead, model.PermissionCan, enforcer, errorResponder),
				alarmAPI.DownloadExport,
			)
			alarmExportRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.PermAlarmRead, model.PermissionCan, enforcer, errorResponder),
				alarmAPI.GetExport,
			)
		}

		exportConfigurationAPI := exportconfiguration.NewApi(primaryDbClient, errorResponder, logger)
		protected.POST(
			"/export-configuration",
			middleware.Authorize(apisecurity.PermExportConfigurations, model.PermissionCan, enforcer, errorResponder),
			exportConfigurationAPI.Export,
		)

		entityStore := entity.NewStore(primaryDbClient, dbExportClient, timezoneConfigProvider, authorProvider, patternfields.NewTransformer(primaryDbClient), json.NewDecoder())
		entityAPI := entity.NewApi(
			entityStore,
			exportTaskExecutor,
			entityCleanerTaskChan,
			entityPublChan,
			metricsEntityMetaUpdater,
			json.NewEncoder(),
			errorResponder,
			logger,
		)

		exportTaskExecutor.RegisterType("entity", entityStore.Export)
		entityExportRouter := protected.Group("/entity-export")
		{
			entityExportRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjEntity, model.PermissionRead, enforcer, errorResponder),
				entityAPI.StartExport,
			)
			entityExportRouter.GET(
				"/:id/download",
				security.GetFileAuthMiddleware(),
				middleware.Authorize(apisecurity.ObjEntity, model.PermissionRead, enforcer, errorResponder),
				entityAPI.DownloadExport,
			)
			entityExportRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjEntity, model.PermissionRead, enforcer, errorResponder),
				entityAPI.GetExport,
			)
		}

		patternAPI := pattern.NewApi(
			pattern.NewStore(primaryDbClient, secondaryDbClient, pbhComputeChan, entityPublChan, stateSettingsUpdatesChan, authorProvider, patternfields.NewTransformer(primaryDbClient), logger),
			userInterfaceConfig, enforcer, errorResponder, patternfields.NewFieldGetter(secondaryDbClient))
		patternRouter := protected.Group("/patterns")
		{
			patternRouter.Use(middleware.OnlyAuth(errorResponder))
			patternRouter.POST(
				"",
				middleware.SetAuthor(errorResponder),
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
				middleware.SetAuthor(errorResponder),
				patternAPI.Update,
			)
			patternRouter.DELETE(
				"/:id",
				patternAPI.Delete,
			)
		}
		protected.POST(
			"/patterns-alarms-count",
			middleware.OnlyAuth(errorResponder),
			patternAPI.CountAlarms,
		)
		protected.POST(
			"/patterns-entities-count",
			middleware.OnlyAuth(errorResponder),
			patternAPI.CountEntities,
		)

		protected.POST(
			"/pbehavior-timespans",
			middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionRead, enforcer, errorResponder),
			pbehaviortimespan.GetTimeSpans(
				pbehaviortimespan.NewService(primaryDbClient, timezoneConfigProvider), errorResponder),
		)
		protected.GET(
			"/pbehavior-ics/:id",
			middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionRead, enforcer, errorResponder),
			pbehaviorics.GetICS(pbehaviorics.NewStore(primaryDbClient, authorProvider),
				pbehaviorics.NewService(timezoneConfigProvider), errorResponder),
		)

		// event-filter API
		eventFilterApi := eventfilter.NewApi(
			eventfilter.NewStore(primaryDbClient, authorProvider, patternfields.NewTransformer(primaryDbClient),
				notifStore, validator.NewValidator(tplExecutor), tplExecutor, templateConfigProvider, externalDataContainer,
				json.NewEncoder(), json.NewDecoder()),
			dbexport.NewExporter(primaryDbClient),
			errorResponder,
		)
		eventFilterRouter := protected.Group("/eventfilter/rules")
		{
			eventFilterRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjEventFilterRule, model.PermissionCreate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				eventFilterApi.Create)
			eventFilterRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjEventFilterRule, model.PermissionRead, enforcer, errorResponder),
				eventFilterApi.Get)
			eventFilterRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjEventFilterRule, model.PermissionDelete, enforcer, errorResponder),
				eventFilterApi.Delete)
			eventFilterRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjEventFilterRule, model.PermissionRead, enforcer, errorResponder),
				eventFilterApi.List)
			eventFilterRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjEventFilterRule, model.PermissionUpdate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				eventFilterApi.Update)
		}
		protected.GET(
			"/eventfilter/:id/failures",
			middleware.Authorize(apisecurity.ObjEventFilterRule, model.PermissionRead, enforcer, errorResponder),
			eventFilterApi.ListFailures)
		protected.PUT(
			"/eventfilter/:id/failures",
			middleware.Authorize(apisecurity.ObjEventFilterRule, model.PermissionCreate, enforcer, errorResponder),
			eventFilterApi.ReadFailures)
		protected.POST(
			"eventfilter-db-export",
			middleware.Authorize(apisecurity.ObjEventFilterRule, model.PermissionRead, enforcer, errorResponder),
			eventFilterApi.DBExport)
		protected.POST(
			"/eventfilter-template-validate",
			middleware.Authorize(apisecurity.ObjEventFilterRule, model.PermissionRead, enforcer, errorResponder),
			eventFilterApi.ValidateTemplates)
		protected.GET(
			"/eventfilter-template-vars",
			middleware.Authorize(apisecurity.ObjEventFilterRule, model.PermissionRead, enforcer, errorResponder),
			eventFilterApi.GetTemplateVars)
		protected.GET(
			"/eventfilter-copy-vars",
			middleware.Authorize(apisecurity.ObjEventFilterRule, model.PermissionRead, enforcer, errorResponder),
			eventFilterApi.GetCopyVars)

		pbehaviorApi := pbehavior.NewApi(
			pbehavior.NewStore(
				primaryDbClient,
				secondaryDbClient,
				lockRedisSession,
				pbhEntityTypeResolver,
				libpbehavior.NewTypeComputer(libpbehavior.NewModelProvider(primaryDbClient, authorProvider), json.NewDecoder()),
				timezoneConfigProvider,
				authorProvider,
				patternfields.NewTransformer(primaryDbClient),
				websocketHub,
				userInterfaceConfig,
			),
			dbexport.NewExporter(primaryDbClient),
			pbhComputeChan,
			workers.NewJobPublisher(jobKeyPbhPatterns, amqpPublisher),
			errorResponder,
		)
		pbehaviorRouter := protected.Group("/pbehaviors")
		{
			pbehaviorRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionCreate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				pbehaviorApi.Create)
			pbehaviorRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionRead, enforcer, errorResponder),
				pbehaviorApi.List)
			pbehaviorRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionRead, enforcer, errorResponder),
				pbehaviorApi.Get)
			pbehaviorRouter.GET(
				"/:id/entities",
				middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionRead, enforcer, errorResponder),
				pbehaviorApi.ListEntities)
			pbehaviorRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionUpdate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				pbehaviorApi.Update)
			pbehaviorRouter.PATCH(
				"/:id",
				middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionUpdate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				pbehaviorApi.Patch)
			pbehaviorRouter.DELETE(
				"",
				middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionDelete, enforcer, errorResponder),
				pbehaviorApi.DeleteByName)
			pbehaviorRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionDelete, enforcer, errorResponder),
				pbehaviorApi.Delete)
		}
		protected.POST(
			"/pbehaviors-db-export",
			middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionRead, enforcer, errorResponder),
			pbehaviorApi.DBExport,
		)
		protected.PUT(
			"/pbehavior-patterns",
			middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionUpdate, enforcer, errorResponder),
			pbehaviorApi.ExecPattern,
		)
		protected.PUT(
			"/all-pbehavior-patterns",
			middleware.Authorize(apisecurity.PermPbhPatterns, model.PermissionCan, enforcer, errorResponder),
			pbehaviorApi.ExecAllPatterns,
		)

		pbehaviorCommentRouter := protected.Group("/pbehavior-comments")
		{
			pbehaviorCommentAPI := pbehaviorcomment.NewApi(
				pbehaviorcomment.NewStore(primaryDbClient, authorProvider), errorResponder)
			pbehaviorCommentRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionUpdate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				pbehaviorCommentAPI.Create,
			)
			pbehaviorCommentRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionUpdate, enforcer, errorResponder),
				pbehaviorCommentAPI.Delete,
			)
		}
		entityRouter := protected.Group("/entities")
		{
			entityRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjEntity, model.PermissionRead, enforcer, errorResponder),
				entityAPI.List,
			)

			entityRouter.POST(
				"/archive-disabled",
				middleware.Authorize(apisecurity.ObjEntity, model.PermissionDelete, enforcer, errorResponder),
				entityAPI.ArchiveDisabled,
			)
			entityRouter.POST(
				"/archive-unlinked",
				middleware.Authorize(apisecurity.ObjEntity, model.PermissionDelete, enforcer, errorResponder),
				entityAPI.ArchiveUnlinked,
			)
			entityRouter.POST(
				"/clean-archived",
				middleware.Authorize(apisecurity.ObjEntity, model.PermissionDelete, enforcer, errorResponder),
				entityAPI.CleanArchived,
			)

			entityRouter.GET(
				"/context-graph",
				middleware.Authorize(apisecurity.ObjEntity, model.PermissionRead, enforcer, errorResponder),
				entityAPI.GetContextGraph,
			)

			entityRouter.POST(
				"/check-state-setting",
				middleware.Authorize(apisecurity.ObjStateSettings, model.PermissionRead, enforcer, errorResponder),
				entityAPI.CheckStateSetting,
			)

			entityRouter.GET(
				"/state-setting",
				middleware.Authorize(apisecurity.ObjStateSettings, model.PermissionRead, enforcer, errorResponder),
				entityAPI.GetStateSetting,
			)

			entityRouter.GET(
				"/pbehaviors",
				middleware.Authorize(apisecurity.ObjEntity, model.PermissionRead, enforcer, errorResponder),
				middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionRead, enforcer, errorResponder),
				pbehaviorApi.ListByEntityID,
			)

			entityRouter.GET(
				"/pbehavior-calendar",
				middleware.Authorize(apisecurity.ObjEntity, model.PermissionRead, enforcer, errorResponder),
				middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionRead, enforcer, errorResponder),
				pbehaviorApi.CalendarByEntityID,
			)
		}

		entitybasicsAPI := entitybasic.NewApi(entitybasic.NewStore(primaryDbClient), entityPublChan,
			metricsEntityMetaUpdater, errorResponder, logger)
		entitybasicsRouter := protected.Group("/entitybasics")
		{
			entitybasicsRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjEntity, model.PermissionRead, enforcer, errorResponder),
				entitybasicsAPI.Get,
			)
			entitybasicsRouter.PUT(
				"",
				middleware.Authorize(apisecurity.ObjEntity, model.PermissionUpdate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				entitybasicsAPI.Update,
			)
			entitybasicsRouter.DELETE(
				"",
				middleware.Authorize(apisecurity.ObjEntity, model.PermissionDelete, enforcer, errorResponder),
				entitybasicsAPI.Delete,
			)
		}

		entityserviceStore := entityservice.NewStore(primaryDbClient, linkGenerator, enableSameServiceNames, authorProvider,
			patternfields.NewTransformer(primaryDbClient), validator.NewValidator(tplExecutor), templateConfigProvider,
			logger)
		entityserviceAPI := entityservice.NewApi(entityserviceStore, entityPublChan, metricsEntityMetaUpdater,
			errorResponder, logger)
		entityserviceRouter := protected.Group("/entityservices")
		{
			entityserviceRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjEntityService, model.PermissionCreate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				entityserviceAPI.Create,
			)
			entityserviceRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjEntityService, model.PermissionRead, enforcer, errorResponder),
				entityserviceAPI.Get,
			)
			entityserviceRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjEntityService, model.PermissionUpdate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				entityserviceAPI.Update,
			)
			entityserviceRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjEntityService, model.PermissionDelete, enforcer, errorResponder),
				entityserviceAPI.Delete,
			)
			protected.GET(
				"/entityservice-dependencies",
				middleware.Authorize(apisecurity.ObjEntityService, model.PermissionRead, enforcer, errorResponder),
				entityserviceAPI.GetDependencies,
			)
			protected.GET(
				"/entityservice-impacts",
				middleware.Authorize(apisecurity.ObjEntityService, model.PermissionRead, enforcer, errorResponder),
				entityserviceAPI.GetImpacts,
			)
			protected.POST(
				"/entityservice-template-validate",
				middleware.Authorize(apisecurity.ObjEntityService, model.PermissionRead, enforcer, errorResponder),
				entityserviceAPI.ValidateTemplates)
			protected.GET(
				"/entityservice-template-vars",
				middleware.Authorize(apisecurity.ObjEntityService, model.PermissionRead, enforcer, errorResponder),
				entityserviceAPI.GetTemplateVars)
		}

		entityCommentRouter := protected.Group("/entity-comments")
		{
			entityCommentAPI := entitycomment.NewApi(entitycomment.NewStore(primaryDbClient, logger), errorResponder)
			entityCommentRouter.POST(
				"",
				middleware.Authorize(apisecurity.PermEntityComment, model.PermissionCan, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				entityCommentAPI.Create,
			)
			entityCommentRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.PermEntityComment, model.PermissionCan, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				entityCommentAPI.Update,
			)
			entityCommentRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjEntity, model.PermissionRead, enforcer, errorResponder),
				entityCommentAPI.List,
			)
		}

		typeRouter := protected.Group("/pbehavior-types")
		{
			pbehaviorTypeApi := pbehaviortype.NewApi(
				pbehaviortype.NewStore(primaryDbClient, authorProvider),
				pbhComputeChan,
				errorResponder,
				logger,
			)
			pbhTypeAuthorizeRead := middleware.Authorize(apisecurity.ObjPbehaviorType, model.PermissionRead, enforcer, errorResponder)
			pbhTypeAuthorizeCreate := middleware.Authorize(apisecurity.ObjPbehaviorType, model.PermissionCreate, enforcer, errorResponder)
			pbhTypeAuthorizeUpdate := middleware.Authorize(apisecurity.ObjPbehaviorType, model.PermissionUpdate, enforcer, errorResponder)
			pbhTypeAuthorizeDelete := middleware.Authorize(apisecurity.ObjPbehaviorType, model.PermissionDelete, enforcer, errorResponder)

			typeRouter.GET("", pbhTypeAuthorizeRead, pbehaviorTypeApi.List)
			typeRouter.POST("", pbhTypeAuthorizeCreate, middleware.SetAuthor(errorResponder), pbehaviorTypeApi.Create)
			typeRouter.GET("/next-priority", pbhTypeAuthorizeRead, pbehaviorTypeApi.GetNextPriority)
			typeRouter.GET("/:id", pbhTypeAuthorizeRead, pbehaviorTypeApi.Get)
			typeRouter.PUT("/:id", pbhTypeAuthorizeUpdate, middleware.SetAuthor(errorResponder), pbehaviorTypeApi.Update)
			typeRouter.DELETE("/:id", pbhTypeAuthorizeDelete, pbehaviorTypeApi.Delete)
		}

		reasonRouter := protected.Group("/pbehavior-reasons")
		{
			reasonAPI := pbehaviorreason.NewApi(
				pbehaviorreason.NewStore(primaryDbClient, authorProvider),
				pbhComputeChan,
				errorResponder,
				logger,
			)
			reasonRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjPbehaviorReason, model.PermissionCreate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				reasonAPI.Create)
			reasonRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjPbehaviorReason, model.PermissionRead, enforcer, errorResponder),
				reasonAPI.List)
			reasonRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjPbehaviorReason, model.PermissionRead, enforcer, errorResponder),
				reasonAPI.Get)
			reasonRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjPbehaviorReason, model.PermissionUpdate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				reasonAPI.Update)
			reasonRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjPbehaviorReason, model.PermissionDelete, enforcer, errorResponder),
				reasonAPI.Delete)
		}
		exceptionAPI := pbehaviorexception.NewApi(
			pbehaviorexception.NewStore(primaryDbClient, timezoneConfigProvider, authorProvider),
			pbhComputeChan,
			conf.File.ImportMaxSize,
			errorResponder,
			logger,
		)
		exceptionRouter := protected.Group("/pbehavior-exceptions")
		{
			exceptionRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjPbehaviorException, model.PermissionCreate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				exceptionAPI.Create)
			exceptionRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjPbehaviorException, model.PermissionRead, enforcer, errorResponder),
				exceptionAPI.List)
			exceptionRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjPbehaviorException, model.PermissionRead, enforcer, errorResponder),
				exceptionAPI.Get)
			exceptionRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjPbehaviorException, model.PermissionUpdate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				exceptionAPI.Update)
			exceptionRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjPbehaviorException, model.PermissionDelete, enforcer, errorResponder),
				exceptionAPI.Delete)
		}

		protected.POST(
			"/pbehavior-exception-import",
			middleware.Authorize(apisecurity.ObjPbehaviorException, model.PermissionCreate, enforcer, errorResponder),
			exceptionAPI.Import,
		)

		weatherRouter := protected.Group("/weather-services")
		{
			weatherAPI := serviceweather.NewApi(serviceweather.NewStore(
				primaryDbClient,
				linkGenerator,
				alarmStore,
				timezoneConfigProvider,
				authorProvider,
				logger,
			), errorResponder)
			weatherRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjEntityService, model.PermissionRead, enforcer, errorResponder),
				weatherAPI.List,
			)
			weatherRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjEntityService, model.PermissionRead, enforcer, errorResponder),
				weatherAPI.EntityList,
			)
		}

		entityCategoryRouter := protected.Group("/entity-categories")
		{
			entityCategoryAPI := entitycategory.NewApi(
				entitycategory.NewStore(primaryDbClient, authorProvider), errorResponder)
			entityCategoryRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjEntityCategory, model.PermissionCreate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				entityCategoryAPI.Create,
			)
			entityCategoryRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjEntityCategory, model.PermissionRead, enforcer, errorResponder),
				entityCategoryAPI.List,
			)
			entityCategoryRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjEntityCategory, model.PermissionRead, enforcer, errorResponder),
				entityCategoryAPI.Get,
			)
			entityCategoryRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjEntityCategory, model.PermissionUpdate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				entityCategoryAPI.Update,
			)
			entityCategoryRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjEntityCategory, model.PermissionDelete, enforcer, errorResponder),
				entityCategoryAPI.Delete,
			)
		}

		eventApi := event.NewApi(amqpPublisher, errorResponder, logger)
		eventRouter := protected.Group("/event")
		{
			eventRouter.POST(
				"",
				middleware.Authorize(apisecurity.PermEvent, model.PermissionCan, enforcer, errorResponder),
				eventApi.Send)
		}

		appInfoApi := appinfo.NewApi(
			appinfo.NewStore(primaryDbClient, maintenanceAdapter, pgPoolProvider, security.GetConfig(), authorProvider),
			errorResponder)
		protected.GET("app-info", appInfoApi.GetAppInfo)
		appInfoRouter := protected.Group("/internal")
		{
			appInfoRouter.PUT(
				"user_interface",
				middleware.Authorize(apisecurity.PermUserInterfaceUpdate, model.PermissionCan, enforcer, errorResponder),
				appInfoApi.UpdateUserInterface,
			)
			appInfoRouter.POST(
				"user_interface",
				middleware.Authorize(apisecurity.PermUserInterfaceUpdate, model.PermissionCan, enforcer, errorResponder),
				appInfoApi.UpdateUserInterface,
			)
		}

		viewAPI := view.NewApi(
			view.NewStore(
				primaryDbClient,
				viewtab.NewStore(primaryDbClient, widget.NewStore(primaryDbClient, authorProvider, enforcer,
					patternfields.NewTransformer(primaryDbClient), validator.NewValidator(tplExecutor),
					templateConfigProvider), authorProvider, enforcer),
				authorProvider,
				enforcer,
			),
			enforcer,
			errorResponder,
		)
		viewRouter := protected.Group("/views")
		{
			viewRouter.POST(
				"",
				middleware.AuthorizeAtLeastOnePerm([]apisecurity.PermCheck{
					{
						Obj: apisecurity.ObjView,
						Act: model.PermissionCreate,
					},
					{
						Obj: apisecurity.PermPrivateViewGroups,
						Act: model.PermissionCan,
					},
				}, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				viewAPI.Create,
				middleware.ReloadEnforcerPolicyOnChange(enforcer, errorResponder),
			)
			viewRouter.GET(
				"/:id",
				middleware.AuthorizeAtLeastOnePerm([]apisecurity.PermCheck{
					{
						Obj: apisecurity.ObjView,
						Act: model.PermissionRead,
					},
					{
						Obj: apisecurity.PermPrivateViewGroups,
						Act: model.PermissionCan,
					},
				}, enforcer, errorResponder),
				middleware.AuthorizeOwnership(apisecurity.NewViewOwnerStrategy(primaryDbClient, enforcer, model.PermissionRead), errorResponder),
				viewAPI.Get,
			)
			viewRouter.PUT(
				"/:id",
				middleware.AuthorizeAtLeastOnePerm([]apisecurity.PermCheck{
					{
						Obj: apisecurity.ObjView,
						Act: model.PermissionUpdate,
					},
					{
						Obj: apisecurity.PermPrivateViewGroups,
						Act: model.PermissionCan,
					},
				}, enforcer, errorResponder),
				middleware.AuthorizeOwnership(apisecurity.NewViewOwnerStrategy(primaryDbClient, enforcer, model.PermissionUpdate), errorResponder),
				middleware.SetAuthor(errorResponder),
				viewAPI.Update,
			)
			viewRouter.DELETE(
				"/:id",
				middleware.AuthorizeAtLeastOnePerm([]apisecurity.PermCheck{
					{
						Obj: apisecurity.ObjView,
						Act: model.PermissionDelete,
					},
					{
						Obj: apisecurity.PermPrivateViewGroups,
						Act: model.PermissionCan,
					},
				}, enforcer, errorResponder),
				middleware.AuthorizeOwnership(apisecurity.NewViewOwnerStrategy(primaryDbClient, enforcer, model.PermissionDelete), errorResponder),
				viewAPI.Delete,
				middleware.ReloadEnforcerPolicyOnChange(enforcer, errorResponder),
			)
		}

		viewTabAPI := viewtab.NewApi(viewtab.NewStore(primaryDbClient, widget.NewStore(primaryDbClient, authorProvider,
			enforcer, patternfields.NewTransformer(primaryDbClient), validator.NewValidator(tplExecutor),
			templateConfigProvider), authorProvider, enforcer), enforcer, errorResponder)
		viewTabRouter := protected.Group("/view-tabs")
		{
			viewTabRouter.POST(
				"",
				middleware.AuthorizeAtLeastOnePerm([]apisecurity.PermCheck{
					{
						Obj: apisecurity.ObjView,
						Act: model.PermissionUpdate,
					},
					{
						Obj: apisecurity.PermPrivateViewGroups,
						Act: model.PermissionCan,
					},
				}, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				viewTabAPI.Create,
			)
			viewTabRouter.GET(
				"/:id",
				middleware.AuthorizeAtLeastOnePerm([]apisecurity.PermCheck{
					{
						Obj: apisecurity.ObjView,
						Act: model.PermissionRead,
					},
					{
						Obj: apisecurity.PermPrivateViewGroups,
						Act: model.PermissionCan,
					},
				}, enforcer, errorResponder),
				middleware.AuthorizeOwnership(apisecurity.NewViewTabOwnershipStrategy(primaryDbClient, enforcer, model.PermissionRead), errorResponder),
				viewTabAPI.Get,
			)
			viewTabRouter.PUT(
				"/:id",
				middleware.AuthorizeAtLeastOnePerm([]apisecurity.PermCheck{
					{
						Obj: apisecurity.ObjView,
						Act: model.PermissionUpdate,
					},
					{
						Obj: apisecurity.PermPrivateViewGroups,
						Act: model.PermissionCan,
					},
				}, enforcer, errorResponder),
				middleware.AuthorizeOwnership(apisecurity.NewViewTabOwnershipStrategy(primaryDbClient, enforcer, model.PermissionUpdate), errorResponder),
				middleware.SetAuthor(errorResponder),
				viewTabAPI.Update,
			)
			viewTabRouter.DELETE(
				"/:id",
				middleware.AuthorizeAtLeastOnePerm([]apisecurity.PermCheck{
					{
						Obj: apisecurity.ObjView,
						Act: model.PermissionUpdate,
					},
					{
						Obj: apisecurity.PermPrivateViewGroups,
						Act: model.PermissionCan,
					},
				}, enforcer, errorResponder),
				middleware.AuthorizeOwnership(apisecurity.NewViewTabOwnershipStrategy(primaryDbClient, enforcer, model.PermissionUpdate), errorResponder),
				viewTabAPI.Delete,
			)
		}

		widgetAPI := widget.NewApi(
			widget.NewStore(primaryDbClient, authorProvider, enforcer, patternfields.NewTransformer(primaryDbClient),
				validator.NewValidator(tplExecutor), templateConfigProvider),
			enforcer,
			errorResponder,
		)
		widgetRouter := protected.Group("/widgets")
		{
			widgetRouter.POST(
				"",
				middleware.AuthorizeAtLeastOnePerm([]apisecurity.PermCheck{
					{
						Obj: apisecurity.ObjView,
						Act: model.PermissionUpdate,
					},
					{
						Obj: apisecurity.PermPrivateViewGroups,
						Act: model.PermissionCan,
					},
				}, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				widgetAPI.Create,
			)
			widgetRouter.GET(
				"/:id",
				middleware.AuthorizeAtLeastOnePerm([]apisecurity.PermCheck{
					{
						Obj: apisecurity.ObjView,
						Act: model.PermissionRead,
					},
					{
						Obj: apisecurity.PermPrivateViewGroups,
						Act: model.PermissionCan,
					},
				}, enforcer, errorResponder),
				middleware.AuthorizeOwnership(apisecurity.NewWidgetOwnershipStrategy(primaryDbClient, enforcer, model.PermissionRead), errorResponder),
				widgetAPI.Get,
			)
			widgetRouter.PUT(
				"/:id",
				middleware.AuthorizeAtLeastOnePerm([]apisecurity.PermCheck{
					{
						Obj: apisecurity.ObjView,
						Act: model.PermissionUpdate,
					},
					{
						Obj: apisecurity.PermPrivateViewGroups,
						Act: model.PermissionCan,
					},
				}, enforcer, errorResponder),
				middleware.AuthorizeOwnership(apisecurity.NewWidgetOwnershipStrategy(primaryDbClient, enforcer, model.PermissionUpdate), errorResponder),
				middleware.SetAuthor(errorResponder),
				widgetAPI.Update,
			)
			widgetRouter.DELETE(
				"/:id",
				middleware.AuthorizeAtLeastOnePerm([]apisecurity.PermCheck{
					{
						Obj: apisecurity.ObjView,
						Act: model.PermissionUpdate,
					},
					{
						Obj: apisecurity.PermPrivateViewGroups,
						Act: model.PermissionCan,
					},
				}, enforcer, errorResponder),
				middleware.AuthorizeOwnership(apisecurity.NewWidgetOwnershipStrategy(primaryDbClient, enforcer, model.PermissionUpdate), errorResponder),
				widgetAPI.Delete,
			)
		}

		widgetFilterAPI := widgetfilter.NewApi(
			widgetfilter.NewStore(primaryDbClient, authorProvider, patternfields.NewTransformer(primaryDbClient)),
			enforcer, errorResponder)
		widgetFilterRouter := protected.Group("/widget-filters")
		{
			widgetFilterRouter.GET(
				"",
				middleware.AuthorizeAtLeastOnePerm([]apisecurity.PermCheck{
					{
						Obj: apisecurity.ObjView,
						Act: model.PermissionRead,
					},
					{
						Obj: apisecurity.PermPrivateViewGroups,
						Act: model.PermissionCan,
					},
				}, enforcer, errorResponder),
				widgetFilterAPI.List,
			)
			widgetFilterRouter.POST(
				"",
				middleware.AuthorizeAtLeastOnePerm([]apisecurity.PermCheck{
					{
						Obj: apisecurity.ObjView,
						Act: model.PermissionRead,
					},
					{
						Obj: apisecurity.PermPrivateViewGroups,
						Act: model.PermissionCan,
					},
				}, enforcer, errorResponder), // keep PermissionRead for private filters
				middleware.SetAuthor(errorResponder),
				widgetFilterAPI.Create,
			)
			widgetFilterRouter.GET(
				"/:id",
				middleware.AuthorizeAtLeastOnePerm([]apisecurity.PermCheck{
					{
						Obj: apisecurity.ObjView,
						Act: model.PermissionRead,
					},
					{
						Obj: apisecurity.PermPrivateViewGroups,
						Act: model.PermissionCan,
					},
				}, enforcer, errorResponder),
				widgetFilterAPI.Get,
			)
			widgetFilterRouter.PUT(
				"/:id",
				middleware.AuthorizeAtLeastOnePerm([]apisecurity.PermCheck{
					{
						Obj: apisecurity.ObjView,
						Act: model.PermissionRead,
					},
					{
						Obj: apisecurity.PermPrivateViewGroups,
						Act: model.PermissionCan,
					},
				}, enforcer, errorResponder), // keep PermissionRead for private filters
				middleware.SetAuthor(errorResponder),
				widgetFilterAPI.Update,
			)
			widgetFilterRouter.DELETE(
				"/:id",
				middleware.AuthorizeAtLeastOnePerm([]apisecurity.PermCheck{
					{
						Obj: apisecurity.ObjView,
						Act: model.PermissionRead,
					},
					{
						Obj: apisecurity.PermPrivateViewGroups,
						Act: model.PermissionCan,
					},
				}, enforcer, errorResponder),
				widgetFilterAPI.Delete,
			)
		}

		protected.PUT(
			"/widget-filter-positions",
			middleware.AuthorizeAtLeastOnePerm([]apisecurity.PermCheck{
				{
					Obj: apisecurity.ObjView,
					Act: model.PermissionRead,
				},
				{
					Obj: apisecurity.PermPrivateViewGroups,
					Act: model.PermissionCan,
				},
			}, enforcer, errorResponder),
			widgetFilterAPI.UpdatePositions,
		)

		widgetTemplateAPI := widgettemplate.NewApi(
			widgettemplate.NewStore(primaryDbClient, authorProvider), errorResponder)
		widgetTemplateRouter := protected.Group("/widget-templates")
		{
			widgetTemplateRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjWidgetTemplate, model.PermissionRead, enforcer, errorResponder),
				widgetTemplateAPI.List,
			)
			widgetTemplateRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjWidgetTemplate, model.PermissionCreate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				widgetTemplateAPI.Create,
			)
			widgetTemplateRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjWidgetTemplate, model.PermissionRead, enforcer, errorResponder),
				widgetTemplateAPI.Get,
			)
			widgetTemplateRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjWidgetTemplate, model.PermissionUpdate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				widgetTemplateAPI.Update,
			)
			widgetTemplateRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjWidgetTemplate, model.PermissionDelete, enforcer, errorResponder),
				widgetTemplateAPI.Delete,
			)
		}
		protected.POST(
			"widget-template-validate",
			middleware.Authorize(apisecurity.ObjView, model.PermissionRead, enforcer, errorResponder),
			widgetAPI.ValidateTemplates)
		protected.GET(
			"widget-template-vars",
			middleware.Authorize(apisecurity.ObjView, model.PermissionRead, enforcer, errorResponder),
			widgetAPI.GetTemplateVars)

		viewGroupAPI := viewgroup.NewApi(viewgroup.NewStore(primaryDbClient, authorProvider), errorResponder)
		viewGroupRouter := protected.Group("/view-groups")
		{
			viewGroupRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjViewGroup, model.PermissionCreate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				viewGroupAPI.Create,
			)
			viewGroupRouter.GET(
				"",
				middleware.ProvideAuthorizedIds(model.PermissionRead, enforcer, apisecurity.NewViewOwnedObjectsProvider(primaryDbClient), errorResponder),
				viewGroupAPI.List,
			)
			viewGroupRouter.GET(
				"/:id",
				middleware.AuthorizeAtLeastOnePerm([]apisecurity.PermCheck{
					{
						Obj: apisecurity.ObjViewGroup,
						Act: model.PermissionRead,
					},
					{
						Obj: apisecurity.PermPrivateViewGroups,
						Act: model.PermissionCan,
					},
				}, enforcer, errorResponder),
				viewGroupAPI.Get,
			)
			viewGroupRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjViewGroup, model.PermissionUpdate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				viewGroupAPI.Update,
			)
			viewGroupRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjViewGroup, model.PermissionDelete, enforcer, errorResponder),
				viewGroupAPI.Delete,
			)
		}

		protected.POST(
			"/view-copy/:id",
			middleware.AuthorizeAtLeastOnePerm([]apisecurity.PermCheck{
				{
					Obj: apisecurity.ObjView,
					Act: model.PermissionUpdate,
				},
				{
					Obj: apisecurity.PermPrivateViewGroups,
					Act: model.PermissionCan,
				},
			}, enforcer, errorResponder),
			middleware.AuthorizeOwnership(apisecurity.NewViewOwnerStrategy(primaryDbClient, enforcer, model.PermissionRead), errorResponder),
			middleware.SetAuthor(errorResponder),
			viewAPI.Copy,
			middleware.ReloadEnforcerPolicyOnChange(enforcer, errorResponder),
		)

		protected.PUT(
			"/view-positions",
			middleware.Authorize(apisecurity.ObjView, model.PermissionUpdate, enforcer, errorResponder),
			middleware.Authorize(apisecurity.ObjViewGroup, model.PermissionUpdate, enforcer, errorResponder),
			viewAPI.UpdatePositions,
		)

		protected.POST(
			"/view-export",
			middleware.Authorize(apisecurity.ObjView, model.PermissionRead, enforcer, errorResponder),
			middleware.Authorize(apisecurity.ObjViewGroup, model.PermissionRead, enforcer, errorResponder),
			viewAPI.Export,
		)

		protected.POST(
			"/view-import",
			middleware.Authorize(apisecurity.ObjView, model.PermissionUpdate, enforcer, errorResponder),
			middleware.Authorize(apisecurity.ObjViewGroup, model.PermissionUpdate, enforcer, errorResponder),
			viewAPI.Import,
			middleware.ReloadEnforcerPolicyOnChange(enforcer, errorResponder),
		)

		protected.POST(
			"/view-tab-copy/:id",
			middleware.AuthorizeAtLeastOnePerm([]apisecurity.PermCheck{
				{
					Obj: apisecurity.ObjView,
					Act: model.PermissionUpdate,
				},
				{
					Obj: apisecurity.PermPrivateViewGroups,
					Act: model.PermissionCan,
				},
			}, enforcer, errorResponder),
			middleware.AuthorizeOwnership(apisecurity.NewViewTabOwnershipStrategy(primaryDbClient, enforcer, model.PermissionRead), errorResponder),
			middleware.SetAuthor(errorResponder),
			viewTabAPI.Copy,
		)

		protected.PUT(
			"/view-tab-positions",
			middleware.AuthorizeAtLeastOnePerm([]apisecurity.PermCheck{
				{
					Obj: apisecurity.ObjView,
					Act: model.PermissionUpdate,
				},
				{
					Obj: apisecurity.PermPrivateViewGroups,
					Act: model.PermissionCan,
				},
			}, enforcer, errorResponder),
			viewTabAPI.UpdatePositions,
		)

		protected.POST(
			"/widget-copy/:id",
			middleware.AuthorizeAtLeastOnePerm([]apisecurity.PermCheck{
				{
					Obj: apisecurity.ObjView,
					Act: model.PermissionUpdate,
				},
				{
					Obj: apisecurity.PermPrivateViewGroups,
					Act: model.PermissionCan,
				},
			}, enforcer, errorResponder),
			middleware.AuthorizeOwnership(apisecurity.NewWidgetOwnershipStrategy(primaryDbClient, enforcer, model.PermissionRead), errorResponder),
			middleware.SetAuthor(errorResponder),
			widgetAPI.Copy,
		)

		protected.PUT(
			"/widget-grid-positions",
			middleware.AuthorizeAtLeastOnePerm([]apisecurity.PermCheck{
				{
					Obj: apisecurity.ObjView,
					Act: model.PermissionUpdate,
				},
				{
					Obj: apisecurity.PermPrivateViewGroups,
					Act: model.PermissionCan,
				},
			}, enforcer, errorResponder),
			widgetAPI.UpdateGridPositions,
		)

		// broadcast message API
		broadcastMessageApi := broadcastmessage.NewAPI(
			broadcastmessage.NewStore(primaryDbClient, maintenanceAdapter, authorProvider),
			broadcastMessageChan,
			websocketHub,
			errorResponder,
		)
		broadcastMessageRouter := protected.Group("/broadcast-message")
		{

			broadcastMessageRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjBroadcastMessage, model.PermissionCreate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				broadcastMessageApi.Create)
			broadcastMessageRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjBroadcastMessage, model.PermissionRead, enforcer, errorResponder),
				broadcastMessageApi.Get)
			broadcastMessageRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjBroadcastMessage, model.PermissionDelete, enforcer, errorResponder),
				broadcastMessageApi.Delete)
			broadcastMessageRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjBroadcastMessage, model.PermissionRead, enforcer, errorResponder),
				broadcastMessageApi.List)
			broadcastMessageRouter.PUT(
				"/:id/read",
				middleware.OnlyAuth(errorResponder),
				broadcastMessageApi.Read)
			broadcastMessageRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjBroadcastMessage, model.PermissionUpdate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				broadcastMessageApi.Update)
			// can not make typical format like /api/v4/broadcast-message/active
			// because it would be failed with conflict error apart of get /:id route
			protected.GET(
				"/active-broadcast-message",
				broadcastMessageApi.GetActive)
		}

		associativeTableApi := associativetable.NewApi(associativetable.NewStore(primaryDbClient), errorResponder)
		associativeRouter := protected.Group("/associativetable")
		{
			associativeRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjAssociativeTable, model.PermissionUpdate, enforcer, errorResponder),
				associativeTableApi.Update,
			)
			associativeRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjAssociativeTable, model.PermissionRead, enforcer, errorResponder),
				associativeTableApi.Get,
			)
			associativeRouter.DELETE(
				"",
				middleware.Authorize(apisecurity.ObjAssociativeTable, model.PermissionDelete, enforcer, errorResponder),
				associativeTableApi.Delete,
			)
		}

		scenarioAPI := scenario.NewApi(
			scenario.NewStore(primaryDbClient, authorProvider, patternfields.NewTransformer(primaryDbClient),
				validator.NewValidator(tplExecutor), tplExecutor, templateConfigProvider, json.NewEncoder(), json.NewDecoder()),
			dbexport.NewExporter(primaryDbClient),
			errorResponder,
		)
		scenarioRouter := protected.Group("/scenarios")
		{
			scenarioRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjAction, model.PermissionCreate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				scenarioAPI.Create,
			)
			scenarioRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjAction, model.PermissionRead, enforcer, errorResponder),
				scenarioAPI.List,
			)
			scenarioRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjAction, model.PermissionRead, enforcer, errorResponder),
				scenarioAPI.Get,
			)
			scenarioRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjAction, model.PermissionUpdate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				scenarioAPI.Update,
			)
			scenarioRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjAction, model.PermissionDelete, enforcer, errorResponder),
				scenarioAPI.Delete,
			)
		}
		protected.POST(
			"scenarios-db-export",
			middleware.Authorize(apisecurity.ObjAction, model.PermissionRead, enforcer, errorResponder),
			scenarioAPI.DBExport)
		protected.POST(
			"/scenario-template-validate",
			middleware.Authorize(apisecurity.ObjAction, model.PermissionRead, enforcer, errorResponder),
			scenarioAPI.ValidateTemplates)
		protected.GET(
			"/scenario-template-vars",
			middleware.Authorize(apisecurity.ObjAction, model.PermissionRead, enforcer, errorResponder),
			scenarioAPI.GetTemplateVars)
		protected.GET(
			"/scenario-pattern-fields",
			middleware.Authorize(apisecurity.ObjAction, model.PermissionRead, enforcer, errorResponder),
			patternAPI.GetPatternFields(mongo.ScenarioCollection))

		contextGraphAPI := contextgraph.NewApi(conf, contextgraph.NewMongoStatusReporter(primaryDbClient),
			workers.NewJobPublisher(jobKeyImport, amqpPublisher), conf.File.ImportMaxSize, errorResponder, logger)
		protected.PUT(
			"contextgraph-import",
			middleware.Authorize(apisecurity.ObjContextGraph, model.PermissionCreate, enforcer, errorResponder),
			contextGraphAPI.ImportAll,
		)
		protected.PUT(
			"contextgraph-import-partial",
			middleware.Authorize(apisecurity.ObjContextGraph, model.PermissionCreate, enforcer, errorResponder),
			contextGraphAPI.ImportPartial,
		)
		protected.GET(
			"contextgraph-import-status/:id",
			middleware.Authorize(apisecurity.ObjContextGraph, model.PermissionRead, enforcer, errorResponder),
			contextGraphAPI.Status,
		)

		stateSettingsRouter := protected.Group("/state-settings")
		{
			stateSettingsApi := statesettings.NewApi(
				statesettings.NewStore(primaryDbClient, stateSettingsUpdatesChan, authorProvider, patternfields.NewTransformer(primaryDbClient)),
				errorResponder)
			stateSettingsRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjStateSettings, model.PermissionCreate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				stateSettingsApi.Create,
			)
			stateSettingsRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjStateSettings, model.PermissionRead, enforcer, errorResponder),
				stateSettingsApi.List,
			)
			stateSettingsRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjStateSettings, model.PermissionRead, enforcer, errorResponder),
				stateSettingsApi.Get,
			)
			stateSettingsRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjStateSettings, model.PermissionUpdate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				stateSettingsApi.Update,
			)
			stateSettingsRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjStateSettings, model.PermissionDelete, enforcer, errorResponder),
				stateSettingsApi.Delete,
			)
		}

		notificationAPI := notification.NewApi(notification.NewStore(primaryDbClient, authorProvider), errorResponder)
		protected.GET("/notifications", notificationAPI.List)
		notifSettingsRouter := protected.Group("/notification-settings")
		{
			notifSettingsRouter.PUT(
				"",
				middleware.Authorize(apisecurity.PermNotification, model.PermissionCan, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				notificationAPI.UpdateSettings,
			)
			notifSettingsRouter.GET(
				"",
				middleware.Authorize(apisecurity.PermNotification, model.PermissionCan, enforcer, errorResponder),
				notificationAPI.GetSettings,
			)
		}

		playlistRouter := protected.Group("/playlists")
		{
			playlistApi := playlist.NewApi(playlist.NewStore(primaryDbClient, authorProvider),
				viewtab.NewStore(primaryDbClient, widget.NewStore(primaryDbClient, authorProvider, enforcer,
					patternfields.NewTransformer(primaryDbClient), validator.NewValidator(tplExecutor),
					templateConfigProvider), authorProvider, enforcer), enforcer, errorResponder)
			playlistRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjPlaylist, model.PermissionCreate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				playlistApi.Create,
				middleware.ReloadEnforcerPolicyOnChange(enforcer, errorResponder),
			)
			playlistRouter.GET(
				"",
				middleware.ProvideAuthorizedIds(model.PermissionRead, enforcer, nil, errorResponder),
				playlistApi.List,
			)
			playlistRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjPlaylist, model.PermissionRead, enforcer, errorResponder),
				middleware.AuthorizeByID(model.PermissionRead, enforcer, errorResponder),
				playlistApi.Get,
			)
			playlistRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjPlaylist, model.PermissionUpdate, enforcer, errorResponder),
				middleware.AuthorizeByID(model.PermissionUpdate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				playlistApi.Update,
			)
			playlistRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjPlaylist, model.PermissionDelete, enforcer, errorResponder),
				middleware.AuthorizeByID(model.PermissionDelete, enforcer, errorResponder),
				playlistApi.Delete,
				middleware.ReloadEnforcerPolicyOnChange(enforcer, errorResponder),
			)
		}

		idleRuleStore := idlerule.NewStore(primaryDbClient, authorProvider, patternfields.NewTransformer(primaryDbClient))
		idleRuleAPI := idlerule.NewApi(idleRuleStore, dbexport.NewExporter(primaryDbClient), errorResponder)
		idleRuleRouter := protected.Group("/idle-rules")
		{
			idleRuleRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjIdleRule, model.PermissionCreate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				idleRuleAPI.Create,
			)
			idleRuleRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjIdleRule, model.PermissionRead, enforcer, errorResponder),
				idleRuleAPI.List,
			)
			idleRuleRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjIdleRule, model.PermissionRead, enforcer, errorResponder),
				idleRuleAPI.Get,
			)
			idleRuleRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjIdleRule, model.PermissionUpdate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				idleRuleAPI.Update,
			)
			idleRuleRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjIdleRule, model.PermissionDelete, enforcer, errorResponder),
				idleRuleAPI.Delete,
			)
		}
		protected.POST(
			"idle-rules-db-export",
			middleware.Authorize(apisecurity.ObjIdleRule, model.PermissionRead, enforcer, errorResponder),
			idleRuleAPI.DBExport)

		linkRuleAPI := linkrule.NewApi(
			linkrule.NewStore(primaryDbClient, authorProvider, patternfields.NewTransformer(primaryDbClient),
				validator.NewValidator(tplExecutor), tplExecutor, templateConfigProvider, externalDataContainer, enforcer),
			dbexport.NewExporter(primaryDbClient),
			errorResponder,
		)
		linkRuleRouter := protected.Group("/link-rules")
		{
			linkRuleRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjLinkRule, model.PermissionCreate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				linkRuleAPI.Create,
			)
			linkRuleRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjLinkRule, model.PermissionRead, enforcer, errorResponder),
				linkRuleAPI.List,
			)
			linkRuleRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjLinkRule, model.PermissionRead, enforcer, errorResponder),
				linkRuleAPI.Get,
			)
			linkRuleRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjLinkRule, model.PermissionUpdate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				linkRuleAPI.Update,
			)
			linkRuleRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjLinkRule, model.PermissionDelete, enforcer, errorResponder),
				linkRuleAPI.Delete,
			)
		}
		protected.POST(
			"link-rules-db-export",
			middleware.Authorize(apisecurity.ObjLinkRule, model.PermissionRead, enforcer, errorResponder),
			linkRuleAPI.DBExport)
		protected.POST(
			"/link-rule-template-validate",
			middleware.Authorize(apisecurity.ObjLinkRule, model.PermissionRead, enforcer, errorResponder),
			middleware.SetAuthor(errorResponder),
			linkRuleAPI.ValidateTemplates)
		protected.GET(
			"/link-rule-template-vars",
			middleware.Authorize(apisecurity.ObjLinkRule, model.PermissionRead, enforcer, errorResponder),
			linkRuleAPI.GetTemplateVars)

		linkCategoryRouter := protected.Group("/link-categories")
		{
			linkCategoryRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjLinkRule, model.PermissionRead, enforcer, errorResponder),
				linkRuleAPI.GetCategories,
			)
		}

		alarmTagAPI := alarmtag.NewApi(
			alarmtag.NewStore(primaryDbClient, authorProvider, patternfields.NewTransformer(primaryDbClient)),
			errorResponder,
		)
		alarmTagRouter := protected.Group("/alarm-tags")
		{
			alarmTagRouter.GET(
				"",
				middleware.Authorize(apisecurity.PermAlarmRead, model.PermissionCan, enforcer, errorResponder),
				alarmTagAPI.List,
			)
			alarmTagRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjAlarmTag, model.PermissionCreate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				alarmTagAPI.Create,
			)
			alarmTagRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjAlarmTag, model.PermissionRead, enforcer, errorResponder),
				alarmTagAPI.Get,
			)
			alarmTagRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjAlarmTag, model.PermissionUpdate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				alarmTagAPI.Update,
			)
			alarmTagRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjAlarmTag, model.PermissionDelete, enforcer, errorResponder),
				alarmTagAPI.Delete,
			)
		}
		protected.GET(
			"alarm-tag-labels",
			middleware.Authorize(apisecurity.PermAlarmRead, model.PermissionCan, enforcer, errorResponder),
			alarmTagAPI.ListLabels,
		)

		colorThemeApi := colortheme.NewApi(
			colortheme.NewStore(primaryDbClient, authorProvider, userInterfaceAdapter),
			errorResponder)
		colorThemeRouter := protected.Group("/color-themes")
		{
			colorThemeRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjColorTheme, model.PermissionCreate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				colorThemeApi.Create,
			)
			colorThemeRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjColorTheme, model.PermissionRead, enforcer, errorResponder),
				colorThemeApi.List,
			)
			colorThemeRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjColorTheme, model.PermissionRead, enforcer, errorResponder),
				colorThemeApi.Get,
			)
			colorThemeRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjColorTheme, model.PermissionUpdate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				colorThemeApi.Update,
			)
			colorThemeRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjColorTheme, model.PermissionDelete, enforcer, errorResponder),
				colorThemeApi.Delete,
			)
		}

		healthcheckRouter := protected.Group("/healthcheck")
		{
			healthcheckApi := healthcheck.NewApi(healthcheckStore, errorResponder)
			healthcheckRouter.GET(
				"",
				middleware.Authorize(apisecurity.PermHealthcheck, model.PermissionCan, enforcer, errorResponder),
				healthcheckApi.Get,
			)
			healthcheckRouter.GET(
				"/live",
				healthcheckApi.IsLive,
			)
			healthcheckRouter.GET(
				"/status",
				middleware.Authorize(apisecurity.PermHealthcheck, model.PermissionCan, enforcer, errorResponder),
				healthcheckApi.GetStatus,
			)
			healthcheckRouter.GET(
				"/engines-order",
				middleware.Authorize(apisecurity.PermHealthcheck, model.PermissionCan, enforcer, errorResponder),
				healthcheckApi.GetEnginesOrder,
			)
			healthcheckRouter.GET(
				"/parameters",
				middleware.Authorize(apisecurity.PermHealthcheck, model.PermissionCan, enforcer, errorResponder),
				healthcheckApi.GetParameters,
			)
			healthcheckRouter.PUT(
				"/parameters",
				middleware.Authorize(apisecurity.PermHealthcheck, model.PermissionCan, enforcer, errorResponder),
				healthcheckApi.UpdateParameters,
			)
		}

		externalDataStore := externaldatatable.NewStore(primaryDbClient, pgPoolProvider, dbExportClient, json.NewDecoder())
		externalDataTableAPI := externaldatatable.NewAPI(externalDataStore, exdataImportWorker,
			conf.File.ImportMaxSize, exportTaskExecutor, json.NewEncoder(), errorResponder)
		externalDataTableRouter := protected.Group("/external-data-tables")
		{
			externalDataTableRouter.POST(
				"/:table/data",
				middleware.Authorize(apisecurity.ObjExternalDataTable, model.PermissionUpdate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				externalDataTableAPI.CreateData,
			)
			externalDataTableRouter.GET(
				"/:table/data",
				middleware.Authorize(apisecurity.ObjExternalDataTable, model.PermissionRead, enforcer, errorResponder),
				externalDataTableAPI.ListData,
			)
			externalDataTableRouter.GET(
				"/:table/data/:id",
				middleware.Authorize(apisecurity.ObjExternalDataTable, model.PermissionRead, enforcer, errorResponder),
				externalDataTableAPI.GetData,
			)
			externalDataTableRouter.PUT(
				"/:table/data/:id",
				middleware.Authorize(apisecurity.ObjExternalDataTable, model.PermissionUpdate, enforcer, errorResponder),
				externalDataTableAPI.UpdateData,
			)
			externalDataTableRouter.DELETE(
				"/:table/data/:id",
				middleware.Authorize(apisecurity.ObjExternalDataTable, model.PermissionUpdate, enforcer, errorResponder),
				externalDataTableAPI.DeleteData,
			)
			externalDataTableRouter.GET(
				"/:table/schema",
				middleware.Authorize(apisecurity.ObjExternalDataTable, model.PermissionRead, enforcer, errorResponder),
				externalDataTableAPI.GetSchema,
			)
			externalDataTableRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjExternalDataTable, model.PermissionCreate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				externalDataTableAPI.Create,
			)
			externalDataTableRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjExternalDataTable, model.PermissionRead, enforcer, errorResponder),
				externalDataTableAPI.List,
			)
			externalDataTableRouter.GET(
				"/:table",
				middleware.Authorize(apisecurity.ObjExternalDataTable, model.PermissionRead, enforcer, errorResponder),
				externalDataTableAPI.Get,
			)
			externalDataTableRouter.PUT(
				"/:table",
				middleware.Authorize(apisecurity.ObjExternalDataTable, model.PermissionUpdate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				externalDataTableAPI.Update,
			)
			externalDataTableRouter.DELETE(
				"/:table",
				middleware.Authorize(apisecurity.ObjExternalDataTable, model.PermissionDelete, enforcer, errorResponder),
				externalDataTableAPI.Delete,
			)
		}

		externalDataImportRouter := protected.Group("/external-data-import")
		{
			externalDataImportRouter.POST(
				"/:table",
				middleware.Authorize(apisecurity.ObjExternalDataTable, model.PermissionUpdate, enforcer, errorResponder),
				externalDataTableAPI.Import,
			)
			externalDataImportRouter.PUT(
				"/:id/preview",
				middleware.Authorize(apisecurity.ObjExternalDataTable, model.PermissionUpdate, enforcer, errorResponder),
				externalDataTableAPI.Preview,
			)
			externalDataImportRouter.GET(
				"/:id/status",
				middleware.Authorize(apisecurity.ObjExternalDataTable, model.PermissionUpdate, enforcer, errorResponder),
				externalDataTableAPI.ImportStatus,
			)
			externalDataImportRouter.GET(
				"/:id/data",
				middleware.Authorize(apisecurity.ObjExternalDataTable, model.PermissionUpdate, enforcer, errorResponder),
				externalDataTableAPI.ImportData,
			)
			externalDataImportRouter.PUT(
				"/:id/complete",
				middleware.Authorize(apisecurity.ObjExternalDataTable, model.PermissionUpdate, enforcer, errorResponder),
				externalDataTableAPI.ImportComplete,
			)
		}

		exportTaskExecutor.RegisterType("externaldata", externalDataStore.Export)
		externalDataExportRouter := protected.Group("/external-data-export")
		{
			externalDataExportRouter.POST(
				":table",
				middleware.Authorize(apisecurity.ObjExternalDataTable, model.PermissionRead, enforcer, errorResponder),
				externalDataTableAPI.Export,
			)
			externalDataExportRouter.GET(
				"/:id/download",
				security.GetFileAuthMiddleware(),
				middleware.Authorize(apisecurity.ObjExternalDataTable, model.PermissionRead, enforcer, errorResponder),
				externalDataTableAPI.ExportDownload,
			)
			externalDataExportRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjExternalDataTable, model.PermissionRead, enforcer, errorResponder),
				externalDataTableAPI.ExportStatus,
			)
		}

		templateAPI := template.NewAPI(template.NewStore(primaryDbClient, authorProvider, enforcer,
			tplTestTypePermMapping, json.NewDecoder()), templateConfigProvider, errorResponder, logger)

		bulkRouter := protected.Group("/bulk")
		{
			patternRouter := bulkRouter.Group("/patterns")
			{
				patternRouter.DELETE(
					"",
					middleware.PreProcessBulk(apiConfigProvider, errorResponder, false),
					patternAPI.BulkDelete,
				)
			}

			scenarioRouter := bulkRouter.Group("/scenarios")
			{
				scenarioRouter.POST(
					"",
					middleware.Authorize(apisecurity.ObjAction, model.PermissionCreate, enforcer, errorResponder),
					middleware.PreProcessBulk(apiConfigProvider, errorResponder, true),
					scenarioAPI.BulkCreate,
				)
				scenarioRouter.PUT(
					"",
					middleware.Authorize(apisecurity.ObjAction, model.PermissionUpdate, enforcer, errorResponder),
					middleware.PreProcessBulk(apiConfigProvider, errorResponder, true),
					scenarioAPI.BulkUpdate,
				)
				scenarioRouter.DELETE(
					"",
					middleware.Authorize(apisecurity.ObjAction, model.PermissionDelete, enforcer, errorResponder),
					middleware.PreProcessBulk(apiConfigProvider, errorResponder, false),
					scenarioAPI.BulkDelete,
				)
			}

			idleruleRouter := bulkRouter.Group("/idle-rules")
			{
				idleruleRouter.POST(
					"",
					middleware.Authorize(apisecurity.ObjIdleRule, model.PermissionCreate, enforcer, errorResponder),
					middleware.PreProcessBulk(apiConfigProvider, errorResponder, true),
					idleRuleAPI.BulkCreate,
				)
				idleruleRouter.PUT(
					"",
					middleware.Authorize(apisecurity.ObjIdleRule, model.PermissionUpdate, enforcer, errorResponder),
					middleware.PreProcessBulk(apiConfigProvider, errorResponder, true),
					idleRuleAPI.BulkUpdate,
				)
				idleruleRouter.DELETE(
					"",
					middleware.Authorize(apisecurity.ObjIdleRule, model.PermissionDelete, enforcer, errorResponder),
					middleware.PreProcessBulk(apiConfigProvider, errorResponder, false),
					idleRuleAPI.BulkDelete,
				)
			}

			eventFilterRouter := bulkRouter.Group("/eventfilters")
			{
				eventFilterRouter.POST(
					"",
					middleware.Authorize(apisecurity.ObjEventFilterRule, model.PermissionCreate, enforcer, errorResponder),
					middleware.PreProcessBulk(apiConfigProvider, errorResponder, true),
					eventFilterApi.BulkCreate,
				)
				eventFilterRouter.PUT(
					"",
					middleware.Authorize(apisecurity.ObjEventFilterRule, model.PermissionUpdate, enforcer, errorResponder),
					middleware.PreProcessBulk(apiConfigProvider, errorResponder, true),
					eventFilterApi.BulkUpdate,
				)
				eventFilterRouter.DELETE(
					"",
					middleware.Authorize(apisecurity.ObjEventFilterRule, model.PermissionDelete, enforcer, errorResponder),
					middleware.PreProcessBulk(apiConfigProvider, errorResponder, false),
					eventFilterApi.BulkDelete,
				)
			}

			entityserviceRouter := bulkRouter.Group("/entityservices")
			{
				entityserviceRouter.POST(
					"",
					middleware.Authorize(apisecurity.ObjEntityService, model.PermissionCreate, enforcer, errorResponder),
					middleware.PreProcessBulk(apiConfigProvider, errorResponder, true),
					entityserviceAPI.BulkCreate,
				)
				entityserviceRouter.PUT(
					"",
					middleware.Authorize(apisecurity.ObjEntityService, model.PermissionUpdate, enforcer, errorResponder),
					middleware.PreProcessBulk(apiConfigProvider, errorResponder, true),
					entityserviceAPI.BulkUpdate,
				)
				entityserviceRouter.DELETE(
					"",
					middleware.Authorize(apisecurity.ObjEntityService, model.PermissionDelete, enforcer, errorResponder),
					middleware.PreProcessBulk(apiConfigProvider, errorResponder, false),
					entityserviceAPI.BulkDelete,
				)
			}

			bulkRouter.PUT(
				"/role-permissions",
				middleware.Authorize(apisecurity.PermAcl, model.PermissionUpdate, enforcer, errorResponder),
				middleware.PreProcessBulk(apiConfigProvider, errorResponder, false),
				roleApi.BulkUpdatePermissions,
				middleware.ReloadEnforcerPolicyOnChange(enforcer, errorResponder),
			)

			userRouter := bulkRouter.Group("/users")
			{
				userRouter.POST(
					"",
					middleware.Authorize(apisecurity.PermAcl, model.PermissionCreate, enforcer, errorResponder),
					middleware.PreProcessBulk(apiConfigProvider, errorResponder, true),
					userApi.BulkCreate,
					middleware.ReloadEnforcerPolicyOnChange(enforcer, errorResponder),
				)
				userRouter.PUT(
					"",
					middleware.Authorize(apisecurity.PermAcl, model.PermissionUpdate, enforcer, errorResponder),
					middleware.PreProcessBulk(apiConfigProvider, errorResponder, true),
					userApi.BulkUpdate,
					middleware.ReloadEnforcerPolicyOnChange(enforcer, errorResponder),
				)
				userRouter.PATCH(
					"",
					middleware.Authorize(apisecurity.PermAcl, model.PermissionUpdate, enforcer, errorResponder),
					middleware.PreProcessBulk(apiConfigProvider, errorResponder, true),
					userApi.BulkPatch,
					middleware.ReloadEnforcerPolicyOnChange(enforcer, errorResponder),
				)
				userRouter.DELETE(
					"",
					middleware.Authorize(apisecurity.PermAcl, model.PermissionDelete, enforcer, errorResponder),
					middleware.PreProcessBulk(apiConfigProvider, errorResponder, false),
					userApi.BulkDelete,
				)
			}

			pbehaviorRouter := bulkRouter.Group("/pbehaviors")
			{
				pbehaviorRouter.POST(
					"",
					middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionCreate, enforcer, errorResponder),
					middleware.PreProcessBulk(apiConfigProvider, errorResponder, true),
					pbehaviorApi.BulkCreate,
				)
				pbehaviorRouter.PUT(
					"",
					middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionUpdate, enforcer, errorResponder),
					middleware.PreProcessBulk(apiConfigProvider, errorResponder, true),
					pbehaviorApi.BulkUpdate,
				)
				pbehaviorRouter.DELETE(
					"",
					middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionDelete, enforcer, errorResponder),
					middleware.PreProcessBulk(apiConfigProvider, errorResponder, false),
					pbehaviorApi.BulkDelete,
				)
			}

			entityPbehaviorRouter := bulkRouter.Group("/entity-pbehaviors")
			{
				entityPbehaviorRouter.POST(
					"",
					middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionCreate, enforcer, errorResponder),
					middleware.PreProcessBulk(apiConfigProvider, errorResponder, true),
					pbehaviorApi.BulkEntityCreate,
				)
				entityPbehaviorRouter.DELETE(
					"",
					middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionDelete, enforcer, errorResponder),
					middleware.PreProcessBulk(apiConfigProvider, errorResponder, false),
					pbehaviorApi.BulkEntityDelete,
				)
			}

			connectorPbehaviorRouter := bulkRouter.Group("/connector-pbehaviors")
			{
				connectorPbehaviorRouter.POST(
					"",
					middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionCreate, enforcer, errorResponder),
					middleware.PreProcessBulk(apiConfigProvider, errorResponder, true),
					pbehaviorApi.BulkConnectorCreate,
				)
				connectorPbehaviorRouter.DELETE(
					"",
					middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionDelete, enforcer, errorResponder),
					middleware.PreProcessBulk(apiConfigProvider, errorResponder, true),
					pbehaviorApi.BulkConnectorDelete,
				)
				connectorPbehaviorRouter.PUT(
					"",
					middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionCreate, enforcer, errorResponder),
					middleware.Authorize(apisecurity.ObjPbehavior, model.PermissionDelete, enforcer, errorResponder),
					middleware.PreProcessBulk(apiConfigProvider, errorResponder, true),
					pbehaviorApi.BulkConnectorEdit,
				)
			}

			entityRouter := bulkRouter.Group("/entities")
			{
				entityRouter.PUT(
					"/enable",
					middleware.Authorize(apisecurity.ObjEntity, model.PermissionUpdate, enforcer, errorResponder),
					middleware.PreProcessBulk(apiConfigProvider, errorResponder, true),
					entityAPI.BulkEnable,
				)
				entityRouter.PUT(
					"/disable",
					middleware.Authorize(apisecurity.ObjEntity, model.PermissionUpdate, enforcer, errorResponder),
					middleware.PreProcessBulk(apiConfigProvider, errorResponder, true),
					entityAPI.BulkDisable,
				)
			}

			linkRuleRouter := bulkRouter.Group("/link-rules")
			{
				linkRuleRouter.DELETE(
					"",
					middleware.Authorize(apisecurity.ObjLinkRule, model.PermissionDelete, enforcer, errorResponder),
					middleware.PreProcessBulk(apiConfigProvider, errorResponder, false),
					linkRuleAPI.BulkDelete,
				)
			}

			alarmRouter := bulkRouter.Group("alarms")
			{
				alarmRouter.PUT(
					"/ack",
					middleware.Authorize(apisecurity.PermAlarmUpdate, model.PermissionCan, enforcer, errorResponder),
					alarmActionAPI.BulkAck,
				)
				alarmRouter.PUT(
					"/ackremove",
					middleware.Authorize(apisecurity.PermAlarmUpdate, model.PermissionCan, enforcer, errorResponder),
					alarmActionAPI.BulkAckRemove,
				)
				alarmRouter.PUT(
					"/snooze",
					middleware.Authorize(apisecurity.PermAlarmUpdate, model.PermissionCan, enforcer, errorResponder),
					alarmActionAPI.BulkSnooze,
				)
				alarmRouter.PUT(
					"/cancel",
					middleware.Authorize(apisecurity.PermAlarmUpdate, model.PermissionCan, enforcer, errorResponder),
					alarmActionAPI.BulkCancel,
				)
				alarmRouter.PUT(
					"/uncancel",
					middleware.Authorize(apisecurity.PermAlarmUpdate, model.PermissionCan, enforcer, errorResponder),
					alarmActionAPI.BulkUncancel,
				)
				alarmRouter.PUT(
					"/assocticket",
					middleware.Authorize(apisecurity.PermAlarmUpdate, model.PermissionCan, enforcer, errorResponder),
					alarmActionAPI.BulkAssocTicket,
				)
				alarmRouter.PUT(
					"/comment",
					middleware.Authorize(apisecurity.PermAlarmUpdate, model.PermissionCan, enforcer, errorResponder),
					alarmActionAPI.BulkComment,
				)
				alarmRouter.PUT(
					"/changestate",
					middleware.Authorize(apisecurity.PermAlarmUpdate, model.PermissionCan, enforcer, errorResponder),
					alarmActionAPI.BulkChangeState,
				)
			}

			alarmTagRouter := bulkRouter.Group("/alarm-tags")
			{
				alarmTagRouter.DELETE(
					"",
					middleware.Authorize(apisecurity.ObjAlarmTag, model.PermissionDelete, enforcer, errorResponder),
					middleware.PreProcessBulk(apiConfigProvider, errorResponder, false),
					alarmTagAPI.BulkDelete,
				)
			}

			colorThemeRouter := bulkRouter.Group("/color-themes")
			{
				colorThemeRouter.DELETE(
					"",
					middleware.Authorize(apisecurity.ObjColorTheme, model.PermissionDelete, enforcer, errorResponder),
					middleware.PreProcessBulk(apiConfigProvider, errorResponder, false),
					colorThemeApi.BulkDelete,
				)
			}

			bulkRouter.DELETE(
				"/external-data-tables/:table/data",
				middleware.Authorize(apisecurity.ObjExternalDataTable, model.PermissionUpdate, enforcer, errorResponder),
				middleware.PreProcessBulk(apiConfigProvider, errorResponder, false),
				externalDataTableAPI.BulkDeleteData,
			)
		}

		dateStorageRouter := protected.Group("data-storage")
		{
			dateStorageAPI := datastorage.NewApi(datastorage.NewStore(primaryDbClient, pgPoolProvider, logger), errorResponder)
			dateStorageRouter.GET(
				"",
				middleware.Authorize(apisecurity.PermDataStorageRead, model.PermissionCan, enforcer, errorResponder),
				dateStorageAPI.Get,
			)
			dateStorageRouter.PUT(
				"",
				middleware.Authorize(apisecurity.PermDataStorageUpdate, model.PermissionCan, enforcer, errorResponder),
				dateStorageAPI.Update,
			)
		}

		messageRateStatsRouter := protected.Group("/message-rate-stats")
		{
			messageRateStatsAPI := messageratestats.NewApi(messageratestats.NewStore(pgPoolProvider), errorResponder)
			messageRateStatsRouter.GET(
				"",
				middleware.Authorize(apisecurity.PermMessageRateStatsRead, model.PermissionCan, enforcer, errorResponder),
				messageRateStatsAPI.List,
			)
		}

		fileRouter := protected.Group("/file")
		{
			fileAPI := file.NewApi(enforcer, file.NewStore(primaryDbClient, libfile.NewStorage(
				filepath.Join(conf.File.Dir, canopsis.SubDirUpload),
				libfile.NewEtagEncoder(),
			), conf.File.UploadMaxSize), errorResponder)
			fileRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjFile, model.PermissionCreate, enforcer, errorResponder),
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
				middleware.Authorize(apisecurity.ObjFile, model.PermissionDelete, enforcer, errorResponder),
				fileAPI.Delete,
			)
		}

		iconsCacheMiddlewareGetter := middleware.NewCacheMiddlewareGetter(cacheExpiration, nil)
		iconsPath := "/icons"
		iconRouter := protected.Group(iconsPath)
		{
			iconStore := icon.NewStore(
				primaryDbClient,
				libfile.NewStorage(filepath.Join(conf.File.Dir, canopsis.SubDirIcons), libfile.NewEtagEncoder()),
			)
			iconApi := icon.NewApi(iconStore, websocketHub, conf.File.IconMaxSize, []string{mimeTypeSvg}, errorResponder)
			iconRouter.POST(
				"",
				middleware.Authorize(apisecurity.PermIcon, model.PermissionCan, enforcer, errorResponder),
				iconApi.Create,
				iconsCacheMiddlewareGetter.ClearCache(BaseUrl+iconsPath),
			)
			iconRouter.GET(
				"",
				iconsCacheMiddlewareGetter.Cache(logger),
				iconApi.List,
			)
			iconRouter.GET(
				"/:id",
				security.GetFileAuthMiddleware(),
				iconsCacheMiddlewareGetter.Cache(logger),
				iconApi.Get,
			)
			iconRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.PermIcon, model.PermissionCan, enforcer, errorResponder),
				iconApi.Delete,
				iconsCacheMiddlewareGetter.ClearCache(BaseUrl+iconsPath),
			)
			iconRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.PermIcon, model.PermissionCan, enforcer, errorResponder),
				iconApi.Update,
				iconsCacheMiddlewareGetter.ClearCache(BaseUrl+iconsPath),
			)
			iconRouter.PATCH(
				"/:id",
				middleware.Authorize(apisecurity.PermIcon, model.PermissionCan, enforcer, errorResponder),
				iconApi.Patch,
				iconsCacheMiddlewareGetter.ClearCache(BaseUrl+iconsPath),
			)
		}

		resolveRuleRouter := protected.Group("/resolve-rules")
		{
			resolveRuleAPI := resolverule.NewApi(
				resolverule.NewStore(primaryDbClient, authorProvider, patternfields.NewTransformer(primaryDbClient)),
				errorResponder,
			)
			resolveRuleRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjResolveRule, model.PermissionCreate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				resolveRuleAPI.Create,
			)
			resolveRuleRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjResolveRule, model.PermissionRead, enforcer, errorResponder),
				resolveRuleAPI.List,
			)
			resolveRuleRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjResolveRule, model.PermissionRead, enforcer, errorResponder),
				resolveRuleAPI.Get,
			)
			resolveRuleRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjResolveRule, model.PermissionUpdate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				resolveRuleAPI.Update,
			)
			resolveRuleRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjResolveRule, model.PermissionDelete, enforcer, errorResponder),
				resolveRuleAPI.Delete,
			)
		}

		flappingRuleAPI := flappingrule.NewApi(
			flappingrule.NewStore(primaryDbClient, authorProvider, patternfields.NewTransformer(primaryDbClient)),
			dbexport.NewExporter(primaryDbClient), errorResponder,
		)
		flappingRuleRouter := protected.Group("/flapping-rules")
		{
			flappingRuleRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjFlappingRule, model.PermissionCreate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				flappingRuleAPI.Create,
			)
			flappingRuleRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjFlappingRule, model.PermissionRead, enforcer, errorResponder),
				flappingRuleAPI.List,
			)
			flappingRuleRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjFlappingRule, model.PermissionRead, enforcer, errorResponder),
				flappingRuleAPI.Get,
			)
			flappingRuleRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjFlappingRule, model.PermissionUpdate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				flappingRuleAPI.Update,
			)
			flappingRuleRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjFlappingRule, model.PermissionDelete, enforcer, errorResponder),
				flappingRuleAPI.Delete,
			)
		}
		protected.POST(
			"flapping-rules-db-export",
			middleware.Authorize(apisecurity.ObjFlappingRule, model.PermissionRead, enforcer, errorResponder),
			flappingRuleAPI.DBExport)

		entityInfoDictionaryApi := entityinfodictionary.NewApi(
			entityinfodictionary.NewStore(primaryDbClient), errorResponder, logger)
		protected.GET("/entity-infos-dictionary/keys",
			middleware.Authorize(apisecurity.ObjEntity, model.PermissionRead, enforcer, errorResponder),
			entityInfoDictionaryApi.ListKeys,
		)
		protected.GET("/entity-infos-dictionary/values",
			middleware.Authorize(apisecurity.ObjEntity, model.PermissionRead, enforcer, errorResponder),
			entityInfoDictionaryApi.ListValues,
		)

		techMetricsAPI := techmetrics.NewApi(techMetricsTaskExecutor, techmetrics.NewStore(primaryDbClient),
			timezoneConfigProvider, errorResponder)
		techMetricsRouter := protected.Group("/tech-metrics-export")
		{
			techMetricsRouter.POST(
				"",
				middleware.Authorize(apisecurity.PermTechMetrics, model.PermissionCan, enforcer, errorResponder),
				techMetricsAPI.StartExport,
			)
			techMetricsRouter.GET(
				"",
				middleware.Authorize(apisecurity.PermTechMetrics, model.PermissionCan, enforcer, errorResponder),
				techMetricsAPI.GetExport,
			)
			techMetricsRouter.GET(
				"/download",
				security.GetFileAuthMiddleware(),
				middleware.Authorize(apisecurity.PermTechMetrics, model.PermissionCan, enforcer, errorResponder),
				techMetricsAPI.DownloadExport,
			)
		}

		techMetricsSettingsRouter := protected.Group("/tech-metrics-settings")
		{
			techMetricsSettingsRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjTechMetricsSettings, model.PermissionRead, enforcer, errorResponder),
				techMetricsAPI.GetSettings,
			)
			techMetricsSettingsRouter.PUT(
				"",
				middleware.Authorize(apisecurity.ObjTechMetricsSettings, model.PermissionUpdate, enforcer, errorResponder),
				techMetricsAPI.UpdateSettings,
			)
		}

		tplDataRouter := protected.Group("/template-data")
		{
			protected.GET(
				"/template-vars",
				templateAPI.GetEnvVars,
			)
			tplDataRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjTemplateData, model.PermissionCreate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				templateAPI.CreateData,
			)
			tplDataRouter.GET(
				"",
				middleware.Authorize(apisecurity.ObjTemplateData, model.PermissionRead, enforcer, errorResponder),
				templateAPI.ListData,
			)
			tplDataRouter.GET(
				"/:id",
				middleware.Authorize(apisecurity.ObjTemplateData, model.PermissionRead, enforcer, errorResponder),
				templateAPI.GetData,
			)
			tplDataRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjTemplateData, model.PermissionUpdate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				templateAPI.UpdateData,
			)
			tplDataRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjTemplateData, model.PermissionDelete, enforcer, errorResponder),
				templateAPI.DeleteData,
			)
		}

		tplTestRouter := protected.Group("/template-test")
		{
			tplTestRouter.POST(
				"",
				middleware.SetAuthor(errorResponder),
				templateAPI.CreateTest,
			)
			tplTestRouter.GET(
				"",
				templateAPI.ListTest,
			)
			tplTestRouter.GET(
				"/:id",
				templateAPI.GetTest,
			)
			tplTestRouter.PUT(
				"/:id",
				middleware.SetAuthor(errorResponder),
				templateAPI.UpdateTest,
			)
			tplTestRouter.DELETE(
				"/:id",
				templateAPI.DeleteTest,
			)
		}

		maintenanceApi := maintenance.NewApi(
			maintenance.NewStore(
				primaryDbClient,
				userprovider.NewMongoProvider(primaryDbClient, apiConfigProvider),
				security.GetTokenService(),
				sessionStore,
			),
			errorResponder,
		)
		protected.PUT(
			"/maintenance",
			middleware.Authorize(apisecurity.PermMaintenance, model.PermissionCan, enforcer, errorResponder),
			maintenanceApi.Maintenance,
		)

		entityInfosPropertyAPI := entityinfosproperty.NewApi(
			entityinfosproperty.NewStore(primaryDbClient, authorProvider), errorResponder)
		entityInfosPropertyRouter := protected.Group("/entity-infos-properties")
		{
			entityInfosPropertyRouter.POST(
				"",
				middleware.Authorize(apisecurity.ObjEntityInfoProperty, model.PermissionCreate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				entityInfosPropertyAPI.Create,
			)

			entityInfosPropertyGetPerms := []apisecurity.PermCheck{
				{
					Obj: apisecurity.ObjEntity,
					Act: model.PermissionRead,
				},
				{
					Obj: apisecurity.ObjEntityInfoProperty,
					Act: model.PermissionRead,
				},
			}

			entityInfosPropertyRouter.GET(
				"",
				middleware.AuthorizeAtLeastOnePerm(entityInfosPropertyGetPerms, enforcer, errorResponder),
				entityInfosPropertyAPI.List,
			)
			entityInfosPropertyRouter.GET(
				"/:id",
				middleware.AuthorizeAtLeastOnePerm(entityInfosPropertyGetPerms, enforcer, errorResponder),
				entityInfosPropertyAPI.Get,
			)
			entityInfosPropertyRouter.PUT(
				"/:id",
				middleware.Authorize(apisecurity.ObjEntityInfoProperty, model.PermissionUpdate, enforcer, errorResponder),
				middleware.SetAuthor(errorResponder),
				entityInfosPropertyAPI.Update,
			)
			entityInfosPropertyRouter.DELETE(
				"/:id",
				middleware.Authorize(apisecurity.ObjEntityInfoProperty, model.PermissionDelete, enforcer, errorResponder),
				entityInfosPropertyAPI.Delete,
			)
		}
		entityInfosPropertyBulkRouter := bulkRouter.Group("/entity-infos-properties")
		{
			entityInfosPropertyBulkRouter.DELETE(
				"",
				middleware.Authorize(apisecurity.ObjEntityInfoProperty, model.PermissionDelete, enforcer, errorResponder),
				middleware.PreProcessBulk(apiConfigProvider, errorResponder, false),
				entityInfosPropertyAPI.BulkDelete,
			)
		}
	}

	return nil
}
