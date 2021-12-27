package entity

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/export"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"net/http"
)

type API interface {
	List(c *gin.Context)
	StartExport(c *gin.Context)
	GetExport(c *gin.Context)
	DownloadExport(c *gin.Context)
	Clean(c *gin.Context)
}

type api struct {
	store               Store
	exportExecutor      export.TaskExecutor
	defaultExportFields export.Fields
	exportSeparators    map[string]rune
	cleanTaskChan       chan<- CleanTask
	logger              zerolog.Logger
}

func NewApi(
	store Store,
	exportExecutor export.TaskExecutor,
	cleanTaskChan chan<- CleanTask,
	logger zerolog.Logger,
) API {
	fields := []string{"_id", "name", "type", "enabled", "depends", "impact"}
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
		cleanTaskChan: cleanTaskChan,
		logger:        logger,
	}
}

// Find all entities
// @Summary Find entities
// @Description Get paginated list of entities
// @Tags entities
// @ID entities-find-all
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param request query ListRequestWithPagination true "request"
// @Success 200 {object} common.PaginatedListResponse{data=[]Entity}
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /entities [get]
func (a *api) List(c *gin.Context) {
	var query ListRequestWithPagination
	query.Query = pagination.GetDefaultQuery()
	if err := c.ShouldBind(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, query))
		return
	}

	entities, err := a.store.Find(c.Request.Context(), query)
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

// Start export entities
// @Summary Start export entities
// @Description Start export entities
// @Tags entities
// @ID entities-export-start
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param request body ExportRequest true "request"
// @Success 200 {object} ExportResponse
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /entity-export [post]
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
	taskID, err := a.exportExecutor.StartExecute(c.Request.Context(), export.Task{
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

// Get status of export entities
// @Summary Get status of export entities
// @Description Get status of export entities
// @Tags entities
// @ID entities-export-get
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "export task id"
// @Success 200 {object} ExportResponse
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /entity-export/{id} [get]
func (a *api) GetExport(c *gin.Context) {
	id := c.Param("id")
	t, err := a.exportExecutor.GetStatus(c.Request.Context(), id)
	if err != nil {
		panic(err)
	}

	if t == nil {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, ExportResponse{
		ID:     id,
		Status: t.Status,
	})
}

// Get result of export entities
// @Summary Get result of export entities
// @Description Get result of export entities
// @Tags entities
// @ID entities-export-download
// @Produce text/csv
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "export task id"
// @Success 200 {object} http.Response
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /entity-export/{id}/download [get]
func (a *api) DownloadExport(c *gin.Context) {
	id := c.Param("id")
	t, err := a.exportExecutor.GetStatus(c.Request.Context(), id)
	if err != nil {
		panic(err)
	}

	if t == nil || t.Status != export.TaskStatusSucceeded {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.Status(http.StatusOK)
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, t.Filename))
	c.Header("Content-Type", "text/csv")
	c.ContentType()
	c.File(t.File)
}

// Clean disabled entities
// @Summary Clean disabled entities
// @Description Clean disabled entities
// @Tags entities
// @ID entities-clean
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param request query CleanRequest true "request"
// @Success 202
// @Failure 400 {object} common.ErrorResponse
// @Router /entities/clean [post]
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
