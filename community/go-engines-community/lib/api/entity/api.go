package entity

import (
	"git.canopsis.net/canopsis/go-engines/lib/api/common"
	"git.canopsis.net/canopsis/go-engines/lib/api/export"
	"git.canopsis.net/canopsis/go-engines/lib/api/pagination"
	"github.com/gin-gonic/gin"
	"net/http"
)

type API interface {
	common.CrudAPI
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
	exportExecutor export.TaskExecutor,
) API {
	return &api{
		store:               store,
		exportExecutor:      exportExecutor,
		defaultExportFields: []string{"_id", "name", "type", "enabled", "depends", "impact"},
		exportSeparators: map[string]rune{"comma": ',', "semicolon": ';',
			"tab": '	', "space": ' '},
	}
}

func (a *api) List(c *gin.Context) {
	panic("not implemented")
}

func (a *api) Get(c *gin.Context) {
	panic("not implemented")
}

func (a *api) Create(c *gin.Context) {
	panic("not implemented")
}

func (a *api) Update(c *gin.Context) {
	panic("not implemented")
}

func (a *api) Delete(c *gin.Context) {
	panic("not implemented")
}

// Start export entities
// @Summary Start export entities
// @Description Start export entities
// @Tags entities
// @ID entities-export-start
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param request query ExportRequest true "request"
// @Success 200 {object} ExportResponse
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /entity-export [post]
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

	taskID, err := a.exportExecutor.StartExecute(export.Task{
		ExportFields: exportFields,
		Separator:    separator,
		DataFetcher: func(page, limit int64) (interface{}, int64, error) {
			res, err := a.store.Find(ListRequestWithPagination{
				Query:       pagination.Query{Paginate: true, Page: page, Limit: limit},
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
	t, err := a.exportExecutor.GetStatus(id)
	if err != nil {
		panic(err)
	}

	if t == nil || t.Status != export.TaskStatusSucceeded {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.Status(http.StatusOK)
	c.Header("Content-Disposition", `attachment; filename="entities.csv"`)
	c.Header("Content-Type", "text/csv")
	c.ContentType()
	c.File(t.File)
}
