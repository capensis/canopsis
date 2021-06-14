package alarm

import (
	"net/http"

	"git.canopsis.net/canopsis/go-engines/lib/api/auth"
	"git.canopsis.net/canopsis/go-engines/lib/api/common"
	"git.canopsis.net/canopsis/go-engines/lib/api/export"
	"git.canopsis.net/canopsis/go-engines/lib/api/pagination"
	"github.com/gin-gonic/gin"
)

type API interface {
	List(c *gin.Context)
	Count(c *gin.Context)
	StartExport(c *gin.Context)
	GetExport(c *gin.Context)
	DownloadExport(c *gin.Context)
}

type api struct {
	store               Store
	exportExecutor      export.TaskExecutor
	defaultExportFields []string
	exportSeparators    map[string]rune
}

func NewApi(
	store Store,
	executor export.TaskExecutor,
) API {
	return &api{
		store:          store,
		exportExecutor: executor,
		defaultExportFields: []string{"_id", "v.connector", "v.connector_name", "v.component",
			"v.resource", "v.output", "v.state.val", "v.status.val"},
		exportSeparators: map[string]rune{"comma": ',', "semicolon": ';',
			"tab": '	', "space": ' '},
	}
}

// Find all alarms
// @Summary Find all alarms
// @Description Get paginated list of alarms
// @Tags alarms
// @ID alarms-find-all
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param request query ListRequestWithPagination true "request"
// @Success 200 {object} common.PaginatedListResponse{data=[]Alarm}
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /alarms [get]
func (a *api) List(c *gin.Context) {
	var r ListRequestWithPagination
	r.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	apiKey := ""
	if key, ok := c.Get(auth.ApiKey); ok {
		apiKey = key.(string)
	}

	aggregationResult, err := a.store.Find(apiKey, r)
	if err != nil {
		panic(err)
	}

	res, err := common.NewPaginatedResponse(r.Query, aggregationResult)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

// Count alarms
// @Summary Count alarms
// @Description Get counts of alarms
// @Tags alarms
// @ID alarms-get-counts
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param request query FilterRequest true "request"
// @Success 200 {object} Count
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /alarm-counters [get]
func (a *api) Count(c *gin.Context) {
	var r FilterRequest

	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	res, err := a.store.Count(r)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, res)
}

// Start export alarms
// @Summary Start export alarms
// @Description Start export alarms
// @Tags alarms
// @ID alarms-export-start
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param request query ExportRequest true "request"
// @Success 200 {object} ExportResponse
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /alarm-export [post]
func (a *api) StartExport(c *gin.Context) {
	var r ExportRequest
	if err := c.ShouldBindQuery(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	separator := a.exportSeparators[r.Separator]
	exportFields := r.SearchBy
	if len(exportFields) == 0 {
		exportFields = a.defaultExportFields
	}

	apiKey := ""
	if key, ok := c.Get(auth.ApiKey); ok {
		apiKey = key.(string)
	}

	taskID, err := a.exportExecutor.StartExecute(export.Task{
		ExportFields: exportFields,
		Separator:    separator,
		DataFetcher: func(page, limit int64) (interface{}, int64, error) {
			res, err := a.store.Find(apiKey, ListRequestWithPagination{
				Query:       pagination.Query{Page: page, Limit: limit, Paginate: true},
				ListRequest: r.ListRequest,
			})
			if err != nil {
				return nil, 0, err
			}

			return res.Data, res.TotalCount, nil
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

// Get status of export alarms
// @Summary Get status of export alarms
// @Description Get status of export alarms
// @Tags alarms
// @ID alarms-export-get
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "export task id"
// @Success 200 {object} ExportResponse
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /alarm-export/{id} [get]
func (a *api) GetExport(c *gin.Context) {
	id := c.Param("id")
	t, err := a.exportExecutor.GetStatus(id)
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

// Get result of export alarms
// @Summary Get result of export alarms
// @Description Get result of export alarms
// @Tags alarms
// @ID alarms-export-download
// @Produce text/csv
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "export task id"
// @Success 200 {object} http.Response
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /alarm-export/{id}/download [get]
func (a *api) DownloadExport(c *gin.Context) {
	id := c.Param("id")
	t, err := a.exportExecutor.GetStatus(id)
	if err != nil {
		panic(err)
	}

	if t == nil || t.Status != export.TaskStatusSucceeded {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.Status(http.StatusOK)
	c.Header("Content-Disposition", `attachment; filename="alarms.csv"`)
	c.Header("Content-Type", "text/csv")
	c.ContentType()
	c.File(t.File)
}
