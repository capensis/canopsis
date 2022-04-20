package alarm

import (
	"fmt"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/export"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
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
	store                  Store
	exportExecutor         export.TaskExecutor
	defaultExportFields    export.Fields
	exportSeparators       map[string]rune
	timezoneConfigProvider config.TimezoneConfigProvider
}

func NewApi(
	store Store,
	executor export.TaskExecutor,
	timezoneConfigProvider config.TimezoneConfigProvider,
) API {
	fields := []string{"_id", "v.connector", "v.connector_name", "v.component",
		"v.resource", "v.output", "v.state.val", "v.status.val"}
	defaultExportFields := make(export.Fields, len(fields))
	for i, field := range fields {
		defaultExportFields[i] = export.Field{
			Name:  field,
			Label: field,
		}
	}

	return &api{
		store:               store,
		exportExecutor:      executor,
		defaultExportFields: defaultExportFields,
		exportSeparators: map[string]rune{"comma": ',', "semicolon": ';',
			"tab": '	', "space": ' '},
		timezoneConfigProvider: timezoneConfigProvider,
	}
}

// List
// @Param request query ListRequestWithPagination true "request"
// @Success 200 {object} common.PaginatedListResponse{data=[]Alarm}
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

	aggregationResult, err := a.store.Find(c.Request.Context(), apiKey, r)
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

// Count
// @Param request query FilterRequest true "request"
// @Success 200 {object} Count
func (a *api) Count(c *gin.Context) {
	var r FilterRequest

	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	res, err := a.store.Count(c.Request.Context(), r)
	if err != nil {
		panic(err)
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

	apiKey := ""
	if key, ok := c.Get(auth.ApiKey); ok {
		apiKey = key.(string)
	}

	taskID, err := a.exportExecutor.StartExecute(c.Request.Context(), export.Task{
		Filename:     "alarms",
		ExportFields: exportFields,
		Separator:    separator,
		DataFetcher: getDataFetcher(a.store, apiKey, r, exportFields.Fields(),
			a.timezoneConfigProvider.Get().Location),
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
