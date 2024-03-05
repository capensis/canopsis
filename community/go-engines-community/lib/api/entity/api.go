package entity

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/bulk"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/export"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type API interface {
	List(c *gin.Context)
	BulkEnable(c *gin.Context)
	BulkDisable(c *gin.Context)
	StartExport(c *gin.Context)
	GetExport(c *gin.Context)
	DownloadExport(c *gin.Context)
	ArchiveDisabled(c *gin.Context)
	ArchiveUnlinked(c *gin.Context)
	CleanArchived(c *gin.Context)
	GetContextGraph(c *gin.Context)
	CheckStateSetting(c *gin.Context)
	GetStateSetting(c *gin.Context)
}

type api struct {
	store                Store
	exportExecutor       export.TaskExecutor
	defaultExportFields  export.Fields
	exportSeparators     map[string]rune
	cleanTaskChan        chan<- CleanTask
	entityChangeListener chan<- entityservice.ChangeEntityMessage
	metricMetaUpdater    metrics.MetaUpdater
	actionLogger         logger.ActionLogger
	encoder              encoding.Encoder
	logger               zerolog.Logger
}

func NewApi(
	store Store,
	exportExecutor export.TaskExecutor,
	cleanTaskChan chan<- CleanTask,
	entityChangeListener chan<- entityservice.ChangeEntityMessage,
	metricMetaUpdater metrics.MetaUpdater,
	actionLogger logger.ActionLogger,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) API {
	fields := []string{"_id", "name", "type", "enabled", "connector", "component", "services"}
	defaultExportFields := make(export.Fields, len(fields))
	for i, field := range fields {
		defaultExportFields[i] = export.Field{
			Name:  field,
			Label: field,
		}
	}

	return &api{
		store:               store,
		exportExecutor:      exportExecutor,
		defaultExportFields: defaultExportFields,
		exportSeparators: map[string]rune{"comma": ',', "semicolon": ';',
			"tab": '	', "space": ' '},
		cleanTaskChan:        cleanTaskChan,
		entityChangeListener: entityChangeListener,
		metricMetaUpdater:    metricMetaUpdater,
		actionLogger:         actionLogger,
		encoder:              encoder,
		logger:               logger,
	}
}

// List
// @Success 200 {object} common.PaginatedListResponse{data=[]Entity}
func (a *api) List(c *gin.Context) {
	var query ListRequestWithPagination
	query.Query = pagination.GetDefaultQuery()
	if err := c.ShouldBind(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, query))
		return
	}

	entities, err := a.store.Find(c, query)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}

		panic(err)
	}

	res, err := common.NewPaginatedResponse(query.Query, entities)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

// StartExport
// @Param request body ExportRequest true "request"
// @Success 200 {object} ExportResponse
func (a *api) StartExport(c *gin.Context) {
	var r ExportRequest
	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	separator := a.exportSeparators[r.Separator]
	if len(r.Fields) == 0 {
		r.Fields = a.defaultExportFields
	}

	params, err := a.encoder.Encode(r.BaseFilterRequest)
	if err != nil {
		panic(err)
	}

	task, err := a.exportExecutor.StartExecute(c, export.TaskParameters{
		Type:           "entity",
		Parameters:     string(params),
		Fields:         r.Fields,
		Separator:      separator,
		FilenamePrefix: "entities",
		UserID:         c.MustGet(auth.UserKey).(string),
	})

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, ExportResponse{
		ID:     task.ID,
		Status: task.Status,
	})
}

// GetExport
// @Success 200 {object} ExportResponse
func (a *api) GetExport(c *gin.Context) {
	id := c.Param("id")
	t, err := a.exportExecutor.Get(c, id)
	if err != nil {
		panic(err)
	}

	if t == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, ExportResponse{
		ID:     id,
		Status: t.Status,
	})
}

func (a *api) DownloadExport(c *gin.Context) {
	id := c.Param("id")
	t, err := a.exportExecutor.Get(c, id)
	if err != nil {
		panic(err)
	}

	if t == nil || t.Status != export.TaskStatusSucceeded {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.Status(http.StatusOK)
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, t.Filename))
	c.Header("Content-Type", "text/csv")
	c.ContentType()
	c.File(t.File)
}

// ArchiveDisabled
// @Param body body ArchiveDisabledRequest true "body"
func (a *api) ArchiveDisabled(c *gin.Context) {
	var r ArchiveDisabledRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	select {
	case a.cleanTaskChan <- CleanTask{
		Type:                    CleanTaskTypeArchiveDisabled,
		ArchiveWithDependencies: r.WithDependencies,
		UserID:                  c.MustGet(auth.UserKey).(string),
	}:
	default:
		a.logger.Debug().Msg("cleaning in progress, skip")
	}

	c.Status(http.StatusAccepted)
}

// ArchiveUnlinked
// @Param body body ArchiveUnlinkedRequest true "body"
func (a *api) ArchiveUnlinked(c *gin.Context) {
	var r ArchiveUnlinkedRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	select {
	case a.cleanTaskChan <- CleanTask{
		Type:          CleanTaskTypeArchiveUnlinked,
		ArchiveBefore: &r.ArchiveBefore,
		UserID:        c.MustGet(auth.UserKey).(string),
	}:
	default:
		a.logger.Debug().Msg("cleaning in progress, skip")
	}

	c.Status(http.StatusAccepted)
}

func (a *api) CleanArchived(c *gin.Context) {
	select {
	case a.cleanTaskChan <- CleanTask{
		Type:   CleanTaskTypeCleanArchived,
		UserID: c.MustGet(auth.UserKey).(string),
	}:
	default:
		a.logger.Debug().Msg("cleaning in progress, skip")
	}

	c.Status(http.StatusAccepted)
}

// BulkEnable
// @Param body body []BulkToggleRequestItem true "body"
func (a *api) BulkEnable(c *gin.Context) {
	a.toggle(c, true)
}

// BulkDisable
// @Param body body []BulkToggleRequestItem true "body"
func (a *api) BulkDisable(c *gin.Context) {
	a.toggle(c, false)
}

// GetContextGraph
// @Success 200 {object} ContextGraphResponse
func (a *api) GetContextGraph(c *gin.Context) {
	var r ContextGraphRequest
	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	res, err := a.store.GetContextGraph(c, r.ID)
	if err != nil {
		panic(err)
	}

	if res == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, res)
}

// CheckStateSetting
// @Param request body CheckStateSettingRequest true "request"
// @Success 200 {object} StateSettingResponse
func (a *api) CheckStateSetting(c *gin.Context) {
	request := CheckStateSettingRequest{}

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	response, err := a.store.CheckStateSetting(c, request)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, response)
}

// GetStateSetting
// @Success 200 {object} StateSettingResponse
func (a *api) GetStateSetting(c *gin.Context) {
	request := ContextGraphRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))

		return
	}

	response, err := a.store.GetStateSetting(c, request.ID)
	if err != nil {
		if errors.Is(err, ErrNoFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)

			return
		}

		panic(err)
	}

	c.JSON(http.StatusOK, response)
}

func (a *api) toggle(c *gin.Context, enabled bool) {
	userId := c.MustGet(auth.UserKey).(string)
	bulk.Handler(c, func(request BulkToggleRequestItem) (string, error) {
		isToggled, simplifiedEntity, err := a.store.Toggle(c, request.ID, enabled)
		if err != nil || simplifiedEntity.ID == "" {
			return "", err
		}

		if isToggled {
			msg := entityservice.ChangeEntityMessage{
				ID:         simplifiedEntity.ID,
				Name:       simplifiedEntity.Name,
				Component:  simplifiedEntity.Component,
				EntityType: simplifiedEntity.Type,
				IsToggled:  isToggled,
			}

			if !enabled && simplifiedEntity.Type == types.EntityTypeComponent {
				msg.Resources = make([]string, len(simplifiedEntity.Resources))
				copy(msg.Resources, simplifiedEntity.Resources)
			}

			a.sendChangeMessage(msg)
		}

		entry := logger.LogEntry{
			Action:    logger.ActionUpdate,
			ValueType: logger.ValueTypeEntity,
			ValueID:   simplifiedEntity.ID,
		}

		if simplifiedEntity.Type == types.EntityTypeService {
			entry.ValueType = logger.ValueTypeEntityService
		}

		err = a.actionLogger.Action(context.Background(), userId, entry)
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}

		a.metricMetaUpdater.UpdateById(c, simplifiedEntity.ID)
		if isToggled && simplifiedEntity.Type == types.EntityTypeComponent {
			a.metricMetaUpdater.UpdateById(c, simplifiedEntity.Resources...)
		}

		return simplifiedEntity.ID, nil
	}, a.logger)
}

func (a *api) sendChangeMessage(msg entityservice.ChangeEntityMessage) {
	select {
	case a.entityChangeListener <- msg:
	default:
		a.logger.Err(errors.New("channel is full")).
			Str("entity", msg.ID).
			Msg("fail to send change entity message")
	}
}
