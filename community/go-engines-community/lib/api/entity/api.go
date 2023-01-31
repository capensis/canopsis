package entity

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/export"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/rs/zerolog"
	"github.com/valyala/fastjson"
)

type API interface {
	List(c *gin.Context)
	BulkEnable(c *gin.Context)
	BulkDisable(c *gin.Context)
	StartExport(c *gin.Context)
	GetExport(c *gin.Context)
	DownloadExport(c *gin.Context)
	Clean(c *gin.Context)
	GetContextGraph(c *gin.Context)
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
	logger               zerolog.Logger
}

func NewApi(
	store Store,
	exportExecutor export.TaskExecutor,
	cleanTaskChan chan<- CleanTask,
	entityChangeListener chan<- entityservice.ChangeEntityMessage,
	metricMetaUpdater metrics.MetaUpdater,
	actionLogger logger.ActionLogger,
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
	exportFields := r.Fields
	if len(exportFields) == 0 {
		exportFields = a.defaultExportFields
	}

	fields := exportFields.Fields()
	taskID, err := a.exportExecutor.StartExecute(c, export.Task{
		Filename:     "entities",
		ExportFields: exportFields,
		Separator:    separator,
		DataFetcher: func(ctx context.Context, page, limit int64) ([]map[string]string, int64, error) {
			res, err := a.store.Find(ctx, ListRequestWithPagination{
				Query: pagination.Query{Paginate: true, Page: page, Limit: limit},
				ListRequest: ListRequest{
					BaseFilterRequest: r.BaseFilterRequest,
					SearchBy:          fields,
				},
			})
			if err != nil {
				return nil, 0, err
			}
			data, err := export.ConvertToMap(res.Data, fields, "", nil)
			if err != nil {
				return nil, 0, err
			}

			return data, res.TotalCount, err
		},
	})

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, ExportResponse{
		ID:     taskID,
		Status: export.TaskStatusRunning,
	})
}

// GetExport
// @Success 200 {object} ExportResponse
func (a *api) GetExport(c *gin.Context) {
	id := c.Param("id")
	t, err := a.exportExecutor.GetStatus(c, id)
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
	t, err := a.exportExecutor.GetStatus(c, id)
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

func (a *api) Clean(c *gin.Context) {
	var r CleanRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	select {
	case a.cleanTaskChan <- CleanTask{
		Archive:             r.Archive,
		ArchiveDependencies: r.ArchiveDependencies,
		UserID:              c.MustGet(auth.UserKey).(string),
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

func (a *api) toggle(c *gin.Context, enabled bool) {
	userId := c.MustGet(auth.UserKey).(string)

	var ar fastjson.Arena

	raw, err := c.GetRawData()
	if err != nil {
		panic(err)
	}

	jsonValue, err := fastjson.ParseBytes(raw)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	rawObjects, err := jsonValue.Array()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	response := ar.NewArray()

	for idx, rawObject := range rawObjects {
		userObject, err := rawObject.Object()
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawObject, ar.NewString(err.Error())))
			continue
		}

		var request BulkToggleRequestItem
		err = json.Unmarshal(userObject.MarshalTo(nil), &request)
		if err != nil {
			a.logger.Err(err).Msg("cannot update entity")
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusInternalServerError, rawObject, ar.NewString(common.InternalServerErrorResponse.Error)))
			continue
		}

		err = binding.Validator.ValidateStruct(request)
		if err != nil {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusBadRequest, rawObject, common.NewValidationErrorFastJsonValue(&ar, err, request)))
			continue
		}

		isToggled, simplifiedEntity, err := a.store.Toggle(c, request.ID, enabled)
		if err != nil {
			a.logger.Err(err).Msg("cannot update entity")
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusInternalServerError, rawObject, ar.NewString(common.InternalServerErrorResponse.Error)))
			continue
		}

		if simplifiedEntity.ID == "" {
			response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, "", http.StatusNotFound, rawObject, ar.NewString("Not found")))
			continue
		}

		if isToggled {
			msg := entityservice.ChangeEntityMessage{
				ID:         simplifiedEntity.ID,
				EntityType: simplifiedEntity.Type,
				IsToggled:  isToggled,
			}

			if !enabled && simplifiedEntity.Type == types.EntityTypeComponent {
				msg.Resources = make([]string, len(simplifiedEntity.Resources))
				copy(msg.Resources, simplifiedEntity.Resources)
			}

			a.sendChangeMessage(msg)
		}

		response.SetArrayItem(idx, common.GetBulkResponseItem(&ar, simplifiedEntity.ID, http.StatusOK, rawObject, nil))

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
	}

	c.Data(http.StatusMultiStatus, gin.MIMEJSON, response.MarshalTo(nil))
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
