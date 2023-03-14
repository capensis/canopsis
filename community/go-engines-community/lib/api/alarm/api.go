package alarm

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/export"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"github.com/valyala/fastjson"
)

type API interface {
	List(c *gin.Context)
	Get(c *gin.Context)
	GetOpen(c *gin.Context)
	GetDetails(c *gin.Context)
	ListByService(c *gin.Context)
	ListByComponent(c *gin.Context)
	ResolvedList(c *gin.Context)
	Count(c *gin.Context)
	StartExport(c *gin.Context)
	GetExport(c *gin.Context)
	DownloadExport(c *gin.Context)
	GetLinks(c *gin.Context)
}

type api struct {
	store                  Store
	exportExecutor         export.TaskExecutor
	defaultExportFields    export.Fields
	exportSeparators       map[string]rune
	timezoneConfigProvider config.TimezoneConfigProvider

	logger zerolog.Logger
}

func NewApi(
	store Store,
	executor export.TaskExecutor,
	timezoneConfigProvider config.TimezoneConfigProvider,
	logger zerolog.Logger,
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
		logger:                 logger,
	}
}

// List
// @Success 200 {object} common.PaginatedListResponse{data=[]Alarm}
func (a *api) List(c *gin.Context) {
	var r ListRequestWithPagination
	r.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	aggregationResult, err := a.store.Find(c, r)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		panic(err)
	}

	res, err := common.NewPaginatedResponse(r.Query, aggregationResult)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

// Get
// @Success 200 {object} Alarm
func (a *api) Get(c *gin.Context) {
	alarm, err := a.store.GetByID(c, c.Param("id"))
	if err != nil {
		panic(err)
	}

	if alarm == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, alarm)
}

// GetOpen
// @Success 200 {object} Alarm
func (a *api) GetOpen(c *gin.Context) {
	r := GetOpenRequest{}
	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}
	alarm, ok, err := a.store.GetOpenByEntityID(c, r.ID)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	if alarm == nil {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, alarm)
}

// GetDetails
// @Param request body []DetailsRequest true "request"
// @Success 200 {array} DetailsResponse
func (a *api) GetDetails(c *gin.Context) {
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

	defaultQuery := pagination.GetDefaultQuery()
	response := make([]DetailsResponse, len(rawObjects))

	for idx, rawObject := range rawObjects {
		object, err := rawObject.Object()
		if err != nil {
			response[idx].Status = http.StatusBadRequest
			response[idx].Error = err.Error()
			continue
		}

		var request DetailsRequest
		err = json.Unmarshal(object.MarshalTo(nil), &request)
		if err != nil {
			response[idx].Status = http.StatusBadRequest
			response[idx].Error = err.Error()
			continue
		}

		if request.Steps != nil {
			request.Steps.Paginate = true
			if request.Steps.Page == 0 {
				request.Steps.Page = defaultQuery.Page
			}
			if request.Steps.Limit == 0 {
				request.Steps.Limit = defaultQuery.Limit
			}
		}

		if request.Children != nil {
			request.Children.Paginate = true
			if request.Children.Page == 0 {
				request.Children.Page = defaultQuery.Page
			}
			if request.Children.Limit == 0 {
				request.Children.Limit = defaultQuery.Limit
			}
		}

		err = binding.Validator.ValidateStruct(request)
		if err != nil {
			response[idx].ID = request.ID
			response[idx].Status = http.StatusBadRequest
			var errs validator.ValidationErrors
			if errors.As(err, &errs) {
				response[idx].Errors = common.TransformValidationErrors(errs, request).Errors
			} else {
				response[idx].Error = "Request has invalid structure"
			}
			continue
		}

		details, err := a.store.GetDetails(c, request)
		if err != nil {
			response[idx].ID = request.ID
			response[idx].Status = http.StatusInternalServerError
			response[idx].Error = common.InternalServerErrorResponse.Error
			a.logger.Err(err).Msg("cannot fetch alarm details")
			continue
		}

		if details == nil {
			response[idx].ID = request.ID
			response[idx].Status = http.StatusNotFound
			response[idx].Error = common.NotFoundResponse.Error
			continue
		}

		response[idx].ID = request.ID
		response[idx].Status = http.StatusOK
		response[idx].Data = *details
	}

	c.JSON(http.StatusMultiStatus, response)
}

// ListByService
// @Success 200 {object} common.PaginatedListResponse{data=[]Alarm}
func (a *api) ListByService(c *gin.Context) {
	r := ListByServiceRequest{
		Query: pagination.GetDefaultQuery(),
	}
	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	aggregationResult, err := a.store.FindByService(c, c.Param("id"), r)
	if err != nil {
		panic(err)
	}

	if aggregationResult == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	res, err := common.NewPaginatedResponse(r.Query, aggregationResult)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

// ListByComponent
// @Success 200 {object} common.PaginatedListResponse{data=[]Alarm}
func (a *api) ListByComponent(c *gin.Context) {
	r := ListByComponentRequest{
		Query: pagination.GetDefaultQuery(),
	}
	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	aggregationResult, err := a.store.FindByComponent(c, r)
	if err != nil {
		panic(err)
	}

	if aggregationResult == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	res, err := common.NewPaginatedResponse(r.Query, aggregationResult)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

// ResolvedList
// @Success 200 {object} common.PaginatedListResponse{data=[]Alarm}
func (a *api) ResolvedList(c *gin.Context) {
	r := ResolvedListRequest{
		Query: pagination.GetDefaultQuery(),
	}
	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	aggregationResult, err := a.store.FindResolved(c, r)
	if err != nil {
		panic(err)
	}

	if aggregationResult == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	res, err := common.NewPaginatedResponse(r.Query, aggregationResult)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

// Count
// @Success 200 {object} Count
func (a *api) Count(c *gin.Context) {
	var r FilterRequest

	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	res, err := a.store.Count(c, r)
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

	taskID, err := a.exportExecutor.StartExecute(c, export.Task{
		Filename:     "alarms",
		ExportFields: exportFields,
		Separator:    separator,
		DataFetcher: getDataFetcher(a.store, r, exportFields.Fields(),
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
	t, err := a.exportExecutor.GetStatus(c, id)
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
	t, err := a.exportExecutor.GetStatus(c, id)
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

// GetLinks
// @Success 200 {array} link.Link
func (a *api) GetLinks(c *gin.Context) {
	var r LinksRequest
	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	links, ok, err := a.store.GetLinks(c, c.Param("id"), r.Ids)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, links)
}
